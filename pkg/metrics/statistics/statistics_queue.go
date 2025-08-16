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
		`rx_bytes\[([0-9]+)\]`,
	},
	"tx_bytes": {
		"tx-([0-9]+).bytes",
		"tx_queue_([0-9]+)_bytes",
		"tx-([0-9]+).tx_bytes",
		`tx_bytes\[([0-9]+)\]`,
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
				regexpLogger := Logger.With("regexp", metricRegexp.String(), "regexpGroup", metricRegexpName, "metricName", srcMetricName)
				matchedMetricRegexp := metricRegexp.FindAllStringSubmatch(srcMetricName, -1)
				if matchedMetricRegexp == nil {
					continue
				}
				if regexpMatched {
					regexpLogger.Error("Queued metric has more than one match, some regexps are overlapping, skipping")
				}
				regexpMatched = true
				regexpLogger.Debug("Metric matches pattern")
				if len(matchedMetricRegexp) > 1 {
					regexpLogger.Error("Regexp matched more than once, regexp or metric is broken, skipping")
					continue
				}

				if len(matchedMetricRegexp[0]) != 2 {
					regexpLogger.Error("Regexp first match does not have 2 matches, regexp or metric is broken, skipping")
					continue
				}
				// We expect to have 1 match, that's why we taking [0] from matchedMetricRegexp
				// and we need the first capture group, which is always second, that's why [1]
				metricIndexString := matchedMetricRegexp[0][1]
				metricIndex64, err := strconv.ParseInt(metricIndexString, 10, 64)
				if err != nil {
					regexpLogger.Error("Cannot parse queue index, skipping", "error", err)
					continue
				}
				metricIndex := int(metricIndex64)
				regexpLogger.Debug("Metric has index", "index", metricIndex)

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
		stats.RxBytes = common.SumFieldsFloat64([]*float64{
			stats.RxUcastBytes,
			stats.RxMcastBytes,
			stats.RxBcastBytes,
		})
	}

	if stats.TxBytes == nil {
		stats.TxBytes = common.SumFieldsFloat64([]*float64{
			stats.TxUcastBytes,
			stats.TxMcastBytes,
			stats.TxBcastBytes,
		})
	}
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
