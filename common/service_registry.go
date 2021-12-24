package common

// ServiceRegistryInterface core utilities for managing services on the ledger
type ServiceRegistryInterface interface {
	// Register create or update a service in the ledger
	Register(service *Service) error

	// Get return a service by its organization ID, device ID, and name
	Get(organizationId string, deviceId string, serviceName string) (*Service, error)

	// GetAll return a list of services by their organization ID and device ID
	GetAll(organizationId string, deviceId string) ([]*Service, error)

	// Deregister remove a service from the ledger
	Deregister(service *Service) error
}
