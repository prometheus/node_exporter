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

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type arpCollector struct {
	entries *prometheus.Desc
	logger  log.Logger
}

func init() {
	registerCollector("arp", defaultEnabled, NewARPCollector)
}

// NewARPCollector returns a new Collector exposing ARP stats.
func NewARPCollector(logger log.Logger) (Collector, error) {
	return &arpCollector{
		entries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "arp", "entries"),
			"ARP entries by device",
			[]string{"device"}, nil,
		),
		logger: logger,
	}, nil
}

func getARPEntries() (map[string]uint32, error) {
	file, err := os.Open(procFilePath("net/arp"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := parseARPEntries(file)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// TODO: This should get extracted to the github.com/prometheus/procfs package
// to support more complete parsing of /proc/net/arp. Instead of adding
// more fields to this function's return values it should get moved and
// changed to support each field.
func parseARPEntries(data io.Reader) (map[string]uint32, error) {
	scanner := bufio.NewScanner(data)
	entries := make(map[string]uint32)

	for scanner.Scan() {
		columns := strings.Fields(scanner.Text())

		if len(columns) < 6 {
			return nil, fmt.Errorf("unexpected ARP table format")
		}

		if columns[0] != "IP" {
			deviceIndex := len(columns) - 1
			entries[columns[deviceIndex]]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse ARP info: %w", err)
	}

	return entries, nil
}

func (c *arpCollector) Update(ch chan<- prometheus.Metric) error {
	entries, err := getARPEntries()
	if err != nil {
		return fmt.Errorf("could not get ARP entries: %w", err)
	}

	for device, entryCount := range entries {
		ch <- prometheus.MustNewConstMetric(
			c.entries, prometheus.GaugeValue, float64(entryCount), device)
	}

	return nil
}
