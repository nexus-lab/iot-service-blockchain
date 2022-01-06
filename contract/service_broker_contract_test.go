package contract

import (
	"fmt"
	"testing"

	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceBrokerContractTestSuite struct {
	suite.Suite
}

func (s *ServiceBrokerContractTestSuite) TestRequest() {
	ctx := &MockTransactionContext{DeviceId: "device1", OrganizationId: "org1"}
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("Request", mock.AnythingOfType("*common.ServiceRequest")).Return(nil)

	contract := new(ServiceBrokerSmartContract)
	err := contract.Request(ctx, fmt.Sprintf("{\"id\":\"request1\",\"time\":\"2021-12-12T17:38:00-05:00\",\"service\":{\"name\":\"service1\",\"organizationId\":\"%s\",\"deviceId\":\"%s\"},\"method\":\"GET\",\"arguments\":[]}", ctx.OrganizationId, ctx.DeviceId))
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Request", mock.AnythingOfType("*common.ServiceRequest"))
	assert.True(s.T(), called, "should put request to service broker")
	request := serviceBroker.Calls[0].Arguments[0].(*common.ServiceRequest)
	assert.Equal(s.T(), "request1", request.Id, "should have correct request ID")
	request, _ = common.DeserializeServiceRequest(ctx.stub.EventPayload)
	assert.Equal(s.T(), fmt.Sprintf("request://%s/%s/%s/%s/request", ctx.OrganizationId, ctx.DeviceId, "service1", "request1"), ctx.stub.EventName, "should emit event with name")
	assert.Equal(s.T(), "request1", request.Id, "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Request(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")
}

func (s *ServiceBrokerContractTestSuite) TestRespond() {
	ctx := &MockTransactionContext{DeviceId: "device1", OrganizationId: "org1"}
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("Respond", mock.AnythingOfType("*common.ServiceResponse")).Return(nil)

	pair := &common.ServiceRequestResponse{
		Request: &common.ServiceRequest{
			Id: "request1",
			Service: common.Service{
				Name:           "service1",
				DeviceId:       ctx.DeviceId,
				OrganizationId: ctx.OrganizationId,
			},
		},
		Response: new(common.ServiceResponse),
	}
	serviceBroker.On("Get", "request1").Return(pair, nil)
	serviceBroker.On("Get", mock.Anything).Return(nil, new(common.NotFoundError))

	contract := new(ServiceBrokerSmartContract)
	err := contract.Respond(ctx, "{\"requestId\":\"request1\",\"time\":\"2021-12-12T17:40:00-05:00\",\"statusCode\":0,\"returnValue\":\"[\\\"1.0\\\", \\\"2.0\\\", \\\"3.0\\\"]\"}")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Respond", mock.AnythingOfType("*common.ServiceResponse"))
	assert.True(s.T(), called, "should put response to service broker")
	response := serviceBroker.Calls[1].Arguments[0].(*common.ServiceResponse)
	assert.Equal(s.T(), "request1", response.RequestId, "should change device ID")
	response, _ = common.DeserializeServiceResponse(ctx.stub.EventPayload)
	assert.Equal(s.T(), fmt.Sprintf("request://%s/%s/%s/%s/respond", ctx.OrganizationId, ctx.DeviceId, "service1", "request1"), ctx.stub.EventName, "should emit event with name")
	assert.Equal(s.T(), "request1", response.RequestId, "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Respond(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")

	err = contract.Respond(ctx, "{\"requestId\":\"request2\"}")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")

	ctx.DeviceId = "device2"
	err = contract.Respond(ctx, "{\"requestId\":\"request1\"}")
	assert.Error(s.T(), err, "should refuse to respond for another device")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")
}

func (s *ServiceBrokerContractTestSuite) TestGet() {
	ctx := &MockTransactionContext{DeviceId: "device2", OrganizationId: "org2"}
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("Get", "request1").Return(new(common.ServiceRequestResponse), nil)

	contract := new(ServiceBrokerSmartContract)
	_, _ = contract.Get(ctx, "request1")
	called := serviceBroker.AssertCalled(s.T(), "Get", "request1")
	assert.True(s.T(), called, "should retrieve request & response from service broker")
}

func (s *ServiceBrokerContractTestSuite) TestGetAll() {
	ctx := &MockTransactionContext{DeviceId: "device2", OrganizationId: "org2"}
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("GetAll", "org1", "device1", "service1").Return([]*common.ServiceRequestResponse{{}, {}}, nil)

	contract := new(ServiceBrokerSmartContract)
	_, _ = contract.GetAll(ctx, "org1", "device1", "service1")
	called := serviceBroker.AssertCalled(s.T(), "GetAll", "org1", "device1", "service1")
	assert.True(s.T(), called, "should retrieve requests & responses from service broker")
}

func (s *ServiceBrokerContractTestSuite) TestRemove() {
	ctx := &MockTransactionContext{DeviceId: "device1", OrganizationId: "org1"}
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("Remove", mock.Anything).Return(nil)

	pair := &common.ServiceRequestResponse{
		Request: &common.ServiceRequest{
			Id: "request1",
			Service: common.Service{
				Name:           "service1",
				DeviceId:       ctx.DeviceId,
				OrganizationId: ctx.OrganizationId,
			},
		},
		Response: new(common.ServiceResponse),
	}
	serviceBroker.On("Get", "request1").Return(pair, nil)
	serviceBroker.On("Get", mock.Anything).Return(nil, new(common.NotFoundError))

	contract := new(ServiceBrokerSmartContract)
	err := contract.Remove(ctx, "request1")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Remove", "request1")
	assert.True(s.T(), called, "should remove request & response from service broker")
	requestId := serviceBroker.Calls[1].Arguments[0].(string)
	assert.Equal(s.T(), "request1", requestId, "should remove the correct device")
	assert.Equal(s.T(), fmt.Sprintf("request://%s/%s/%s/%s/remove", ctx.OrganizationId, ctx.DeviceId, "service1", "request1"), ctx.stub.EventName, "should emit event with name")
	assert.Equal(s.T(), "request1", string(ctx.stub.EventPayload), "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Remove(ctx, "request2")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")

	ctx.DeviceId = "device2"
	err = contract.Remove(ctx, "request1")
	assert.Error(s.T(), err, "should refuse to respond for another device")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")
}

func TestServiceBrokerContractTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceBrokerContractTestSuite))
}
