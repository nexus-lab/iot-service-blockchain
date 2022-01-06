package contract

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// ServiceRegistrySmartContract smart contract for managing IoT services on the ledger
type ServiceRegistrySmartContract struct {
	contractapi.Contract
}

// Register create or update an IoT service in the ledger
func (s *ServiceRegistrySmartContract) Register(ctx TransactionContextInterface, data string) error {
	var err error
	var organizationId, deviceId string

	service, err := common.DeserializeService([]byte(data))
	if err != nil {
		return err
	}

	if organizationId, err = ctx.GetClientIdentity().GetMSPID(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetClientIdentity().GetID(); err != nil {
		return err
	}

	if service.OrganizationId != organizationId || service.DeviceId != deviceId {
		return fmt.Errorf("cannot register a service other than one of the requested device")
	}

	return ctx.GetServiceRegistry().Register(service)
}

// Get return a device by its organization ID, device ID, and name
func (s *ServiceRegistrySmartContract) Get(ctx TransactionContextInterface, organizationId string, deviceId string, name string) (*common.Service, error) {
	return ctx.GetServiceRegistry().Get(organizationId, deviceId, name)
}

// GetAll return a list of devices by their organization ID and device ID
func (s *ServiceRegistrySmartContract) GetAll(ctx TransactionContextInterface, organizationId string, deviceId string) ([]*common.Service, error) {
	return ctx.GetServiceRegistry().GetAll(organizationId, deviceId)
}

// Deregister remove an IoT service and its request/responses from the ledger
func (s *ServiceRegistrySmartContract) Deregister(ctx TransactionContextInterface, data string) error {
	var err error
	var organizationId, deviceId string

	service, err := common.DeserializeService([]byte(data))
	if err != nil {
		return err
	}

	if organizationId, err = ctx.GetClientIdentity().GetMSPID(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetClientIdentity().GetID(); err != nil {
		return err
	}

	if service.OrganizationId != organizationId || service.DeviceId != deviceId {
		return fmt.Errorf("cannot deregister a service other than one of the requested device")
	}

	return ctx.GetServiceRegistry().Deregister(service)
}
