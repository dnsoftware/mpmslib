// Package configloader Загрузка конфигурационных данных из распределенной системы конфигураций
// Если она не доступна - загружаем из локального файла
// Локальный файл формируется из данных, полученных из DCS (когда она доступна)
// Настроить триггер отслеживания изменений в удаленном конфиге!

package configloader

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/dnsoftware/mpmslib/pkg/dcs"
	"github.com/dnsoftware/mpmslib/pkg/utils"
)

// ConfigLoader получение конфигурационных данных
type ConfigLoader struct {
	dcs             *dcs.DCS // удаленный конфиг
	localConfigFile string   // путь к локальному конфиг файлу
	logger          *zap.Logger
}

type Option func(loader *ConfigLoader)

func WithLogger(logger *zap.Logger) Option {
	return func(s *ConfigLoader) {
		s.logger = logger
	}
}

func NewConfigLoader(dcs *dcs.DCS, localConfigFile string, options ...Option) (*ConfigLoader, error) {

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	cfgLoader := &ConfigLoader{
		dcs:             dcs,
		localConfigFile: localConfigFile,
		logger:          logger,
	}

	for _, option := range options {
		option(cfgLoader)
	}

	return cfgLoader, nil
}

func (c *ConfigLoader) LoadRemoteConfig() (string, error) {

	cfgData, err := c.dcs.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("LoadRemoteConfig: %w", err)
	}

	if !utils.FileExists(c.localConfigFile) {
		err = utils.CreateFileWithDirs(c.localConfigFile)
		if err != nil {
			return "", fmt.Errorf("CreateFileWithDirs: %w", err)
		}
	}

	err = utils.SaveTextToFile(c.localConfigFile, cfgData)
	if err != nil {
		return "", fmt.Errorf("SaveTextToFile: %w", err)
	}

	return cfgData, nil
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
