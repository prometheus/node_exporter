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

//go:build (darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris) && !noloadavg
// +build darwin dragonfly freebsd linux netbsd openbsd solaris
// +build !noloadavg

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type loadavgCollector struct {
	metric []typedDesc
	logger log.Logger
}

func init() {
	registerCollector("loadavg", defaultEnabled, NewLoadavgCollector)
}

// NewLoadavgCollector returns a new Collector exposing load average stats.
func NewLoadavgCollector(logger log.Logger) (Collector, error) {
	return &loadavgCollector{
		metric: []typedDesc{
			{prometheus.NewDesc(namespace+"_load1", "1m load average.", nil, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(namespace+"_load5", "5m load average.", nil, nil), prometheus.GaugeValue},
			{prometheus.NewDesc(namespace+"_load15", "15m load average.", nil, nil), prometheus.GaugeValue},
		},
		logger: logger,
	}, nil
}

func (c *loadavgCollector) Update(ch chan<- prometheus.Metric) error {
	loads, err := getLoad()
	if err != nil {
		return fmt.Errorf("couldn't get load: %w", err)
	}
	for i, load := range loads {
		level.Debug(c.logger).Log("msg", "return load", "index", i, "load", load)
		ch <- c.metric[i].mustNewConstMetric(load)
	}
	return err
}
