// Statistics, eg `ethtool -S ethX`
package statistics

import (
	"github.com/newrushbolt/go-ethtool-metrics/internal"
)

func ParseInfo(rawInfo string) *StatisticsInfo {
	if rawInfo == "" {
		return nil
	}

	generalStatisticsMap, _ := internal.ParseAbstractColonData(rawInfo, "", true)
	var generalStatistics GeneralStatisticsInfo
	internal.ParseAbstractDataObject(&generalStatisticsMap, &generalStatistics, "general_statistics")

	statistics := StatisticsInfo{generalStatistics}
	return &statistics
}
