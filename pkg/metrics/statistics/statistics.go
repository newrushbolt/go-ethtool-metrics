// Statistics, eg `ethtool -S ethX`
package statistics

import (
	"github.com/newrushbolt/go-ethtool-metrics/internal"
)

func ParseInfo(rawInfo string) *StatisticsInfo {
	generalStatisticsMap, _ := internal.ParseAbstractColonData(rawInfo, "", true)
	var generalStatistics GeneralStatisticsInfo
	internal.ParseAbstractDataObject(&generalStatisticsMap, &generalStatistics, "general_statistics")

	statistics := StatisticsInfo{generalStatistics}
	return &statistics
	// deviceInfoMap, _ := internal.ParseAbstractColonData(rawInfo, "", true)
	// var device_info DriverInfo
	// internal.ParseAbstractDataObject(&deviceInfoMap, &device_info, "driver")
	// device_info.Features = &features
	// return &device_info
}
