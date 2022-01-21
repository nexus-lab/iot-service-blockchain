package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceRequestTestSuite struct {
	suite.Suite
}

func (s *ServiceRequestTestSuite) TestGetKeyComponents() {
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
	assert.Equal(s.T(), []string{request.Id}, request.GetKeyComponents(), "should return correct key components")
}

func (s *ServiceRequestTestSuite) TestSerialize() {
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
	serialized := "{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\"," +
		"\"service\":{\"name\":\"service1\",\"deviceId\":\"device1\",\"organizationId\":\"org1\",\"version\":0," +
		"\"description\":\"\",\"lastUpdateTime\":\"0001-01-01T00:00:00Z\"},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]}"

	data, err := request.Serialize()
	assert.Equal(s.T(), serialized, string(data), "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceRequestTestSuite) TestValidate() {
	request := ServiceRequest{}

	request.Id = "123456"
	assert.Error(s.T(), request.Validate(), "should error on invalid ID")
	request.Id = "ffbc9005-c62a-4563-a8f7-b32bba27d707"

	assert.Error(s.T(), request.Validate(), "should error on empty serivce organization ID")
	request.Service.OrganizationId = "org1"

	assert.Error(s.T(), request.Validate(), "should error on empty service device ID")
	request.Service.DeviceId = "device1"

	assert.Error(s.T(), request.Validate(), "should error on empty service name")
	request.Service.Name = "service1"

	assert.Error(s.T(), request.Validate(), "should error on empty method")
	request.Method = "GET"

	assert.Error(s.T(), request.Validate(), "should error on empty last update time")
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	request.Time = updateTime

	assert.Nil(s.T(), request.Validate(), "should return no error")
}

func (s *ServiceRequestTestSuite) TestDeserializeService() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	expected := &ServiceRequest{
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
	serialized := "{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\"," +
		"\"service\":{\"name\":\"service1\",\"deviceId\":\"device1\",\"organizationId\":\"org1\",\"version\":0," +
		"\"description\":\"\",\"lastUpdateTime\":\"0001-01-01T00:00:00Z\"},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]}"

	actual, err := DeserializeServiceRequest([]byte(serialized))
	assert.Equal(s.T(), expected, actual, "should return parsed service")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeServiceRequest([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestServiceRequestTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceRequestTestSuite))
}
