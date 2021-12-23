package contract

import (
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// DeviceRegistry core utilities for managing devices on the ledger
type DeviceRegistry struct {
	stateRegistry StateRegistry
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
	ctx := r.stateRegistry.Ctx.(*TransactionContext)

	// deregister services of the device
	services, err := ctx.GetServiceRegistry().GetAll(device.OrganizationId, device.Id)
	if err != nil {
		return err
	}
	for _, service := range services {
		if err = ctx.GetServiceRegistry().Deregister(service); err != nil {
			return err
		}
	}

	return r.stateRegistry.RemoveState(device)
}

// CreateDeviceRegistry create a new device registry from transaction context
func CreateDeviceRegistry(ctx TransactionContextInterface) *DeviceRegistry {
	stateRegistry := new(StateRegistry)
	stateRegistry.Ctx = ctx
	stateRegistry.Name = "devices"
	stateRegistry.Deserialize = func(data []byte) (common.StateInterface, error) {
		return common.DeserializeDevice(data)
	}

	registry := new(DeviceRegistry)
	registry.stateRegistry = *stateRegistry

	return registry
}
