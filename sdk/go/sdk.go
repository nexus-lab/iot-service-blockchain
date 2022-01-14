package sdk

import (
	"crypto/x509"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/nexus-lab/iot-service-blockchain/common"
)

// Sdk the iot service blockchain sdk
type Sdk struct {
	gw              *client.Gateway
	organizationId  string
	deviceId        string
	deviceRegistry  DeviceRegistryInterface
	serviceRegistry ServiceRegistryInterface
	serviceBroker   ServiceBrokerInterface
}

// SdkOptions SDK initialization options
type SdkOptions struct {
	// OrganizationId organization/MSP ID
	OrganizationId string

	// Certificate PEM-formated X509 client certificate
	Certificate []byte

	// PrivateKey PEM-formated client private key
	PrivateKey []byte

	// GatewayPeerEndpoint network address of the gateway peer
	GatewayPeerEndpoint string

	// GatewayPeerServerName server name of the gateway peer
	GatewayPeerServerName string

	// GatewayPeerTLSCertificate PEM-formated X509 TLS certificate of the gateway peer
	GatewayPeerTLSCertificate []byte

	// NetworkName blockchain network channel name
	NetworkName string

	// ChaincodeId name of the chaincode
	ChaincodeId string
}

func newIdentity(organizationId string, certificate []byte) (*identity.X509Identity, error) {
	cert, err := identity.CertificateFromPEM(certificate)
	if err != nil {
		return nil, err
	}

	id, err := identity.NewX509Identity(organizationId, cert)
	if err != nil {
		return nil, err
	}

	return id, err
}

func newSign(privateKey []byte) (identity.Sign, error) {
	key, err := identity.PrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	sign, err := identity.NewPrivateKeySign(key)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

func newGrpcConnection(endpoint, serverName string, certificate []byte) (*grpc.ClientConn, error) {
	cert, err := identity.CertificateFromPEM(certificate)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AddCert(cert)
	credentials := credentials.NewClientTLSFromCert(pool, serverName)

	connection, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials))
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// NewSdk create a new SDK instance from options
func NewSdk(options *SdkOptions) (*Sdk, error) {
	sdk := new(Sdk)

	id, err := newIdentity(options.OrganizationId, options.Certificate)
	if err != nil {
		return nil, err
	}

	sign, err := newSign(options.PrivateKey)
	if err != nil {
		return nil, err
	}

	conn, err := newGrpcConnection(
		options.GatewayPeerEndpoint,
		options.GatewayPeerServerName,
		options.GatewayPeerTLSCertificate,
	)
	if err != nil {
		return nil, err
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(conn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	sdk.gw = gw

	network := gw.GetNetwork(options.NetworkName)

	sdk.connectSmartContracts(network, options.ChaincodeId)
	if err = sdk.setIdentity(options.OrganizationId, options.Certificate); err != nil {
		return nil, err
	}

	return sdk, nil
}

func (s *Sdk) setIdentity(organizationId string, certificate []byte) error {
	cert, err := common.ParseCertificate(certificate)
	if err != nil {
		return err
	}

	clientId, err := common.GetClientId(cert)
	if err != nil {
		return err
	}

	s.deviceId = clientId
	s.organizationId = organizationId

	return nil
}

func (s *Sdk) connectSmartContracts(network *client.Network, chaincodeId string) {
	s.deviceRegistry = createDeviceRegistry(network, chaincodeId)
	s.serviceRegistry = createServiceRegistry(network, chaincodeId)
	s.serviceBroker = createServiceBroker(network, chaincodeId)
}

// GetDeviceId return the device/client ID of the current calling application
func (s *Sdk) GetDeviceId() string {
	return s.deviceId
}

// GetOrganizationId return the organization ID of the current calling application
func (s *Sdk) GetOrganizationId() string {
	return s.organizationId
}

// GetDeviceRegistry return the device registry
func (s *Sdk) GetDeviceRegistry() DeviceRegistryInterface {
	return s.deviceRegistry
}

// GetServiceRegistry return the service registry
func (s *Sdk) GetServiceRegistry() ServiceRegistryInterface {
	return s.serviceRegistry
}

// GetServiceBroker return the service broker
func (s *Sdk) GetServiceBroker() ServiceBrokerInterface {
	return s.serviceBroker
}

// Close close connection to the Hyperledger Fabric gateway
func (s *Sdk) Close() {
	if s.gw != nil {
		s.gw.Close()
	}
}
