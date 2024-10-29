package configloader

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Tls struct {
	CaFile   string `yaml:"ca_file"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type Etcd struct {
	Endpoints string `yaml:"endpoints"`
	Auth      Auth   `yaml:"auth"`
	Tls       Tls    `yaml:"tls"`
}

type StartConfig struct {
	Etcd Etcd `yaml:"etcd"`
}

// LoadStartConfig загрузка стартового конфига для получения доступа к кластеру etcd, где хранится основная конфигурация
func LoadStartConfig(fullPath string) (*StartConfig, error) {

	log.Printf("loading config @ `%s`", fullPath)
	rawCfg, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	cfg := StartConfig{}
	if err := yaml.Unmarshal(rawCfg, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
