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
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testDMMultipathCollector struct {
	mc Collector
}

func (c testDMMultipathCollector) Collect(ch chan<- prometheus.Metric) {
	c.mc.Update(ch)
}

func (c testDMMultipathCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func newTestDMMultipathCollector(logger *slog.Logger) (prometheus.Collector, error) {
	mc, err := NewDMMultipathCollector(logger)
	if err != nil {
		return testDMMultipathCollector{}, err
	}
	return &testDMMultipathCollector{mc}, nil
}

func TestDMMultipathMetrics(t *testing.T) {
	*procPath = "fixtures/proc"
	*sysPath = "fixtures/sys"

	testcase := `# HELP node_dmmultipath_device_active Whether the multipath device-mapper device is active (1) or suspended (0).
# TYPE node_dmmultipath_device_active gauge
node_dmmultipath_device_active{device="mpathA",sysfs_name="dm-5"} 1
node_dmmultipath_device_active{device="mpathB",sysfs_name="dm-6"} 1
# HELP node_dmmultipath_device_info Non-numeric information about a DM-multipath device.
# TYPE node_dmmultipath_device_info gauge
node_dmmultipath_device_info{device="mpathA",sysfs_name="dm-5",uuid="mpath-3600508b1001c1234567890abcdef1234"} 1
node_dmmultipath_device_info{device="mpathB",sysfs_name="dm-6",uuid="mpath-3600508b1001cabcdef4567890123456"} 1
# HELP node_dmmultipath_device_paths Number of paths for a multipath device.
# TYPE node_dmmultipath_device_paths gauge
node_dmmultipath_device_paths{device="mpathA",sysfs_name="dm-5"} 4
node_dmmultipath_device_paths{device="mpathB",sysfs_name="dm-6"} 2
# HELP node_dmmultipath_device_paths_active Number of paths in active state (SCSI running or NVMe live) for a multipath device.
# TYPE node_dmmultipath_device_paths_active gauge
node_dmmultipath_device_paths_active{device="mpathA",sysfs_name="dm-5"} 3
node_dmmultipath_device_paths_active{device="mpathB",sysfs_name="dm-6"} 2
# HELP node_dmmultipath_device_paths_failed Number of paths not in active state for a multipath device.
# TYPE node_dmmultipath_device_paths_failed gauge
node_dmmultipath_device_paths_failed{device="mpathA",sysfs_name="dm-5"} 1
node_dmmultipath_device_paths_failed{device="mpathB",sysfs_name="dm-6"} 0
# HELP node_dmmultipath_device_size_bytes Size of the multipath device in bytes, read from /sys/block/<dm>/size.
# TYPE node_dmmultipath_device_size_bytes gauge
node_dmmultipath_device_size_bytes{device="mpathA",sysfs_name="dm-5"} 5.36870912e+10
node_dmmultipath_device_size_bytes{device="mpathB",sysfs_name="dm-6"} 1.073741824e+11
# HELP node_dmmultipath_path_state Reports the underlying device state for a multipath path, as read from /sys/block/<dev>/device/state.
# TYPE node_dmmultipath_path_state gauge
node_dmmultipath_path_state{device="mpathA",path="sdi",state="running"} 1
node_dmmultipath_path_state{device="mpathA",path="sdj",state="running"} 1
node_dmmultipath_path_state{device="mpathA",path="sdk",state="running"} 1
node_dmmultipath_path_state{device="mpathA",path="sdl",state="offline"} 1
node_dmmultipath_path_state{device="mpathB",path="sdm",state="running"} 1
node_dmmultipath_path_state{device="mpathB",path="sdn",state="running"} 1
`
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
	}))
	c, err := newTestDMMultipathCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
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
