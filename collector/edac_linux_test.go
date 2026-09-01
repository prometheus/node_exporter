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

//go:build !noedac

package collector

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
)

// newTestEdacRegistry points the collector at the sysfs fixture tree and returns
// a registry it has been registered against.
func newTestEdacRegistry(t *testing.T) *prometheus.Registry {
	t.Helper()

	if _, err := kingpin.CommandLine.Parse([]string{
		"--path.sysfs", "fixtures/sys",
		"--collector.edac",
	}); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewEdacCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(&testEdacCollector{ec: c})
	return reg
}

// Fixture mc1 registers dimm0, dimm1 and dimm3 to cover sparse indices; mc2 uses
// the rank* spelling and omits dimm_ue_count on rank1.
func TestEdacDimmLayer(t *testing.T) {
	reg := newTestEdacRegistry(t)

	expected := `
# HELP node_edac_dimm_correctable_errors_total Total correctable memory errors for this DIMM.
# TYPE node_edac_dimm_correctable_errors_total counter
node_edac_dimm_correctable_errors_total{controller="1",dimm="0",dimm_label="CPU_SrcID#0_MC#1_Chan#0_DIMM#0"} 100
node_edac_dimm_correctable_errors_total{controller="1",dimm="1",dimm_label="CPU_SrcID#0_MC#1_Chan#1_DIMM#0"} 200
node_edac_dimm_correctable_errors_total{controller="1",dimm="3",dimm_label="CPU_SrcID#0_MC#1_Chan#3_DIMM#0"} 300
node_edac_dimm_correctable_errors_total{controller="2",dimm="0",dimm_label="PROC 1 DIMM 8"} 7
node_edac_dimm_correctable_errors_total{controller="2",dimm="1",dimm_label="PROC 2 DIMM 10"} 9
# HELP node_edac_dimm_uncorrectable_errors_total Total uncorrectable memory errors for this DIMM. Best effort: errors the controller cannot localize to a slot are counted in ue_noinfo_count instead.
# TYPE node_edac_dimm_uncorrectable_errors_total counter
node_edac_dimm_uncorrectable_errors_total{controller="1",dimm="0",dimm_label="CPU_SrcID#0_MC#1_Chan#0_DIMM#0"} 0
node_edac_dimm_uncorrectable_errors_total{controller="1",dimm="1",dimm_label="CPU_SrcID#0_MC#1_Chan#1_DIMM#0"} 0
node_edac_dimm_uncorrectable_errors_total{controller="1",dimm="3",dimm_label="CPU_SrcID#0_MC#1_Chan#3_DIMM#0"} 1
node_edac_dimm_uncorrectable_errors_total{controller="2",dimm="0",dimm_label="PROC 1 DIMM 8"} 0
`

	if err := testutil.GatherAndCompare(reg, strings.NewReader(expected),
		"node_edac_dimm_correctable_errors_total",
		"node_edac_dimm_uncorrectable_errors_total",
	); err != nil {
		t.Fatal(err)
	}
}

// Fixture mc3 models amd64_edac_mod, where rank* repeats the csrow counts. Its
// rank* dirs hold 999 so double-counting shows up loudly.
func TestEdacCsrowTakesPriorityOverDimmLayer(t *testing.T) {
	reg := newTestEdacRegistry(t)

	families, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}

	for _, family := range families {
		if !strings.HasPrefix(family.GetName(), "node_edac_dimm_") {
			continue
		}
		for _, metric := range family.GetMetric() {
			for _, label := range metric.GetLabel() {
				if label.GetName() == "controller" && label.GetValue() == "3" {
					t.Errorf("%s reported for controller 3, which exposes csrow* alongside rank*: "+
						"the DIMM layer must be skipped there or counts are doubled (value %v, dimm_label %q)",
						family.GetName(), metric.GetCounter().GetValue(), edacLabelValue(metric, "dimm_label"))
				}
			}
		}
	}
}

func edacLabelValue(metric *dto.Metric, name string) string {
	for _, label := range metric.GetLabel() {
		if label.GetName() == name {
			return label.GetValue()
		}
	}
	return ""
}

// testEdacCollector adapts the Collector interface for use with a registry.
type testEdacCollector struct {
	ec Collector
}

func (tc *testEdacCollector) Collect(ch chan<- prometheus.Metric) {
	sink := make(chan prometheus.Metric)
	go func() {
		if err := tc.ec.Update(sink); err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	for m := range sink {
		ch <- m
	}
}

func (tc *testEdacCollector) Describe(_ chan<- *prometheus.Desc) {
	// No-op for testing.
}
