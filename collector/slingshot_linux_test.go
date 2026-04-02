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
	tmpDir := t.TempDir()
	sysRoot := filepath.Join(tmpDir, "sys")
	rootfsRoot := filepath.Join(tmpDir, "rootfs")

	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "properties", "nid"), "42\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "properties", "pid_granule"), "64\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "properties", "current_esm_link_speed"), "32 GT/s\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "current_link_width"), "x16\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "port", "mtu"), "9000\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "port", "link_layer_retry"), "enabled\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "port", "loopback"), "none\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "port", "media"), "copper\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "port", "speed"), "100G\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "port", "link"), "up\n")
	setTestFile(t, filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device", "telemetry", "counter_packets"), "7\n")
	setTestFile(t, filepath.Join(sysRoot, "class", "net", "hsn0", "device", "cxi", "cxi0"), "")
	setTestFile(t, filepath.Join(sysRoot, "class", "net", "hsn0", "address"), "aa:bb:cc:dd:ee:ff\n")

	// This target is outside telemetry root and should never be followed.
	setTestFile(t, filepath.Join(sysRoot, "external", "counter_packets"), "12345\n")

	classDir := filepath.Join(sysRoot, "class", "cxi")
	if err := os.MkdirAll(classDir, 0o755); err != nil {
		t.Fatalf("mkdir class cxi: %v", err)
	}
	deviceTarget := filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0")
	if err := os.Symlink(deviceTarget, filepath.Join(classDir, "cxi0")); err != nil {
		t.Fatalf("create class symlink: %v", err)
	}

	telemetryDir := filepath.Join(deviceTarget, "device", "telemetry")
	if err := os.Symlink(filepath.Join(sysRoot, "external"), filepath.Join(telemetryDir, "subsystem")); err != nil {
		t.Fatalf("create telemetry external symlink: %v", err)
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

	if _, ok := metricByName["node_slingshot_info"]; !ok {
		t.Fatalf("missing node_slingshot_info metric")
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_nid"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 42 {
		t.Fatalf("expected node_slingshot_nid{device=\"cxi0\",interface=\"hsn0\"}=42, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_pid_granule"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 64 {
		t.Fatalf("expected node_slingshot_pid_granule{device=\"cxi0\",interface=\"hsn0\"}=64, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_link_mtu"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 9000 {
		t.Fatalf("expected node_slingshot_link_mtu{device=\"cxi0\",interface=\"hsn0\"}=9000, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_link_speed"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 100e9 {
		t.Fatalf("expected node_slingshot_link_speed{device=\"cxi0\",interface=\"hsn0\"}=100e9, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_pcie_info"], map[string]string{"device": "cxi0", "interface": "hsn0", "slot": "0000:00:00.0"}); !found || got != 1 {
		t.Fatalf("expected node_slingshot_pcie_info{device=\"cxi0\",interface=\"hsn0\",slot=\"0000:00:00.0\"}=1, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_pcie_speed_gts"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 32 {
		t.Fatalf("expected node_slingshot_pcie_speed_gts{device=\"cxi0\",interface=\"hsn0\"}=32, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_pcie_width"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 16 {
		t.Fatalf("expected node_slingshot_pcie_width{device=\"cxi0\",interface=\"hsn0\"}=16, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_link_info"], map[string]string{"device": "cxi0", "interface": "hsn0", "state": "up", "link_layer_retry": "enabled", "loopback": "none", "media": "copper"}); !found || got != 1 {
		t.Fatalf("expected node_slingshot_link_info{device=\"cxi0\",interface=\"hsn0\",state=\"up\",link_layer_retry=\"enabled\",loopback=\"none\",media=\"copper\"}=1, got %v (found=%v)", got, found)
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

	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_counter_packets"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 7 {
		t.Fatalf("expected node_slingshot_telemetry_counter_packets{device=\"cxi0\",interface=\"hsn0\"}=7, got %v (found=%v)", got, found)
	}
	if _, ok := metricByName["node_slingshot_telemetry_subsystem_counter_packets"]; ok {
		t.Fatalf("unexpected telemetry metric from external symlinked directory")
	}
	if _, ok := metricByName["node_slingshot_telemetry_properties_nid"]; ok {
		t.Fatalf("unexpected telemetry metric from non-telemetry device tree")
	}

	if got, ok := gaugeValueByLabels(metricByName["node_slingshot_scrape_errors"], map[string]string{"source": slingshotSourceTelemetry, "interface": "all"}); !ok || got != 0 {
		t.Fatalf("expected telemetry scrape errors to be 0, got %v (found=%v)", got, ok)
	}
}

func TestParseMetricValueSupportsIntegerEncodings(t *testing.T) {
	value, ts, hasTimestamp, err := parseMetricValue("0x10@1000000000")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if value != 16 {
		t.Fatalf("expected value 16, got %v", value)
	}
	if !hasTimestamp || ts != 1000000000 {
		t.Fatalf("unexpected timestamp parse result: ts=%v hasTimestamp=%v", ts, hasTimestamp)
	}

	_, _, _, err = parseMetricValue("not-a-number")
	if err == nil {
		t.Fatal("expected parse error for non-numeric value")
	}
	if !errors.Is(err, errNonNumericMetricValue) {
		t.Fatalf("expected non-numeric parse error, got: %v", err)
	}
}

func TestTelemetryMetricNameCollisionsUseNumericSuffixes(t *testing.T) {
	tmpDir := t.TempDir()
	sysRoot := filepath.Join(tmpDir, "sys")
	rootfsRoot := filepath.Join(tmpDir, "rootfs")

	deviceRoot := filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device")
	setTestFile(t, filepath.Join(deviceRoot, "properties", "nid"), "1\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "a-b"), "11\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "a_b"), "22\n")
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

	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_a_b"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 11 {
		t.Fatalf("expected node_slingshot_telemetry_a_b=11, got %v (found=%v)", got, found)
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_a_b_1"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 22 {
		t.Fatalf("expected node_slingshot_telemetry_a_b_1=22, got %v (found=%v)", got, found)
	}
}

func TestSlingshotInfoCollectorOnlyEmitsInfoMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	sysRoot := filepath.Join(tmpDir, "sys")
	rootfsRoot := filepath.Join(tmpDir, "rootfs")

	deviceRoot := filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device")
	setTestFile(t, filepath.Join(deviceRoot, "properties", "nid"), "42\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "counter_packets"), "7\n")
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
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_nid"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 42 {
		t.Fatalf("expected node_slingshot_nid{device=\"cxi0\",interface=\"hsn0\"}=42, got %v (found=%v)", got, found)
	}
	if _, ok := metricByName["node_slingshot_telemetry_counter_packets"]; ok {
		t.Fatalf("unexpected telemetry metric from info-only collector")
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceInfo, "interface": "all"}); !found || got == 0 {
		t.Fatalf("expected info scraped metrics > 0, got %v (found=%v)", got, found)
	}
	if _, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceTelemetry, "interface": "all"}); found {
		t.Fatalf("unexpected telemetry source accounting from info-only collector")
	}
}

func TestSlingshotMetricsCollectorOnlyEmitsTelemetryMetrics(t *testing.T) {
	tmpDir := t.TempDir()
	sysRoot := filepath.Join(tmpDir, "sys")
	rootfsRoot := filepath.Join(tmpDir, "rootfs")

	deviceRoot := filepath.Join(sysRoot, "devices", "pci0000:00", "0000:00:00.0", "cxi0", "device")
	setTestFile(t, filepath.Join(deviceRoot, "properties", "nid"), "42\n")
	setTestFile(t, filepath.Join(deviceRoot, "telemetry", "counter_packets"), "7\n")
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

	if got, found := gaugeValueByLabels(metricByName["node_slingshot_telemetry_counter_packets"], map[string]string{"device": "cxi0", "interface": "hsn0"}); !found || got != 7 {
		t.Fatalf("expected node_slingshot_telemetry_counter_packets{device=\"cxi0\",interface=\"hsn0\"}=7, got %v (found=%v)", got, found)
	}
	if _, ok := metricByName["node_slingshot_info"]; ok {
		t.Fatalf("unexpected info metric from metrics-only collector")
	}
	if _, ok := metricByName["node_slingshot_nid"]; ok {
		t.Fatalf("unexpected node_slingshot_nid metric from metrics-only collector")
	}
	if got, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceTelemetry, "interface": "all"}); !found || got == 0 {
		t.Fatalf("expected telemetry scraped metrics > 0, got %v (found=%v)", got, found)
	}
	if _, found := gaugeValueByLabels(metricByName["node_slingshot_scraped_metrics"], map[string]string{"source": slingshotSourceInfo, "interface": "all"}); found {
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

func setTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for test file %q: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file %q: %v", path, err)
	}
}
