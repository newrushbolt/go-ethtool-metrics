package statistics

type StatisticsInfo struct {
	General GeneralStatisticsInfo
}

type GeneralStatisticsInfo struct {
	TxBytes uint64 `general_statistics:"tx_bytes"`
	RxBytes uint64 `general_statistics:"rx_bytes"`

	RxErrors uint64 `general_statistics:"rx_errors"` // bnxt:missing
	TxErrors uint64 `general_statistics:"tx_errors,tx_err"`

	RxDiscards   uint64 `general_statistics:"rx_discards,veb.rx_discards,rx_stat_discard"`
	TxDiscards   uint64 `general_statistics:"tx_discards,veb.tx_discards,tx_stat_discard"`
	TxCollisions uint64 `general_statistics:"tx_collisions,tx_total_collisions,collisions"`
	RxCrcErrors  uint64 `general_statistics:"rx_crc_errors,rx_crc_errors.nic"` // Only exists in Intel
}

// Only exists in bnxt, should probably have separate struct
// tx_pause_frames: 0
// rx_pause_frames: 0

// Only exists in intel, should probably have separate struct
// rx_dropped: 0
// tx_dropped: 0
// rx_unknown_protocol
// rx_alloc_fail: 0
// rx_pg_alloc_fail: 0
// port.tx_dropped_link_down: 0
// port.illegal_bytes: 0
// port.mac_local_faults: 2
// port.mac_remote_faults: 2
// port.tx_timeout: 0
// port.rx_length_errors: 0
// port.rx_oversize: 0

// Only exists in i40
// WTF, rx_csum_bad is different from rx_crc_errors
// testdata/intel/i40e/02_sfp_10_or_25g_sr/src/statistics:563:     rx_csum_bad.nic: 4757
// testdata/intel/i40e/00_sfp_10g_sr85/src/statistics:601:     port.rx_csum_bad: 1778
// testdata/intel/i40e/01_int_tp/src/statistics:232:     port.rx_csum_bad: 0
