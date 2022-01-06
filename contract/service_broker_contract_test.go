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
	ctx := new(MockTransactionContext)
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	deviceId, _ := ctx.GetClientIdentity().GetID()
	organizationId, _ := ctx.GetClientIdentity().GetMSPID()
	serviceBroker.On("Request", mock.AnythingOfType("*common.ServiceRequest")).Return(nil)

	contract := new(ServiceBrokerSmartContract)
	err := contract.Request(ctx, fmt.Sprintf("{\"id\":\"Request1\",\"time\":\"2021-12-12T17:38:00-05:00\",\"service\":{\"name\":\"Service1\",\"organizationId\":\"%s\",\"deviceId\":\"%s\"},\"method\":\"GET\",\"arguments\":[]}", organizationId, deviceId))
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Request", mock.AnythingOfType("*common.ServiceRequest"))
	assert.True(s.T(), called, "should put request to service broker")
	request := serviceBroker.Calls[0].Arguments[0].(*common.ServiceRequest)
	assert.Equal(s.T(), "Request1", request.Id, "should have correct request ID")
	request, _ = common.DeserializeServiceRequest(ctx.stub.eventPayload)
	assert.Equal(s.T(), fmt.Sprintf("request://%s/%s/%s/%s/request", organizationId, deviceId, "Service1", "Request1"), ctx.stub.eventName, "should emit event with name")
	assert.Equal(s.T(), "Request1", request.Id, "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Request(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
	assert.Empty(s.T(), ctx.stub.eventName, "should not emit event")
}

func (s *ServiceBrokerContractTestSuite) TestRespond() {
	ctx := new(MockTransactionContext)
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	deviceId, _ := ctx.GetClientIdentity().GetID()
	organizationId, _ := ctx.GetClientIdentity().GetMSPID()
	serviceBroker.On("Respond", mock.AnythingOfType("*common.ServiceResponse")).Return(nil)

	pair := &common.ServiceRequestResponse{
		Request: &common.ServiceRequest{
			Id: "Request1",
			Service: common.Service{
				Name:           "Service1",
				DeviceId:       deviceId,
				OrganizationId: organizationId,
			},
		},
		Response: new(common.ServiceResponse),
	}
	serviceBroker.On("Get", "Request1").Return(pair, nil)
	serviceBroker.On("Get", mock.Anything).Return(nil, new(common.NotFoundError))

	contract := new(ServiceBrokerSmartContract)
	err := contract.Respond(ctx, "{\"requestId\":\"Request1\",\"time\":\"2021-12-12T17:40:00-05:00\",\"statusCode\":0,\"returnValue\":\"[\\\"1.0\\\", \\\"2.0\\\", \\\"3.0\\\"]\"}")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Respond", mock.AnythingOfType("*common.ServiceResponse"))
	assert.True(s.T(), called, "should put response to service broker")
	response := serviceBroker.Calls[1].Arguments[0].(*common.ServiceResponse)
	assert.Equal(s.T(), "Request1", response.RequestId, "should change device ID")
	response, _ = common.DeserializeServiceResponse(ctx.stub.eventPayload)
	assert.Equal(s.T(), fmt.Sprintf("request://%s/%s/%s/%s/respond", organizationId, deviceId, "Service1", "Request1"), ctx.stub.eventName, "should emit event with name")
	assert.Equal(s.T(), "Request1", response.RequestId, "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Respond(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
	assert.Empty(s.T(), ctx.stub.eventName, "should not emit event")

	err = contract.Respond(ctx, "{\"requestId\":\"Request2\"}")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
	assert.Empty(s.T(), ctx.stub.eventName, "should not emit event")

	ctx.clientId = &MockClientIdentity{Id: "Device2", MspId: "Org2MSP"}
	err = contract.Respond(ctx, "{\"requestId\":\"Request1\"}")
	assert.Error(s.T(), err, "should refuse to respond for another device")
	assert.Empty(s.T(), ctx.stub.eventName, "should not emit event")
}

func (s *ServiceBrokerContractTestSuite) TestGet() {
	ctx := new(MockTransactionContext)
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("Get", "Request1").Return(new(common.ServiceRequestResponse), nil)

	contract := new(ServiceBrokerSmartContract)
	_, _ = contract.Get(ctx, "Request1")
	called := serviceBroker.AssertCalled(s.T(), "Get", "Request1")
	assert.True(s.T(), called, "should retrieve request & response from service broker")
}

func (s *ServiceBrokerContractTestSuite) TestGetAll() {
	ctx := new(MockTransactionContext)
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("GetAll", "Org1MSP", "Device1", "Service1").Return([]*common.ServiceRequestResponse{{}, {}}, nil)

	contract := new(ServiceBrokerSmartContract)
	_, _ = contract.GetAll(ctx, "Org1MSP", "Device1", "Service1")
	called := serviceBroker.AssertCalled(s.T(), "GetAll", "Org1MSP", "Device1", "Service1")
	assert.True(s.T(), called, "should retrieve requests & responses from service broker")
}

func (s *ServiceBrokerContractTestSuite) TestRemove() {
	ctx := new(MockTransactionContext)
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	deviceId, _ := ctx.GetClientIdentity().GetID()
	organizationId, _ := ctx.GetClientIdentity().GetMSPID()
	serviceBroker.On("Remove", mock.Anything).Return(nil)

	pair := &common.ServiceRequestResponse{
		Request: &common.ServiceRequest{
			Id: "Request1",
			Service: common.Service{
				Name:           "Service1",
				DeviceId:       deviceId,
				OrganizationId: organizationId,
			},
		},
		Response: new(common.ServiceResponse),
	}
	serviceBroker.On("Get", "Request1").Return(pair, nil)
	serviceBroker.On("Get", mock.Anything).Return(nil, new(common.NotFoundError))

	contract := new(ServiceBrokerSmartContract)
	err := contract.Remove(ctx, "Request1")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Remove", "Request1")
	assert.True(s.T(), called, "should remove request & response from service broker")
	requestId := serviceBroker.Calls[1].Arguments[0].(string)
	assert.Equal(s.T(), "Request1", requestId, "should remove the correct device")
	assert.Equal(s.T(), fmt.Sprintf("request://%s/%s/%s/%s/remove", organizationId, deviceId, "Service1", "Request1"), ctx.stub.eventName, "should emit event with name")
	assert.Equal(s.T(), "Request1", string(ctx.stub.eventPayload), "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Remove(ctx, "Request2")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
	assert.Empty(s.T(), ctx.stub.eventName, "should not emit event")

	ctx.clientId = &MockClientIdentity{Id: "Device2", MspId: "Org2MSP"}
	err = contract.Remove(ctx, "Request1")
	assert.Error(s.T(), err, "should refuse to respond for another device")
	assert.Empty(s.T(), ctx.stub.eventName, "should not emit event")
}

func TestServiceBrokerContractTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceBrokerContractTestSuite))
}
