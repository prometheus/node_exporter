// Copyright 2019 The Prometheus Authors
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

// +build !noshedstat

package collector

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

const nsPerSec = 1e9

var (
	runningSecondsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "schedstat", "running_seconds_total"),
		"Number of seconds CPU spent running a process.",
		[]string{"cpu"},
		nil,
	)

	waitingSecondsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "schedstat", "waiting_seconds_total"),
		"Number of seconds spent by processing waiting for this CPU.",
		[]string{"cpu"},
		nil,
	)

	timeslicesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "schedstat", "timeslices_total"),
		"Number of timeslices executed by CPU.",
		[]string{"cpu"},
		nil,
	)
)

// NewSchedstatCollector returns a new Collector exposing task scheduler statistics
func NewSchedstatCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &schedstatCollector{fs, logger}, nil
}

type schedstatCollector struct {
	fs     procfs.FS
	logger log.Logger
}

func init() {
	registerCollector("schedstat", defaultEnabled, NewSchedstatCollector)
}

func (c *schedstatCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := c.fs.Schedstat()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "schedstat file does not exist")
			return ErrNoData
		}
		return err
	}

	for _, cpu := range stats.CPUs {
		ch <- prometheus.MustNewConstMetric(
			runningSecondsTotal,
			prometheus.CounterValue,
			float64(cpu.RunningNanoseconds)/nsPerSec,
			cpu.CPUNum,
		)

		ch <- prometheus.MustNewConstMetric(
			waitingSecondsTotal,
			prometheus.CounterValue,
			float64(cpu.WaitingNanoseconds)/nsPerSec,
			cpu.CPUNum,
		)

		ch <- prometheus.MustNewConstMetric(
			timeslicesTotal,
			prometheus.CounterValue,
			float64(cpu.RunTimeslices),
			cpu.CPUNum,
		)
	}

	return nil
}
