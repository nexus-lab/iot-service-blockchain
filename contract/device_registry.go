package contract

import (
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// DeviceRegistryInterface core utilities for managing devices on the ledger
type DeviceRegistryInterface interface {
	// Register create or update a device in the ledger
	Register(device *common.Device) error

	// Get return a device by its organization ID and device ID
	Get(organizationId string, deviceId string) (*common.Device, error)

	// GetAll return a list of devices by their organization ID
	GetAll(organizationId string) ([]*common.Device, error)

	// Deregister remove a device from the ledger
	Deregister(device *common.Device) error
}

// DeviceRegistry core utilities for managing devices on the ledger
type DeviceRegistry struct {
	ctx           TransactionContextInterface
	stateRegistry StateRegistryInterface
}

// Register create or update a device in the ledger
func (r *DeviceRegistry) Register(device *common.Device) error {
	return r.stateRegistry.PutState(device)
}

// Get return a device by its organization ID and device ID
func (r *DeviceRegistry) Get(organizationId string, deviceId string) (*common.Device, error) {
	state, err := r.stateRegistry.GetState(organizationId, deviceId)
	if err != nil {
		return nil, err
	}

	return state.(*common.Device), nil
}

// GetAll return a list of devices by their organization ID
func (r *DeviceRegistry) GetAll(organizationId string) ([]*common.Device, error) {
	states, err := r.stateRegistry.GetStates(organizationId)
	if err != nil {
		return nil, err
	}

	devices := make([]*common.Device, 0)
	for _, state := range states {
		devices = append(devices, state.(*common.Device))
	}

	return devices, err
}

// Deregister remove a device from the ledger
func (r *DeviceRegistry) Deregister(device *common.Device) error {
	// deregister services of the device
	services, err := r.ctx.GetServiceRegistry().GetAll(device.OrganizationId, device.Id)
	if err != nil {
		return err
	}
	for _, service := range services {
		if err = r.ctx.GetServiceRegistry().Deregister(service); err != nil {
			return err
		}
	}

	return r.stateRegistry.RemoveState(device)
}

func createDeviceRegistry(ctx TransactionContextInterface) *DeviceRegistry {
	stateRegistry := new(StateRegistry)
	stateRegistry.ctx = ctx
	stateRegistry.Name = "devices"
	stateRegistry.Deserialize = func(data []byte) (StateInterface, error) {
		return common.DeserializeDevice(data)
	}

	registry := new(DeviceRegistry)
	registry.ctx = ctx
	registry.stateRegistry = stateRegistry

	return registry
}
