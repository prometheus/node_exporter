// Copyright 2020 The Prometheus Authors
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

package collector

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

const zoneinfoSubsystem = "zoneinfo"

type zoneinfoCollector struct {
	metricDescs map[string]*prometheus.Desc
	logger      log.Logger
	fs          procfs.FS
}

func init() {
	registerCollector("zoneinfo", defaultDisabled, NewZoneinfoCollector)
}

// NewZoneinfoCollector returns a new Collector exposing zone stats.
func NewZoneinfoCollector(logger log.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	return &zoneinfoCollector{
		metricDescs: map[string]*prometheus.Desc{},
		logger:      logger,
		fs:          fs,
	}, nil
}

func (c *zoneinfoCollector) Update(ch chan<- prometheus.Metric) error {
	metrics, err := c.fs.Zoneinfo()
	if err != nil {
		return fmt.Errorf("couldn't get zoneinfo: %w", err)
	}
	for _, metric := range metrics {
		node := metric.Node
		zone := metric.Zone
		metricStruct := reflect.ValueOf(metric)
		typeOfMetricStruct := metricStruct.Type()
		for i := 0; i < metricStruct.NumField(); i++ {
			value := reflect.Indirect(metricStruct.Field(i))
			if value.Kind() != reflect.Int64 {
				continue
			}
			metricName := toSnakeCase(typeOfMetricStruct.Field(i).Name)
			desc, ok := c.metricDescs[metricName]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, zoneinfoSubsystem, metricName),
					fmt.Sprintf("Zoneinfo information field %s.", metricName),
					[]string{"node", "zone"}, nil)
				c.metricDescs[metricName] = desc
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
				float64(reflect.Indirect(metricStruct.Field(i)).Int()),
				node, zone)
		}
		for i, value := range metric.Protection {
			metricName := fmt.Sprintf("protection_%d", i)
			desc, ok := c.metricDescs[metricName]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, zoneinfoSubsystem, metricName),
					fmt.Sprintf("Zoneinfo information field %s.", metricName),
					[]string{"node", "zone"}, nil)
				c.metricDescs[metricName] = desc
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue,
				float64(*value), node, zone)
		}

	}
	return nil
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
