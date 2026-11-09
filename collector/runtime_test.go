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

package collector

import (
	"io"
	"log/slog"
	"maps"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/prometheus/node_exporter/config"
)

// The tests use textfile because it is enabled by default and builds on every
// OS. Collectors such as filesystem leave goroutines that race with later
// tests that parse flags.

func TestNewRuntimeEnabledCollectors(t *testing.T) {
	parseCollectorFlags(t)
	tests := []struct {
		name            string
		disableDefaults bool
		wantTextfile    bool
	}{
		{name: "defaults", wantTextfile: true},
		{name: "disable defaults", disableDefaults: true, wantTextfile: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.NewConfigWithDefaults()
			cfg.CollectorDisableDefaults = test.disableDefaults

			runtime, err := NewRuntime(cfg, testLogger())
			if err != nil {
				t.Fatal(err)
			}
			if got := slices.Contains(runtime.EnabledCollectors(), "textfile"); got != test.wantTextfile {
				t.Fatalf("textfile enabled: want %t, got %t (%v)", test.wantTextfile, got, runtime.EnabledCollectors())
			}
		})
	}
}

func TestNewRuntimeCollectorsOverrideDefaults(t *testing.T) {
	parseCollectorFlags(t)
	tests := []struct {
		name            string
		disableDefaults bool
		collectors      map[string]bool
		wantTextfile    bool
	}{
		{
			name:            "enable with defaults disabled",
			disableDefaults: true,
			collectors:      map[string]bool{"textfile": true},
			wantTextfile:    true,
		},
		{
			name:         "disable a default collector",
			collectors:   map[string]bool{"textfile": false},
			wantTextfile: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.NewConfigWithDefaults()
			cfg.CollectorDisableDefaults = test.disableDefaults
			cfg.Collectors = test.collectors

			runtime, err := NewRuntime(cfg, testLogger())
			if err != nil {
				t.Fatal(err)
			}
			if got := slices.Contains(runtime.EnabledCollectors(), "textfile"); got != test.wantTextfile {
				t.Fatalf("textfile enabled: want %t, got %t (%v)", test.wantTextfile, got, runtime.EnabledCollectors())
			}
		})
	}
}

func TestNewRuntimeCollectorsTakePrecedenceOverFlags(t *testing.T) {
	parseCollectorFlags(t)
	original := *collectorState["textfile"]
	*collectorState["textfile"] = false
	forcedCollectors["textfile"] = true
	t.Cleanup(func() {
		*collectorState["textfile"] = original
		delete(forcedCollectors, "textfile")
	})

	runtime, err := NewRuntime(textfileOnlyConfig(), testLogger())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := runtime.EnabledCollectors(), []string{"textfile"}; !slices.Equal(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestNewRuntimeRejectsUnknownCollector(t *testing.T) {
	cfg := config.NewConfigWithDefaults()
	cfg.Collectors = map[string]bool{"doesnotexist": true}

	if _, err := NewRuntime(cfg, testLogger()); err == nil {
		t.Fatal("want error for an unknown collector, got nil")
	}
}

func TestNewRuntimeDisableDefaultsDoesNotAffectOtherRuntimes(t *testing.T) {
	parseCollectorFlags(t)
	cfg := config.NewConfigWithDefaults()
	cfg.CollectorDisableDefaults = true
	if _, err := NewRuntime(cfg, testLogger()); err != nil {
		t.Fatal(err)
	}

	runtime, err := NewRuntime(config.NewConfigWithDefaults(), testLogger())
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(runtime.EnabledCollectors(), "textfile") {
		t.Fatalf("want textfile enabled, got %v", runtime.EnabledCollectors())
	}
}

func TestNewRuntimeUsesDefaultsWithoutParsedFlags(t *testing.T) {
	// Before kingpin.Parse, every collector flag holds false.
	original := *collectorState["textfile"]
	*collectorState["textfile"] = false
	t.Cleanup(func() { *collectorState["textfile"] = original })

	runtime, err := NewRuntime(config.NewConfigWithDefaults(), testLogger())
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(runtime.EnabledCollectors(), "textfile") {
		t.Fatalf("want textfile enabled, got %v", runtime.EnabledCollectors())
	}
}

func TestRuntimeExposesCollectorMetrics(t *testing.T) {
	parseCollectorFlags(t)
	runtime, err := NewRuntime(textfileOnlyConfig(), testLogger())
	if err != nil {
		t.Fatal(err)
	}

	collectorsRegistry := prometheus.NewRegistry()
	collectorsRegistry.MustRegister(runtime.Collectors()...)
	runtimeRegistry, err := runtime.Registry()
	if err != nil {
		t.Fatal(err)
	}

	want := `# HELP node_scrape_collector_success node_exporter: Whether a collector succeeded.
# TYPE node_scrape_collector_success gauge
node_scrape_collector_success{collector="textfile"} 1
`
	for name, registry := range map[string]*prometheus.Registry{
		"Collectors": collectorsRegistry,
		"Registry":   runtimeRegistry,
	} {
		if err := testutil.GatherAndCompare(registry, strings.NewReader(want), "node_scrape_collector_success"); err != nil {
			t.Errorf("%s: %s", name, err)
		}
	}
}

func TestRuntimeFiltered(t *testing.T) {
	parseCollectorFlags(t)
	runtime, err := NewRuntime(textfileOnlyConfig(), testLogger())
	if err != nil {
		t.Fatal(err)
	}

	filtered, err := runtime.Filtered("textfile")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := filtered.EnabledCollectors(), []string{"textfile"}; !slices.Equal(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}

	// meminfo is enabled by default, but not in this runtime.
	if _, err := runtime.Filtered("meminfo"); err == nil {
		t.Fatal("want error for a collector outside the runtime, got nil")
	}
}

func TestDeprecatedNodeCollectorWrappers(t *testing.T) {
	parseCollectorFlags(t)
	original := make(map[string]bool, len(collectorState))
	for c, state := range collectorState {
		original[c] = *state
	}
	t.Cleanup(func() {
		for c, state := range collectorState {
			*state = original[c]
		}
		delete(forcedCollectors, "textfile")
	})

	*collectorState["textfile"] = true
	forcedCollectors["textfile"] = true
	DisableDefaultCollectors()

	first, err := NewNodeCollector(testLogger())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := slices.Sorted(maps.Keys(first.Collectors)), []string{"textfile"}; !slices.Equal(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}

	second, err := NewNodeCollector(testLogger())
	if err != nil {
		t.Fatal(err)
	}
	if first.Collectors["textfile"] != second.Collectors["textfile"] {
		t.Fatal("want callers to share collector instances")
	}
}

// textfileOnlyConfig returns a Config that enables only the textfile collector.
func textfileOnlyConfig() config.Config {
	cfg := config.NewConfigWithDefaults()
	cfg.CollectorDisableDefaults = true
	cfg.Collectors = map[string]bool{"textfile": true}
	return cfg
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

var parseCollectorFlagsOnce sync.Once

// parseCollectorFlags applies the kingpin defaults once, so the collector
// flags hold their default values when a test runs on its own.
func parseCollectorFlags(t *testing.T) {
	t.Helper()
	var err error
	parseCollectorFlagsOnce.Do(func() {
		_, err = kingpin.CommandLine.Parse([]string{})
	})
	if err != nil {
		t.Fatal(err)
	}
}
