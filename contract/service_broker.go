package contract

import (
	"encoding/json"
	"fmt"

	"github.com/nexus-lab/iot-service-blockchain/common"
)

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
	requestRegistry  StateRegistry
	responseRegistry StateRegistry
	indexRegistry    StateRegistry
}

// Request make a request to an IoT service
func (b *ServiceBroker) Request(request *common.ServiceRequest) error {
	// check if request already exists
	request_, err := b.GetRequest(request.Id)
	if request_ != nil {
		return fmt.Errorf("Request already exists")
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
	_, err := b.GetRequest(response.RequestId)
	if err != nil {
		return err
	}

	// check if response already exists
	response_, err := b.GetResponse(response.RequestId)
	if _, ok := err.(*common.NotFoundError); err != nil && !ok {
		return err
	}
	if response_ != nil {
		return fmt.Errorf("Response already exists")
	}

	return b.responseRegistry.PutState(response)
}

// GetRequest return an IoT service request by its ID
func (b *ServiceBroker) GetRequest(requestId string) (*common.ServiceRequest, error) {
	request, err := b.requestRegistry.GetState(requestId)
	if err != nil {
		return nil, err
	}

	return request.(*common.ServiceRequest), nil
}

// GetResponse return an IoT service response by its request ID
func (b *ServiceBroker) GetResponse(requestId string) (*common.ServiceResponse, error) {
	response, err := b.responseRegistry.GetState(requestId)
	if err != nil {
		return nil, err
	}

	return response.(*common.ServiceResponse), nil
}

// Get return an IoT service request and its response by the request ID
func (b *ServiceBroker) Get(requestId string) (*common.ServiceRequestResponse, error) {
	request, err := b.GetRequest(requestId)
	if err != nil {
		return nil, err
	}

	response, err := b.GetResponse(requestId)
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
		request, err := b.GetRequest(index.RequestId)
		if err != nil {
			return nil, err
		}
		response, err := b.GetResponse(index.RequestId)
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
	response, err := b.GetResponse(requestId)
	if _, ok := err.(*common.NotFoundError); err != nil && !ok {
		return err
	}
	if response != nil {
		if err = b.responseRegistry.RemoveState(response); err != nil {
			return err
		}
	}

	// remove request from global state
	request, err := b.GetRequest(requestId)
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

// CreateServiceBroker create a new service broker from transaction context
func CreateServiceBroker(ctx TransactionContextInterface) *ServiceBroker {
	requestRegistry := new(StateRegistry)
	requestRegistry.Ctx = ctx
	requestRegistry.Name = "requests"
	requestRegistry.Deserialize = func(data []byte) (common.StateInterface, error) {
		return common.DeserializeServiceRequest(data)
	}

	responseRegistry := new(StateRegistry)
	responseRegistry.Ctx = ctx
	responseRegistry.Name = "responses"
	responseRegistry.Deserialize = func(data []byte) (common.StateInterface, error) {
		return common.DeserializeServiceResponse(data)
	}

	indexRegistry := new(StateRegistry)
	indexRegistry.Ctx = ctx
	indexRegistry.Name = "request_indices"
	indexRegistry.Deserialize = func(data []byte) (common.StateInterface, error) {
		return deserializeServiceRequestIndex(data)
	}

	broker := new(ServiceBroker)
	broker.requestRegistry = *requestRegistry
	broker.responseRegistry = *responseRegistry
	broker.indexRegistry = *indexRegistry

	return broker
}
