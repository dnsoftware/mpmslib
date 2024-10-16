package tlscerts

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/mpmslib/pkg/utils"
)

func TestNewTLSCerts(t *testing.T) {
	projectRoot, err := utils.GetProjectRoot("readme.md")
	require.NoError(t, err)

	caPath := ""
	publicPath := ""
	privatePath := ""
	_, err = NewTLSCerts(caPath, publicPath, privatePath)

	caPath = projectRoot + "/tests/assets/certs/ca.crt"
	_, err = NewTLSCerts(caPath, publicPath, privatePath)
	require.Error(t, err)

	publicPath = projectRoot + "/tests/assets/certs/cert.crt"
	_, err = NewTLSCerts(caPath, publicPath, privatePath)
	require.Error(t, err)

	privatePath = projectRoot + "/tests/assets/certs/cert.key"
	tlsCerts, err := NewTLSCerts(caPath, publicPath, privatePath)
	require.NoError(t, err)

	_ = tlsCerts
}
