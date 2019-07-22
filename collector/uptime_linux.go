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

// +build !nouptime

package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"syscall"
)

type UptimeCollector struct {
	desc *prometheus.Desc
}

func init() {
	registerCollector("uptime", defaultEnabled, NewUptimeCollector)
}

// NewUptimeCollector returns a new Collector exposing the current node uptime in seconds.
func NewUptimeCollector() (Collector, error) {
	return &UptimeCollector{
		desc: prometheus.NewDesc(
			namespace+"_uptime_seconds",
			"Node uptime in seconds.",
			nil, nil,
		),
	}, nil
}

func (c *UptimeCollector) Update(ch chan<- prometheus.Metric) error {
	s := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(s)
	if err != nil {
		log.Errorf("Error reading uptime %s", err)
	}

	uptime := float64(s.Uptime)
	log.Debugf("Return uptime: %f", uptime)
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, uptime)
	return nil
}
