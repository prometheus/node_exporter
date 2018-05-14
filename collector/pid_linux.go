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

// +build !nopid

package collector

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	pidStatSubsystem = "pid"
)

type pidCollector struct{}

func init() {
	registerCollector(pidStatSubsystem, defaultEnabled, NewPidStatCollector)
}

// NewPidCollector returns new pidCollector exposing pid usages.
func NewPidStatCollector() (Collector, error) {
	return &pidCollector{}, nil
}

func (c pidCollector) Update(ch chan<- prometheus.Metric) error {
	pidUsed, err := getPidUsage()
	if err != nil {
		return fmt.Errorf("couldn't get PID total: %s", err)
	}

	pidMax, err := getPidMax()
	if err != nil {
		return fmt.Errorf("couldn't get PID max: %s", err)
	}

	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "pid_total", "count"),
		"pid counts.", nil, nil), prometheus.GaugeValue, pidUsed)
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "pid_max", "count"),
		"pid max limit.", nil, nil), prometheus.GaugeValue, pidMax)
	return nil
}

func getPidUsage() (float64, error) {
	var (
		count    int    = 0
		dir_name uint64 = 0
	)
	paths, err := ioutil.ReadDir(*procPath)
	if err != nil {
		return 0, fmt.Errorf("read \"/proc\" directory error")
	}

	for _, s := range paths {
		if !s.IsDir() {
			continue
		}
		dir_name, err = strconv.ParseUint(s.Name(), 10, 64)
		if err != nil {
			log.Debugf("%s", err)
			continue
		}
		if dir_name != 0 {
			count++
		}
	}

	total := float64(count)
	return total, nil
}

func getPidMax() (float64, error) {
	filename := procFilePath("sys/kernel/pid_max")
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseFloat(string(content[:len(content)-1]), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid value %s in pid_max: %s", string(content), err)
	}

	return v, nil
}
