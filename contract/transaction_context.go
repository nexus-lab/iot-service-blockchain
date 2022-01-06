package contract

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/nexus-lab/iot-service-blockchain/common"
)

// TransactionContextInterface extends the default transaction context with specific services
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface

	// GetOrganizationId return the organization MSP ID
	GetOrganizationId() (string, error)

	// GetDeviceId returns the ID associated with the invoking identity which is unique within the MSP
	GetDeviceId() (string, error)

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

// GetOrganizationId return the organization MSP ID
func (c *TransactionContext) GetOrganizationId() (string, error) {
	return c.GetClientIdentity().GetMSPID()
}

// GetDeviceId returns the ID associated with the invoking identity which is unique within the MSP
func (c *TransactionContext) GetDeviceId() (string, error) {
	cert, err := c.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return "", err
	}
	return common.GetClientId(cert)
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
