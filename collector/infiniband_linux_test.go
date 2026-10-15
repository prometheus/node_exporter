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

//go:build !noinfiniband

package collector

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// writeInfiniBandTestDevice creates the minimal set of sysfs files under
// sysRoot/class/infiniband/<name> needed for sysfs.FS to parse a device. When
// withFwVer is false, the required fw_ver file is omitted so that parsing the
// device fails, simulating a restricted/firmware-managed port whose counters
// cannot be read (see https://github.com/prometheus/node_exporter/issues/3823).
func writeInfiniBandTestDevice(t *testing.T, sysRoot, name string, withFwVer bool) {
	t.Helper()

	devPath := filepath.Join(sysRoot, "class", "infiniband", name)
	portPath := filepath.Join(devPath, "ports", "1")
	if err := os.MkdirAll(portPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if withFwVer {
		if err := os.WriteFile(filepath.Join(devPath, "fw_ver"), []byte("1.2.3\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	portFiles := map[string]string{
		"link_layer": "InfiniBand\n",
		"state":      "4: ACTIVE\n",
		"phys_state": "5: LinkUp\n",
		"rate":       "100 Gb/sec (4X EDR)\n",
	}
	for f, content := range portFiles {
		if err := os.WriteFile(filepath.Join(portPath, f), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestInfiniBandCollectorSkipsSysfsReadsForExcludedDevices verifies that a
// device matched by --collector.infiniband.device-exclude is never parsed
// from sysfs. Before the fix, the collector unconditionally parsed every
// device (via sysfs.FS.InfiniBandClass) before applying the device filter, so
// a single unreadable/excluded device would fail the whole scrape.
func TestInfiniBandCollectorSkipsSysfsReadsForExcludedDevices(t *testing.T) {
	sysRoot := t.TempDir()
	// mlx5_0 is a healthy device that should still be scraped normally.
	writeInfiniBandTestDevice(t, sysRoot, "mlx5_0", true)
	// mlx5_12 simulates a restricted port: its fw_ver file cannot be read, so
	// parsing it errors out. It must never be parsed once excluded.
	writeInfiniBandTestDevice(t, sysRoot, "mlx5_12", false)

	oldSysPath, oldExclude := *sysPath, *infinibandDeviceExclude
	*sysPath = sysRoot
	*infinibandDeviceExclude = "^mlx5_12$"
	defer func() {
		*sysPath = oldSysPath
		*infinibandDeviceExclude = oldExclude
	}()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector, err := NewInfiniBandCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 100)
	if err := collector.Update(ch); err != nil {
		t.Fatalf("Update() failed, meaning the excluded device's sysfs files were still read: %v", err)
	}
	close(ch)

	sawExcluded, sawIncluded := false, false
	for m := range ch {
		var out dto.Metric
		if err := m.Write(&out); err != nil {
			t.Fatal(err)
		}
		for _, l := range out.GetLabel() {
			if l.GetName() != "device" {
				continue
			}
			switch l.GetValue() {
			case "mlx5_12":
				sawExcluded = true
			case "mlx5_0":
				sawIncluded = true
			}
		}
	}

	if sawExcluded {
		t.Error("metrics for excluded device mlx5_12 were emitted")
	}
	if !sawIncluded {
		t.Error("expected metrics for included device mlx5_0, got none")
	}
}
