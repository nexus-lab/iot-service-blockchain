package contract

import (
	"crypto/x509"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MockClientIdentity struct {
	Id    string
	MspId string
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

type MockChaincodeStub struct {
	shim.ChaincodeStub
	eventName    string
	eventPayload []byte
}

func (s *MockChaincodeStub) SetEvent(name string, payload []byte) error {
	s.eventName = name
	s.eventPayload = payload
	return nil
}

func (s *MockChaincodeStub) ResetEvent() {
	s.eventName = ""
	s.eventPayload = nil
}

type MockTransactionContext struct {
	contractapi.TransactionContext
	clientId        *MockClientIdentity
	stub            *MockChaincodeStub
	deviceRegistry  DeviceRegistryInterface
	serviceRegistry ServiceRegistryInterface
	serviceBroker   ServiceBrokerInterface
}

func (c *MockTransactionContext) GetDeviceRegistry() DeviceRegistryInterface {
	return c.deviceRegistry
}

func (c *MockTransactionContext) GetServiceRegistry() ServiceRegistryInterface {
	return c.serviceRegistry
}

func (c *MockTransactionContext) GetServiceBroker() ServiceBrokerInterface {
	return c.serviceBroker
}

func (c *MockTransactionContext) GetClientIdentity() cid.ClientIdentity {
	if c.clientId == nil {
		c.clientId = &MockClientIdentity{Id: "Device1", MspId: "Org1MSP"}
	}
	return c.clientId
}

func (c *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	if c.stub == nil {
		c.stub = &MockChaincodeStub{}
	}
	return c.stub
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
