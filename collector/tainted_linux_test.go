// Copyright 2024 The Prometheus Authors
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

//go:build !notainted

package collector

import (
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestTaintedCollector(t *testing.T) {
	*procPath = "fixtures/proc"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewTaintedCollector(logger)
	if err != nil {
		t.Fatalf("failed to create tainted collector: %v", err)
	}

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(&taintedCollectorWrapper{c.(*taintedCollector)})

	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather failed: %v", err)
	}
	if len(mfs) != 1 {
		t.Fatalf("expected 1 metric family, got %d", len(mfs))
	}

	mf := mfs[0]
	if got := mf.GetName(); got != "node_kernel_tainted" {
		t.Errorf("metric name: want node_kernel_tainted, got %s", got)
	}

	// Expect one series per known taint bit (20 defined by the kernel).
	const wantBits = 20
	if got := len(mf.GetMetric()); got != wantBits {
		t.Errorf("metric count: want %d, got %d", wantBits, got)
	}

	// Build bit → value map for assertion.
	// Fixture is 12288 = bit 12 (O) + bit 13 (E).
	// Build flag → value map for assertion (labels: bit, flag).
	flagVals := make(map[string]float64)
	for _, m := range mf.GetMetric() {
		// Each metric has exactly 2 labels: bit and flag.
		for _, lp := range m.GetLabel() {
			if lp.GetName() == "flag" {
				flagVals[lp.GetValue()] = m.GetGauge().GetValue()
			}
		}
	}

	// Fixture is 12288 = bit 12 (O) + bit 13 (E).
	for _, tc := range []struct {
		flag string
		want float64
	}{
		{"O", 1}, // Externally-built (out-of-tree) module — set
		{"E", 1}, // Unsigned module — set
		{"L", 0}, // Soft lockup — must be clear
		{"P", 0},
		{"T", 0},
	} {
		got, ok := flagVals[tc.flag]
		if !ok {
			t.Errorf("flag %q not found in metrics", tc.flag)
			continue
		}
		if got != tc.want {
			t.Errorf("flag %q: want %.0f, got %.0f", tc.flag, tc.want, got)
		}
	}
}

// taintedCollectorWrapper adapts taintedCollector to prometheus.Collector.
type taintedCollectorWrapper struct {
	c *taintedCollector
}

func (w *taintedCollectorWrapper) Describe(ch chan<- *prometheus.Desc) {
	ch <- w.c.desc
}

func (w *taintedCollectorWrapper) Collect(ch chan<- prometheus.Metric) {
	_ = w.c.Update(ch)
}
