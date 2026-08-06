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

//go:build !noconntrack

package collector

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// procWithoutConntrackStat builds a minimal procfs tree that has the conntrack
// sysctl entries but no net/stat/nf_conntrack, mimicking a kernel built
// without CONFIG_NF_CONNTRACK_PROCFS.
func procWithoutConntrackStat(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	sysctl := filepath.Join(root, "sys", "net", "netfilter")
	if err := os.MkdirAll(sysctl, 0755); err != nil {
		t.Fatal(err)
	}
	for name, value := range map[string]string{
		"nf_conntrack_count": "123\n",
		"nf_conntrack_max":   "65536\n",
	} {
		if err := os.WriteFile(filepath.Join(sysctl, name), []byte(value), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// net/stat exists but holds no nf_conntrack file.
	if err := os.MkdirAll(filepath.Join(root, "net", "stat"), 0755); err != nil {
		t.Fatal(err)
	}

	return root
}

// collectConntrack runs one scrape and returns the fully qualified names of the
// metrics that were emitted, along with the error the collector returned.
func collectConntrack(t *testing.T, c Collector) ([]string, error) {
	t.Helper()

	ch := make(chan prometheus.Metric, 32)
	err := c.Update(ch)
	close(ch)

	var names []string
	for metric := range ch {
		desc := metric.Desc().String()
		_, rest, found := strings.Cut(desc, `fqName: "`)
		if !found {
			t.Fatalf("cannot parse metric description: %s", desc)
		}
		name, _, found := strings.Cut(rest, `"`)
		if !found {
			t.Fatalf("cannot parse metric description: %s", desc)
		}
		names = append(names, name)
	}

	return names, err
}

func newConntrackCollectorForTest(t *testing.T, logger *slog.Logger, args ...string) Collector {
	t.Helper()

	if _, err := kingpin.CommandLine.Parse(args); err != nil {
		t.Fatal(err)
	}
	c, err := NewConntrackCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	return c
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestConntrackStatsCollectedWhenProcfsAvailable(t *testing.T) {
	c := newConntrackCollectorForTest(t, discardLogger(), "--path.procfs", "fixtures/proc")

	names, err := collectConntrack(t, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := len(names), 10; got != want {
		t.Fatalf("expected %d metrics, got %d: %v", want, got, names)
	}
	if !slices.Contains(names, "node_nf_conntrack_stat_insert_failed") {
		t.Fatalf("expected stat metrics to be present, got: %v", names)
	}
}

func TestConntrackStatsMissingProcfsKeepsSysctlMetrics(t *testing.T) {
	c := newConntrackCollectorForTest(t, discardLogger(), "--path.procfs", procWithoutConntrackStat(t))

	names, err := collectConntrack(t, c)
	if err != ErrNoData {
		t.Fatalf("expected ErrNoData, got: %v", err)
	}
	want := []string{"node_nf_conntrack_entries", "node_nf_conntrack_entries_limit"}
	if len(names) != len(want) {
		t.Fatalf("expected only %v, got: %v", want, names)
	}
	for _, name := range want {
		if !slices.Contains(names, name) {
			t.Fatalf("expected %s to be present, got: %v", name, names)
		}
	}
}

func TestConntrackStatsDisabledByFlag(t *testing.T) {
	c := newConntrackCollectorForTest(t, discardLogger(),
		"--path.procfs", "fixtures/proc", "--no-collector.conntrack.stats")

	names, err := collectConntrack(t, c)
	if err != nil {
		t.Fatalf("expected no error when stats are disabled, got: %v", err)
	}
	if got, want := len(names), 2; got != want {
		t.Fatalf("expected %d metrics, got %d: %v", want, got, names)
	}
	for _, name := range names {
		if strings.HasPrefix(name, "node_nf_conntrack_stat_") {
			t.Fatalf("stat metrics should not be collected when disabled, got: %v", names)
		}
	}
}

func TestConntrackStatsMissingProcfsWarnsOnce(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn}))
	c := newConntrackCollectorForTest(t, logger, "--path.procfs", procWithoutConntrackStat(t))

	for i := range 3 {
		if _, err := collectConntrack(t, c); err != ErrNoData {
			t.Fatalf("scrape %d: expected ErrNoData, got: %v", i, err)
		}
	}

	if got := strings.Count(buf.String(), "conntrack statistics unavailable"); got != 1 {
		t.Fatalf("expected exactly one warning across three scrapes, got %d:\n%s", got, buf.String())
	}
}
