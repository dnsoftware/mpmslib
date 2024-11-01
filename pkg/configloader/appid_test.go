package configloader

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppID(t *testing.T) {
	code := "ALPH_POOLER_PROP_EU_BEGET_1"

	appID, err := NewAppID(code, "_")
	require.NoError(t, err)

	full := appID.GetFullID()
	require.Equal(t, code, full)

	coinLevel := appID.GetCoinID()
	require.Equal(t, "ALPH", coinLevel)
}
