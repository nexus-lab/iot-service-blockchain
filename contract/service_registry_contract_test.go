package contract

import (
	"testing"

	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceRegistryContractTestSuite struct {
	suite.Suite
}

func (s *ServiceRegistryContractTestSuite) TestRegister() {
	ctx := new(MockTransactionContext)
	serviceRegistry := new(MockServiceRegistry)
	ctx.serviceRegistry = serviceRegistry

	deviceId, _ := ctx.GetClientIdentity().GetID()
	organizationId, _ := ctx.GetClientIdentity().GetMSPID()
	serviceRegistry.On("Register", mock.AnythingOfType("*common.Service")).Return(nil)

	contract := new(ServiceRegistrySmartContract)
	err := contract.Register(ctx, "{\"name\":\"Service1\",\"Version\":1,\"description\":\"Service of Device1\",\"organizationId\":\"Org1MSP\",\"deviceId\":\"Device1Id\",\"lastUpdateTime\":\"2021-12-12T17:36:00-05:00\"}")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceRegistry.AssertCalled(s.T(), "Register", mock.AnythingOfType("*common.Service"))
	assert.True(s.T(), called, "should put service to service registry")
	service := serviceRegistry.Calls[0].Arguments[0].(*common.Service)
	assert.Equal(s.T(), deviceId, service.DeviceId, "should change device ID")
	assert.Equal(s.T(), organizationId, service.OrganizationId, "should change organization ID")

	err = contract.Register(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
}

func (s *ServiceRegistryContractTestSuite) TestGet() {
	ctx := new(MockTransactionContext)
	serviceRegistry := new(MockServiceRegistry)
	ctx.serviceRegistry = serviceRegistry

	serviceRegistry.On("Get", "Org1MSP", "Device1Id", "Service1").Return(new(common.Service), nil)

	contract := new(ServiceRegistrySmartContract)
	_, _ = contract.Get(ctx, "Org1MSP", "Device1Id", "Service1")
	called := serviceRegistry.AssertCalled(s.T(), "Get", "Org1MSP", "Device1Id", "Service1")
	assert.True(s.T(), called, "should retrieve service from service registry")
}

func (s *ServiceRegistryContractTestSuite) TestGetAll() {
	ctx := new(MockTransactionContext)
	serviceRegistry := new(MockServiceRegistry)
	ctx.serviceRegistry = serviceRegistry

	serviceRegistry.On("GetAll", "Org1MSP", "Device1Id").Return([]*common.Service{{}, {}}, nil)

	contract := new(ServiceRegistrySmartContract)
	_, _ = contract.GetAll(ctx, "Org1MSP", "Device1Id")
	called := serviceRegistry.AssertCalled(s.T(), "GetAll", "Org1MSP", "Device1Id")
	assert.True(s.T(), called, "should retrieve services from service registry")
}

func (s *ServiceRegistryContractTestSuite) TestDeregister() {
	ctx := new(MockTransactionContext)
	serviceRegistry := new(MockServiceRegistry)
	ctx.serviceRegistry = serviceRegistry

	deviceId, _ := ctx.GetClientIdentity().GetID()
	organizationId, _ := ctx.GetClientIdentity().GetMSPID()
	expected := new(common.Service)
	serviceRegistry.On("Get", organizationId, deviceId, "Service1").Return(expected, nil)
	serviceRegistry.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, new(common.NotFoundError))
	serviceRegistry.On("Deregister", expected).Return(nil)

	contract := new(ServiceRegistrySmartContract)
	err := contract.Deregister(ctx, "Service1")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceRegistry.AssertCalled(s.T(), "Deregister", expected)
	assert.True(s.T(), called, "should remove service to service registry")
	actual := serviceRegistry.Calls[1].Arguments[0].(*common.Service)
	assert.Equal(s.T(), expected, actual, "should remove the correct service")

	ctx.clientId = &MockClientIdentity{Id: "Device2Id", MspId: "Org2MSP"}
	err = contract.Deregister(ctx, "Service2")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func TestServiceRegistryContractTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceRegistryContractTestSuite))
}
