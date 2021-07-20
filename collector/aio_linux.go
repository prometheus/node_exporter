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

// +build !noaio

package collector

import (
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type aioCollector struct {
	current *prometheus.Desc
	limit   *prometheus.Desc
	logger  log.Logger
}

func init() {
	registerCollector("aio", defaultDisabled, NewAioCollector)
}

// NewAioCollector returns a new Collector exposing aio stats.
func NewAioCollector(logger log.Logger) (Collector, error) {
	return &aioCollector{
		current: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "aio_nr"),
			"Number of currently active asynchronous io contexts.",
			nil, nil,
		),
		limit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "aio_max_nr"),
			"Maximum size of asynchronous io contexts.",
			nil, nil,
		),
		logger: logger,
	}, nil
}

func (c *aioCollector) Update(ch chan<- prometheus.Metric) error {
	value, err := readUintFromFile(procFilePath("sys/fs/aio-nr"))
	if err != nil {
		return nil
	}
	ch <- prometheus.MustNewConstMetric(
		c.current, prometheus.GaugeValue, float64(value))

	value, err = readUintFromFile(procFilePath("sys/fs/aio-max-nr"))
	if err != nil {
		return nil
	}
	ch <- prometheus.MustNewConstMetric(
		c.limit, prometheus.GaugeValue, float64(value))

	return nil
}
