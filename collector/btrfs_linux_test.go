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

//go:build !nobtrfs
// +build !nobtrfs

package collector

import (
	"strings"
	"testing"

	"github.com/prometheus/procfs/btrfs"
)

var expectedBtrfsMetrics = [][]btrfsMetric{
	{
		{name: "info", value: 1, extraLabel: []string{"label"}, extraLabelValue: []string{"fixture"}},
		{name: "global_rsv_size_bytes", value: 1.6777216e+07},
		{name: "reserved_bytes", value: 0, extraLabel: []string{"block_group_type"}, extraLabelValue: []string{"data"}},
		{name: "used_bytes", value: 8.08189952e+08, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"data", "raid0"}},
		{name: "size_bytes", value: 2.147483648e+09, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"data", "raid0"}},
		{name: "allocation_ratio", value: 1, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"data", "raid0"}},
		{name: "reserved_bytes", value: 0, extraLabel: []string{"block_group_type"}, extraLabelValue: []string{"metadata"}},
		{name: "used_bytes", value: 933888, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"metadata", "raid1"}},
		{name: "size_bytes", value: 1.073741824e+09, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"metadata", "raid1"}},
		{name: "allocation_ratio", value: 2, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"metadata", "raid1"}},
		{name: "reserved_bytes", value: 0, extraLabel: []string{"block_group_type"}, extraLabelValue: []string{"system"}},
		{name: "used_bytes", value: 16384, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"system", "raid1"}},
		{name: "size_bytes", value: 8.388608e+06, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"system", "raid1"}},
		{name: "allocation_ratio", value: 2, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"system", "raid1"}},
		{name: "device_size_bytes", value: 1.073741824e+10, extraLabel: []string{"device"}, extraLabelValue: []string{"loop25"}},
		{name: "device_size_bytes", value: 1.073741824e+10, extraLabel: []string{"device"}, extraLabelValue: []string{"loop26"}},
	},
	{
		{name: "info", value: 1, extraLabel: []string{"label"}, extraLabelValue: []string{""}},
		{name: "global_rsv_size_bytes", value: 1.6777216e+07},
		{name: "reserved_bytes", value: 0, extraLabel: []string{"block_group_type"}, extraLabelValue: []string{"data"}},
		{name: "used_bytes", value: 0, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"data", "raid5"}},
		{name: "size_bytes", value: 6.44087808e+08, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"data", "raid5"}},
		{name: "allocation_ratio", value: 1.3333333333333333, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"data", "raid5"}},
		{name: "reserved_bytes", value: 0, extraLabel: []string{"block_group_type"}, extraLabelValue: []string{"metadata"}},
		{name: "used_bytes", value: 114688, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"metadata", "raid6"}},
		{name: "size_bytes", value: 4.29391872e+08, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"metadata", "raid6"}},
		{name: "allocation_ratio", value: 2, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"metadata", "raid6"}},
		{name: "reserved_bytes", value: 0, extraLabel: []string{"block_group_type"}, extraLabelValue: []string{"system"}},
		{name: "used_bytes", value: 16384, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"system", "raid6"}},
		{name: "size_bytes", value: 1.6777216e+07, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"system", "raid6"}},
		{name: "allocation_ratio", value: 2, extraLabel: []string{"block_group_type", "mode"}, extraLabelValue: []string{"system", "raid6"}},
		{name: "device_size_bytes", value: 1.073741824e+10, extraLabel: []string{"device"}, extraLabelValue: []string{"loop22"}},
		{name: "device_size_bytes", value: 1.073741824e+10, extraLabel: []string{"device"}, extraLabelValue: []string{"loop23"}},
		{name: "device_size_bytes", value: 1.073741824e+10, extraLabel: []string{"device"}, extraLabelValue: []string{"loop24"}},
		{name: "device_size_bytes", value: 1.073741824e+10, extraLabel: []string{"device"}, extraLabelValue: []string{"loop25"}},
	},
}

func checkMetric(exp, got *btrfsMetric) bool {
	if exp.name != got.name ||
		exp.value != got.value ||
		len(exp.extraLabel) != len(got.extraLabel) ||
		len(exp.extraLabelValue) != len(got.extraLabelValue) {
		return false
	}

	for i := range exp.extraLabel {
		if exp.extraLabel[i] != got.extraLabel[i] {
			return false
		}

		// Devices (loopXX) can appear in random order, so just check the first 4 characters.
		if strings.HasPrefix(got.extraLabelValue[i], "loop") &&
			exp.extraLabelValue[i][:4] == got.extraLabelValue[i][:4] {
			continue
		}

		if exp.extraLabelValue[i] != got.extraLabelValue[i] {
			return false
		}
	}

	return true
}

func TestBtrfs(t *testing.T) {
	fs, _ := btrfs.NewFS("fixtures/sys")
	collector := &btrfsCollector{fs: fs}

	stats, err := collector.fs.Stats()
	if err != nil {
		t.Fatalf("Failed to retrieve Btrfs stats: %v", err)
	}
	if len(stats) != len(expectedBtrfsMetrics) {
		t.Fatalf("Unexpected number of Btrfs stats: expected %v, got %v", len(expectedBtrfsMetrics), len(stats))
	}

	for i, s := range stats {
		metrics := collector.getMetrics(s, nil)
		if len(metrics) != len(expectedBtrfsMetrics[i]) {
			t.Fatalf("Unexpected number of Btrfs metrics: expected %v, got %v", len(expectedBtrfsMetrics[i]), len(metrics))
		}

		for j, m := range metrics {
			exp := expectedBtrfsMetrics[i][j]
			if !checkMetric(&exp, &m) {
				t.Errorf("Incorrect btrfs metric: expected %#v, got: %#v", exp, m)
			}
		}
	}
}
