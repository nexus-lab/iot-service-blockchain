package common

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"sort"
	"strings"
)

var dnOrder = []string{"CN", "SERIALNUMBER", "C", "L", "ST", "STREET", "O", "OU", "POSTALCODE"}
var oidMap = map[string]string{
	"2.5.4.3":  "CN",
	"2.5.4.5":  "SERIALNUMBER",
	"2.5.4.6":  "C",
	"2.5.4.7":  "L",
	"2.5.4.8":  "ST",
	"2.5.4.9":  "STREET",
	"2.5.4.10": "O",
	"2.5.4.11": "OU",
	"2.5.4.17": "POSTALCODE",
}

// ParseCertificate parse an X509 certificate from PEM string
func ParseCertificate(certificate []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certificate)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return x509.ParseCertificate(block.Bytes)
}

func escapeDN(dn string) string {
	escaped := make([]rune, 0, len(dn))

	for k, c := range dn {
		escape := false

		switch c {
		case ',', '+', '"', '\\', '<', '>', ';':
			escape = true

		case ' ':
			escape = k == 0 || k == len(dn)-1

		case '#':
			escape = k == 0
		}

		if escape {
			escaped = append(escaped, '\\', c)
		} else {
			escaped = append(escaped, c)
		}
	}

	return string(escaped)
}

func mapDN(rdns pkix.RDNSequence) map[string][]string {
	dnMap := make(map[string][]string)

	for _, rdn := range rdns {
		for _, tv := range rdn {
			typeName, ok := oidMap[tv.Type.String()]
			if !ok {
				continue
			}

			_, ok = dnMap[typeName]
			// CN and SERIALNUMBER are single-value field
			if !ok || typeName == "CN" || typeName == "SERIALNUMBER" {
				dnMap[typeName] = make([]string, 0)
			}

			dnMap[typeName] = append(dnMap[typeName], fmt.Sprint(tv.Value))
		}
	}

	for _, values := range dnMap {
		sort.Strings(values)
	}

	return dnMap
}

/**
 * Returns a string representation of the distinguished name,
 * roughly following the RFC 2253 Distinguished Names syntax.
 * Distinguished Names are sorted by their OID name.
 */
func formatDN(rdns pkix.RDNSequence) string {
	dnMap := mapDN(rdns)

	allValues := make([]string, 0)
	for _, typeName := range dnOrder {
		if _, ok := dnMap[typeName]; !ok {
			continue
		}

		values := make([]string, 0)
		for _, value := range dnMap[typeName] {
			values = append(values, fmt.Sprintf("%s=%s", typeName, escapeDN(value)))
		}

		allValues = append(allValues, strings.Join(values, "+"))
	}

	return strings.Join(allValues, ",")
}

// GetClientId return unique client ID from client certificate.
func GetClientId(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", fmt.Errorf("cannot determine identity")
	}
	id := fmt.Sprintf("x509::%s::%s", formatDN(cert.Subject.ToRDNSequence()), formatDN(cert.Issuer.ToRDNSequence()))
	return base64.StdEncoding.EncodeToString([]byte(id)), nil
}
