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

//go:build !nosoftirqs
// +build !nosoftirqs

package collector

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	softirqLabelNames = []string{"cpu", "type"}
)

func (c *softirqsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	softirqs, err := c.fs.Softirqs()
	if err != nil {
		return fmt.Errorf("couldn't get softirqs: %w", err)
	}

	for cpuNo, value := range softirqs.Hi {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "HI")
	}
	for cpuNo, value := range softirqs.Timer {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "TIMER")
	}
	for cpuNo, value := range softirqs.NetTx {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "NET_TX")
	}
	for cpuNo, value := range softirqs.NetRx {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "NET_RX")
	}
	for cpuNo, value := range softirqs.Block {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "BLOCK")
	}
	for cpuNo, value := range softirqs.IRQPoll {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "IRQ_POLL")
	}
	for cpuNo, value := range softirqs.Tasklet {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "TASKLET")
	}
	for cpuNo, value := range softirqs.Sched {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "SCHED")
	}
	for cpuNo, value := range softirqs.HRTimer {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "HRTIMER")
	}
	for cpuNo, value := range softirqs.RCU {
		ch <- c.desc.mustNewConstMetric(float64(value), strconv.Itoa(cpuNo), "RCU")
	}

	return err
}
