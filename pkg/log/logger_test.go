package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	logger := NewLogger("debug", "console")
	require.NotNil(t, logger)
}
