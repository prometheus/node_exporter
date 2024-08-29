// Copyright 2023 The Prometheus Authors
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

//go:build !nonetisr
// +build !nonetisr

package collector

import (
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

type netisrCollector struct {
	sysctls []bsdSysctl
	logger  *slog.Logger
}

const (
	netisrCollectorSubsystem = "netisr"
)

func init() {
	registerCollector("netisr", defaultEnabled, NewNetisrCollector)
}

func NewNetisrCollector(logger *slog.Logger) (Collector, error) {
	return &netisrCollector{
		sysctls: []bsdSysctl{
			{
				name:        "numthreads",
				description: "netisr current thread count",
				mib:         "net.isr.numthreads",
				dataType:    bsdSysctlTypeUint32,
				valueType:   prometheus.GaugeValue,
			},
			{
				name:        "maxprot",
				description: "netisr maximum protocols",
				mib:         "net.isr.maxprot",
				dataType:    bsdSysctlTypeUint32,
				valueType:   prometheus.GaugeValue,
			},
			{
				name:        "defaultqlimit",
				description: "netisr default queue limit",
				mib:         "net.isr.defaultqlimit",
				dataType:    bsdSysctlTypeUint32,
				valueType:   prometheus.GaugeValue,
			},
			{
				name:        "maxqlimit",
				description: "netisr maximum queue limit",
				mib:         "net.isr.maxqlimit",
				dataType:    bsdSysctlTypeUint32,
				valueType:   prometheus.GaugeValue,
			},
			{
				name:        "bindthreads",
				description: "netisr threads bound to CPUs",
				mib:         "net.isr.bindthreads",
				dataType:    bsdSysctlTypeUint32,
				valueType:   prometheus.GaugeValue,
			},
			{
				name:        "maxthreads",
				description: "netisr maximum thread count",
				mib:         "net.isr.maxthreads",
				dataType:    bsdSysctlTypeUint32,
				valueType:   prometheus.GaugeValue,
			},
		},
		logger: logger,
	}, nil
}

func (c *netisrCollector) Update(ch chan<- prometheus.Metric) error {
	for _, m := range c.sysctls {
		v, err := m.Value()
		if err != nil {
			return fmt.Errorf("couldn't get sysctl: %w", err)
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, netisrCollectorSubsystem, m.name),
				m.description,
				nil, nil,
			), m.valueType, v)
	}

	return nil
}
