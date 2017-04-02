// Copyright 2017 The Prometheus Authors
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

// +build !noarp

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type arpCollector struct {
	count *prometheus.Desc
}

func init() {
	Factories["arp"] = NewArpCollector
}

// NewArpCollector returns a new Collector exposing ARP stats.
func NewArpCollector() (Collector, error) {
	return &arpCollector{
		count: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "arp",
				"count"), "ARP entries by device",
			[]string{"device"}, nil,
		),
	}, nil
}

func getArpEntries() (map[string]uint32, error) {
	file, err := os.Open(procFilePath("net/arp"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := parseArpEntries(file)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func parseArpEntries(data io.Reader) (map[string]uint32, error) {
	scanner := bufio.NewScanner(data)
	entries := make(map[string]uint32)

	for scanner.Scan() {
		columns := strings.Fields(scanner.Text())

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to parse ARP info: %s", err)
		}

		if columns[0] != "IP" {
			deviceIndex := len(columns) - 1
			entries[columns[deviceIndex]]++
		}
	}

	return entries, nil
}

func (c *arpCollector) Update(ch chan<- prometheus.Metric) error {
	entries, err := getArpEntries()
	if err != nil {
		return fmt.Errorf("could not get ARP entries: %s", err)
	}

	for device, entryCount := range entries {
		ch <- prometheus.MustNewConstMetric(
			c.count, prometheus.GaugeValue, float64(entryCount), device)
	}

	return nil
}
