package dcs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/dnsoftware/mpmslib/pkg/tlscerts"
)

// DCS (Distributed Configuration System)
type DCS struct {
	clusterNodes []string // list of IP addresses and ports of the cluster node (ip:port, ip:port ...)
	dcsUsername  string   // DCS access user name
	dcsPassword  string   // DCS access password
	certs        tlscerts.TLSCerts
	client       *clientv3.Client
}

// NewSecureDCS working with a remote config over a TLS secure connection
func NewSecureDCS(clusterNodes []string, dcsUsername string, dcsPassword string, certs tlscerts.TLSCerts) (*DCS, error) {

	// Root certificate CA load
	caCert, err := os.ReadFile(certs.CA)
	if err != nil {
		return nil, fmt.Errorf("error reading CA certificate: %w", err)
	}

	// Make certs pool
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load cert and key
	cert, err := tls.LoadX509KeyPair(certs.PublicKey, certs.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error certificates loading: %w", err)
	}
	// Make tls.Config with a configured root certificate
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		// Enabling verification of the server certificate
		InsecureSkipVerify: false,
	}
	_ = tlsConfig

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   clusterNodes,
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
		Username:    dcsUsername,
		Password:    dcsPassword,
	})
	if err != nil {
		return nil, err
	}

	dcs := &DCS{
		clusterNodes: clusterNodes,
		dcsPassword:  dcsPassword,
		certs:        certs,
		client:       cli,
	}

	return dcs, nil
}

// LoadConfig - getting config data by key
//
//	remoteDataKey string   // путь к конфигурационным данным в DCS (Distributed Configuration System)
func (d *DCS) LoadConfig(remoteDataKey string) (string, error) {
	timeout := time.Second * 1

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	respFrom, err := d.client.Get(ctx, remoteDataKey)
	cancel()

	if err != nil {
		return "", err
	}

	data := ""
	if len(respFrom.Kvs) == 1 {
		data = string(respFrom.Kvs[0].Value)
	} else {
		return "", fmt.Errorf("bad range of config response")
	}

	defer d.client.Close()

	return data, nil
}
