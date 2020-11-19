// Copyright 2020 The Prometheus Authors
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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/process"
)

const (
	subsystem = "processes"
)

var (
	// Remake labels
	r = strings.NewReplacer(".exe", "", "-", "_", ".", "_")
)

type processCollector struct {
	procfs  []*process.Process
	pidUsed *prometheus.Desc
	logger  log.Logger
}

func init() {
	registerCollector("processes", defaultDisabled, NewProcessStatCollector)
}

// NewProcessStatCollector returns a new Collector exposing process data read from the proc filesystem.
func NewProcessStatCollector(logger log.Logger) (Collector, error) {
	pfs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get procfs: %w", err)
	}
	return &processCollector{
		procfs: pfs,
		pidUsed: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "pids"),
			"Number of PIDs", nil, nil,
		),
		logger: logger,
	}, nil
}
func (c *processCollector) Update(ch chan<- prometheus.Metric) error {
	// handle error is not required
	// We need update process
	c.procfs, _ = process.Processes()

	// Update number pids
	ch <- prometheus.MustNewConstMetric(c.pidUsed, prometheus.GaugeValue, float64(len(c.procfs)))

	procStats, err := c.getAllocatedProcesses()
	if err != nil {
		return fmt.Errorf("unable to retrieve number of allocated processes: %q", err)
	}

	for procName, procInfo := range procStats {
		// Update memory information
		for k, v := range procInfo["mem"] {
			var key = "mem_bytes"
			if k == "used_percent" {
				key = "mem_percent"
			}

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, key),
					"Memory information field kinds of processes.",
					[]string{"process_name", "parameter_name"}, nil,
				),
				prometheus.GaugeValue, v, procName, k,
			)
		}

		// Update iocounters information
		for k, v := range procInfo["iocounters"] {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "iocounters_bytes"),
					fmt.Sprintf("IOCounters information field kinds of processes."),
					[]string{"process_name", "parameter_name"}, nil,
				),
				prometheus.GaugeValue, v, procName, k,
			)
		}
	}
	return nil
}

func (c *processCollector) getAllocatedProcesses() (map[string]map[string]map[string]float64, error) {
	gaProcesses := map[string]map[string]map[string]float64{}

	for _, process := range c.procfs {

		// Ignore root pid
		if process.Pid == 0 {
			continue
		}

		// Ignore process can't get name
		pName, err := process.Name()
		if err != nil {
			continue
		}

		// remove post.fix .exe before save to map
		// get memory info that process using
		kName := r.Replace(pName)

		// Make a map store data
		memInfo, err := getMemoryInfo(process)
		if err != nil {
			continue
		}

		iocInfo, err := getIOCountersInfo(process)
		if err != nil {
			continue
		}

		// Because processes have multiple threads
		// So we must counter with sum
		if _, ok := gaProcesses[kName]; !ok {
			// Make a new map store data
			gaProcesses[kName] = map[string]map[string]float64{}

			gaProcesses[kName]["mem"] = memInfo
			gaProcesses[kName]["iocounters"] = iocInfo
		} else {
			for k, v := range memInfo {
				gaProcesses[kName]["mem"][k] = gaProcesses[kName]["mem"][k] + v
			}

			for k, v := range iocInfo {
				gaProcesses[kName]["iocounters"][k] = gaProcesses[kName]["iocounters"][k] + v
			}
		}

	}

	return gaProcesses, nil
}

func getMemoryInfo(proc *process.Process) (map[string]float64, error) {
	gMI := map[string]float64{}
	memInfo, err := proc.MemoryInfo()
	if err != nil {
		return nil, err
	}

	// parse data
	memBytes, err := json.Marshal(memInfo)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(memBytes, &gMI)

	// get memory percent
	memPerc, err := proc.MemoryPercent()
	if err != nil {
		return nil, err
	}
	gMI["used_percent"] = float64(memPerc)
	return gMI, nil

}

func getIOCountersInfo(proc *process.Process) (map[string]float64, error) {
	gIOC := map[string]float64{}
	iocInfo, err := proc.IOCounters()
	if err != nil {
		return nil, err
	}

	// parse data
	iocBytes, err := json.Marshal(iocInfo)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(iocBytes, &gIOC)

	return gIOC, nil
}
