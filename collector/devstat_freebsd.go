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

// +build !nodevstat

package collector

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// #cgo LDFLAGS: -ldevstat -lkvm
// #include "devstat_freebsd.h"
import "C"

const (
	devstatSubsystem = "devstat"
)

type devstatCollector struct {
	bytes       typedDesc
	bytes_total typedDesc
	transfers   typedDesc
	duration    typedDesc
	busyTime    typedDesc
	blocks      typedDesc
}

func init() {
	Factories["devstat"] = NewDevstatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// Device stats.
func NewDevstatCollector() (Collector, error) {
	return &devstatCollector{
		bytes: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, devstatSubsystem, "bytes_total"),
			"The total number of bytes in transactions.",
			[]string{"device", "type"}, nil,
		), prometheus.CounterValue},
		transfers: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, devstatSubsystem, "transfers_total"),
			"The total number of transactions.",
			[]string{"device", "type"}, nil,
		), prometheus.CounterValue},
		duration: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, devstatSubsystem, "duration_seconds_total"),
			"The total duration of transactions in seconds.",
			[]string{"device", "type"}, nil,
		), prometheus.CounterValue},
		busyTime: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, devstatSubsystem, "busy_time_seconds_total"),
			"Total time the device had one or more transactions outstanding in seconds.",
			[]string{"device"}, nil,
		), prometheus.CounterValue},
		blocks: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, devstatSubsystem, "blocks_transferred_total"),
			"The total number of blocks transferred.",
			[]string{"device"}, nil,
		), prometheus.CounterValue},
	}, nil
}

func (c *devstatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	count := C._get_ndevs()
	if count == -1 {
		return errors.New("devstat_getdevs() failed")
	}

	for i := C.int(0); i < count; i++ {
		stats := C._get_stats(i)
		device := fmt.Sprintf("%s%d", C.GoString(&stats.device[0]), stats.unit)
		ch <- c.bytes.mustNewConstMetric(float64(stats.bytes.read), device, "read")
		ch <- c.bytes.mustNewConstMetric(float64(stats.bytes.write), device, "write")
		ch <- c.transfers.mustNewConstMetric(float64(stats.transfers.other), device, "other")
		ch <- c.transfers.mustNewConstMetric(float64(stats.transfers.read), device, "read")
		ch <- c.transfers.mustNewConstMetric(float64(stats.transfers.write), device, "write")
		ch <- c.duration.mustNewConstMetric(float64(stats.duration.other), device, "other")
		ch <- c.duration.mustNewConstMetric(float64(stats.duration.read), device, "read")
		ch <- c.duration.mustNewConstMetric(float64(stats.duration.write), device, "write")
		ch <- c.busyTime.mustNewConstMetric(float64(stats.busyTime), device)
		ch <- c.blocks.mustNewConstMetric(float64(stats.blocks), device)
	}
	return err
}
