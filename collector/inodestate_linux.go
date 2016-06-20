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

// +build !noinodestate

package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	inodeStateSubsystem = "inodestate"
)

type inodeStateCollector struct {
	metrics map[string]prometheus.Gauge
}

func init() {
	Factories[inodeStateSubsystem] = NewInodeStateCollector
}

// NewInodeStateCollector returns a new Collector exposing inode-state stats.
func NewInodeStateCollector() (Collector, error) {
	return &inodeStateCollector{
		metrics: map[string]prometheus.Gauge{},
	}, nil
}

func (c *inodeStateCollector) Update(ch chan<- prometheus.Metric) (err error) {
	inode, err := getInodeStateStats(procFilePath("sys/fs/inode-state"))
	if err != nil {
		return fmt.Errorf("couldn't get inode-state: %s", err)
	}
	for name, value := range inode {
		if _, ok := c.metrics[name]; !ok {
			c.metrics[name] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: Namespace,
					Subsystem: inodeStateSubsystem,
					Name:      name,
					Help:      fmt.Sprintf("inode-state statistics: %s.", name),
				},
			)
		}
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value %s in inode-state: %s", value, err)
		}
		c.metrics[name].Set(v)
	}
	for _, m := range c.metrics {
		m.Collect(ch)
	}
	return err
}

func getInodeStateStats(fileName string) (map[string]string, error) {
	inodeState, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	// The inode-state proc file is separated by tabs, not spaces.
	// It has 7 fields but only the first three are used.
	parts := strings.Fields(string(inodeState))
	if len(parts) < 3 {
		return nil, fmt.Errorf("Unexpected number of fields in %s: %d",
			fileName, len(parts))
	}
	var inodeStateStats = map[string]string{}
	// The inode-state proc is only 1 line with 7 values.
	inodeStateStats["allocated"] = parts[0]
	inodeStateStats["free"] = parts[1]
	inodeStateStats["preshrink"] = parts[2]

	return inodeStateStats, nil
}
