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
	"strings"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testDiskStatsCollector struct {
	dsc Collector
}

func (c testDiskStatsCollector) Collect(ch chan<- prometheus.Metric) {
	c.dsc.Update(ch)
}

func (c testDiskStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func NewTestDiskStatsCollector(logger log.Logger) (prometheus.Collector, error) {
	dsc, err := NewDiskstatsCollector(logger)
	if err != nil {
		return testDiskStatsCollector{}, err
	}
	return testDiskStatsCollector{
		dsc: dsc,
	}, err
}

func TestDiskStats(t *testing.T) {
	*sysPath = "fixtures/sys"
	*procPath = "fixtures/proc"
	*ignoredDevices = "^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$"
	testcase := `# HELP node_disk_discard_time_seconds_total This is the total number of seconds spent by all discards.
# TYPE node_disk_discard_time_seconds_total counter
node_disk_discard_time_seconds_total{device="sdb"} 11.13
node_disk_discard_time_seconds_total{device="sdc"} 11.13
# HELP node_disk_discarded_sectors_total The total number of sectors discarded successfully.
# TYPE node_disk_discarded_sectors_total counter
node_disk_discarded_sectors_total{device="sdb"} 1.925173784e+09
node_disk_discarded_sectors_total{device="sdc"} 1.25173784e+08
# HELP node_disk_discards_completed_total The total number of discards completed successfully.
# TYPE node_disk_discards_completed_total counter
node_disk_discards_completed_total{device="sdb"} 68851
node_disk_discards_completed_total{device="sdc"} 18851
# HELP node_disk_discards_merged_total The total number of discards merged.
# TYPE node_disk_discards_merged_total counter
node_disk_discards_merged_total{device="sdb"} 0
node_disk_discards_merged_total{device="sdc"} 0
# HELP node_disk_flush_requests_time_seconds_total This is the total number of seconds spent by all flush requests.
# TYPE node_disk_flush_requests_time_seconds_total counter
node_disk_flush_requests_time_seconds_total{device="sdc"} 1.944
# HELP node_disk_flush_requests_total The total number of flush requests completed successfully
# TYPE node_disk_flush_requests_total counter
node_disk_flush_requests_total{device="sdc"} 1555
# HELP node_disk_info Info of /sys/block/<block_device>.
# TYPE node_disk_info gauge
node_disk_info{device="dm-0",major="252",minor="0"} 1
node_disk_info{device="dm-1",major="252",minor="1"} 1
node_disk_info{device="dm-2",major="252",minor="2"} 1
node_disk_info{device="dm-3",major="252",minor="3"} 1
node_disk_info{device="dm-4",major="252",minor="4"} 1
node_disk_info{device="dm-5",major="252",minor="5"} 1
node_disk_info{device="mmcblk0",major="179",minor="0"} 1
node_disk_info{device="mmcblk0p1",major="179",minor="1"} 1
node_disk_info{device="mmcblk0p2",major="179",minor="2"} 1
node_disk_info{device="nvme0n1",major="259",minor="0"} 1
node_disk_info{device="sda",major="8",minor="0"} 1
node_disk_info{device="sdb",major="8",minor="0"} 1
node_disk_info{device="sdc",major="8",minor="0"} 1
node_disk_info{device="sr0",major="11",minor="0"} 1
node_disk_info{device="vda",major="254",minor="0"} 1
# HELP node_disk_io_now The number of I/Os currently in progress.
# TYPE node_disk_io_now gauge
node_disk_io_now{device="dm-0"} 0
node_disk_io_now{device="dm-1"} 0
node_disk_io_now{device="dm-2"} 0
node_disk_io_now{device="dm-3"} 0
node_disk_io_now{device="dm-4"} 0
node_disk_io_now{device="dm-5"} 0
node_disk_io_now{device="mmcblk0"} 0
node_disk_io_now{device="mmcblk0p1"} 0
node_disk_io_now{device="mmcblk0p2"} 0
node_disk_io_now{device="nvme0n1"} 0
node_disk_io_now{device="sda"} 0
node_disk_io_now{device="sdb"} 0
node_disk_io_now{device="sdc"} 0
node_disk_io_now{device="sr0"} 0
node_disk_io_now{device="vda"} 0
# HELP node_disk_io_time_seconds_total Total seconds spent doing I/Os.
# TYPE node_disk_io_time_seconds_total counter
node_disk_io_time_seconds_total{device="dm-0"} 11325.968
node_disk_io_time_seconds_total{device="dm-1"} 0.076
node_disk_io_time_seconds_total{device="dm-2"} 65.4
node_disk_io_time_seconds_total{device="dm-3"} 0.016
node_disk_io_time_seconds_total{device="dm-4"} 0.024
node_disk_io_time_seconds_total{device="dm-5"} 58.848
node_disk_io_time_seconds_total{device="mmcblk0"} 0.136
node_disk_io_time_seconds_total{device="mmcblk0p1"} 0.024
node_disk_io_time_seconds_total{device="mmcblk0p2"} 0.068
node_disk_io_time_seconds_total{device="nvme0n1"} 222.766
node_disk_io_time_seconds_total{device="sda"} 9653.880000000001
node_disk_io_time_seconds_total{device="sdb"} 60.730000000000004
node_disk_io_time_seconds_total{device="sdc"} 10.73
node_disk_io_time_seconds_total{device="sr0"} 0
node_disk_io_time_seconds_total{device="vda"} 41614.592000000004
# HELP node_disk_io_time_weighted_seconds_total The weighted # of seconds spent doing I/Os.
# TYPE node_disk_io_time_weighted_seconds_total counter
node_disk_io_time_weighted_seconds_total{device="dm-0"} 1.206301256e+06
node_disk_io_time_weighted_seconds_total{device="dm-1"} 0.084
node_disk_io_time_weighted_seconds_total{device="dm-2"} 129.416
node_disk_io_time_weighted_seconds_total{device="dm-3"} 0.10400000000000001
node_disk_io_time_weighted_seconds_total{device="dm-4"} 0.044
node_disk_io_time_weighted_seconds_total{device="dm-5"} 105.632
node_disk_io_time_weighted_seconds_total{device="mmcblk0"} 0.156
node_disk_io_time_weighted_seconds_total{device="mmcblk0p1"} 0.024
node_disk_io_time_weighted_seconds_total{device="mmcblk0p2"} 0.068
node_disk_io_time_weighted_seconds_total{device="nvme0n1"} 1032.546
node_disk_io_time_weighted_seconds_total{device="sda"} 82621.804
node_disk_io_time_weighted_seconds_total{device="sdb"} 67.07000000000001
node_disk_io_time_weighted_seconds_total{device="sdc"} 17.07
node_disk_io_time_weighted_seconds_total{device="sr0"} 0
node_disk_io_time_weighted_seconds_total{device="vda"} 2.0778722280000001e+06
# HELP node_disk_read_bytes_total The total number of bytes read successfully.
# TYPE node_disk_read_bytes_total counter
node_disk_read_bytes_total{device="dm-0"} 5.13708655616e+11
node_disk_read_bytes_total{device="dm-1"} 1.589248e+06
node_disk_read_bytes_total{device="dm-2"} 1.578752e+08
node_disk_read_bytes_total{device="dm-3"} 1.98144e+06
node_disk_read_bytes_total{device="dm-4"} 529408
node_disk_read_bytes_total{device="dm-5"} 4.3150848e+07
node_disk_read_bytes_total{device="mmcblk0"} 798720
node_disk_read_bytes_total{device="mmcblk0p1"} 81920
node_disk_read_bytes_total{device="mmcblk0p2"} 389120
node_disk_read_bytes_total{device="nvme0n1"} 2.377714176e+09
node_disk_read_bytes_total{device="sda"} 5.13713216512e+11
node_disk_read_bytes_total{device="sdb"} 4.944782848e+09
node_disk_read_bytes_total{device="sdc"} 8.48782848e+08
node_disk_read_bytes_total{device="sr0"} 0
node_disk_read_bytes_total{device="vda"} 1.6727491584e+10
# HELP node_disk_read_time_seconds_total The total number of seconds spent by all reads.
# TYPE node_disk_read_time_seconds_total counter
node_disk_read_time_seconds_total{device="dm-0"} 46229.572
node_disk_read_time_seconds_total{device="dm-1"} 0.084
node_disk_read_time_seconds_total{device="dm-2"} 6.5360000000000005
node_disk_read_time_seconds_total{device="dm-3"} 0.10400000000000001
node_disk_read_time_seconds_total{device="dm-4"} 0.028
node_disk_read_time_seconds_total{device="dm-5"} 0.924
node_disk_read_time_seconds_total{device="mmcblk0"} 0.156
node_disk_read_time_seconds_total{device="mmcblk0p1"} 0.024
node_disk_read_time_seconds_total{device="mmcblk0p2"} 0.068
node_disk_read_time_seconds_total{device="nvme0n1"} 21.650000000000002
node_disk_read_time_seconds_total{device="sda"} 18492.372
node_disk_read_time_seconds_total{device="sdb"} 0.084
node_disk_read_time_seconds_total{device="sdc"} 0.014
node_disk_read_time_seconds_total{device="sr0"} 0
node_disk_read_time_seconds_total{device="vda"} 8655.768
# HELP node_disk_reads_completed_total The total number of reads completed successfully.
# TYPE node_disk_reads_completed_total counter
node_disk_reads_completed_total{device="dm-0"} 5.9910002e+07
node_disk_reads_completed_total{device="dm-1"} 388
node_disk_reads_completed_total{device="dm-2"} 11571
node_disk_reads_completed_total{device="dm-3"} 3870
node_disk_reads_completed_total{device="dm-4"} 392
node_disk_reads_completed_total{device="dm-5"} 3729
node_disk_reads_completed_total{device="mmcblk0"} 192
node_disk_reads_completed_total{device="mmcblk0p1"} 17
node_disk_reads_completed_total{device="mmcblk0p2"} 95
node_disk_reads_completed_total{device="nvme0n1"} 47114
node_disk_reads_completed_total{device="sda"} 2.5354637e+07
node_disk_reads_completed_total{device="sdb"} 326552
node_disk_reads_completed_total{device="sdc"} 126552
node_disk_reads_completed_total{device="sr0"} 0
node_disk_reads_completed_total{device="vda"} 1.775784e+06
# HELP node_disk_reads_merged_total The total number of reads merged.
# TYPE node_disk_reads_merged_total counter
node_disk_reads_merged_total{device="dm-0"} 0
node_disk_reads_merged_total{device="dm-1"} 0
node_disk_reads_merged_total{device="dm-2"} 0
node_disk_reads_merged_total{device="dm-3"} 0
node_disk_reads_merged_total{device="dm-4"} 0
node_disk_reads_merged_total{device="dm-5"} 0
node_disk_reads_merged_total{device="mmcblk0"} 3
node_disk_reads_merged_total{device="mmcblk0p1"} 3
node_disk_reads_merged_total{device="mmcblk0p2"} 0
node_disk_reads_merged_total{device="nvme0n1"} 4
node_disk_reads_merged_total{device="sda"} 3.4367663e+07
node_disk_reads_merged_total{device="sdb"} 841
node_disk_reads_merged_total{device="sdc"} 141
node_disk_reads_merged_total{device="sr0"} 0
node_disk_reads_merged_total{device="vda"} 15386
# HELP node_disk_write_time_seconds_total This is the total number of seconds spent by all writes.
# TYPE node_disk_write_time_seconds_total counter
node_disk_write_time_seconds_total{device="dm-0"} 1.1585578e+06
node_disk_write_time_seconds_total{device="dm-1"} 0
node_disk_write_time_seconds_total{device="dm-2"} 122.884
node_disk_write_time_seconds_total{device="dm-3"} 0
node_disk_write_time_seconds_total{device="dm-4"} 0.016
node_disk_write_time_seconds_total{device="dm-5"} 104.684
node_disk_write_time_seconds_total{device="mmcblk0"} 0
node_disk_write_time_seconds_total{device="mmcblk0p1"} 0
node_disk_write_time_seconds_total{device="mmcblk0p2"} 0
node_disk_write_time_seconds_total{device="nvme0n1"} 1011.053
node_disk_write_time_seconds_total{device="sda"} 63877.96
node_disk_write_time_seconds_total{device="sdb"} 5.007
node_disk_write_time_seconds_total{device="sdc"} 1.0070000000000001
node_disk_write_time_seconds_total{device="sr0"} 0
node_disk_write_time_seconds_total{device="vda"} 2.069221364e+06
# HELP node_disk_writes_completed_total The total number of writes completed successfully.
# TYPE node_disk_writes_completed_total counter
node_disk_writes_completed_total{device="dm-0"} 3.9231014e+07
node_disk_writes_completed_total{device="dm-1"} 74
node_disk_writes_completed_total{device="dm-2"} 153522
node_disk_writes_completed_total{device="dm-3"} 0
node_disk_writes_completed_total{device="dm-4"} 38
node_disk_writes_completed_total{device="dm-5"} 98918
node_disk_writes_completed_total{device="mmcblk0"} 0
node_disk_writes_completed_total{device="mmcblk0p1"} 0
node_disk_writes_completed_total{device="mmcblk0p2"} 0
node_disk_writes_completed_total{device="nvme0n1"} 1.07832e+06
node_disk_writes_completed_total{device="sda"} 2.8444756e+07
node_disk_writes_completed_total{device="sdb"} 41822
node_disk_writes_completed_total{device="sdc"} 11822
node_disk_writes_completed_total{device="sr0"} 0
node_disk_writes_completed_total{device="vda"} 6.038856e+06
# HELP node_disk_writes_merged_total The number of writes merged.
# TYPE node_disk_writes_merged_total counter
node_disk_writes_merged_total{device="dm-0"} 0
node_disk_writes_merged_total{device="dm-1"} 0
node_disk_writes_merged_total{device="dm-2"} 0
node_disk_writes_merged_total{device="dm-3"} 0
node_disk_writes_merged_total{device="dm-4"} 0
node_disk_writes_merged_total{device="dm-5"} 0
node_disk_writes_merged_total{device="mmcblk0"} 0
node_disk_writes_merged_total{device="mmcblk0p1"} 0
node_disk_writes_merged_total{device="mmcblk0p2"} 0
node_disk_writes_merged_total{device="nvme0n1"} 43950
node_disk_writes_merged_total{device="sda"} 1.1134226e+07
node_disk_writes_merged_total{device="sdb"} 2895
node_disk_writes_merged_total{device="sdc"} 1895
node_disk_writes_merged_total{device="sr0"} 0
node_disk_writes_merged_total{device="vda"} 2.0711856e+07
# HELP node_disk_written_bytes_total The total number of bytes written successfully.
# TYPE node_disk_written_bytes_total counter
node_disk_written_bytes_total{device="dm-0"} 2.5891680256e+11
node_disk_written_bytes_total{device="dm-1"} 303104
node_disk_written_bytes_total{device="dm-2"} 2.607828992e+09
node_disk_written_bytes_total{device="dm-3"} 0
node_disk_written_bytes_total{device="dm-4"} 70144
node_disk_written_bytes_total{device="dm-5"} 5.89664256e+08
node_disk_written_bytes_total{device="mmcblk0"} 0
node_disk_written_bytes_total{device="mmcblk0p1"} 0
node_disk_written_bytes_total{device="mmcblk0p2"} 0
node_disk_written_bytes_total{device="nvme0n1"} 2.0199236096e+10
node_disk_written_bytes_total{device="sda"} 2.58916880384e+11
node_disk_written_bytes_total{device="sdb"} 1.01012736e+09
node_disk_written_bytes_total{device="sdc"} 8.852736e+07
node_disk_written_bytes_total{device="sr0"} 0
node_disk_written_bytes_total{device="vda"} 1.0938236928e+11
`

	logger := log.NewLogfmtLogger(os.Stderr)
	collector, err := NewDiskstatsCollector(logger)
	if err != nil {
		panic(err)
	}
	c, err := NewTestDiskStatsCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
		if err != nil {
			panic(fmt.Errorf("failed to update collector: %s", err))
		}
		close(sink)
	}()

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
}
