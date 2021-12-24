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
	request, _ := ctx.GetServiceBroker().(*ServiceBroker).GetRequest(response.RequestId)
	if request == nil {
		return fmt.Errorf("cannot get corresponding request")
	}

	// check if the client creating the response is the client requested for service
	if organizationId, err = ctx.GetClientIdentity().GetMSPID(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetClientIdentity().GetID(); err != nil {
		return err
	}
	if request.Service.OrganizationId != organizationId || request.Service.DeviceId != deviceId {
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
