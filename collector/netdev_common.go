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

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	netdevDeviceInclude    = kingpin.Flag("collector.netdev.device-include", "Regexp of net devices to include (mutually exclusive to device-exclude).").String()
	oldNetdevDeviceInclude = kingpin.Flag("collector.netdev.device-whitelist", "DEPRECATED: Use collector.netdev.device-include").Hidden().String()
	netdevDeviceExclude    = kingpin.Flag("collector.netdev.device-exclude", "Regexp of net devices to exclude (mutually exclusive to device-include).").String()
	oldNetdevDeviceExclude = kingpin.Flag("collector.netdev.device-blacklist", "DEPRECATED: Use collector.netdev.device-exclude").Hidden().String()
)

type netDevCollector struct {
	subsystem    string
	deviceFilter netDevFilter
	metricDescs  map[string]*prometheus.Desc
	logger       log.Logger
}

type netDevStats map[string]map[string]uint64

func init() {
	registerCollector("netdev", defaultEnabled, NewNetDevCollector)
}

// NewNetDevCollector returns a new Collector exposing network device stats.
func NewNetDevCollector(logger log.Logger) (Collector, error) {
	if *oldNetdevDeviceInclude != "" {
		if *netdevDeviceInclude == "" {
			level.Warn(logger).Log("msg", "--collector.netdev.device-whitelist is DEPRECATED and will be removed in 2.0.0, use --collector.netdev.device-include")
			*netdevDeviceInclude = *oldNetdevDeviceInclude
		} else {
			return nil, errors.New("--collector.netdev.device-whitelist and --collector.netdev.device-include are mutually exclusive")
		}
	}

	if *oldNetdevDeviceExclude != "" {
		if *netdevDeviceExclude == "" {
			level.Warn(logger).Log("msg", "--collector.netdev.device-blacklist is DEPRECATED and will be removed in 2.0.0, use --collector.netdev.device-exclude")
			*netdevDeviceExclude = *oldNetdevDeviceExclude
		} else {
			return nil, errors.New("--collector.netdev.device-blacklist and --collector.netdev.device-exclude are mutually exclusive")
		}
	}

	if *netdevDeviceExclude != "" && *netdevDeviceInclude != "" {
		return nil, errors.New("device-exclude & device-include are mutually exclusive")
	}

	if *netdevDeviceExclude != "" {
		level.Info(logger).Log("msg", "Parsed flag --collector.netdev.device-exclude", "flag", *netdevDeviceExclude)
	}

	if *netdevDeviceInclude != "" {
		level.Info(logger).Log("msg", "Parsed Flag --collector.netdev.device-include", "flag", *netdevDeviceInclude)
	}

	return &netDevCollector{
		subsystem:    "network",
		deviceFilter: newNetDevFilter(*netdevDeviceExclude, *netdevDeviceInclude),
		metricDescs:  map[string]*prometheus.Desc{},
		logger:       logger,
	}, nil
}

func (c *netDevCollector) Update(ch chan<- prometheus.Metric) error {
	netDev, err := getNetDevStats(&c.deviceFilter, c.logger)
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %w", err)
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
			ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, float64(value), dev)
		}
	}
	return nil
}
