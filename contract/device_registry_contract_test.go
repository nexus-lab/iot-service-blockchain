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
	ctx := &MockTransactionContext{DeviceId: "device1", OrganizationId: "org1"}
	deviceRegistry := new(MockDeviceRegistry)
	ctx.deviceRegistry = deviceRegistry

	deviceRegistry.On("Register", mock.AnythingOfType("*common.Device")).Return(nil)

	contract := new(DeviceRegistrySmartContract)
	err := contract.Register(ctx, fmt.Sprintf("{\"id\":\"%s\",\"organizationId\":\"%s\",\"name\":\"Device1\",\"description\":\"Device of Org1 User1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}", ctx.DeviceId, ctx.OrganizationId))
	assert.Nil(s.T(), err, "should return no error")
	called := deviceRegistry.AssertCalled(s.T(), "Register", mock.AnythingOfType("*common.Device"))
	assert.True(s.T(), called, "should put device to device registry")
	device, _ := common.DeserializeDevice(ctx.stub.EventPayload)
	assert.Equal(s.T(), fmt.Sprintf("device://%s/%s/register", ctx.OrganizationId, ctx.DeviceId), ctx.stub.EventName, "should emit event with name")
	assert.Equal(s.T(), ctx.DeviceId, device.Id, "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Register(ctx, "{\"id\":\"device2\",\"organizationId\":\"org2\",\"name\":\"device2\",\"description\":\"Device of Org2 User1\",\"lastUpdateTime\":\"2021-12-12T17:34:00-05:00\"}")
	assert.Error(s.T(), err, "should return mismatch device ID and organization ID error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")

	err = contract.Register(ctx, "[]")
	assert.Error(s.T(), err, "should return deserialization error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")
}

func (s *DeviceRegistryContractTestSuite) TestGet() {
	ctx := &MockTransactionContext{DeviceId: "device2", OrganizationId: "org2"}
	deviceRegistry := new(MockDeviceRegistry)
	ctx.deviceRegistry = deviceRegistry

	deviceRegistry.On("Get", "org1", "device1").Return(new(common.Device), nil)

	contract := new(DeviceRegistrySmartContract)
	_, _ = contract.Get(ctx, "org1", "device1")
	called := deviceRegistry.AssertCalled(s.T(), "Get", "org1", "device1")
	assert.True(s.T(), called, "should retrieve device from device registry")
}

func (s *DeviceRegistryContractTestSuite) TestGetAll() {
	ctx := &MockTransactionContext{DeviceId: "device2", OrganizationId: "org2"}
	deviceRegistry := new(MockDeviceRegistry)
	ctx.deviceRegistry = deviceRegistry

	deviceRegistry.On("GetAll", "org1").Return([]*common.Device{{}, {}}, nil)

	contract := new(DeviceRegistrySmartContract)
	_, _ = contract.GetAll(ctx, "org1")
	called := deviceRegistry.AssertCalled(s.T(), "GetAll", "org1")
	assert.True(s.T(), called, "should retrieve devices from device registry")
}

func (s *DeviceRegistryContractTestSuite) TestDeregister() {
	ctx := &MockTransactionContext{DeviceId: "device1", OrganizationId: "org1"}
	deviceRegistry := new(MockDeviceRegistry)
	ctx.deviceRegistry = deviceRegistry

	deviceRegistry.On("Deregister", mock.MatchedBy(func(device *common.Device) bool {
		return device.Id == "device1" && device.OrganizationId == "org1"
	})).Return(nil)
	deviceRegistry.On("Deregister", mock.Anything).Return(new(common.NotFoundError))

	contract := new(DeviceRegistrySmartContract)
	err := contract.Deregister(ctx, fmt.Sprintf("{\"id\":\"%s\",\"organizationId\":\"%s\"}", ctx.DeviceId, ctx.OrganizationId))
	assert.Nil(s.T(), err, "should return no error")
	actual := deviceRegistry.Calls[0].Arguments[0].(*common.Device)
	assert.Equal(s.T(), ctx.DeviceId, actual.Id, "should remove the correct device")
	assert.Equal(s.T(), ctx.OrganizationId, actual.OrganizationId, "should remove the correct device")
	device, _ := common.DeserializeDevice(ctx.stub.EventPayload)
	assert.Equal(s.T(), fmt.Sprintf("device://%s/%s/deregister", ctx.OrganizationId, ctx.DeviceId), ctx.stub.EventName, "should emit event with name")
	assert.Equal(s.T(), ctx.DeviceId, device.Id, "should emit event with payload")
	ctx.stub.ResetEvent()

	err = contract.Deregister(ctx, "{\"id\":\"device2\",\"organizationId\":\"org2\"}")
	assert.Error(s.T(), err, "should return mismatch device ID and organization ID error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")

	ctx.DeviceId = "device2"
	err = contract.Deregister(ctx, fmt.Sprintf("{\"id\":\"%s\",\"organizationId\":\"%s\"}", ctx.DeviceId, ctx.OrganizationId))
	assert.IsType(s.T(), new(common.NotFoundError), err, "should return not found error")
	assert.Empty(s.T(), ctx.stub.EventName, "should not emit event")
}

func TestDeviceRegistryContractTestSuite(t *testing.T) {
	suite.Run(t, new(DeviceRegistryContractTestSuite))
}
