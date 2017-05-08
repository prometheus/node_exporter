// Copyright 2017 The Prometheus Authors
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
	"reflect"
	"testing"
)

func TestBcacheStats(t *testing.T) {
	// test dehumanize
	dehumanizeTests := []struct {
		in  []byte
		out float64
	}{
		{
			in:  []byte("542k"),
			out: float64(555008),
		},
		{
			in:  []byte("322M"),
			out: float64(337641472),
		},
	}
	for _, tst := range dehumanizeTests {
		got := dehumanize(tst.in)
		if got != tst.out {
			t.Errorf("want dehumanize %s, want %f got %f", tst.in, tst.out, got)
		}
	}

	// test priorityStats
	priorityStatsTests := []struct {
		in        string
		out_key   string
		out_value float64
		err       error
	}{
		{
			in:        "Unused:         99%",
			out_key:   "priority_stats_unused_percent",
			out_value: float64(99),
			err:       nil,
		},
		{
			in:        "Metadata:       0%",
			out_key:   "priority_stats_metadata_percent",
			out_value: float64(0),
			err:       nil,
		},
	}
	for _, tst := range priorityStatsTests {
		got_key, got_value, got_err := parsePriorityStats(tst.in)
		if got_key != tst.out_key || got_value != tst.out_value || got_err != tst.err {
			t.Errorf("want parsePriorityStats %s, got %s %f %v", tst.in, got_key, got_value, got_err)
		}
	}

	// test getCacheStats
	want := map[string]float64{
		"metadata_written":                float64(512),
		"priority_stats_unused_percent":   float64(99),
		"priority_stats_metadata_percent": float64(0),
		"written":                         float64(0),
		"io_errors":                       float64(0),
	}
	got, _ := getCacheStats("fixtures/sys/fs/bcache/deaddd54-c735-46d5-868e-f331c5fd7c74/cache0")
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want parsePriorityStats want %+v, got %v", want, got)
	}
}
