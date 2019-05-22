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
	"strconv"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestPerfCollector(t *testing.T) {
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
	collector, err := NewPerfCollector()
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
