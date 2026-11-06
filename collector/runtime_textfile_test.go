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

//go:build !notextfile

package collector

import (
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/prometheus/node_exporter/config"
)

func TestRuntimeKeepsTextfileDirectoriesFromCreation(t *testing.T) {
	parseCollectorFlags(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	original := *textFileDirectories
	t.Cleanup(func() { *textFileDirectories = original })

	cfg := config.NewConfigWithDefaults()
	cfg.EnabledCollectors = []string{"textfile"}

	*textFileDirectories = []string{"fixtures/textfile/two_metric_files"}
	runtime, err := NewRuntime(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	*textFileDirectories = []string{"fixtures/textfile/no_metric_files"}
	filtered, err := runtime.Filtered("textfile")
	if err != nil {
		t.Fatal(err)
	}
	later, err := NewRuntime(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name    string
		runtime *Runtime
		want    int
	}{
		{name: "filtered runtime reads the directory set at creation", runtime: filtered, want: 1},
		{name: "new runtime reads the current directory", runtime: later, want: 0},
	} {
		registry, err := test.runtime.Registry()
		if err != nil {
			t.Fatal(err)
		}
		got, err := testutil.GatherAndCount(registry, "testmetric1_1")
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Errorf("%s: want %d testmetric1_1 series, got %d", test.name, test.want, got)
		}
	}
}
