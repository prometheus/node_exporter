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

// +build !nobuddyinfo
// +build !windows,!netbsd

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type buddyInfo map[string]map[string][]float64

const (
	buddyInfoSubsystem = "buddyinfo"
)

type buddyinfoCollector struct {
	desc *prometheus.Desc
}

func init() {
	Factories["buddyinfo"] = NewBuddyinfoCollector
}

// NewBuddyinfoCollector returns a new Collector exposing buddyinfo stats.
func NewBuddyinfoCollector() (Collector, error) {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, buddyInfoSubsystem, "count"),
		"Count of free blocks according to size.",
		[]string{"node", "zone", "size"}, nil,
	)
	return &buddyinfoCollector{desc}, nil
}

// Update calls (*buddyinfoCollector).getBuddyInfo to get the platform specific
// buddyinfo metrics.
func (c *buddyinfoCollector) Update(ch chan<- prometheus.Metric) (err error) {
	buddyInfo, err := c.getBuddyInfo()
	if err != nil {
		return fmt.Errorf("couldn't get buddyinfo: %s", err)
	}
	log.Debugf("Set node_buddy: %#v", buddyInfo)
	for node, zones := range buddyInfo {
		for zone, values := range zones {
			for size, value := range values {
				ch <- prometheus.MustNewConstMetric(
					c.desc,
					prometheus.GaugeValue, value,
					node, zone, strconv.Itoa(size),
				)
			}
		}
	}
	return nil
}
func (c *buddyinfoCollector) getBuddyInfo() (buddyInfo, error) {
	file, err := os.Open(procFilePath("buddyinfo"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseBuddyInfo(file)
}

func parseBuddyInfo(r io.Reader) (buddyInfo, error) {
	var (
		buddyInfo = buddyInfo{}
		scanner   = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		var err error
		line := scanner.Text()
		parts := strings.Fields(string(line))
		node := strings.TrimRight(parts[1], ",")
		zone := strings.TrimRight(parts[3], ",")
		arraySize := len(parts[4:])
		sizes := make([]float64, arraySize)
		for i := 0; i < arraySize; i++ {
			sizes[i], err = strconv.ParseFloat(parts[i+4], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid value in buddyinfo: %s", err)
			}
		}

		if _, ok := buddyInfo[node]; !ok {
			buddyInfo[node] = make(map[string][]float64)
		}
		buddyInfo[node][zone] = sizes
	}

	return buddyInfo, nil
}
