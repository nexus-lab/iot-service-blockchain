package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceResponseTestSuite struct {
	suite.Suite
	response *ServiceResponse
	json     string
}

func (s *ServiceResponseTestSuite) SetupTest() {
	time, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")

	s.json = "{\"requestId\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\",\"statusCode\":0,\"returnValue\":\"[\\\"a\\\",\\\"b\\\",\\\"c\\\"]\"}"
	s.response = &ServiceResponse{
		RequestId:   "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Time:        time,
		StatusCode:  0,
		ReturnValue: "[\"a\",\"b\",\"c\"]",
	}
}

func (s *ServiceResponseTestSuite) TestGetKeyComponents() {
	assert.Equal(s.T(), []string{s.response.RequestId}, s.response.GetKeyComponents(), "should return correct key components")
}

func (s *ServiceResponseTestSuite) TestSerialize() {
	data, err := s.response.Serialize()
	assert.Equal(s.T(), string(data), s.json, "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceResponseTestSuite) TestValidate() {
	response := ServiceResponse{}

	response.RequestId = "123456"
	assert.Error(s.T(), response.Validate(), "should error on invalid ID")
	response.RequestId = s.response.RequestId

	assert.Error(s.T(), response.Validate(), "should error on empty last update time")
	response.Time = s.response.Time

	assert.Nil(s.T(), response.Validate(), "should return no error")
}

func (s *ServiceResponseTestSuite) TestDeserializeService() {
	data := s.json
	response, err := DeserializeServiceResponse([]byte(data))
	assert.Equal(s.T(), response, s.response, "should return parsed service")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeServiceResponse([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestServiceResponseTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceResponseTestSuite))
}
