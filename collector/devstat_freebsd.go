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

//go:build !nodevstat
// +build !nodevstat

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
)

// #cgo LDFLAGS: -ldevstat -lkvm -lelf
// #include "devstat_freebsd.h"
import "C"

const (
	devstatSubsystem = "devstat"
)

type devstatCollector struct {
	mu      sync.Mutex
	devinfo *C.struct_devinfo

	bytes     typedDesc
	transfers typedDesc
	duration  typedDesc
	busyTime  typedDesc
	blocks    typedDesc
	logger    *slog.Logger
}

func init() {
	registerCollector("devstat", defaultDisabled, NewDevstatCollector)
}

// NewDevstatCollector returns a new Collector exposing Device stats.
func NewDevstatCollector(logger *slog.Logger) (Collector, error) {
	return &devstatCollector{
		devinfo: &C.struct_devinfo{},
		bytes: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, devstatSubsystem, "bytes_total"),
			"The total number of bytes in transactions.",
			[]string{"device", "type"}, nil,
		), prometheus.CounterValue},
		transfers: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, devstatSubsystem, "transfers_total"),
			"The total number of transactions.",
			[]string{"device", "type"}, nil,
		), prometheus.CounterValue},
		duration: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, devstatSubsystem, "duration_seconds_total"),
			"The total duration of transactions in seconds.",
			[]string{"device", "type"}, nil,
		), prometheus.CounterValue},
		busyTime: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, devstatSubsystem, "busy_time_seconds_total"),
			"Total time the device had one or more transactions outstanding in seconds.",
			[]string{"device"}, nil,
		), prometheus.CounterValue},
		blocks: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, devstatSubsystem, "blocks_transferred_total"),
			"The total number of blocks transferred.",
			[]string{"device"}, nil,
		), prometheus.CounterValue},
		logger: logger,
	}, nil
}

func (c *devstatCollector) Update(ch chan<- prometheus.Metric) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var stats *C.Stats
	n := C._get_stats(c.devinfo, &stats)
	if n == -1 {
		return errors.New("devstat_getdevs failed")
	}

	base := unsafe.Pointer(stats)
	for i := C.int(0); i < n; i++ {
		offset := i * C.int(C.sizeof_Stats)
		stat := (*C.Stats)(unsafe.Pointer(uintptr(base) + uintptr(offset)))

		device := fmt.Sprintf("%s%d", C.GoString(&stat.device[0]), stat.unit)
		ch <- c.bytes.mustNewConstMetric(float64(stat.bytes.read), device, "read")
		ch <- c.bytes.mustNewConstMetric(float64(stat.bytes.write), device, "write")
		ch <- c.transfers.mustNewConstMetric(float64(stat.transfers.other), device, "other")
		ch <- c.transfers.mustNewConstMetric(float64(stat.transfers.read), device, "read")
		ch <- c.transfers.mustNewConstMetric(float64(stat.transfers.write), device, "write")
		ch <- c.duration.mustNewConstMetric(float64(stat.duration.other), device, "other")
		ch <- c.duration.mustNewConstMetric(float64(stat.duration.read), device, "read")
		ch <- c.duration.mustNewConstMetric(float64(stat.duration.write), device, "write")
		ch <- c.busyTime.mustNewConstMetric(float64(stat.busyTime), device)
		ch <- c.blocks.mustNewConstMetric(float64(stat.blocks), device)
	}
	C.free(unsafe.Pointer(stats))
	return nil
}
