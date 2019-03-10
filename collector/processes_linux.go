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
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

type processCollector struct {
	threadAlloc  *prometheus.Desc
	threadLimit  *prometheus.Desc
	procsState   *prometheus.Desc
	pidUsed      *prometheus.Desc
	pidMax       *prometheus.Desc
	perProcUsage *prometheus.Desc
}

type procResUsage struct {
	rss     int
	vsize   uint
	cpuTime float64
}

func init() {
	registerCollector("processes", defaultDisabled, NewProcessStatCollector)
}

// NewProcessStatCollector returns a new Collector exposing process data read from the proc filesystem.
func NewProcessStatCollector() (Collector, error) {
	subsystem := "processes"
	return &processCollector{
		threadAlloc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "threads"),
			"Allocated threads in system",
			nil, nil,
		),
		threadLimit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "max_threads"),
			"Limit of threads in the system",
			nil, nil,
		),
		procsState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "state"),
			"Number of processes in each state.",
			[]string{"state"}, nil,
		),
		pidUsed: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "pids"),
			"Number of PIDs", nil, nil,
		),
		pidMax: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_processes"),
			"Number of max PIDs limit", nil, nil,
		),
		perProcUsage: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "usage"),
			"Per process usage of system resources",
			[]string{"pid", "resource"}, nil,
		),
	}, nil
}
func (t *processCollector) Update(ch chan<- prometheus.Metric) error {
	pids, states, threads, pidToUsage, err := getAllocatedThreads()
	if err != nil {
		return fmt.Errorf("unable to retrieve number of allocated threads: %q", err)
	}

	ch <- prometheus.MustNewConstMetric(t.threadAlloc, prometheus.GaugeValue, float64(threads))
	maxThreads, err := readUintFromFile(procFilePath("sys/kernel/threads-max"))
	if err != nil {
		return fmt.Errorf("unable to retrieve limit number of threads: %q", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadLimit, prometheus.GaugeValue, float64(maxThreads))

	for state := range states {
		ch <- prometheus.MustNewConstMetric(t.procsState, prometheus.GaugeValue, float64(states[state]), state)
	}

	for pid, usage := range pidToUsage {
		pidStr := strconv.Itoa(pid)
		ch <- prometheus.MustNewConstMetric(t.perProcUsage, prometheus.GaugeValue, float64(usage.rss), pidStr, "rss")
		ch <- prometheus.MustNewConstMetric(t.perProcUsage, prometheus.GaugeValue, float64(usage.vsize), pidStr, "vsize")
		ch <- prometheus.MustNewConstMetric(t.perProcUsage, prometheus.GaugeValue, usage.cpuTime, pidStr, "cpu_time")
	}

	pidM, err := readUintFromFile(procFilePath("sys/kernel/pid_max"))
	if err != nil {
		return fmt.Errorf("unable to retrieve limit number of maximum pids alloved: %q", err)
	}
	ch <- prometheus.MustNewConstMetric(t.pidUsed, prometheus.GaugeValue, float64(pids))
	ch <- prometheus.MustNewConstMetric(t.pidMax, prometheus.GaugeValue, float64(pidM))

	return nil
}

func getAllocatedThreads() (int, map[string]int32, int, map[int]procResUsage, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return 0, nil, 0, nil, err
	}
	p, err := fs.AllProcs()
	if err != nil {
		return 0, nil, 0, nil, err
	}
	pids := 0
	thread := 0
	procStates := make(map[string]int32)
	pidToUsage := make(map[int]procResUsage)
	for _, pid := range p {
		stat, err := pid.NewStat()
		// PIDs can vanish between getting the list and getting stats.
		if os.IsNotExist(err) {
			log.Debugf("file not found when retrieving stats: %q", err)
			continue
		}
		if err != nil {
			return 0, nil, 0, nil, err
		}
		pids++
		procStates[stat.State]++
		pidToUsage[pid.PID] = procResUsage{rss: stat.ResidentMemory(), vsize: stat.VirtualMemory(), cpuTime: stat.CPUTime()}
		thread += stat.NumThreads
	}
	return pids, procStates, thread, pidToUsage, nil
}
