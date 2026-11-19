// Copyright 2026 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux && !noslingshot

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	slingshotTelemetryMetricsInclude = kingpin.Flag(
		"collector.slingshot.telemetry-metrics-include",
		"Regexp of telemetry metric names to include from /sys/class/cxi/<dev>/device/telemetry.",
	).Default(".*").String()
	slingshotTelemetryMetricsExclude = kingpin.Flag(
		"collector.slingshot.telemetry-metrics-exclude",
		"Regexp of telemetry metric names to exclude from /sys/class/cxi/<dev>/device/telemetry.",
	).Default("").String()
	slingshotExportTimestamps = kingpin.Flag(
		"collector.slingshot.telemetry-export-timestamps",
		"Export telemetry timestamps as *_timestamp_seconds metrics when telemetry values are formatted as value@timestamp.",
	).Default("false").Bool()
	pcieSpeedValueExpr = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)`)
	slingshotTelemetryMetricDescriptionsByName = map[string]string{
		"atu_cache_evictions":                 "Number of times a translation tag was evicted from NIC translation cache.",
		"atu_cache_miss_ee":                    "Number of translation cache misses attributed to Event Engine clients.",
		"atu_cache_miss_ixe":                   "Number of translation cache misses attributed to Input Transfer Engine clients.",
		"atu_cache_miss_oxe":                   "Number of translation cache misses attributed to Output Transfer Engine clients.",
		"atu_client_req_ee":                    "Number of translation requests issued by Event Engine clients.",
		"atu_client_req_ixe":                   "Number of translation requests issued by Input Transfer Engine clients.",
		"atu_client_req_oxe":                   "Number of translation requests issued by Output Transfer Engine clients.",
		"hni_llr_rx_replay_event":              "Number of receive-side LLR replay events.",
		"hni_llr_tx_replay_event":              "Number of transmit-side LLR replay events.",
		"hni_pcs_good_cw":                      "Number of physical layer codewords received with no errors.",
		"hni_pcs_uncorrected_cw":               "Number of uncorrected physical layer codewords.",
		"hni_rx_ok_36_to_63":                   "Number of receive packets with size 36-63 bytes.",
		"hni_rx_ok_65_to_127":                  "Number of receive packets with size 65-127 bytes.",
		"hni_rx_ok_128_to_255":                 "Number of receive packets with size 128-255 bytes.",
		"hni_rx_ok_256_to_511":                 "Number of receive packets with size 256-511 bytes.",
		"hni_rx_ok_512_to_1023":                "Number of receive packets with size 512-1023 bytes.",
		"hni_rx_ok_1024_to_2047":               "Number of receive packets with size 1024-2047 bytes.",
		"hni_rx_ok_2048_to_4095":               "Number of receive packets with size 2048-4095 bytes.",
		"hni_rx_ok_4096_to_8191":               "Number of receive packets with size 4096-8191 bytes.",
		"hni_rx_ok_8192_to_max":                "Number of receive packets with size 8192 bytes and larger.",
		"hni_tx_ok_36_to_63":                   "Number of transmit packets with size 36-63 bytes.",
		"hni_tx_ok_65_to_127":                  "Number of transmit packets with size 65-127 bytes.",
		"hni_tx_ok_128_to_255":                 "Number of transmit packets with size 128-255 bytes.",
		"hni_tx_ok_256_to_511":                 "Number of transmit packets with size 256-511 bytes.",
		"hni_tx_ok_512_to_1023":                "Number of transmit packets with size 512-1023 bytes.",
		"hni_tx_ok_1024_to_2047":               "Number of transmit packets with size 1024-2047 bytes.",
		"hni_tx_ok_2048_to_4095":               "Number of transmit packets with size 2048-4095 bytes.",
		"hni_tx_ok_4096_to_8191":               "Number of transmit packets with size 4096-8191 bytes.",
		"hni_tx_ok_8192_to_max":                "Number of transmit packets with size 8192 bytes and larger.",
		"parbs_tarb_pi_non_posted_blocked_cnt": "Number of cycles non-posted PCIe path was blocked.",
		"parbs_tarb_pi_non_posted_pkts":        "Number of PCIe packets transferred using non-posted path.",
		"parbs_tarb_pi_posted_blocked_cnt":     "Number of cycles posted PCIe path was blocked.",
		"parbs_tarb_pi_posted_pkts":            "Number of PCIe packets transferred using posted path.",
		"pct_bad_seq_nacks":                    "Number of sequence-error NACKs.",
		"pct_conn_sct_open":                    "Number of open source connection tracking entries.",
		"pct_no_matching_tct":                  "Number of NACKs caused by no matching target connection.",
		"pct_no_tct_nacks":                     "Number of NACKs caused by target connection table exhaustion.",
		"pct_no_trs_nacks":                     "Number of NACKs caused by target result store exhaustion.",
		"pct_req_ordered":                      "Number of ordered requests.",
		"pct_req_unordered":                    "Number of unordered requests.",
		"pct_resource_busy":                    "Number of NACKs caused by target resource busy conditions.",
		"pct_responses_received":               "Number of responses received.",
		"pct_retry_srb_requests":               "Number of retry requests issued by source retry buffer.",
		"pct_spt_timeouts":                     "Number of source packet tracking timeouts.",
	}
	slingshotTelemetryMetricDescriptionsByPrefix = []slingshotTelemetryMetricPrefixDescription{
		{prefix: "atu_cache_hit_base_page_size", description: "Number of translation cache hits on base page size."},
		{prefix: "atu_cache_hit_derivative1_page_size", description: "Number of translation cache hits on derivative1 page size."},
		{prefix: "atu_cache_hit_derivative2_page_size", description: "Number of translation cache hits on derivative2 page size."},
		{prefix: "atu_cache_miss", description: "Number of translation cache misses."},
		{prefix: "atu_nta_trans_latency", description: "NTA translation latency histogram bucket count."},
		{prefix: "cq_cycles_blocked", description: "Number of cycles command queue could not make progress due to arbitration/blocking."},
		{prefix: "cq_dma_cmd_counts", description: "Number of DMA commands by opcode bucket."},
		{prefix: "cq_num_cq_cmds", description: "Number of successfully parsed command queue commands."},
		{prefix: "cq_num_ll_ops_rejected", description: "Number of low-latency operations rejected from immediate-send path."},
		{prefix: "cq_num_ll_ops_split", description: "Number of low-latency operations split into multiple writes."},
		{prefix: "cq_num_ll_ops_successful", description: "Number of low-latency operations accepted successfully."},
		{prefix: "ee_deferred_eq_switch_cntr", description: "Number of deferred event queue buffer switches due to insufficient space."},
		{prefix: "ee_events_dropped_ordinary_cntr", description: "Number of ordinary full events dropped because event queue was full."},
		{prefix: "ee_events_enqueued_cntr", description: "Number of full events enqueued to event queues."},
		{prefix: "hni_pkts_recv_by_tc", description: "Number of packets received by traffic class."},
		{prefix: "hni_pkts_sent_by_tc", description: "Number of packets sent by traffic class."},
		{prefix: "hni_rx_ok", description: "Number of packets received by packet-size bucket."},
		{prefix: "hni_rx_paused", description: "Number of cycles receive path was paused for a traffic class."},
		{prefix: "hni_tx_ok", description: "Number of packets sent by packet-size bucket."},
		{prefix: "hni_tx_paused", description: "Number of cycles transmit path was paused for a traffic class."},
		{prefix: "lpe_net_match_overflow", description: "Number of messages delivered to overflow list because no priority match existed."},
		{prefix: "lpe_net_match_priority", description: "Number of messages matched on priority list (receive posted before arrival)."},
		{prefix: "lpe_net_match_request", description: "Number of requests matched on request/software list."},
		{prefix: "lpe_net_match_requests", description: "Number of total network match requests observed."},
		{prefix: "lpe_net_match_success", description: "Number of network requests successfully completed by list-processing engines."},
		{prefix: "lpe_rndzv_puts", description: "Number of rendezvous put operations received."},
		{prefix: "lpe_rndzv_puts_offloaded", description: "Number of rendezvous put operations successfully offloaded."},
		{prefix: "lpe_sts_app_match_attempts", description: "Number of application-side match attempts."},
		{prefix: "lpe_sts_net_match_attempts", description: "Number of network-side match attempts."},
		{prefix: "oxe_ptl_tx_put_pkts_tsc", description: "Number of portal put payload packets by traffic shaping class."},
		{prefix: "pct_host_access_latency", description: "Host access latency histogram bucket count."},
		{prefix: "pct_req_rsp_latency", description: "Request/response latency histogram bucket count."},
		{prefix: "tou_num_trig_cmds", description: "Number of triggered commands."},
	}
)

type slingshotTelemetryMetricPrefixDescription struct {
	prefix      string
	description string
}

type slingshotTelemetryMetricSpec struct {
	metricName  string
	description string
	index       string
}

type slingshotMetricFilter struct {
	include *regexp.Regexp
	exclude *regexp.Regexp
}

type slingshotCollector struct {
	logger                          *slog.Logger
	enableInfo                      bool
	enableTelemetry                 bool
	exportTelemetryTS               bool
	telemetryMetricFilter           slingshotMetricFilter
	telemetryMetricDescriptions     map[string]string
	infoDesc                        *prometheus.Desc
	nidDesc                         *prometheus.Desc
	pidGranuleDesc                  *prometheus.Desc
	linkMTUDesc                     *prometheus.Desc
	pcieInfoDesc                    *prometheus.Desc
	pcieSpeedGTSDesc                *prometheus.Desc
	pcieWidthDesc                   *prometheus.Desc
	linkInfoDesc                    *prometheus.Desc
	linkSpeedDesc                   *prometheus.Desc
	scrapedMetricsDesc              *prometheus.Desc
	scrapeErrorsDesc                *prometheus.Desc
	telemetryMetricDescByName       map[string]*prometheus.Desc
	telemetryTimestampDescByName    map[string]*prometheus.Desc
	uniqueTelemetryNameByRaw        map[string]string
	rawTelemetryByUniqueMetricName  map[string]string
	mu                              sync.RWMutex
}

const (
	slingshotSourceInfo      = "info"
	slingshotSourceTelemetry = "telemetry"
)

func init() {
	registerCollector("slingshot", defaultDisabled, NewSlingshotCollector)
	registerCollector("slingshot_info", defaultDisabled, NewSlingshotInfoCollector)
	registerCollector("slingshot_metrics", defaultDisabled, NewSlingshotMetricsCollector)
}

func NewSlingshotCollector(logger *slog.Logger) (Collector, error) {
	return newSlingshotCollector(logger, true, true)
}

func NewSlingshotInfoCollector(logger *slog.Logger) (Collector, error) {
	return newSlingshotCollector(logger, true, false)
}

func NewSlingshotMetricsCollector(logger *slog.Logger) (Collector, error) {
	return newSlingshotCollector(logger, false, true)
}

func newSlingshotCollector(logger *slog.Logger, enableInfo, enableTelemetry bool) (Collector, error) {
	telemetryFilter := slingshotMetricFilter{}
	if enableTelemetry {
		var err error
		telemetryFilter, err = newSlingshotMetricFilter(*slingshotTelemetryMetricsInclude, *slingshotTelemetryMetricsExclude)
		if err != nil {
			return nil, fmt.Errorf("invalid slingshot telemetry metric filter: %w", err)
		}
	}

	return &slingshotCollector{
		logger:                logger,
		enableInfo:            enableInfo,
		enableTelemetry:       enableTelemetry,
		exportTelemetryTS:     *slingshotExportTimestamps,
		telemetryMetricFilter: telemetryFilter,
		telemetryMetricDescriptions: slingshotTelemetryMetricDescriptionsByName,
		infoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "info"),
			"Non-numeric Slingshot NIC metadata. Value is always 1.",
			[]string{
				"device",
				"interface",
				"fru_description",
				"part_number",
				"serial_number",
				"firmware_version",
				"mac",
			},
			nil,
		),
		nidDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "nid"),
			"Slingshot network ID (NID).",
			[]string{"device", "interface"},
			nil,
		),
		pidGranuleDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "pid_granule"),
			"Slingshot PID granule value.",
			[]string{"device", "interface"},
			nil,
		),
		linkMTUDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "link_mtu"),
			"Slingshot link MTU.",
			[]string{"device", "interface"},
			nil,
		),
		pcieInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "pcie_info"),
			"Slingshot PCIe metadata. Value is always 1.",
			[]string{"device", "interface", "slot"},
			nil,
		),
		pcieSpeedGTSDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "pcie_speed_gts"),
			"Slingshot PCIe link speed in GT/s.",
			[]string{"device", "interface"},
			nil,
		),
		pcieWidthDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "pcie_width"),
			"Slingshot PCIe link width in lanes.",
			[]string{"device", "interface"},
			nil,
		),
		linkInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "link_info"),
			"Slingshot link metadata. Value is always 1.",
			[]string{"device", "interface", "state", "link_layer_retry", "loopback", "media"},
			nil,
		),
		linkSpeedDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "link_speed"),
			"Slingshot link speed in bits per second.",
			[]string{"device", "interface"},
			nil,
		),
		scrapedMetricsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "scraped_metrics"),
			"Number of slingshot metrics exported in the current scrape by source.",
			[]string{"source", "interface"},
			nil,
		),
		scrapeErrorsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "slingshot", "scrape_errors"),
			"Number of slingshot metric read/parse errors in the current scrape by source.",
			[]string{"source", "interface"},
			nil,
		),
		telemetryMetricDescByName:      make(map[string]*prometheus.Desc),
		telemetryTimestampDescByName:   make(map[string]*prometheus.Desc),
		uniqueTelemetryNameByRaw:       make(map[string]string),
		rawTelemetryByUniqueMetricName: make(map[string]string),
	}, nil
}

func (c *slingshotCollector) Update(ch chan<- prometheus.Metric) error {
	if !c.enableInfo && !c.enableTelemetry {
		return ErrNoData
	}

	classPath := sysFilePath(filepath.Join("class", "cxi"))
	entries, err := os.ReadDir(classPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("slingshot sysfs path not found, skipping", "path", classPath)
			return ErrNoData
		}
		return fmt.Errorf("read slingshot class directory %q: %w", classPath, err)
	}

	devices := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "cxi") {
			continue
		}

		entryPath := filepath.Join(classPath, name)
		info, err := os.Stat(entryPath)
		if err != nil || !info.IsDir() {
			continue
		}
		devices = append(devices, name)
	}
	if len(devices) == 0 {
		return ErrNoData
	}
	sort.Strings(devices)

	scraped := map[string]float64{
		slingshotSourceInfo:      0,
		slingshotSourceTelemetry: 0,
	}
	errorsBySource := map[string]float64{
		slingshotSourceInfo:      0,
		slingshotSourceTelemetry: 0,
	}

	sources := make([]string, 0, 2)
	if c.enableInfo {
		sources = append(sources, slingshotSourceInfo)
	}
	if c.enableTelemetry {
		sources = append(sources, slingshotSourceTelemetry)
	}

	for _, device := range devices {
		iface := lookupInterfaceName(device)

		if c.enableInfo {
			if err := c.collectInfoMetrics(ch, device, iface); err != nil {
				errorsBySource[slingshotSourceInfo]++
				c.logger.Debug("failed to collect slingshot info", "device", device, "err", err)
			} else {
				scraped[slingshotSourceInfo]++
			}
		}

		if c.enableTelemetry {
			telemetryRoot := sysFilePath(filepath.Join("class", "cxi", device, "device", "telemetry"))
			telemetryScraped, telemetryErrs := c.collectTelemetryMetrics(ch, device, iface, telemetryRoot)
			scraped[slingshotSourceTelemetry] += float64(telemetryScraped)
			errorsBySource[slingshotSourceTelemetry] += float64(telemetryErrs)
		}
	}

	for _, source := range sources {
		ch <- prometheus.MustNewConstMetric(c.scrapedMetricsDesc, prometheus.GaugeValue, scraped[source], source, "all")
		ch <- prometheus.MustNewConstMetric(c.scrapeErrorsDesc, prometheus.GaugeValue, errorsBySource[source], source, "all")
	}

	totalScraped := 0.0
	for _, source := range sources {
		totalScraped += scraped[source]
	}
	if totalScraped == 0 {
		return ErrNoData
	}
	return nil
}

func (c *slingshotCollector) collectInfoMetrics(ch chan<- prometheus.Metric, device, iface string) error {
	devicePath := sysFilePath(filepath.Join("class", "cxi", device, "device"))

	if _, err := os.Stat(devicePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoData
		}
		return fmt.Errorf("stat device path %q: %w", devicePath, err)
	}

	mac := ""
	if iface != "" {
		mac = readTrimmedFileOrEmpty(sysFilePath(filepath.Join("class", "net", iface, "address")))
	}

	fruDescription := readTrimmedFileOrEmpty(filepath.Join(devicePath, "fru", "description"))
	partNumber := readTrimmedFileOrEmpty(filepath.Join(devicePath, "fru", "part_number"))
	serialNumber := readTrimmedFileOrEmpty(filepath.Join(devicePath, "fru", "serial_number"))
	firmwareVersion := readTrimmedFileOrEmpty(filepath.Join(devicePath, "uc", "qspi_blob_version"))
	pcieSlot := readPCIESlot(devicePath)
	pcieSpeedGTS, hasPCIeSpeedGTS := readPCIESpeedGTS(devicePath)
	pcieWidth, hasPCIeWidth := readPCIELinkWidth(devicePath)
	linkLayerRetry := readTrimmedFileOrEmpty(filepath.Join(devicePath, "port", "link_layer_retry"))
	linkLoopback := readTrimmedFileOrEmpty(filepath.Join(devicePath, "port", "loopback"))
	linkMedia := readTrimmedFileOrEmpty(filepath.Join(devicePath, "port", "media"))
	linkSpeed := readTrimmedFileOrEmpty(filepath.Join(devicePath, "port", "speed"))
	linkState := readTrimmedFileOrEmpty(filepath.Join(devicePath, "port", "link"))

	ch <- prometheus.MustNewConstMetric(
		c.infoDesc,
		prometheus.GaugeValue,
		1,
		device,
		iface,
		fruDescription,
		partNumber,
		serialNumber,
		firmwareVersion,
		mac,
	)

	if pcieSlot != "" {
		ch <- prometheus.MustNewConstMetric(c.pcieInfoDesc, prometheus.GaugeValue, 1, device, iface, pcieSlot)
	}
	if hasPCIeSpeedGTS {
		ch <- prometheus.MustNewConstMetric(c.pcieSpeedGTSDesc, prometheus.GaugeValue, pcieSpeedGTS, device, iface)
	}
	if hasPCIeWidth {
		ch <- prometheus.MustNewConstMetric(c.pcieWidthDesc, prometheus.GaugeValue, pcieWidth, device, iface)
	}
	if linkState != "" || linkLayerRetry != "" || linkLoopback != "" || linkMedia != "" {
		ch <- prometheus.MustNewConstMetric(c.linkInfoDesc, prometheus.GaugeValue, 1, device, iface, linkState, linkLayerRetry, linkLoopback, linkMedia)
	}
	if linkSpeedBps, ok := parseLinkSpeedBps(linkSpeed); ok {
		ch <- prometheus.MustNewConstMetric(c.linkSpeedDesc, prometheus.GaugeValue, linkSpeedBps, device, iface)
	}

	if nid, ok := readFloatFile(filepath.Join(devicePath, "properties", "nid")); ok {
		ch <- prometheus.MustNewConstMetric(c.nidDesc, prometheus.GaugeValue, nid, device, iface)
	}
	if pidGranule, ok := readFloatFile(filepath.Join(devicePath, "properties", "pid_granule")); ok {
		ch <- prometheus.MustNewConstMetric(c.pidGranuleDesc, prometheus.GaugeValue, pidGranule, device, iface)
	}
	if linkMTU, ok := readFloatFile(filepath.Join(devicePath, "port", "mtu")); ok {
		ch <- prometheus.MustNewConstMetric(c.linkMTUDesc, prometheus.GaugeValue, linkMTU, device, iface)
	}

	return nil
}

func (c *slingshotCollector) collectTelemetryMetrics(ch chan<- prometheus.Metric, device, iface, root string) (int, int) {
	rootInfo, err := os.Stat(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, 0
		}
		c.logger.Debug("failed to stat slingshot telemetry root", "path", root, "err", err)
		return 0, 1
	}
	if !rootInfo.IsDir() {
		c.logger.Debug("slingshot telemetry root is not a directory", "path", root)
		return 0, 1
	}

	scraped := 0
	parseErrors := 0

	rootResolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		rootResolved = root
	}

	visited := map[string]struct{}{filepath.Clean(rootResolved): {}}
	var walk func(absPath, relPrefix string)
	walk = func(absPath, relPrefix string) {
		entries, err := os.ReadDir(absPath)
		if err != nil {
			parseErrors++
			return
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name() < entries[j].Name()
		})

		for _, entry := range entries {
			entryName := entry.Name()
			entryAbsPath := filepath.Join(absPath, entryName)
			relPath := entryName
			if relPrefix != "" {
				relPath = filepath.Join(relPrefix, entryName)
			}

			entryMode := entry.Type()
			if entryMode&os.ModeSymlink != 0 {
				resolvedPath, err := filepath.EvalSymlinks(entryAbsPath)
				if err != nil {
					parseErrors++
					continue
				}
				resolvedInfo, err := os.Stat(resolvedPath)
				if err != nil {
					parseErrors++
					continue
				}

				if resolvedInfo.IsDir() {
					if !isPathWithinRoot(rootResolved, resolvedPath) {
						continue
					}
					canonical := filepath.Clean(resolvedPath)
					if _, seen := visited[canonical]; seen {
						continue
					}
					visited[canonical] = struct{}{}
					walk(resolvedPath, relPath)
					continue
				}

				if !resolvedInfo.Mode().IsRegular() {
					continue
				}

				emitted, errs := c.collectTelemetryMetricFile(ch, device, iface, relPath, entryAbsPath)
				scraped += emitted
				parseErrors += errs
				continue
			}

			if entry.IsDir() {
				walk(entryAbsPath, relPath)
				continue
			}

			entryType := entry.Type()
			// In virtual filesystems like sysfs, d_type may be unknown (mode 0)
			// for regular files. Only skip entries when the type is known and
			// definitely not a regular file.
			if entryType&os.ModeType != 0 && !entryType.IsRegular() {
				continue
			}

			emitted, errs := c.collectTelemetryMetricFile(ch, device, iface, relPath, entryAbsPath)
			scraped += emitted
			parseErrors += errs
		}
	}

	walk(rootResolved, "")
	return scraped, parseErrors
}

func (c *slingshotCollector) collectTelemetryMetricFile(
	ch chan<- prometheus.Metric,
	device, iface, relPath, filePath string,
) (int, int) {
	rawMetricName := strings.ReplaceAll(relPath, string(filepath.Separator), "_")
	spec, ok := c.telemetryMetricSpec(rawMetricName)
	if !ok {
		return 0, 0
	}
	if !c.telemetryMetricFilter.Allow(rawMetricName) {
		return 0, 0
	}

	rawValue, err := readTrimmedFile(filePath)
	if err != nil {
		return 0, 1
	}

	value, ts, hasTimestamp, err := parseMetricValue(rawValue)
	if err != nil {
		if errors.Is(err, errNonNumericMetricValue) {
			return 0, 0
		}
		return 0, 1
	}

	desc := c.telemetryMetricDesc(rawMetricName, spec)
	if spec.index == "" {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, device, iface)
	} else {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, device, iface, spec.index)
	}
	scraped := 1

	if c.exportTelemetryTS && hasTimestamp {
		tsDesc := c.telemetryTimestampDesc(rawMetricName, spec)
		if spec.index == "" {
			ch <- prometheus.MustNewConstMetric(tsDesc, prometheus.GaugeValue, normalizeEpochToSeconds(ts), device, iface)
		} else {
			ch <- prometheus.MustNewConstMetric(tsDesc, prometheus.GaugeValue, normalizeEpochToSeconds(ts), device, iface, spec.index)
		}
		scraped++
	}

	return scraped, 0
}

func (c *slingshotCollector) telemetryMetricDesc(rawMetricName string, spec slingshotTelemetryMetricSpec) *prometheus.Desc {
	c.mu.RLock()
	desc, ok := c.telemetryMetricDescByName[rawMetricName]
	c.mu.RUnlock()
	if ok {
		return desc
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if desc, ok := c.telemetryMetricDescByName[rawMetricName]; ok {
		return desc
	}

	uniqueName := c.uniqueTelemetryMetricNameLocked(spec.metricName)
	labelNames := []string{"device", "interface"}
	if spec.index != "" {
		labelNames = append(labelNames, "index")
	}
	desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "slingshot", "telemetry_"+uniqueName),
		spec.description,
		labelNames,
		nil,
	)
	c.telemetryMetricDescByName[rawMetricName] = desc
	return desc
}

func (c *slingshotCollector) telemetryTimestampDesc(rawMetricName string, spec slingshotTelemetryMetricSpec) *prometheus.Desc {
	c.mu.RLock()
	desc, ok := c.telemetryTimestampDescByName[rawMetricName]
	c.mu.RUnlock()
	if ok {
		return desc
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if desc, ok := c.telemetryTimestampDescByName[rawMetricName]; ok {
		return desc
	}

	uniqueName := c.uniqueTelemetryMetricNameLocked(spec.metricName)
	labelNames := []string{"device", "interface"}
	if spec.index != "" {
		labelNames = append(labelNames, "index")
	}
	desc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "slingshot", "telemetry_"+uniqueName+"_timestamp_seconds"),
		fmt.Sprintf("Timestamp for slingshot telemetry metric %q in seconds since Unix epoch.", spec.metricName),
		labelNames,
		nil,
	)
	c.telemetryTimestampDescByName[rawMetricName] = desc
	return desc
}

func (c *slingshotCollector) uniqueTelemetryMetricNameLocked(rawMetricName string) string {
	if name, ok := c.uniqueTelemetryNameByRaw[rawMetricName]; ok {
		return name
	}

	base := sanitizeMetricName(rawMetricName)
	candidate := base
	seq := 0

	for {
		if existingRaw, ok := c.rawTelemetryByUniqueMetricName[candidate]; !ok || existingRaw == rawMetricName {
			c.rawTelemetryByUniqueMetricName[candidate] = rawMetricName
			c.uniqueTelemetryNameByRaw[rawMetricName] = candidate
			return candidate
		}

		seq++
		candidate = fmt.Sprintf("%s_%d", base, seq)
	}
}

func isPathWithinRoot(root, candidate string) bool {
	root = filepath.Clean(root)
	candidate = filepath.Clean(candidate)

	rel, err := filepath.Rel(root, candidate)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func newSlingshotMetricFilter(includeExpr, excludeExpr string) (slingshotMetricFilter, error) {
	include, err := regexp.Compile(includeExpr)
	if err != nil {
		return slingshotMetricFilter{}, fmt.Errorf("compile include regexp %q: %w", includeExpr, err)
	}

	excludePattern := excludeExpr
	if excludePattern == "" {
		excludePattern = "$^"
	}
	exclude, err := regexp.Compile(excludePattern)
	if err != nil {
		return slingshotMetricFilter{}, fmt.Errorf("compile exclude regexp %q: %w", excludeExpr, err)
	}

	return slingshotMetricFilter{include: include, exclude: exclude}, nil
}

func (c *slingshotCollector) telemetryMetricSpec(rawMetricName string) (slingshotTelemetryMetricSpec, bool) {
	if description, ok := c.telemetryMetricDescriptions[rawMetricName]; ok {
		return slingshotTelemetryMetricSpec{
			metricName:  rawMetricName,
			description: description,
		}, true
	}

	for _, prefixDescription := range slingshotTelemetryMetricDescriptionsByPrefix {
		prefix := prefixDescription.prefix + "_"
		if !strings.HasPrefix(rawMetricName, prefix) {
			continue
		}

		index := strings.TrimPrefix(rawMetricName, prefix)
		if !isNumericMetricSuffix(index) {
			continue
		}

		return slingshotTelemetryMetricSpec{
			metricName:  prefixDescription.prefix,
			description: prefixDescription.description,
			index:       index,
		}, true
	}

	return slingshotTelemetryMetricSpec{}, false
}

func isNumericMetricSuffix(raw string) bool {
	if raw == "" {
		return false
	}
	_, err := strconv.ParseUint(raw, 10, 64)
	return err == nil
}

func (f slingshotMetricFilter) Allow(metricName string) bool {
	if f.include != nil && !f.include.MatchString(metricName) {
		return false
	}
	if f.exclude != nil && f.exclude.MatchString(metricName) {
		return false
	}
	return true
}

func readTrimmedFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func readTrimmedFileOrEmpty(path string) string {
	s, err := readTrimmedFile(path)
	if err != nil {
		return ""
	}
	return s
}

func readFloatFile(path string) (float64, bool) {
	raw, err := readTrimmedFile(path)
	if err != nil || raw == "" {
		return 0, false
	}

	v, err := parseNumericToken(raw)
	if err != nil {
		return 0, false
	}
	return v, true
}

var errNonNumericMetricValue = errors.New("non-numeric slingshot metric value")

func parseMetricValue(raw string) (value float64, timestamp float64, hasTimestamp bool, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, 0, false, errNonNumericMetricValue
	}

	parts := strings.SplitN(raw, "@", 2)
	value, err = parseNumericToken(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, false, fmt.Errorf("parse metric value %q: %w", raw, errNonNumericMetricValue)
	}
	if len(parts) == 1 {
		return value, 0, false, nil
	}

	timestamp, err = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, false, fmt.Errorf("parse metric timestamp from %q: %w", raw, err)
	}
	return value, timestamp, true, nil
}

func parseNumericToken(raw string) (float64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, errNonNumericMetricValue
	}

	if v, err := strconv.ParseFloat(raw, 64); err == nil {
		return v, nil
	}
	if v, err := strconv.ParseInt(raw, 0, 64); err == nil {
		return float64(v), nil
	}
	if v, err := strconv.ParseUint(raw, 0, 64); err == nil {
		return float64(v), nil
	}

	return 0, errNonNumericMetricValue
}

func normalizeEpochToSeconds(v float64) float64 {
	// Handle either epoch seconds (possibly fractional) or epoch nanoseconds.
	if v >= 1e12 {
		return v / 1e9
	}
	return v
}

func parseLinkSpeedBps(raw string) (float64, bool) {
	raw = strings.TrimSpace(strings.ToUpper(raw))
	if raw == "" {
		return 0, false
	}

	// Normalize common suffixes: "100G", "100 GB/S", "100GBPS", "100000000000".
	raw = strings.ReplaceAll(raw, " ", "")
	raw = strings.TrimSuffix(raw, "B/S")
	raw = strings.TrimSuffix(raw, "BPS")
	raw = strings.TrimSuffix(raw, "BIT/S")
	raw = strings.TrimSuffix(raw, "BITS/S")
	raw = strings.TrimSuffix(raw, "PS")

	multiplier := 1.0
	if len(raw) > 0 {
		switch raw[len(raw)-1] {
		case 'T':
			multiplier = 1e12
			raw = raw[:len(raw)-1]
		case 'G':
			multiplier = 1e9
			raw = raw[:len(raw)-1]
		case 'M':
			multiplier = 1e6
			raw = raw[:len(raw)-1]
		case 'K':
			multiplier = 1e3
			raw = raw[:len(raw)-1]
		}
	}

	v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return 0, false
	}

	return v * multiplier, true
}

func sanitizeMetricName(metric string) string {
	metric = strings.ToLower(metric)
	var b strings.Builder
	b.Grow(len(metric))

	prevUnderscore := false
	for i := 0; i < len(metric); i++ {
		ch := metric[i]
		isAlpha := ch >= 'a' && ch <= 'z'
		isDigit := ch >= '0' && ch <= '9'
		if isAlpha || isDigit {
			b.WriteByte(ch)
			prevUnderscore = false
			continue
		}
		if !prevUnderscore {
			b.WriteByte('_')
			prevUnderscore = true
		}
	}

	out := strings.Trim(b.String(), "_")
	if out == "" {
		out = "metric"
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "metric_" + out
	}
	return out
}

func lookupInterfaceName(device string) string {
	pattern := sysFilePath(filepath.Join("class", "net", "*", "device", "cxi", device))
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}
	sort.Strings(matches)

	match := matches[0]
	netDevicePath := filepath.Dir(filepath.Dir(filepath.Dir(match)))
	return filepath.Base(netDevicePath)
}

func readPCIESlot(devicePath string) string {
	target, err := filepath.EvalSymlinks(devicePath)
	if err != nil {
		return ""
	}

	for dir := target; dir != "." && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		slot := filepath.Base(dir)
		if strings.Contains(slot, ":") {
			return slot
		}
		if next := filepath.Dir(dir); next == dir {
			break
		}
	}
	return ""
}

func readPCIESpeed(devicePath string) string {
	speed := readTrimmedFileOrEmpty(filepath.Join(devicePath, "properties", "current_esm_link_speed"))
	if speed == "" || strings.EqualFold(speed, "Absent") || strings.EqualFold(speed, "Disabled") {
		speed = readTrimmedFileOrEmpty(filepath.Join(devicePath, "current_link_speed"))
	}
	if speed == "" {
		return ""
	}
	return speed
}

func readPCIESpeedGTS(devicePath string) (float64, bool) {
	speed := readPCIESpeed(devicePath)
	if speed == "" {
		return 0, false
	}

	match := pcieSpeedValueExpr.FindStringSubmatch(speed)
	if len(match) < 2 {
		return 0, false
	}

	v, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, false
	}

	return v, true
}

func readPCIELinkWidth(devicePath string) (float64, bool) {
	width := readTrimmedFileOrEmpty(filepath.Join(devicePath, "current_link_width"))
	if width == "" {
		return 0, false
	}

	cleanWidth := strings.TrimSpace(strings.TrimPrefix(width, "x"))
	cleanWidth = strings.TrimSpace(strings.TrimPrefix(cleanWidth, "X"))
	if cleanWidth == "" {
		return 0, false
	}

	v, err := strconv.ParseFloat(cleanWidth, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}