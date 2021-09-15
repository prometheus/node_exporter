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

package collector

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func TestDiskStats(t *testing.T) {
	var diskLabelNames = []string{"device"}
	*sysPath = "fixtures/sys"
	*procPath = "fixtures/proc"
	*ignoredDevices = "^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$"

	testcases := []string{
		prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "info"),
			"Info of /sys/block/<block_device>.",
			[]string{"device", "major", "minor"},
			nil,
		).String(),
		readsCompletedDesc.String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "reads_merged_total"),
			"The total number of reads merged.",
			diskLabelNames,
			nil,
		).String(),
		readBytesDesc.String(),
		readTimeSecondsDesc.String(),
		writesCompletedDesc.String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "writes_merged_total"),
			"The number of writes merged.",
			diskLabelNames,
			nil,
		).String(),
		writtenBytesDesc.String(),
		writeTimeSecondsDesc.String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "io_now"),
			"The number of I/Os currently in progress.",
			diskLabelNames,
			nil,
		).String(),
		ioTimeSecondsDesc.String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "io_time_weighted_seconds_total"),
			"The weighted # of seconds spent doing I/Os.",
			diskLabelNames,
			nil,
		).String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "discards_completed_total"),
			"The total number of discards completed successfully.",
			diskLabelNames,
			nil,
		).String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "discards_merged_total"),
			"The total number of discards merged.",
			diskLabelNames,
			nil,
		).String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "discarded_sectors_total"),
			"The total number of sectors discarded successfully.",
			diskLabelNames,
			nil,
		).String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "discard_time_seconds_total"),
			"This is the total number of seconds spent by all discards.",
			diskLabelNames,
			nil,
		).String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "flush_requests_total"),
			"The total number of flush requests completed successfully",
			diskLabelNames,
			nil,
		).String(),
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "flush_requests_time_seconds_total"),
			"This is the total number of seconds spent by all flush requests.",
			diskLabelNames,
			nil,
		).String(),
	}
	collector, err := NewDiskstatsCollector(log.NewLogfmtLogger(os.Stderr))
	if err != nil {
		panic(err)
	}

	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	for _, expected := range testcases {
		metric := (<-sink)
		if metric == nil {
			t.Fatalf("Expected '%s' but got nothing (nil).", expected)
		}

		got := metric.Desc().String()
		metric.Desc()
		if expected != got {
			t.Errorf("Expected '%s' but got '%s'", expected, got)
		} else {
			t.Logf("Successfully got '%s'", got)
		}
	}
}
