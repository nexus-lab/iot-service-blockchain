package contract

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// ServiceBrokerSmartContract smart contract for managing IoT service requests and responses
type ServiceBrokerSmartContract struct {
	contractapi.Contract
}

// Request make a request to an IoT service
func (s *ServiceBrokerSmartContract) Request(ctx TransactionContextInterface, data string) error {
	request, err := common.DeserializeServiceRequest([]byte(data))
	if err != nil {
		return err
	}

	return ctx.GetServiceBroker().Request(request)
}

// Respond respond to an IoT service request
func (s *ServiceBrokerSmartContract) Respond(ctx TransactionContextInterface, data string) error {
	var err error
	var organizationId, deviceId string
	var response *common.ServiceResponse

	if response, err = common.DeserializeServiceResponse([]byte(data)); err != nil {
		return err
	}

	// check if corresponding request exists
	pair, err := ctx.GetServiceBroker().Get(response.RequestId)
	if err != nil {
		return err
	}

	// check if the client creating the response is the client requested for service
	if organizationId, err = ctx.GetClientIdentity().GetMSPID(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetClientIdentity().GetID(); err != nil {
		return err
	}
	if pair.Request.Service.OrganizationId != organizationId || pair.Request.Service.DeviceId != deviceId {
		return fmt.Errorf("cannot create response from a device other than the requested device")
	}

	return ctx.GetServiceBroker().Respond(response)
}

// Get return an IoT service request and its response by the request ID
func (s *ServiceBrokerSmartContract) Get(ctx TransactionContextInterface, requestId string) (*common.ServiceRequestResponse, error) {
	return ctx.GetServiceBroker().Get(requestId)
}

// GetAll return a list of IoT service requests and their responses by their organization ID, device ID, and service name
func (s *ServiceBrokerSmartContract) GetAll(ctx TransactionContextInterface, organizationId string, deviceId string, serviceName string) ([]*common.ServiceRequestResponse, error) {
	return ctx.GetServiceBroker().GetAll(organizationId, deviceId, serviceName)
}

// Remove remove a (request, response) pair from the ledger
func (s *ServiceBrokerSmartContract) Remove(ctx TransactionContextInterface, requestId string) error {
	var err error
	var organizationId, deviceId string

	// check if corresponding request exists
	pair, err := ctx.GetServiceBroker().Get(requestId)
	if err != nil {
		return err
	}

	// check if the client creating the response is the client requested for service
	if organizationId, err = ctx.GetClientIdentity().GetMSPID(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetClientIdentity().GetID(); err != nil {
		return err
	}
	if pair.Request.Service.OrganizationId != organizationId || pair.Request.Service.DeviceId != deviceId {
		return fmt.Errorf("cannot remove response from a device other than the requested device")
	}

	return ctx.GetServiceBroker().Remove(requestId)
}
