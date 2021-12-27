package contract

import (
	"crypto/x509"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MockClientIdentity struct {
	Id    string
	MspId string
}

type MockTransactionContext struct {
	contractapi.TransactionContext
	clientId        *MockClientIdentity
	deviceRegistry  common.DeviceRegistryInterface
	serviceRegistry common.ServiceRegistryInterface
	serviceBroker   common.ServiceBrokerInterface
}

func (c *MockTransactionContext) GetDeviceRegistry() common.DeviceRegistryInterface {
	return c.deviceRegistry
}

func (c *MockTransactionContext) GetServiceRegistry() common.ServiceRegistryInterface {
	return c.serviceRegistry
}

func (c *MockTransactionContext) GetServiceBroker() common.ServiceBrokerInterface {
	return c.serviceBroker
}

func (i *MockClientIdentity) GetID() (string, error) {
	return i.Id, nil
}

func (i *MockClientIdentity) GetMSPID() (string, error) {
	return i.MspId, nil
}

func (i *MockClientIdentity) GetAttributeValue(attrName string) (value string, found bool, err error) {
	return "", false, nil
}

func (i *MockClientIdentity) AssertAttributeValue(attrName, attrValue string) error {
	return nil
}

func (i *MockClientIdentity) GetX509Certificate() (*x509.Certificate, error) {
	return nil, nil
}

func (c *MockTransactionContext) GetClientIdentity() cid.ClientIdentity {
	if c.clientId == nil {
		c.clientId = &MockClientIdentity{Id: "Device1Id", MspId: "Org1MSP"}
	}
	return c.clientId
}

type TransactionContextTestSuite struct {
	suite.Suite
	ctx TransactionContext
}

func (s *TransactionContextTestSuite) SetupTest() {
	s.ctx = TransactionContext{}
}

func (s *TransactionContextTestSuite) TestGetDeviceRegistry() {
	expected := createDeviceRegistry(&s.ctx)
	actual := s.ctx.GetDeviceRegistry().(*DeviceRegistry)
	assert.Equal(s.T(), expected.stateRegistry.(*StateRegistry).Name, actual.stateRegistry.(*StateRegistry).Name, "should return device registry")
}

func (s *TransactionContextTestSuite) TestGetServiceRegistry() {
	expected := createServiceRegistry(&s.ctx)
	actual := s.ctx.GetServiceRegistry().(*ServiceRegistry)
	assert.Equal(s.T(), expected.stateRegistry.(*StateRegistry).Name, actual.stateRegistry.(*StateRegistry).Name, "should return service registry")
}

func (s *TransactionContextTestSuite) TestGetServiceBroker() {
	expected := createServiceBroker(&s.ctx)
	actual := s.ctx.GetServiceBroker().(*ServiceBroker)
	assert.Equal(s.T(), expected.requestRegistry.(*StateRegistry).Name, actual.requestRegistry.(*StateRegistry).Name, "should return service broker")
}

func TestTransactionContextTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionContextTestSuite))
}
