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

// +build !nomeminfo

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	memInfoSubsystem = "memory"
)

type meminfoCollector struct{}

func init() {
	Factories["meminfo"] = NewMeminfoCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// memory stats.
func NewMeminfoCollector() (Collector, error) {
	return &meminfoCollector{}, nil
}

func (c *meminfoCollector) Update(ch chan<- prometheus.Metric) (err error) {
	memInfo, err := getMemInfo()
	if err != nil {
		return fmt.Errorf("couldn't get meminfo: %s", err)
	}
	log.Debugf("Set node_mem: %#v", memInfo)
	for k, v := range memInfo {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, memInfoSubsystem, k),
				fmt.Sprintf("Memory information field %s.", k),
				nil, nil,
			),
			prometheus.GaugeValue, v,
		)
	}
	return nil
}

func getMemInfo() (map[string]float64, error) {
	file, err := os.Open(procFilePath("meminfo"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseMemInfo(file)
}

func parseMemInfo(r io.Reader) (map[string]float64, error) {
	var (
		memInfo = map[string]float64{}
		scanner = bufio.NewScanner(r)
		re      = regexp.MustCompile("\\((.*)\\)")
	)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(string(line))
		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in meminfo: %s", err)
		}
		switch len(parts) {
		case 2: // no unit
		case 3: // has unit, we presume kB
			fv *= 1024
		default:
			return nil, fmt.Errorf("invalid line in meminfo: %s", line)
		}
		key := parts[0][:len(parts[0])-1] // remove trailing : from key
		// Active(anon) -> Active_anon
		key = re.ReplaceAllString(key, "_${1}")
		memInfo[key] = fv
	}

	return memInfo, nil
}
