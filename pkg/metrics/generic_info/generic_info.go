// Generic info, eg `ethtool ethX`
package generic_info

import (
	"log/slog"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/newrushbolt/go-ethtool-metrics/internal"
)

func _DropHeaderLine(input string) string {
	secondLineIndex := strings.Index(input, "\n")
	if (strings.HasPrefix(input, "Settings for ")) && (secondLineIndex > 0) {
		output := input[secondLineIndex:]
		return output
	} else {
		return input
	}
}

func _GetPortSpeedBytes(input string) (speedBytes float64) {
	speedRe := regexp.MustCompile(`(\d+)(.+)`)
	result_slice := speedRe.FindAllStringSubmatch(input, -1)
	var rawSpeedBytes float64
	var speedSuffix string
	result := result_slice[0]
	if len(result) != 3 {
		return math.NaN()
	}
	var err error
	rawSpeedBytes, err = strconv.ParseFloat(result[1], 64)
	if err != nil {
		slog.Error("Cannot get float64 from speed string", "speed_string", input)
		return math.NaN()
	}
	speedSuffix = result[2]

	var speedMultiplier float64
	// Doing straight metric convesion, not 2^x
	switch speedSuffix {
	case "Mb/s":
		speedMultiplier = 1000 * 1000
	case "Gb/s":
		speedMultiplier = 1000 * 1000 * 1000
	default:
		slog.Error("Cannot get speed units from string, must have 'Gb/s' or 'Mb/s'", "speed_string", "input")
		return math.NaN()
	}
	speedBytes = rawSpeedBytes * speedMultiplier
	return speedBytes
}

func _ParseSupportedSettings(input string) *AvaliableSettings {
	var output AvaliableSettings
	inputMap := internal.ParseAbstractColonData(input, "Supported ", false)
	internal.ParseAbstractDataObject(&inputMap, &output, "generic_info_avaliable_settings")
	return &output
}

func _ParseAdvertisedSettings(input string) *AvaliableSettings {
	var output AvaliableSettings
	inputMap := internal.ParseAbstractColonData(input, "Advertised ", false)
	internal.ParseAbstractDataObject(&inputMap, &output, "generic_info_avaliable_settings")
	return &output
}

func _ParseSettings(input string) *Settings {
	var output Settings
	inputMap := internal.ParseAbstractColonData(input, "", true)
	internal.ParseAbstractDataObject(&inputMap, &output, "generic_info_settings")
	return &output
}

func ParseInfo(rawInfo string, config *CollectConfig) *GenericInfo {
	slog.SetLogLoggerLevel(internal.GetLogLevel())

	if rawInfo == "" {
		slog.Info("Module got empty ethtool data, skipping", "module", "generic_info")
		return nil
	}

	cleanInput := _DropHeaderLine(rawInfo)

	var supportedSetting *AvaliableSettings
	if config.CollectSupportedSettings {
		supportedSetting = _ParseSupportedSettings(cleanInput)
	}

	var advertisedSettings *AvaliableSettings
	if config.CollectAdvertisedSettings {
		advertisedSettings = _ParseAdvertisedSettings(cleanInput)
	}

	var settings *Settings
	if config.CollectSettings {
		settings = _ParseSettings(cleanInput)
	}

	commonInfo := GenericInfo{
		SupportedSettings:  supportedSetting,
		AdvertisedSettings: advertisedSettings,
		Settings:           settings,
	}
	if (commonInfo.Settings.Speed != "Unknown!") && (commonInfo.Settings.Speed != "") {
		speedBytes := _GetPortSpeedBytes(commonInfo.Settings.Speed)
		commonInfo.Settings.SpeedBytes = &speedBytes
	}
	return &commonInfo
}
