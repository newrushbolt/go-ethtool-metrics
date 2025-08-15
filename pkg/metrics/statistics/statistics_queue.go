// Statistics, eg `ethtool -S ethX`
package statistics

import (
	"regexp"
	"strconv"

	"github.com/newrushbolt/go-ethtool-metrics/common"
)

type queuedMetrics map[int]map[string]string

var rawQueuedRegexps = map[string][]string{
	"rx_bytes": {
		"rx-([0-9]+).bytes",
		"rx_queue_([0-9]+)_bytes",
		"rx-([0-9]+).rx_bytes",
	},
	"tx_bytes": {
		"tx-([0-9]+).bytes",
		"tx_queue_([0-9]+)_bytes",
		"tx-([0-9]+).tx_bytes",
	},
	"rx_ucast_bytes": {
		`\[([0-9]+)\]: rx_ucast_bytes`,
	},
	"rx_mcast_bytes": {
		`\[([0-9]+)\]: rx_mcast_bytes`,
	},
	"rx_bcast_bytes": {
		`\[([0-9]+)\]: rx_bcast_bytes`,
	},
	"tx_ucast_bytes": {
		`\[([0-9]+)\]: tx_ucast_bytes`,
	},
	"tx_mcast_bytes": {
		`\[([0-9]+)\]: tx_mcast_bytes`,
	},
	"tx_bcast_bytes": {
		`\[([0-9]+)\]: tx_bcast_bytes`,
	},
	"tpa_bytes": {
		`\[([0-9]+)\]: tpa_bytes`,
	},
}

var queuedRegexps map[string][]*regexp.Regexp

func init() {
	queuedRegexps = compileQueuedRegexps(rawQueuedRegexps)
}

func compileQueuedRegexps(rawQueuedRegexps map[string][]string) map[string][]*regexp.Regexp {
	queuedRegexps := make(map[string][]*regexp.Regexp, len(rawQueuedRegexps))
	for regexName, regexStrings := range rawQueuedRegexps {
		var compiledRegexps []*regexp.Regexp
		for _, regexString := range regexStrings {
			compiledRegex := regexp.MustCompile(regexString)
			compiledRegexps = append(compiledRegexps, compiledRegex)
		}
		queuedRegexps[regexName] = compiledRegexps
	}
	return queuedRegexps
}

func extractQueuedMetrics(srcMetrics map[string]string) queuedMetrics {
	queuedMetricsMap := queuedMetrics{}
	for srcMetricName, srcMetricvalue := range srcMetrics {
		for metricRegexpName, possibleMetricRegexps := range queuedRegexps {
			regexpMatched := false
			for _, metricRegexp := range possibleMetricRegexps {
				matchedMetricRegexp := metricRegexp.FindAllStringSubmatch(srcMetricName, -1)
				if matchedMetricRegexp == nil {
					continue
				}
				if regexpMatched {
					Logger.Error("Queued metric has more than one match, some regexps are overlapping", "metric", srcMetricName, "patternRegexp", metricRegexp.String(), "pattern", metricRegexpName)
				}
				regexpMatched = true
				Logger.Debug("Metric matches pattern", "metric", srcMetricName, "patternRegexp", metricRegexp.String(), "pattern", metricRegexpName)
				// TODO: check if only one result matched
				// switch len(matchedMetricRegexp){
				// case 1:

				// We expect to have 1 match, that's why we taking [0] from matchedMetricRegexp
				// and we need the first capture group, which is always second, that's why [1]
				metricIndexString := matchedMetricRegexp[0][1]
				metricIndex64, err := strconv.ParseInt(metricIndexString, 10, 64)
				if err != nil {
					continue
				}
				metricIndex := int(metricIndex64)
				Logger.Debug("Metric has index", "metric", srcMetricName, "index", metricIndex)

				currentIndexMap := queuedMetricsMap[metricIndex]
				if currentIndexMap == nil {
					newCurrentIndexMap := map[string]string{
						metricRegexpName: srcMetricvalue,
					}
					queuedMetricsMap[metricIndex] = newCurrentIndexMap
					continue
				}
				currentIndexMap[metricRegexpName] = srcMetricvalue
			}
		}
	}
	return queuedMetricsMap
}

func queueRemovePerTypeBytes(stats *QueueStatistics) {
	stats.TxUcastBytes = nil
	stats.TxMcastBytes = nil
	stats.TxBcastBytes = nil
	stats.RxUcastBytes = nil
	stats.RxMcastBytes = nil
	stats.RxBcastBytes = nil
}

func queueGenerateMissingBytesMetrics(stats *QueueStatistics) {
	if stats.RxBytes == nil {
		stats.RxBytes = sumBytesFields([]*float64{
			stats.RxUcastBytes,
			stats.RxMcastBytes,
			stats.RxBcastBytes,
		})
	}

	if stats.TxBytes == nil {
		stats.TxBytes = sumBytesFields([]*float64{
			stats.TxUcastBytes,
			stats.TxMcastBytes,
			stats.TxBcastBytes,
		})
	}
}

func sumBytesFields(fields []*float64) *float64 {
	var sum float64
	var exists bool
	for _, value := range fields {
		if value != nil {
			sum += *value
			exists = true
		}
	}
	if exists {
		return &sum
	}
	return nil
}

func parseQueuedInfo(statisticsMap map[string]string, config CollectConfig) *PerQueueStatistics {
	allQueuedMetrics := extractQueuedMetrics(statisticsMap)
	perQueueStatistics := make(PerQueueStatistics, len(allQueuedMetrics))
	for queue, queueMetricsMap := range allQueuedMetrics {
		var queueStatistics QueueStatistics
		common.ParseAbstractDataObject(&queueMetricsMap, &queueStatistics, "queue_statistics")
		if config.PerQueueGenerateMissingBytesMetrics {
			queueGenerateMissingBytesMetrics(&queueStatistics)
		}
		if !config.PerQueuePerTypeBytes {
			queueRemovePerTypeBytes(&queueStatistics)
		}

		perQueueStatistics[queue] = queueStatistics
	}
	return &perQueueStatistics
}
