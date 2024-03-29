package contract

import (
	"encoding/json"
	"fmt"

	"github.com/nexus-lab/iot-service-blockchain/common"
)

// ServiceBrokerInterface core utilities for managing service requests on ledger
type ServiceBrokerInterface interface {
	// Request make a request to an IoT service
	Request(request *common.ServiceRequest) error

	// Respond respond to an IoT service request
	Respond(response *common.ServiceResponse) error

	// Get return an IoT service request and its response by the request ID
	Get(requestId string) (*common.ServiceRequestResponse, error)

	// GetAll return a list of IoT service requests and their responses by their organization ID, device ID, and service name
	GetAll(organizationId string, deviceId string, serviceName string) ([]*common.ServiceRequestResponse, error)

	// Remove remove a (request, response) pair from the ledger
	Remove(requestId string) error
}

// Dummy index object
type serviceRequestIndex struct {
	OrganizationId string `json:"organizationId"`
	DeviceId       string `json:"deviceId"`
	ServiceName    string `json:"serviceName"`
	RequestId      string `json:"requestId"`
}

func (i *serviceRequestIndex) GetKeyComponents() []string {
	return []string{i.OrganizationId, i.DeviceId, i.ServiceName, i.RequestId}
}

func (i *serviceRequestIndex) Serialize() ([]byte, error) {
	return json.Marshal(i)
}

func (i *serviceRequestIndex) Validate() error {
	return nil
}

func deserializeServiceRequestIndex(data []byte) (*serviceRequestIndex, error) {
	index := new(serviceRequestIndex)

	if err := json.Unmarshal(data, index); err != nil {
		return nil, err
	}

	return index, nil
}

// ServiceBroker core utilities for managing IoT service requests and responses on the ledger
type ServiceBroker struct {
	ctx              TransactionContextInterface
	requestRegistry  StateRegistryInterface
	responseRegistry StateRegistryInterface
	indexRegistry    StateRegistryInterface
}

// Request make a request to an IoT service
func (b *ServiceBroker) Request(request *common.ServiceRequest) error {
	// check if service exists
	service := request.Service
	_, err := b.ctx.GetServiceRegistry().Get(service.OrganizationId, service.DeviceId, service.Name)
	if err != nil {
		return err
	}

	// check if request already exists
	request_, err := b.getRequest(request.Id)
	if _, ok := err.(*common.NotFoundError); err != nil && !ok {
		return err
	}
	if request_ != nil {
		return fmt.Errorf("request already exists")
	}

	if err = b.requestRegistry.PutState(request); err != nil {
		return err
	}

	return b.indexRegistry.PutState(
		&serviceRequestIndex{
			OrganizationId: request.Service.OrganizationId,
			DeviceId:       request.Service.DeviceId,
			ServiceName:    request.Service.Name,
			RequestId:      request.Id,
		},
	)
}

// Respond respond to an IoT service request
func (b *ServiceBroker) Respond(response *common.ServiceResponse) error {
	// check if the request exists
	_, err := b.getRequest(response.RequestId)
	if err != nil {
		return err
	}

	// check if response already exists
	response_, err := b.getResponse(response.RequestId)
	if _, ok := err.(*common.NotFoundError); err != nil && !ok {
		return err
	}
	if response_ != nil {
		return fmt.Errorf("response already exists")
	}

	return b.responseRegistry.PutState(response)
}

func (b *ServiceBroker) getRequest(requestId string) (*common.ServiceRequest, error) {
	request, err := b.requestRegistry.GetState(requestId)
	if err != nil {
		return nil, err
	}

	return request.(*common.ServiceRequest), nil
}

func (b *ServiceBroker) getResponse(requestId string) (*common.ServiceResponse, error) {
	response, err := b.responseRegistry.GetState(requestId)
	if err != nil {
		return nil, err
	}

	return response.(*common.ServiceResponse), nil
}

// Get return an IoT service request and its response by the request ID
func (b *ServiceBroker) Get(requestId string) (*common.ServiceRequestResponse, error) {
	request, err := b.getRequest(requestId)
	if err != nil {
		return nil, err
	}

	response, err := b.getResponse(requestId)
	if _, ok := err.(*common.NotFoundError); err != nil && !ok {
		return nil, err
	}

	return &common.ServiceRequestResponse{Request: request, Response: response}, nil
}

// GetAll return a list of IoT service requests and their responses by their organization ID, device ID, and service name
func (b *ServiceBroker) GetAll(organizationId string, deviceId string, serviceName string) ([]*common.ServiceRequestResponse, error) {
	states, err := b.indexRegistry.GetStates(organizationId, deviceId, serviceName)
	if err != nil {
		return nil, err
	}

	results := make([]*common.ServiceRequestResponse, 0)

	for _, state := range states {
		index := state.(*serviceRequestIndex)
		request, err := b.getRequest(index.RequestId)
		if err != nil {
			return nil, err
		}
		response, err := b.getResponse(index.RequestId)
		if _, ok := err.(*common.NotFoundError); err != nil && !ok {
			return nil, err
		}

		results = append(results, &common.ServiceRequestResponse{Request: request, Response: response})
	}

	return results, err
}

// Remove remove a (request, response) pair from the ledger
func (b *ServiceBroker) Remove(requestId string) error {
	// remove response from global state, if exists
	response, err := b.getResponse(requestId)
	if _, ok := err.(*common.NotFoundError); err != nil && !ok {
		return err
	}
	if response != nil {
		if err = b.responseRegistry.RemoveState(response); err != nil {
			return err
		}
	}

	// remove request from global state
	request, err := b.getRequest(requestId)
	if err != nil {
		return err
	}
	if err = b.requestRegistry.RemoveState(request); err != nil {
		return err
	}

	// remove index from global state
	index := &serviceRequestIndex{
		OrganizationId: request.Service.OrganizationId,
		DeviceId:       request.Service.DeviceId,
		ServiceName:    request.Service.Name,
		RequestId:      request.Id,
	}
	return b.indexRegistry.RemoveState(index)
}

func createServiceBroker(ctx TransactionContextInterface) *ServiceBroker {
	requestRegistry := new(StateRegistry)
	requestRegistry.ctx = ctx
	requestRegistry.Name = "requests"
	requestRegistry.Deserialize = func(data []byte) (StateInterface, error) {
		return common.DeserializeServiceRequest(data)
	}

	responseRegistry := new(StateRegistry)
	responseRegistry.ctx = ctx
	responseRegistry.Name = "responses"
	responseRegistry.Deserialize = func(data []byte) (StateInterface, error) {
		return common.DeserializeServiceResponse(data)
	}

	indexRegistry := new(StateRegistry)
	indexRegistry.ctx = ctx
	indexRegistry.Name = "request_indices"
	indexRegistry.Deserialize = func(data []byte) (StateInterface, error) {
		return deserializeServiceRequestIndex(data)
	}

	broker := new(ServiceBroker)
	broker.ctx = ctx
	broker.requestRegistry = requestRegistry
	broker.responseRegistry = responseRegistry
	broker.indexRegistry = indexRegistry

	return broker
}
