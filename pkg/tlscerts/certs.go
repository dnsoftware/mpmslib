package tlscerts

import (
	"errors"

	"github.com/dnsoftware/mpmslib/pkg/utils"
)

// TLSCerts сертификаты для TLS соединений
type TLSCerts struct {
	CA         string // центр сертификации (ca.crt)
	PublicKey  string // публичный ключ cert.crt
	PrivateKey string // приватный ключ cert.key
}

func NewTLSCerts(caPath, publicKeyPAth, privateKeyPath string) (TLSCerts, error) {
	tlsCerts := TLSCerts{}

	if !utils.FileExists(caPath) {
		return tlsCerts, errors.New("CA cert does not exist")
	}

	if !utils.FileExists(publicKeyPAth) {
		return tlsCerts, errors.New("public cert does not exist")
	}

	if !utils.FileExists(privateKeyPath) {
		return tlsCerts, errors.New("private cert does not exist")
	}

	tlsCerts.CA = caPath
	tlsCerts.PublicKey = publicKeyPAth
	tlsCerts.PrivateKey = privateKeyPath

	return tlsCerts, nil
}
