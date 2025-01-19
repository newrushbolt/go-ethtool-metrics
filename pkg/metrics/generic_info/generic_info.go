// Generic info, eg `ethtool ethX`
package generic_info

import (
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

func _GetPortSpeedBytes(input string) (speedBytes uint64) {
	speedRe := regexp.MustCompile(`(\d+)(.+)`)
	result_slice := speedRe.FindAllStringSubmatch(input, -1)
	var rawSpeedBytes uint64
	var speedSuffix string
	result := result_slice[0]
	if len(result) == 3 {
		var _ error
		rawSpeedBytes, _ = strconv.ParseUint(result[1], 10, 64)
		// TODO: log error
		speedSuffix = result[2]

	} else {
		return uint64(0)
	}
	speedMultiplier := uint64(0)
	// Doing straight metric convesion, not 2^x
	switch speedSuffix {
	case "Mb/s":
		speedMultiplier = 1000 * 1000
	case "Gb/s":
		speedMultiplier = 1000 * 1000 * 1000
	}
	// TODO: log unit error
	speedBytes = rawSpeedBytes * speedMultiplier
	return speedBytes
}

func _ParseSupportedSettings(input string) *AvaliableSettings {
	inputMap, _ := internal.ParseAbstractColonData(input, "Supported ", false)
	var output AvaliableSettings
	internal.ParseAbstractDataObject(&inputMap, &output, "generic_info_avaliable_settings")
	return &output
}

func _ParseAdvertisedSettings(input string) *AvaliableSettings {
	inputMap, _ := internal.ParseAbstractColonData(input, "Advertised ", false)
	var output AvaliableSettings
	internal.ParseAbstractDataObject(&inputMap, &output, "generic_info_avaliable_settings")
	return &output
}

func _ParseSettings(input string) *Settings {
	inputMap, _ := internal.ParseAbstractColonData(input, "", true)
	var output Settings
	internal.ParseAbstractDataObject(&inputMap, &output, "generic_info_settings")
	return &output
}

func ParseInfo(input string) *GenericInfo {
	cleanInput := _DropHeaderLine(input)

	// TODO: add config
	// var supportedSetting *AvaliableSettings
	// Add if based on config later
	supportedSetting := _ParseSupportedSettings(cleanInput)

	// TODO: add config
	// var advertisedSettings *AvaliableSettings
	// Add if based on config later
	advertisedSettings := _ParseAdvertisedSettings(cleanInput)

	// TODO: add config
	// var advertisedSettings *AvaliableSettings
	// Add if based on config later
	settings := _ParseSettings(cleanInput)

	commonInfo := GenericInfo{
		SupportedSettings:  supportedSetting,
		AdvertisedSettings: advertisedSettings,
		Settings:           settings,
	}
	if (commonInfo.Settings.Speed != "Unknown!") && (commonInfo.Settings.Speed != "") {
		commonInfo.Settings.SpeedBytes = _GetPortSpeedBytes(commonInfo.Settings.Speed)
	}
	return &commonInfo
}
