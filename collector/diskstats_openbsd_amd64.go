// Copyright 2020 The Prometheus Authors
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
	"unsafe"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	DS_DISKNAMELEN = 16
)

type DiskStats struct {
	Name       [DS_DISKNAMELEN]int8
	Busy       int32
	Rxfer      uint64
	Wxfer      uint64
	Seek       uint64
	Rbytes     uint64
	Wbytes     uint64
	Attachtime unix.Timeval
	Timestamp  unix.Timeval
	Time       unix.Timeval
}

type diskstatsCollector struct {
	rxfer  typedDesc
	rbytes typedDesc
	wxfer  typedDesc
	wbytes typedDesc
	time   typedDesc
	logger log.Logger
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector(logger log.Logger) (Collector, error) {
	return &diskstatsCollector{
		rxfer:  typedDesc{readsCompletedDesc, prometheus.CounterValue},
		rbytes: typedDesc{readBytesDesc, prometheus.CounterValue},
		wxfer:  typedDesc{writesCompletedDesc, prometheus.CounterValue},
		wbytes: typedDesc{writtenBytesDesc, prometheus.CounterValue},
		time:   typedDesc{ioTimeSecondsDesc, prometheus.CounterValue},
		logger: logger,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	diskstatsb, err := unix.SysctlRaw("hw.diskstats")
	if err != nil {
		return err
	}

	ndisks := len(diskstatsb) / int(unsafe.Sizeof(DiskStats{}))
	diskstats := *(*[]DiskStats)(unsafe.Pointer(&diskstatsb))

	for i := 0; i < ndisks; i++ {
		dn := *(*[DS_DISKNAMELEN]int8)(unsafe.Pointer(&diskstats[i].Name[0]))
		diskname := int8ToString(dn[:])

		ch <- c.rxfer.mustNewConstMetric(float64(diskstats[i].Rxfer), diskname)
		ch <- c.rbytes.mustNewConstMetric(float64(diskstats[i].Rbytes), diskname)
		ch <- c.wxfer.mustNewConstMetric(float64(diskstats[i].Wxfer), diskname)
		ch <- c.wbytes.mustNewConstMetric(float64(diskstats[i].Wbytes), diskname)
		time := float64(diskstats[i].Time.Sec) + float64(diskstats[i].Time.Usec)/1000000
		ch <- c.time.mustNewConstMetric(time, diskname)
	}
	return nil
}
