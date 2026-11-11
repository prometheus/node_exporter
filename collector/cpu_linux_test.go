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
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/prometheus/procfs"
)

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

func TestParseCPUList(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []uint64
		hasErr   bool
	}{
		{
			name:     "empty",
			input:    "",
			expected: []uint64{},
		},
		{
			name:     "newline only",
			input:    "\n",
			expected: []uint64{},
		},
		{
			name:     "null",
			input:    "(null)",
			expected: []uint64{},
		},
		{
			name:     "whitespace padded null",
			input:    "      (null)",
			expected: []uint64{},
		},
		{
			name:     "null with newline",
			input:    "(null)\n",
			expected: []uint64{},
		},
		{
			name:     "whitespace padded null with newline",
			input:    "      (null)\n",
			expected: []uint64{},
		},
		{
			name:     "kernel formatted padded null",
			input:    "                                                                                                                                                                                          (null)\n",
			expected: []uint64{},
		},
		{
			name:     "single cpu",
			input:    "0",
			expected: []uint64{0},
		},
		{
			name:     "single cpu with newline",
			input:    "5\n",
			expected: []uint64{5},
		},
		{
			name:     "multiple cpus and ranges",
			input:    "1,3-5,9\n",
			expected: []uint64{1, 3, 4, 5, 9},
		},
		{
			name:     "range with multiple items",
			input:    "0-2,4,6-7",
			expected: []uint64{0, 1, 2, 4, 6, 7},
		},
		{
			name:   "invalid range backwards",
			input:  "5-2",
			hasErr: true,
		},
		{
			name:   "invalid string",
			input:  "foo",
			hasErr: true,
		},
		{
			name:   "invalid range format",
			input:  "1-2-3",
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseCPUList(tc.input)
			if tc.hasErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("parseCPUList(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestGetIsolatedCPUs(t *testing.T) {
	tempDir := t.TempDir()
	cpuDir := filepath.Join(tempDir, "devices", "system", "cpu")
	if err := os.MkdirAll(cpuDir, 0755); err != nil {
		t.Fatal(err)
	}

	isolatedPath := filepath.Join(cpuDir, "isolated")

	originalSysPath := *sysPath
	*sysPath = tempDir
	defer func() {
		*sysPath = originalSysPath
	}()

	// Test (null) padded with whitespace
	if err := os.WriteFile(isolatedPath, []byte("      (null)\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cpus, err := getIsolatedCPUs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(cpus, []uint64{}) {
		t.Fatalf("getIsolatedCPUs() = %v, want empty slice", cpus)
	}

	// Test valid cpu range
	if err := os.WriteFile(isolatedPath, []byte("1,3-5,9\n"), 0644); err != nil {
		t.Fatal(err)
	}
	cpus, err = getIsolatedCPUs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(cpus, []uint64{1, 3, 4, 5, 9}) {
		t.Fatalf("getIsolatedCPUs() = %v, want [1, 3, 4, 5, 9]", cpus)
	}

	// Test non-existent file
	if err := os.Remove(isolatedPath); err != nil {
		t.Fatal(err)
	}
	_, err = getIsolatedCPUs()
	if err == nil || !os.IsNotExist(err) {
		t.Fatalf("expected not exist error, got %v", err)
	}
}
