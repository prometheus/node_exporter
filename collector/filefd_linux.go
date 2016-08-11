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

// +build !nofilefd

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	fileFDStatSubsystem = "filefd"
)

type fileFDStatCollector struct{}

func init() {
	Factories[fileFDStatSubsystem] = NewFileFDStatCollector
}

// NewFileFDStatCollector returns a new Collector exposing file-nr stats.
func NewFileFDStatCollector() (Collector, error) {
	return &fileFDStatCollector{}, nil
}

func (c *fileFDStatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	fileFDStat, err := getFileFDStats(procFilePath("sys/fs/file-nr"))
	if err != nil {
		return fmt.Errorf("couldn't get file-nr: %s", err)
	}
	for name, value := range fileFDStat {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value %s in file-nr: %s", value, err)
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, fileFDStatSubsystem, name),
				fmt.Sprintf("File descriptor statistics: %s.", name),
				nil, nil,
			),
			prometheus.GaugeValue, v,
		)
	}
	return nil
}

func getFileFDStats(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseFileFDStats(file, fileName)
}

func parseFileFDStats(r io.Reader, fileName string) (map[string]string, error) {
	var scanner = bufio.NewScanner(r)
	scanner.Scan()
	// The file-nr proc file is separated by tabs, not spaces.
	line := strings.Split(string(scanner.Text()), "\u0009")
	var fileFDStat = map[string]string{}
	// The file-nr proc is only 1 line with 3 values.
	fileFDStat["allocated"] = line[0]
	// The second value is skipped as it will always be zero in linux 2.6.
	fileFDStat["maximum"] = line[2]

	return fileFDStat, nil
}
