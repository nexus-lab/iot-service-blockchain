package contract

import (
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

	serviceBroker.On("Request", mock.AnythingOfType("*common.ServiceRequest")).Return(nil)

	contract := new(ServiceBrokerSmartContract)
	err := contract.Request(ctx, "{\"id\":\"Request1Id\",\"time\":\"2021-12-12T17:38:00-05:00\",\"service\":{\"name\":\"Service1\",\"organizationId\":\"Org1MSP\",\"deviceId\":\"Device1Id\"},\"method\":\"GET\",\"arguments\":[]}")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Request", mock.AnythingOfType("*common.ServiceRequest"))
	assert.True(s.T(), called, "should put request to service broker")
	request := serviceBroker.Calls[0].Arguments[0].(*common.ServiceRequest)
	assert.Equal(s.T(), "Request1Id", request.Id, "should have correct request ID")

	err = contract.Request(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
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
			Service: common.Service{
				DeviceId:       deviceId,
				OrganizationId: organizationId,
			},
		},
		Response: new(common.ServiceResponse),
	}
	serviceBroker.On("Get", "Request1Id").Return(pair, nil)
	serviceBroker.On("Get", mock.Anything).Return(nil, new(common.NotFoundError))

	contract := new(ServiceBrokerSmartContract)
	err := contract.Respond(ctx, "{\"requestId\":\"Request1Id\",\"time\":\"2021-12-12T17:40:00-05:00\",\"statusCode\":0,\"returnValue\":\"[\\\"1.0\\\", \\\"2.0\\\", \\\"3.0\\\"]\"}")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Respond", mock.AnythingOfType("*common.ServiceResponse"))
	assert.True(s.T(), called, "should put response to service broker")
	response := serviceBroker.Calls[1].Arguments[0].(*common.ServiceResponse)
	assert.Equal(s.T(), "Request1Id", response.RequestId, "should change device ID")

	err = contract.Respond(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")

	err = contract.Respond(ctx, "{\"requestId\":\"Request2Id\"}")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")

	ctx.clientId = &MockClientIdentity{Id: "Device2Id", MspId: "Org2MSP"}
	err = contract.Respond(ctx, "{\"requestId\":\"Request1Id\"}")
	assert.Error(s.T(), err, "should refuse to respond for another device")
}

func (s *ServiceBrokerContractTestSuite) TestGet() {
	ctx := new(MockTransactionContext)
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("Get", "Request1Id").Return(new(common.ServiceRequestResponse), nil)

	contract := new(ServiceBrokerSmartContract)
	_, _ = contract.Get(ctx, "Request1Id")
	called := serviceBroker.AssertCalled(s.T(), "Get", "Request1Id")
	assert.True(s.T(), called, "should retrieve request & response from service broker")
}

func (s *ServiceBrokerContractTestSuite) TestGetAll() {
	ctx := new(MockTransactionContext)
	serviceBroker := new(MockServiceBroker)
	ctx.serviceBroker = serviceBroker

	serviceBroker.On("GetAll", "Org1MSP", "Device1Id", "Service1").Return([]*common.ServiceRequestResponse{{}, {}}, nil)

	contract := new(ServiceBrokerSmartContract)
	_, _ = contract.GetAll(ctx, "Org1MSP", "Device1Id", "Service1")
	called := serviceBroker.AssertCalled(s.T(), "GetAll", "Org1MSP", "Device1Id", "Service1")
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
			Service: common.Service{
				DeviceId:       deviceId,
				OrganizationId: organizationId,
			},
		},
		Response: new(common.ServiceResponse),
	}
	serviceBroker.On("Get", "Request1Id").Return(pair, nil)
	serviceBroker.On("Get", mock.Anything).Return(nil, new(common.NotFoundError))

	contract := new(ServiceBrokerSmartContract)
	err := contract.Remove(ctx, "Request1Id")
	assert.Nil(s.T(), err, "should return no error")
	called := serviceBroker.AssertCalled(s.T(), "Remove", "Request1Id")
	assert.True(s.T(), called, "should remove request & response from service broker")
	requestId := serviceBroker.Calls[1].Arguments[0].(string)
	assert.Equal(s.T(), "Request1Id", requestId, "should remove the correct device")

	err = contract.Remove(ctx, "Request2Id")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")

	ctx.clientId = &MockClientIdentity{Id: "Device2Id", MspId: "Org2MSP"}
	err = contract.Remove(ctx, "Request1Id")
	assert.Error(s.T(), err, "should refuse to respond for another device")
}

func TestServiceBrokerContractTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceBrokerContractTestSuite))
}
