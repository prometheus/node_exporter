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

// +build !noprocesses

package collector

import (
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"

	"github.com/prometheus/client_golang/prometheus"
)

func canTestPerf(t *testing.T) {
	paranoidBytes, err := ioutil.ReadFile("/proc/sys/kernel/perf_event_paranoid")
	if err != nil {
		t.Skip("Procfs not mounted, skipping perf tests")
	}
	paranoidStr := strings.Replace(string(paranoidBytes), "\n", "", -1)
	paranoid, err := strconv.Atoi(paranoidStr)
	if err != nil {
		t.Fatalf("Expected perf_event_paranoid to be an int, got: %s", paranoidStr)
	}
	if paranoid >= 1 {
		t.Skip("Skipping perf tests, set perf_event_paranoid to 0")
	}
}

func TestPerfCollector(t *testing.T) {
	canTestPerf(t)
	collector, err := NewPerfCollector(log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}

	// Setup background goroutine to capture metrics.
	metrics := make(chan prometheus.Metric)
	defer close(metrics)
	go func() {
		for range metrics {
		}
	}()
	if err := collector.Update(metrics); err != nil {
		t.Fatal(err)
	}
}

func TestPerfCollectorStride(t *testing.T) {
	canTestPerf(t)

	tests := []struct {
		name   string
		flag   string
		exCPUs []int
	}{
		{
			name:   "valid single CPU",
			flag:   "1",
			exCPUs: []int{1},
		},
		{
			name:   "valid range CPUs",
			flag:   "1-5",
			exCPUs: []int{1, 2, 3, 4, 5},
		},
		{
			name:   "valid stride",
			flag:   "1-8:2",
			exCPUs: []int{1, 3, 5, 7},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ncpu := runtime.NumCPU()
			for _, cpu := range test.exCPUs {
				if cpu > ncpu {
					t.Skipf("Skipping test because runtime.NumCPU < %d", cpu)
				}
			}
			perfCPUsFlag = &test.flag
			collector, err := NewPerfCollector(log.NewNopLogger())
			if err != nil {
				t.Fatal(err)
			}

			c := collector.(*perfCollector)
			for _, cpu := range test.exCPUs {
				if _, ok := c.perfHwProfilers[cpu]; !ok {
					t.Fatalf("Expected CPU %v in hardware profilers", cpu)
				}
				if _, ok := c.perfSwProfilers[cpu]; !ok {
					t.Fatalf("Expected CPU %v in software profilers", cpu)
				}
				if _, ok := c.perfCacheProfilers[cpu]; !ok {
					t.Fatalf("Expected CPU %v in cache profilers", cpu)
				}
			}
		})
	}
}

func TestPerfCPUFlagToCPUs(t *testing.T) {
	tests := []struct {
		name   string
		flag   string
		exCpus []int
		errStr string
	}{
		{
			name:   "valid single CPU",
			flag:   "1",
			exCpus: []int{1},
		},
		{
			name:   "valid range CPUs",
			flag:   "1-5",
			exCpus: []int{1, 2, 3, 4, 5},
		},
		{
			name:   "valid double digit",
			flag:   "10",
			exCpus: []int{10},
		},
		{
			name:   "valid double digit range",
			flag:   "10-12",
			exCpus: []int{10, 11, 12},
		},
		{
			name:   "valid double digit stride",
			flag:   "10-20:5",
			exCpus: []int{10, 15, 20},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpus, err := perfCPUFlagToCPUs(test.flag)
			if test.errStr != "" {
				if err != nil {
					t.Fatal("expected error to not be nil")
				}
				if test.errStr != err.Error() {
					t.Fatalf(
						"expected error %q, got %q",
						test.errStr,
						err.Error(),
					)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(cpus) != len(test.exCpus) {
				t.Fatalf(
					"expected CPUs %v, got %v",
					test.exCpus,
					cpus,
				)
			}
			for i := range cpus {
				if test.exCpus[i] != cpus[i] {
					t.Fatalf(
						"expected CPUs %v, got %v",
						test.exCpus[i],
						cpus[i],
					)
				}
			}
		})
	}
}

func TestPerfTracepointFlagToTracepoints(t *testing.T) {
	tests := []struct {
		name          string
		flag          []string
		exTracepoints []*perfTracepoint
		errStr        string
	}{
		{
			name: "valid single tracepoint",
			flag: []string{"sched:sched_kthread_stop"},
			exTracepoints: []*perfTracepoint{
				{
					subsystem: "sched",
					event:     "sched_kthread_stop",
				},
			},
		},
		{
			name: "valid multiple tracepoints",
			flag: []string{"sched:sched_kthread_stop", "sched:sched_process_fork"},
			exTracepoints: []*perfTracepoint{
				{
					subsystem: "sched",
					event:     "sched_kthread_stop",
				},
				{
					subsystem: "sched",
					event:     "sched_process_fork",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tracepoints, err := perfTracepointFlagToTracepoints(test.flag)
			if test.errStr != "" {
				if err != nil {
					t.Fatal("expected error to not be nil")
				}
				if test.errStr != err.Error() {
					t.Fatalf(
						"expected error %q, got %q",
						test.errStr,
						err.Error(),
					)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			for i := range tracepoints {
				if test.exTracepoints[i].event != tracepoints[i].event &&
					test.exTracepoints[i].subsystem != tracepoints[i].subsystem {
					t.Fatalf(
						"expected tracepoint %v, got %v",
						test.exTracepoints[i],
						tracepoints[i],
					)
				}
			}
		})
	}
}
