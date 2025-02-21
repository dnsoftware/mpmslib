package servicediscovery

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/dnsoftware/mpmslib/pkg/utils"
)

// поднять экземпляр etcd
func TestServiceDiscovery(t *testing.T) {
	etcdEndpoints := strings.Split("127.0.0.1:2379,127.0.0.1:2479,127.0.0.1:2579", ",")
	serviceKey := "serviceTest"

	projectRootAnchorFile := ".env"
	basePath, err := utils.GetProjectRoot(projectRootAnchorFile)
	require.NoError(t, err)

	caPath := basePath + "/certs/ca.crt"
	publicPath := basePath + "/certs/client.crt"
	privatePath := basePath + "/certs/client.key"

	// Root certificate CA load
	caCert, err := os.ReadFile(caPath)
	require.NoError(t, err)

	// Make certs pool
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load cert and key
	cert, err := tls.LoadX509KeyPair(publicPath, privatePath)
	require.NoError(t, err)

	// Make tls.Config with a configured root certificate
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		// Enabling verification of the server certificate
		InsecureSkipVerify: false,
	}

	cfg := clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
		Username:    "root",
		Password:    "securepassword",
	}

	serviceAddr := "127.0.0.1:8080"
	sd, err := NewServiceDiscovery(cfg, "/service_discovery/services", serviceKey, serviceAddr, 5, 10)
	require.NoError(t, err)
	allServices, err := sd.DiscoverAllServices()
	require.NoError(t, err)
	for key, val := range allServices {
		fmt.Println(key, val)
	}

	srv, err := sd.DiscoverService(serviceKey)
	require.NoError(t, err)
	require.Equal(t, serviceAddr, srv)

	sd.RegisterService(serviceKey+":grpc", "127.0.0.1:7878")
	srv, err = sd.DiscoverService(serviceKey + ":grpc")
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1:7878", srv)

	sd.RegisterService(serviceKey+":grpc2", "127.0.0.1:4444")
	srv, err = sd.DiscoverService(serviceKey + ":grpc2")
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1:4444", srv)

}
