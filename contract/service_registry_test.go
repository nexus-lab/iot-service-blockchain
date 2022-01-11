package contract

import (
	"testing"

	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockServiceRegistry struct {
	mock.Mock
}

func (r *MockServiceRegistry) Register(service *common.Service) error {
	args := r.Called(service)
	return args.Error(0)
}

func (r *MockServiceRegistry) Get(organizationId string, deviceId string, serviceName string) (*common.Service, error) {
	args := r.Called(organizationId, deviceId, serviceName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common.Service), args.Error(1)
}

func (r *MockServiceRegistry) GetAll(organizationId string, deviceId string) ([]*common.Service, error) {
	args := r.Called(organizationId, deviceId)
	return args.Get(0).([]*common.Service), args.Error(1)
}

func (r *MockServiceRegistry) Deregister(service *common.Service) error {
	args := r.Called(service)
	return args.Error(0)
}

type ServiceRegistryTestSuite struct {
	suite.Suite
}

func (s *ServiceRegistryTestSuite) TestRegister() {
	stateRegistry := new(MockStateRegistry)
	deviceRegistry := new(MockDeviceRegistry)
	transactionContext := new(MockTransactionContext)

	transactionContext.deviceRegistry = deviceRegistry

	serviceRegistry := new(ServiceRegistry)
	serviceRegistry.ctx = transactionContext
	serviceRegistry.stateRegistry = stateRegistry

	service := new(common.Service)
	service.OrganizationId = "org1"
	service.DeviceId = "device1"
	service.Name = "service1"

	stateRegistry.On("PutState", service).Return(nil)
	deviceRegistry.On("Get", "org1", "device1").Return(new(common.Device), nil)
	deviceRegistry.On("Get", mock.Anything, mock.Anything).Return(nil, new(common.NotFoundError))

	err := serviceRegistry.Register(service)
	called := stateRegistry.AssertCalled(s.T(), "PutState", service)
	assert.True(s.T(), called, "should put service to state registry")
	assert.Nil(s.T(), err, "should return no error")

	service = new(common.Service)
	service.OrganizationId = "org2"
	service.DeviceId = "device2"

	err = serviceRegistry.Register(service)
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func (s *ServiceRegistryTestSuite) TestGet() {
	stateRegistry := new(MockStateRegistry)

	serviceRegistry := new(ServiceRegistry)
	serviceRegistry.ctx = new(MockTransactionContext)
	serviceRegistry.stateRegistry = stateRegistry

	service := new(common.Service)
	stateRegistry.On("GetState", []string{"org1", "device1", "service1"}).Return(service, nil)
	stateRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))

	result, err := serviceRegistry.Get("org1", "device1", "service1")
	assert.Equal(s.T(), service, result, "should return the correct service")
	assert.Nil(s.T(), err, "should return no error")

	result, err = serviceRegistry.Get("org2", "device2", "service2")
	assert.Nil(s.T(), result, "should return no service")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func (s *ServiceRegistryTestSuite) TestGetAll() {
	stateRegistry := new(MockStateRegistry)

	serviceRegistry := new(ServiceRegistry)
	serviceRegistry.ctx = new(MockTransactionContext)
	serviceRegistry.stateRegistry = stateRegistry

	services := []StateInterface{new(common.Service), new(common.Service)}
	stateRegistry.On("GetStates", []string{"org1", "device1"}).Return(services, nil)
	stateRegistry.On("GetStates", mock.Anything).Return([]StateInterface{}, nil)

	results, err := serviceRegistry.GetAll("org1", "device1")
	assert.Equal(s.T(), len(services), len(results), "should return the correct number of services")
	assert.Nil(s.T(), err, "should return no error")
	for i := range results {
		assert.Equal(s.T(), services[i], results[i], "should return correct service")
	}

	results, err = serviceRegistry.GetAll("org2", "device2")
	assert.Zero(s.T(), len(results), "should return no service")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceRegistryTestSuite) TestDeregister() {
	stateRegistry := new(MockStateRegistry)
	serviceBroker := new(MockServiceBroker)
	transactionContext := new(MockTransactionContext)

	transactionContext.serviceBroker = serviceBroker

	serviceRegistry := new(ServiceRegistry)
	serviceRegistry.ctx = transactionContext
	serviceRegistry.stateRegistry = stateRegistry

	service := new(common.Service)
	service.DeviceId = "device1"
	service.OrganizationId = "org1"
	service.Name = "service1"

	pairs := []*common.ServiceRequestResponse{
		{Request: &common.ServiceRequest{Id: "request1"}},
		{Request: &common.ServiceRequest{Id: "request2"}},
	}

	serviceBroker.On("GetAll", "org1", "device1", "service1").Return(pairs, nil)
	serviceBroker.On("Remove", mock.Anything).Return(nil)
	stateRegistry.On("RemoveState", service).Return(nil)

	err := serviceRegistry.Deregister(service)
	called := stateRegistry.AssertCalled(s.T(), "RemoveState", service)
	assert.True(s.T(), called, "should remove service from state registry")
	assert.Nil(s.T(), err, "should return no error")

	called = serviceBroker.AssertCalled(s.T(), "Remove", "request2")
	assert.True(s.T(), called, "should remove service (request, response) pairs by the service broker")
}

func TestServiceRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceRegistryTestSuite))
}
