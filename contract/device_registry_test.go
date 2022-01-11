package contract

import (
	"testing"

	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockDeviceRegistry struct {
	mock.Mock
}

func (r *MockDeviceRegistry) Register(device *common.Device) error {
	args := r.Called(device)
	return args.Error(0)
}

func (r *MockDeviceRegistry) Get(organizationId string, deviceId string) (*common.Device, error) {
	args := r.Called(organizationId, deviceId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common.Device), args.Error(1)
}

func (r *MockDeviceRegistry) GetAll(organizationId string) ([]*common.Device, error) {
	args := r.Called(organizationId)
	return args.Get(0).([]*common.Device), args.Error(1)
}

func (r *MockDeviceRegistry) Deregister(device *common.Device) error {
	args := r.Called(device)
	return args.Error(0)
}

type DeviceRegistryTestSuite struct {
	suite.Suite
}

func (s *DeviceRegistryTestSuite) TestRegister() {
	stateRegistry := new(MockStateRegistry)

	deviceRegistry := new(DeviceRegistry)
	deviceRegistry.ctx = new(MockTransactionContext)
	deviceRegistry.stateRegistry = stateRegistry

	device := new(common.Device)
	stateRegistry.On("PutState", device).Return(nil)

	err := deviceRegistry.Register(device)
	called := stateRegistry.AssertCalled(s.T(), "PutState", device)
	assert.True(s.T(), called, "should put device to state registry")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *DeviceRegistryTestSuite) TestGet() {
	stateRegistry := new(MockStateRegistry)

	deviceRegistry := new(DeviceRegistry)
	deviceRegistry.ctx = new(MockTransactionContext)
	deviceRegistry.stateRegistry = stateRegistry

	device := new(common.Device)
	stateRegistry.On("GetState", []string{"org1", "device1"}).Return(device, nil)
	stateRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))

	result, err := deviceRegistry.Get("org1", "device1")
	assert.Equal(s.T(), device, result, "should return the correct device")
	assert.Nil(s.T(), err, "should return no error")

	result, err = deviceRegistry.Get("org2", "device2")
	assert.Nil(s.T(), result, "should return no device")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func (s *DeviceRegistryTestSuite) TestGetAll() {
	stateRegistry := new(MockStateRegistry)

	deviceRegistry := new(DeviceRegistry)
	deviceRegistry.ctx = new(MockTransactionContext)
	deviceRegistry.stateRegistry = stateRegistry

	devices := []StateInterface{new(common.Device), new(common.Device)}
	stateRegistry.On("GetStates", []string{"org1"}).Return(devices, nil)
	stateRegistry.On("GetStates", mock.Anything).Return([]StateInterface{}, nil)

	results, err := deviceRegistry.GetAll("org1")
	assert.Equal(s.T(), len(devices), len(results), "should return the correct number of devices")
	assert.Nil(s.T(), err, "should return no error")
	for i := range results {
		assert.Equal(s.T(), devices[i], results[i], "should return correct device")
	}

	results, err = deviceRegistry.GetAll("org2")
	assert.Zero(s.T(), len(results), "should return no device")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *DeviceRegistryTestSuite) TestDeregister() {
	stateRegistry := new(MockStateRegistry)
	serviceRegistry := new(MockServiceRegistry)
	transactionContext := new(MockTransactionContext)

	transactionContext.serviceRegistry = serviceRegistry

	deviceRegistry := new(DeviceRegistry)
	deviceRegistry.ctx = transactionContext
	deviceRegistry.stateRegistry = stateRegistry

	device := new(common.Device)
	device.Id = "device1"
	device.OrganizationId = "org1"

	services := []*common.Service{new(common.Service), new(common.Service)}

	serviceRegistry.On("GetAll", "org1", "device1").Return(services, nil)
	serviceRegistry.On("Deregister", mock.AnythingOfType("*common.Service")).Return(nil)
	stateRegistry.On("RemoveState", device).Return(nil)

	err := deviceRegistry.Deregister(device)
	called := stateRegistry.AssertCalled(s.T(), "RemoveState", device)
	assert.True(s.T(), called, "should remove device from state registry")
	assert.Nil(s.T(), err, "should return no error")

	called = serviceRegistry.AssertCalled(s.T(), "Deregister", services[1])
	assert.True(s.T(), called, "should deregister service by the service registry")
}

func TestDeviceRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceRegistryTestSuite))
}
