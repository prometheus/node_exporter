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
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestDMMultipathMetrics(t *testing.T) {
	*procPath = "fixtures/proc"
	*sysPath = "fixtures/sys"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewDMMultipathCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	c := coll.(*dmMultipathCollector)

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

	assertGaugeValue(t, metrics, "device_active", labelMap{"device": "mpathA", "sysfs_name": "dm-5"}, 1)
	assertGaugeValue(t, metrics, "device_active", labelMap{"device": "mpathB", "sysfs_name": "dm-6"}, 1)
	assertGaugeValue(t, metrics, "device_size_bytes", labelMap{"device": "mpathA", "sysfs_name": "dm-5"}, 53687091200)
	assertGaugeValue(t, metrics, `device_paths"`, labelMap{"device": "mpathA", "sysfs_name": "dm-5"}, 4)
	assertGaugeValue(t, metrics, `device_paths"`, labelMap{"device": "mpathB", "sysfs_name": "dm-6"}, 2)

	// mpathA: sdi, sdj, sdk are running; sdl is offline → 3 active, 1 failed.
	assertGaugeValue(t, metrics, "device_paths_active", labelMap{"device": "mpathA", "sysfs_name": "dm-5"}, 3)
	assertGaugeValue(t, metrics, "device_paths_failed", labelMap{"device": "mpathA", "sysfs_name": "dm-5"}, 1)

	// mpathB: sdm, sdn are both running → 2 active, 0 failed.
	assertGaugeValue(t, metrics, "device_paths_active", labelMap{"device": "mpathB", "sysfs_name": "dm-6"}, 2)
	assertGaugeValue(t, metrics, "device_paths_failed", labelMap{"device": "mpathB", "sysfs_name": "dm-6"}, 0)

	assertGaugeValue(t, metrics, "path_state",
		labelMap{"device": "mpathA", "path": "sdi", "state": "running"}, 1)
	assertGaugeValue(t, metrics, "path_state",
		labelMap{"device": "mpathA", "path": "sdl", "state": "offline"}, 1)
}

func TestDMMultipathNoDevices(t *testing.T) {
	*procPath = "fixtures/proc"
	*sysPath = t.TempDir()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewDMMultipathCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	c := coll.(*dmMultipathCollector)

	ch := make(chan prometheus.Metric, 200)
	err = c.Update(ch)
	close(ch)

	if err != ErrNoData {
		t.Fatalf("expected ErrNoData, got %v", err)
	}
}

func TestIsPathActive(t *testing.T) {
	tests := []struct {
		state  string
		active bool
	}{
		{"running", true},
		{"live", true},
		{"offline", false},
		{"blocked", false},
		{"transport-offline", false},
		{"dead", false},
		{"unknown", false},
		{"", false},
	}
	for _, tc := range tests {
		got := isPathActive(tc.state)
		if got != tc.active {
			t.Errorf("isPathActive(%q) = %v, want %v", tc.state, got, tc.active)
		}
	}
}

