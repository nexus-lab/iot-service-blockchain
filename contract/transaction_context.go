package contract

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface extends the default transaction context with specific services
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface

	// GetDeviceRegistry get the default instance of device registry
	GetDeviceRegistry() DeviceRegistryInterface

	// GetServiceRegistry get the default instance of service registry
	GetServiceRegistry() ServiceRegistryInterface

	// GetServiceBroker get the default instance of service broker
	GetServiceBroker() ServiceBrokerInterface
}

// TransactionContext an implementation of TransactionContextInterface
type TransactionContext struct {
	contractapi.TransactionContext
	deviceRegistry  DeviceRegistryInterface
	serviceRegistry ServiceRegistryInterface
	serviceBroker   ServiceBrokerInterface
}

// GetDeviceRegistry get the device registry instance
func (c *TransactionContext) GetDeviceRegistry() DeviceRegistryInterface {
	if c.deviceRegistry == nil {
		c.deviceRegistry = createDeviceRegistry(c)
	}

	return c.deviceRegistry
}

// GetServiceRegistry get the device registry instance
func (c *TransactionContext) GetServiceRegistry() ServiceRegistryInterface {
	if c.serviceRegistry == nil {
		c.serviceRegistry = createServiceRegistry(c)
	}

	return c.serviceRegistry
}

// GetServiceBroker get the device broker instance
func (c *TransactionContext) GetServiceBroker() ServiceBrokerInterface {
	if c.serviceBroker == nil {
		c.serviceBroker = createServiceBroker(c)
	}

	return c.serviceBroker
}
