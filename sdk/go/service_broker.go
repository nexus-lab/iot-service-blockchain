package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// ServiceRequestEvent an event emitted by the service broker contract notifying a service request/response update
type ServiceRequestEvent struct {
	// Action name of the action performed on the service request
	Action string

	// OrganizationId organization ID of the requested service
	OrganizationId string

	// DeviceId device ID of the requested service
	DeviceId string

	// ServiceName name of the requested service
	ServiceName string

	// RequestId ID of the request
	RequestId string

	// Payload custom event payload
	Payload interface{}
}

// ServiceBrokerInterface core utilities for managing service requests on ledger
type ServiceBrokerInterface interface {
	// Request make a request to an IoT service
	Request(request *common.ServiceRequest) error

	// Respond respond to an IoT service request
	Respond(response *common.ServiceResponse) error

	// Get return an IoT service request and its response (if any) by the request ID
	Get(requestId string) (*common.ServiceRequestResponse, error)

	// GetAll return a list of IoT service requests and their responses (if any) by their service organization ID, service device ID, and service name
	GetAll(organizationId string, deviceId string, serviceName string) ([]*common.ServiceRequestResponse, error)

	// Remove remove a service request and its response (if any) from the ledger
	Remove(requestId string) error

	// RegisterEvent registers for service request events
	RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *ServiceRequestEvent, context.CancelFunc, error)
}

// ServiceBroker core utilities for managing IoT service requests and responses on the ledger
type ServiceBroker struct {
	contract ContractInterface
}

// Request make a request to an IoT service
func (r *ServiceBroker) Request(request *common.ServiceRequest) error {
	if request == nil {
		return fmt.Errorf("cannot send an empty request")
	}

	data, err := request.Serialize()
	if err != nil {
		return err
	}

	_, err = r.contract.SubmitTransaction("Request", string(data))
	return err
}

// Respond respond to an IoT service request
func (r *ServiceBroker) Respond(response *common.ServiceResponse) error {
	if response == nil {
		return fmt.Errorf("cannot send an empty response")
	}

	data, err := response.Serialize()
	if err != nil {
		return err
	}

	_, err = r.contract.SubmitTransaction("Respond", string(data))
	return err
}

// Get return an IoT service request and its response by the request ID
func (r *ServiceBroker) Get(requestId string) (*common.ServiceRequestResponse, error) {
	data, err := r.contract.SubmitTransaction("Get", requestId)
	if err != nil {
		return nil, err
	}

	return common.DeserializeServiceRequestResponse(data)
}

// GetAll return a list of IoT service requests and their responses by their organization ID, device ID, and service name
func (r *ServiceBroker) GetAll(organizationId string, deviceId string, serviceName string) ([]*common.ServiceRequestResponse, error) {
	data, err := r.contract.SubmitTransaction("GetAll", organizationId, deviceId, serviceName)
	if err != nil {
		return nil, err
	}

	results := make([]*common.ServiceRequestResponse, 0)
	if err = json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Remove remove a (request, response) pair from the ledger
func (r *ServiceBroker) Remove(requestId string) error {
	_, err := r.contract.SubmitTransaction("Remove", requestId)
	return err
}

// RegisterEvent registers for service request events
func (r *ServiceBroker) RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *ServiceRequestEvent, context.CancelFunc, error) {
	dest := make(chan *ServiceRequestEvent)
	source, cancel, err := r.contract.RegisterEvent(options...)
	pattern := regexp.MustCompile(`^request:\/\/(.+?)\/(.+?)\/(.+?)\/(.+?)\/(.+?)$`)

	go func() {
		defer close(dest)

		for event := range source {
			matches := pattern.FindStringSubmatch(event.EventName)
			if len(matches) != 6 {
				continue
			}

			serviceRequestEvent := &ServiceRequestEvent{
				OrganizationId: matches[1],
				DeviceId:       matches[2],
				ServiceName:    matches[3],
				RequestId:      matches[4],
				Action:         matches[5],
			}

			if serviceRequestEvent.Action == "request" {
				request, err := common.DeserializeServiceRequest(event.Payload)
				if err != nil {
					log.Printf("bad service request event payload %#v, action is %s\n", event.Payload, serviceRequestEvent.Action)
					continue
				}
				serviceRequestEvent.Payload = request
			} else if serviceRequestEvent.Action == "respond" {
				response, err := common.DeserializeServiceResponse(event.Payload)
				if err != nil {
					log.Printf("bad service response event payload %#v, action is %s\n", event.Payload, serviceRequestEvent.Action)
					continue
				}
				serviceRequestEvent.Payload = response
			} else if serviceRequestEvent.Action == "remove" {
				serviceRequestEvent.Payload = string(event.Payload)
			} else {
				serviceRequestEvent.Payload = event.Payload
			}

			dest <- serviceRequestEvent
		}
	}()

	return dest, cancel, err
}

// CreateServiceBroker the default factory for creating service brokers
func CreateServiceBroker(network *client.Network, chaincodeId string) ServiceBrokerInterface {
	return &ServiceBroker{
		contract: &Contract{
			network:      network,
			chaincodeId:  chaincodeId,
			contractName: "service_broker",
		},
	}
}
