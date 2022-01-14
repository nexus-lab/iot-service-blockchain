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

type ServiceBrokerTestSuite struct {
	suite.Suite
}

func (s *ServiceBrokerTestSuite) TestRequest() {
	contract := new(MockContract)
	serviceBroker := &ServiceBroker{contract}

	request := &common.ServiceRequest{Id: "request1"}
	data, _ := request.Serialize()
	contract.On("SubmitTransaction", "Request", string(data)).Return(nil, nil)

	err := serviceBroker.Request(request)
	assert.Nil(s.T(), err, "should return no error")

	err = serviceBroker.Request(nil)
	assert.Error(s.T(), err, "should return error if input is null")

	request = &common.ServiceRequest{Id: "request2"}
	data, _ = request.Serialize()
	contract.On("SubmitTransaction", "Request", string(data)).Return(nil, errors.New(""))

	err = serviceBroker.Request(request)
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceBrokerTestSuite) TestRespond() {
	contract := new(MockContract)
	serviceBroker := &ServiceBroker{contract}

	response := &common.ServiceResponse{RequestId: "request1"}
	data, _ := response.Serialize()
	contract.On("SubmitTransaction", "Respond", string(data)).Return(nil, nil)

	err := serviceBroker.Respond(response)
	assert.Nil(s.T(), err, "should return no error")

	err = serviceBroker.Respond(nil)
	assert.Error(s.T(), err, "should return error if input is null")

	response = &common.ServiceResponse{RequestId: "request2"}
	data, _ = response.Serialize()
	contract.On("SubmitTransaction", "Respond", string(data)).Return(nil, errors.New(""))

	err = serviceBroker.Respond(response)
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceBrokerTestSuite) TestGet() {
	contract := new(MockContract)
	serviceBroker := &ServiceBroker{contract}

	expected := &common.ServiceRequestResponse{}
	data, _ := expected.Serialize()
	contract.On("SubmitTransaction", "Get", "request1").Return(data, nil)

	actual, err := serviceBroker.Get("request1")
	assert.Equal(s.T(), expected, actual, "should return correct request & response")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "Get", "request2").Return(nil, new(common.NotFoundError))

	actual, err = serviceBroker.Get("request2")
	assert.Nil(s.T(), actual, "should return no request / response")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")

	contract.On("SubmitTransaction", "Get", "request3").Return(nil, errors.New(""))

	_, err = serviceBroker.Get("request3")
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceBrokerTestSuite) TestGetAll() {
	contract := new(MockContract)
	serviceBroker := &ServiceBroker{contract}

	expected := []*common.ServiceRequestResponse{new(common.ServiceRequestResponse), new(common.ServiceRequestResponse)}
	data, _ := json.Marshal(expected)
	contract.On("SubmitTransaction", "GetAll", "org1", "device1", "service1").Return(data, nil)

	actual, err := serviceBroker.GetAll("org1", "device1", "service1")
	assert.Equal(s.T(), expected, actual, "should return correct requests & responses")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "GetAll", "org2", "device2", "service2").Return([]byte("[]"), nil)

	actual, err = serviceBroker.GetAll("org2", "device2", "service2")
	assert.Zero(s.T(), len(actual), "should return no request / response")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "GetAll", "org3", "device3", "service3").Return(nil, errors.New(""))

	_, err = serviceBroker.GetAll("org3", "device3", "service3")
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceBrokerTestSuite) TestRemove() {
	contract := new(MockContract)
	serviceBroker := &ServiceBroker{contract}

	contract.On("SubmitTransaction", "Remove", "request1").Return(nil, nil)

	err := serviceBroker.Remove("request1")
	assert.Nil(s.T(), err, "should return no error")

	contract.On("SubmitTransaction", "Remove", "request2").Return(nil, errors.New(""))

	err = serviceBroker.Remove("request2")
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func (s *ServiceBrokerTestSuite) TestRegisterEvent() {
	contract := new(MockContract)
	serviceBroker := &ServiceBroker{contract}

	eventChannel := make(chan *client.ChaincodeEvent)
	go func() {
		for i := 0; i < 2; i++ {
			data, _ := (&common.ServiceRequest{Id: fmt.Sprintf("request%d", i)}).Serialize()
			eventChannel <- &client.ChaincodeEvent{
				EventName: fmt.Sprintf("request://org%d/device%d/service%d/request%d/request", i, i, i, i),
				Payload:   data,
			}
		}

		for i := 0; i < 2; i++ {
			data, _ := (&common.ServiceResponse{RequestId: fmt.Sprintf("request%d", i)}).Serialize()
			eventChannel <- &client.ChaincodeEvent{
				EventName: fmt.Sprintf("request://org%d/device%d/service%d/request%d/respond", i, i, i, i),
				Payload:   data,
			}
		}

		for i := 0; i < 2; i++ {
			eventChannel <- &client.ChaincodeEvent{
				EventName: fmt.Sprintf("request://org%d/device%d/service%d/request%d/remove", i, i, i, i),
				Payload:   []byte(fmt.Sprintf("request%d", i)),
			}
		}
	}()

	var cancelFunc context.CancelFunc = func() {
		close(eventChannel)
	}

	contract.On("RegisterEvent", mock.Anything).Return(eventChannel, cancelFunc, nil)

	source, cancel, err := serviceBroker.RegisterEvent()
	defer cancel()
	assert.Nil(s.T(), err, "should return no error")
	assert.IsType(s.T(), *new(context.CancelFunc), cancel, "should return correct cancel function")

	for i := 0; i < 2; i++ {
		event := <-source
		assert.Equal(s.T(), "request", event.Action, "should return correct action")
		assert.Equal(s.T(), fmt.Sprintf("org%d", i), event.OrganizationId, "should return correct organization ID")
		assert.Equal(s.T(), fmt.Sprintf("device%d", i), event.DeviceId, "should return correct device ID")
		assert.Equal(s.T(), fmt.Sprintf("service%d", i), event.ServiceName, "should return correct service name")
		assert.Equal(s.T(), fmt.Sprintf("request%d", i), event.RequestId, "should return correct request ID")
		assert.IsType(s.T(), new(common.ServiceRequest), event.Payload, "should return parsed service request as event payload")
		assert.Equal(s.T(), fmt.Sprintf("request%d", i), event.Payload.(*common.ServiceRequest).Id, "should return correct event payload")
	}

	for i := 0; i < 2; i++ {
		event := <-source
		assert.Equal(s.T(), "respond", event.Action, "should return correct action")
		assert.Equal(s.T(), fmt.Sprintf("org%d", i), event.OrganizationId, "should return correct organization ID")
		assert.Equal(s.T(), fmt.Sprintf("device%d", i), event.DeviceId, "should return correct device ID")
		assert.Equal(s.T(), fmt.Sprintf("service%d", i), event.ServiceName, "should return correct service name")
		assert.Equal(s.T(), fmt.Sprintf("request%d", i), event.RequestId, "should return correct request ID")
		assert.IsType(s.T(), new(common.ServiceResponse), event.Payload, "should return parsed service request as event payload")
		assert.Equal(s.T(), fmt.Sprintf("request%d", i), event.Payload.(*common.ServiceResponse).RequestId, "should return correct event payload")
	}

	for i := 0; i < 2; i++ {
		event := <-source
		assert.Equal(s.T(), "remove", event.Action, "should return correct action")
		assert.Equal(s.T(), fmt.Sprintf("org%d", i), event.OrganizationId, "should return correct organization ID")
		assert.Equal(s.T(), fmt.Sprintf("device%d", i), event.DeviceId, "should return correct device ID")
		assert.Equal(s.T(), fmt.Sprintf("service%d", i), event.ServiceName, "should return correct service name")
		assert.Equal(s.T(), fmt.Sprintf("request%d", i), event.RequestId, "should return correct request ID")
		assert.Equal(s.T(), fmt.Sprintf("request%d", i), event.Payload, "should return correct event payload")
	}

	contract = new(MockContract)
	serviceBroker = &ServiceBroker{contract}
	contract.On("RegisterEvent", mock.Anything).Return(nil, nil, errors.New(""))

	_, _, err = serviceBroker.RegisterEvent()
	assert.Error(s.T(), err, "should return error when sdk or smart contract fails")
}

func TestServiceBrokerTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceBrokerTestSuite))
}
