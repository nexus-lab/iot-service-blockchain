package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeviceTestSuite struct {
	suite.Suite
}

func (s *DeviceTestSuite) TestGetKeyComponents() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	device := &Device{
		Id:             "device1",
		OrganizationId: "org1",
		Name:           "device1",
		Description:    "Device of Org1 User1",
		LastUpdateTime: updateTime,
	}
	assert.Equal(s.T(), []string{device.OrganizationId, device.Id}, device.GetKeyComponents(), "should return correct key components")
}

func (s *DeviceTestSuite) TestSerialize() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	device := &Device{
		Id:             "device1",
		OrganizationId: "org1",
		Name:           "device1",
		Description:    "Device of Org1 User1",
		LastUpdateTime: updateTime,
	}
	serialized := "{\"id\":\"device1\",\"organizationId\":\"org1\",\"name\":\"device1\",\"description\":\"Device of Org1 User1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}"

	data, err := device.Serialize()
	assert.Equal(s.T(), string(data), serialized, "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *DeviceTestSuite) TestValidate() {
	device := Device{}

	assert.Error(s.T(), device.Validate(), "should error on empty ID")
	device.Id = "device1"

	assert.Error(s.T(), device.Validate(), "should error on empty organization ID")
	device.OrganizationId = "org1"

	assert.Error(s.T(), device.Validate(), "should error on empty name")
	device.Name = "device1"

	assert.Error(s.T(), device.Validate(), "should error on empty last update time")
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	device.LastUpdateTime = updateTime

	assert.Nil(s.T(), device.Validate(), "should return no error")
}

func (s *DeviceTestSuite) TestDeserializeDevice() {
	updateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")
	expected := &Device{
		Id:             "device1",
		OrganizationId: "org1",
		Name:           "device1",
		Description:    "Device of Org1 User1",
		LastUpdateTime: updateTime,
	}
	serialized := "{\"id\":\"device1\",\"organizationId\":\"org1\",\"name\":\"device1\",\"description\":\"Device of Org1 User1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}"

	actual, err := DeserializeDevice([]byte(serialized))
	assert.Equal(s.T(), expected, actual, "should return parsed device")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeDevice([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestDeviceTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}
