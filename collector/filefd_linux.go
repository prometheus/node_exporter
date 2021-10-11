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

//go:build !nofilefd
// +build !nofilefd

package collector

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	fileFDStatSubsystem = "filefd"
)

type fileFDStatCollector struct {
	logger log.Logger
}

func init() {
	registerCollector(fileFDStatSubsystem, defaultEnabled, NewFileFDStatCollector)
}

// NewFileFDStatCollector returns a new Collector exposing file-nr stats.
func NewFileFDStatCollector(logger log.Logger) (Collector, error) {
	return &fileFDStatCollector{logger}, nil
}

func (c *fileFDStatCollector) Update(ch chan<- prometheus.Metric) error {
	fileFDStat, err := parseFileFDStats(procFilePath("sys/fs/file-nr"))
	if err != nil {
		return fmt.Errorf("couldn't get file-nr: %w", err)
	}
	for name, value := range fileFDStat {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value %s in file-nr: %w", value, err)
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, fileFDStatSubsystem, name),
				fmt.Sprintf("File descriptor statistics: %s.", name),
				nil, nil,
			),
			prometheus.GaugeValue, v,
		)
	}
	return nil
}

func parseFileFDStats(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	parts := bytes.Split(bytes.TrimSpace(content), []byte("\u0009"))
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected number of file stats in %q", filename)
	}

	var fileFDStat = map[string]string{}
	// The file-nr proc is only 1 line with 3 values.
	fileFDStat["allocated"] = string(parts[0])
	// The second value is skipped as it will always be zero in linux 2.6.
	fileFDStat["maximum"] = string(parts[2])

	return fileFDStat, nil
}
