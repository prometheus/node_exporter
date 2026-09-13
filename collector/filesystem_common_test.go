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

//go:build !nofilesystem && (linux || freebsd || netbsd || openbsd || darwin || dragonfly || aix)

package collector

import (
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// testFilesystemCollector exposes collectStats through the Collector interface
// so metric emission can be checked without reading a real mount table.
type testFilesystemCollector struct {
	collector *filesystemCollector
	stats     []filesystemStats
}

func (c testFilesystemCollector) Collect(ch chan<- prometheus.Metric) {
	c.collector.collectStats(ch, c.stats)
}

func (c testFilesystemCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func newFilesystemCollectorForStats(t *testing.T) *filesystemCollector {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewFilesystemCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	return coll.(*filesystemCollector)
}

// gatherFilesystemStats gathers stats through a pedantic registry, which
// rejects a label set emitted twice exactly the way a scrape does. It returns
// the number of series gathered per metric name.
func gatherFilesystemStats(t *testing.T, c *filesystemCollector, stats []filesystemStats) map[string]int {
	t.Helper()
	registry := prometheus.NewPedanticRegistry()
	if err := registry.Register(testFilesystemCollector{collector: c, stats: stats}); err != nil {
		t.Fatal(err)
	}
	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("gathering filesystem metrics: %s", err)
	}
	seriesPerMetric := make(map[string]int, len(families))
	for _, family := range families {
		seriesPerMetric[family.GetName()] = len(family.GetMetric())
	}
	return seriesPerMetric
}

func nfsMountStat(fsType, minor string) filesystemStats {
	return filesystemStats{
		labels: filesystemLabels{
			device:     "storagesystem:/exports/home",
			mountPoint: "/home",
			fsType:     fsType,
			major:      "0",
			minor:      minor,
		},
		size:      536870912000,
		free:      525147832320,
		avail:     525147832320,
		files:     33554432,
		filesFree: 33216512,
		purgeable: -1,
	}
}

// everyFilesystemMetric lists the metrics emitted for a mount that reports no
// device error. purgeable is left out because it is only emitted on platforms
// that report it.
var everyFilesystemMetric = []string{
	"node_filesystem_avail_bytes",
	"node_filesystem_device_error",
	"node_filesystem_files",
	"node_filesystem_files_free",
	"node_filesystem_free_bytes",
	"node_filesystem_mount_info",
	"node_filesystem_readonly",
	"node_filesystem_size_bytes",
}

func checkFilesystemSeries(t *testing.T, stats []filesystemStats, expected map[string]int) {
	t.Helper()
	c := newFilesystemCollectorForStats(t)
	series := gatherFilesystemStats(t, c, stats)
	for _, name := range everyFilesystemMetric {
		want, tracked := expected[name]
		if !tracked {
			t.Errorf("%s: not accounted for in the test expectations", name)
			continue
		}
		if series[name] != want {
			t.Errorf("%s: expected %d series, got %d", name, want, series[name])
		}
	}
	for name := range series {
		if _, tracked := expected[name]; !tracked {
			t.Errorf("%s: unexpected metric emitted with %d series", name, series[name])
		}
	}
}

// A multihomed NFS export is listed once per server address, so the same device
// at the same mount point shows up several times with the same fstype but a
// device number of its own per address. Only node_filesystem_mount_info carries
// major and minor; every other filesystem metric describes the same series each
// time, which used to fail the whole scrape.
func TestFilesystemDeduplicatesRepeatMountWithDifferentDeviceNumbers(t *testing.T) {
	expected := map[string]int{"node_filesystem_mount_info": 2}
	for _, name := range everyFilesystemMetric {
		if _, ok := expected[name]; !ok {
			expected[name] = 1
		}
	}

	checkFilesystemSeries(t, []filesystemStats{
		nfsMountStat("nfs", "41"),
		nfsMountStat("nfs", "42"),
	}, expected)
}

// The same device numbers reported under two filesystem types are two series
// for the per-fstype metrics but a single series for mount_info, which does not
// carry fstype. This is the mirror image of the case above and needs its own
// deduplication key.
func TestFilesystemDeduplicatesMountInfoAcrossFilesystemTypes(t *testing.T) {
	expected := map[string]int{"node_filesystem_mount_info": 1}
	for _, name := range everyFilesystemMetric {
		if _, ok := expected[name]; !ok {
			expected[name] = 2
		}
	}

	checkFilesystemSeries(t, []filesystemStats{
		nfsMountStat("nfs", "41"),
		nfsMountStat("nfs4", "41"),
	}, expected)
}

// Mount table entries that repeat a filesystem without adding anything a metric
// carries are reported once, which is what the existing deduplication intended.
func TestFilesystemDeduplicatesIdenticalRepeatMount(t *testing.T) {
	expected := map[string]int{}
	for _, name := range everyFilesystemMetric {
		expected[name] = 1
	}

	checkFilesystemSeries(t, []filesystemStats{
		nfsMountStat("nfs", "41"),
		nfsMountStat("nfs", "41"),
	}, expected)
}
