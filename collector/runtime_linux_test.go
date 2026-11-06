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
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/prometheus/node_exporter/config"
)

func TestRuntimeKeepsCPUConfigFromCreation(t *testing.T) {
	parseCollectorFlags(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	originalProc, originalSys := *procPath, *sysPath
	originalGuest, originalInfo := *enableCPUGuest, *enableCPUInfo
	originalFlags, originalBugs := *flagsInclude, *bugsInclude
	t.Cleanup(func() {
		*procPath, *sysPath = originalProc, originalSys
		*enableCPUGuest, *enableCPUInfo = originalGuest, originalInfo
		*flagsInclude, *bugsInclude = originalFlags, originalBugs
	})

	cfg := config.NewConfigWithDefaults()
	cfg.EnabledCollectors = []string{"cpu"}

	*procPath, *sysPath = "fixtures/proc", "fixtures/sys"
	*enableCPUGuest, *enableCPUInfo = false, false
	*flagsInclude, *bugsInclude = "", ""
	runtime, err := NewRuntime(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	*enableCPUGuest, *enableCPUInfo = true, true
	filtered, err := runtime.Filtered("cpu")
	if err != nil {
		t.Fatal(err)
	}
	later, err := NewRuntime(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name      string
		runtime   *Runtime
		wantGuest bool
		wantInfo  bool
	}{
		{name: "filtered runtime keeps the config set at creation", runtime: filtered, wantGuest: false, wantInfo: false},
		{name: "new runtime uses the current config", runtime: later, wantGuest: true, wantInfo: true},
	} {
		registry, err := test.runtime.Registry()
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
			t.Errorf("%s: want guest metrics %t, got %d series", test.name, test.wantGuest, guest)
		}
		if got := info > 0; got != test.wantInfo {
			t.Errorf("%s: want cpu info metrics %t, got %d series", test.name, test.wantInfo, info)
		}
	}
}
