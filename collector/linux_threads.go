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

// build+ !nothreads

package collector

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type threadsCollector struct {
	threadAlloc *prometheus.Desc
}
func init() {
	registerCollector("threads", defaultDisabled, NewThreadsCollector)
}

func NewThreadsCollector() (Collector, error) {
	return &threadsCollector{
		threadAlloc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace,"", "threads"),
			"Allocated thread",
			nil, nil,
		),
	}, nil
}
func (t *threadsCollector) Update(ch chan<- prometheus.Metric) error {
	val, err := readProcessStatus()
	if err != nil {
		return fmt.Errorf("Unable to retrieve number of threads %v\n", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadAlloc, prometheus.GaugeValue, float64(val))
	return nil
}

func readProcessStatus() (int, error) {
	processDir, err := regexp.Compile("([0-9]){1,8}")
	if err != nil {
		return 0, err
	}
	folders, err := ioutil.ReadDir("/proc")
	if err != nil {
		return 0, err
	}
	threads := 0
	for _, f := range folders {
		if f.IsDir() && processDir.MatchString(f.Name()) {
			file, err := ioutil.ReadFile("/proc/" + f.Name() + "/status")
			if err != nil {
				return 0, err
			}
			line := strings.Split(string(file), "\n")
			if err != nil {
				return 0, err
			}
			for _, l := range line {
				if strings.Contains(string(l), "Threads:") {
					threadStr := strings.Split(string(l), ":")
					tread, err :=  strconv.Atoi(threadStr[1])
					if err != nil {
						return 0, err
					}
					threads += tread
				}
			}
		}
	}
	return threads, nil
}