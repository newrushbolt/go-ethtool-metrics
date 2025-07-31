package generic_info

type CollectConfig struct {
	CollectAdvertisedSettings bool
	CollectSupportedSettings  bool
	CollectSettings           bool
}

func (config CollectConfig) Default() *CollectConfig {
	return &CollectConfig{
		CollectAdvertisedSettings: false,
		CollectSupportedSettings:  false,
		CollectSettings:           true,
	}
}

type GenericInfo struct {
	SupportedSettings  *AvaliableSettings
	AdvertisedSettings *AvaliableSettings
	Settings           *Settings
}

type AvaliableSettings struct {
	LinkModes     []string `generic_info_avaliable_settings:"link modes"`
	PauseFrameUse string   `generic_info_avaliable_settings:"pause frame use"`
	FecModes      string   `generic_info_avaliable_settings:"FEC modes"`
}

type Settings struct {
	Speed      string
	SpeedBytes *float64
	// There is a silent conflict of units between
	// prometheus-oriented products (eg node_exporter), using BYTES,
	// and network-oriented tools (eg ethtool), using BITS
	//
	// Both worlds have their reasons to do so,
	// so for ease of use I decided to present both metrics
	//
	// If having 2 units bothers you in terms of storing extra metrics or any other way,
	// feel free to drop excessive metric on metric_relabel_config
	SpeedBits *float64

	Duplex          string
	Port            string
	Transceiver     string
	AutoNegotiation bool `generic_info_settings:"Auto-negotiation"`
	LinkDetected    bool `generic_info_settings:"Link detected"`
}
