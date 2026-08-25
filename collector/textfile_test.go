// Copyright 2015 The Prometheus Authors
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

//go:build !notextfile

package collector

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
)

type collectorAdapter struct {
	Collector
}

// Describe implements the prometheus.Collector interface.
func (a collectorAdapter) Describe(ch chan<- *prometheus.Desc) {
	// We have to send *some* metric in Describe, but we don't know which ones
	// we're going to get, so just send a dummy metric.
	ch <- prometheus.NewDesc("dummy_metric", "Dummy metric.", nil, nil)
}

// Collect implements the prometheus.Collector interface.
func (a collectorAdapter) Collect(ch chan<- prometheus.Metric) {
	if err := a.Update(ch); err != nil {
		panic(fmt.Sprintf("failed to update collector: %v", err))
	}
}

func TestTextfileCollector(t *testing.T) {
	tests := []struct {
		paths []string
		out   string
	}{
		{
			paths: []string{"fixtures/textfile/no_metric_files"},
			out:   "fixtures/textfile/no_metric_files.out",
		},
		{
			paths: []string{"fixtures/textfile/two_metric_files"},
			out:   "fixtures/textfile/two_metric_files.out",
		},
		{
			paths: []string{"fixtures/textfile/nonexistent_path"},
			out:   "fixtures/textfile/nonexistent_path.out",
		},
		{
			paths: []string{"fixtures/textfile/client_side_timestamp"},
			out:   "fixtures/textfile/client_side_timestamp.out",
		},
		{
			paths: []string{"fixtures/textfile/different_metric_types"},
			out:   "fixtures/textfile/different_metric_types.out",
		},
		{
			paths: []string{"fixtures/textfile/inconsistent_metrics"},
			out:   "fixtures/textfile/inconsistent_metrics.out",
		},
		{
			paths: []string{"fixtures/textfile/histogram"},
			out:   "fixtures/textfile/histogram.out",
		},
		{
			paths: []string{"fixtures/textfile/histogram_extra_dimension"},
			out:   "fixtures/textfile/histogram_extra_dimension.out",
		},
		{
			paths: []string{"fixtures/textfile/summary"},
			out:   "fixtures/textfile/summary.out",
		},
		{
			paths: []string{"fixtures/textfile/summary_extra_dimension"},
			out:   "fixtures/textfile/summary_extra_dimension.out",
		},
		{
			paths: []string{
				"fixtures/textfile/histogram_extra_dimension",
				"fixtures/textfile/summary_extra_dimension",
			},
			out: "fixtures/textfile/glob_extra_dimension.out",
		},
		{
			paths: []string{"fixtures/textfile/*_extra_dimension"},
			out:   "fixtures/textfile/glob_extra_dimension.out",
		},
		{
			paths: []string{"fixtures/textfile/metrics_merge_empty_help"},
			out:   "fixtures/textfile/metrics_merge_empty_help.out",
		},
		{
			paths: []string{"fixtures/textfile/metrics_merge_no_help"},
			out:   "fixtures/textfile/metrics_merge_no_help.out",
		},
		{
			paths: []string{"fixtures/textfile/metrics_merge_same_help"},
			out:   "fixtures/textfile/metrics_merge_same_help.out",
		},
		{
			paths: []string{"fixtures/textfile/metrics_merge_different_help"},
			out:   "fixtures/textfile/metrics_merge_different_help.out",
		},
	}

	for i, test := range tests {
		mtime := 1.0
		c := &textFileCollector{
			paths:  test.paths,
			mtime:  &mtime,
			logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		}

		// Suppress a log message about `nonexistent_path` not existing, this is
		// expected and clutters the test output.
		promslogConfig := &promslog.Config{}
		flag.AddFlags(kingpin.CommandLine, promslogConfig)
		if _, err := kingpin.CommandLine.Parse([]string{"--log.level", "debug"}); err != nil {
			t.Fatal(err)
		}

		registry := prometheus.NewRegistry()
		registry.MustRegister(collectorAdapter{c})

		rw := httptest.NewRecorder()
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(rw, &http.Request{})
		got := string(rw.Body.String())

		want, err := os.ReadFile(test.out)
		if err != nil {
			t.Fatalf("%d. error reading fixture file %s: %s", i, test.out, err)
		}

		if string(want) != got {
			t.Fatalf("%d.%q want:\n\n%s\n\ngot:\n\n%s", i, test.paths, string(want), got)
		}
	}
}

func TestTextfileCollectorSupportsUTF8Names(t *testing.T) {
	dir := t.TempDir()
	const input = `# HELP "my.dotted.metric" A metric with UTF-8 compatible names.
# TYPE "my.dotted.metric" gauge
{"my.dotted.metric", "error.message"="Not Found"} 1
`
	if err := os.WriteFile(filepath.Join(dir, "test.prom"), []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}

	c := &textFileCollector{
		paths:  []string{dir},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectorAdapter{c})

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("failed to gather textfile with UTF-8 names: %v", err)
	}

	for _, family := range families {
		if family.GetName() != "my.dotted.metric" {
			continue
		}
		if got := family.Metric[0].Label[0].GetName(); got != "error.message" {
			t.Fatalf("expected dotted label name, got %q", got)
		}
		return
	}
	t.Fatalf("expected dotted metric name, got %v", families)
}
