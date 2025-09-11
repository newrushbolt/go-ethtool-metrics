// Statistics, eg `ethtool -S ethX`
package statistics

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/newrushbolt/go-ethtool-metrics/common"
)

type queuedMetrics map[int]map[string]string

// These regexps must not overlap
// Multiple matches will cause an error and none of those matches will make it to the final metrics
var rawQueuedRegexps = map[string][]string{
	// General
	"rx_bytes": {
		"rx-([0-9]+)\\.bytes",
		"rx_queue_([0-9]+)_bytes",
		"rx-([0-9]+)\\.rx_bytes",
		`rx_bytes\[([0-9]+)\]`,
	},
	"tx_bytes": {
		"tx-([0-9]+)\\.bytes",
		"tx_queue_([0-9]+)_bytes",
		"tx-([0-9]+)\\.tx_bytes",
		`tx_bytes\[([0-9]+)\]`,
	},
	"rx_drops": {
		`rx([0-9]+)_drops`,
		`rx_queue_([0-9]+)_drops`,
	},
	"rx_drop_cnt": {
		// Not sure if we can treat `rx_queue_drop_cnt[0]` for gve driver as rx_drops,
		// because there are also `rx_drops_packet_over_mru`, `rx_drops_invalid_checksum` and `rx_dropped_pkt`
		// TODO: figure this out, probably by loadtesting with real data and reading gve driver source
		`rx_queue_drop_cnt\[([0-9]+)\]`,
	},
	// PerType
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
	// XDP RX
	"rx_xdp_drops": {
		`rx([0-9]+)_xdp_drops`,
		`rx_queue_([0-9]+)_xdp_drops`,
		`rx_xdp_drop\[([0-9]+)\]`,
	},
	"rx_xdp_tx": {
		`rx([0-9]+)_xdp_tx`,
		`rx_queue_([0-9]+)_xdp_tx`,
		`rx_xdp_tx\[([0-9]+)\]`,
	},
	"rx_xdp_tx_errors": {
		`rx_xdp_tx_errors\[([0-9]+)\]`,
	},
	"rx_xdp_aborted": {
		`rx_xdp_aborted\[([0-9]+)\]`,
	},
	"rx_xdp_pass": {
		`rx_xdp_pass\[([0-9]+)\]`,
	},
	"rx_xdp_redirect": {
		`rx_xdp_redirect\[([0-9]+)\]`,
	},
	"rx_xdp_redirects": {
		`rx([0-9]+)_xdp_redirects`,
		`rx_queue_([0-9]+)_xdp_redirects`,
	},
	"rx_xdp_redirect_errors": {
		`rx_xdp_redirect_errors\[([0-9]+)\]`,
	},
	"rx_xdp_alloc_fails": {
		`rx_xdp_alloc_fails\[([0-9]+)\]`,
	},
	// XDP TX
	"tx_xdp_tx": {
		`tx([0-9]+)_xdp_tx`,
		`tx_queue_([0-9]+)_xdp_tx`,
	},
	"tx_xdp_xmit": {
		`tx_xdp_xmit\[([0-9]+)\]`,
	},
	"tx_xdp_xmit_errors": {
		`tx_xdp_xmit_errors\[([0-9]+)\]`,
	},
	"tx_xdp_tx_drops": {
		`tx([0-9]+)_xdp_tx_drops`,
		`tx_queue_([0-9]+)_xdp_tx_drops`,
	},
}

var queuedRegexps map[string][]*regexp.Regexp

func init() {
	queuedRegexps = compileQueuedRegexps(rawQueuedRegexps)
}

func compileQueuedRegexps(rawQueuedRegexps map[string][]string) map[string][]*regexp.Regexp {
	queuedRegexps := make(map[string][]*regexp.Regexp, len(rawQueuedRegexps))
	for regexpName, regexpStrings := range rawQueuedRegexps {
		var compiledRegexps []*regexp.Regexp
		for _, regexpString := range regexpStrings {
			anchoredRegexpString := fmt.Sprintf("^%s$", regexpString)
			compiledRegex := regexp.MustCompile(anchoredRegexpString)
			compiledRegexps = append(compiledRegexps, compiledRegex)
		}
		queuedRegexps[regexpName] = compiledRegexps
	}
	return queuedRegexps
}

type queueMatchedMetric struct {
	Value   string
	Matches []queueMatchedMetricsMatch
}
type queueMatchedMetricsMatch struct {
	Group   string
	Pattern *regexp.Regexp
	Queue   int
}

func cleanEmptyQueues(metricsMap *queuedMetrics) {
	for queue, metrics := range *metricsMap {
		if len(metrics) == 0 {
			delete(*metricsMap, queue)
		}
	}
}

// This function fetches per-queue metrics due to regexp-rules from rawQueuedRegexps
// If the metric matched more than one regexp, we skip the current match and log an error
func extractQueuedMetrics(srcMetrics map[string]string) queuedMetrics {
	matchedMetrics := map[string]queueMatchedMetric{}
	for srcMetricName, srcMetricvalue := range srcMetrics {
		for metricRegexpName, possibleMetricRegexps := range queuedRegexps {
			for _, metricRegexp := range possibleMetricRegexps {
				regexpLogger := Logger.With("regexp", metricRegexp.String(), "regexpGroup", metricRegexpName, "metricName", srcMetricName)

				matchedMetricRegexp := metricRegexp.FindAllStringSubmatch(srcMetricName, -1)
				if matchedMetricRegexp == nil {
					continue
				}
				regexpLogger.Debug("Metric matches pattern")
				if len(matchedMetricRegexp) > 1 {
					regexpLogger.Error("Regexp matched more than once for one line, regexp or metric is broken, skipping")
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

				newMetricMatches := []queueMatchedMetricsMatch{
					{
						Group:   metricRegexpName,
						Pattern: metricRegexp,
						Queue:   metricIndex,
					},
				}

				var currentMetricObject queueMatchedMetric
				var ok bool
				currentMetricObject, ok = matchedMetrics[srcMetricName]
				if ok {
					currentMetricObject.Matches = append(currentMetricObject.Matches, newMetricMatches...)
				} else {
					currentMetricObject = queueMatchedMetric{
						Value:   srcMetricvalue,
						Matches: newMetricMatches,
					}
				}
				matchedMetrics[srcMetricName] = currentMetricObject

			}
		}
	}

	queuedMetricsMap := queuedMetrics{}
	for metricName, metricObject := range matchedMetrics {
		metricLogger := Logger.With("metricName", metricName, "metricObject", metricObject)
		switch len(metricObject.Matches) {
		case 0:
			metricLogger.Debug("No per-queue regexp matched, skipping metric")
		case 1:
			match := metricObject.Matches[0]
			currentIndexMap := queuedMetricsMap[match.Queue]
			if currentIndexMap == nil {
				newCurrentIndexMap := map[string]string{}
				queuedMetricsMap[match.Queue] = newCurrentIndexMap
			} else if _, ok := queuedMetricsMap[match.Queue][match.Group]; ok {
				delete(queuedMetricsMap[match.Queue], match.Group)
				metricLogger.Error("This metric is already matched by other regexp, skipping it at all. Regexp rules are probably broken.")
				break
			}
			queuedMetricsMap[match.Queue][match.Group] = metricObject.Value
			metricLogger.Debug("Metric added")
		default:
			metricLogger.Error("Metric has more than one match, skipping it at all. Regexp rules are probably broken.")
		}
	}
	cleanEmptyQueues(&queuedMetricsMap)
	return queuedMetricsMap
}

func queueGenerateMissingBytesMetrics(stats *QueueStatistics) {
	if stats.General.RxBytes == nil {
		stats.General.RxBytes = common.SumFieldsFloat64([]*float64{
			stats.PerType.RxUcastBytes,
			stats.PerType.RxMcastBytes,
			stats.PerType.RxBcastBytes,
		})
	}

	if stats.General.TxBytes == nil {
		stats.General.TxBytes = common.SumFieldsFloat64([]*float64{
			stats.PerType.TxUcastBytes,
			stats.PerType.TxMcastBytes,
			stats.PerType.TxBcastBytes,
		})
	}
}

func parseQueuedInfo(statisticsMap map[string]string, config CollectConfig) *PerQueueStatistics {
	allQueuedMetrics := extractQueuedMetrics(statisticsMap)
	perQueueStatistics := make(PerQueueStatistics, len(allQueuedMetrics))
	for queue, queueMetricsMap := range allQueuedMetrics {
		var queueStatistics QueueStatistics

		if config.PerQueueGeneral {
			var general QueueStatisticsGeneral
			common.ParseAbstractDataObject(&queueMetricsMap, &general, "queue_statistics_general")
			queueStatistics.General = &general
		}

		if config.PerQueueXdp {
			var xdp QueueStatisticsXdp
			common.ParseAbstractDataObject(&queueMetricsMap, &xdp, "queue_statistics_xdp")
			queueStatistics.Xdp = &xdp
		}

		if config.PerQueueGenerateMissingBytesMetrics || config.PerQueuePerType {
			var perType QueueStatisticsPerType
			common.ParseAbstractDataObject(&queueMetricsMap, &perType, "queue_statistics_per_type")
			queueStatistics.PerType = &perType
		}
		if config.PerQueueGenerateMissingBytesMetrics {
			queueGenerateMissingBytesMetrics(&queueStatistics)
		}
		// Deleting per-queue metrics if they were only needed to calculate missing general metrics
		if config.PerQueueGenerateMissingBytesMetrics && !config.PerQueuePerType {
			queueStatistics.PerType = nil
		}

		perQueueStatistics[queue] = queueStatistics
	}
	return &perQueueStatistics
}
