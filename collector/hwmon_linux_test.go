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

//go:build !nohwmon

package collector

import (
	"io"
	"log/slog"
	"slices"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

// gatherChipLabels runs the collector through a registry and returns the
// observed `chip` label values across all metrics. It also surfaces any
// gather error — which is the failure mode #3637 reported.
func gatherChipLabels(t *testing.T, c *hwMonCollector) ([]string, error) {
	t.Helper()
	reg := prometheus.NewRegistry()
	if err := reg.Register(testHwmonCollector{c: c}); err != nil {
		t.Fatalf("register: %v", err)
	}
	families, err := reg.Gather()
	if err != nil {
		return nil, err
	}
	var chips []string
	for _, fam := range families {
		for _, m := range fam.Metric {
			for _, l := range m.Label {
				if l.GetName() == "chip" {
					chips = append(chips, l.GetValue())
				}
			}
		}
	}
	return chips, nil
}

// testHwmonCollector adapts hwMonCollector.Update for prometheus.Registry
// so we can exercise duplicate detection at the gather step.
type testHwmonCollector struct {
	c *hwMonCollector
}

func (t testHwmonCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(t, ch)
}

func (t testHwmonCollector) Collect(ch chan<- prometheus.Metric) {
	if err := t.c.Update(ch); err != nil {
		panic(err)
	}
}

func newTestHwmonCollector() *hwMonCollector {
	return &hwMonCollector{logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
}

// Two hwmon entries sharing the same parent device — the configuration
// that triggered #3637 on ASUS WMI laptops — must produce distinct chip
// labels and not error during gather.
func TestHwmonDuplicateChipNamesAreDisambiguated(t *testing.T) {
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = "fixtures/sys/hwmon/duplicate_chip_names"

	chips, err := gatherChipLabels(t, newTestHwmonCollector())
	if err != nil {
		t.Fatalf("gather: %v", err)
	}

	if !slices.Contains(chips, "platform_asus_nb_wmi_asus") {
		t.Errorf("expected disambiguated chip 'platform_asus_nb_wmi_asus', got %v", uniq(chips))
	}
	if !slices.Contains(chips, "platform_asus_nb_wmi_asus_wmi_sensors") {
		t.Errorf("expected disambiguated chip 'platform_asus_nb_wmi_asus_wmi_sensors', got %v", uniq(chips))
	}
	for _, chip := range chips {
		if chip == "platform_asus_nb_wmi" {
			t.Errorf("colliding chip should not be emitted with bare base name, got %q", chip)
		}
	}
}

// When chip names are already unique, the collector must leave them alone
// — no surprise suffixes for unaffected users.
func TestHwmonUniqueChipNamesAreUnchanged(t *testing.T) {
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = "fixtures/sys/hwmon/unique_chip_names"

	chips, err := gatherChipLabels(t, newTestHwmonCollector())
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	if !slices.Contains(chips, "platform_coretemp_0") {
		t.Errorf("expected platform_coretemp_0, got %v", uniq(chips))
	}
	if !slices.Contains(chips, "platform_coretemp_1") {
		t.Errorf("expected platform_coretemp_1, got %v", uniq(chips))
	}
}

// When colliding entries share the same `name` file content, the `name`
// disambiguator collapses too. We must still emit unique chip labels by
// falling back to the hwmonX basename.
func TestHwmonDuplicateChipNamesWithSameNameFile(t *testing.T) {
	prev := *sysPath
	t.Cleanup(func() { *sysPath = prev })
	*sysPath = "fixtures/sys/hwmon/same_name_file"

	chips, err := gatherChipLabels(t, newTestHwmonCollector())
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	if !slices.Contains(chips, "platform_asus_nb_wmi_hwmon3") {
		t.Errorf("expected platform_asus_nb_wmi_hwmon3, got %v", uniq(chips))
	}
	if !slices.Contains(chips, "platform_asus_nb_wmi_hwmon4") {
		t.Errorf("expected platform_asus_nb_wmi_hwmon4, got %v", uniq(chips))
	}
}

func uniq(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
