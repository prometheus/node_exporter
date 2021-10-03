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

//go:build !nomeminfo_numa
// +build !nomeminfo_numa

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	memInfoNumaSubsystem = "memory_numa"
)

var meminfoNodeRE = regexp.MustCompile(`.*devices/system/node/node([0-9]*)`)

type meminfoMetric struct {
	metricName string
	metricType prometheus.ValueType
	numaNode   string
	value      float64
}

type meminfoNumaCollector struct {
	metricDescs map[string]*prometheus.Desc
	logger      log.Logger
}

func init() {
	registerCollector("meminfo_numa", defaultDisabled, NewMeminfoNumaCollector)
}

// NewMeminfoNumaCollector returns a new Collector exposing memory stats.
func NewMeminfoNumaCollector(logger log.Logger) (Collector, error) {
	return &meminfoNumaCollector{
		metricDescs: map[string]*prometheus.Desc{},
		logger:      logger,
	}, nil
}

func (c *meminfoNumaCollector) Update(ch chan<- prometheus.Metric) error {
	metrics, err := getMemInfoNuma()
	if err != nil {
		return fmt.Errorf("couldn't get NUMA meminfo: %w", err)
	}
	for _, v := range metrics {
		desc, ok := c.metricDescs[v.metricName]
		if !ok {
			desc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, memInfoNumaSubsystem, v.metricName),
				fmt.Sprintf("Memory information field %s.", v.metricName),
				[]string{"node"}, nil)
			c.metricDescs[v.metricName] = desc
		}
		ch <- prometheus.MustNewConstMetric(desc, v.metricType, v.value, v.numaNode)
	}
	return nil
}

func getMemInfoNuma() ([]meminfoMetric, error) {
	var (
		metrics []meminfoMetric
	)

	nodes, err := filepath.Glob(sysFilePath("devices/system/node/node[0-9]*"))
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		meminfoFile, err := os.Open(filepath.Join(node, "meminfo"))
		if err != nil {
			return nil, err
		}
		defer meminfoFile.Close()

		numaInfo, err := parseMemInfoNuma(meminfoFile)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, numaInfo...)

		numastatFile, err := os.Open(filepath.Join(node, "numastat"))
		if err != nil {
			return nil, err
		}
		defer numastatFile.Close()

		nodeNumber := meminfoNodeRE.FindStringSubmatch(node)
		if nodeNumber == nil {
			return nil, fmt.Errorf("device node string didn't match regexp: %s", node)
		}

		numaStat, err := parseMemInfoNumaStat(numastatFile, nodeNumber[1])
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, numaStat...)
	}

	return metrics, nil
}

func parseMemInfoNuma(r io.Reader) ([]meminfoMetric, error) {
	var (
		memInfo []meminfoMetric
		scanner = bufio.NewScanner(r)
		re      = regexp.MustCompile(`\((.*)\)`)
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)

		fv, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in meminfo: %w", err)
		}
		switch l := len(parts); {
		case l == 4: // no unit
		case l == 5 && parts[4] == "kB": // has unit
			fv *= 1024
		default:
			return nil, fmt.Errorf("invalid line in meminfo: %s", line)
		}
		metric := strings.TrimRight(parts[2], ":")

		// Active(anon) -> Active_anon
		metric = re.ReplaceAllString(metric, "_${1}")
		memInfo = append(memInfo, meminfoMetric{metric, prometheus.GaugeValue, parts[1], fv})
	}

	return memInfo, scanner.Err()
}

func parseMemInfoNumaStat(r io.Reader, nodeNumber string) ([]meminfoMetric, error) {
	var (
		numaStat []meminfoMetric
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line scan did not return 2 fields: %s", line)
		}

		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in numastat: %w", err)
		}

		numaStat = append(numaStat, meminfoMetric{parts[0] + "_total", prometheus.CounterValue, nodeNumber, fv})
	}
	return numaStat, scanner.Err()
}
