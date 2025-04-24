package internal

import (
	"log/slog"
	"os"
)

func GetLogLevel() slog.Level {
	var level slog.Level
	envLevel := os.Getenv("GO_ETHTOOL_METRICS_LOG_LEVEL")
	err := level.UnmarshalText([]byte(envLevel))
	if err != nil {
		return slog.LevelDebug
	}
	return level
}
