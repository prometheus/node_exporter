// Copyright 2016 The Prometheus Authors
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

//go:build !novmstat

package collector

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestVmstatValueType(t *testing.T) {
	tests := []struct {
		name string
		want prometheus.ValueType
	}{
		{"nr_free_pages", prometheus.GaugeValue},
		{"nr_dirty", prometheus.GaugeValue},
		{"nr_writeback", prometheus.GaugeValue},
		{"oom_kill", prometheus.CounterValue},
		{"pgpgin", prometheus.CounterValue},
		{"pswpin", prometheus.CounterValue},
		{"pgfault", prometheus.CounterValue},
		{"pgmajfault", prometheus.CounterValue},
		{"nr_dirtied", prometheus.CounterValue},
		{"numa_hit", prometheus.UntypedValue},
		{"workingset_refault", prometheus.UntypedValue},
	}
	for _, tt := range tests {
		if got := vmstatValueType(tt.name); got != tt.want {
			t.Fatalf("%s: got %v want %v", tt.name, got, tt.want)
		}
	}
}
