package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TODO: organization msp ?= clientIdentity.msp
func GetOrganizationId(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("Failed to get the organization identify of the calling client")
	}
	return id, nil
}

func GetDeviceId(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to get the client identify of the calling client")
	}

	return id, nil
}
