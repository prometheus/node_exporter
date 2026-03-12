// Copyright The Prometheus Authors
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

//go:build !nonvmesubsystem

package collector

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestNVMeSubsystemMetrics(t *testing.T) {
	*sysPath = "fixtures/sys"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewNVMeSubsystemCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	c := coll.(*nvmeSubsystemCollector)

	ch := make(chan prometheus.Metric, 200)
	if err := c.Update(ch); err != nil {
		t.Fatal(err)
	}
	close(ch)

	metrics := make(map[string][]*dto.Metric)
	for m := range ch {
		d := &dto.Metric{}
		if err := m.Write(d); err != nil {
			t.Fatal(err)
		}
		desc := m.Desc().String()
		metrics[desc] = append(metrics[desc], d)
	}

	assertGaugeValue(t, metrics, `paths"`, labelMap{"subsystem": "nvme-subsys0"}, 4)
	assertGaugeValue(t, metrics, "paths_live", labelMap{"subsystem": "nvme-subsys0"}, 3)

	assertGaugeValue(t, metrics, "path_state",
		labelMap{"subsystem": "nvme-subsys0", "controller": "nvme0", "transport": "fc", "state": "live"}, 1)
	assertGaugeValue(t, metrics, "path_state",
		labelMap{"subsystem": "nvme-subsys0", "controller": "nvme0", "transport": "fc", "state": "dead"}, 0)
	assertGaugeValue(t, metrics, "path_state",
		labelMap{"subsystem": "nvme-subsys0", "controller": "nvme3", "transport": "fc", "state": "dead"}, 1)
	assertGaugeValue(t, metrics, "path_state",
		labelMap{"subsystem": "nvme-subsys0", "controller": "nvme3", "transport": "fc", "state": "live"}, 0)
}

func TestNVMeSubsystemNoDevices(t *testing.T) {
	*sysPath = t.TempDir()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewNVMeSubsystemCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	c := coll.(*nvmeSubsystemCollector)

	ch := make(chan prometheus.Metric, 200)
	err = c.Update(ch)
	close(ch)

	if err != ErrNoData {
		t.Fatalf("expected ErrNoData, got %v", err)
	}
}

func TestNormalizeControllerState(t *testing.T) {
	tests := []struct {
		raw      string
		expected string
	}{
		{"live", "live"},
		{"connecting", "connecting"},
		{"resetting", "resetting"},
		{"dead", "dead"},
		{"deleting", "deleting"},
		{"deleting (no IO)", "deleting (no IO)"},
		{"new", "new"},
		{"", "unknown"},
		{"something-else", "unknown"},
	}
	for _, tc := range tests {
		got := normalizeControllerState(tc.raw)
		if got != tc.expected {
			t.Errorf("normalizeControllerState(%q) = %q, want %q", tc.raw, got, tc.expected)
		}
	}
}

type labelMap map[string]string

func assertGaugeValue(t *testing.T, metrics map[string][]*dto.Metric, metricSubstring string, labels labelMap, expected float64) {
	t.Helper()
	for desc, ms := range metrics {
		if !strings.Contains(desc, metricSubstring) {
			continue
		}
		for _, m := range ms {
			if matchLabels(m.GetLabel(), labels) {
				got := m.GetGauge().GetValue()
				if got != expected {
					t.Errorf("%s%v: got %v, want %v", metricSubstring, labels, got, expected)
				}
				return
			}
		}
	}
	t.Errorf("metric %s%v not found", metricSubstring, labels)
}

func matchLabels(pairs []*dto.LabelPair, want labelMap) bool {
	if want == nil {
		return len(pairs) == 0
	}
	found := 0
	for _, lp := range pairs {
		if v, ok := want[lp.GetName()]; ok && v == lp.GetValue() {
			found++
		}
	}
	return found == len(want)
}
