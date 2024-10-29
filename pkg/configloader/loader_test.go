package configloader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	mpmslogger "github.com/dnsoftware/mpmslib/pkg/logger"
	"github.com/dnsoftware/mpmslib/pkg/utils"
)

func setup(remoteDataKey string) (*ConfigLoader, error) {
	// Указываем адреса рабочих нод кластера
	clusterNode := []string{"31.128.39.18:2379", "31.129.98.136:2379", "45.147.179.134:2379"}

	projectRoot, err := utils.GetProjectRoot("readme.md")
	if err != nil {
		return nil, err
	}

	// Прописываем пути к сертификатам (копируем рабочие сертификаты сгенерированные в процессе развертывания etcd кластера)
	caPath := projectRoot + "/tests/assets/certs/ca.crt"
	publicPath := projectRoot + "/tests/assets/certs/cert.crt"
	privatePath := projectRoot + "/tests/assets/certs/cert.key"
	localConfigPath := projectRoot + "/tests/configloader/config.yaml"
	dcsUsername := "root"
	dcsPassword := "etcdpassword"

	logPath := projectRoot + "/tests/configloader/log.log"
	logger, err := mpmslogger.NewLogger(logPath, zapcore.ErrorLevel)
	if err != nil {
		return nil, err
	}

	confLoader, err := NewConfigLoader(remoteDataKey, clusterNode, caPath, publicPath, privatePath, localConfigPath, dcsUsername, dcsPassword, WithLogger(logger))
	if err != nil {
		return nil, err
	}

	return confLoader, nil
}

func TestLoadConfig(t *testing.T) {
	confLoader, err := setup("/testnamespace/badconfig.yaml")
	require.NoError(t, err)

	confData, err := confLoader.LoadRemoteConfig()
	require.Error(t, err)

	confLoader, err = setup("/testnamespace/config.yaml")
	require.NoError(t, err)

	confData, err = confLoader.LoadRemoteConfig()
	require.NotEmpty(t, confData)

	// Если не удалось загрузить удаленный конфиг - логируем этот факт и загружаем из локального файла
	if err != nil {
		confLoader.logger.Error("Remote config does not load: " + err.Error())

		confData, err = confLoader.LoadLocalConfig()
		require.NoError(t, err)
		require.NotEmpty(t, confData)
	}

	// Тестирование наблюдателя
	changedConfig := make(chan string)
	confLoader.dcs.ActivateWatcher(changedConfig)

	newdata := ""
	err = confLoader.dcs.SaveConfig("olddata")
	err = confLoader.dcs.SaveConfig("newdata")
	require.NoError(t, err)

	newdata = <-changedConfig
	fmt.Println(newdata)
}
