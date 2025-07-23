package metrics_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/newrushbolt/go-ethtool-metrics/pkg/metrics/driver_info"
	"github.com/newrushbolt/go-ethtool-metrics/pkg/metrics/generic_info"
	"github.com/newrushbolt/go-ethtool-metrics/pkg/metrics/module_info"
	"github.com/newrushbolt/go-ethtool-metrics/pkg/metrics/statistics"

	"github.com/stretchr/testify/assert"
)

func ReadFixturePair(fixtureName string, ethtoolMode string, settingsMode string) (src string, result string) {
	fixtureSourcePath := fmt.Sprintf("../testdata/%s/src/%s", fixtureName, ethtoolMode)
	fixtureResultPath := fmt.Sprintf("../testdata/%s/results/%s.%s.json", fixtureName, ethtoolMode, settingsMode)

	fixtureSourceData, err := os.ReadFile(fixtureSourcePath)
	// Allow missing source
	// Probably should catch only `FileNotFound` error here
	if err != nil {
		// Return empty string if source file was not found
		fixtureSourceData = []byte{}
	}

	// Result file must exist anyway. If src is empty or missing, just put `null` in result file
	fixtureResultData, err := os.ReadFile(fixtureResultPath)
	if err != nil {
		panic(err)
	}

	return string(fixtureSourceData), string(fixtureResultData)
}

func GetFixtureList() []string {
	return []string{
		"intel/i40e/00_sfp_10g_sr85",
		"intel/i40e/01_int_tp",
		"intel/i40e/02_sfp_10_or_25g_sr",
		"intel/igb/00_int_tp",
		"broadcom/bnxt_en/00_sfp_10gsr85",
	}
}

func TestModuleInfoDefault(t *testing.T) {
	testModule := "module_info"
	config := module_info.CollectConfig{}.Default()
	for _, fixture := range GetFixtureList() {
		t.Run(fixture, func(t *testing.T) {
			srcFile, resultFile := ReadFixturePair(fixture, testModule, "default")
			info := module_info.ParseInfo(srcFile, config)
			infoJson, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resultFile, string(infoJson))
		})
	}
}

func TestModuleInfoFull(t *testing.T) {
	testModule := "module_info"
	config := module_info.CollectConfig{
		CollectDiagnosticsAlarms:   true,
		CollectDiagnosticsValues:   true,
		CollectDiagnosticsWarnings: true,
		CollectVendor:              true,
	}
	for _, fixture := range GetFixtureList() {
		t.Run(fixture, func(t *testing.T) {
			srcFile, resultFile := ReadFixturePair(fixture, testModule, "full")
			info := module_info.ParseInfo(srcFile, &config)
			infoJson, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resultFile, string(infoJson))
		})
	}
}

func TestDriverInfoDefault(t *testing.T) {
	testModule := "driver_info"
	for _, fixture := range GetFixtureList() {
		t.Run(fixture, func(t *testing.T) {
			srcFile, resultFile := ReadFixturePair(fixture, testModule, "default")
			config := driver_info.CollectConfig{}.Default()
			info := driver_info.ParseInfo(srcFile, config)
			infoJson, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resultFile, string(infoJson))
		})
	}
}

func TestDriverInfoFull(t *testing.T) {
	testModule := "driver_info"
	for _, fixture := range GetFixtureList() {
		t.Run(fixture, func(t *testing.T) {
			srcFile, resultFile := ReadFixturePair(fixture, testModule, "full")
			config := driver_info.CollectConfig{
				CollectFeatures: true,
				CollectCommon:   true,
			}
			info := driver_info.ParseInfo(srcFile, &config)
			infoJson, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resultFile, string(infoJson))
		})
	}
}

func TestGenericInfoDefault(t *testing.T) {
	testModule := "generic_info"
	for _, fixture := range GetFixtureList() {
		t.Run(fixture, func(t *testing.T) {
			srcFile, resultFile := ReadFixturePair(fixture, testModule, "default")
			config := generic_info.CollectConfig{}.Default()
			info := generic_info.ParseInfo(srcFile, config)
			infoJson, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resultFile, string(infoJson))
		})
	}
}

func TestGenericInfoFull(t *testing.T) {
	testModule := "generic_info"
	for _, fixture := range GetFixtureList() {
		t.Run(fixture, func(t *testing.T) {
			srcFile, resultFile := ReadFixturePair(fixture, testModule, "full")
			config := generic_info.CollectConfig{
				CollectAdvertisedSettings: true,
				CollectSupportedSettings:  true,
				CollectSettings:           true,
			}
			info := generic_info.ParseInfo(srcFile, &config)
			infoJson, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resultFile, string(infoJson))
		})
	}
}

func TestStatisticsDefault(t *testing.T) {
	testModule := "statistics"
	for _, fixture := range GetFixtureList() {
		t.Run(fixture, func(t *testing.T) {
			config := statistics.CollectConfig{}.Default()
			srcFile, resultFile := ReadFixturePair(fixture, testModule, "default")
			info := statistics.ParseInfo(srcFile, config)
			infoJson, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resultFile, string(infoJson))
		})
	}
}

func TestStatisticsFull(t *testing.T) {
	testModule := "statistics"
	for _, fixture := range GetFixtureList() {
		t.Run(fixture, func(t *testing.T) {
			config := statistics.CollectConfig{
				General:  true,
				PerQueue: true,
			}
			srcFile, resultFile := ReadFixturePair(fixture, testModule, "full")
			info := statistics.ParseInfo(srcFile, &config)
			infoJson, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, resultFile, string(infoJson))
		})
	}
}
