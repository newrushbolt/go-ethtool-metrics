package statistics

type CollectConfig struct {
	General  bool
	PerQueue bool
	// Some drivers (at least bnxt_en) does not provide per-queue tx_bytes\rx_bytes metrics.
	// This flag allows us to generate them using formula `bytes=ucast_bytes+mcast_bytes+bcast_bytes`
	PerQueueGenerateMissingBytesMetrics bool
	PerQueuePerTypeBytes                bool
}

func (config CollectConfig) Default() *CollectConfig {
	return &CollectConfig{
		General:                             true,
		PerQueue:                            false,
		PerQueueGenerateMissingBytesMetrics: true,
		PerQueuePerTypeBytes:                false,
	}
}

type StatisticsInfo struct {
	General  *GeneralStatistics
	PerQueue *PerQueueStatistics
}

type PerQueueStatistics []QueueStatistics
type QueueStatistics struct {
	TxBytes      *float64 `queue_statistics:"tx_bytes"`
	RxBytes      *float64 `queue_statistics:"rx_bytes"`
	RxUcastBytes *float64 `queue_statistics:"rx_ucast_bytes"`
	RxMcastBytes *float64 `queue_statistics:"rx_mcast_bytes"`
	RxBcastBytes *float64 `queue_statistics:"rx_bcast_bytes"`
	TxUcastBytes *float64 `queue_statistics:"tx_ucast_bytes"`
	TxMcastBytes *float64 `queue_statistics:"tx_mcast_bytes"`
	TxBcastBytes *float64 `queue_statistics:"tx_bcast_bytes"`
	// Other "yet to be supported metrics"
	// First prio
	// [0]: tx_errors: 0
	// [0]: tx_discards: 0
	// [0]: rx_discards: 0
	// [0]: rx_errors: 0
	// Second prio
	// [0]: rx_ucast_packets: 276347403
	// [0]: rx_mcast_packets: 1829492
	// [0]: rx_bcast_packets: 30964
	// [0]: tx_ucast_packets: 5860187670
	// [0]: tx_mcast_packets: 1
	// [0]: tx_bcast_packets: 0
	// [0]: tpa_packets: 42605077
	// [0]: tpa_bytes: 94005443409
	// [0]: tpa_events: 19894790
	// [0]: tpa_aborts: 8731049
	// [0]: rx_l4_csum_errors: 0
	// [0]: rx_resets: 0
	// [0]: rx_buf_errors: 0
}

type GeneralStatistics struct {
	TxBytes *float64 `general_statistics:"tx_bytes"`
	RxBytes *float64 `general_statistics:"rx_bytes"`

	RxErrors *float64 `general_statistics:"rx_errors"` // bnxt:missing
	TxErrors *float64 `general_statistics:"tx_errors,tx_err"`

	RxDiscards   *float64 `general_statistics:"rx_discards,veb.rx_discards,rx_stat_discard"`
	TxDiscards   *float64 `general_statistics:"tx_discards,veb.tx_discards,tx_stat_discard"`
	TxCollisions *float64 `general_statistics:"tx_collisions,tx_total_collisions,collisions"`
	RxCrcErrors  *float64 `general_statistics:"rx_crc_errors,rx_crc_errors.nic"` // Only exists in Intel
}

// Other possible variables
// TODO: Add testdata for "virtio"

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
