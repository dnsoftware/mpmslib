package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectRoot(t *testing.T) {
	path, err := GetProjectRoot("readme.md")

	assert.NoError(t, err)
	assert.NotEmpty(t, path)
}
