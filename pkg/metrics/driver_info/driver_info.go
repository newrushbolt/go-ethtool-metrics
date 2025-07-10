// Driver info, eg `ethtool -i ethX`
package driver_info

import (
	"log/slog"
	"os"

	"github.com/newrushbolt/go-ethtool-metrics/internal"
)

var (
	Logger *slog.Logger
)

func ParseInfo(rawInfo string, config *CollectConfig) *DriverInfo {
	loggerLever := internal.GetLogLevel()
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: loggerLever}))

	if rawInfo == "" {
		Logger.Info("Module got empty ethtool data, skipping", "module", "driver_info")
		return nil
	}

	deviceInfoMap := internal.ParseAbstractColonData(Logger, rawInfo, "", true)
	var device_info DriverInfo
	internal.ParseAbstractDataObject(Logger, &deviceInfoMap, &device_info, "driver")

	if config.DriverFeatures {
		var features DriverFeatures
		internal.ParseAbstractDataObject(Logger, &deviceInfoMap, &features, "driver_supports")
		device_info.Features = &features
	}
	return &device_info
}
