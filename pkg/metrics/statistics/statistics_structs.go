package statistics

type CollectConfig struct {
	General  bool
	PerQueue bool
}

func (config CollectConfig) Default() *CollectConfig {
	return &CollectConfig{
		General:  true,
		PerQueue: false,
	}
}

type StatisticsInfo struct {
	General  *GeneralStatistics
	PerQueue *PerQueueStatistics
}

type PerQueueStatistics []QueueStatistics
type QueueStatistics struct {
	// TODO: move to *float64 later
	TxBytes float64 `queue_statistics:"tx_bytes"`
	RxBytes float64 `queue_statistics:"rx_bytes"`
	// Needs bnxt_en and Nan support first. Fix after 0.0.4
	// We also need to calulate RxBytes|TxBytes for bnxt_en by summing all types of bytes
	// RxUcastBytes *float64 `queue_statistics:"rx_ucast_bytes"`
	// RxMcastBytes *float64 `queue_statistics:"rx_mcast_bytes"`
	// RxBcastBytes *float64 `queue_statistics:"rx_bcast_bytes"`
	// TxUcastBytes *float64 `queue_statistics:"tx_ucast_bytes"`
	// TxMcastBytes *float64 `queue_statistics:"tx_mcast_bytes"`
	// TxBcastBytes *float64 `queue_statistics:"tx_bcast_bytes"`
	// TpaBytes     *float64 `queue_statistics:"tpa_bytes"`
}

type GeneralStatistics struct {
	TxBytes uint64 `general_statistics:"tx_bytes"`
	RxBytes uint64 `general_statistics:"rx_bytes"`

	RxErrors uint64 `general_statistics:"rx_errors"` // bnxt:missing
	TxErrors uint64 `general_statistics:"tx_errors,tx_err"`

	RxDiscards   uint64 `general_statistics:"rx_discards,veb.rx_discards,rx_stat_discard"`
	TxDiscards   uint64 `general_statistics:"tx_discards,veb.tx_discards,tx_stat_discard"`
	TxCollisions uint64 `general_statistics:"tx_collisions,tx_total_collisions,collisions"`
	RxCrcErrors  uint64 `general_statistics:"rx_crc_errors,rx_crc_errors.nic"` // Only exists in Intel
}

// Other possible variables
// TODO: Add testdata for virtio
// rx_queue_0_packets: 99116794
// rx_queue_0_bytes: 708784361629
// rx_queue_0_drops: 0
// rx_queue_0_xdp_packets: 0
// rx_queue_0_xdp_tx: 0
// rx_queue_0_xdp_redirects: 0
// rx_queue_0_xdp_drops: 0
// rx_queue_0_kicks: 4060

// tx_queue_0_packets: 91340287
// tx_queue_0_bytes: 1484407675445
// tx_queue_0_xdp_tx: 0
// tx_queue_0_xdp_tx_drops: 0
// tx_queue_0_kicks: 89560041

// Only exists in bnxt, should probably have separate struct
// tx_pause_frames: 0
// rx_pause_frames: 0
