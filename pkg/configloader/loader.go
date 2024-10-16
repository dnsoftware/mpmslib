// Package configloader Загрузка конфигурационных данных из распределенной системы конфигураций
// Если она не доступна - загружаем из локального файла
// Локальный файл формируется из данных, полученных из DCS (когда она доступна)
// Настроить триггер отслеживания изменений в удаленном конфиге!

package configloader

import (
	"github.com/dnsoftware/mpmslib/pkg/dcs"
)

// ConfigLoader получение конфигурационных данных
type ConfigLoader struct {
	dcs             dcs.DCS // удаленный конфиг
	localConfigFile string  // путь к локальному конфиг файлу
}

func NewConfigLoader(dcs dcs.DCS, localConfigFile string) (*ConfigLoader, error) {

	cfgLoader := &ConfigLoader{
		dcs:             dcs,
		localConfigFile: localConfigFile,
	}

	return cfgLoader, nil
}

func (c *ConfigLoader) LoadConfigData() (error, string) {

	return nil, ""
}
