// Statistics, eg `ethtool -S ethX`
package statistics

import (
	"log/slog"
	"regexp"
	"strconv"

	"github.com/newrushbolt/go-ethtool-metrics/internal"
)

type queuedMetrics map[int]map[string]string

// func calculateBroadcomPerQueueBytes

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

func extractQueuedMetricsV2(srcMetrics map[string]string) queuedMetrics {
	rawQueuedRegexps := map[string][]string{
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

	queuedRegexps := compileQueuedRegexps(rawQueuedRegexps)
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
					slog.Error("Queued metric has more than one match, some regexps are overlapping", "metric", srcMetricName, "patternRegexp", metricRegexp.String(), "pattern", metricRegexpName)
				}
				regexpMatched = true
				slog.Debug("Metric matches pattern", "metric", srcMetricName, "patternRegexp", metricRegexp.String(), "pattern", metricRegexpName)
				// TODO: check if only one result matched
				// switch len(matchedMetricRegexp){
				// case 1:

				// We expect to have 1 match, that's why we taking [0] from matchedMetricRegexp
				// and we need the first capture group, which is always second, that's why [1]
				metricIndexString := matchedMetricRegexp[0][1]
				metricIndex64, _ := strconv.ParseInt(metricIndexString, 10, 64)
				metricIndex := int(metricIndex64)
				slog.Debug("Metric has index", "metric", srcMetricName, "index", metricIndex)

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

// func calculateBroadcomPerQueueBytes
// ./testdata/broadcom/bnxt_en/00_sfp_10gsr85/src/statistics:7:		[0]: rx_ucast_bytes: 594284585291
// ./testdata/broadcom/bnxt_en/00_sfp_10gsr85/src/statistics:8:		[0]: rx_mcast_bytes: 226857008
// ./testdata/broadcom/bnxt_en/00_sfp_10gsr85/src/statistics:9:		[0]: rx_bcast_bytes: 1857840
// ./testdata/broadcom/bnxt_en/00_sfp_10gsr85/src/statistics:15:		[0]: tx_ucast_bytes: 1790928332018
// ./testdata/broadcom/bnxt_en/00_sfp_10gsr85/src/statistics:16:		[0]: tx_mcast_bytes: 86
// ./testdata/broadcom/bnxt_en/00_sfp_10gsr85/src/statistics:17:		[0]: tx_bcast_bytes: 0
// ./testdata/broadcom/bnxt_en/00_sfp_10gsr85/src/statistics:19:		[0]: tpa_bytes:      94005443409

func parseQueuedInfo(statisticsMap map[string]string) *PerQueueStatistics {
	allQueuedMetrics := extractQueuedMetricsV2(statisticsMap)
	perQueueStatistics := make(PerQueueStatistics, len(allQueuedMetrics))
	for queue, queueMetricsMap := range allQueuedMetrics {
		var queueStatistics QueueStatistics
		internal.ParseAbstractDataObject(&queueMetricsMap, &queueStatistics, "queue_statistics")
		perQueueStatistics[queue] = queueStatistics
	}
	return &perQueueStatistics
}

func ParseInfo(rawInfo string, config *CollectConfig) *StatisticsInfo {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if rawInfo == "" {
		slog.Info("Module got empty ethtool data, skipping", "module", "Statistics")
		return nil
	}

	statistics := StatisticsInfo{}
	generalStatisticsMap, _ := internal.ParseAbstractColonData(rawInfo, "", true)

	if config.PerQueue {
		statistics.PerQueue = parseQueuedInfo(generalStatisticsMap)
	}

	if config.General {
		var generalStatistics GeneralStatistics
		internal.ParseAbstractDataObject(&generalStatisticsMap, &generalStatistics, "general_statistics")
		statistics.General = &generalStatistics
	}
	return &statistics
}
