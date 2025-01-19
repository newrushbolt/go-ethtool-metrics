package driver_info

type CollectConfig struct {
	DriverFeatures bool
}

func (config CollectConfig) Default() *CollectConfig {
	return &CollectConfig{
		DriverFeatures: false,
	}
}

type DriverInfo struct {
	DriverName           string   `driver:"driver"`
	DriverVersion        string   `driver:"version"`
	FirmwareVersion      string   `driver:"firmware-version"`
	FirmwareVersionParts []string `driver:"firmware-version"`
	BusAddress           string   `driver:"bus-info"`
	Features             *DriverFeatures
}

type DriverFeatures struct {
	EepromAccess bool `driver_supports:"supports-eeprom-access"`
	PrivFlags    bool `driver_supports:"supports-priv-flags"`
	RegisterDump bool `driver_supports:"supports-register-dump"`
	Statistics   bool `driver_supports:"supports-statistics"`
	Test         bool `driver_supports:"supports-test"`
}
