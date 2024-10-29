package configloader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnsoftware/mpmslib/pkg/utils"
)

func TestLoadStartConfig(t *testing.T) {
	path, err := utils.GetProjectRoot("readme.md")
	require.NoError(t, err)

	fullPath := path + "/tests/configloader/startconf.yaml"
	startConf, err := LoadStartConfig(fullPath)
	require.NoError(t, err)

	fmt.Println(startConf)
}
