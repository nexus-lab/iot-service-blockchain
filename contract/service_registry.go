package contract

import (
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// ServiceRegistryInterface core utilities for managing services on the ledger
type ServiceRegistryInterface interface {
	// Register create or update a service in the ledger
	Register(service *common.Service) error

	// Get return a service by its organization ID, device ID, and name
	Get(organizationId string, deviceId string, serviceName string) (*common.Service, error)

	// GetAll return a list of services by their organization ID and device ID
	GetAll(organizationId string, deviceId string) ([]*common.Service, error)

	// Deregister remove a service from the ledger
	Deregister(service *common.Service) error
}

// ServiceRegistry core utilities for managing services on the ledger
type ServiceRegistry struct {
	ctx           TransactionContextInterface
	stateRegistry StateRegistryInterface
}

// Register create or update a service in the ledger
func (r *ServiceRegistry) Register(service *common.Service) error {
	// check if device exists
	_, err := r.ctx.GetDeviceRegistry().Get(service.OrganizationId, service.DeviceId)
	if err != nil {
		return err
	}

	return r.stateRegistry.PutState(service)
}

// Get return a service by its organization ID, device ID, and name
func (r *ServiceRegistry) Get(organizationId string, deviceId string, name string) (*common.Service, error) {
	state, err := r.stateRegistry.GetState(organizationId, deviceId, name)
	if err != nil {
		return nil, err
	}

	return state.(*common.Service), nil
}

// GetAll return a list of services by their organization ID and device ID
func (r *ServiceRegistry) GetAll(organizationId string, deviceId string) ([]*common.Service, error) {
	states, err := r.stateRegistry.GetStates(organizationId, deviceId)
	if err != nil {
		return nil, err
	}

	services := make([]*common.Service, 0)
	for _, state := range states {
		services = append(services, state.(*common.Service))
	}

	return services, err
}

// Deregister remove a service from the ledger
func (r *ServiceRegistry) Deregister(service *common.Service) error {
	// remove related requests and responses
	pairs, err := r.ctx.GetServiceBroker().GetAll(service.OrganizationId, service.DeviceId, service.Name)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		if err = r.ctx.GetServiceBroker().Remove(pair.Request.Id); err != nil {
			return err
		}
	}

	return r.stateRegistry.RemoveState(service)
}

func createServiceRegistry(ctx TransactionContextInterface) *ServiceRegistry {
	stateRegistry := new(StateRegistry)
	stateRegistry.ctx = ctx
	stateRegistry.Name = "services"
	stateRegistry.Deserialize = func(data []byte) (StateInterface, error) {
		return common.DeserializeService(data)
	}

	registry := new(ServiceRegistry)
	registry.ctx = ctx
	registry.stateRegistry = stateRegistry

	return registry
}
