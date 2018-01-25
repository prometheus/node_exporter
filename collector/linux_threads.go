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

// +build !nothreads

package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
	"io/ioutil"
	"strconv"
	"strings"
)

type threadsCollector struct {
	threadAlloc *prometheus.Desc
	threadMax *prometheus.Desc
}

func init() {
	registerCollector("threads", defaultDisabled, NewThreadsCollector)
}

func NewThreadsCollector() (Collector, error) {
	return &threadsCollector{
		threadAlloc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "threads"),
			"Allocated threads in system",
			nil, nil,
		),
		threadMax: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "threads", "max"),
			"Limit of threads in the system",
			nil, nil,
		),
	}, nil
}
func (t *threadsCollector) Update(ch chan<- prometheus.Metric) error {
	val, err := getAllocatedThreads()
	if err != nil {
		return fmt.Errorf("Unable to retrieve number of allocated threads %v\n", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadAlloc, prometheus.GaugeValue, float64(val))
	maxThreads, err := getMaxThreads()
	if err != nil {
		return fmt.Errorf("Unable to retrieve limit number of threads %v\n", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadMax, prometheus.GaugeValue, float64(maxThreads))
	return nil
}

func getAllocatedThreads() (int, error) {
	p, err := procfs.AllProcs()
	if err != nil {
		return 0, err
	}
	thread := 0
	for _, pid := range p {
		stat, err := pid.NewStat()
		if err != nil {
			return 0, err
		}
		thread += stat.NumThreads

	}
	return thread, nil
}

func getMaxThreads() (int, error) {
	rawData, err := ioutil.ReadFile("/proc/sys/kernel/threads-max")
	if err != nil {
		return 0, nil
	}
	maxThreads, err := strconv.Atoi(strings.Trim(string(rawData),"\n"))
	if err != nil {
		return 0, nil
	}
	return maxThreads, nil

}