// Module info, eg `ethtool -m ethX`
package module_info

import "github.com/newrushbolt/go-ethtool-metrics/internal"

func _ParseVendorInfo(rawInfo string) *VendorInfo {
	vendorInfoMap, _ := internal.ParseAbstractColonData(rawInfo, "Vendor", false)
	var vendorInfo VendorInfo
	internal.ParseAbstractDataObject(&vendorInfoMap, &vendorInfo, "vendor")
	return &vendorInfo
}

func _ParseDiagnosticsValues(rawInfo string) *DiagnosticsValues {
	diagnosticsValuesMap, _ := internal.ParseAbstractColonData(rawInfo, "", false)
	var diagnosticsValues DiagnosticsValues
	internal.ParseAbstractDataObject(&diagnosticsValuesMap, &diagnosticsValues, "diag_values")
	return &diagnosticsValues
}

func _ParseDiagnosticsAlarms(rawInfo string) *DiagnosticsAlarms {
	diagnosticsAlarmsMap, _ := internal.ParseAbstractColonData(rawInfo, "", false)
	var diagnosticsAlarms DiagnosticsAlarms
	internal.ParseAbstractDataObject(&diagnosticsAlarmsMap, &diagnosticsAlarms, "diag_alarms")
	return &diagnosticsAlarms
}

func _ParseDiagnosticsWarnings(rawInfo string) *DiagnosticsWarnings {
	diagnosticsWarningsMap, _ := internal.ParseAbstractColonData(rawInfo, "", false)
	var diagnosticsWarnings DiagnosticsWarnings
	internal.ParseAbstractDataObject(&diagnosticsWarningsMap, &diagnosticsWarnings, "diag_warnings")
	return &diagnosticsWarnings
}

func ParseInfo(rawInfo string, config *ModuleInfoConfig) *ModuleInfo {
	// Empty string means we got an error getting raw info from ethtool
	// This is pretty common for module info `ethtool -m ethX`
	// TODO: better data detection
	// Maybe we should look into the common info `ethtool ethX`
	// and decide weather we should even try to get module info
	if rawInfo == "" {
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
