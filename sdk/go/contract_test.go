package sdk

import (
	"context"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/stretchr/testify/mock"
)

type MockContract struct {
	mock.Mock
}

func (c *MockContract) SubmitTransaction(name string, args_ ...string) ([]byte, error) {
	args__ := make([]interface{}, 0)
	args__ = append(args__, name)
	for _, arg := range args_ {
		args__ = append(args__, arg)
	}

	args := c.Called(args__...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}

func (c *MockContract) RegisterEvent(options ...client.ChaincodeEventsOption) (<-chan *client.ChaincodeEvent, context.CancelFunc, error) {
	options_ := make([]interface{}, 0)
	for _, option := range options {
		options_ = append(options_, option)
	}
	args := c.Called(options_...)

	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}

	if args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}

	return args.Get(0).(chan *client.ChaincodeEvent), args.Get(1).(context.CancelFunc), args.Error(2)
}
