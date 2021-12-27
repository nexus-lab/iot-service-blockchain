package common

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeviceTestSuite struct {
	suite.Suite
	device *Device
	json   string
}

func (s *DeviceTestSuite) SetupTest() {
	lastUpdateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")

	id := "eDUwOTo6Q049dXNlcjEsT1U9Y2xpZW50LE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUzo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw=="
	s.json = fmt.Sprintf("{\"id\":\"%s\",\"organizationId\":\"Org1MSP\",\"name\":\"device1\",\"description\":\"Device of Org1 User1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}", id)
	s.device = &Device{
		Id:             id,
		OrganizationId: "Org1MSP",
		Name:           "device1",
		Description:    "Device of Org1 User1",
		LastUpdateTime: lastUpdateTime,
	}
}

func (s *DeviceTestSuite) TestGetKeyComponents() {
	assert.Equal(s.T(), []string{s.device.OrganizationId, s.device.Id}, s.device.GetKeyComponents(), "should return correct key components")
}

func (s *DeviceTestSuite) TestSerialize() {
	data, err := s.device.Serialize()
	assert.Equal(s.T(), string(data), s.json, "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *DeviceTestSuite) TestValidate() {
	device := Device{}

	assert.Error(s.T(), device.Validate(), "should error on empty ID")
	device.Id = s.device.Id

	assert.Error(s.T(), device.Validate(), "should error on empty organization ID")
	device.OrganizationId = s.device.OrganizationId

	assert.Error(s.T(), device.Validate(), "should error on empty name")
	device.Name = s.device.Name

	assert.Error(s.T(), device.Validate(), "should error on empty last update time")
	device.LastUpdateTime = s.device.LastUpdateTime

	assert.Nil(s.T(), device.Validate(), "should return no error")
}

func (s *DeviceTestSuite) TestDeserializeDevice() {
	device, err := DeserializeDevice([]byte(s.json))
	assert.Equal(s.T(), device, s.device, "should return parsed device")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeDevice([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestDeviceTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}
