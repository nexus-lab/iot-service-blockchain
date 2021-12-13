package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(new(DeviceRegistrySmartContract), new(ServiceRegistrySmartContract), new(ServiceBrokerSmartContract))

	if err != nil {
		log.Panicf("Failed to create chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Failed to start chaincode: %v", err)
	}
}
