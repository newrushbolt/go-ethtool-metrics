package common

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogLevelDefault(t *testing.T) {
	t.Setenv("GO_ETHTOOL_METRICS_LOG_LEVEL", "")
	level := GetLogLevel()
	assert.Equal(t, level, slog.LevelInfo)
}

func TestGetLogLevel(t *testing.T) {
	t.Setenv("GO_ETHTOOL_METRICS_LOG_LEVEL", "ERROR")
	level := GetLogLevel()
	assert.Equal(t, level, slog.LevelError)
}
