package contract

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// DeviceRegistrySmartContract smart contract for managing devices on the ledger
type DeviceRegistrySmartContract struct {
	contractapi.Contract
}

// Register create or update a device in the ledger
func (s *DeviceRegistrySmartContract) Register(ctx TransactionContextInterface, data string) error {
	var err error
	var organizationId, deviceId string

	device, err := common.DeserializeDevice([]byte(data))
	if err != nil {
		return err
	}

	if organizationId, err = ctx.GetClientIdentity().GetMSPID(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetClientIdentity().GetID(); err != nil {
		return err
	}

	device.OrganizationId = organizationId
	device.Id = deviceId

	return ctx.GetDeviceRegistry().Register(device)
}

// Get return a device by its organization ID and device ID
func (s *DeviceRegistrySmartContract) Get(ctx TransactionContextInterface, organizationId string, deviceId string) (*common.Device, error) {
	return ctx.GetDeviceRegistry().Get(organizationId, deviceId)
}

// GetAll return a list of devices by their organization ID
func (s *DeviceRegistrySmartContract) GetAll(ctx TransactionContextInterface, organizationId string) ([]*common.Device, error) {
	return ctx.GetDeviceRegistry().GetAll(organizationId)
}

// Deregister remove a device and its services from the ledger
func (s *DeviceRegistrySmartContract) Deregister(ctx TransactionContextInterface) error {
	var err error
	var device *common.Device
	var organizationId, deviceId string

	// check if the device already exists
	if organizationId, err = ctx.GetClientIdentity().GetMSPID(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetClientIdentity().GetID(); err != nil {
		return err
	}
	if device, err = ctx.GetDeviceRegistry().Get(organizationId, deviceId); err != nil {
		return err
	}

	return ctx.GetDeviceRegistry().Deregister(device)
}
