package common

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceRequestTestSuite struct {
	suite.Suite
	request *ServiceRequest
	json    string
}

func (s *ServiceRequestTestSuite) SetupTest() {
	time, _ := time.Parse(time.RFC3339, "2021-12-12T17:34:00-05:00")

	deviceId := "eDUwOTo6Q049dXNlcjEsT1U9Y2xpZW50LE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUzo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw=="
	s.json = fmt.Sprintf("{\"id\":\"ffbc9005-c62a-4563-a8f7-b32bba27d707\",\"time\":\"2021-12-12T17:34:00-05:00\",\"service\":{\"name\":\"service1\",\"deviceId\":\"%s\",\"organizationId\":\"Org1MSP\",\"version\":0,\"description\":\"\",\"lastUpdateTime\":\"0001-01-01T00:00:00Z\"},\"method\":\"GET\",\"arguments\":[\"1\",\"2\",\"3\"]}", deviceId)
	s.request = &ServiceRequest{
		Id: "ffbc9005-c62a-4563-a8f7-b32bba27d707",
		Service: Service{
			Name:           "service1",
			DeviceId:       deviceId,
			OrganizationId: "Org1MSP",
		},
		Method:    "GET",
		Arguments: []string{"1", "2", "3"},
		Time:      time,
	}
}

func (s *ServiceRequestTestSuite) TestGetKeyComponents() {
	assert.Equal(s.T(), []string{s.request.Id}, s.request.GetKeyComponents(), "should return correct key components")
}

func (s *ServiceRequestTestSuite) TestSerialize() {
	data, err := s.request.Serialize()
	assert.Equal(s.T(), string(data), s.json, "should serialize to JSON")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *ServiceRequestTestSuite) TestValidate() {
	request := ServiceRequest{}

	request.Id = "123456"
	assert.Error(s.T(), request.Validate(), "should error on invalid ID")
	request.Id = s.request.Id

	assert.Error(s.T(), request.Validate(), "should error on empty serivce organization ID")
	request.Service.OrganizationId = s.request.Service.OrganizationId

	assert.Error(s.T(), request.Validate(), "should error on empty service device ID")
	request.Service.DeviceId = s.request.Service.DeviceId

	assert.Error(s.T(), request.Validate(), "should error on empty service")
	request.Service.Name = s.request.Service.Name

	assert.Error(s.T(), request.Validate(), "should error on empty method")
	request.Method = s.request.Method

	assert.Error(s.T(), request.Validate(), "should error on empty last update time")
	request.Time = s.request.Time

	assert.Nil(s.T(), request.Validate(), "should return no error")
}

func (s *ServiceRequestTestSuite) TestDeserializeService() {
	data := s.json
	request, err := DeserializeServiceRequest([]byte(data))
	assert.Equal(s.T(), request, s.request, "should return parsed service")
	assert.Nil(s.T(), err, "should return no error")

	_, err = DeserializeServiceRequest([]byte{0x00})
	assert.Error(s.T(), err, "should return an error")
}

func TestServiceRequestTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceRequestTestSuite))
}
