package statistics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var statisticsMap = map[string]string{
	"[0]: rx_ucast_bytes": "1316",
	"[0]: rx_mcast_bytes": "1613",
	"[0]: rx_bcast_bytes": "2",
	"[0]: tx_ucast_bytes": "1613",
	"[0]: tx_mcast_bytes": "1316",
	"[0]: tx_bcast_bytes": "1",
}

func ptr[T any](v T) *T { return &v }

func TestPerTypeBytesDeletion(t *testing.T) {
	config := CollectConfig{
		PerQueue:                            true,
		PerQueueGenerateMissingBytesMetrics: true,
		PerQueuePerTypeBytes:                false,
	}

	expectedParseResult := PerQueueStatistics{
		QueueStatistics{
			TxBytes:      ptr(2930.0),
			RxBytes:      ptr(2931.0),
			RxUcastBytes: nil,
			RxMcastBytes: nil,
			RxBcastBytes: nil,
			TxUcastBytes: nil,
			TxMcastBytes: nil,
			TxBcastBytes: nil,
		},
	}

	parseResult := parseQueuedInfo(statisticsMap, config)
	assert.Equal(t, &expectedParseResult, parseResult)
}

func TestPerTypeBytesKeep(t *testing.T) {
	config := CollectConfig{
		PerQueue:                            true,
		PerQueueGenerateMissingBytesMetrics: true,
		PerQueuePerTypeBytes:                true,
	}

	expectedParseResult := PerQueueStatistics{
		QueueStatistics{
			TxBytes:      ptr(2930.0),
			RxBytes:      ptr(2931.0),
			RxUcastBytes: ptr(1316.0),
			RxMcastBytes: ptr(1613.0),
			RxBcastBytes: ptr(2.0),
			TxUcastBytes: ptr(1613.0),
			TxMcastBytes: ptr(1316.0),
			TxBcastBytes: ptr(1.0),
		},
	}

	parseResult := parseQueuedInfo(statisticsMap, config)
	assert.Equal(t, &expectedParseResult, parseResult)
}
