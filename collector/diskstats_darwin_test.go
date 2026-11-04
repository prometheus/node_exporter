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

//go:build !nodiskstats

package collector

import (
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/lufia/iostat"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type darwinDiskstatsTestCollector struct {
	dsc   *diskstatsCollector
	stats []*iostat.DriveStats
}

func (c darwinDiskstatsTestCollector) Collect(ch chan<- prometheus.Metric) {
	c.dsc.updateFromStats(ch, c.stats)
}

func (c darwinDiskstatsTestCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func TestDiskStatsDarwin(t *testing.T) {
	dsc, err := NewDiskstatsCollector(slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	c := darwinDiskstatsTestCollector{
		dsc: dsc.(*diskstatsCollector),
		stats: []*iostat.DriveStats{
			{
				Name:           "disk0",
				BlockSize:      4096,
				BytesRead:      16384,
				BytesWritten:   8192000,
				NumRead:        3,
				NumWrite:       1000,
				TotalReadTime:  1500 * time.Millisecond,
				TotalWriteTime: 2 * time.Second,
				ReadErrors:     1,
				WriteErrors:    2,
				ReadRetries:    3,
				WriteRetries:   4,
			},
		},
	}

	expected := `# HELP node_disk_read_bytes_total The total number of bytes read successfully.
# TYPE node_disk_read_bytes_total counter
node_disk_read_bytes_total{device="disk0"} 16384
# HELP node_disk_read_errors_total The total number of read errors.
# TYPE node_disk_read_errors_total counter
node_disk_read_errors_total{device="disk0"} 1
# HELP node_disk_read_retries_total The total number of read retries.
# TYPE node_disk_read_retries_total counter
node_disk_read_retries_total{device="disk0"} 3
# HELP node_disk_read_time_seconds_total The total number of seconds spent by all reads.
# TYPE node_disk_read_time_seconds_total counter
node_disk_read_time_seconds_total{device="disk0"} 1.5
# HELP node_disk_reads_completed_total The total number of reads completed successfully.
# TYPE node_disk_reads_completed_total counter
node_disk_reads_completed_total{device="disk0"} 3
# HELP node_disk_write_errors_total The total number of write errors.
# TYPE node_disk_write_errors_total counter
node_disk_write_errors_total{device="disk0"} 2
# HELP node_disk_write_retries_total The total number of write retries.
# TYPE node_disk_write_retries_total counter
node_disk_write_retries_total{device="disk0"} 4
# HELP node_disk_write_time_seconds_total This is the total number of seconds spent by all writes.
# TYPE node_disk_write_time_seconds_total counter
node_disk_write_time_seconds_total{device="disk0"} 2
# HELP node_disk_writes_completed_total The total number of writes completed successfully.
# TYPE node_disk_writes_completed_total counter
node_disk_writes_completed_total{device="disk0"} 1000
# HELP node_disk_written_bytes_total The total number of bytes written successfully.
# TYPE node_disk_written_bytes_total counter
node_disk_written_bytes_total{device="disk0"} 8.192e+06
`

	// CollectAndCompare checks every collected metric, so it also fails if the
	// sector families, which divided operation counts by the block size, return.
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected)); err != nil {
		t.Fatal(err)
	}
}
