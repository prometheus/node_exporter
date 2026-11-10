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

//go:build !noksmd

package collector

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testKsmdCollector struct {
	kc Collector
}

func (c testKsmdCollector) Collect(ch chan<- prometheus.Metric) {
	c.kc.Update(ch)
}

func (c testKsmdCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// TestKsmdMetrics reads the fixture, whose general_profit is negative (the
// kernel prints it as a signed value), and checks that it does not prevent the
// other ksmd metrics from being collected.
func TestKsmdMetrics(t *testing.T) {
	*sysPath = "fixtures/sys"

	kc, err := NewKsmdCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	// Update returning an error means collection stopped early.
	if err := kc.Update(make(chan prometheus.Metric, 32)); err != nil {
		t.Fatalf("Update() returned an error: %v", err)
	}

	expected := `# HELP node_ksmd_full_scans_total How many times all mergeable areas have been scanned.
# TYPE node_ksmd_full_scans_total counter
node_ksmd_full_scans_total 323
# HELP node_ksmd_general_profit_bytes How effective is KSM. ksm_saved_pages * sizeof(page) - (all_rmap_items) * sizeof(rmap_item)
# TYPE node_ksmd_general_profit_bytes gauge
node_ksmd_general_profit_bytes -81920
# HELP node_ksmd_pages_sharing How many more sites are sharing them i.e. how much saved.
# TYPE node_ksmd_pages_sharing gauge
node_ksmd_pages_sharing 255
# HELP node_ksmd_run The status of ksmd. stopped = 0, running = 1
# TYPE node_ksmd_run gauge
node_ksmd_run 1
`
	reg := prometheus.NewRegistry()
	reg.MustRegister(testKsmdCollector{kc})

	err = testutil.GatherAndCompare(reg, strings.NewReader(expected),
		"node_ksmd_full_scans_total",
		"node_ksmd_general_profit_bytes",
		"node_ksmd_pages_sharing",
		"node_ksmd_run",
	)
	if err != nil {
		t.Fatal(err)
	}
}
