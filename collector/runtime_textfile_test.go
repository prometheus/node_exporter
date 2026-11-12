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
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRuntimeReadsTextfileDirectoriesFromConfig(t *testing.T) {
	for _, test := range []struct {
		name      string
		directory string
		want      int
	}{
		{name: "directory with metrics", directory: "fixtures/textfile/two_metric_files", want: 1},
		{name: "directory without metrics", directory: "fixtures/textfile/no_metric_files", want: 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg := textfileOnlyConfig()
			cfg.Textfile.Directories = []string{test.directory}
			runtime, err := NewRuntime(cfg, testLogger())
			if err != nil {
				t.Fatal(err)
			}
			registry, err := runtime.Registry()
			if err != nil {
				t.Fatal(err)
			}
			got, err := testutil.GatherAndCount(registry, "testmetric1_1")
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("want %d testmetric1_1 series, got %d", test.want, got)
			}
		})
	}
}
