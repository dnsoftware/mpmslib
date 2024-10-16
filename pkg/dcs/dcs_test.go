package dcs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/mpmslib/pkg/tlscerts"
	"github.com/dnsoftware/mpmslib/pkg/utils"
)

func TestLoadConfig(t *testing.T) {

	clusterNode := []string{"31.128.39.18:2379", "31.129.98.136:2379", "45.147.179.134:2379"}

	projectRoot, err := utils.GetProjectRoot("readme.md")
	require.NoError(t, err)

	caPath := projectRoot + "/tests/assets/certs/ca.crt"
	publicPath := projectRoot + "/tests/assets/certs/cert.crt"
	privatePath := projectRoot + "/tests/assets/certs/cert.key"

	tlsCerts, err := tlscerts.NewTLSCerts(caPath, publicPath, privatePath)
	require.NoError(t, err)

	dcs, err := NewSecureDCS(clusterNode, "root", "etcdpassword", tlsCerts)
	require.NoError(t, err)

	confData, err := dcs.LoadConfig("/testnamespace/config.yaml")
	require.NoError(t, err)

	fmt.Println(confData)
}
