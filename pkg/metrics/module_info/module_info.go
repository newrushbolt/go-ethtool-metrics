// Module info, eg `ethtool -m ethX`
package module_info

import (
	"log/slog"

	"github.com/newrushbolt/go-ethtool-metrics/internal"
)

func _ParseVendorInfo(rawInfo string) *VendorInfo {
	var vendorInfo VendorInfo
	vendorInfoMap := internal.ParseAbstractColonData(rawInfo, "Vendor", false)
	internal.ParseAbstractDataObject(&vendorInfoMap, &vendorInfo, "vendor")
	return &vendorInfo
}

func _ParseDiagnosticsValues(rawInfo string) *DiagnosticsValues {
	var diagnosticsValues DiagnosticsValues
	diagnosticsValuesMap := internal.ParseAbstractColonData(rawInfo, "", false)
	internal.ParseAbstractDataObject(&diagnosticsValuesMap, &diagnosticsValues, "diag_values")
	return &diagnosticsValues
}

func _ParseDiagnosticsAlarms(rawInfo string) *DiagnosticsAlarms {
	var diagnosticsAlarms DiagnosticsAlarms
	diagnosticsAlarmsMap := internal.ParseAbstractColonData(rawInfo, "", false)
	internal.ParseAbstractDataObject(&diagnosticsAlarmsMap, &diagnosticsAlarms, "diag_alarms")
	return &diagnosticsAlarms
}

func _ParseDiagnosticsWarnings(rawInfo string) *DiagnosticsWarnings {
	var diagnosticsWarnings DiagnosticsWarnings
	diagnosticsWarningsMap := internal.ParseAbstractColonData(rawInfo, "", false)
	internal.ParseAbstractDataObject(&diagnosticsWarningsMap, &diagnosticsWarnings, "diag_warnings")
	return &diagnosticsWarnings
}

func ParseInfo(rawInfo string, config *CollectConfig) *ModuleInfo {
	slog.SetLogLoggerLevel(internal.GetLogLevel())
	// Empty string means we got an error getting raw info from ethtool
	// This is pretty common for module info `ethtool -m ethX`
	//
	// TODO: better data detection
	// Maybe we should look into the common info `ethtool ethX`
	// and decide whether we should even try to get module info
	if rawInfo == "" {
		slog.Info("Module got empty ethtool data, skipping", "module", "ModuleInfo")
		return nil
	}

	var alarms *DiagnosticsAlarms
	if config.CollectDiagnosticsAlarms {
		alarms = _ParseDiagnosticsAlarms(rawInfo)
	}

	var values *DiagnosticsValues
	if config.CollectDiagnosticsValues {
		values = _ParseDiagnosticsValues(rawInfo)
	}

	var warnings *DiagnosticsWarnings
	if config.CollectDiagnosticsWarnings {
		warnings = _ParseDiagnosticsWarnings(rawInfo)
	}

	var diagnostics Diagnostics
	if (alarms != nil) || (values != nil) || (warnings != nil) {
		diagnostics = Diagnostics{
			Values:   values,
			Alarms:   alarms,
			Warnings: warnings,
		}
	}

	var vendorInfo *VendorInfo
	if config.CollectVendor {
		vendorInfo = _ParseVendorInfo(rawInfo)
	}
	moduleInfo := ModuleInfo{
		Vendor:      vendorInfo,
		Diagnostics: &diagnostics,
	}
	return &moduleInfo
}
