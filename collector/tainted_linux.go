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

// +build linux
// +build !notainted

package collector

import (
	"fmt"
	"math"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type taintedCollector struct {
	bits []typedDesc
}

func init() {
	registerCollector("tainted", defaultDisabled, NewTaintedCollector)
}

// NewTaintedCollector returns a new Collector exposing tainted bits.
// Bits are listed in least-significant order, ie bit 0, bit 1, ..., bit N.
// See https://www.kernel.org/doc/html/latest/admin-guide/tainted-kernels.html
// for more information on tainted bits.
func NewTaintedCollector() (Collector, error) {
	return &taintedCollector{
		bits: []typedDesc{
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "P"),
					"Proprietary module was loaded.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "F"),
					"Module was force loaded.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "S"),
					"SMP kernel oops on an officially SMP incapable processor.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "R"),
					"Module was force unloaded.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "M"),
					"Processor reported a Machine Check Exception (MCE).",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "B"),
					"Bad page referenced or some unexpected page flags.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "U"),
					"Taint requested by userspace application.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "D"),
					"Kernel died recently, i.e. there was an OOPS or BUG.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "A"),
					"An ACPI table was overridden by user.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "W"),
					"Kernel issued warning.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "C"),
					"Staging driver was loaded.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "I"),
					"Workaround for bug in platform firmware applied.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "O"),
					"Externally-built (\"out-of-tree\") module was loaded.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "E"),
					"Unsigned module was loaded.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "L"),
					"Soft lockup occurred.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "K"),
					"Kernel has been live patched.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "X"),
					"Auxiliary taint, defined and used by for distros.",
					nil, nil),
				prometheus.GaugeValue,
			},
			{
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "tainted", "T"),
					"The kernel was built with the struct randomization plugin.",
					nil, nil),
				prometheus.GaugeValue,
			},
		},
	}, nil
}

func (c *taintedCollector) Update(ch chan<- prometheus.Metric) (err error) {
	bits, err := getBits(len(c.bits))
	if err != nil {
		return fmt.Errorf("couldn't get tainted: %s", err)
	}
	for i, bit := range bits {
		log.Debugf("Set tainted bit %d: %0.0f", i, bit)
		ch <- c.bits[i].mustNewConstMetric(bit)
	}
	return err
}

// Read tainted value from /proc/sys/kernel/tainted and return least
// significant bits as an array.
func getBits(count int) (bits []float64, err error) {
	tainted, err := readUintFromFile(procFilePath("sys/kernel/tainted"))
	if err != nil {
		return nil, fmt.Errorf("couldn't get tainted value: %s", err)
	}
	return parseBits(tainted, count), nil
}

// Return "count" least significant bits of "number".
func parseBits(number uint64, count int) (bits []float64) {
	bits = make([]float64, count)
	for i := 0; i < count; i++ {
		val := number & uint64(math.Pow(2, float64(i)))
		if val != 0 {
			bits[i] = 1.0
		}
	}
	return bits
}
