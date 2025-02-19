package generic_info

type GenericInfo struct {
	SupportedSettings  *AvaliableSettings
	AdvertisedSettings *AvaliableSettings
	Settings           *Settings
}

// Cut off `Supported ` prefix
type AvaliableSettings struct {
	LinkModes     []string `generic_info_avaliable_settings:"link modes"`
	PauseFrameUse string   `generic_info_avaliable_settings:"pause frame use"`
	FecModes      string   `generic_info_avaliable_settings:"FEC modes"`
}

type Settings struct {
	Speed           string
	SpeedBytes      uint64
	Duplex          string
	Port            string
	Transceiver     string
	AutoNegotiation bool `generic_info_settings:"Auto-negotiation"`
	LinkDetected    bool `generic_info_settings:"Link detected"`
}
