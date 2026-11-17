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

//go:build !nopcidevice

package collector

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

func TestPCICollectorWithNameResolution(t *testing.T) {
	// Test the PCI collector with name resolution enabled and compare against expected output
	setPCICollectorFlags(t, "fixtures/sys", true, "", "fixtures/pci.ids")

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewPcideviceCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(&testPCICollector{pc: c})

	// Read expected output from fixture file
	expectedOutput, err := os.ReadFile("fixtures/pcidevice-names-output.txt")
	if err != nil {
		t.Fatal(err)
	}

	err = testutil.GatherAndCompare(reg, strings.NewReader(string(expectedOutput)))
	if err != nil {
		t.Fatal(err)
	}
}

func setPCICollectorFlags(t *testing.T, path string, names bool, include, idsFile string) {
	t.Helper()
	oldPath, oldNames, oldInclude, oldIDs := *sysPath, *pciNames, *pciClassInclude, *pciIDsFile
	t.Cleanup(func() {
		*sysPath, *pciNames, *pciClassInclude, *pciIDsFile = oldPath, oldNames, oldInclude, oldIDs
	})
	*sysPath, *pciNames, *pciClassInclude, *pciIDsFile = path, names, include, idsFile
}

func gatherPCIMetrics(t *testing.T) []*dto.MetricFamily {
	t.Helper()
	c, err := NewPcideviceCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(&testPCICollector{pc: c})
	families, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}
	return families
}

func pciMetricLabel(m *dto.Metric, name string) string {
	for _, label := range m.GetLabel() {
		if label.GetName() == name {
			return label.GetValue()
		}
	}
	return ""
}

func TestPCICollectorClassInclude(t *testing.T) {
	// Include display subclasses and vendors, a separate GPU audio function,
	// and a non-display class containing "03" to check anchoring.
	devices := []struct {
		address, class, vendor, device string
	}{
		{"0000:01:00.0", "0x030200", "0x10de", "0x2231"},
		{"0000:02:00.0", "0x030000", "0x1a03", "0x2000"},
		{"0000:03:00.0", "0x038000", "0x15ad", "0x0405"},
		{"0000:04:00.0", "0x030001", "0x102b", "0x0534"},
		{"0000:01:00.1", "0x040300", "0x10de", "0x1aef"},
		{"0000:05:00.0", "0x060400", "0x1022", "0x1634"},
		{"0000:06:00.0", "0x060403", "0x1022", "0x1634"},
	}
	path := t.TempDir()
	links := filepath.Join(path, "bus/pci/devices")
	if err := os.MkdirAll(links, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		dir := filepath.Join(path, "devices/pci0000:00", device.address)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join("../../../devices/pci0000:00", device.address), filepath.Join(links, device.address)); err != nil {
			t.Fatal(err)
		}
		for name, value := range map[string]string{
			"class": device.class, "vendor": device.vendor, "device": device.device,
			"subsystem_vendor": device.vendor, "subsystem_device": device.device, "revision": "0x01",
			"max_link_speed": "16.0 GT/s PCIe", "current_link_speed": "8.0 GT/s PCIe",
			"max_link_width": "16", "current_link_width": "8", "numa_node": "0",
			"power_state": "D0", "d3cold_allowed": "1", "sriov_drivers_autoprobe": "1",
			"sriov_numvfs": "2", "sriov_totalvfs": "4", "sriov_vf_total_msix": "8",
		} {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(value+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	idsFile := filepath.Join(path, "pci.ids")
	if err := os.WriteFile(idsFile, []byte("10de  NVIDIA\n\t2231  RTX A5000\nC 03  Display controller\n\t02  3D controller\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, names := range []bool{false, true} {
		for _, missingIDs := range []bool{false, true} {
			t.Run(fmt.Sprintf("names=%t/missingIDs=%t", names, missingIDs), func(t *testing.T) {
				file := idsFile
				if missingIDs {
					file = filepath.Join(path, "missing.ids")
				}
				setPCICollectorFlags(t, path, names, "", file)
				baseline := gatherPCIMetrics(t)
				// All 12 families and all 18 samples per device must be present.
				if len(baseline) != 12 {
					t.Fatalf("expected 12 families, got %d", len(baseline))
				}
				total := 0
				for _, family := range baseline {
					total += len(family.GetMetric())
				}
				if total != 18*len(devices) {
					t.Fatalf("expected %d samples, got %d", 18*len(devices), total)
				}
				for _, tc := range []struct {
					name, include string
					selected      int
				}{
					{"empty", "", len(devices)},
					{"all", ".*", len(devices)},
					{"display", "^0x03[0-9a-f]{4}$", 4},
					{"none", "^0xffffff$", 0},
				} {
					t.Run(tc.name, func(t *testing.T) {
						*pciClassInclude = tc.include
						actual := gatherPCIMetrics(t)
						if tc.selected == 0 {
							if len(actual) != 0 {
								t.Fatalf("no-match filter emitted %d families", len(actual))
							}
							return
						}
						if len(actual) != len(baseline) {
							t.Fatalf("expected %d families, got %d", len(baseline), len(actual))
						}
						for i, family := range actual {
							expected := proto.Clone(baseline[i]).(*dto.MetricFamily)
							expected.Metric = nil
							for _, metric := range baseline[i].GetMetric() {
								bus := pciMetricLabel(metric, "bus")
								function := pciMetricLabel(metric, "function")
								if tc.selected == len(devices) || (bus >= "01" && bus <= "04" && function == "0") {
									expected.Metric = append(expected.Metric, metric)
								}
							}
							if !proto.Equal(family, expected) {
								t.Fatalf("%s: filtered samples differ from selected baseline samples", family.GetName())
							}
							if family.GetName() != "node_pcidevice_info" {
								continue
							}
							if len(family.GetMetric()) != tc.selected {
								t.Fatalf("expected %d identities, got %d", tc.selected, len(family.GetMetric()))
							}
							for _, metric := range family.GetMetric() {
								address := fmt.Sprintf("%s:%s:%s.%s", pciMetricLabel(metric, "segment"), pciMetricLabel(metric, "bus"), pciMetricLabel(metric, "device"), pciMetricLabel(metric, "function"))
								for _, device := range devices {
									if device.address == address {
										for label, want := range map[string]string{"class_id": device.class, "vendor_id": device.vendor, "device_id": device.device} {
											if got := pciMetricLabel(metric, label); got != want {
												t.Fatalf("%s %s: want %s, got %s", address, label, want, got)
											}
										}
									}
								}
								if names && address == "0000:01:00.0" {
									want := "RTX A5000"
									if missingIDs {
										want = "2231"
									}
									if got := pciMetricLabel(metric, "device_name"); got != want {
										t.Fatalf("expected device name %q, got %q", want, got)
									}
								}
							}
						}
					})
				}
			})
		}
	}
}

func TestPCICollectorInvalidClassInclude(t *testing.T) {
	setPCICollectorFlags(t, "fixtures/sys", false, "[", "")
	c, err := NewPcideviceCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if c != nil || err == nil || !strings.Contains(err.Error(), "--collector.pcidevice.class-include") {
		t.Fatalf("expected a class-include configuration error, got collector=%v error=%v", c, err)
	}
}

// testPCICollector wraps the PCI collector for testing
type testPCICollector struct {
	pc Collector
}

func (tc *testPCICollector) Collect(ch chan<- prometheus.Metric) {
	sink := make(chan prometheus.Metric)
	go func() {
		err := tc.pc.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	for m := range sink {
		ch <- m
	}
}

func (tc *testPCICollector) Describe(_ chan<- *prometheus.Desc) {
	// No-op for testing
}
