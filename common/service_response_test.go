package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceResponseTestSuite struct {
	suite.Suite
}

func (s *ServiceResponseTestSuite) TestGetKeyComponents() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	response := &ServiceResponse{
		RequestId:   "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Time:        updateTime,
		StatusCode:  0,
		ReturnValue: "[\"a\",\"b\",\"c\"]",
	}
	assert.Equal(s.T(), []string{response.RequestId}, response.GetKeyComponents(), "should return correct key components")
}

func (s *ServiceResponseTestSuite) TestSerialize() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	response := &ServiceResponse{
		RequestId:   "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Time:        updateTime,
		StatusCode:  0,
		ReturnValue: "[\"a\",\"b\",\"c\"]",
	}
	serialized := "{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\"," +
		"\"statusCode\":0,\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}"

	data, err := response.Serialize()
	assert.Equal(s.T(), serialized, string(data), "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceResponseTestSuite) TestValidate() {
	response := ServiceResponse{}

	response.RequestId = "123456"
	assert.Error(s.T(), response.Validate(), "should error on invalid ID")
	response.RequestId = "ffbc9005-c62a-4563-a8f7-b32bba27d707"

	assert.Error(s.T(), response.Validate(), "should error on empty last update time")
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	response.Time = updateTime

	assert.Nil(s.T(), response.Validate(), "should return no error")
}

func (s *ServiceResponseTestSuite) TestDeserializeService() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	expected := &ServiceResponse{
		RequestId:   "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Time:        updateTime,
		StatusCode:  0,
		ReturnValue: "[\"a\",\"b\",\"c\"]",
	}
	serialized := "{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\"," +
		"\"statusCode\":0,\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}"

	actual, err := DeserializeServiceResponse([]byte(serialized))
	assert.Equal(s.T(), expected, actual, "should return parsed service")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeServiceResponse([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestServiceResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceResponseTestSuite))
}
