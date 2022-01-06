package contract

import (
	"fmt"
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
	err := contract.Register(ctx, fmt.Sprintf("{\"name\":\"Service1\",\"Version\":1,\"description\":\"Service of Device1\",\"organizationId\":\"%s\",\"deviceId\":\"%s\",\"lastUpdateTime\":\"2021-12-12T17:36:00-05:00\"}", organizationId, deviceId))
	assert.Nil(s.T(), err, "should return no error")
	called := serviceRegistry.AssertCalled(s.T(), "Register", mock.AnythingOfType("*common.Service"))
	assert.True(s.T(), called, "should put service to service registry")

	err = contract.Register(ctx, "{\"name\":\"Service2\",\"Version\":1,\"description\":\"Service of Device2\",\"organizationId\":\"Org2MSP\",\"deviceId\":\"Device2Id\",\"lastUpdateTime\":\"2021-12-12T17:36:00-05:00\"}")
	assert.Error(s.T(), err, "should return mismatch device ID and organization ID error")

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
	serviceRegistry.On("Deregister", mock.MatchedBy(func(service *common.Service) bool {
		return service.DeviceId == deviceId && service.OrganizationId == organizationId && service.Name == "Service1"
	})).Return(nil)
	serviceRegistry.On("Deregister", mock.Anything).Return(new(common.NotFoundError))

	contract := new(ServiceRegistrySmartContract)
	err := contract.Deregister(ctx, fmt.Sprintf("{\"name\":\"Service1\",\"organizationId\":\"%s\",\"deviceId\":\"%s\"}", organizationId, deviceId))
	assert.Nil(s.T(), err, "should return no error")
	actual := serviceRegistry.Calls[0].Arguments[0].(*common.Service)
	assert.Equal(s.T(), deviceId, actual.DeviceId, "should remove the correct service")
	assert.Equal(s.T(), organizationId, actual.OrganizationId, "should remove the correct service")
	assert.Equal(s.T(), "Service1", actual.Name, "should remove the correct service")

	err = contract.Deregister(ctx, "{\"name\":\"Service2\",\"organizationId\":\"Org2MSP\",\"deviceId\":\"Device2Id\"}")
	assert.Error(s.T(), err, "should return mismatch device ID and organization ID error")

	ctx.clientId = &MockClientIdentity{Id: "Device2Id", MspId: "Org2MSP"}
	err = contract.Deregister(ctx, "{\"name\":\"Service2\",\"organizationId\":\"Org2MSP\",\"deviceId\":\"Device2Id\"}")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func TestServiceRegistryContractTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceRegistryContractTestSuite))
}
