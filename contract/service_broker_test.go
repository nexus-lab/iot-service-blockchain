package contract

import (
	"testing"

	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockServiceBroker struct {
	mock.Mock
}

func (r *MockServiceBroker) Request(request *common.ServiceRequest) error {
	args := r.Called(request)
	return args.Error(0)
}

func (r *MockServiceBroker) Respond(response *common.ServiceResponse) error {
	args := r.Called(response)
	return args.Error(0)
}

func (r *MockServiceBroker) Get(requestId string) (*common.ServiceRequestResponse, error) {
	args := r.Called(requestId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common.ServiceRequestResponse), args.Error(1)
}

func (r *MockServiceBroker) GetAll(organizationId string, deviceId string, serviceName string) ([]*common.ServiceRequestResponse, error) {
	args := r.Called(organizationId, deviceId, serviceName)
	return args.Get(0).([]*common.ServiceRequestResponse), args.Error(1)
}

func (r *MockServiceBroker) Remove(requestId string) error {
	args := r.Called(requestId)
	return args.Error(0)
}

type ServiceBrokerTestSuite struct {
	suite.Suite
}

func (s *ServiceBrokerTestSuite) TestRequest() {
	requestRegistry := new(MockStateRegistry)
	indexRegistry := new(MockStateRegistry)
	serviceRegistry := new(MockServiceRegistry)
	transactionContext := new(MockTransactionContext)

	transactionContext.serviceRegistry = serviceRegistry

	serviceBroker := new(ServiceBroker)
	serviceBroker.ctx = transactionContext
	serviceBroker.indexRegistry = indexRegistry
	serviceBroker.requestRegistry = requestRegistry

	request := &common.ServiceRequest{
		Id: "request1",
		Service: common.Service{
			OrganizationId: "org1",
			DeviceId:       "device1",
			Name:           "service1",
		},
	}

	requestRegistry.On("GetState", []string{"request1"}).Return(nil, new(common.NotFoundError))
	requestRegistry.On("GetState", mock.Anything).Return(new(common.ServiceRequest), nil)
	requestRegistry.On("PutState", request).Return(nil)
	indexRegistry.On("PutState", mock.Anything).Return(nil)
	serviceRegistry.On("Get", "org1", "device1", "service1").Return(new(common.Service), nil)
	serviceRegistry.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, new(common.NotFoundError))

	err := serviceBroker.Request(request)
	called := requestRegistry.AssertCalled(s.T(), "PutState", request)
	assert.True(s.T(), called, "should put request to state registry")
	assert.Nil(s.T(), err, "should return no error")
	called = indexRegistry.AssertCalled(s.T(), "PutState", mock.Anything)
	assert.True(s.T(), called, "should put index to state registry")
	index := indexRegistry.Calls[0].Arguments[0].(*serviceRequestIndex)
	assert.Equal(s.T(), request.Id, index.RequestId, "should set request ID of the index")

	request = &common.ServiceRequest{
		Id: "request1",
		Service: common.Service{
			OrganizationId: "org2",
			DeviceId:       "device2",
			Name:           "service2",
		},
	}
	err = serviceBroker.Request(request)
	notCalled := requestRegistry.AssertNotCalled(s.T(), "PutState", request)
	assert.True(s.T(), notCalled, "should not put request to state registry")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return service not found error")

	request = &common.ServiceRequest{
		Id: "request2",
		Service: common.Service{
			OrganizationId: "org1",
			DeviceId:       "device1",
			Name:           "service1",
		},
	}
	err = serviceBroker.Request(request)
	notCalled = requestRegistry.AssertNotCalled(s.T(), "PutState", request)
	assert.True(s.T(), notCalled, "should not put request to state registry")
	assert.EqualError(s.T(), err, "request already exists", "should return request already exists error")
}

func (s *ServiceBrokerTestSuite) TestRespond() {
	requestRegistry := new(MockStateRegistry)
	responseRegistry := new(MockStateRegistry)
	transactionContext := new(MockTransactionContext)

	serviceBroker := new(ServiceBroker)
	serviceBroker.ctx = transactionContext
	serviceBroker.requestRegistry = requestRegistry
	serviceBroker.responseRegistry = responseRegistry

	response := &common.ServiceResponse{RequestId: "request1"}

	requestRegistry.On("GetState", []string{"request1"}).Return(new(common.ServiceRequest), nil)
	requestRegistry.On("GetState", []string{"request2"}).Return(new(common.ServiceRequest), nil)
	requestRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))
	responseRegistry.On("GetState", []string{"request1"}).Return(nil, new(common.NotFoundError))
	responseRegistry.On("GetState", mock.Anything).Return(new(common.ServiceResponse), nil)
	responseRegistry.On("PutState", response).Return(nil)

	err := serviceBroker.Respond(response)
	called := responseRegistry.AssertCalled(s.T(), "PutState", response)
	assert.True(s.T(), called, "should put response to state registry")
	assert.Nil(s.T(), err, "should return no error")

	response = &common.ServiceResponse{RequestId: "request2"}
	err = serviceBroker.Respond(response)
	notCalled := responseRegistry.AssertNotCalled(s.T(), "PutState", response)
	assert.True(s.T(), notCalled, "should not put response to state registry")
	assert.EqualError(s.T(), err, "response already exists", "should return response already exists error")

	response = &common.ServiceResponse{RequestId: "request3"}
	err = serviceBroker.Respond(response)
	notCalled = responseRegistry.AssertNotCalled(s.T(), "PutState", response)
	assert.True(s.T(), notCalled, "should not put response to state registry")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return request not found error")
}

func (s *ServiceBrokerTestSuite) TestGet() {
	requestRegistry := new(MockStateRegistry)
	responseRegistry := new(MockStateRegistry)
	transactionContext := new(MockTransactionContext)

	serviceBroker := new(ServiceBroker)
	serviceBroker.ctx = transactionContext
	serviceBroker.requestRegistry = requestRegistry
	serviceBroker.responseRegistry = responseRegistry

	request1 := new(common.ServiceRequest)
	request2 := new(common.ServiceRequest)
	response1 := new(common.ServiceResponse)
	requestRegistry.On("GetState", []string{"request1"}).Return(request1, nil)
	requestRegistry.On("GetState", []string{"request2"}).Return(request2, nil)
	requestRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))
	responseRegistry.On("GetState", []string{"request1"}).Return(response1, nil)
	responseRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))

	result, err := serviceBroker.Get("request1")
	assert.Equal(s.T(), request1, result.Request, "should return the correct request")
	assert.Equal(s.T(), response1, result.Response, "should return the correct response")
	assert.Nil(s.T(), err, "should return no error")

	result, err = serviceBroker.Get("request2")
	assert.Equal(s.T(), request2, result.Request, "should return the correct request")
	assert.Nil(s.T(), result.Response, "should return no response")
	assert.Nil(s.T(), err, "should return no error")

	result, err = serviceBroker.Get("request3")
	assert.Nil(s.T(), result, "should return no request and response")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func (s *ServiceBrokerTestSuite) TestGetAll() {
	indexRegistry := new(MockStateRegistry)
	requestRegistry := new(MockStateRegistry)
	responseRegistry := new(MockStateRegistry)
	transactionContext := new(MockTransactionContext)

	serviceBroker := new(ServiceBroker)
	serviceBroker.ctx = transactionContext
	serviceBroker.indexRegistry = indexRegistry
	serviceBroker.requestRegistry = requestRegistry
	serviceBroker.responseRegistry = responseRegistry

	indices := []StateInterface{
		&serviceRequestIndex{RequestId: "request1"},
		&serviceRequestIndex{RequestId: "request2"},
	}
	request1 := new(common.ServiceRequest)
	request2 := new(common.ServiceRequest)
	response1 := new(common.ServiceResponse)

	indexRegistry.On("GetStates", []string{"org1", "device1", "service1"}).Return(indices, nil)
	indexRegistry.On("GetStates", mock.Anything).Return([]StateInterface{}, nil)
	requestRegistry.On("GetState", []string{"request1"}).Return(request1, nil)
	requestRegistry.On("GetState", []string{"request2"}).Return(request2, nil)
	requestRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))
	responseRegistry.On("GetState", []string{"request1"}).Return(response1, nil)
	responseRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))

	results, err := serviceBroker.GetAll("org1", "device1", "service1")
	assert.Equal(s.T(), len(indices), len(results), "should return the correct number of requests/responses")
	assert.Nil(s.T(), err, "should return no error")
	assert.Equal(s.T(), request1, results[0].Request, "should return the correct request")
	assert.Equal(s.T(), response1, results[0].Response, "should return the correct response")
	assert.Equal(s.T(), request2, results[1].Request, "should return the correct request")
	assert.Nil(s.T(), results[1].Response, "should return the correct response")

	results, err = serviceBroker.GetAll("org2", "device2", "service2")
	assert.Zero(s.T(), len(results), "should return no device")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceBrokerTestSuite) TestRemove() {
	indexRegistry := new(MockStateRegistry)
	requestRegistry := new(MockStateRegistry)
	responseRegistry := new(MockStateRegistry)
	transactionContext := new(MockTransactionContext)

	serviceBroker := new(ServiceBroker)
	serviceBroker.ctx = transactionContext
	serviceBroker.indexRegistry = indexRegistry
	serviceBroker.requestRegistry = requestRegistry
	serviceBroker.responseRegistry = responseRegistry

	request1 := new(common.ServiceRequest)
	request2 := new(common.ServiceRequest)
	response1 := new(common.ServiceResponse)

	requestRegistry.On("GetState", []string{"request1"}).Return(request1, nil)
	requestRegistry.On("GetState", []string{"request2"}).Return(request2, nil)
	requestRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))
	responseRegistry.On("GetState", []string{"request1"}).Return(response1, nil)
	responseRegistry.On("GetState", mock.Anything).Return(nil, new(common.NotFoundError))
	requestRegistry.On("RemoveState", mock.Anything).Return(nil)
	responseRegistry.On("RemoveState", mock.Anything).Return(nil)
	indexRegistry.On("RemoveState", mock.Anything).Return(nil)

	err := serviceBroker.Remove("request1")
	called := requestRegistry.AssertCalled(s.T(), "RemoveState", request1)
	assert.True(s.T(), called, "should remove request from state registry")
	called = responseRegistry.AssertCalled(s.T(), "RemoveState", response1)
	assert.True(s.T(), called, "should remove response from state registry")
	called = indexRegistry.AssertCalled(s.T(), "RemoveState", mock.Anything)
	assert.True(s.T(), called, "should remove index from state registry")
	index := indexRegistry.Calls[0].Arguments[0].(*serviceRequestIndex)
	assert.Equal(s.T(), request1.Id, index.RequestId, "should remove correct index from state registry")
	assert.Nil(s.T(), err, "should return no error")

	err = serviceBroker.Remove("request2")
	called = requestRegistry.AssertCalled(s.T(), "RemoveState", request2)
	assert.True(s.T(), called, "should remove request from state registry")
	notCalled := responseRegistry.AssertNumberOfCalls(s.T(), "RemoveState", 1)
	assert.True(s.T(), notCalled, "should not remove response from state registry")
	called = indexRegistry.AssertNumberOfCalls(s.T(), "RemoveState", 2)
	assert.True(s.T(), called, "should remove index from state registry")
	index = indexRegistry.Calls[1].Arguments[0].(*serviceRequestIndex)
	assert.Equal(s.T(), request2.Id, index.RequestId, "should remove correct index from state registry")
	assert.Nil(s.T(), err, "should return no error")

	err = serviceBroker.Remove("request3")
	notCalled = requestRegistry.AssertNumberOfCalls(s.T(), "RemoveState", 2)
	assert.True(s.T(), notCalled, "should not remove request from state registry")
	notCalled = responseRegistry.AssertNumberOfCalls(s.T(), "RemoveState", 1)
	assert.True(s.T(), notCalled, "should not remove response from state registry")
	notCalled = indexRegistry.AssertNumberOfCalls(s.T(), "RemoveState", 2)
	assert.True(s.T(), notCalled, "should not remove index from state registry")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func TestServiceBrokerTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceBrokerTestSuite))
}
