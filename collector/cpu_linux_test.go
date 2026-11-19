// Copyright 2021 The Prometheus Authors
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

//go:build !nocpu

package collector

import (
	"io"
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/prometheus/procfs"
)

type testCPUCollector struct {
	mc Collector
}

func (c testCPUCollector) Collect(ch chan<- prometheus.Metric) {
	c.mc.Update(ch)
}

func (c testCPUCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func copyStats(d, s map[int64]procfs.CPUStat) {
	for k := range s {
		v := s[k]
		d[k] = v
	}
}

func makeTestCPUCollector(s map[int64]procfs.CPUStat) *cpuCollector {
	dup := make(map[int64]procfs.CPUStat, len(s))
	copyStats(dup, s)
	return &cpuCollector{
		logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
		cpuStats: dup,
	}
}

func TestCPU(t *testing.T) {
	firstCPUStat := map[int64]procfs.CPUStat{
		0: {
			User:      100.0,
			Nice:      100.0,
			System:    100.0,
			Idle:      100.0,
			Iowait:    100.0,
			IRQ:       100.0,
			SoftIRQ:   100.0,
			Steal:     100.0,
			Guest:     100.0,
			GuestNice: 100.0,
		},
	}

	c := makeTestCPUCollector(firstCPUStat)
	want := map[int64]procfs.CPUStat{
		0: {
			User:      101.0,
			Nice:      101.0,
			System:    101.0,
			Idle:      101.0,
			Iowait:    101.0,
			IRQ:       101.0,
			SoftIRQ:   101.0,
			Steal:     101.0,
			Guest:     101.0,
			GuestNice: 101.0,
		},
	}
	c.updateCPUStats(want)
	got := c.cpuStats
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %v CPU Stat: got %v", want, got)
	}

	c = makeTestCPUCollector(firstCPUStat)
	jumpBack := map[int64]procfs.CPUStat{
		0: {
			User:      99.9,
			Nice:      99.9,
			System:    99.9,
			Idle:      99.9,
			Iowait:    99.9,
			IRQ:       99.9,
			SoftIRQ:   99.9,
			Steal:     99.9,
			Guest:     99.9,
			GuestNice: 99.9,
		},
	}
	c.updateCPUStats(jumpBack)
	got = c.cpuStats
	if reflect.DeepEqual(jumpBack, got) {
		t.Fatalf("should have %v CPU Stat: got %v", firstCPUStat, got)
	}

	c = makeTestCPUCollector(firstCPUStat)
	resetIdle := map[int64]procfs.CPUStat{
		0: {
			User:      102.0,
			Nice:      102.0,
			System:    102.0,
			Idle:      1.0,
			Iowait:    102.0,
			IRQ:       102.0,
			SoftIRQ:   102.0,
			Steal:     102.0,
			Guest:     102.0,
			GuestNice: 102.0,
		},
	}
	c.updateCPUStats(resetIdle)
	got = c.cpuStats
	if !reflect.DeepEqual(resetIdle, got) {
		t.Fatalf("should have %v CPU Stat: got %v", resetIdle, got)
	}
}

func TestCPUOffline(t *testing.T) {
	// CPU 1 goes offline.
	firstCPUStat := map[int64]procfs.CPUStat{
		0: {
			User:      100.0,
			Nice:      100.0,
			System:    100.0,
			Idle:      100.0,
			Iowait:    100.0,
			IRQ:       100.0,
			SoftIRQ:   100.0,
			Steal:     100.0,
			Guest:     100.0,
			GuestNice: 100.0,
		},
		1: {
			User:      101.0,
			Nice:      101.0,
			System:    101.0,
			Idle:      101.0,
			Iowait:    101.0,
			IRQ:       101.0,
			SoftIRQ:   101.0,
			Steal:     101.0,
			Guest:     101.0,
			GuestNice: 101.0,
		},
	}

	c := makeTestCPUCollector(firstCPUStat)
	want := map[int64]procfs.CPUStat{
		0: {
			User:      100.0,
			Nice:      100.0,
			System:    100.0,
			Idle:      100.0,
			Iowait:    100.0,
			IRQ:       100.0,
			SoftIRQ:   100.0,
			Steal:     100.0,
			Guest:     100.0,
			GuestNice: 100.0,
		},
	}
	c.updateCPUStats(want)
	got := c.cpuStats
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %v CPU Stat: got %v", want, got)
	}

	// CPU 1 comes back online.
	want = map[int64]procfs.CPUStat{
		0: {
			User:      100.0,
			Nice:      100.0,
			System:    100.0,
			Idle:      100.0,
			Iowait:    100.0,
			IRQ:       100.0,
			SoftIRQ:   100.0,
			Steal:     100.0,
			Guest:     100.0,
			GuestNice: 100.0,
		},
		1: {
			User:      101.0,
			Nice:      101.0,
			System:    101.0,
			Idle:      101.0,
			Iowait:    101.0,
			IRQ:       101.0,
			SoftIRQ:   101.0,
			Steal:     101.0,
			Guest:     101.0,
			GuestNice: 101.0,
		},
	}
	c.updateCPUStats(want)
	got = c.cpuStats
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("should have %v CPU Stat: got %v", want, got)
	}
}

func TestCPUNUMAInfo(t *testing.T) {
	origProcPath, origSysPath := *procPath, *sysPath
	t.Cleanup(func() { *procPath, *sysPath = origProcPath, origSysPath })
	*procPath = "fixtures/proc"
	*sysPath = "fixtures/sys"

	mc, err := NewCPUCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	// node2 in the fixtures is a memory-only node and has no CPUs.
	expected := `# HELP node_cpu_numa_info NUMA node of each online CPU, information from /sys/devices/system/node/node*/cpulist.
# TYPE node_cpu_numa_info gauge
node_cpu_numa_info{cpu="0",node="0"} 1
node_cpu_numa_info{cpu="1",node="0"} 1
node_cpu_numa_info{cpu="2",node="1"} 1
node_cpu_numa_info{cpu="3",node="1"} 1
`
	if err := testutil.CollectAndCompare(testCPUCollector{mc}, strings.NewReader(expected), "node_cpu_numa_info"); err != nil {
		t.Fatal(err)
	}
}
