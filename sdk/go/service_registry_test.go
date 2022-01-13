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

type ServiceRegistryTestSuite struct {
	suite.Suite
}

func (s *ServiceRegistryTestSuite) TestRegister() {
	contract := new(MockContract)
	serviceRegistry := &ServiceRegistry{contract}

	service := &common.Service{Name: "service1"}
	data, _ := service.Serialize()
	contract.On("SubmitTransaction", "Register", string(data)).Return(nil, nil)

	err := serviceRegistry.Register(service)
	assert.Nil(s.T(), err, "should return no error")

	err = serviceRegistry.Register(nil)
	assert.Error(s.T(), err, "should return error if input is null")

	service = &common.Service{Name: "service2"}
	data, _ = service.Serialize()
	contract.On("SubmitTransaction", "Register", string(data)).Return(nil, errors.New(""))

	err = serviceRegistry.Register(service)
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceRegistryTestSuite) TestGet() {
	contract := new(MockContract)
	serviceRegistry := &ServiceRegistry{contract}

	expected := new(common.Service)
	data, _ := expected.Serialize()
	contract.On("SubmitTransaction", "Get", "org1", "device1", "service1").Return(data, nil)

	actual, err := serviceRegistry.Get("org1", "device1", "service1")
	assert.Equal(s.T(), expected, actual, "should return correct service")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "Get", "org2", "device2", "service2").Return(nil, new(common.NotFoundError))

	actual, err = serviceRegistry.Get("org2", "device2", "service2")
	assert.Nil(s.T(), actual, "should return no service")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")

	contract.On("SubmitTransaction", "Get", "org3", "device3", "service3").Return(nil, errors.New(""))

	_, err = serviceRegistry.Get("org3", "device3", "service3")
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceRegistryTestSuite) TestGetAll() {
	contract := new(MockContract)
	serviceRegistry := &ServiceRegistry{contract}

	expected := []*common.Service{new(common.Service), new(common.Service)}
	data, _ := json.Marshal(expected)
	contract.On("SubmitTransaction", "GetAll", "org1", "device1").Return(data, nil)

	actual, err := serviceRegistry.GetAll("org1", "device1")
	assert.Equal(s.T(), expected, actual, "should return correct services")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "GetAll", "org2", "device2").Return([]byte("[]"), nil)

	actual, err = serviceRegistry.GetAll("org2", "device2")
	assert.Zero(s.T(), len(actual), "should return no service")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "GetAll", "org3", "device3").Return(nil, errors.New(""))

	_, err = serviceRegistry.GetAll("org3", "device3")
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceRegistryTestSuite) TestDeregister() {
	contract := new(MockContract)
	serviceRegistry := &ServiceRegistry{contract}

	service := &common.Service{Name: "service1"}
	data, _ := service.Serialize()
	contract.On("SubmitTransaction", "Deregister", string(data)).Return(nil, nil)

	err := serviceRegistry.Deregister(service)
	assert.Nil(s.T(), err, "should return no error")

	err = serviceRegistry.Deregister(nil)
	assert.Error(s.T(), err, "should return error if input is null")

	service = &common.Service{Name: "service2"}
	data, _ = service.Serialize()
	contract.On("SubmitTransaction", "Deregister", string(data)).Return(nil, errors.New(""))

	err = serviceRegistry.Deregister(service)
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceRegistryTestSuite) TestRegisterEvent() {
	contract := new(MockContract)
	serviceRegistry := &ServiceRegistry{contract}

	eventChannel := make(chan *client.ChaincodeEvent)
	go func() {
		for i := 0; i < 5; i++ {
			data, _ := (&common.Service{Name: fmt.Sprintf("service%d", i)}).Serialize()
			eventChannel <- &client.ChaincodeEvent{
				EventName: fmt.Sprintf("service://org%d/device%d/service%d/register", i, i, i),
				Payload:   data,
			}
		}
	}()

	var cancelFunc context.CancelFunc = func() {
		close(eventChannel)
	}

	contract.On("RegisterEvent", mock.Anything).Return(eventChannel, cancelFunc, nil)

	source, cancel, err := serviceRegistry.RegisterEvent()
	assert.Nil(s.T(), err, "should return no error")
	assert.IsType(s.T(), *new(context.CancelFunc), cancel, "should return correct cancel function")

	for i := 0; i < 5; i++ {
		event := <-source
		assert.Equal(s.T(), "register", event.Action, "should return correct action")
		assert.Equal(s.T(), fmt.Sprintf("org%d", i), event.OrganizationId, "should return correct organization ID")
		assert.Equal(s.T(), fmt.Sprintf("device%d", i), event.DeviceId, "should return correct device ID")
		assert.Equal(s.T(), fmt.Sprintf("service%d", i), event.ServiceName, "should return correct service name")
		assert.IsType(s.T(), new(common.Service), event.Payload, "should return parsed service as event payload")
		assert.Equal(s.T(), fmt.Sprintf("service%d", i), event.Payload.(*common.Service).Name, "should return correct event payload")
	}

	cancel()

	contract = new(MockContract)
	serviceRegistry = &ServiceRegistry{contract}
	contract.On("RegisterEvent", mock.Anything).Return(nil, nil, errors.New(""))

	_, _, err = serviceRegistry.RegisterEvent()
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func TestServiceRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceRegistryTestSuite))
}
