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

//go:build !nopowermeter

package collector

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// buildFakePowerMeterSysfs writes a minimal /sys tree containing ACPI power
// meter devices under <root>/bus/acpi/drivers/power_meter/ and returns root.
// Each meter is described by its ACPI device name (e.g. "ACPI000D:00") and a
// map of sysfs attribute file basename -> content (without trailing newline).
// Optional measures entries are given as "<linkname>:<target-basename>" and
// materialized as symlinks into a shared <root>/devices/LNXSYSTM:00/<target>.
func buildFakePowerMeterSysfs(t *testing.T, meters map[string]map[string]string, measures map[string][]string) string {
	t.Helper()

	root := t.TempDir()
	base := filepath.Join(root, "bus", "acpi", "drivers", "power_meter")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatalf("mkdir base: %v", err)
	}

	devicesDir := filepath.Join(root, "devices", "LNXSYSTM:00")
	if err := os.MkdirAll(devicesDir, 0o755); err != nil {
		t.Fatalf("mkdir devices: %v", err)
	}

	for name, attrs := range meters {
		meterDir := filepath.Join(base, name)
		if err := os.MkdirAll(meterDir, 0o755); err != nil {
			t.Fatalf("mkdir meter %s: %v", name, err)
		}
		for fname, content := range attrs {
			if err := os.WriteFile(filepath.Join(meterDir, fname), []byte(content+"\n"), 0o644); err != nil {
				t.Fatalf("write %s/%s: %v", name, fname, err)
			}
		}

		measureLinks := measures[name]
		if len(measureLinks) > 0 {
			measuresDir := filepath.Join(meterDir, "measures")
			if err := os.MkdirAll(measuresDir, 0o755); err != nil {
				t.Fatalf("mkdir measures %s: %v", name, err)
			}
			for i, entry := range measureLinks {
				parts := strings.SplitN(entry, ":", 2)
				linkName, target := parts[0], parts[1]
				targetDir := filepath.Join(devicesDir, target)
				if err := os.MkdirAll(targetDir, 0o755); err != nil {
					t.Fatalf("mkdir target %s: %v", target, err)
				}
				// The procfs parsePowerMeterMeasures reads symlinks and takes
				// the basename of the target. Relative symlink target with
				// enough ../ to reach devices/LNXSYSTM:00/<target>.
				rel := filepath.Join("..", "..", "..", "..", "devices", "LNXSYSTM:00", target)
				if err := os.Symlink(rel, filepath.Join(measuresDir, linkName)); err != nil {
					t.Fatalf("symlink measures[%d]: %v", i, err)
				}
			}
		} else {
			// Empty measures/ dir so parsePowerMeterMeasures sees no entries.
			if err := os.MkdirAll(filepath.Join(meterDir, "measures"), 0o755); err != nil {
				t.Fatalf("mkdir empty measures: %v", err)
			}
		}
	}

	return root
}

// gatheredMetric is a (family name, metric) pair returned from gatherMetrics.
type gatheredMetric struct {
	Family string
	Metric *dto.Metric
}

// gatherMetrics runs the collector through a fresh registry and returns every
// emitted metric as a dto.Metric plus its family name.
func gatherMetrics(t *testing.T, c *powerMeterCollector) []gatheredMetric {
	t.Helper()
	reg := prometheus.NewRegistry()
	if err := reg.Register(testPowerMeterCollector{c: c}); err != nil {
		t.Fatalf("register: %v", err)
	}
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	var out []gatheredMetric
	for _, fam := range families {
		for _, m := range fam.Metric {
			out = append(out, gatheredMetric{Family: fam.GetName(), Metric: m})
		}
	}
	return out
}

// label returns the value of label `name` on m, or "" if absent.
func label(m *dto.Metric, name string) string {
	for _, l := range m.Label {
		if l.GetName() == name {
			return l.GetValue()
		}
	}
	return ""
}

// findMetric returns the first metric whose family == family AND whose labels
// match the provided map exactly. Returns nil if none matches.
func findMetric(metrics []gatheredMetric, family string, labels map[string]string) *dto.Metric {
	for _, m := range metrics {
		if m.Family != family {
			continue
		}
		match := true
		for k, v := range labels {
			if label(m.Metric, k) != v {
				match = false
				break
			}
		}
		if match && len(m.Metric.Label) == len(labels) {
			return m.Metric
		}
	}
	return nil
}

type testPowerMeterCollector struct {
	c *powerMeterCollector
}

func (t testPowerMeterCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(t, ch)
}

func (t testPowerMeterCollector) Collect(ch chan<- prometheus.Metric) {
	if err := t.c.Update(ch); err != nil {
		panic(err)
	}
}

func TestPowerMeterCollector_FullField(t *testing.T) {
	meters := map[string]map[string]string{
		"ACPI000D:00": {
			"power1_average":              "15000000", // 15 W
			"power1_average_min":          "0",
			"power1_average_max":          "60000000", // 60 W
			"power1_average_interval":     "1000",     // 1 s
			"power1_average_interval_min": "100",      // 0.1 s
			"power1_average_interval_max": "10000",    // 10 s
			"power1_alarm":                "0",
			"power1_cap":                  "25000000",  // 25 W
			"power1_cap_min":              "1000000",   // 1 W
			"power1_cap_max":              "100000000", // 100 W
			"power1_cap_hyst":             "500000",    // 0.5 W
			"power1_accuracy":             "1.50%",     // 1.50
			"power1_is_battery":           "0",
			"power1_model_number":         "ACME PM01",
			"power1_serial_number":        "SN12345",
			"power1_oem_info":             "ACME Corp",
		},
	}
	measures := map[string][]string{
		"ACPI000D:00": {"cpu0:LNXCPU:00", "mem0:LNXMEM:00"},
	}

	root := buildFakePowerMeterSysfs(t, meters, measures)
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = root

	c, err := NewPowerMeterCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new collector: %v", err)
	}

	metrics := gatherMetrics(t, c.(*powerMeterCollector))

	// --- Gauge metrics (power, cap, interval, state, accuracy) ---
	wantGauges := map[string]float64{
		"node_power_meter_average_watts":                15.0,
		"node_power_meter_average_min_watts":            0.0,
		"node_power_meter_average_max_watts":            60.0,
		"node_power_meter_average_interval_seconds":     1.0,
		"node_power_meter_average_interval_min_seconds": 0.1,
		"node_power_meter_average_interval_max_seconds": 10.0,
		"node_power_meter_alarm":                        0.0,
		"node_power_meter_cap_watts":                    25.0,
		"node_power_meter_cap_min_watts":                1.0,
		"node_power_meter_cap_max_watts":                100.0,
		"node_power_meter_cap_hysteresis_watts":         0.5,
		"node_power_meter_is_battery":                   0.0,
		"node_power_meter_accuracy_percent":             1.50,
	}
	for family, want := range wantGauges {
		m := findMetric(metrics, family, map[string]string{"meter": "ACPI000D:00"})
		if m == nil {
			t.Errorf("missing metric %s{meter=ACPI000D:00}", family)
			continue
		}
		if m.GetGauge().GetValue() != want {
			t.Errorf("%s: want %v, got %v", family, want, m.GetGauge().GetValue())
		}
	}

	// --- Info metric ---
	info := findMetric(metrics, "node_power_meter_info", map[string]string{
		"meter":         "ACPI000D:00",
		"model_number":  "ACME PM01",
		"serial_number": "SN12345",
		"oem_info":      "ACME Corp",
	})
	if info == nil {
		t.Fatalf("missing node_power_meter_info{meter=ACPI000D:00}")
	}
	if info.GetGauge().GetValue() != 1.0 {
		t.Errorf("info metric value: want 1.0, got %v", info.GetGauge().GetValue())
	}

	// --- Measures info metric (one series per measured device) ---
	for _, dev := range []string{"LNXCPU:00", "LNXMEM:00"} {
		m := findMetric(metrics, "node_power_meter_measures_info", map[string]string{
			"meter":  "ACPI000D:00",
			"device": dev,
		})
		if m == nil {
			t.Errorf("missing node_power_meter_measures_info{meter=ACPI000D:00, device=%s}", dev)
			continue
		}
		if m.GetGauge().GetValue() != 1.0 {
			t.Errorf("measures_info[%s] value: want 1.0, got %v", dev, m.GetGauge().GetValue())
		}
	}
}

func TestPowerMeterCollector_NoDevices(t *testing.T) {
	root := buildFakePowerMeterSysfs(t, nil, nil)
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = root

	c, err := NewPowerMeterCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new collector: %v", err)
	}
	pc := c.(*powerMeterCollector)

	ch := make(chan prometheus.Metric, 16)
	err = pc.Update(ch)
	if !errors.Is(err, ErrNoData) {
		t.Fatalf("expected ErrNoData when no devices exist, got %v", err)
	}
	close(ch)
	for m := range ch {
		t.Errorf("unexpected metric emitted on ErrNoData path: %v", m)
	}
}

func TestPowerMeterCollector_PartialFields(t *testing.T) {
	// Only average and alarm; every other field is nil/empty.
	meters := map[string]map[string]string{
		"ACPI000D:01": {
			"power1_average": "5000000", // 5 W
			"power1_alarm":   "1",
		},
	}
	root := buildFakePowerMeterSysfs(t, meters, nil)
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = root

	c, err := NewPowerMeterCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new collector: %v", err)
	}
	metrics := gatherMetrics(t, c.(*powerMeterCollector))

	// Present metrics.
	if m := findMetric(metrics, "node_power_meter_average_watts", map[string]string{"meter": "ACPI000D:01"}); m == nil {
		t.Errorf("missing average_watts")
	} else if m.GetGauge().GetValue() != 5.0 {
		t.Errorf("average_watts: want 5, got %v", m.GetGauge().GetValue())
	}
	if m := findMetric(metrics, "node_power_meter_alarm", map[string]string{"meter": "ACPI000D:01"}); m == nil {
		t.Errorf("missing alarm")
	} else if m.GetGauge().GetValue() != 1.0 {
		t.Errorf("alarm: want 1, got %v", m.GetGauge().GetValue())
	}

	// Absent optional metrics.
	for _, family := range []string{
		"node_power_meter_average_min_watts",
		"node_power_meter_average_max_watts",
		"node_power_meter_cap_watts",
		"node_power_meter_cap_min_watts",
		"node_power_meter_cap_max_watts",
		"node_power_meter_cap_hysteresis_watts",
		"node_power_meter_average_interval_seconds",
		"node_power_meter_average_interval_min_seconds",
		"node_power_meter_average_interval_max_seconds",
		"node_power_meter_accuracy_percent",
		"node_power_meter_is_battery",
	} {
		if m := findMetric(metrics, family, map[string]string{"meter": "ACPI000D:01"}); m != nil {
			t.Errorf("did not expect %s when underlying field is absent", family)
		}
	}

	// Info metric still present with only the `meter` label.
	info := findMetric(metrics, "node_power_meter_info", map[string]string{"meter": "ACPI000D:01"})
	if info == nil {
		t.Fatalf("missing info metric for partial meter")
	}
	// No metadata labels beyond `meter`.
	if got, want := len(info.Label), 1; got != want {
		t.Errorf("info metric label count: want %d, got %d", want, got)
	}

	// No measures for this meter → no measures_info series.
	for _, m := range metrics {
		if m.Family == "node_power_meter_measures_info" {
			t.Errorf("did not expect measures_info for meter with empty measures/: got %v", m.Metric)
		}
	}
}

func TestPowerMeterCollector_Filter(t *testing.T) {
	meters := map[string]map[string]string{
		"ACPI000D:00": {"power1_average": "15000000"},
		"ACPI000D:01": {"power1_average": "5000000"},
	}
	root := buildFakePowerMeterSysfs(t, meters, nil)
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = root

	// Configure filter to exclude ACPI000D:00.
	prevFilter := *powerMeterIgnoredMeters
	t.Cleanup(func() { *powerMeterIgnoredMeters = prevFilter })
	*powerMeterIgnoredMeters = `^ACPI000D:00$`

	c, err := NewPowerMeterCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("new collector: %v", err)
	}
	metrics := gatherMetrics(t, c.(*powerMeterCollector))

	if m := findMetric(metrics, "node_power_meter_average_watts", map[string]string{"meter": "ACPI000D:00"}); m != nil {
		t.Errorf("ACPI000D:00 should have been filtered out, but metric is present: %v", m)
	}
	if m := findMetric(metrics, "node_power_meter_average_watts", map[string]string{"meter": "ACPI000D:01"}); m == nil {
		t.Errorf("ACPI000D:01 should still be present after filtering ACPI000D:00")
	} else if m.GetGauge().GetValue() != 5.0 {
		t.Errorf("ACPI000D:01 average: want 5, got %v", m.GetGauge().GetValue())
	}
}
