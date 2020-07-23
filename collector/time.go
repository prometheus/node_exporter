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

// +build !notime

package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type timeCollector struct {
	desc *prometheus.Desc
}

func init() {
	registerCollector("time", defaultEnabled, NewTimeCollector)
}

// NewTimeCollector returns a new Collector exposing the current system time in
// seconds since epoch.
func NewTimeCollector() (Collector, error) {
	return &timeCollector{
		desc: prometheus.NewDesc(
			namespace+"_time_seconds",
			"System time in seconds since epoch (1970).",
			nil, nil,
		),
	}, nil
}

func (c *timeCollector) Update(ch chan<- prometheus.Metric) error {
	now := float64(time.Now().UnixNano()) / 1e9
	log.Debugf("Return time: %f", now)
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, now)
	return nil
}
