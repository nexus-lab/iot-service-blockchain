package common

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	service *Service
	json    string
}

func (s *ServiceTestSuite) SetupTest() {
	lastUpdateTime, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")

	deviceId := "eDUwOTo6Q049dXNlcjEsT1U9Y2xpZW50LE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUzo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw=="
	s.json = fmt.Sprintf("{\"name\":\"service1\",\"deviceId\":\"%s\",\"organizationId\":\"Org1MSP\",\"version\":1,\"description\":\"Service of Device1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}", deviceId)
	s.service = &Service{
		Name:           "service1",
		DeviceId:       deviceId,
		OrganizationId: "Org1MSP",
		Version:        1,
		Description:    "Service of Device1",
		LastUpdateTime: lastUpdateTime,
	}
}

func (s *ServiceTestSuite) TestGetKeyComponents() {
	assert.Equal(s.T(), []string{s.service.OrganizationId, s.service.DeviceId, s.service.Name}, s.service.GetKeyComponents(), "should return correct key components")
}

func (s *ServiceTestSuite) TestSerialize() {
	data, err := s.service.Serialize()
	assert.Equal(s.T(), string(data), s.json, "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceTestSuite) TestValidate() {
	service := Service{}

	assert.Error(s.T(), service.Validate(), "should error on empty name")
	service.Name = s.service.Name

	assert.Error(s.T(), service.Validate(), "should error on empty device ID")
	service.DeviceId = s.service.DeviceId

	assert.Error(s.T(), service.Validate(), "should error on empty organization ID")
	service.OrganizationId = s.service.OrganizationId

	assert.Error(s.T(), service.Validate(), "should error on empty version")
	service.Version = s.service.Version

	assert.Error(s.T(), service.Validate(), "should error on empty last update time")
	service.LastUpdateTime = s.service.LastUpdateTime

	assert.Nil(s.T(), service.Validate(), "should return no error")
}

func (s *ServiceTestSuite) TestDeserializeService() {
	data := s.json
	service, err := DeserializeService([]byte(data))
	assert.Equal(s.T(), service, s.service, "should return parsed service")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeService([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
