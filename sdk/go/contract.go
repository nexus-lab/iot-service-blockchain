package sdk

import (
	"context"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type ContractInterface interface {
	// SubmitTransaction submit a transaction to the ledger
	SubmitTransaction(name string, args ...string) ([]byte, error)

	// RegisterEvent register for chaincode events
	RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *client.ChaincodeEvent, context.CancelFunc, error)
}

type Contract struct {
	network      *client.Network
	chaincodeId  string
	contractName string
}

// SubmitTransaction submit a transaction to the ledger
func (c *Contract) SubmitTransaction(name string, args ...string) ([]byte, error) {
	return c.network.GetContractWithName(c.chaincodeId, c.contractName).SubmitTransaction(name, args...)
}

// RegisterEvent register for chaincode events
func (c *Contract) RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *client.ChaincodeEvent, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	source, err := c.network.ChaincodeEvents(ctx, c.chaincodeId, options...)
	return source, cancel, err
}
