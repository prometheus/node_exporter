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

//go:build linux && !nocpu

package collector

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/prometheus/node_exporter/config"
)

func TestRuntimeReadsCPUConfigFromConfig(t *testing.T) {
	// The cpu flags hold the opposite of the non-default cases, so a runtime
	// that reads them fails the test.
	parseFlags(t, "--path.procfs=fixtures/proc", "--path.sysfs=fixtures/sys", "--collector.cpu.guest", "--no-collector.cpu.info")

	for _, test := range []struct {
		name      string
		cpu       config.CPUConfig
		wantGuest bool
		wantInfo  bool
	}{
		{name: "defaults", cpu: config.NewConfigWithDefaults().CPU, wantGuest: true},
		{name: "guest off, info on", cpu: config.CPUConfig{EnableInfo: true}, wantInfo: true},
		{name: "info enabled by an include filter", cpu: config.CPUConfig{BugsInclude: "^$"}, wantInfo: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.NewConfigWithDefaults()
			cfg.CollectorDisableDefaults = true
			cfg.Collectors = map[string]bool{"cpu": true}
			cfg.CPU = test.cpu
			runtime, err := NewRuntime(cfg, testLogger())
			if err != nil {
				t.Fatal(err)
			}
			registry, err := runtime.Registry()
			if err != nil {
				t.Fatal(err)
			}
			guest, err := testutil.GatherAndCount(registry, "node_cpu_guest_seconds_total")
			if err != nil {
				t.Fatal(err)
			}
			info, err := testutil.GatherAndCount(registry, "node_cpu_info")
			if err != nil {
				t.Fatal(err)
			}
			if got := guest > 0; got != test.wantGuest {
				t.Errorf("want guest metrics %t, got %d series", test.wantGuest, guest)
			}
			if got := info > 0; got != test.wantInfo {
				t.Errorf("want cpu info metrics %t, got %d series", test.wantInfo, info)
			}
		})
	}
}
