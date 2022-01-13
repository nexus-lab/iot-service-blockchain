package common

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceRequestResponseTestSuite struct {
	suite.Suite
}

func (s *ServiceRequestResponseTestSuite) TestSerialize() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	request := &ServiceRequest{
		Id: "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Service: Service{
			Name:           "service1",
			DeviceId:       "device1",
			OrganizationId: "org1",
		},
		Method:    "GET",
		Arguments: []string{"1", "2", "3"},
		Time:      updateTime,
	}
	response := &ServiceResponse{
		RequestId:   "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Time:        updateTime,
		StatusCode:  0,
		ReturnValue: "[\"a\",\"b\",\"c\"]",
	}
	pair := &ServiceRequestResponse{Request: request, Response: response}
	serializedRequest := "{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\"," +
		"\"service\":{\"name\":\"service1\",\"deviceId\":\"device1\",\"organizationId\":\"org1\",\"version\":0," +
		"\"description\":\"\",\"lastUpdateTime\":\"0001-01-01T00:00:00Z\"},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]}"
	serializedResponse := "{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\"," +
		"\"statusCode\":0,\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}"
	serialized := fmt.Sprintf("{\"request\":%s,\"response\":%s}", serializedRequest, serializedResponse)

	data, err := pair.Serialize()
	assert.Equal(s.T(), serialized, string(data), "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceRequestResponseTestSuite) TestDeserializeService() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	request := &ServiceRequest{
		Id: "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Service: Service{
			Name:           "service1",
			DeviceId:       "device1",
			OrganizationId: "org1",
		},
		Method:    "GET",
		Arguments: []string{"1", "2", "3"},
		Time:      updateTime,
	}
	response := &ServiceResponse{
		RequestId:   "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Time:        updateTime,
		StatusCode:  0,
		ReturnValue: "[\"a\",\"b\",\"c\"]",
	}
	expected := &ServiceRequestResponse{Request: request, Response: response}
	serializedRequest := "{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\"," +
		"\"service\":{\"name\":\"service1\",\"deviceId\":\"device1\",\"organizationId\":\"org1\",\"version\":0," +
		"\"description\":\"\",\"lastUpdateTime\":\"0001-01-01T00:00:00Z\"},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]}"
	serializedResponse := "{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\"," +
		"\"statusCode\":0,\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}"
	serialized := fmt.Sprintf("{\"request\":%s,\"response\":%s}", serializedRequest, serializedResponse)

	actual, err := DeserializeServiceRequestResponse([]byte(serialized))
	assert.Equal(s.T(), expected, actual, "should return parsed service")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeServiceRequestResponse([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestServiceRequestResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceRequestResponseTestSuite))
}
