// Copyright 2018 The Prometheus Authors
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

// +build !noprocesses

package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

type processCollector struct {
	threadAlloc *prometheus.Desc
	threadLimit *prometheus.Desc
	procsState  *prometheus.Desc
	pidUsed     *prometheus.Desc
	pidMax      *prometheus.Desc
}

func init() {
	registerCollector("processes", defaultDisabled, NewProcessStatCollector)
}

func NewProcessStatCollector() (Collector, error) {
	subsystem := "processes"
	return &processCollector{
		threadAlloc: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "threads"),
			"Allocated threads in system",
			nil, nil,
		),
		threadLimit: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "max_threads"),
			"Limit of threads in the system",
			nil, nil,
		),
		procsState: PrometheusNewDesc(
			prometheus.BuildFQName(namespace, subsystem, "state"),
			"Number of processes in each state.",
			[]string{"state"}, nil,
		),
		pidUsed: PrometheusNewDesc(prometheus.BuildFQName(namespace, subsystem, "pids"),
			"Number of PIDs", nil, nil,
		),
		pidMax: PrometheusNewDesc(prometheus.BuildFQName(namespace, subsystem, "max_processes"),
			"Number of max PIDs limit", nil, nil,
		),
	}, nil
}
func (t *processCollector) Update(ch chan<- prometheus.Metric) error {
	pids, states, threads, err := getAllocatedThreads()
	if err != nil {
		return fmt.Errorf("Unable to retrieve number of allocated threads %v\n", err)
	}

	ch <- prometheus.MustNewConstMetric(t.threadAlloc, prometheus.GaugeValue, float64(threads))
	maxThreads, err := readUintFromFile(procFilePath("sys/kernel/threads-max"))
	if err != nil {
		return fmt.Errorf("Unable to retrieve limit number of threads %v\n", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadLimit, prometheus.GaugeValue, float64(maxThreads))

	for state := range states {
		ch <- prometheus.MustNewConstMetric(t.procsState, prometheus.GaugeValue, float64(states[state]), state)
	}

	pidM, err := readUintFromFile(procFilePath("sys/kernel/pid_max"))
	if err != nil {
		return fmt.Errorf("Unable to retrieve limit number of maximum pids alloved %v\n", err)
	}
	ch <- prometheus.MustNewConstMetric(t.pidUsed, prometheus.GaugeValue, float64(pids))
	ch <- prometheus.MustNewConstMetric(t.pidMax, prometheus.GaugeValue, float64(pidM))

	return nil
}

func getAllocatedThreads() (int, map[string]int32, int, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return 0, nil, 0, err
	}
	p, err := fs.AllProcs()
	if err != nil {
		return 0, nil, 0, err
	}
	thread := 0
	procStates := make(map[string]int32)
	for _, pid := range p {
		stat, err := pid.NewStat()
		if err != nil {
			return 0, nil, 0, err
		}
		procStates[stat.State] += 1
		thread += stat.NumThreads
	}
	return len(p), procStates, thread, nil
}
