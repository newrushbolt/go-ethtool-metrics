package module_info

type ModuleInfo struct {
	Vendor      *VendorInfo
	Diagnostics *Diagnostics
}

type CollectConfig struct {
	CollectDiagnosticsAlarms   bool
	CollectDiagnosticsValues   bool
	CollectDiagnosticsWarnings bool
	CollectVendor              bool
}

func (config CollectConfig) Default() *CollectConfig {
	return &CollectConfig{
		CollectDiagnosticsAlarms:   true,
		CollectDiagnosticsValues:   false,
		CollectDiagnosticsWarnings: false,
		CollectVendor:              false,
	}
}

type VendorInfo struct {
	Name         string `vendor:"name"`
	OUI          string
	PartNumber   string `vendor:"PN"`
	Revision     string `vendor:"rev"`
	SerialNumber string `vendor:"SN"`
}

type Diagnostics struct {
	Values   *DiagnosticsValues
	Alarms   *DiagnosticsAlarms
	Warnings *DiagnosticsWarnings
}

type DiagnosticsAlarms struct {
	BiasHigh        bool `diag_alarms:"Laser bias current high alarm"`
	BiasLow         bool `diag_alarms:"Laser bias current low alarm"`
	OutputPowerHigh bool `diag_alarms:"Laser output power high alarm"`
	OutputLow       bool `diag_alarms:"Laser output power low alarm"`
	TemperatureHigh bool `diag_alarms:"Module temperature high alarm"`
	TemperatureLow  bool `diag_alarms:"Module temperature low alarm"`
	VoltageHigh     bool `diag_alarms:"Module voltage high alarm"`
	VoltageLow      bool `diag_alarms:"Module voltage low alarm"`
	InputPowerHigh  bool `diag_alarms:"Laser rx power high alarm"`
	InputPowerLow   bool `diag_alarms:"Laser rx power low alarm"`
}

type DiagnosticsWarnings struct {
	BiasHigh        bool `diag_warnings:"Laser bias current high warning"`
	BiasLow         bool `diag_warnings:"Laser bias current low warning"`
	OutputPowerHigh bool `diag_warnings:"Laser output power high warning"`
	OutputLow       bool `diag_warnings:"Laser output power low warning"`
	TemperatureHigh bool `diag_warnings:"Module temperature high warning"`
	TemperatureLow  bool `diag_warnings:"Module temperature low warning"`
	VoltageHigh     bool `diag_warnings:"Module voltage high warning"`
	VoltageLow      bool `diag_warnings:"Module voltage low warning"`
	InputPowerHigh  bool `diag_warnings:"Laser rx power high warning"`
	InputPowerLow   bool `diag_warnings:"Laser rx power low warning"`
}

type DiagnosticsValues struct {
	BiasMilliAmps         float32 `diag_values:"Laser bias current"`
	OutputPowerMilliWatts float32 `diag_values:"Laser output power"`
	InputPowerMilliWatts  float32 `diag_values:"Receiver signal average optical power"`
	TemperatureCelcius    float32 `diag_values:"Module temperature"`
	Voltage               float32 `diag_values:"Module voltage"`
}
