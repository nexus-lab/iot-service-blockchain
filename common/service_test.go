package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (s *ServiceTestSuite) TestGetKeyComponents() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	service := &Service{
		Name:           "service1",
		DeviceId:       "device1",
		OrganizationId: "org1",
		Version:        1,
		Description:    "Service of Device1",
		LastUpdateTime: updateTime,
	}
	assert.Equal(s.T(), []string{service.OrganizationId, service.DeviceId, service.Name}, service.GetKeyComponents(), "should return correct key components")
}

func (s *ServiceTestSuite) TestSerialize() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	service := &Service{
		Name:           "service1",
		DeviceId:       "device1",
		OrganizationId: "org1",
		Version:        1,
		Description:    "Service of Device1",
		LastUpdateTime: updateTime,
	}
	serialized := "{\"name\":\"service1\",\"deviceId\":\"device1\",\"organizationId\":\"org1\",\"version\":1,\"description\":\"Service of Device1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}"

	data, err := service.Serialize()
	assert.Equal(s.T(), serialized, string(data), "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceTestSuite) TestValidate() {
	service := Service{}

	assert.Error(s.T(), service.Validate(), "should error on empty name")
	service.Name = "service1"

	assert.Error(s.T(), service.Validate(), "should error on empty device ID")
	service.DeviceId = "device1"

	assert.Error(s.T(), service.Validate(), "should error on empty organization ID")
	service.OrganizationId = "org1"

	assert.Error(s.T(), service.Validate(), "should error on empty version")
	service.Version = 1

	assert.Error(s.T(), service.Validate(), "should error on empty last update time")
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	service.LastUpdateTime = updateTime

	assert.Nil(s.T(), service.Validate(), "should return no error")
}

func (s *ServiceTestSuite) TestDeserializeService() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	expected := &Service{
		Name:           "service1",
		DeviceId:       "device1",
		OrganizationId: "org1",
		Version:        1,
		Description:    "Service of Device1",
		LastUpdateTime: updateTime,
	}
	serialized := "{\"name\":\"service1\",\"deviceId\":\"device1\",\"organizationId\":\"org1\",\"version\":1,\"description\":\"Service of Device1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}"

	actual, err := DeserializeService([]byte(serialized))
	assert.Equal(s.T(), expected, actual, "should return parsed service")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeService([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
