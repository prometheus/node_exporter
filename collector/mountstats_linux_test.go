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

//go:build !nomountstats

package collector

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testMountStatsCollector struct {
	c Collector
}

func (c testMountStatsCollector) Collect(ch chan<- prometheus.Metric) {
	c.c.Update(ch)
}

func (c testMountStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func NewTestMountStatsCollector(logger *slog.Logger) (prometheus.Collector, error) {
	c, err := NewMountStatsCollector(logger)
	if err != nil {
		return testMountStatsCollector{}, err
	}
	return testMountStatsCollector{c: c}, nil
}

func TestMountStatsMountPointLabel(t *testing.T) {
	*procPath = "fixtures/proc"

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewTestMountStatsCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(c)

	expected := `# HELP node_mountstats_nfs_age_seconds_total The age of the NFS mount in seconds.
# TYPE node_mountstats_nfs_age_seconds_total counter
node_mountstats_nfs_age_seconds_total{export="192.168.1.1:/srv/test",mountaddr="192.168.1.1",mountpoint="/mnt/nfs/test",protocol="tcp"} 13968
node_mountstats_nfs_age_seconds_total{export="192.168.1.1:/srv/test",mountaddr="192.168.1.1",mountpoint="/mnt/nfs/test-dupe",protocol="tcp"} 13968
node_mountstats_nfs_age_seconds_total{export="192.168.1.1:/srv/test",mountaddr="192.168.1.1",mountpoint="/mnt/nfs/test-dupe",protocol="udp"} 13968
`

	if err := testutil.GatherAndCompare(reg, strings.NewReader(expected), "node_mountstats_nfs_age_seconds_total"); err != nil {
		t.Fatal(err)
	}
}
