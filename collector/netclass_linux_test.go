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

//go:build !nonetclass && linux

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
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testNetClassCollector struct {
	ncc Collector
}

func (c testNetClassCollector) Collect(ch chan<- prometheus.Metric) {
	c.ncc.Update(ch)
}

func (c testNetClassCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// newNetClassCollectorForSysfs points the netclass collector at a synthetic
// sysfs tree rooted at sysDir.
func newNetClassCollectorForSysfs(t *testing.T, sysDir string) *netClassCollector {
	t.Helper()

	*sysPath = sysDir
	*netclassNetlink = false
	*netclassIgnoredDevices = "^$"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewNetClassCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	return coll.(*netClassCollector)
}

// writeNetClassIface creates the attribute files of an interface that is
// present in the synthetic sysfs tree.
func writeNetClassIface(t *testing.T, netDir, iface string, attrs map[string]string) {
	t.Helper()

	dir := filepath.Join(netDir, iface)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	for name, value := range attrs {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(value+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestNetClassDeviceVanishedDuringScan covers the race behind issue #1915.
// /sys/class/net is listed first and every interface is read afterwards, so an
// interface removed in between used to fail the collector outright and drop the
// metrics of all the interfaces that were still there.
func TestNetClassDeviceVanishedDuringScan(t *testing.T) {
	sysDir := t.TempDir()
	netDir := filepath.Join(sysDir, "class", "net")
	if err := os.MkdirAll(netDir, 0o755); err != nil {
		t.Fatal(err)
	}

	writeNetClassIface(t, netDir, "eth0", map[string]string{
		"address":   "01:02:03:04:05:06",
		"broadcast": "ff:ff:ff:ff:ff:ff",
		"duplex":    "full",
		"flags":     "0x1003",
		"mtu":       "1500",
		"operstate": "up",
	})

	// Entries in /sys/class/net are symlinks into /sys/devices, so a dangling
	// symlink is what an interface looks like once the kernel removed it.
	if err := os.Symlink(filepath.Join(sysDir, "devices", "virtual", "veth0"), filepath.Join(netDir, "veth0")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.ReadDir(filepath.Join(netDir, "veth0")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("fixture is supposed to read as missing, got err: %v", err)
	}

	c := newNetClassCollectorForSysfs(t, sysDir)

	ch := make(chan prometheus.Metric, 100)
	err := c.Update(ch)
	close(ch)
	if err != nil {
		t.Fatalf("expected the vanished interface to be skipped, got err: %v", err)
	}
	for m := range ch {
		if strings.Contains(m.Desc().String(), "veth0") {
			t.Fatalf("vanished interface should not be reported, got %s", m.Desc())
		}
	}

	expected := `# HELP node_network_info Non-numeric data from /sys/class/net/<iface>, value is always 1.
# TYPE node_network_info gauge
node_network_info{address="01:02:03:04:05:06",adminstate="up",broadcast="ff:ff:ff:ff:ff:ff",device="eth0",duplex="full",ifalias="",operstate="up"} 1
# HELP node_network_up Value is 1 if operstate is 'up', 0 otherwise.
# TYPE node_network_up gauge
node_network_up{device="eth0"} 1
`
	if err := testutil.CollectAndCompare(testNetClassCollector{ncc: c}, strings.NewReader(expected), "node_network_info", "node_network_up"); err != nil {
		t.Fatal(err)
	}
}

// TestNetClassMissingSysfsReportsNoData guards the behaviour added in #1986:
// when /sys/class/net itself cannot be read there is nothing to report, which
// is different from one interface going away mid-scan.
func TestNetClassMissingSysfsReportsNoData(t *testing.T) {
	c := newNetClassCollectorForSysfs(t, t.TempDir())

	ch := make(chan prometheus.Metric, 100)
	err := c.Update(ch)
	close(ch)

	if !errors.Is(err, ErrNoData) {
		t.Fatalf("expected ErrNoData when /sys/class/net is absent, got: %v", err)
	}
	for m := range ch {
		t.Fatalf("expected no metrics, got %s", m.Desc())
	}
}
