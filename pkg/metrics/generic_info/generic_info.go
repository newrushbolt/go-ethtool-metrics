// Generic info, eg `ethtool ethX`
package generic_info

import (
	"log/slog"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/newrushbolt/go-ethtool-metrics/common"
)

var (
	Logger  *slog.Logger
	speedRe = regexp.MustCompile(`(\d+)(.+)`)
)

func init() {
	loggerLever := common.GetLogLevel()
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: loggerLever}))
}

func _DropHeaderLine(input string) string {
	secondLineIndex := strings.Index(input, "\n")
	if (strings.HasPrefix(input, "Settings for ")) && (secondLineIndex > 0) {
		output := input[secondLineIndex:]
		return output
	} else {
		return input
	}
}

func _GetPortSpeedBits(rawInput string) float64 {
	// TODO: rewrite parse logic in a more elegant way
	var rawSpeedBytes float64

	input := strings.Trim(rawInput, " \t")
	resultSlice := speedRe.FindAllStringSubmatch(input, -1)
	if resultSlice == nil {
		Logger.Error("Cannot get speed units from string", "speed_string", rawInput)
		return math.NaN()
	}
	result := resultSlice[0]

	// We expect exactly three matches, because first match is full match which we never use
	// The second is numeric value, and the third is postfix that defines the unit
	if len(result) != 3 {
		Logger.Error("Cannot get speed units from string", "speed_string", rawInput)
		return math.NaN()
	}
	var err error
	rawSpeedBytes, err = strconv.ParseFloat(result[1], 64)
	if err != nil {
		Logger.Error("Cannot get float64 from speed string", "speed_string", rawInput)
		return math.NaN()
	}

	var speedMultiplier float64
	speedPostfix := result[2]
	switch speedPostfix {
	case "Mb/s":
		// Doing straight metric conversion, not 2^x
		speedMultiplier = 1000 * 1000
	case "Gb/s":
		speedMultiplier = 1000 * 1000 * 1000
	default:
		Logger.Error("Cannot get speed units from string, must have 'Gb/s' or 'Mb/s'", "speed_string", input)
		return math.NaN()
	}
	speedBits := rawSpeedBytes * speedMultiplier
	return speedBits
}

func _ParseSupportedSettings(input string) *AvaliableSettings {
	var output AvaliableSettings
	inputMap := common.ParseAbstractColonData(input, "Supported ", false)
	common.ParseAbstractDataObject(&inputMap, &output, "generic_info_avaliable_settings")
	return &output
}

func _ParseAdvertisedSettings(input string) *AvaliableSettings {
	var output AvaliableSettings
	inputMap := common.ParseAbstractColonData(input, "Advertised ", false)
	common.ParseAbstractDataObject(&inputMap, &output, "generic_info_avaliable_settings")
	return &output
}

func _ParseSettings(input string) *Settings {
	var output Settings
	inputMap := common.ParseAbstractColonData(input, "", true)
	common.ParseAbstractDataObject(&inputMap, &output, "generic_info_settings")
	return &output
}

func ParseInfo(rawInfo string, config *CollectConfig) *GenericInfo {
	loggerLever := common.GetLogLevel()
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: loggerLever}))

	if rawInfo == "" {
		Logger.Info("Module got empty ethtool data, skipping", "module", "generic_info")
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
		speedBits := _GetPortSpeedBits(commonInfo.Settings.Speed)
		speedBytes := speedBits / 8
		commonInfo.Settings.SpeedBits = &speedBits
		commonInfo.Settings.SpeedBytes = &speedBytes
	}
	return &commonInfo
}
