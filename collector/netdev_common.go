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

// +build !nonetdev
// +build linux freebsd openbsd dragonfly darwin

package collector

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	netdevIgnoredDevices = kingpin.Flag("collector.netdev.device-blacklist", "Regexp of net devices to blacklist (mutually exclusive to device-whitelist).").String()
	netdevAcceptDevices  = kingpin.Flag("collector.netdev.device-whitelist", "Regexp of net devices to whitelist (mutually exclusive to device-blacklist).").String()
)

type netDevCollector struct {
	subsystem             string
	ignoredDevicesPattern *regexp.Regexp
	acceptDevicesPattern  *regexp.Regexp
	metricDescs           map[string]*prometheus.Desc
}

func init() {
	registerCollector("netdev", defaultEnabled, NewNetDevCollector)
}

// NewNetDevCollector returns a new Collector exposing network device stats.
func NewNetDevCollector() (Collector, error) {
	if *netdevIgnoredDevices != "" && *netdevAcceptDevices != "" {
		return nil, errors.New("device-blacklist & accept-devices are mutually exclusive")
	}

	var ignorePattern *regexp.Regexp = nil
	if *netdevIgnoredDevices != "" {
		ignorePattern = regexp.MustCompile(*netdevIgnoredDevices)
	}

	var acceptPattern *regexp.Regexp = nil
	if *netdevAcceptDevices != "" {
		acceptPattern = regexp.MustCompile(*netdevAcceptDevices)
	}

	return &netDevCollector{
		subsystem:             "network",
		ignoredDevicesPattern: ignorePattern,
		acceptDevicesPattern:  acceptPattern,
		metricDescs:           map[string]*prometheus.Desc{},
	}, nil
}

func (c *netDevCollector) Update(ch chan<- prometheus.Metric) error {
	netDev, err := getNetDevStats(c.ignoredDevicesPattern, c.acceptDevicesPattern)
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %s", err)
	}
	for dev, devStats := range netDev {
		for key, value := range devStats {
			desc, ok := c.metricDescs[key]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, c.subsystem, key+"_total"),
					fmt.Sprintf("Network device statistic %s.", key),
					[]string{"device"},
					nil,
				)
				c.metricDescs[key] = desc
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in netstats: %s", value, err)
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, v, dev)
		}
	}
	return nil
}
