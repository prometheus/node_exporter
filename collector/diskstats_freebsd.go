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

// +build freebsd
// +build !nodiskstats

package collector

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	BINTIME_SCALE = 5.42101086242752217003726400434970855712890625e-20
)

type typedDescFunc struct {
	typedDesc
	value func(stat *diskStat) float64
}

type diskstatsCollector struct {
	descs  []typedDescFunc
	logger log.Logger
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector(logger log.Logger) (Collector, error) {
	return &diskstatsCollector{
		descs: []typedDescFunc{
			{
				typedDesc: typedDesc{
					desc:      readsCompletedDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *diskStat) float64 {
					return stat.readsCompleted
				},
			},
			{
				typedDesc: typedDesc{
					desc:      readTimeSecondsDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *diskStat) float64 {
					return stat.readTimeSeconds
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writesCompletedDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *diskStat) float64 {
					return stat.writesCompleted
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writeTimeSecondsDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *diskStat) float64 {
					return stat.writeTimeSeconds
				},
			},
			{
				typedDesc: typedDesc{
					desc:      readBytesDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *diskStat) float64 {
					return stat.readBytes
				},
			},
			{
				typedDesc: typedDesc{
					desc:      writtenBytesDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *diskStat) float64 {
					return stat.writtenBytes
				},
			},
			{
				typedDesc: typedDesc{
					desc:      ioTimeSecondsDesc,
					valueType: prometheus.CounterValue,
				},
				value: func(stat *diskStat) float64 {
					return stat.ioTimeSeconds
				},
			},
		},
		logger: logger,
	}, nil
}

type diskStat struct {
	readsCompleted   float64
	writesCompleted  float64
	readBytes        float64
	writtenBytes     float64
	readTimeSeconds  float64
	writeTimeSeconds float64
	ioTimeSeconds    float64
	name             string
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	r, err := unix.Sysctl("kern.devstat.all")
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %w", err)
	}
	buf := []byte(r)
	length := len(buf)

	count := int(uint64(length) / uint64(sizeOfDevstat))

	buf = buf[8:] // devstat.all has version in the head.
	// parse buf to Devstat
	for i := 0; i < count; i++ {
		b := buf[i*sizeOfDevstat : i*sizeOfDevstat+sizeOfDevstat]

		sizeOfDeviceName := deviceNameStop - deviceNameStart
		deviceName := make([]byte, sizeOfDeviceName)
		for i, char := range b[deviceNameStart:deviceNameStop] {
			if char == 0 {
				sizeOfDeviceName = i
				break
			}
			deviceName[i] = byte(char)
		}
		ds := diskStat{
			readsCompleted:   float64(binary.LittleEndian.Uint64(b[operationsReadStart:operationsReadStop])),
			writesCompleted:  float64(binary.LittleEndian.Uint64(b[operationsWriteStart:operationsWriteStop])),
			readBytes:        float64(binary.LittleEndian.Uint64(b[bytesReadStart:bytesReadStop])),
			writtenBytes:     float64(binary.LittleEndian.Uint64(b[bytesWriteStart:bytesWriteStop])),
			readTimeSeconds:  float64(binary.LittleEndian.Uint64(b[durationReadSecStart:durationReadSecStop])) + float64(binary.LittleEndian.Uint64(b[durationReadFracStart:durationReadFracStop]))*BINTIME_SCALE*1000,
			writeTimeSeconds: float64(binary.LittleEndian.Uint64(b[durationWriteSecStart:durationWriteSecStop])) + float64(binary.LittleEndian.Uint64(b[durationWriteFracStart:durationWriteFracStop]))*BINTIME_SCALE*1000,
			ioTimeSeconds:    float64(binary.LittleEndian.Uint64(b[busyTimeSecStart:busyTimeSecStop])) + float64(binary.LittleEndian.Uint64(b[busyTimeFracStart:busyTimeFracStop]))*BINTIME_SCALE*1000,
			name:             string(deviceName[0:sizeOfDeviceName]) + strconv.Itoa(int(binary.LittleEndian.Uint32(b[unitNumberStart:unitNumberStop]))),
		}
		for _, desc := range c.descs {
			v := desc.value(&ds)
			ch <- desc.mustNewConstMetric(v, ds.name)
		}
	}

	return nil
}
