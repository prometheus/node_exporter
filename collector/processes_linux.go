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

	"github.com/peterloeffler/procfs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	mof_threshold = kingpin.Flag("collector.processes.mof_threshold", "Threshold for max open files in %.").Default("0").String()
	mofLabelNames = []string{"pid", "process_name", "max_open_files", "cur_open_files"}
)

type openFiles struct {
	pid          string
	name         string
	maxOpenFiles string
	curOpenFiles string
	percent      float64
}

type processCollector struct {
	threadAlloc *prometheus.Desc
	threadLimit *prometheus.Desc
	procsState  *prometheus.Desc
	pidUsed     *prometheus.Desc
	pidMax      *prometheus.Desc
	fhMOFPct    *prometheus.Desc
}

func init() {
	registerCollector("processes", defaultDisabled, NewProcessStatCollector)
}

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
		fhMOFPct: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_open_files_percentage"),
			"Process that has reached more than "+string(*mof_threshold)+"% (tunable via --collector.processes.mof_threshold) of max open files according to it's hard limit.", mofLabelNames, nil,
		),
	}, nil
}
func (t *processCollector) Update(ch chan<- prometheus.Metric) error {
	pids, states, threads, ofr, err := getAllocatedThreads()
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

	pidM, err := readUintFromFile(procFilePath("sys/kernel/pid_max"))
	if err != nil {
		return fmt.Errorf("unable to retrieve limit number of maximum pids alloved: %q", err)
	}
	ch <- prometheus.MustNewConstMetric(t.pidUsed, prometheus.GaugeValue, float64(pids))
	ch <- prometheus.MustNewConstMetric(t.pidMax, prometheus.GaugeValue, float64(pidM))

	for _, p := range ofr {
		ch <- prometheus.MustNewConstMetric(t.fhMOFPct, prometheus.GaugeValue, p.percent, p.pid, p.name, p.maxOpenFiles, p.curOpenFiles)
	}

	return nil
}

func getAllocatedThreads() (int, map[string]int32, int, []openFiles, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return 0, nil, 0, []openFiles{}, err
	}
	p, err := fs.AllProcs()
	if err != nil {
		return 0, nil, 0, []openFiles{}, err
	}
	pids := 0
	thread := 0
	procStates := make(map[string]int32)
	openFilesReached := []openFiles{}
	for _, pid := range p {
		stat, err := pid.NewStat()
		// PIDs can vanish between getting the list and getting stats.
		if os.IsNotExist(err) {
			log.Debugf("file not found when retrieving stats: %q", err)
			continue
		}
		if err != nil {
			return 0, nil, 0, []openFiles{}, err
		}
		pids += 1
		procStates[stat.State] += 1
		thread += stat.NumThreads
		percentage := (float64(stat.CurrentOpenFiles) * 100) / float64(stat.MaxOpenFiles)
		if mof_threshold, err := strconv.ParseFloat(*mof_threshold, 64); err == nil {
			if percentage >= mof_threshold {
				// provide metric for processes that have more open files than mof_threshold% of it's max open files
				openFilesReached = append(openFilesReached, openFiles{
					pid:          strconv.Itoa(stat.PID),
					name:         stat.Comm,
					maxOpenFiles: strconv.Itoa(stat.MaxOpenFiles),
					curOpenFiles: strconv.Itoa(stat.CurrentOpenFiles),
					percent:      percentage,
				})
			}
		}
		if err != nil {
			return 0, nil, 0, []openFiles{}, err
		}
	}
	return pids, procStates, thread, openFilesReached, nil
}
