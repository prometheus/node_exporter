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

//go:build !nodmmultipath

package collector

import (
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestDMMultipathScan(t *testing.T) {
	devices, err := scanDMMultipathDevices(filepath.Join("fixtures", "sys"))
	if err != nil {
		t.Fatal(err)
	}

	if len(devices) != 2 {
		t.Fatalf("expected 2 multipath devices, got %d", len(devices))
	}

	devA := devices[0]
	if devA.Name != "mpathA" {
		t.Errorf("expected mpathA, got %s", devA.Name)
	}
	if devA.SysfsName != "dm-5" {
		t.Errorf("expected dm-5, got %s", devA.SysfsName)
	}
	if !strings.HasPrefix(devA.UUID, "mpath-") {
		t.Errorf("expected mpath- UUID prefix, got %s", devA.UUID)
	}
	if devA.Suspended {
		t.Error("expected device not suspended")
	}
	if devA.SizeBytes != 53687091200 {
		t.Errorf("expected 53687091200 bytes (50 GiB), got %d", devA.SizeBytes)
	}
	if len(devA.Paths) != 4 {
		t.Fatalf("expected 4 paths, got %d", len(devA.Paths))
	}

	runningCount := 0
	for _, p := range devA.Paths {
		if p.State == "running" {
			runningCount++
		}
	}
	if runningCount != 3 {
		t.Errorf("expected 3 running paths, got %d", runningCount)
	}

	devB := devices[1]
	if devB.Name != "mpathB" {
		t.Errorf("expected mpathB, got %s", devB.Name)
	}
	if len(devB.Paths) != 2 {
		t.Fatalf("expected 2 paths for mpathB, got %d", len(devB.Paths))
	}
}

func TestDMMultipathSkipsNonMultipath(t *testing.T) {
	devices, err := scanDMMultipathDevices(filepath.Join("fixtures", "sys"))
	if err != nil {
		t.Fatal(err)
	}

	for _, dev := range devices {
		if dev.SysfsName == "dm-7" {
			t.Error("dm-7 (LVM) should have been filtered out")
		}
	}
}

func TestDMMultipathMetrics(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewDMMultipathCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	c := coll.(*dmMultipathCollector)
	c.scanDevices = func() ([]dmMultipathDevice, error) {
		return scanDMMultipathDevices(filepath.Join("fixtures", "sys"))
	}

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

	assertGaugeValue(t, metrics, "device_active", labelMap{"device": "mpathA"}, 1)
	assertGaugeValue(t, metrics, "device_active", labelMap{"device": "mpathB"}, 1)
	assertGaugeValue(t, metrics, "device_size_bytes", labelMap{"device": "mpathA"}, 53687091200)
	assertGaugeValue(t, metrics, "device_paths_total", labelMap{"device": "mpathA"}, 4)
	assertGaugeValue(t, metrics, "device_paths_total", labelMap{"device": "mpathB"}, 2)

	// mpathA: sdi, sdj, sdk are running; sdl is offline → 3 active, 1 failed.
	assertGaugeValue(t, metrics, "device_paths_active", labelMap{"device": "mpathA"}, 3)
	assertGaugeValue(t, metrics, "device_paths_failed", labelMap{"device": "mpathA"}, 1)

	// mpathB: sdm, sdn are both running → 2 active, 0 failed.
	assertGaugeValue(t, metrics, "device_paths_active", labelMap{"device": "mpathB"}, 2)
	assertGaugeValue(t, metrics, "device_paths_failed", labelMap{"device": "mpathB"}, 0)

	assertGaugeValue(t, metrics, "path_state",
		labelMap{"device": "mpathA", "path": "sdi", "state": "running"}, 1)
	assertGaugeValue(t, metrics, "path_state",
		labelMap{"device": "mpathA", "path": "sdi", "state": "offline"}, 0)
	assertGaugeValue(t, metrics, "path_state",
		labelMap{"device": "mpathA", "path": "sdl", "state": "offline"}, 1)
	assertGaugeValue(t, metrics, "path_state",
		labelMap{"device": "mpathA", "path": "sdl", "state": "running"}, 0)
}

func TestDMMultipathNoDevices(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewDMMultipathCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	c := coll.(*dmMultipathCollector)
	c.scanDevices = func() ([]dmMultipathDevice, error) {
		return nil, nil
	}

	ch := make(chan prometheus.Metric, 200)
	if err := c.Update(ch); err != nil {
		t.Fatal(err)
	}
	close(ch)

	count := 0
	for range ch {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 metrics when no devices, got %d", count)
	}
}

func TestNormalizeDMPathState(t *testing.T) {
	tests := []struct {
		raw      string
		expected string
	}{
		{"running", "running"},
		{"offline", "offline"},
		{"blocked", "blocked"},
		{"transport-offline", "transport-offline"},
		{"created", "created"},
		{"quiesce", "quiesce"},
		{"", "unknown"},
		{"something-else", "unknown"},
	}
	for _, tc := range tests {
		got := normalizeDMPathState(tc.raw)
		if got != tc.expected {
			t.Errorf("normalizeDMPathState(%q) = %q, want %q", tc.raw, got, tc.expected)
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
