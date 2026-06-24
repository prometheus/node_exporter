// Copyright 2019 The Prometheus Authors
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

//go:build !noext4

package collector

import (
	"testing"

	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/ext4"
)

var expectedExt4Metrics = [][]ext4Metric{
	{
		{name: "errors_total", value: 12},
		{name: "warnings_total", value: 34},
		{name: "messages_total", value: 567},
	},
}

func checkExt4Metric(exp, got *ext4Metric) bool {
	if exp.name != got.name ||
		exp.value != got.value {
		return false
	}
	return true
}

func TestExt4(t *testing.T) {
	fs, err := ext4.NewFS(procfs.DefaultMountPoint, "fixtures/sys")
	if err != nil {
		t.Fatal(err)
	}
	collector := &ext4Collector{fs: fs}

	stats, err := collector.fs.ProcStat()
	if err != nil {
		t.Fatalf("Failed to retrieve ext4 stats: %v", err)
	}
	if len(stats) != len(expectedExt4Metrics) {
		t.Fatalf("Unexpected number of ext4 stats: expected %v, got %v", len(expectedExt4Metrics), len(stats))
	}

	for i, s := range stats {
		metrics := collector.getMetrics(s)
		if len(metrics) != len(expectedExt4Metrics[i]) {
			t.Fatalf("Unexpected number of ext4 metrics: expected %v, got %v", len(expectedExt4Metrics[i]), len(metrics))
		}

		for j, m := range metrics {
			exp := expectedExt4Metrics[i][j]
			if !checkExt4Metric(&exp, &m) {
				t.Errorf("Incorrect ext4 metric: expected %#v, got: %#v", exp, m)
			}
		}
	}
}
