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

// +build !nostat

package collector

import (
	"fmt"

	"github.com/prometheus/procfs"

	"github.com/prometheus/client_golang/prometheus"

	gops "github.com/mitchellh/go-ps"
)

var (
	// R  Running
	// S  Sleeping in an interruptible wait
	// D  Waiting in uninterruptible disk sleep
	// Z  Zombie
	// T  Stopped (on a signal) or (before Linux 2.6.33) trace stopped
	// t  Tracing stop (Linux 2.6.33 onward)
	// W  Paging (only before Linux 2.6.0)
	// X  Dead (from Linux 2.6.0 onward)
	// x  Dead (Linux 2.6.33 to 3.13 only)
	// K  Wakekill (Linux 2.6.33 to 3.13 only)
	// W  Waking (Linux 2.6.33 to 3.13 only)
	knownStates = [...]string{"R", "S", "D", "Z", "T", "t", "W", "X", "x", "K", "W", "P"}
)

type statCollector struct {
	cpu          *prometheus.Desc
	intr         *prometheus.Desc
	ctxt         *prometheus.Desc
	forks        *prometheus.Desc
	btime        *prometheus.Desc
	procsRunning *prometheus.Desc
	procsBlocked *prometheus.Desc
	procsState   *prometheus.Desc
}

func init() {
	registerCollector("stat", defaultEnabled, NewStatCollector)
}

// NewStatCollector returns a new Collector exposing kernel/system statistics.
func NewStatCollector() (Collector, error) {
	return &statCollector{
		cpu: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "cpu"),
			"Seconds the cpus spent in each mode.",
			[]string{"cpu", "mode"}, nil,
		),
		intr: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "intr"),
			"Total number of interrupts serviced.",
			nil, nil,
		),
		ctxt: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "context_switches"),
			"Total number of context switches.",
			nil, nil,
		),
		forks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "forks"),
			"Total number of forks.",
			nil, nil,
		),
		btime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "boot_time"),
			"Node boot time, in unixtime.",
			nil, nil,
		),
		procsRunning: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "procs_running"),
			"Number of processes in runnable state.",
			nil, nil,
		),
		procsBlocked: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "procs_blocked"),
			"Number of processes blocked waiting for I/O to complete.",
			nil, nil,
		),
		procsState: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "procs_state"),
			"Number of processes in each state.",
			[]string{"state"}, nil,
		),
	}, nil
}

// Update implements Collector and exposes kernel and system statistics.
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return fmt.Errorf("failed to open procfs: %v", err)
	}
	stats, err := fs.NewStat()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.intr, prometheus.CounterValue, float64(stats.IRQTotal))
	ch <- prometheus.MustNewConstMetric(c.ctxt, prometheus.CounterValue, float64(stats.ContextSwitches))
	ch <- prometheus.MustNewConstMetric(c.forks, prometheus.CounterValue, float64(stats.ProcessCreated))

	ch <- prometheus.MustNewConstMetric(c.btime, prometheus.GaugeValue, float64(stats.BootTime))

	ch <- prometheus.MustNewConstMetric(c.procsRunning, prometheus.GaugeValue, float64(stats.ProcessesRunning))
	ch <- prometheus.MustNewConstMetric(c.procsBlocked, prometheus.GaugeValue, float64(stats.ProcessesBlocked))

	processes, err := gops.Processes()
	if err != nil {
		return err
	}
	processStates := make(map[string]uint64)
	for _, process := range processes {
		proc, err := procfs.NewProc(process.Pid())
		if err != nil {
			continue
		}
		procStat, err := proc.NewStat()
		if err != nil {
			return err
		}
		processStates[procStat.State]++
	}

	for _, state := range knownStates {
		ch <- prometheus.MustNewConstMetric(c.procsState, prometheus.GaugeValue, float64(processStates[state]), state)
	}

	return nil
}
