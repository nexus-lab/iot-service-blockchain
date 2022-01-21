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

// ServiceEvent an event emitted by the service registry contract notifying a service update
type ServiceEvent struct {
	// Action name of the action performed on the service
	Action string

	// OrganizationId organization ID of the service
	OrganizationId string

	// DeviceId ID of the device to which the service belongs
	DeviceId string

	// ServiceName name of the service
	ServiceName string

	// Payload custom event payload
	Payload interface{}
}

// ServiceRegistryInterface core utilities for managing services on the ledger
type ServiceRegistryInterface interface {
	// Register create or update a service in the ledger
	Register(service *common.Service) error

	// Get return a service by its organization ID, device ID, and name
	Get(organizationId string, deviceId string, serviceName string) (*common.Service, error)

	// GetAll return a list of services by their organization ID and device ID
	GetAll(organizationId string, deviceId string) ([]*common.Service, error)

	// Deregister remove a service from the ledger
	Deregister(service *common.Service) error

	// RegisterEvent registers for service registry events
	RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *ServiceEvent, context.CancelFunc, error)
}

// ServiceRegistry core utilities for managing services on the ledger
type ServiceRegistry struct {
	contract ContractInterface
}

// Register create or update a service in the ledger
func (r *ServiceRegistry) Register(service *common.Service) error {
	if service == nil {
		return fmt.Errorf("cannot register an empty service")
	}

	data, err := service.Serialize()
	if err != nil {
		return err
	}

	_, err = r.contract.SubmitTransaction("Register", string(data))
	return err
}

// Get return a service by its organization ID, device ID, and name
func (r *ServiceRegistry) Get(organizationId string, deviceId string, serviceName string) (*common.Service, error) {
	data, err := r.contract.SubmitTransaction("Get", organizationId, deviceId, serviceName)
	if err != nil {
		return nil, err
	}

	return common.DeserializeService(data)
}

// GetAll return a list of services by their organization ID and device ID
func (r *ServiceRegistry) GetAll(organizationId string, deviceId string) ([]*common.Service, error) {
	data, err := r.contract.SubmitTransaction("GetAll", organizationId, deviceId)
	if err != nil {
		return nil, err
	}

	results := make([]*common.Service, 0)
	if err = json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Deregister remove a service from the ledger
func (r *ServiceRegistry) Deregister(service *common.Service) error {
	if service == nil {
		return fmt.Errorf("cannot deregister an empty device")
	}

	data, err := service.Serialize()
	if err != nil {
		return err
	}

	_, err = r.contract.SubmitTransaction("Deregister", string(data))
	return err
}

// RegisterEvent registers for service registry events
func (r *ServiceRegistry) RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *ServiceEvent, context.CancelFunc, error) {
	dest := make(chan *ServiceEvent)
	source, cancel, err := r.contract.RegisterEvent(options...)
	pattern := regexp.MustCompile(`^service:\/\/(.+?)\/(.+?)\/(.+?)\/(.+?)$`)

	go func() {
		defer close(dest)

		for event := range source {
			matches := pattern.FindStringSubmatch(event.EventName)
			if len(matches) != 5 {
				continue
			}

			serviceEvent := &ServiceEvent{
				OrganizationId: matches[1],
				DeviceId:       matches[2],
				ServiceName:    matches[3],
				Action:         matches[4],
			}

			if serviceEvent.Action == "register" || serviceEvent.Action == "deregister" {
				service, err := common.DeserializeService(event.Payload)
				if err != nil {
					log.Printf("bad service event payload %#v, action is %s\n", event.Payload, serviceEvent.Action)
					continue
				}
				serviceEvent.Payload = service
			} else {
				serviceEvent.Payload = event.Payload
			}

			dest <- serviceEvent
		}
	}()

	return dest, cancel, err
}

func createServiceRegistry(network *client.Network, chaincodeId string) ServiceRegistryInterface {
	return &ServiceRegistry{
		contract: &Contract{
			network:      network,
			chaincodeId:  chaincodeId,
			contractName: "service_registry",
		},
	}
}
