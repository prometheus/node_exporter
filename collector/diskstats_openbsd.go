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

// +build !nodiskstats

package collector

import (
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

/*
#include <sys/types.h>
#include <sys/disk.h>
*/
import "C"

const (
	diskSubsystem = "disk"
)

type diskstatsCollector struct {
	rxfer  typedDesc
	rbytes typedDesc
	wxfer  typedDesc
	wbytes typedDesc
	time   typedDesc
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector() (Collector, error) {
	var diskLabelNames = []string{"device"}

	return &diskstatsCollector{
		rxfer: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "reads_completed_total"),
			"The total number of reads completed successfully.",
			diskLabelNames, nil,
		), prometheus.CounterValue},
		rbytes: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "read_bytes_total"),
			"The total number of bytes read successfully.",
			diskLabelNames, nil,
		), prometheus.CounterValue},
		wxfer: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "writes_completed_total"),
			"The total number of writes completed successfully.",
			diskLabelNames, nil,
		), prometheus.CounterValue},
		wbytes: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "written_bytes_total"),
			"The total number of bytes written successfully.",
			diskLabelNames, nil,
		), prometheus.CounterValue},
		time: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, diskSubsystem, "io_time_seconds_total"),
			"The total number of seconds spent by all IO.",
			diskLabelNames, nil,
		), prometheus.CounterValue},
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	diskstatsb, err := unix.SysctlRaw("hw.diskstats")
	if err != nil {
		return err
	}

	ndisks := len(diskstatsb) / C.sizeof_struct_diskstats
	diskstats := *(*[]C.struct_diskstats)(unsafe.Pointer(&diskstatsb))

	for i := 0; i < ndisks; i++ {
		diskname := C.GoString(&diskstats[i].ds_name[0])

		ch <- c.rxfer.mustNewConstMetric(float64(diskstats[i].ds_rxfer), diskname)
		ch <- c.rbytes.mustNewConstMetric(float64(diskstats[i].ds_rbytes), diskname)
		ch <- c.wxfer.mustNewConstMetric(float64(diskstats[i].ds_wxfer), diskname)
		ch <- c.wbytes.mustNewConstMetric(float64(diskstats[i].ds_wbytes), diskname)
		time := float64(diskstats[i].ds_time.tv_sec) + float64(diskstats[i].ds_time.tv_usec)/1000000
		ch <- c.time.mustNewConstMetric(time, diskname)
	}
	return nil
}
