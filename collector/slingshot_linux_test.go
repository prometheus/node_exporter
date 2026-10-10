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

//go:build linux && !noslingshot

package collector

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"
)

var slingshotMetricFQNameRE = regexp.MustCompile(`fqName: "([^"]+)"`)

func TestSlingshotCollectorCollectsTelemetryAndInfoFields(t *testing.T) {
	setSlingshotTestPaths(t, "fixtures/sys")

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewSlingshotCollector(logger)
	if err != nil {
		t.Fatalf("new collector: %v", err)
	}

	metricByName, updateErr := collectSlingshotMetrics(t, c)
	if updateErr != nil {
		t.Fatalf("update failed: %v", updateErr)
	}

	if got, found := gaugeValueByExactLabels(metricByName["node_slingshot_info"], map[string]string{
		"device":           "cxi0",
		"interface":        "hsn0",
		"fru_description":  "HPE Slingshot NIC",
		"part_number":      "P12345",
		"serial_number":    "SN12345",
		"firmware_version": "1.2.3",
		"mac":              "aa:bb:cc:dd:ee:ff",
	}); !found || got != 1 {
		t.Fatalf("expected node_slingshot_info with complete metadata labels, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_nid"], map[string]string{"device": "cxi0"}); !found || got != 42 {
		t.Fatalf("expected node_slingshot_nid{device=\"cxi0\"}=42, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_pid_granule"], map[string]string{"device": "cxi0"}); !found || got != 64 {
		t.Fatalf("expected node_slingshot_pid_granule{device=\"cxi0\"}=64, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_link_mtu"], map[string]string{"device": "cxi0"}); !found || got != 9000 {
		t.Fatalf("expected node_slingshot_link_mtu{device=\"cxi0\"}=9000, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_link_speed"], map[string]string{"device": "cxi0"}); !found || got != 100e9 {
		t.Fatalf("expected node_slingshot_link_speed{device=\"cxi0\"}=100e9, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_pcie_info"], map[string]string{"device": "cxi0", "slot": "0000:00:00.0"}); !found || got != 1 {
		t.Fatalf("expected node_slingshot_pcie_info{device=\"cxi0\",slot=\"0000:00:00.0\"}=1, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_pcie_speed_gts"], map[string]string{"device": "cxi0"}); !found || got != 32 {
		t.Fatalf("expected node_slingshot_pcie_speed_gts{device=\"cxi0\"}=32, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_pcie_width"], map[string]string{"device": "cxi0"}); !found || got != 16 {
		t.Fatalf("expected node_slingshot_pcie_width{device=\"cxi0\"}=16, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_link_info"], map[string]string{"device": "cxi0", "state": "up", "link_layer_retry": "enabled", "loopback": "none", "media": "copper"}); !found || got != 1 {
		t.Fatalf("expected node_slingshot_link_info{device=\"cxi0\",state=\"up\",link_layer_retry=\"enabled\",loopback=\"none\",media=\"copper\"}=1, got %v (found=%v)", got, found)
	}
	if _, ok := metricByName["node_slingshot_link_layer_retry_info"]; ok {
		t.Fatalf("unexpected legacy node_slingshot_link_layer_retry_info metric")
	}
	if _, ok := metricByName["node_slingshot_link_loopback_info"]; ok {
		t.Fatalf("unexpected legacy node_slingshot_link_loopback_info metric")
	}
	if _, ok := metricByName["node_slingshot_link_media_info"]; ok {
		t.Fatalf("unexpected legacy node_slingshot_link_media_info metric")
	}
	if _, ok := metricByName["node_slingshot_link_state_info"]; ok {
		t.Fatalf("unexpected legacy node_slingshot_link_state_info metric")
	}
	if metricHasLabelName(metricByName["node_slingshot_info"], "pcie_speed") || metricHasLabelName(metricByName["node_slingshot_info"], "pcie_slot") || metricHasLabelName(metricByName["node_slingshot_info"], "link_state") {
		t.Fatalf("node_slingshot_info unexpectedly contains old PCIe/link labels")
	}
	if _, ok := metricByName["node_slingshot_pcie_speed_info"]; ok {
		t.Fatalf("unexpected legacy node_slingshot_pcie_speed_info metric")
	}
	if _, ok := metricByName["node_slingshot_pcie_slot_info"]; ok {
		t.Fatalf("unexpected legacy node_slingshot_pcie_slot_info metric")
	}

	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_pct_spt_timeouts"], map[string]string{"device": "cxi0"}); !found || got != 7 {
		t.Fatalf("expected node_slingshot_telemetry_pct_spt_timeouts{device=\"cxi0\"}=7, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_hni_tx_paused"], map[string]string{"device": "cxi0", "index": "7"}); !found || got != 13 {
		t.Fatalf("expected node_slingshot_telemetry_hni_tx_paused{device=\"cxi0\",index=\"7\"}=13, got %v (found=%v)", got, found)
	}
	if _, ok := metricByName["node_slingshot_telemetry_counter_packets"]; ok {
		t.Fatalf("unexpected telemetry metric without slingshot telemetry description")
	}
	if _, ok := metricByName["node_slingshot_telemetry_properties_nid"]; ok {
		t.Fatalf("unexpected telemetry metric from non-telemetry device tree")
	}

	if got, ok := gaugeValueByLabels(metricByName["node_slingshot_scrape_errors"], map[string]string{"source": slingshotSourceTelemetry}); !ok || got != 0 {
		t.Fatalf("expected telemetry scrape errors to be 0, got %v (found=%v)", got, ok)
	}
}

func TestSlingshotTelemetryDoesNotFollowExternalSymlink(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryRoot := filepath.Join(tmpDir, "telemetry")
	externalRoot := filepath.Join(tmpDir, "external")
	setTestFile(t, filepath.Join(externalRoot, "pct_spt_timeouts"), "12345\n")
	if err := os.MkdirAll(telemetryRoot, 0o755); err != nil {
		t.Fatalf("create telemetry directory: %v", err)
	}
	if err := os.Symlink(externalRoot, filepath.Join(telemetryRoot, "subsystem")); err != nil {
		t.Fatalf("create external telemetry symlink: %v", err)
	}

	collector, err := NewSlingshotMetricsCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new metrics collector: %v", err)
	}
	slingshotCollector := collector.(*slingshotCollector)
	ch := make(chan prometheus.Metric, 1)
	scraped, scrapeErrors := slingshotCollector.collectTelemetryMetrics(ch, "cxi0", "hsn0", telemetryRoot)
	close(ch)

	if scraped != 0 || scrapeErrors != 0 || len(ch) != 0 {
		t.Fatalf("external symlink collection = (%d scraped, %d errors, %d metrics), want all zero", scraped, scrapeErrors, len(ch))
	}
}

func TestParseMetricValueSupportsIntegerEncodings(t *testing.T) {
	tests := []struct {
		name             string
		raw              string
		wantValue        float64
		wantTimestamp    float64
		wantHasTimestamp bool
		wantNonNumeric   bool
		wantError        bool
	}{
		{name: "decimal", raw: "42", wantValue: 42},
		{name: "signed", raw: "-7", wantValue: -7},
		{name: "float", raw: "1.25", wantValue: 1.25},
		{name: "scientific", raw: "6.022e23", wantValue: 6.022e23},
		{name: "hex with timestamp", raw: "0x10@1000000000", wantValue: 16, wantTimestamp: 1000000000, wantHasTimestamp: true},
		{name: "empty", raw: "  ", wantNonNumeric: true, wantError: true},
		{name: "non numeric", raw: "not-a-number", wantNonNumeric: true, wantError: true},
		{name: "invalid timestamp", raw: "7@not-a-timestamp", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, timestamp, hasTimestamp, err := parseMetricValue(test.raw)
			if (err != nil) != test.wantError {
				t.Fatalf("parseMetricValue(%q) error = %v, wantError=%v", test.raw, err, test.wantError)
			}
			if test.wantNonNumeric != errors.Is(err, errNonNumericMetricValue) {
				t.Fatalf("parseMetricValue(%q) non-numeric error = %v, want %v", test.raw, err, test.wantNonNumeric)
			}
			if err == nil && (value != test.wantValue || timestamp != test.wantTimestamp || hasTimestamp != test.wantHasTimestamp) {
				t.Fatalf("parseMetricValue(%q) = (%v, %v, %v), want (%v, %v, %v)", test.raw, value, timestamp, hasTimestamp, test.wantValue, test.wantTimestamp, test.wantHasTimestamp)
			}
		})
	}
}

func TestIsPathWithinRoot(t *testing.T) {
	tests := []struct {
		name      string
		root      string
		candidate string
		want      bool
	}{
		{name: "same path", root: "/a/b", candidate: "/a/b", want: true},
		{name: "child path", root: "/a/b", candidate: "/a/b/c", want: true},
		{name: "sibling path", root: "/a/b", candidate: "/a/c", want: false},
		{name: "parent path", root: "/a/b", candidate: "/a", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isPathWithinRoot(test.root, test.candidate); got != test.want {
				t.Fatalf("isPathWithinRoot(%q, %q) = %v, want %v", test.root, test.candidate, got, test.want)
			}
		})
	}
}

func TestNewSlingshotCollectorRejectsInvalidTelemetryFilters(t *testing.T) {
	oldInclude := *slingshotTelemetryMetricsInclude
	oldExclude := *slingshotTelemetryMetricsExclude
	t.Cleanup(func() {
		*slingshotTelemetryMetricsInclude = oldInclude
		*slingshotTelemetryMetricsExclude = oldExclude
	})

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	*slingshotTelemetryMetricsInclude = "["
	*slingshotTelemetryMetricsExclude = ""
	if _, err := NewSlingshotMetricsCollector(logger); err == nil {
		t.Fatal("expected invalid include regexp to fail collector construction")
	}

	*slingshotTelemetryMetricsInclude = ".*"
	*slingshotTelemetryMetricsExclude = "["
	if _, err := NewSlingshotMetricsCollector(logger); err == nil {
		t.Fatal("expected invalid exclude regexp to fail collector construction")
	}

	if _, err := NewSlingshotInfoCollector(logger); err != nil {
		t.Fatalf("info-only collector should ignore telemetry filters: %v", err)
	}
}

func TestSlingshotMetricFilterIncludeExclude(t *testing.T) {
	filter, err := newSlingshotMetricFilter(`^(atu|hni)_`, `miss`)
	if err != nil {
		t.Fatalf("new metric filter: %v", err)
	}

	tests := []struct {
		metric string
		want   bool
	}{
		{metric: "atu_cache_evictions", want: true},
		{metric: "atu_cache_miss_ee", want: false},
		{metric: "hni_tx_paused_7", want: true},
		{metric: "pct_spt_timeouts", want: false},
	}
	for _, test := range tests {
		if got := filter.Allow(test.metric); got != test.want {
			t.Errorf("filter.Allow(%q) = %v, want %v", test.metric, got, test.want)
		}
	}
}

func TestSlingshotCollectorReturnsNoDataWithoutDevices(t *testing.T) {
	tests := []struct {
		name           string
		createClassDir bool
	}{
		{name: "missing class directory"},
		{name: "empty class directory", createClassDir: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sysRoot := filepath.Join(t.TempDir(), "sys")
			if test.createClassDir {
				if err := os.MkdirAll(filepath.Join(sysRoot, "class", "cxi"), 0o755); err != nil {
					t.Fatalf("create class directory: %v", err)
				}
			}
			setSlingshotTestPaths(t, sysRoot)

			collector, err := NewSlingshotCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err != nil {
				t.Fatalf("new collector: %v", err)
			}
			if _, updateErr := collectSlingshotMetrics(t, collector); !errors.Is(updateErr, ErrNoData) {
				t.Fatalf("Update() error = %v, want ErrNoData", updateErr)
			}
		})
	}
}

func TestSlingshotTelemetryTimestampExport(t *testing.T) {
	sysRoot := filepath.Join(t.TempDir(), "sys")
	deviceRoot := createSlingshotTestDevice(t, sysRoot, "cxi0", "0000:00:00.0")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "pct_spt_timeouts"), "7@1609459200000000000\n")
	setSlingshotTestPaths(t, sysRoot)

	oldExportTimestamps := *slingshotExportTimestamps
	t.Cleanup(func() { *slingshotExportTimestamps = oldExportTimestamps })
	*slingshotExportTimestamps = true

	collector, err := NewSlingshotMetricsCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new metrics collector: %v", err)
	}
	metricByName, updateErr := collectSlingshotMetrics(t, collector)
	if updateErr != nil {
		t.Fatalf("update failed: %v", updateErr)
	}

	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_pct_spt_timeouts_timestamp_seconds"], map[string]string{"device": "cxi0"}); !found || got != 1609459200 {
		t.Fatalf("expected normalized timestamp 1609459200, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceTelemetry, "interface": "all"}); !found || got != 2 {
		t.Fatalf("expected two exported telemetry metrics, got %v (found=%v)", got, found)
	}
}

func TestSlingshotMalformedTelemetryIncrementsScrapeErrors(t *testing.T) {
	sysRoot := filepath.Join(t.TempDir(), "sys")
	deviceRoot := createSlingshotTestDevice(t, sysRoot, "cxi0", "0000:00:00.0")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "pct_spt_timeouts"), "7@invalid\n")
	setSlingshotTestPaths(t, sysRoot)

	collector, err := NewSlingshotMetricsCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new metrics collector: %v", err)
	}
	metricByName, updateErr := collectSlingshotMetrics(t, collector)
	if !errors.Is(updateErr, ErrNoData) {
		t.Fatalf("Update() error = %v, want ErrNoData when no telemetry metric is exported", updateErr)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_scrape_errors"], map[string]string{"source": slingshotSourceTelemetry}); !found || got != 1 {
		t.Fatalf("expected one telemetry scrape error, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceTelemetry}); !found || got != 0 {
		t.Fatalf("expected zero exported telemetry metrics, got %v (found=%v)", got, found)
	}
}

func TestSlingshotMetricsCollectorCollectsMultipleDevices(t *testing.T) {
	sysRoot := filepath.Join(t.TempDir(), "sys")
	device0Root := createSlingshotTestDevice(t, sysRoot, "cxi0", "0000:00:00.0")
	device1Root := createSlingshotTestDevice(t, sysRoot, "cxi1", "0000:00:00.1")
	setTestFile(t, filepath.Join(device0Root, "telemetry", "pct_spt_timeouts"), "10\n")
	setTestFile(t, filepath.Join(device1Root, "telemetry", "pct_spt_timeouts"), "20\n")
	setSlingshotTestPaths(t, sysRoot)

	collector, err := NewSlingshotMetricsCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new metrics collector: %v", err)
	}
	metricByName, updateErr := collectSlingshotMetrics(t, collector)
	if updateErr != nil {
		t.Fatalf("update failed: %v", updateErr)
	}

	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_pct_spt_timeouts"], map[string]string{"device": "cxi0"}); !found || got != 10 {
		t.Fatalf("expected cxi0 telemetry value 10, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_pct_spt_timeouts"], map[string]string{"device": "cxi1"}); !found || got != 20 {
		t.Fatalf("expected cxi1 telemetry value 20, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceTelemetry}); !found || got != 2 {
		t.Fatalf("expected two exported telemetry metrics, got %v (found=%v)", got, found)
	}
}

func TestSlingshotTelemetrySymlinkCycleDoesNotLoop(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryRoot := filepath.Join(tmpDir, "telemetry")
	setTestFile(t, filepath.Join(telemetryRoot, "pct_spt_timeouts"), "11\n")
	if err := os.MkdirAll(filepath.Join(telemetryRoot, "nested"), 0o755); err != nil {
		t.Fatalf("create nested telemetry dir: %v", err)
	}
	if err := os.Symlink("..", filepath.Join(telemetryRoot, "nested", "loop")); err != nil {
		t.Fatalf("create loop symlink: %v", err)
	}

	collector, err := NewSlingshotMetricsCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new metrics collector: %v", err)
	}
	slingshotCollector := collector.(*slingshotCollector)

	ch := make(chan prometheus.Metric, 8)
	scraped, scrapeErrors := slingshotCollector.collectTelemetryMetrics(ch, "cxi0", "hsn0", telemetryRoot)
	close(ch)

	if scraped != 1 {
		t.Fatalf("scraped = %d, want 1", scraped)
	}
	if scrapeErrors != 0 {
		t.Fatalf("scrapeErrors = %d, want 0", scrapeErrors)
	}
}

func TestUniqueTelemetryMetricNameCollisionIsDeterministic(t *testing.T) {
	c := &slingshotCollector{
		uniqueTelemetryNameByRaw:       make(map[string]string),
		rawTelemetryByUniqueMetricName: make(map[string]string),
	}

	c.mu.Lock()
	first := c.uniqueTelemetryMetricNameLocked("foo-bar")
	second := c.uniqueTelemetryMetricNameLocked("foo_bar")
	repeatedFirst := c.uniqueTelemetryMetricNameLocked("foo-bar")
	repeatedSecond := c.uniqueTelemetryMetricNameLocked("foo_bar")
	c.mu.Unlock()

	if first != "foo_bar" {
		t.Fatalf("first unique name = %q, want %q", first, "foo_bar")
	}
	if second != "foo_bar_1" {
		t.Fatalf("second unique name = %q, want %q", second, "foo_bar_1")
	}
	if repeatedFirst != first {
		t.Fatalf("repeat for first raw metric changed: got %q, want %q", repeatedFirst, first)
	}
	if repeatedSecond != second {
		t.Fatalf("repeat for second raw metric changed: got %q, want %q", repeatedSecond, second)
	}
}

func TestSlingshotValueParsers(t *testing.T) {
	t.Run("normalize epoch", func(t *testing.T) {
		tests := []struct {
			value float64
			want  float64
		}{
			{value: 1609459200, want: 1609459200},
			{value: 1609459200000000000, want: 1609459200},
		}
		for _, test := range tests {
			if got := normalizeEpochToSeconds(test.value); got != test.want {
				t.Errorf("normalizeEpochToSeconds(%v) = %v, want %v", test.value, got, test.want)
			}
		}
	})

	t.Run("link speed", func(t *testing.T) {
		tests := []struct {
			raw    string
			want   float64
			wantOK bool
		}{
			{raw: "100G", want: 100e9, wantOK: true},
			{raw: "100 GB/S", want: 100e9, wantOK: true},
			{raw: "25.6Gbps", want: 25.6e9, wantOK: true},
			{raw: "100000000000", want: 100e9, wantOK: true},
			{raw: "", wantOK: false},
			{raw: "unknown", wantOK: false},
		}
		for _, test := range tests {
			got, ok := parseLinkSpeedBps(test.raw)
			if got != test.want || ok != test.wantOK {
				t.Errorf("parseLinkSpeedBps(%q) = (%v, %v), want (%v, %v)", test.raw, got, ok, test.want, test.wantOK)
			}
		}
	})
}

func TestSlingshotPCIeParsers(t *testing.T) {
	devicePath := filepath.Join(t.TempDir(), "device")
	setTestFile(t, filepath.Join(devicePath, "properties", "current_esm_link_speed"), "Disabled\n")
	setTestFile(t, filepath.Join(devicePath, "current_link_speed"), "16 GT/s\n")
	setTestFile(t, filepath.Join(devicePath, "current_link_width"), "X8\n")

	if got, ok := readPCIESpeedGTS(devicePath); !ok || got != 16 {
		t.Fatalf("readPCIESpeedGTS() = (%v, %v), want (16, true)", got, ok)
	}
	if got, ok := readPCIELinkWidth(devicePath); !ok || got != 8 {
		t.Fatalf("readPCIELinkWidth() = (%v, %v), want (8, true)", got, ok)
	}
}

func TestTelemetryDescriptionsControlExactAndIndexedMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	sysRoot := filepath.Join(tmpDir, "sys")
	rootfsRoot := filepath.Join(tmpDir, "rootfs")

	deviceRoot := filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device")
	setTestFile(t, filepath.Join(deviceRoot, "properties", "nid"), "1\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "hni_rx_ok_36_to_63"), "99\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "hni_tx_paused_0"), "11\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "hni_tx_paused_7"), "22\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "hni_tx_paused_extra"), "33\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "pct_spt_timeouts"), "44\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "counter_packets"), "55\n")
	setTestFile(t, filepath.Join(sysRoot, "class", "net", "hsn0", "device", "cxi", "cxi0"), "")

	classDir := filepath.Join(sysRoot, "class", "cxi")
	if err := os.MkdirAll(classDir, 0o755); err != nil {
		t.Fatalf("mkdir class cxi: %v", err)
	}
	if err := os.Symlink(filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0"), filepath.Join(classDir, "cxi0")); err != nil {
		t.Fatalf("create class symlink: %v", err)
	}

	oldSysPath := *sysPath
	oldRootfsPath := *rootfsPath
	t.Cleanup(func() {
		*sysPath = oldSysPath
		*rootfsPath = oldRootfsPath
	})
	*sysPath = sysRoot
	*rootfsPath = rootfsRoot

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewSlingshotCollector(logger)
	if err != nil {
		t.Fatalf("new collector: %v", err)
	}

	metricByName, updateErr := collectSlingshotMetrics(t, c)
	if updateErr != nil {
		t.Fatalf("update failed: %v", updateErr)
	}

	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_hni_tx_paused"], map[string]string{"device": "cxi0", "index": "0"}); !found || got != 11 {
		t.Fatalf("expected node_slingshot_telemetry_hni_tx_paused{index=\"0\"}=11, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_hni_rx_ok_36_to_63"], map[string]string{"device": "cxi0"}); !found || got != 99 {
		t.Fatalf("expected node_slingshot_telemetry_hni_rx_ok_36_to_63=99, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_hni_tx_paused"], map[string]string{"device": "cxi0", "index": "7"}); !found || got != 22 {
		t.Fatalf("expected node_slingshot_telemetry_hni_tx_paused{index=\"7\"}=22, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_pct_spt_timeouts"], map[string]string{"device": "cxi0"}); !found || got != 44 {
		t.Fatalf("expected node_slingshot_telemetry_pct_spt_timeouts=44, got %v (found=%v)", got, found)
	}
	if _, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_hni_tx_paused"], map[string]string{"device": "cxi0", "index": "extra"}); found {
		t.Fatalf("unexpected non-numeric indexed suffix metric")
	}
	if _, ok := metricByName["node_slingshot_telemetry_counter_packets"]; ok {
		t.Fatalf("unexpected telemetry metric without slingshot telemetry description")
	}
}

func TestSlingshotInfoCollectorOnlyEmitsInfoMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	sysRoot := filepath.Join(tmpDir, "sys")
	rootfsRoot := filepath.Join(tmpDir, "rootfs")

	deviceRoot := filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device")
	setTestFile(t, filepath.Join(deviceRoot, "properties", "nid"), "42\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "pct_spt_timeouts"), "7\n")
	setTestFile(t, filepath.Join(sysRoot, "class", "net", "hsn0", "device", "cxi", "cxi0"), "")

	classDir := filepath.Join(sysRoot, "class", "cxi")
	if err := os.MkdirAll(classDir, 0o755); err != nil {
		t.Fatalf("mkdir class cxi: %v", err)
	}
	if err := os.Symlink(filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0"), filepath.Join(classDir, "cxi0")); err != nil {
		t.Fatalf("create class symlink: %v", err)
	}

	oldSysPath := *sysPath
	oldRootfsPath := *rootfsPath
	t.Cleanup(func() {
		*sysPath = oldSysPath
		*rootfsPath = oldRootfsPath
	})
	*sysPath = sysRoot
	*rootfsPath = rootfsRoot

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewSlingshotInfoCollector(logger)
	if err != nil {
		t.Fatalf("new info collector: %v", err)
	}

	metricByName, updateErr := collectSlingshotMetrics(t, c)
	if updateErr != nil {
		t.Fatalf("update failed: %v", updateErr)
	}

	if _, ok := metricByName["node_slingshot_info"]; !ok {
		t.Fatalf("missing node_slingshot_info metric")
	}
	if _, ok := metricByName["node_slingshot_nid"]; !ok {
		t.Fatalf("missing node_slingshot_nid metric")
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_nid"], map[string]string{"device": "cxi0"}); !found || got != 42 {
		t.Fatalf("expected node_slingshot_nid{device=\"cxi0\"}=42, got %v (found=%v)", got, found)
	}
	if _, ok := metricByName["node_slingshot_telemetry_pct_spt_timeouts"]; ok {
		t.Fatalf("unexpected telemetry metric from info-only collector")
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceInfo}); !found || got == 0 {
		t.Fatalf("expected info scraped metrics > 0, got %v (found=%v)", got, found)
	}
	if _, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceTelemetry}); found {
		t.Fatalf("unexpected telemetry source accounting from info-only collector")
	}
}

func TestSlingshotMetricsCollectorOnlyEmitsTelemetryMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	sysRoot := filepath.Join(tmpDir, "sys")
	rootfsRoot := filepath.Join(tmpDir, "rootfs")

	deviceRoot := filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device")
	setTestFile(t, filepath.Join(deviceRoot, "properties", "nid"), "42\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "pct_spt_timeouts"), "7\n")
	setTestFile(t, filepath.Join(sysRoot, "class", "net", "hsn0", "device", "cxi", "cxi0"), "")

	classDir := filepath.Join(sysRoot, "class", "cxi")
	if err := os.MkdirAll(classDir, 0o755); err != nil {
		t.Fatalf("mkdir class cxi: %v", err)
	}
	if err := os.Symlink(filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0"), filepath.Join(classDir, "cxi0")); err != nil {
		t.Fatalf("create class symlink: %v", err)
	}

	oldSysPath := *sysPath
	oldRootfsPath := *rootfsPath
	t.Cleanup(func() {
		*sysPath = oldSysPath
		*rootfsPath = oldRootfsPath
	})
	*sysPath = sysRoot
	*rootfsPath = rootfsRoot

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewSlingshotMetricsCollector(logger)
	if err != nil {
		t.Fatalf("new metrics collector: %v", err)
	}

	metricByName, updateErr := collectSlingshotMetrics(t, c)
	if updateErr != nil {
		t.Fatalf("update failed: %v", updateErr)
	}

	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_pct_spt_timeouts"], map[string]string{"device": "cxi0"}); !found || got != 7 {
		t.Fatalf("expected node_slingshot_telemetry_pct_spt_timeouts{device=\"cxi0\"}=7, got %v (found=%v)", got, found)
	}
	if _, ok := metricByName["node_slingshot_info"]; ok {
		t.Fatalf("unexpected info metric from metrics-only collector")
	}
	if _, ok := metricByName["node_slingshot_nid"]; ok {
		t.Fatalf("unexpected node_slingshot_nid metric from metrics-only collector")
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceTelemetry}); !found || got == 0 {
		t.Fatalf("expected telemetry scraped metrics > 0, got %v (found=%v)", got, found)
	}
	if _, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceInfo}); found {
		t.Fatalf("unexpected info source accounting from metrics-only collector")
	}
}

func collectSlingshotMetrics(t *testing.T, c Collector) (map[string][]*dto.Metric, error) {
	t.Helper()

	ch := make(chan prometheus.Metric, 8192)
	err := c.Update(ch)
	close(ch)

	metrics := make(map[string][]*dto.Metric)
	for m := range ch {
		desc := m.Desc().String()
		match := slingshotMetricFQNameRE.FindStringSubmatch(desc)
		if len(match) != 2 {
			t.Fatalf("unable to parse metric fqName from desc: %s", desc)
		}

		pb := &dto.Metric{}
		if writeErr := m.Write(pb); writeErr != nil {
			t.Fatalf("write metric protobuf: %v", writeErr)
		}
		metrics[match[1]] = append(metrics[match[1]], pb)
	}

	return metrics, err
}

func gaugeValueByLabels(metrics []*dto.Metric, labels map[string]string) (float64, bool) {
	for _, metric := range metrics {
		if metric == nil || metric.Gauge == nil {
			continue
		}
		if metricLabelsMatch(metric, labels) {
			return metric.Gauge.GetValue(), true
		}
	}
	return 0, false
}

func gaugeValueByExactLabels(metrics []*dto.Metric, labels map[string]string) (float64, bool) {
	for _, metric := range metrics {
		if metric == nil || metric.Gauge == nil || len(metric.GetLabel()) != len(labels) {
			continue
		}
		if metricLabelsMatch(metric, labels) {
			return metric.Gauge.GetValue(), true
		}
	}
	return 0, false
}

func metricLabelsMatch(metric *dto.Metric, expected map[string]string) bool {
	for key, value := range expected {
		found := false
		for _, label := range metric.GetLabel() {
			if label.GetName() == key && label.GetValue() == value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func metricHasLabelName(metrics []*dto.Metric, name string) bool {
	for _, metric := range metrics {
		for _, label := range metric.GetLabel() {
			if label.GetName() == name {
				return true
			}
		}
	}
	return false
}

func createSlingshotTestDevice(t *testing.T, sysRoot, device, slot string) string {
	t.Helper()

	deviceTarget := filepath.Join(sysRoot, "devices", "pci0000:00", slot, device)
	deviceRoot := filepath.Join(deviceTarget, "device")
	if err := os.MkdirAll(deviceRoot, 0o755); err != nil {
		t.Fatalf("create device directory: %v", err)
	}

	classDir := filepath.Join(sysRoot, "class", "cxi")
	if err := os.MkdirAll(classDir, 0o755); err != nil {
		t.Fatalf("create class directory: %v", err)
	}
	if err := os.Symlink(deviceTarget, filepath.Join(classDir, device)); err != nil {
		t.Fatalf("create class symlink for %s: %v", device, err)
	}
	return deviceRoot
}

func setSlingshotTestPaths(t *testing.T, sysRoot string) {
	t.Helper()

	oldSysPath := *sysPath
	oldRootfsPath := *rootfsPath
	t.Cleanup(func() {
		*sysPath = oldSysPath
		*rootfsPath = oldRootfsPath
	})
	*sysPath = sysRoot
	*rootfsPath = filepath.Join(filepath.Dir(sysRoot), "rootfs")
}

func setTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for test file %q: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file %q: %v", path, err)
	}
}
