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
	procsState   *prometheus.Desc
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
		procsState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "procs_state"),
			"Number of processes in each state.",
			[]string{"state"}, nil,
		),
	}, nil
}
func (t *threadsCollector) Update(ch chan<- prometheus.Metric) error {
	states ,threads, err := getAllocatedThreads()
	if err != nil {
		return fmt.Errorf("Unable to retrieve number of allocated threads %v\n", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadAlloc, prometheus.GaugeValue, float64(threads))
	maxThreads, err := getMaxThreads()
	if err != nil {
		return fmt.Errorf("Unable to retrieve limit number of threads %v\n", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadMax, prometheus.GaugeValue, float64(maxThreads))
	for state := range states {
		ch <- prometheus.MustNewConstMetric(t.procsState, prometheus.GaugeValue, float64(states[state]), state)
	}
	return nil
}

func getAllocatedThreads() (map[string]int32,int, error) {
	p, err := procfs.AllProcs()
	if err != nil {
		return nil, 0, err
	}
	thread := 0
	procStates := make(map[string]int32)
	for _, pid := range p {
		stat, err := pid.NewStat()
		if err != nil {
			return nil, 0, err
		}
		procStates[stat.State] += 1
		thread += stat.NumThreads

	}
	return procStates,thread, nil
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