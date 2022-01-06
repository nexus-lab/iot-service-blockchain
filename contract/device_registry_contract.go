package contract

import (
	"fmt"

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

	if organizationId, err = ctx.GetOrganizationId(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetDeviceId(); err != nil {
		return err
	}

	if device.OrganizationId != organizationId || device.Id != deviceId {
		return fmt.Errorf("cannot register a device other than the requested device")
	}

	err = ctx.GetDeviceRegistry().Register(device)

	// notify listening clients of the update
	if err == nil {
		event := fmt.Sprintf("device://%s/%s/register", device.OrganizationId, device.Id)
		payload, _ := device.Serialize()
		err = ctx.GetStub().SetEvent(event, payload)
	}

	return err
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
func (s *DeviceRegistrySmartContract) Deregister(ctx TransactionContextInterface, data string) error {
	var err error
	var organizationId, deviceId string

	device, err := common.DeserializeDevice([]byte(data))
	if err != nil {
		return err
	}

	if organizationId, err = ctx.GetOrganizationId(); err != nil {
		return err
	}
	if deviceId, err = ctx.GetDeviceId(); err != nil {
		return err
	}

	if device.OrganizationId != organizationId || device.Id != deviceId {
		return fmt.Errorf("cannot deregister a device other than the requested device")
	}

	err = ctx.GetDeviceRegistry().Deregister(device)

	// notify listening clients of the update
	if err == nil {
		event := fmt.Sprintf("device://%s/%s/deregister", device.OrganizationId, device.Id)
		payload, _ := device.Serialize()
		err = ctx.GetStub().SetEvent(event, payload)
	}

	return err
}
