// Copyright 2016 The Prometheus Authors
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

// +build !novmstat

package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	vmStatSubsystem = "vmstat"
)

type vmStatCollector struct{}

func init() {
	Factories["vmstat"] = NewvmStatCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// vmstat stats.
func NewvmStatCollector() (Collector, error) {
	return &vmStatCollector{}, nil
}

func (c *vmStatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	file, err := os.Open(procFilePath("vmstat"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		value, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return err
		}

		metric := prometheus.NewUntyped(prometheus.UntypedOpts{
			Namespace: Namespace,
			Subsystem: vmStatSubsystem,
			Name:      parts[0],
			Help:      fmt.Sprintf("/proc/vmstat information field %s.", parts[0]),
		})
		metric.Set(value)
		metric.Collect(ch)
	}
	return err
}
