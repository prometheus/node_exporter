// Copyright 2026 The Prometheus Authors
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

//go:build !noinfiniband
// +build !noinfiniband

package collector

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs/sysfs"
)

// newTestCollector returns a minimally-initialized infinibandCollector rooted
// at sysPath. We only need the fields the EFA helpers touch (sysPath, logger);
// metricDescs and fs are unused by isEFADevice / readEFAHWCounter.
func newTestCollector(sysPath string) *infinibandCollector {
	return &infinibandCollector{
		sysPath: sysPath,
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

// writeFile creates parent dirs and writes content. Test helper to keep
// individual test cases readable.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestIsEFADevice(t *testing.T) {
	cases := []struct {
		name       string
		vendorFile string  // empty => don't create file
		vendorBody string
		want       bool
	}{
		{
			name:       "EFA vendor matches",
			vendorFile: "vendor",
			vendorBody: "0x1d0f\n", // trailing newline is normal sysfs format
			want:       true,
		},
		{
			name:       "EFA vendor without newline still matches",
			vendorFile: "vendor",
			vendorBody: "0x1d0f",
			want:       true,
		},
		{
			name:       "Mellanox vendor does not match",
			vendorFile: "vendor",
			vendorBody: "0x15b3\n",
			want:       false,
		},
		{
			name: "missing vendor file returns false (not an error)",
			// vendorFile empty -> file not created
			want: false,
		},
		{
			name:       "empty vendor file does not match",
			vendorFile: "vendor",
			vendorBody: "",
			want:       false,
		},
		{
			name:       "vendor with extra whitespace still matches after trim",
			vendorFile: "vendor",
			vendorBody: "  0x1d0f  \n",
			want:       true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sysPath := t.TempDir()
			devName := "test_dev"
			if tc.vendorFile != "" {
				writeFile(t,
					filepath.Join(sysPath, "class", "infiniband", devName, "device", tc.vendorFile),
					tc.vendorBody)
			}

			got := newTestCollector(sysPath).isEFADevice(devName)
			if got != tc.want {
				t.Errorf("isEFADevice(%q) = %v, want %v", devName, got, tc.want)
			}
		})
	}
}

func TestReadEFAHWCounter(t *testing.T) {
	sysPath := t.TempDir()
	devName := "rdmap113s0"
	port := uint(1)
	hwDir := filepath.Join(sysPath, "class", "infiniband", devName,
		"ports", "1", "hw_counters")

	// Set up a realistic hw_counters/ subset matching what an EFA NIC exposes.
	writeFile(t, filepath.Join(hwDir, "tx_bytes"), "123456789\n")
	writeFile(t, filepath.Join(hwDir, "rx_bytes"), "987654321")
	writeFile(t, filepath.Join(hwDir, "tx_pkts"), "  42  \n")
	writeFile(t, filepath.Join(hwDir, "rx_drops"), "0\n")
	writeFile(t, filepath.Join(hwDir, "garbage"), "not-a-number\n")
	writeFile(t, filepath.Join(hwDir, "negative"), "-1\n") // uint parse must reject

	c := newTestCollector(sysPath)

	cases := []struct {
		counter string
		wantPtr bool   // true => expect non-nil pointer
		wantVal uint64 // only checked when wantPtr is true
	}{
		{counter: "tx_bytes", wantPtr: true, wantVal: 123456789},
		{counter: "rx_bytes", wantPtr: true, wantVal: 987654321},
		{counter: "tx_pkts", wantPtr: true, wantVal: 42}, // whitespace trim
		{counter: "rx_drops", wantPtr: true, wantVal: 0},
		{counter: "missing", wantPtr: false}, // file does not exist
		{counter: "garbage", wantPtr: false}, // unparseable
		{counter: "negative", wantPtr: false}, // strconv.ParseUint rejects negative
	}

	for _, tc := range cases {
		t.Run(tc.counter, func(t *testing.T) {
			got := c.readEFAHWCounter(devName, port, tc.counter)
			if tc.wantPtr {
				if got == nil {
					t.Fatalf("readEFAHWCounter(%q) = nil, want non-nil", tc.counter)
				}
				if *got != tc.wantVal {
					t.Errorf("readEFAHWCounter(%q) = %d, want %d", tc.counter, *got, tc.wantVal)
				}
			} else if got != nil {
				t.Errorf("readEFAHWCounter(%q) = %d, want nil", tc.counter, *got)
			}
		})
	}
}

// TestEFAVendorIDConstant guards against typos in the vendor ID literal —
// 0x1d0f is the AWS vendor ID assigned by PCI-SIG. A regression here would
// silently make every EFA device fall through the IB path with zero bytes.
func TestEFAVendorIDConstant(t *testing.T) {
	if !strings.HasPrefix(efaVendorID, "0x") {
		t.Errorf("efaVendorID = %q, want 0x-prefixed hex string", efaVendorID)
	}
	if efaVendorID != "0x1d0f" {
		t.Errorf("efaVendorID = %q, want 0x1d0f (AWS PCI vendor ID)", efaVendorID)
	}
}

// newTestCollectorWithDescs returns a collector with the minimum metric
// descriptions wired up for pushEFACounter to find them. Only a handful of
// metric names are exercised in pushEFACounter tests, so we register just
// those rather than reproducing the full NewInfiniBandCollector init block.
func newTestCollectorWithDescs(sysPath string, metricNames ...string) *infinibandCollector {
	c := newTestCollector(sysPath)
	c.subsystem = "infiniband"
	c.metricDescs = make(map[string]*prometheus.Desc)
	for _, name := range metricNames {
		c.metricDescs[name] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, name),
			"test",
			[]string{"device", "port"},
			nil,
		)
	}
	return c
}

// metricValue extracts the float value from a prometheus.Metric for assertion.
func metricValue(t *testing.T, m prometheus.Metric) float64 {
	t.Helper()
	pb := &dto.Metric{}
	if err := m.Write(pb); err != nil {
		t.Fatalf("metric write: %v", err)
	}
	if pb.Counter != nil {
		return pb.Counter.GetValue()
	}
	if pb.Gauge != nil {
		return pb.Gauge.GetValue()
	}
	t.Fatalf("metric has neither Counter nor Gauge: %v", pb)
	return 0
}

func TestPushEFACounter(t *testing.T) {
	sysPath := t.TempDir()
	devName := "rdmap113s0"
	port := uint(1)
	hwDir := filepath.Join(sysPath, "class", "infiniband", devName,
		"ports", "1", "hw_counters")

	writeFile(t, filepath.Join(hwDir, "tx_bytes"), "9999\n")
	// Intentionally do NOT create rx_bytes — pushEFACounter must emit nothing
	// when the underlying read returns nil. This is the silent-skip contract
	// existing pushCounter callers rely on.

	c := newTestCollectorWithDescs(sysPath,
		"port_data_transmitted_bytes_total",
		"port_data_received_bytes_total",
	)

	t.Run("present counter emits one metric with expected value", func(t *testing.T) {
		ch := make(chan prometheus.Metric, 1)
		c.pushEFACounter(ch, "port_data_transmitted_bytes_total", "tx_bytes", devName, port, "1")
		close(ch)

		var collected []prometheus.Metric
		for m := range ch {
			collected = append(collected, m)
		}
		if len(collected) != 1 {
			t.Fatalf("got %d metrics, want 1", len(collected))
		}
		if v := metricValue(t, collected[0]); v != 9999 {
			t.Errorf("got value %v, want 9999", v)
		}
	})

	t.Run("missing counter file emits nothing", func(t *testing.T) {
		ch := make(chan prometheus.Metric, 1)
		c.pushEFACounter(ch, "port_data_received_bytes_total", "rx_bytes", devName, port, "1")
		close(ch)

		var collected []prometheus.Metric
		for m := range ch {
			collected = append(collected, m)
		}
		if len(collected) != 0 {
			t.Errorf("got %d metrics, want 0 (file missing should silently skip)", len(collected))
		}
	})
}

// --------------------------------------------------------------------------
// Mock-sysfs E2E test for Update(): builds a minimal but procfs-compatible
// /sys/class/infiniband tree containing both an EFA-style device and a
// Mellanox-style IB device, then verifies that Update() routes each through
// the correct code path.
// --------------------------------------------------------------------------

// ibDeviceSpec describes one InfiniBand-class device worth of sysfs files
// to materialize under sysPath. The procfs library requires a specific
// minimal file set per device/port; this struct keeps the fixture setup
// declarative so test cases stay readable.
type ibDeviceSpec struct {
	name        string
	vendor      string // device/vendor (empty = no file => non-EFA)
	fwVer       string // fw_ver (required by procfs)
	port        uint
	state       string // e.g. "4: ACTIVE"
	physState   string // e.g. "5: LinkUp"
	rate        string // e.g. "100 Gb/sec (4X EDR)"
	linkLayer   string // e.g. "InfiniBand" or "Ethernet"
	ibCounters  map[string]string // counters/<file> => content (IB devices)
	efaHWCounts map[string]string // hw_counters/<file> => content (EFA devices)
}

// writeSpec materializes a device spec under sysPath/class/infiniband/<name>/.
// counters/ directory is always created (procfs requires it) even when empty.
func writeSpec(t *testing.T, sysPath string, s ibDeviceSpec) {
	t.Helper()
	devDir := filepath.Join(sysPath, "class", "infiniband", s.name)
	portDir := filepath.Join(devDir, "ports", "1")
	countersDir := filepath.Join(portDir, "counters")

	// Mandatory device-level files. fw_ver is required by procfs.
	writeFile(t, filepath.Join(devDir, "fw_ver"), s.fwVer)

	if s.vendor != "" {
		writeFile(t, filepath.Join(devDir, "device", "vendor"), s.vendor)
	}

	// Mandatory port-level files.
	writeFile(t, filepath.Join(portDir, "link_layer"), s.linkLayer)
	writeFile(t, filepath.Join(portDir, "state"), s.state)
	writeFile(t, filepath.Join(portDir, "phys_state"), s.physState)
	writeFile(t, filepath.Join(portDir, "rate"), s.rate)

	// counters/ must exist as a directory even if no counters are populated.
	if err := os.MkdirAll(countersDir, 0o755); err != nil {
		t.Fatalf("mkdir counters: %v", err)
	}
	for f, v := range s.ibCounters {
		writeFile(t, filepath.Join(countersDir, f), v)
	}

	// EFA devices store bytes/packets in hw_counters/, IB error metrics in
	// the same directory. Only populate if the test wants to exercise that
	// path.
	for f, v := range s.efaHWCounts {
		writeFile(t, filepath.Join(portDir, "hw_counters", f), v)
	}
}

// newE2ECollector builds a fully-wired infinibandCollector against the given
// sysfs root with all the metric descriptions our Update() emits. Mirrors
// what NewInfiniBandCollector would do, but pointed at a fake sysfs and with
// no side-effects on global *sysPath flag.
func newE2ECollector(t *testing.T, sysPath string) *infinibandCollector {
	t.Helper()
	fs, err := sysfs.NewFS(sysPath)
	if err != nil {
		t.Fatalf("sysfs.NewFS(%q): %v", sysPath, err)
	}

	// Same metric descriptions wired by NewInfiniBandCollector. We only need
	// names; descriptions and help text don't affect behavior under test.
	names := []string{
		// always emitted
		"state_id", "physical_state_id", "rate_bytes_per_second",
		// shared between EFA and IB paths
		"port_data_transmitted_bytes_total", "port_data_received_bytes_total",
		"port_packets_transmitted_total", "port_packets_received_total",
		// IB-only counters our Update path touches
		"legacy_multicast_packets_received_total", "legacy_multicast_packets_transmitted_total",
		"legacy_data_received_bytes_total", "legacy_packets_received_total",
		"legacy_unicast_packets_received_total", "legacy_unicast_packets_transmitted_total",
		"legacy_data_transmitted_bytes_total", "legacy_packets_transmitted_total",
		"excessive_buffer_overrun_errors_total", "link_downed_total",
		"link_error_recovery_total", "local_link_integrity_errors_total",
		"multicast_packets_received_total", "multicast_packets_transmitted_total",
		"port_constraint_errors_received_total", "port_constraint_errors_transmitted_total",
		"port_discards_received_total", "port_discards_transmitted_total",
		"port_errors_received_total", "port_transmit_wait_total",
		"unicast_packets_received_total", "unicast_packets_transmitted_total",
		"port_receive_remote_physical_errors_total", "port_receive_switch_relay_errors_total",
		"symbol_error_total", "vl15_dropped_total",
		// EFA-only counters
		"efa_rx_drops_total", "efa_retrans_packets_total", "efa_retrans_bytes_total",
		"efa_retrans_timeout_events_total", "efa_unresponsive_remote_events_total",
		"efa_impaired_remote_conn_events_total",
		"efa_rdma_read_bytes_total", "efa_rdma_write_bytes_total",
	}
	descs := make(map[string]*prometheus.Desc, len(names))
	for _, n := range names {
		descs[n] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "infiniband", n),
			"test",
			[]string{"device", "port"},
			nil,
		)
	}
	return &infinibandCollector{
		fs:          fs,
		sysPath:     sysPath,
		metricDescs: descs,
		logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		subsystem:   "infiniband",
	}
}

// collectedMetric is a flattened representation of a single emitted metric
// for assertion lookup, keyed by (metric name, device, port).
type collectedMetric struct {
	name   string
	device string
	port   string
	value  float64
}

func collectAll(t *testing.T, c *infinibandCollector) []collectedMetric {
	t.Helper()
	ch := make(chan prometheus.Metric, 256)
	if err := c.Update(ch); err != nil {
		t.Fatalf("Update: %v", err)
	}
	close(ch)

	var out []collectedMetric
	for m := range ch {
		pb := &dto.Metric{}
		if err := m.Write(pb); err != nil {
			t.Fatalf("metric write: %v", err)
		}
		var value float64
		switch {
		case pb.Counter != nil:
			value = pb.Counter.GetValue()
		case pb.Gauge != nil:
			value = pb.Gauge.GetValue()
		default:
			continue // skip non-counter/gauge (info metric etc.)
		}
		name := m.Desc().String()
		// Desc().String() is verbose ("Desc{fqName: ...}"); extract metric
		// name by parsing the FQName from labels would be cleaner, but for a
		// keyed lookup we just normalize via the label map.
		var device, port string
		for _, lp := range pb.Label {
			switch lp.GetName() {
			case "device":
				device = lp.GetValue()
			case "port":
				port = lp.GetValue()
			}
		}
		out = append(out, collectedMetric{name: name, device: device, port: port, value: value})
	}
	return out
}

// findMetric returns the value of the first metric whose Desc string contains
// the given metric short name and labels match. Returns (value, true) on
// match, (0, false) otherwise.
func findMetric(metrics []collectedMetric, shortName, device, port string) (float64, bool) {
	for _, m := range metrics {
		if m.device == device && m.port == port && strings.Contains(m.name, shortName) {
			return m.value, true
		}
	}
	return 0, false
}

func TestUpdate_EFAReadsHWCounters(t *testing.T) {
	sysPath := t.TempDir()

	// EFA device with realistic vendor and hw_counters/.
	writeSpec(t, sysPath, ibDeviceSpec{
		name:      "rdmap_test",
		vendor:    "0x1d0f\n",
		fwVer:     "1.0.0",
		port:      1,
		state:     "4: ACTIVE",
		physState: "5: LinkUp",
		rate:      "400 Gb/sec (4X NDR)",
		linkLayer: "Ethernet", // EFA reports Ethernet link_layer
		efaHWCounts: map[string]string{
			"tx_bytes":                  "111111\n",
			"rx_bytes":                  "222222\n",
			"tx_pkts":                   "333\n",
			"rx_pkts":                   "444\n",
			"rx_drops":                  "5\n",
			"retrans_pkts":              "6\n",
			"retrans_bytes":             "7\n",
			"retrans_timeout_events":    "8\n",
			"unresponsive_remote_events": "9\n",
			"impaired_remote_conn_events": "10\n",
			"rdma_read_bytes":           "1000\n",
			"rdma_write_bytes":          "2000\n",
		},
	})

	c := newE2ECollector(t, sysPath)
	metrics := collectAll(t, c)

	cases := []struct {
		metric string
		want   float64
	}{
		// Bytes/packets emitted under existing IB-equivalent metric names
		// (this is the key behavior — EFA data appears in normal IB dashboards).
		{"port_data_transmitted_bytes_total", 111111},
		{"port_data_received_bytes_total", 222222},
		{"port_packets_transmitted_total", 333},
		{"port_packets_received_total", 444},
		// EFA-only diagnostic counters under efa_* prefix.
		{"efa_rx_drops_total", 5},
		{"efa_retrans_packets_total", 6},
		{"efa_retrans_bytes_total", 7},
		{"efa_retrans_timeout_events_total", 8},
		{"efa_unresponsive_remote_events_total", 9},
		{"efa_impaired_remote_conn_events_total", 10},
		{"efa_rdma_read_bytes_total", 1000},
		{"efa_rdma_write_bytes_total", 2000},
	}
	for _, tc := range cases {
		got, ok := findMetric(metrics, tc.metric, "rdmap_test", "1")
		if !ok {
			t.Errorf("%s missing from emitted metrics", tc.metric)
			continue
		}
		if got != tc.want {
			t.Errorf("%s = %v, want %v", tc.metric, got, tc.want)
		}
	}

	// state/physical/rate should always be emitted regardless of device type.
	for _, name := range []string{"state_id", "physical_state_id", "rate_bytes_per_second"} {
		if _, ok := findMetric(metrics, name, "rdmap_test", "1"); !ok {
			t.Errorf("EFA device missing %s metric", name)
		}
	}
}

func TestUpdate_IBDeviceUnaffectedByEFAPath(t *testing.T) {
	sysPath := t.TempDir()

	// Mellanox IB device — has vendor file but with non-EFA vendor ID,
	// confirming isEFADevice's negative path: presence of vendor file alone
	// must not trigger EFA mode.
	writeSpec(t, sysPath, ibDeviceSpec{
		name:      "mlx5_test",
		vendor:    "0x15b3\n", // Mellanox PCI vendor
		fwVer:     "20.36.1010",
		port:      1,
		state:     "4: ACTIVE",
		physState: "5: LinkUp",
		rate:      "100 Gb/sec (4X EDR)",
		linkLayer: "InfiniBand",
		ibCounters: map[string]string{
			"port_xmit_data":    "999",
			"port_rcv_data":     "888",
			"port_xmit_packets": "77",
			"port_rcv_packets":  "66",
			"link_downed":       "1",
			"symbol_error":      "2",
			"port_xmit_wait":    "3",
		},
		// Intentionally also populate hw_counters/tx_bytes — if the EFA
		// branch were mistakenly taken, port_data_transmitted_bytes_total
		// would equal 99999 instead of being driven by port_xmit_data. This
		// is the regression guard.
		efaHWCounts: map[string]string{
			"tx_bytes": "99999",
			"rx_bytes": "88888",
		},
	})

	c := newE2ECollector(t, sysPath)
	metrics := collectAll(t, c)

	// IB device must use port.Counters.PortXmitData, NOT hw_counters/tx_bytes.
	// Note: procfs library does not apply ×4 to port_xmit_data; the value is
	// whatever the kernel reports. So we assert "value comes from counters/
	// path" by asserting it does NOT equal the hw_counters/ value we planted.
	if got, ok := findMetric(metrics, "port_data_transmitted_bytes_total", "mlx5_test", "1"); !ok {
		t.Error("port_data_transmitted_bytes_total missing for IB device")
	} else if got == 99999 {
		t.Errorf("port_data_transmitted_bytes_total = %v — EFA branch was taken for Mellanox device (vendor=0x15b3), this is a regression", got)
	}

	// IB-only error counters must still be emitted (EFA path skips these).
	for _, name := range []string{
		"link_downed_total",
		"symbol_error_total",
		"port_transmit_wait_total",
	} {
		if _, ok := findMetric(metrics, name, "mlx5_test", "1"); !ok {
			t.Errorf("%s missing for IB device — EFA branch may have swallowed it", name)
		}
	}

	// EFA-only diagnostic metrics must NOT be emitted for an IB device.
	for _, name := range []string{
		"efa_rx_drops_total",
		"efa_retrans_packets_total",
		"efa_rdma_read_bytes_total",
	} {
		if _, ok := findMetric(metrics, name, "mlx5_test", "1"); ok {
			t.Errorf("%s emitted for non-EFA device — efa_* metrics leaked into IB path", name)
		}
	}
}

func TestUpdate_MixedFleet_EFAAndIBCoexist(t *testing.T) {
	sysPath := t.TempDir()

	writeSpec(t, sysPath, ibDeviceSpec{
		name: "rdmap_a", vendor: "0x1d0f\n", fwVer: "1.0",
		port: 1, state: "4: ACTIVE", physState: "5: LinkUp",
		rate: "400 Gb/sec (4X NDR)", linkLayer: "Ethernet",
		efaHWCounts: map[string]string{"tx_bytes": "5000"},
	})
	writeSpec(t, sysPath, ibDeviceSpec{
		name: "mlx5_a", vendor: "0x15b3\n", fwVer: "20.0",
		port: 1, state: "4: ACTIVE", physState: "5: LinkUp",
		rate: "100 Gb/sec (4X EDR)", linkLayer: "InfiniBand",
		ibCounters: map[string]string{"port_xmit_data": "6000"},
	})

	c := newE2ECollector(t, sysPath)
	metrics := collectAll(t, c)

	// EFA device reads from hw_counters; IB device from counters. Both must
	// appear in the same Update() pass.
	if v, ok := findMetric(metrics, "port_data_transmitted_bytes_total", "rdmap_a", "1"); !ok || v != 5000 {
		t.Errorf("EFA device port_data_transmitted_bytes_total = %v (ok=%v), want 5000", v, ok)
	}
	if v, ok := findMetric(metrics, "port_data_transmitted_bytes_total", "mlx5_a", "1"); !ok {
		t.Error("IB device port_data_transmitted_bytes_total missing")
	} else if v == 5000 {
		t.Errorf("IB device leaked EFA-path value: got %v", v)
	}
}
