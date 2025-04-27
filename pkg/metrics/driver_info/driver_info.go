// Driver info, eg `ethtool -i ethX`
package driver_info

import (
	"log/slog"

	"github.com/newrushbolt/go-ethtool-metrics/internal"
)

func ParseInfo(rawInfo string, config *CollectConfig) *DriverInfo {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if rawInfo == "" {
		slog.Info("Module got empty ethtool data, skipping", "module", "driver_info")
		return nil
	}

	deviceInfoMap, _ := internal.ParseAbstractColonData(rawInfo, "", true)
	var device_info DriverInfo
	internal.ParseAbstractDataObject(&deviceInfoMap, &device_info, "driver")

	if config.DriverFeatures {
		var features DriverFeatures
		internal.ParseAbstractDataObject(&deviceInfoMap, &features, "driver_supports")
		device_info.Features = &features
	}
	return &device_info
}
