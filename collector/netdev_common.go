// Copyright 2015 The Prometheus Authors
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

//go:build !nonetdev && (linux || freebsd || openbsd || dragonfly || darwin || aix)

package collector

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/node_exporter/config"
)

type netDevCollector struct {
	subsystem        string
	deviceFilter     deviceFilter
	metricDescsMutex sync.Mutex
	metricDescs      map[string]*prometheus.Desc
	logger           *slog.Logger
	config           config.NetDevConfig
	paths            pathConfig
}

type netDevStats map[string]map[string]uint64

func init() {
	registerCollector("netdev", defaultEnabled, NewNetDevCollector)
}

// NewNetDevCollector returns a new Collector exposing network device stats.
func NewNetDevCollector(cfg config.Config, logger *slog.Logger) (Collector, error) {
	if _, _, err := cfg.NetDev.DeviceFilter(); err != nil {
		return nil, err
	}

	if cfg.NetDev.DeviceExclude != "" {
		logger.Info("Parsed flag --collector.netdev.device-exclude", "flag", cfg.NetDev.DeviceExclude)
	}

	if cfg.NetDev.DeviceInclude != "" {
		logger.Info("Parsed Flag --collector.netdev.device-include", "flag", cfg.NetDev.DeviceInclude)
	}

	return &netDevCollector{
		paths:        pathConfig(cfg.Path),
		subsystem:    "network",
		deviceFilter: newDeviceFilter(cfg.NetDev.DeviceExclude, cfg.NetDev.DeviceInclude),
		metricDescs:  map[string]*prometheus.Desc{},
		logger:       logger,
		config:       cfg.NetDev,
	}, nil
}

func (c *netDevCollector) metricDesc(key string, labels []string) *prometheus.Desc {
	c.metricDescsMutex.Lock()
	defer c.metricDescsMutex.Unlock()

	if _, ok := c.metricDescs[key]; !ok {
		c.metricDescs[key] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, key+"_total"),
			fmt.Sprintf("Network device statistic %s.", key),
			labels,
			nil,
		)
	}

	return c.metricDescs[key]
}

func (c *netDevCollector) Update(_ context.Context, ch chan<- prometheus.Metric) error {
	netDev, err := getNetDevStats(c.config, c.paths, &c.deviceFilter, c.logger)
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %w", err)
	}

	netDevLabels, err := getNetDevLabels(c.config, c.paths)
	if err != nil {
		return fmt.Errorf("couldn't get netdev labels: %w", err)
	}

	for dev, devStats := range netDev {
		if !c.config.DetailedMetrics {
			legacy(devStats)
		}

		labels := []string{"device"}
		labelValues := []string{dev}
		if devLabels, exists := netDevLabels[dev]; exists {
			for labelName, labelValue := range devLabels {
				labels = append(labels, labelName)
				labelValues = append(labelValues, labelValue)
			}
		}

		for key, value := range devStats {
			desc := c.metricDesc(key, labels)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, float64(value), labelValues...)
		}
	}
	if c.config.AddressInfo {
		interfaces, err := net.Interfaces()
		if err != nil {
			return fmt.Errorf("could not get network interfaces: %w", err)
		}

		desc := prometheus.NewDesc(prometheus.BuildFQName(namespace, "network_address",
			"info"), "node network address by device",
			[]string{"device", "address", "netmask", "scope"}, nil)

		for _, addr := range getAddrsInfo(interfaces, &c.deviceFilter, c.logger) {
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1,
				addr.device, addr.addr, addr.netmask, addr.scope)
		}
	}
	return nil
}

type addrInfo struct {
	device  string
	addr    string
	scope   string
	netmask string
}

func scope(ip net.IP) string {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return "link-local"
	}

	if ip.IsInterfaceLocalMulticast() {
		return "interface-local"
	}

	if ip.IsGlobalUnicast() {
		return "global"
	}

	return ""
}

// getAddrsInfo returns interface name, address, scope and netmask for all interfaces.
func getAddrsInfo(interfaces []net.Interface, filter *deviceFilter, logger *slog.Logger) []addrInfo {
	var res []addrInfo

	for _, ifs := range interfaces {
		if filter.ignored(ifs.Name) {
			logger.Debug("Ignoring device", "device", ifs.Name)
			continue
		}
		addrs, _ := ifs.Addrs()
		for _, addr := range addrs {
			ip, ipNet, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			size, _ := ipNet.Mask.Size()

			res = append(res, addrInfo{
				device:  ifs.Name,
				addr:    ip.String(),
				scope:   scope(ip),
				netmask: strconv.Itoa(size),
			})
		}
	}

	return res
}

// https://github.com/torvalds/linux/blob/master/net/core/net-procfs.c#L75-L97
func legacy(metrics map[string]uint64) {
	if metric, ok := pop(metrics, "receive_errors"); ok {
		metrics["receive_errs"] = metric
	}
	if metric, ok := pop(metrics, "receive_dropped"); ok {
		metrics["receive_drop"] = metric + popz(metrics, "receive_missed_errors")
	}
	if metric, ok := pop(metrics, "receive_fifo_errors"); ok {
		metrics["receive_fifo"] = metric
	}
	if metric, ok := pop(metrics, "receive_frame_errors"); ok {
		metrics["receive_frame"] = metric + popz(metrics, "receive_length_errors") + popz(metrics, "receive_over_errors") + popz(metrics, "receive_crc_errors")
	}
	if metric, ok := pop(metrics, "multicast"); ok {
		metrics["receive_multicast"] = metric
	}
	if metric, ok := pop(metrics, "transmit_errors"); ok {
		metrics["transmit_errs"] = metric
	}
	if metric, ok := pop(metrics, "transmit_dropped"); ok {
		metrics["transmit_drop"] = metric
	}
	if metric, ok := pop(metrics, "transmit_fifo_errors"); ok {
		metrics["transmit_fifo"] = metric
	}
	if metric, ok := pop(metrics, "multicast"); ok {
		metrics["receive_multicast"] = metric
	}
	if metric, ok := pop(metrics, "collisions"); ok {
		metrics["transmit_colls"] = metric
	}
	if metric, ok := pop(metrics, "transmit_carrier_errors"); ok {
		metrics["transmit_carrier"] = metric + popz(metrics, "transmit_aborted_errors") + popz(metrics, "transmit_heartbeat_errors") + popz(metrics, "transmit_window_errors")
	}
}

func pop(m map[string]uint64, key string) (uint64, bool) {
	value, ok := m[key]
	delete(m, key)
	return value, ok
}

func popz(m map[string]uint64, key string) uint64 {
	if value, ok := m[key]; ok {
		delete(m, key)
		return value
	}
	return 0
}
