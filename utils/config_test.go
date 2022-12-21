package common_utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPath(t *testing.T) {
	path := GetPath()
	_, err := os.Stat(path)
	require.NoError(t, err)
}

func TestLoad(t *testing.T) {
	path := GetPath()
	conf, err := LoadConfig(path)
	require.NoError(t, err)
	require.NotNil(t, conf)
}
