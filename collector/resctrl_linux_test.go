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

//go:build !noresctrl

package collector

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/prometheus/procfs/resctrlfs"

	"github.com/prometheus/node_exporter/config"
)

type testResctrlCollector struct {
	rc Collector
}

func (c testResctrlCollector) Collect(ch chan<- prometheus.Metric) {
	c.rc.Update(ch)
}

func (c testResctrlCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func newTestResctrlCollector(path string, logger *slog.Logger) (prometheus.Collector, error) {
	cfg := flagConfig()
	cfg.Resctrl.Path = path
	rc, err := NewResctrlCollector(cfg, logger)
	if err != nil {
		return testResctrlCollector{}, err
	}
	return testResctrlCollector{rc}, nil
}

func TestResctrlMetrics(t *testing.T) {
	// Domain 1 reports "Unavailable" for mbm_local_bytes, so that sample is
	// absent.
	testcase := `# HELP node_resctrl_llc_occupancy_bytes Last level cache bytes occupied in this domain.
# TYPE node_resctrl_llc_occupancy_bytes gauge
node_resctrl_llc_occupancy_bytes{domain="0"} 4.4040192e+07
node_resctrl_llc_occupancy_bytes{domain="1"} 8.388608e+06
# HELP node_resctrl_memory_bandwidth_bytes_total Bytes moved between the last level cache and memory. scope=total includes remote sockets, scope=local is memory attached to this domain.
# TYPE node_resctrl_memory_bandwidth_bytes_total counter
node_resctrl_memory_bandwidth_bytes_total{domain="0",scope="local"} 2.10196273664e+11
node_resctrl_memory_bandwidth_bytes_total{domain="0",scope="total"} 3.15294410752e+11
node_resctrl_memory_bandwidth_bytes_total{domain="1",scope="total"} 1.05098136832e+11
`

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := newTestResctrlCollector("fixtures/resctrl", logger)
	if err != nil {
		t.Fatal(err)
	}

	if err := testutil.CollectAndCompare(c, strings.NewReader(testcase)); err != nil {
		t.Fatal(err)
	}
}

func TestResctrlNotMounted(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := flagConfig()
	cfg.Resctrl.Path = t.TempDir()
	c, err := NewResctrlCollector(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric, 16)
	err = c.Update(ch)
	close(ch)

	if err != ErrNoData {
		t.Fatalf("expected ErrNoData, got %v", err)
	}
}

func TestDefaultResctrlPathMatchesResctrlfs(t *testing.T) {
	if config.DefaultResctrlPath != resctrlfs.DefaultMountPoint {
		t.Fatalf("config.DefaultResctrlPath is %q, resctrlfs.DefaultMountPoint is %q", config.DefaultResctrlPath, resctrlfs.DefaultMountPoint)
	}
}
