// Copyright 2021 The Prometheus Authors
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

//go:build !noethtool
// +build !noethtool

// The hard work of collecting data from the kernel via the ethtool interfaces is done by
// https://github.com/safchain/ethtool/
// by Sylvain Afchain. Used under the Apache license.

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
	"github.com/safchain/ethtool"
	"golang.org/x/sys/unix"
)

var (
	ethtoolDeviceInclude   = kingpin.Flag("collector.ethtool.device-include", "Regexp of ethtool devices to include (mutually exclusive to device-exclude).").String()
	ethtoolDeviceExclude   = kingpin.Flag("collector.ethtool.device-exclude", "Regexp of ethtool devices to exclude (mutually exclusive to device-include).").String()
	ethtoolIncludedMetrics = kingpin.Flag("collector.ethtool.metrics-include", "Regexp of ethtool stats to include.").Default(".*").String()
	ethtoolReceivedRegex   = regexp.MustCompile(`(^|_)rx(_|$)`)
	ethtoolTransmitRegex   = regexp.MustCompile(`(^|_)tx(_|$)`)
)

type Ethtool interface {
	DriverInfo(string) (ethtool.DrvInfo, error)
	Stats(string) (map[string]uint64, error)
	LinkInfo(string) (ethtool.EthtoolCmd, error)
}

type ethtoolLibrary struct {
	ethtool *ethtool.Ethtool
}

func (e *ethtoolLibrary) DriverInfo(intf string) (ethtool.DrvInfo, error) {
	return e.ethtool.DriverInfo(intf)
}

func (e *ethtoolLibrary) Stats(intf string) (map[string]uint64, error) {
	return e.ethtool.Stats(intf)
}

func (e *ethtoolLibrary) LinkInfo(intf string) (ethtool.EthtoolCmd, error) {
	var ethtoolCmd ethtool.EthtoolCmd
	_, err := ethtoolCmd.CmdGet(intf)
	return ethtoolCmd, err
}

type ethtoolCollector struct {
	fs             sysfs.FS
	entries        map[string]*prometheus.Desc
	entriesMutex   sync.Mutex
	ethtool        Ethtool
	deviceFilter   deviceFilter
	infoDesc       *prometheus.Desc
	metricsPattern *regexp.Regexp
	logger         *slog.Logger
}

// makeEthtoolCollector is the internal constructor for EthtoolCollector.
// This allows NewEthtoolTestCollector to override its .ethtool interface
// for testing.
func makeEthtoolCollector(logger *slog.Logger) (*ethtoolCollector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ethtool library: %w", err)
	}

	if *ethtoolDeviceInclude != "" {
		logger.Info("Parsed flag --collector.ethtool.device-include", "flag", *ethtoolDeviceInclude)
	}
	if *ethtoolDeviceExclude != "" {
		logger.Info("Parsed flag --collector.ethtool.device-exclude", "flag", *ethtoolDeviceExclude)
	}
	if *ethtoolIncludedMetrics != "" {
		logger.Info("Parsed flag --collector.ethtool.metrics-include", "flag", *ethtoolIncludedMetrics)
	}

	// Pre-populate some common ethtool metrics.
	return &ethtoolCollector{
		fs:             fs,
		ethtool:        &ethtoolLibrary{e},
		deviceFilter:   newDeviceFilter(*ethtoolDeviceExclude, *ethtoolDeviceInclude),
		metricsPattern: regexp.MustCompile(*ethtoolIncludedMetrics),
		logger:         logger,
		entries: map[string]*prometheus.Desc{
			"rx_bytes": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "ethtool", "received_bytes_total"),
				"Network interface bytes received",
				[]string{"device"}, nil,
			),
			"rx_dropped": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "ethtool", "received_dropped_total"),
				"Number of received frames dropped",
				[]string{"device"}, nil,
			),
			"rx_errors": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "ethtool", "received_errors_total"),
				"Number of received frames with errors",
				[]string{"device"}, nil,
			),
			"rx_packets": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "ethtool", "received_packets_total"),
				"Network interface packets received",
				[]string{"device"}, nil,
			),
			"tx_bytes": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "ethtool", "transmitted_bytes_total"),
				"Network interface bytes sent",
				[]string{"device"}, nil,
			),
			"tx_errors": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "ethtool", "transmitted_errors_total"),
				"Number of sent frames with errors",
				[]string{"device"}, nil,
			),
			"tx_packets": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "ethtool", "transmitted_packets_total"),
				"Network interface packets sent",
				[]string{"device"}, nil,
			),

			// link info
			"supported_port": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "supported_port_info"),
				"Type of ports or PHYs supported by network device",
				[]string{"device", "type"}, nil,
			),
			"supported_speed": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "supported_speed_bytes"),
				"Combination of speeds and features supported by network device",
				[]string{"device", "duplex", "mode"}, nil,
			),
			"supported_autonegotiate": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "autonegotiate_supported"),
				"If this port device supports autonegotiate",
				[]string{"device"}, nil,
			),
			"supported_pause": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "pause_supported"),
				"If this port device supports pause frames",
				[]string{"device"}, nil,
			),
			"supported_asymmetricpause": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "asymmetricpause_supported"),
				"If this port device supports asymmetric pause frames",
				[]string{"device"}, nil,
			),
			"advertised_speed": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "advertised_speed_bytes"),
				"Combination of speeds and features offered by network device",
				[]string{"device", "duplex", "mode"}, nil,
			),
			"advertised_autonegotiate": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "autonegotiate_advertised"),
				"If this port device offers autonegotiate",
				[]string{"device"}, nil,
			),
			"advertised_pause": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "pause_advertised"),
				"If this port device offers pause capability",
				[]string{"device"}, nil,
			),
			"advertised_asymmetricpause": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "asymmetricpause_advertised"),
				"If this port device offers asymmetric pause capability",
				[]string{"device"}, nil,
			),
			"autonegotiate": prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "network", "autonegotiate"),
				"If this port is using autonegotiate",
				[]string{"device"}, nil,
			),
		},
		infoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "ethtool", "info"),
			"A metric with a constant '1' value labeled by bus_info, device, driver, expansion_rom_version, firmware_version, version.",
			[]string{"bus_info", "device", "driver", "expansion_rom_version", "firmware_version", "version"}, nil,
		),
	}, nil
}

func init() {
	registerCollector("ethtool", defaultDisabled, NewEthtoolCollector)
}

// Generate the fully-qualified metric name for the ethool metric.
func buildEthtoolFQName(metric string) string {
	metricName := strings.TrimLeft(strings.ToLower(SanitizeMetricName(metric)), "_")
	metricName = ethtoolReceivedRegex.ReplaceAllString(metricName, "${1}received${2}")
	metricName = ethtoolTransmitRegex.ReplaceAllString(metricName, "${1}transmitted${2}")
	return prometheus.BuildFQName(namespace, "ethtool", metricName)
}

// NewEthtoolCollector returns a new Collector exposing ethtool stats.
func NewEthtoolCollector(logger *slog.Logger) (Collector, error) {
	return makeEthtoolCollector(logger)
}

// updatePortCapabilities generates metrics for autonegotiate, pause and asymmetricpause.
// The bit offsets here correspond to ethtool_link_mode_bit_indices in linux/include/uapi/linux/ethtool.h
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/uapi/linux/ethtool.h
func (c *ethtoolCollector) updatePortCapabilities(ch chan<- prometheus.Metric, prefix string, device string, linkModes uint32) {
	var (
		autonegotiate   = 0.0
		pause           = 0.0
		asymmetricPause = 0.0
	)
	if linkModes&(1<<unix.ETHTOOL_LINK_MODE_Autoneg_BIT) != 0 {
		autonegotiate = 1.0
	}
	if linkModes&(1<<unix.ETHTOOL_LINK_MODE_Pause_BIT) != 0 {
		pause = 1.0
	}
	if linkModes&(1<<unix.ETHTOOL_LINK_MODE_Asym_Pause_BIT) != 0 {
		asymmetricPause = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.entry(fmt.Sprintf("%s_autonegotiate", prefix)), prometheus.GaugeValue, autonegotiate, device)
	ch <- prometheus.MustNewConstMetric(c.entry(fmt.Sprintf("%s_pause", prefix)), prometheus.GaugeValue, pause, device)
	ch <- prometheus.MustNewConstMetric(c.entry(fmt.Sprintf("%s_asymmetricpause", prefix)), prometheus.GaugeValue, asymmetricPause, device)
}

// updatePortInfo generates port type metrics to indicate if the network devices supports Twisted Pair, optical fiber, etc.
// The bit offsets here correspond to ethtool_link_mode_bit_indices in linux/include/uapi/linux/ethtool.h
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/uapi/linux/ethtool.h
func (c *ethtoolCollector) updatePortInfo(ch chan<- prometheus.Metric, device string, linkModes uint32) {
	for name, bit := range map[string]int{
		"TP":        unix.ETHTOOL_LINK_MODE_TP_BIT,
		"AUI":       unix.ETHTOOL_LINK_MODE_AUI_BIT,
		"MII":       unix.ETHTOOL_LINK_MODE_MII_BIT,
		"FIBRE":     unix.ETHTOOL_LINK_MODE_FIBRE_BIT,
		"BNC":       unix.ETHTOOL_LINK_MODE_BNC_BIT,
		"Backplane": unix.ETHTOOL_LINK_MODE_Backplane_BIT,
	} {
		if linkModes&(1<<bit) != 0 {
			ch <- prometheus.MustNewConstMetric(c.entry("supported_port"), prometheus.GaugeValue, 1.0, device, name)
		}

	}
}

// updateSpeeds generates metrics corresponding to the speeds and duplex modes supported or advertised by the network device.
// The bit offsets here correspond to ethtool_link_mode_bit_indices in linux/include/uapi/linux/ethtool.h
// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/include/uapi/linux/ethtool.h
func (c *ethtoolCollector) updateSpeeds(ch chan<- prometheus.Metric, prefix string, device string, linkModes uint32) {
	linkMode := fmt.Sprintf("%s_speed", prefix)
	const (
		full = "full"
		half = "half"
		// This is in _bytes_ to match bytes-per-second speeds from netclass.
		Mbps = 1000000.0 / 8.0
	)

	for bit, labels := range map[int]struct {
		speed  int
		duplex string
		phy    string
	}{
		unix.ETHTOOL_LINK_MODE_10baseT_Half_BIT:               {10, half, "T"},
		unix.ETHTOOL_LINK_MODE_10baseT_Full_BIT:               {10, full, "T"},
		unix.ETHTOOL_LINK_MODE_100baseT_Half_BIT:              {100, half, "T"},
		unix.ETHTOOL_LINK_MODE_100baseT_Full_BIT:              {100, full, "T"},
		unix.ETHTOOL_LINK_MODE_1000baseT_Half_BIT:             {1000, half, "T"},
		unix.ETHTOOL_LINK_MODE_1000baseT_Full_BIT:             {1000, full, "T"},
		unix.ETHTOOL_LINK_MODE_10000baseT_Full_BIT:            {10000, full, "T"},
		unix.ETHTOOL_LINK_MODE_2500baseT_Full_BIT:             {2500, full, "T"},
		unix.ETHTOOL_LINK_MODE_1000baseKX_Full_BIT:            {1000, full, "KX"},
		unix.ETHTOOL_LINK_MODE_10000baseKX4_Full_BIT:          {10000, full, "KX4"},
		unix.ETHTOOL_LINK_MODE_10000baseKR_Full_BIT:           {10000, full, "KR"},
		unix.ETHTOOL_LINK_MODE_10000baseR_FEC_BIT:             {10000, full, "R_FEC"},
		unix.ETHTOOL_LINK_MODE_20000baseMLD2_Full_BIT:         {20000, full, "MLD2"},
		unix.ETHTOOL_LINK_MODE_20000baseKR2_Full_BIT:          {20000, full, "KR2"},
		unix.ETHTOOL_LINK_MODE_40000baseKR4_Full_BIT:          {40000, full, "KR4"},
		unix.ETHTOOL_LINK_MODE_40000baseCR4_Full_BIT:          {40000, full, "CR4"},
		unix.ETHTOOL_LINK_MODE_40000baseSR4_Full_BIT:          {40000, full, "SR4"},
		unix.ETHTOOL_LINK_MODE_40000baseLR4_Full_BIT:          {40000, full, "LR4"},
		unix.ETHTOOL_LINK_MODE_56000baseKR4_Full_BIT:          {56000, full, "KR4"},
		unix.ETHTOOL_LINK_MODE_56000baseCR4_Full_BIT:          {56000, full, "CR4"},
		unix.ETHTOOL_LINK_MODE_56000baseSR4_Full_BIT:          {56000, full, "SR4"},
		unix.ETHTOOL_LINK_MODE_56000baseLR4_Full_BIT:          {56000, full, "LR4"},
		unix.ETHTOOL_LINK_MODE_25000baseCR_Full_BIT:           {25000, full, "CR"},
		unix.ETHTOOL_LINK_MODE_25000baseKR_Full_BIT:           {25000, full, "KR"},
		unix.ETHTOOL_LINK_MODE_25000baseSR_Full_BIT:           {25000, full, "SR"},
		unix.ETHTOOL_LINK_MODE_50000baseCR2_Full_BIT:          {50000, full, "CR2"},
		unix.ETHTOOL_LINK_MODE_50000baseKR2_Full_BIT:          {50000, full, "KR2"},
		unix.ETHTOOL_LINK_MODE_100000baseKR4_Full_BIT:         {100000, full, "KR4"},
		unix.ETHTOOL_LINK_MODE_100000baseSR4_Full_BIT:         {100000, full, "SR4"},
		unix.ETHTOOL_LINK_MODE_100000baseCR4_Full_BIT:         {100000, full, "CR4"},
		unix.ETHTOOL_LINK_MODE_100000baseLR4_ER4_Full_BIT:     {100000, full, "R4_ER4"},
		unix.ETHTOOL_LINK_MODE_50000baseSR2_Full_BIT:          {50000, full, "SR2"},
		unix.ETHTOOL_LINK_MODE_1000baseX_Full_BIT:             {1000, full, "X"},
		unix.ETHTOOL_LINK_MODE_10000baseCR_Full_BIT:           {10000, full, "CR"},
		unix.ETHTOOL_LINK_MODE_10000baseSR_Full_BIT:           {10000, full, "SR"},
		unix.ETHTOOL_LINK_MODE_10000baseLR_Full_BIT:           {10000, full, "LR"},
		unix.ETHTOOL_LINK_MODE_10000baseLRM_Full_BIT:          {10000, full, "LRM"},
		unix.ETHTOOL_LINK_MODE_10000baseER_Full_BIT:           {10000, full, "ER"},
		unix.ETHTOOL_LINK_MODE_5000baseT_Full_BIT:             {5000, full, "T"},
		unix.ETHTOOL_LINK_MODE_50000baseKR_Full_BIT:           {50000, full, "KR"},
		unix.ETHTOOL_LINK_MODE_50000baseSR_Full_BIT:           {50000, full, "SR"},
		unix.ETHTOOL_LINK_MODE_50000baseCR_Full_BIT:           {50000, full, "CR"},
		unix.ETHTOOL_LINK_MODE_50000baseLR_ER_FR_Full_BIT:     {50000, full, "LR_ER_FR"},
		unix.ETHTOOL_LINK_MODE_50000baseDR_Full_BIT:           {50000, full, "DR"},
		unix.ETHTOOL_LINK_MODE_100000baseKR2_Full_BIT:         {100000, full, "KR2"},
		unix.ETHTOOL_LINK_MODE_100000baseSR2_Full_BIT:         {100000, full, "SR2"},
		unix.ETHTOOL_LINK_MODE_100000baseCR2_Full_BIT:         {100000, full, "CR2"},
		unix.ETHTOOL_LINK_MODE_100000baseLR2_ER2_FR2_Full_BIT: {100000, full, "LR2_ER2_FR2"},
		unix.ETHTOOL_LINK_MODE_100000baseDR2_Full_BIT:         {100000, full, "DR2"},
		unix.ETHTOOL_LINK_MODE_200000baseKR4_Full_BIT:         {200000, full, "KR4"},
		unix.ETHTOOL_LINK_MODE_200000baseSR4_Full_BIT:         {200000, full, "SR4"},
		unix.ETHTOOL_LINK_MODE_200000baseLR4_ER4_FR4_Full_BIT: {200000, full, "LR4_ER4_FR4"},
		unix.ETHTOOL_LINK_MODE_200000baseDR4_Full_BIT:         {200000, full, "DR4"},
		unix.ETHTOOL_LINK_MODE_200000baseCR4_Full_BIT:         {200000, full, "CR4"},
		unix.ETHTOOL_LINK_MODE_100baseT1_Full_BIT:             {100, full, "T1"},
		unix.ETHTOOL_LINK_MODE_1000baseT1_Full_BIT:            {1000, full, "T1"},
		unix.ETHTOOL_LINK_MODE_400000baseKR8_Full_BIT:         {400000, full, "KR8"},
		unix.ETHTOOL_LINK_MODE_400000baseSR8_Full_BIT:         {400000, full, "SR8"},
		unix.ETHTOOL_LINK_MODE_400000baseLR8_ER8_FR8_Full_BIT: {400000, full, "LR8_ER8_FR8"},
		unix.ETHTOOL_LINK_MODE_400000baseDR8_Full_BIT:         {400000, full, "DR8"},
		unix.ETHTOOL_LINK_MODE_400000baseCR8_Full_BIT:         {400000, full, "CR8"},
		unix.ETHTOOL_LINK_MODE_100000baseKR_Full_BIT:          {100000, full, "KR"},
		unix.ETHTOOL_LINK_MODE_100000baseSR_Full_BIT:          {100000, full, "SR"},
		unix.ETHTOOL_LINK_MODE_100000baseLR_ER_FR_Full_BIT:    {100000, full, "LR_ER_FR"},
		unix.ETHTOOL_LINK_MODE_100000baseCR_Full_BIT:          {100000, full, "CR"},
		unix.ETHTOOL_LINK_MODE_100000baseDR_Full_BIT:          {100000, full, "DR"},
		unix.ETHTOOL_LINK_MODE_200000baseKR2_Full_BIT:         {200000, full, "KR2"},
		unix.ETHTOOL_LINK_MODE_200000baseSR2_Full_BIT:         {200000, full, "SR2"},
		unix.ETHTOOL_LINK_MODE_200000baseLR2_ER2_FR2_Full_BIT: {200000, full, "LR2_ER2_FR2"},
		unix.ETHTOOL_LINK_MODE_200000baseDR2_Full_BIT:         {200000, full, "DR2"},
		unix.ETHTOOL_LINK_MODE_200000baseCR2_Full_BIT:         {200000, full, "CR2"},
		unix.ETHTOOL_LINK_MODE_400000baseKR4_Full_BIT:         {400000, full, "KR4"},
		unix.ETHTOOL_LINK_MODE_400000baseSR4_Full_BIT:         {400000, full, "SR4"},
		unix.ETHTOOL_LINK_MODE_400000baseLR4_ER4_FR4_Full_BIT: {400000, full, "LR4_ER4_FR4"},
		unix.ETHTOOL_LINK_MODE_400000baseDR4_Full_BIT:         {400000, full, "DR4"},
		unix.ETHTOOL_LINK_MODE_400000baseCR4_Full_BIT:         {400000, full, "CR4"},
		unix.ETHTOOL_LINK_MODE_100baseFX_Half_BIT:             {100, half, "FX"},
		unix.ETHTOOL_LINK_MODE_100baseFX_Full_BIT:             {100, full, "FX"},
	} {
		if linkModes&(1<<bit) != 0 {
			ch <- prometheus.MustNewConstMetric(c.entry(linkMode), prometheus.GaugeValue,
				float64(labels.speed)*Mbps, device, labels.duplex, fmt.Sprintf("%dbase%s", labels.speed, labels.phy))
		}
	}
}

func (c *ethtoolCollector) Update(ch chan<- prometheus.Metric) error {
	netClass, err := c.fs.NetClassDevices()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			c.logger.Debug("Could not read netclass file", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("could not get net class info: %w", err)
	}

	if len(netClass) == 0 {
		return fmt.Errorf("no network devices found")
	}

	for _, device := range netClass {
		var stats map[string]uint64
		var err error

		if c.deviceFilter.ignored(device) {
			continue
		}

		linkInfo, err := c.ethtool.LinkInfo(device)
		if err == nil {
			c.updateSpeeds(ch, "supported", device, linkInfo.Supported)
			c.updatePortInfo(ch, device, linkInfo.Supported)
			c.updatePortCapabilities(ch, "supported", device, linkInfo.Supported)
			c.updateSpeeds(ch, "advertised", device, linkInfo.Advertising)
			c.updatePortCapabilities(ch, "advertised", device, linkInfo.Advertising)
			ch <- prometheus.MustNewConstMetric(c.entry("autonegotiate"), prometheus.GaugeValue, float64(linkInfo.Autoneg), device)
		} else {
			if errno, ok := err.(syscall.Errno); ok {
				if err == unix.EOPNOTSUPP {
					c.logger.Debug("ethtool link info error", "err", err, "device", device, "errno", uint(errno))
				} else if errno != 0 {
					c.logger.Error("ethtool link info error", "err", err, "device", device, "errno", uint(errno))
				}
			} else {
				c.logger.Error("ethtool link info error", "err", err, "device", device)
			}
		}

		drvInfo, err := c.ethtool.DriverInfo(device)

		if err == nil {
			ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1.0,
				drvInfo.BusInfo, device, drvInfo.Driver, drvInfo.EromVersion, drvInfo.FwVersion, drvInfo.Version)
		} else {
			if errno, ok := err.(syscall.Errno); ok {
				if err == unix.EOPNOTSUPP {
					c.logger.Debug("ethtool driver info error", "err", err, "device", device, "errno", uint(errno))
				} else if errno != 0 {
					c.logger.Error("ethtool driver info error", "err", err, "device", device, "errno", uint(errno))
				}
			} else {
				c.logger.Error("ethtool driver info error", "err", err, "device", device)
			}
		}

		stats, err = c.ethtool.Stats(device)

		// If Stats() returns EOPNOTSUPP it doesn't support ethtool stats. Log that only at Debug level.
		// Otherwise log it at Error level.
		if err != nil {
			if errno, ok := err.(syscall.Errno); ok {
				if err == unix.EOPNOTSUPP {
					c.logger.Debug("ethtool stats error", "err", err, "device", device, "errno", uint(errno))
				} else if errno != 0 {
					c.logger.Error("ethtool stats error", "err", err, "device", device, "errno", uint(errno))
				}
			} else {
				c.logger.Error("ethtool stats error", "err", err, "device", device)
			}
		}

		if stats == nil || len(stats) < 1 {
			// No stats returned; device does not support ethtool stats.
			continue
		}

		// Sanitizing the metric names can lead to duplicate metric names. Therefore check for clashes beforehand.
		metricFQNames := make(map[string]string)
		for metric := range stats {
			metricName := SanitizeMetricName(metric)
			if !c.metricsPattern.MatchString(metricName) {
				continue
			}
			metricFQName := buildEthtoolFQName(metricName)
			existingMetric, exists := metricFQNames[metricFQName]
			if exists {
				c.logger.Debug("dropping duplicate metric name", "device", device,
					"metricFQName", metricFQName, "metric1", existingMetric, "metric2", metricName)
				// Keep the metricName as "deleted" in the dict in case there are 3 duplicates.
				metricFQNames[metricFQName] = ""
			} else {
				metricFQNames[metricFQName] = metricName
			}
		}

		// Sort metric names so that the test fixtures will match up
		keys := make([]string, 0, len(metricFQNames))
		for k := range metricFQNames {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, metricFQName := range keys {
			metric := metricFQNames[metricFQName]
			if metric == "" {
				// Skip the "deleted" duplicate metrics
				continue
			}

			val := stats[metric]

			// Check to see if this metric exists; if not then create it and store it in c.entries.
			entry := c.entryWithCreate(metric, metricFQName)
			ch <- prometheus.MustNewConstMetric(
				entry, prometheus.UntypedValue, float64(val), device)
		}
	}

	return nil
}

func (c *ethtoolCollector) entryWithCreate(key, metricFQName string) *prometheus.Desc {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()

	if _, ok := c.entries[key]; !ok {
		c.entries[key] = prometheus.NewDesc(
			metricFQName,
			fmt.Sprintf("Network interface %s", key),
			[]string{"device"}, nil,
		)
	}

	return c.entries[key]
}

func (c *ethtoolCollector) entry(key string) *prometheus.Desc {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()
	return c.entries[key]
}
