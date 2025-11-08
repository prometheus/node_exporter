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

//go:build !nodiskstats
// +build !nodiskstats

package collector

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"

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

func NewTestDiskStatsCollector(logger *slog.Logger) (prometheus.Collector, error) {
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
	*udevDataPath = "fixtures/udev/data"
	*diskstatsDeviceExclude = "^(z?ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$"
	testcase := `# HELP node_disk_ata_rotation_rate_rpm ATA disk rotation rate in RPMs (0 for SSDs).
# TYPE node_disk_ata_rotation_rate_rpm gauge
node_disk_ata_rotation_rate_rpm{device="sda"} 7200
node_disk_ata_rotation_rate_rpm{device="sdb"} 0
node_disk_ata_rotation_rate_rpm{device="sdc"} 0
# HELP node_disk_ata_write_cache ATA disk has a write cache.
# TYPE node_disk_ata_write_cache gauge
node_disk_ata_write_cache{device="sda"} 1
node_disk_ata_write_cache{device="sdb"} 1
node_disk_ata_write_cache{device="sdc"} 1
# HELP node_disk_ata_write_cache_enabled ATA disk has its write cache enabled.
# TYPE node_disk_ata_write_cache_enabled gauge
node_disk_ata_write_cache_enabled{device="sda"} 0
node_disk_ata_write_cache_enabled{device="sdb"} 1
node_disk_ata_write_cache_enabled{device="sdc"} 0
# HELP node_disk_device_mapper_info Info about disk device mapper.
# TYPE node_disk_device_mapper_info gauge
node_disk_device_mapper_info{device="dm-0",lv_layer="",lv_name="",name="nvme0n1_crypt",uuid="CRYPT-LUKS2-jolaulot80fy9zsiobkxyxo7y2dqeho2-nvme0n1_crypt",vg_name=""} 1
node_disk_device_mapper_info{device="dm-1",lv_layer="",lv_name="swap_1",name="system-swap_1",uuid="LVM-wbGqQEBL9SxrW2DLntJwgg8fAv946hw3Tvjqh0v31fWgxEtD4BoHO0lROWFUY65T",vg_name="system"} 1
node_disk_device_mapper_info{device="dm-2",lv_layer="",lv_name="root",name="system-root",uuid="LVM-NWEDo8q5ABDyJuC3F8veKNyWfYmeIBfFMS4MF3HakzUhkk7ekDm6fJTHkl2fYHe7",vg_name="system"} 1
node_disk_device_mapper_info{device="dm-3",lv_layer="",lv_name="var",name="system-var",uuid="LVM-hrxHo0rlZ6U95ku5841Lpd17bS1Z7V7lrtEE60DVgE6YEOCdS9gcDGyonWim4hGP",vg_name="system"} 1
node_disk_device_mapper_info{device="dm-4",lv_layer="",lv_name="tmp",name="system-tmp",uuid="LVM-XTNGOHjPWLHcxmJmVu5cWTXEtuzqDeBkdEHAZW5q9LxWQ2d4mb5CchUQzUPJpl8H",vg_name="system"} 1
node_disk_device_mapper_info{device="dm-5",lv_layer="",lv_name="home",name="system-home",uuid="LVM-MtoJaWTpjWRXlUnNFlpxZauTEuYlMvGFutigEzCCrfj8CNh6jCRi5LQJXZCpLjPf",vg_name="system"} 1
# HELP node_disk_discard_time_seconds_total This is the total number of seconds spent by all discards.
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
# HELP node_disk_filesystem_info Info about disk filesystem.
# TYPE node_disk_filesystem_info gauge
node_disk_filesystem_info{device="dm-0",type="LVM2_member",usage="raid",uuid="c3C3uW-gD96-Yw69-c1CJ-5MwT-6ysM-mST0vB",version="LVM2 001"} 1
node_disk_filesystem_info{device="dm-1",type="swap",usage="other",uuid="5272bb60-04b5-49cd-b730-be57c7604450",version="1"} 1
node_disk_filesystem_info{device="dm-2",type="ext4",usage="filesystem",uuid="3deafd0d-faff-4695-8d15-51061ae1f51b",version="1.0"} 1
node_disk_filesystem_info{device="dm-3",type="ext4",usage="filesystem",uuid="5c772222-f7d4-4c8e-87e8-e97df6b7a45e",version="1.0"} 1
node_disk_filesystem_info{device="dm-4",type="ext4",usage="filesystem",uuid="a9479d44-60e1-4015-a1e5-bb065e6dd11b",version="1.0"} 1
node_disk_filesystem_info{device="dm-5",type="ext4",usage="filesystem",uuid="b05b726a-c718-4c4d-8641-7c73a7696d83",version="1.0"} 1
node_disk_filesystem_info{device="mmcblk0p1",type="vfat",usage="filesystem",uuid="6284-658D",version="FAT32"} 1
node_disk_filesystem_info{device="mmcblk0p2",type="ext4",usage="filesystem",uuid="83324ce8-a6f3-4e35-ad64-dbb3d6b87a32",version="1.0"} 1
node_disk_filesystem_info{device="sda",type="LVM2_member",usage="raid",uuid="cVVv6j-HSA2-IY33-1Jmj-dO2H-YL7w-b4Oxqw",version="LVM2 001"} 1
node_disk_filesystem_info{device="sdc",type="LVM2_member",usage="raid",uuid="QFy9W7-Brj3-hQ6v-AF8i-3Zqg-n3Vs-kGY4vb",version="LVM2 001"} 1
# HELP node_disk_flush_requests_time_seconds_total This is the total number of seconds spent by all flush requests.
# TYPE node_disk_flush_requests_time_seconds_total counter
node_disk_flush_requests_time_seconds_total{device="sdc"} 1.944
# HELP node_disk_flush_requests_total The total number of flush requests completed successfully
# TYPE node_disk_flush_requests_total counter
node_disk_flush_requests_total{device="sdc"} 1555
# HELP node_disk_info Info of /sys/block/<block_device>.
# TYPE node_disk_info gauge
node_disk_info{device="dm-0",major="252",minor="0",model="",path="",revision="",rotational="0",serial="",wwn=""} 1
node_disk_info{device="dm-1",major="252",minor="1",model="",path="",revision="",rotational="0",serial="",wwn=""} 1
node_disk_info{device="dm-2",major="252",minor="2",model="",path="",revision="",rotational="0",serial="",wwn=""} 1
node_disk_info{device="dm-3",major="252",minor="3",model="",path="",revision="",rotational="0",serial="",wwn=""} 1
node_disk_info{device="dm-4",major="252",minor="4",model="",path="",revision="",rotational="0",serial="",wwn=""} 1
node_disk_info{device="dm-5",major="252",minor="5",model="",path="",revision="",rotational="0",serial="",wwn=""} 1
node_disk_info{device="mmcblk0",major="179",minor="0",model="",path="platform-df2969f3.mmc",revision="",rotational="0",serial="0x83e36d93",wwn=""} 1
node_disk_info{device="mmcblk0p1",major="179",minor="1",model="",path="platform-df2969f3.mmc",revision="",rotational="0",serial="0x83e36d93",wwn=""} 1
node_disk_info{device="mmcblk0p2",major="179",minor="2",model="",path="platform-df2969f3.mmc",revision="",rotational="0",serial="0x83e36d93",wwn=""} 1
node_disk_info{device="nvme0n1",major="259",minor="0",model="SAMSUNG EHFTF55LURSY-000Y9",path="pci-0000:02:00.0-nvme-1",revision="4NBTUY95",rotational="0",serial="S252B6CU1HG3M1",wwn="eui.p3vbbiejx5aae2r3"} 1
node_disk_info{device="sda",major="8",minor="0",model="TOSHIBA_KSDB4U86",path="pci-0000:3b:00.0-sas-phy7-lun-0",revision="0102",rotational="1",serial="2160A0D5FVGG",wwn="0x7c72382b8de36a64"} 1
node_disk_info{device="sdb",major="8",minor="16",model="SuperMicro_SSD",path="pci-0000:00:1f.2-ata-1",revision="0R",rotational="0",serial="SMC0E1B87ABBB16BD84E",wwn="0xe1b87abbb16bd84e"} 1
node_disk_info{device="sdc",major="8",minor="32",model="INTEL_SSDS9X9SI0",path="pci-0000:00:1f.2-ata-4",revision="0100",rotational="0",serial="3EWB5Y25CWQWA7EH1U",wwn="0x58907ddc573a5de"} 1
node_disk_info{device="sr0",major="11",minor="0",model="Virtual_CDROM0",path="pci-0000:00:14.0-usb-0:1.1:1.0-scsi-0:0:0:0",revision="1.00",rotational="0",serial="AAAABBBBCCCC1",wwn=""} 1
node_disk_info{device="vda",major="254",minor="0",model="",path="pci-0000:00:06.0",revision="",rotational="0",serial="",wwn=""} 1
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

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	collector, err := NewDiskstatsCollector(logger)
	if err != nil {
		t.Fatal(err)
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
