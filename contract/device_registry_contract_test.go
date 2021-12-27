package contract

import (
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
	err := contract.Register(ctx, "{\"id\":\"Device1Id\",\"organizationId\":\"Org1Id\",\"name\":\"Device1\",\"description\":\"Device of Org1 User1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}")
	assert.Nil(s.T(), err, "should return no error")
	called := deviceRegistry.AssertCalled(s.T(), "Register", mock.AnythingOfType("*common.Device"))
	assert.True(s.T(), called, "should put device to device registry")
	device := deviceRegistry.Calls[0].Arguments[0].(*common.Device)
	assert.Equal(s.T(), deviceId, device.Id, "should change device ID")
	assert.Equal(s.T(), organizationId, device.OrganizationId, "should change organization ID")

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
	expected := new(common.Device)
	deviceRegistry.On("Get", organizationId, deviceId).Return(expected, nil)
	deviceRegistry.On("Get", mock.Anything, mock.Anything).Return(nil, new(common.NotFoundError))
	deviceRegistry.On("Deregister", expected).Return(nil)

	contract := new(DeviceRegistrySmartContract)
	err := contract.Deregister(ctx)
	assert.Nil(s.T(), err, "should return no error")
	called := deviceRegistry.AssertCalled(s.T(), "Deregister", expected)
	assert.True(s.T(), called, "should remove device to device registry")
	actual := deviceRegistry.Calls[1].Arguments[0].(*common.Device)
	assert.Equal(s.T(), expected, actual, "should remove the correct device")

	ctx.clientId = &MockClientIdentity{Id: "Device2Id", MspId: "Org2MSP"}
	err = contract.Deregister(ctx)
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
}

func TestDeviceRegistryContractTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceRegistryContractTestSuite))
}
