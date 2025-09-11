package statistics

import (
	"reflect"
	"regexp"
	"strings"
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
		PerQueueGeneral:                     true,
		PerQueueGenerateMissingBytesMetrics: true,
		PerQueuePerType:                     false,
	}

	expectedParseResult := PerQueueStatistics{
		QueueStatistics{
			General: &QueueStatisticsGeneral{
				TxBytes: ptr(2930.0),
				RxBytes: ptr(2931.0),
				RxDrops: nil,
			},
			PerType: nil,
		},
	}

	parseResult := parseQueuedInfo(statisticsMap, config)
	assert.Equal(t, &expectedParseResult, parseResult)
}

func TestPerTypeBytesKeep(t *testing.T) {
	config := CollectConfig{
		PerQueueGeneral:                     true,
		PerQueueGenerateMissingBytesMetrics: true,
		PerQueuePerType:                     true,
	}

	expectedParseResult := PerQueueStatistics{
		QueueStatistics{
			General: &QueueStatisticsGeneral{
				TxBytes: ptr(2930.0),
				RxBytes: ptr(2931.0),
			},
			PerType: &QueueStatisticsPerType{
				RxUcastBytes: ptr(1316.0),
				RxMcastBytes: ptr(1613.0),
				RxBcastBytes: ptr(2.0),
				TxUcastBytes: ptr(1613.0),
				TxMcastBytes: ptr(1316.0),
				TxBcastBytes: ptr(1.0),
			},
		},
	}

	parseResult := parseQueuedInfo(statisticsMap, config)
	assert.Equal(t, &expectedParseResult, parseResult)
}

func TestExtractQueuedMetricsMultipleRegexMatch(t *testing.T) {
	metrics := map[string]string{
		"rx-0.bytes":       "100",
		"rx_queue_0_bytes": "200",
	}
	result := extractQueuedMetrics(metrics)
	assert.Equal(t, queuedMetrics{}, result)
}

func TestExtractQueuedMetricsInvalidQueueIndex(t *testing.T) {
	metrics := map[string]string{
		"rx-abc.bytes": "400",
	}
	result := extractQueuedMetrics(metrics)
	assert.Equal(t, queuedMetrics{}, result)
}

func TestExtractQueuedMetricsNoMatch(t *testing.T) {
	metrics := map[string]string{
		"not_a_queue_metric": "999",
	}
	result := extractQueuedMetrics(metrics)
	assert.Equal(t, queuedMetrics{}, result)
}

func TestExtractQueuedMetricsOverlappingRegexps(t *testing.T) {
	// backup original and restore regexps at the end
	orig := queuedRegexps
	defer func() { queuedRegexps = orig }()
	queuedRegexps = map[string][]*regexp.Regexp{
		"tx_bytes": {
			regexp.MustCompile(`tx[0-9a-z]+_([0-9])[0-9a-z]+`),
			regexp.MustCompile(`tx([0-9])_[0-9a-z]+`),
		},
		"tx_bytes_wrong": {
			regexp.MustCompile(`tx([0-9])_[0-9a-z]+`),
		},
	}

	metrics := map[string]string{"tx0_0bytes": "123"}
	result := extractQueuedMetrics(metrics)
	assert.Equal(t, queuedMetrics{}, result)
}

func TestQueueStatisticsStructTagsHaveRegexps(t *testing.T) {
	findMissingTags := func(tagPrefix string, structType reflect.Type) (missingTags []string) {
		for i := 0; i < structType.NumField(); i++ {
			f := structType.Field(i)
			tagListString := f.Tag.Get(tagPrefix)
			if tagListString == "" {
				continue
			}
			tagList := strings.Split(tagListString, ",")
			for _, tag := range tagList {
				tag = strings.TrimSpace(tag)
				if tag == "" {
					continue
				}
				if _, ok := rawQueuedRegexps[tag]; !ok {
					missingTags = append(missingTags, tag)
				}
			}
		}
		return
	}

	assert.Empty(t, findMissingTags("queue_statistics_general", reflect.TypeOf(QueueStatisticsGeneral{})))
	assert.Empty(t, findMissingTags("queue_statistics_per_type", reflect.TypeOf(QueueStatisticsPerType{})))
	assert.Empty(t, findMissingTags("queue_statistics_xdp", reflect.TypeOf(QueueStatisticsXdp{})))
}
