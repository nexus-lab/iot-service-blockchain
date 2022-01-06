package common

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// ParseCertificate parse an X509 certificate from PEM string
func ParseCertificate(rawCert string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(rawCert))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return x509.ParseCertificate(block.Bytes)
}

// GetClientId return unique client ID from client certificate.
// See https://pkg.go.dev/github.com/hyperledger/fabric-chaincode-go@v0.0.0-20210718160520-38d29fabecb9/pkg/cid#ClientID.GetID
func GetClientId(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", fmt.Errorf("cannot determine identity")
	}
	id := fmt.Sprintf("x509::%s::%s", cert.Subject.ToRDNSequence().String(), cert.Issuer.ToRDNSequence().String())
	return base64.StdEncoding.EncodeToString([]byte(id)), nil
}
