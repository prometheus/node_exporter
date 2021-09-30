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
	"os"
	"regexp"
	"sort"
	"strings"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
	"github.com/safchain/ethtool"
	"golang.org/x/sys/unix"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	ethtoolIgnoredDevices  = kingpin.Flag("collector.ethtool.ignored-devices", "Regexp of net devices to ignore for ethtool collector.").Default("^$").String()
	ethtoolIncludedMetrics = kingpin.Flag("collector.ethtool.metrics-include", "Regexp of ethtool stats to include.").Default(".*").String()
	metricNameRegex        = regexp.MustCompile(`_*[^0-9A-Za-z_]+_*`)
	receivedRegex          = regexp.MustCompile(`(^|_)rx(_|$)`)
	transmittedRegex       = regexp.MustCompile(`(^|_)tx(_|$)`)
)

type Ethtool interface {
	DriverInfo(string) (ethtool.DrvInfo, error)
	Stats(string) (map[string]uint64, error)
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

type ethtoolCollector struct {
	fs                    sysfs.FS
	entries               map[string]*prometheus.Desc
	ethtool               Ethtool
	ignoredDevicesPattern *regexp.Regexp
	infoDesc              *prometheus.Desc
	metricsPattern        *regexp.Regexp
	logger                log.Logger
}

// makeEthtoolCollector is the internal constructor for EthtoolCollector.
// This allows NewEthtoolTestCollector to override its .ethtool interface
// for testing.
func makeEthtoolCollector(logger log.Logger) (*ethtoolCollector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ethtool library: %w", err)
	}

	// Pre-populate some common ethtool metrics.
	return &ethtoolCollector{
		fs:                    fs,
		ethtool:               &ethtoolLibrary{e},
		ignoredDevicesPattern: regexp.MustCompile(*ethtoolIgnoredDevices),
		metricsPattern:        regexp.MustCompile(*ethtoolIncludedMetrics),
		logger:                logger,
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

// Sanitize the given metric name by replacing invalid characters by underscores.
//
// OpenMetrics and the Prometheus exposition format require the metric name
// to consist only of alphanumericals and "_", ":" and they must not start
// with digits. Since colons in MetricFamily are reserved to signal that the
// MetricFamily is the result of a calculation or aggregation of a general
// purpose monitoring system, colons will be replaced as well.
//
// Note: If not subsequently prepending a namespace and/or subsystem (e.g.,
// with prometheus.BuildFQName), the caller must ensure that the supplied
// metricName does not begin with a digit.
func SanitizeMetricName(metricName string) string {
	return metricNameRegex.ReplaceAllString(metricName, "_")
}

// Generate the fully-qualified metric name for the ethool metric.
func buildEthtoolFQName(metric string) string {
	metricName := strings.TrimLeft(strings.ToLower(SanitizeMetricName(metric)), "_")
	metricName = receivedRegex.ReplaceAllString(metricName, "${1}received${2}")
	metricName = transmittedRegex.ReplaceAllString(metricName, "${1}transmitted${2}")
	return prometheus.BuildFQName(namespace, "ethtool", metricName)
}

// NewEthtoolCollector returns a new Collector exposing ethtool stats.
func NewEthtoolCollector(logger log.Logger) (Collector, error) {
	return makeEthtoolCollector(logger)
}

func (c *ethtoolCollector) Update(ch chan<- prometheus.Metric) error {
	netClass, err := c.fs.NetClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			level.Debug(c.logger).Log("msg", "Could not read netclass file", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("could not get net class info: %w", err)
	}

	if len(netClass) == 0 {
		return fmt.Errorf("no network devices found")
	}

	for device := range netClass {
		var stats map[string]uint64
		var err error

		if c.ignoredDevicesPattern.MatchString(device) {
			continue
		}

		drvInfo, err := c.ethtool.DriverInfo(device)

		if err == nil {
			ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1.0,
				drvInfo.BusInfo, device, drvInfo.Driver, drvInfo.EromVersion, drvInfo.FwVersion, drvInfo.Version)
		} else {
			if errno, ok := err.(syscall.Errno); ok {
				if err == unix.EOPNOTSUPP {
					level.Debug(c.logger).Log("msg", "ethtool driver info error", "err", err, "device", device, "errno", uint(errno))
				} else if errno != 0 {
					level.Error(c.logger).Log("msg", "ethtool driver info error", "err", err, "device", device, "errno", uint(errno))
				}
			} else {
				level.Error(c.logger).Log("msg", "ethtool driver info error", "err", err, "device", device)
			}
		}

		stats, err = c.ethtool.Stats(device)

		// If Stats() returns EOPNOTSUPP it doesn't support ethtool stats. Log that only at Debug level.
		// Otherwise log it at Error level.
		if err != nil {
			if errno, ok := err.(syscall.Errno); ok {
				if err == unix.EOPNOTSUPP {
					level.Debug(c.logger).Log("msg", "ethtool stats error", "err", err, "device", device, "errno", uint(errno))
				} else if errno != 0 {
					level.Error(c.logger).Log("msg", "ethtool stats error", "err", err, "device", device, "errno", uint(errno))
				}
			} else {
				level.Error(c.logger).Log("msg", "ethtool stats error", "err", err, "device", device)
			}
		}

		if stats == nil || len(stats) < 1 {
			// No stats returned; device does not support ethtool stats.
			continue
		}

		// Sort metric names so that the test fixtures will match up
		keys := make([]string, 0, len(stats))
		for k := range stats {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, metric := range keys {
			if !c.metricsPattern.MatchString(metric) {
				continue
			}
			val := stats[metric]
			metricFQName := buildEthtoolFQName(metric)

			// Check to see if this metric exists; if not then create it and store it in c.entries.
			entry, exists := c.entries[metric]
			if !exists {
				entry = prometheus.NewDesc(
					metricFQName,
					fmt.Sprintf("Network interface %s", metric),
					[]string{"device"}, nil,
				)
				c.entries[metric] = entry
			}
			ch <- prometheus.MustNewConstMetric(
				entry, prometheus.UntypedValue, float64(val), device)
		}
	}

	return nil
}
