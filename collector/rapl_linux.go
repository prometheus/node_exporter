// Copyright 2019 The Prometheus Authors
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

// +build !norapl

package collector

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var raplDir = "/sys/class/powercap/"

func init() {
	// disable this collector by default, if rapl files are not present
	defaultRegistration := defaultDisabled

	_, err := ioutil.ReadDir(raplDir)
	if err == nil {
		defaultRegistration = defaultEnabled
	}

	registerCollector("rapl", defaultRegistration, NewRaplCollector)
}

type raplMetric struct {
	microjouleFileName string           // path to the sysfs file with the microjoule value
	descriptor         *prometheus.Desc // prometheus descriptor
	time               time.Time        // last time metric was read
	microjoules        int64            // last read microjoule value from sysfs
	maxMicrojoules     int64            // max microjoule value for the sysfs value
}

type raplCollector struct {
	metrics []raplMetric
}

// NewRaplCollector returns a new Collector exposing RAPL metrics.
func NewRaplCollector() (Collector, error) {
	collector := raplCollector{}

	files, err := ioutil.ReadDir(raplDir)
	if err != nil {
		return nil, errors.New("no sysfs powercap / RAPL power metrics files found")
	}

	countNameUsages := make(map[string]int)

	// loop through directory files searching for file "name" from subdirs
	for _, f := range files {
		nameFile := filepath.Join(raplDir, f.Name(), "/name")
		dat, err := ioutil.ReadFile(nameFile)
		if err == nil { // add new metric since name file was found
			name := strings.TrimSpace(string(dat))

			// make RAPL name prometheus-compatible
			modifiedPromName := strings.Replace(name, "-", "_", -1)

			// store into map how many times this name has been used
			// there can be e.g. multiple "dram" instances, which prometheus will not accept without unique names
			usedName := modifiedPromName
			count, ok := countNameUsages[usedName]
			if ok {
				// make it unique with a "_" + count
				modifiedPromName = modifiedPromName + "_" + strconv.Itoa(count)
				count++
			} else {
				count = 1 // first use
			}
			countNameUsages[usedName] = count

			prometheusDescriptor := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "rapl", modifiedPromName+"_watts"),
				"Current RAPL value in watts",
				nil, nil,
			)

			maxMicrojouleFileName := filepath.Join(raplDir, f.Name(), "/max_energy_range_uj")
			microjouleFileName := filepath.Join(raplDir, f.Name(), "/energy_uj")

			metric := raplMetric{
				microjouleFileName: microjouleFileName,
				descriptor:         prometheusDescriptor,
				time:               time.Now(),
				microjoules:        readRaplValue(microjouleFileName),
				maxMicrojoules:     readRaplValue(maxMicrojouleFileName),
			}

			collector.metrics = append(collector.metrics, metric)
		}
	}

	return &collector, nil
}

func readRaplValue(fileName string) int64 {
	raplValue, err := ioutil.ReadFile(fileName)
	if err == nil {
		value, err2 := strconv.ParseInt(strings.TrimSpace(string(raplValue)), 10, 64)
		if err2 == nil {
			return value
		}

		log.Errorf("numeric value (" + string(raplValue) + ") parsing error:" + err2.Error())
	} else {
		log.Errorf("error reading file \"" + fileName + "\"")
	}

	return -1
}

func calcMicrojouleDifference(metric *raplMetric, newMicrojoules int64) int64 {
	oldMicrojoules := metric.microjoules
	if newMicrojoules < oldMicrojoules {
		// counter must have overflowed, so first take the part before overflow
		deltaMicrojoules := metric.maxMicrojoules - oldMicrojoules
		// then add new value
		deltaMicrojoules += newMicrojoules

		return deltaMicrojoules
	}

	return newMicrojoules - oldMicrojoules // no overflow
}

// Update implements Collector and exposes RAPL related metrics.
func (c *raplCollector) Update(ch chan<- prometheus.Metric) error {
	if c != nil {
		metricCount := len(c.metrics)
		for i := 0; i < metricCount; i++ {
			metric := &c.metrics[i]

			now := time.Now()
			interval := now.Sub(metric.time)

			newMicrojoules := readRaplValue(metric.microjouleFileName)
			deltaMicrojoules := calcMicrojouleDifference(metric, newMicrojoules)

			power := float64(deltaMicrojoules) / 1000000.0 / interval.Seconds()

			// store values for next round
			metric.microjoules = newMicrojoules
			metric.time = now

			ch <- prometheus.MustNewConstMetric(
				metric.descriptor,
				prometheus.GaugeValue,
				power,
			)
		}
		return nil
	}

	return errors.New("nil collector given")
}
