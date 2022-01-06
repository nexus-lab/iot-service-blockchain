package contract

import (
	"crypto/x509"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/nexus-lab/iot-service-blockchain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	CERTIFICATE = `
-----BEGIN CERTIFICATE-----
MIICmjCCAkCgAwIBAgIUd/uzCIgYnvr5IVrGgnVXIF/JvWMwCgYIKoZIzj0EAwIw
bDELMAkGA1UEBhMCVUsxEjAQBgNVBAgTCUhhbXBzaGlyZTEQMA4GA1UEBxMHSHVy
c2xleTEZMBcGA1UEChMQb3JnMi5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eub3Jn
Mi5leGFtcGxlLmNvbTAeFw0yMjAxMDYwMjE1MDBaFw0yMzAxMDYwMjIwMDBaMF0x
CzAJBgNVBAYTAlVTMRcwFQYDVQQIEw5Ob3J0aCBDYXJvbGluYTEUMBIGA1UEChML
SHlwZXJsZWRnZXIxDzANBgNVBAsTBmNsaWVudDEOMAwGA1UEAxMFdXNlcjEwWTAT
BgcqhkjOPQIBBggqhkjOPQMBBwNCAARe9edmNbHEx0pQJP3jfGgjtIDp0a/dmzR4
fi74zEQMKYz8E0nt/BTCGC8Uv9SRvBHI7biYW1k8WXfkCoPmPTjuo4HOMIHLMA4G
A1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQE+JosrrNPvToO
byzv7BkPFxp1QTAfBgNVHSMEGDAWgBQaaZoL4EglLGspr66g1a2vf83MvDARBgNV
HREECjAIggZIb21lUEMwWAYIKgMEBQYHCAEETHsiYXR0cnMiOnsiaGYuQWZmaWxp
YXRpb24iOiIiLCJoZi5FbnJvbGxtZW50SUQiOiJ1c2VyMSIsImhmLlR5cGUiOiJj
bGllbnQifX0wCgYIKoZIzj0EAwIDSAAwRQIhALxDGVIsgP3VxXMzrv+l0ijGgX4T
/AmTkI+tB0LZqzprAiAm3oeXhmFmxUXTnFXbumz7xelcodKByxXLHyAkucX/NA==
-----END CERTIFICATE-----`
	MSP_ID    = "Org1MSP"
	CLIENT_ID = "eDUwOTo6Q049dXNlcjEsT1U9Y2xpZW50LE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUzo6Q049Y2Eub3JnMi5leGFtcGxlLmNvbSxPPW9yZzIuZXhhbXBsZS5jb20sTD1IdXJzbGV5LFNUPUhhbXBzaGlyZSxDPVVL"
)

type mockClientIdentity struct{}

func (i *mockClientIdentity) GetID() (string, error) {
	return CLIENT_ID, nil
}

func (i *mockClientIdentity) GetMSPID() (string, error) {
	return MSP_ID, nil
}

func (i *mockClientIdentity) GetAttributeValue(attrName string) (value string, found bool, err error) {
	return "", false, nil
}

func (i *mockClientIdentity) AssertAttributeValue(attrName, attrValue string) error {
	return nil
}

func (i *mockClientIdentity) GetX509Certificate() (*x509.Certificate, error) {
	return common.ParseCertificate(CERTIFICATE)
}

type mockChaincodeStub struct {
	shim.ChaincodeStub
	EventName    string
	EventPayload []byte
}

func (s *mockChaincodeStub) SetEvent(name string, payload []byte) error {
	s.EventName = name
	s.EventPayload = payload
	return nil
}

func (s *mockChaincodeStub) ResetEvent() {
	s.EventName = ""
	s.EventPayload = nil
}

type MockTransactionContext struct {
	contractapi.TransactionContext
	identity        *mockClientIdentity
	stub            *mockChaincodeStub
	deviceRegistry  DeviceRegistryInterface
	serviceRegistry ServiceRegistryInterface
	serviceBroker   ServiceBrokerInterface

	DeviceId       string
	OrganizationId string
}

func (c *MockTransactionContext) GetOrganizationId() (string, error) {
	return c.OrganizationId, nil
}

func (c *MockTransactionContext) GetDeviceId() (string, error) {
	return c.DeviceId, nil
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
	if c.identity == nil {
		c.identity = new(mockClientIdentity)
	}
	return c.identity
}

func (c *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	if c.stub == nil {
		c.stub = new(mockChaincodeStub)
	}
	return c.stub
}

type TransactionContextTestSuite struct {
	suite.Suite
	ctx *TransactionContext
}

func (s *TransactionContextTestSuite) SetupTest() {
	s.ctx = new(TransactionContext)
	s.ctx.SetClientIdentity(new(mockClientIdentity))
}

func (s *TransactionContextTestSuite) TestGetOrganizationId() {
	organizationId, err := s.ctx.GetOrganizationId()
	assert.Equal(s.T(), MSP_ID, organizationId, "should return correct organization ID")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *TransactionContextTestSuite) TestGetDeviceId() {
	deviceId, err := s.ctx.GetDeviceId()
	assert.Equal(s.T(), CLIENT_ID, deviceId, "should return correct device ID")
	assert.Nil(s.T(), err, "should return no error")
}

func (s *TransactionContextTestSuite) TestGetDeviceRegistry() {
	expected := createDeviceRegistry(s.ctx)
	actual := s.ctx.GetDeviceRegistry().(*DeviceRegistry)
	assert.Equal(s.T(), expected.stateRegistry.(*StateRegistry).Name, actual.stateRegistry.(*StateRegistry).Name, "should return device registry")
}

func (s *TransactionContextTestSuite) TestGetServiceRegistry() {
	expected := createServiceRegistry(s.ctx)
	actual := s.ctx.GetServiceRegistry().(*ServiceRegistry)
	assert.Equal(s.T(), expected.stateRegistry.(*StateRegistry).Name, actual.stateRegistry.(*StateRegistry).Name, "should return service registry")
}

func (s *TransactionContextTestSuite) TestGetServiceBroker() {
	expected := createServiceBroker(s.ctx)
	actual := s.ctx.GetServiceBroker().(*ServiceBroker)
	assert.Equal(s.T(), expected.requestRegistry.(*StateRegistry).Name, actual.requestRegistry.(*StateRegistry).Name, "should return service broker")
}

func TestTransactionContextTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionContextTestSuite))
}
