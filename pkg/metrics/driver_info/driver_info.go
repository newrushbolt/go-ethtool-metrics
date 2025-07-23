// Driver info, eg `ethtool -i ethX`
package driver_info

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

func parseCommonInfo(input string) *DriverInfoCommon {
	var output DriverInfoCommon
	inputMap := common.ParseAbstractColonData(input, "", true)
	common.ParseAbstractDataObject(&inputMap, &output, "driver")
	return &output
}

func parseFeatures(input string) *DriverFeatures {
	var output DriverFeatures
	inputMap := common.ParseAbstractColonData(input, "", true)
	common.ParseAbstractDataObject(&inputMap, &output, "driver_supports")
	return &output
}

func ParseInfo(rawInfo string, config *CollectConfig) *DriverInfo {
	loggerLever := common.GetLogLevel()
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: loggerLever}))

	if rawInfo == "" {
		Logger.Info("Module got empty ethtool data, skipping", "module", "driver_info")
		return nil
	}

	var driverInfoCommon *DriverInfoCommon
	if config.CollectCommon {
		driverInfoCommon = parseCommonInfo(rawInfo)
	}

	var driverFeatures *DriverFeatures
	if config.CollectFeatures {
		driverFeatures = parseFeatures(rawInfo)
	}

	driverInfo := DriverInfo{
		Common:   driverInfoCommon,
		Features: driverFeatures,
	}
	return &driverInfo
}
