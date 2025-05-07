package internal_test

import (
	"log/slog"
	"testing"

	"github.com/newrushbolt/go-ethtool-metrics/internal"
	"github.com/stretchr/testify/assert"
)

func TestGetLogLevelDefault(t *testing.T) {
	t.Setenv("GO_ETHTOOL_METRICS_LOG_LEVEL", "")
	level := internal.GetLogLevel()
	assert.Equal(t, level, slog.LevelDebug)
}

func TestGetLogLevel(t *testing.T) {
	t.Setenv("GO_ETHTOOL_METRICS_LOG_LEVEL", "ERROR")
	level := internal.GetLogLevel()
	assert.Equal(t, level, slog.LevelError)
}
