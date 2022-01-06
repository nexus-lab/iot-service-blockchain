package contract

import (
	"fmt"
	"testing"

	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DeviceRegistryContractTestSuite struct {
	suite.Suite
}

func (s *DeviceRegistryContractTestSuite) TestRegister() {
	ctx := new(MockTransactionContext)
	deviceRegistry := new(MockDeviceRegistry)
	ctx.deviceRegistry = deviceRegistry

	deviceId, _ := ctx.GetClientIdentity().GetID()
	organizationId, _ := ctx.GetClientIdentity().GetMSPID()
	deviceRegistry.On("Register", mock.AnythingOfType("*common.Device")).Return(nil)

	contract := new(DeviceRegistrySmartContract)
	err := contract.Register(ctx, fmt.Sprintf("{\"id\":\"%s\",\"organizationId\":\"%s\",\"name\":\"Device1\",\"description\":\"Device of Org1 User1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}", deviceId, organizationId))
	assert.Nil(s.T(), err, "should return no error")
	called := deviceRegistry.AssertCalled(s.T(), "Register", mock.AnythingOfType("*common.Device"))
	assert.True(s.T(), called, "should put device to device registry")

	err = contract.Register(ctx, "{\"id\":\"Device2Id\",\"organizationId\":\"Org2MSP\",\"name\":\"Device1\",\"description\":\"Device of Org1 User1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}")
	assert.Error(s.T(), err, "should return mismatch device ID and organization ID error")

	err = contract.Register(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
}

func (s *DeviceRegistryContractTestSuite) TestGet() {
	ctx := new(MockTransactionContext)
	deviceRegistry := new(MockDeviceRegistry)
	ctx.deviceRegistry = deviceRegistry

	deviceRegistry.On("Get", "Org1MSP", "Device1Id").Return(new(common.Device), nil)

	contract := new(DeviceRegistrySmartContract)
	_, _ = contract.Get(ctx, "Org1MSP", "Device1Id")
	called := deviceRegistry.AssertCalled(s.T(), "Get", "Org1MSP", "Device1Id")
	assert.True(s.T(), called, "should retrieve device from device registry")
}

func (s *DeviceRegistryContractTestSuite) TestGetAll() {
	ctx := new(MockTransactionContext)
	deviceRegistry := new(MockDeviceRegistry)
	ctx.deviceRegistry = deviceRegistry

	deviceRegistry.On("GetAll", "Org1MSP").Return([]*common.Device{{}, {}}, nil)

	contract := new(DeviceRegistrySmartContract)
	_, _ = contract.GetAll(ctx, "Org1MSP")
	called := deviceRegistry.AssertCalled(s.T(), "GetAll", "Org1MSP")
	assert.True(s.T(), called, "should retrieve devices from device registry")
}

func (s *DeviceRegistryContractTestSuite) TestDeregister() {
	ctx := new(MockTransactionContext)
	deviceRegistry := new(MockDeviceRegistry)
	ctx.deviceRegistry = deviceRegistry

	deviceId, _ := ctx.GetClientIdentity().GetID()
	organizationId, _ := ctx.GetClientIdentity().GetMSPID()
	deviceRegistry.On("Deregister", mock.MatchedBy(func(device *common.Device) bool {
		return device.Id == deviceId && device.OrganizationId == organizationId
	})).Return(nil)
	deviceRegistry.On("Deregister", mock.Anything).Return(new(common.NotFoundError))

	contract := new(DeviceRegistrySmartContract)
	err := contract.Deregister(ctx, fmt.Sprintf("{\"id\":\"%s\",\"organizationId\":\"%s\"}", deviceId, organizationId))
	assert.Nil(s.T(), err, "should return no error")
	actual := deviceRegistry.Calls[0].Arguments[0].(*common.Device)
	assert.Equal(s.T(), deviceId, actual.Id, "should remove the correct device")
	assert.Equal(s.T(), organizationId, actual.OrganizationId, "should remove the correct device")

	err = contract.Deregister(ctx, "{\"id\":\"Device2Id\",\"organizationId\":\"Org2MSP\"}")
	assert.Error(s.T(), err, "should return mismatch device ID and organization ID error")

	ctx.clientId = &MockClientIdentity{Id: "Device2Id", MspId: "Org2MSP"}
	err = contract.Deregister(ctx, "{\"id\":\"Device2Id\",\"organizationId\":\"Org2MSP\"}")
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func TestDeviceRegistryContractTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceRegistryContractTestSuite))
}
