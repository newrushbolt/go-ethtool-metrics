// Statistics, eg `ethtool -S ethX`
package statistics

import (
	"log/slog"
	"os"

	"github.com/newrushbolt/go-ethtool-metrics/common"
)

var (
	Logger *slog.Logger
)

func init() {
	loggerLever := common.GetLogLevel()
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: loggerLever}))
}

func ParseInfo(rawInfo string, config *CollectConfig) *StatisticsInfo {
	loggerLever := common.GetLogLevel()
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: loggerLever}))

	if rawInfo == "" {
		Logger.Info("Module got empty ethtool data, skipping", "module", "Statistics")
		return nil
	}

	statistics := StatisticsInfo{}
	generalStatisticsMap := common.ParseAbstractColonData(rawInfo, "", true)

	if config.PerQueue {
		statistics.PerQueue = parseQueuedInfo(generalStatisticsMap, *config)
	}

	if config.General {
		var generalStatistics GeneralStatistics
		common.ParseAbstractDataObject(&generalStatisticsMap, &generalStatistics, "general_statistics")
		statistics.General = &generalStatistics
	}
	return &statistics
}
