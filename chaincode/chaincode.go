package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/nexus-lab/iot-service-blockchain/contract"
)

func main() {
	deviceRegistryContract := new(contract.DeviceRegistrySmartContract)
	deviceRegistryContract.TransactionContextHandler = new(contract.TransactionContext)
	deviceRegistryContract.Name = "device_registry"

	serviceRegistryContract := new(contract.ServiceRegistrySmartContract)
	serviceRegistryContract.TransactionContextHandler = new(contract.TransactionContext)
	serviceRegistryContract.Name = "service_registry"

	serviceBrokerContract := new(contract.ServiceBrokerSmartContract)
	serviceBrokerContract.TransactionContextHandler = new(contract.TransactionContext)
	serviceBrokerContract.Name = "service_broker"

	chaincode, err := contractapi.NewChaincode(deviceRegistryContract, serviceRegistryContract, serviceBrokerContract)

	if err != nil {
		log.Panicf("Failed to create chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Failed to start chaincode: %v", err)
	}
}
