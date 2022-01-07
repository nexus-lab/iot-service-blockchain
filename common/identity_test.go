package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	PUBLIC_KEY = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
-----END PUBLIC KEY-----`
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
	CLIENT_ID = "eDUwOTo6Q049dXNlcjEsT1U9Y2xpZW50LE89SHlwZXJsZWRnZXIsU1Q" +
		"9Tm9ydGggQ2Fyb2xpbmEsQz1VUzo6Q049Y2Eub3JnMi5leGFtcGxlLmNvbSxPPW9yZz" +
		"IuZXhhbXBsZS5jb20sTD1IdXJzbGV5LFNUPUhhbXBzaGlyZSxDPVVL"
)

type IdentityTestSuite struct {
	suite.Suite
}

func (s *IdentityTestSuite) TestParseCertificate() {
	cert, err := ParseCertificate("-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----")
	assert.Error(s.T(), err, "should return error if PEM is invalid")
	assert.Nil(s.T(), cert, "should not return certicate on error")

	cert, err = ParseCertificate(PUBLIC_KEY)
	assert.Error(s.T(), err, "should return error if certificate is invalid")
	assert.Nil(s.T(), cert, "should not return certicate on error")

	cert, err = ParseCertificate(CERTIFICATE)
	assert.NotNil(s.T(), cert, "should return certicate")
	assert.Nil(s.T(), err, "should return no error if certificate is valid")
}

func (s *IdentityTestSuite) TestGetClientId() {
	cert, _ := ParseCertificate(CERTIFICATE)
	id, err := GetClientId(cert)
	assert.Equal(s.T(), CLIENT_ID, id, "should return correct ID")
	assert.Nil(s.T(), err, "should return no error if certificate is valid")
}

func TestIdentityTestSuite(t *testing.T) {
	suite.Run(t, new(IdentityTestSuite))
}
