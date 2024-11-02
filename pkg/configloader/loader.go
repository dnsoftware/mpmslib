// Package configloader Загрузка конфигурационных данных из распределенной системы конфигураций
// Если она не доступна - загружаем из локального файла
// Локальный файл формируется из данных, полученных из DCS (когда она доступна)
// Настроить триггер отслеживания изменений в удаленном конфиге!

package configloader

import (
	"fmt"
	"os"
	"strings"

	"github.com/dnsoftware/mpmslib/pkg/dcs"
	"github.com/dnsoftware/mpmslib/pkg/tlscerts"
	"github.com/dnsoftware/mpmslib/pkg/utils"
)

// ConfigLoader получение конфигурационных данных с удаленного сервера etcd
type ConfigLoader struct {
	dcs             *dcs.DCS // удаленный конфиг
	localConfigFile string   // путь к локальному конфиг файлу
}

type Option func(loader *ConfigLoader)

// NewConfigLoader формирует структуру получения удаленных конфигурационных данных
// clusterNode список адресов нод кластера etcd с портами: []string{"31.128.39.18:2379", "31.129.98.136:2379", "45.147.179.134:2379"}
// caPath, publicPath, privatePath - пути к корневому сертификату, к сертификату и приватному ключу
// remoteDataKey - ключ на удаленном сервере, данные по которому мы хотим получить
// localConfigPath - путь сохранения полученного удаленного конфига локально
func NewConfigLoader(
	clusterNode []string,
	caPath, publicPath, privatePath string,
	localConfigPath string,
	dcsUsername, dcsPassword string) (*ConfigLoader, error) {

	tlsCerts, err := tlscerts.NewTLSCerts(caPath, publicPath, privatePath)
	if err != nil {
		return nil, err
	}

	// Вносим пароль для юзера root, указанный в процессе развертывания etcd кластера
	dcsConf, err := dcs.NewSecureDCS(clusterNode, dcsUsername, dcsPassword, tlsCerts)
	if err != nil {
		return nil, err
	}

	cfgLoader := &ConfigLoader{
		dcs:             dcsConf,
		localConfigFile: localConfigPath,
	}

	return cfgLoader, nil
}

// LoadRemoteConfig Загрузка удаленного конфига по ключу
func (c *ConfigLoader) LoadRemoteConfig(remoteDataKey string) (string, error) {

	cfgData, err := c.dcs.LoadConfig(remoteDataKey)
	if err != nil {
		return "", fmt.Errorf("LoadRemoteConfig: %w", err)
	}

	return cfgData, nil
}

// MultiloadRemoteConfig загрузка данных по нескольким ключам и объединение в одну строку в порядке следования элементов в массиве ключей
func (c *ConfigLoader) MultiloadRemoteConfig(remoteDataKeys []string) (string, error) {
	var data []string
	for _, val := range remoteDataKeys {
		cfgData, err := c.dcs.LoadConfig(val)
		if err != nil {
			return "", fmt.Errorf("MultiloadRemoteConfig error load, key: %s", val)
		}
		data = append(data, cfgData)
	}

	return strings.Join(data, "\n"), nil
}

func (c *ConfigLoader) SaveConfigToFile(cfgData string) error {

	if !utils.FileExists(c.localConfigFile) {
		err := utils.CreateFileWithDirs(c.localConfigFile)
		if err != nil {
			return fmt.Errorf("CreateFileWithDirs: %w", err)
		}
	}

	err := utils.SaveTextToFile(c.localConfigFile, cfgData)
	if err != nil {
		return fmt.Errorf("SaveConfigToFile error: %w", err)
	}

	return nil
}

func (c *ConfigLoader) LoadLocalConfig() (string, error) {

	if !utils.FileExists(c.localConfigFile) {
		return "", fmt.Errorf("config file does not exist at the specified path")
	}

	data, err := os.ReadFile(c.localConfigFile)
	if err != nil {
		return "", fmt.Errorf("LoadLocalConfig: %w", err)
	}

	return string(data), nil
}
