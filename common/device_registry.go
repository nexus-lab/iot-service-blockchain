package common

// DeviceRegistryInterface core utilities for managing devices on the ledger
type DeviceRegistryInterface interface {
	// Register create or update a device in the ledger
	Register(device *Device) error

	// Get return a device by its organization ID and device ID
	Get(organizationId string, deviceId string) (*Device, error)

	// GetAll return a list of devices by their organization ID
	GetAll(organizationId string) ([]*Device, error)

	// Deregister remove a device from the ledger
	Deregister(device *Device) error
}
