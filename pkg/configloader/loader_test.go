package configloader

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/dnsoftware/mpmslib/pkg/dcs"
	mpmslogger "github.com/dnsoftware/mpmslib/pkg/logger"
	"github.com/dnsoftware/mpmslib/pkg/tlscerts"
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

	tlsCerts, err := tlscerts.NewTLSCerts(caPath, publicPath, privatePath)
	if err != nil {
		return nil, err
	}

	// Вносим пароль для юзера root, указанный в процессе развертывания etcd кластера
	dcsConf, err := dcs.NewSecureDCS(clusterNode, "root", "etcdpassword", tlsCerts, remoteDataKey)
	if err != nil {
		return nil, err
	}

	logPath := projectRoot + "/tests/configloader/log.log"
	logger, err := mpmslogger.NewLogger(logPath, zapcore.ErrorLevel)
	if err != nil {
		return nil, err
	}

	localConfigPath := projectRoot + "/tests/configloader/config.yaml"
	confLoader, err := NewConfigLoader(dcsConf, localConfigPath, WithLogger(logger))
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

	changedConfig := make(chan string)
	go confLoader.dcs.ActivateWatcher(changedConfig)

	newdata := ""
	err = confLoader.dcs.SaveConfig("olddata")
	time.Sleep(2 * time.Second)
	err = confLoader.dcs.SaveConfig("newdata")
	require.NoError(t, err)

	newdata = <-changedConfig

	require.Equal(t, "newdata", newdata)
	fmt.Println(newdata)
}
