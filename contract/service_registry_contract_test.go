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
	ctx := &MockTransactionContext{DeviceId: "device1", OrganizationId: "org1"}
	serviceRegistry := new(MockServiceRegistry)
	ctx.serviceRegistry = serviceRegistry

	serviceRegistry.On("Register", mock.AnythingOfType("*common.Service")).Return(nil)

	contract := new(ServiceRegistrySmartContract)
	err := contract.Register(ctx, fmt.Sprintf("{\"name\":\"service1\",\"Version\":1,\"description\":\"Service of Device1\",\"organizationId\":\"%s\",\"deviceId\":\"%s\",\"lastUpdateTime\":\"2021-12-12T17:36:00-05:00\"}", ctx.OrganizationId, ctx.DeviceId))
	assert.Nil(s.T(), err, "should return no error")
	called := serviceRegistry.AssertCalled(s.T(), "Register", mock.AnythingOfType("*common.Service"))
	assert.True(s.T(), called, "should put service to service registry")
	service, _ := common.DeserializeService(ctx.stub.EventPayload)
	assert.Equal(s.T(), fmt.Sprintf("service://%s/%s/%s/register", ctx.OrganizationId, ctx.DeviceId, "service1"), ctx.stub.EventName, "should emit event with name")
	assert.Equal(s.T(), "service1", service.Name, "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Register(ctx, "{\"name\":\"service2\",\"Version\":1,\"description\":\"Service of Device2\",\"organizationId\":\"org2\",\"deviceId\":\"device2\",\"lastUpdateTime\":\"2021-12-12T17:36:00-05:00\"}")
	assert.Error(s.T(), err, "should return mismatch device ID and organization ID error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")

	err = contract.Register(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")
}

func (s *ServiceRegistryContractTestSuite) TestGet() {
	ctx := &MockTransactionContext{DeviceId: "device2", OrganizationId: "org2"}
	serviceRegistry := new(MockServiceRegistry)
	ctx.serviceRegistry = serviceRegistry

	serviceRegistry.On("Get", "org1", "device1", "service1").Return(new(common.Service), nil)

	contract := new(ServiceRegistrySmartContract)
	_, _ = contract.Get(ctx, "org1", "device1", "service1")
	called := serviceRegistry.AssertCalled(s.T(), "Get", "org1", "device1", "service1")
	assert.True(s.T(), called, "should retrieve service from service registry")
}

func (s *ServiceRegistryContractTestSuite) TestGetAll() {
	ctx := &MockTransactionContext{DeviceId: "device2", OrganizationId: "org2"}
	serviceRegistry := new(MockServiceRegistry)
	ctx.serviceRegistry = serviceRegistry

	serviceRegistry.On("GetAll", "org1", "device1").Return([]*common.Service{{}, {}}, nil)

	contract := new(ServiceRegistrySmartContract)
	_, _ = contract.GetAll(ctx, "org1", "device1")
	called := serviceRegistry.AssertCalled(s.T(), "GetAll", "org1", "device1")
	assert.True(s.T(), called, "should retrieve services from service registry")
}

func (s *ServiceRegistryContractTestSuite) TestDeregister() {
	ctx := &MockTransactionContext{DeviceId: "device1", OrganizationId: "org1"}
	serviceRegistry := new(MockServiceRegistry)
	ctx.serviceRegistry = serviceRegistry

	serviceRegistry.On("Deregister", mock.MatchedBy(func(service *common.Service) bool {
		return service.DeviceId == "device1" && service.OrganizationId == "org1" && service.Name == "service1"
	})).Return(nil)
	serviceRegistry.On("Deregister", mock.Anything).Return(new(common.NotFoundError))

	contract := new(ServiceRegistrySmartContract)
	err := contract.Deregister(ctx, fmt.Sprintf("{\"name\":\"service1\",\"organizationId\":\"%s\",\"deviceId\":\"%s\"}", ctx.OrganizationId, ctx.DeviceId))
	assert.Nil(s.T(), err, "should return no error")
	actual := serviceRegistry.Calls[0].Arguments[0].(*common.Service)
	assert.Equal(s.T(), ctx.DeviceId, actual.DeviceId, "should remove the correct service")
	assert.Equal(s.T(), ctx.OrganizationId, actual.OrganizationId, "should remove the correct service")
	assert.Equal(s.T(), "service1", actual.Name, "should remove the correct service")
	service, _ := common.DeserializeService(ctx.stub.EventPayload)
	assert.Equal(s.T(), fmt.Sprintf("service://%s/%s/%s/deregister", ctx.OrganizationId, ctx.DeviceId, "service1"), ctx.stub.EventName, "should emit event with name")
	assert.Equal(s.T(), "service1", service.Name, "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Deregister(ctx, "{\"name\":\"service2\",\"organizationId\":\"org2\",\"deviceId\":\"device2\"}")
	assert.Error(s.T(), err, "should return mismatch device ID and organization ID error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")

	ctx.DeviceId = "device2"
	err = contract.Deregister(ctx, "{\"name\":\"service2\",\"organizationId\":\"org1\",\"deviceId\":\"device2\"}")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")
}

func TestServiceRegistryContractTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceRegistryContractTestSuite))
}
