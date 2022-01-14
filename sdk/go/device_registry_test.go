package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DeviceRegistryTestSuite struct {
	suite.Suite
}

func (s *DeviceRegistryTestSuite) TestRegister() {
	contract := new(MockContract)
	deviceRegistry := &DeviceRegistry{contract}

	device := &common.Device{Name: "device1"}
	data, _ := device.Serialize()
	contract.On("SubmitTransaction", "Register", string(data)).Return(nil, nil)

	err := deviceRegistry.Register(device)
	assert.Nil(s.T(), err, "should return no error")

	err = deviceRegistry.Register(nil)
	assert.Error(s.T(), err, "should return error if input is null")

	device = &common.Device{Name: "device2"}
	data, _ = device.Serialize()
	contract.On("SubmitTransaction", "Register", string(data)).Return(nil, errors.New(""))

	err = deviceRegistry.Register(device)
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *DeviceRegistryTestSuite) TestGet() {
	contract := new(MockContract)
	deviceRegistry := &DeviceRegistry{contract}

	expected := new(common.Device)
	data, _ := expected.Serialize()
	contract.On("SubmitTransaction", "Get", "org1", "device1").Return(data, nil)

	actual, err := deviceRegistry.Get("org1", "device1")
	assert.Equal(s.T(), expected, actual, "should return correct device")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "Get", "org2", "device2").Return(nil, new(common.NotFoundError))

	actual, err = deviceRegistry.Get("org2", "device2")
	assert.Nil(s.T(), actual, "should return no device")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")

	contract.On("SubmitTransaction", "Get", "org3", "device3").Return(nil, errors.New(""))

	_, err = deviceRegistry.Get("org3", "device3")
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *DeviceRegistryTestSuite) TestGetAll() {
	contract := new(MockContract)
	deviceRegistry := &DeviceRegistry{contract}

	expected := []*common.Device{new(common.Device), new(common.Device)}
	data, _ := json.Marshal(expected)
	contract.On("SubmitTransaction", "GetAll", "org1").Return(data, nil)

	actual, err := deviceRegistry.GetAll("org1")
	assert.Equal(s.T(), expected, actual, "should return correct devices")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "GetAll", "org2").Return([]byte("[]"), nil)

	actual, err = deviceRegistry.GetAll("org2")
	assert.Zero(s.T(), len(actual), "should return no device")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "GetAll", "org3").Return(nil, errors.New(""))

	_, err = deviceRegistry.GetAll("org3")
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *DeviceRegistryTestSuite) TestDeregister() {
	contract := new(MockContract)
	deviceRegistry := &DeviceRegistry{contract}

	device := &common.Device{Name: "device1"}
	data, _ := device.Serialize()
	contract.On("SubmitTransaction", "Deregister", string(data)).Return(nil, nil)

	err := deviceRegistry.Deregister(device)
	assert.Nil(s.T(), err, "should return no error")

	err = deviceRegistry.Deregister(nil)
	assert.Error(s.T(), err, "should return error if input is null")

	device = &common.Device{Name: "device2"}
	data, _ = device.Serialize()
	contract.On("SubmitTransaction", "Deregister", string(data)).Return(nil, errors.New(""))

	err = deviceRegistry.Deregister(device)
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *DeviceRegistryTestSuite) TestRegisterEvent() {
	contract := new(MockContract)
	deviceRegistry := &DeviceRegistry{contract}

	eventChannel := make(chan *client.ChaincodeEvent)
	go func() {
		for i := 0; i < 5; i++ {
			data, _ := (&common.Device{Name: fmt.Sprintf("device%d", i)}).Serialize()
			eventChannel <- &client.ChaincodeEvent{
				EventName: fmt.Sprintf("device://org%d/device%d/register", i, i),
				Payload:   data,
			}
		}
	}()

	var cancelFunc context.CancelFunc = func() {
		close(eventChannel)
	}

	contract.On("RegisterEvent", mock.Anything).Return(eventChannel, cancelFunc, nil)

	source, cancel, err := deviceRegistry.RegisterEvent()
	defer cancel()
	assert.Nil(s.T(), err, "should return no error")
	assert.IsType(s.T(), *new(context.CancelFunc), cancel, "should return correct cancel function")

	for i := 0; i < 5; i++ {
		event := <-source
		assert.Equal(s.T(), "register", event.Action, "should return correct action")
		assert.Equal(s.T(), fmt.Sprintf("org%d", i), event.OrganizationId, "should return correct organization ID")
		assert.Equal(s.T(), fmt.Sprintf("device%d", i), event.DeviceId, "should return correct device ID")
		assert.IsType(s.T(), new(common.Device), event.Payload, "should return parsed device as event payload")
		assert.Equal(s.T(), fmt.Sprintf("device%d", i), event.Payload.(*common.Device).Name, "should return correct event payload")
	}

	contract = new(MockContract)
	deviceRegistry = &DeviceRegistry{contract}
	contract.On("RegisterEvent", mock.Anything).Return(nil, nil, errors.New(""))

	_, _, err = deviceRegistry.RegisterEvent()
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func TestDeviceRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceRegistryTestSuite))
}
