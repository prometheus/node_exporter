// Copyright 2022 The Prometheus Authors
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
	"log/slog"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs"
)

var (
	sysctlInclude     = kingpin.Flag("collector.sysctl.include", "Select sysctl metrics to include").Strings()
	sysctlIncludeInfo = kingpin.Flag("collector.sysctl.include-info", "Select sysctl metrics to include as info metrics").Strings()

	sysctlInfoDesc = prometheus.NewDesc(prometheus.BuildFQName(namespace, "sysctl", "info"), "sysctl info", []string{"name", "value", "index"}, nil)
)

type sysctlCollector struct {
	fs      procfs.FS
	logger  *slog.Logger
	sysctls []*sysctl
}

func init() {
	registerCollector("sysctl", defaultDisabled, NewSysctlCollector)
}

func NewSysctlCollector(logger *slog.Logger) (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	c := &sysctlCollector{
		logger:  logger,
		fs:      fs,
		sysctls: []*sysctl{},
	}

	for _, include := range *sysctlInclude {
		sysctl, err := newSysctl(include, true)
		if err != nil {
			return nil, err
		}
		c.sysctls = append(c.sysctls, sysctl)
	}

	for _, include := range *sysctlIncludeInfo {
		sysctl, err := newSysctl(include, false)
		if err != nil {
			return nil, err
		}
		c.sysctls = append(c.sysctls, sysctl)
	}
	return c, nil
}

func (c *sysctlCollector) Update(ch chan<- prometheus.Metric) error {
	for _, sysctl := range c.sysctls {
		metrics, err := c.newMetrics(sysctl)
		if err != nil {
			return err
		}

		for _, metric := range metrics {
			ch <- metric
		}
	}
	return nil
}

func (c *sysctlCollector) newMetrics(s *sysctl) ([]prometheus.Metric, error) {
	var (
		values interface{}
		length int
		err    error
	)

	if s.numeric {
		values, err = c.fs.SysctlInts(s.name)
		if err != nil {
			return nil, fmt.Errorf("error obtaining sysctl info: %w", err)
		}
		length = len(values.([]int))
	} else {
		values, err = c.fs.SysctlStrings(s.name)
		if err != nil {
			return nil, fmt.Errorf("error obtaining sysctl info: %w", err)
		}
		length = len(values.([]string))
	}

	switch length {
	case 0:
		return nil, fmt.Errorf("sysctl %s has no values", s.name)
	case 1:
		if len(s.keys) > 0 {
			return nil, fmt.Errorf("sysctl %s has only one value, but expected %v", s.name, s.keys)
		}
		return []prometheus.Metric{s.newConstMetric(values)}, nil

	default:

		if len(s.keys) == 0 {
			return s.newIndexedMetrics(values), nil
		}

		if length != len(s.keys) {
			return nil, fmt.Errorf("sysctl %s has %d keys but only %d defined in f lag", s.name, length, len(s.keys))
		}

		return s.newMappedMetrics(values)
	}
}

type sysctl struct {
	numeric bool
	name    string
	keys    []string
}

func newSysctl(include string, numeric bool) (*sysctl, error) {
	parts := strings.SplitN(include, ":", 2)
	s := &sysctl{
		numeric: numeric,
		name:    parts[0],
	}
	if len(parts) == 2 {
		s.keys = strings.Split(parts[1], ",")
		s.name = parts[0]
	}
	return s, nil
}

func (s *sysctl) metricName() string {
	return SanitizeMetricName(s.name)
}

func (s *sysctl) newConstMetric(v interface{}) prometheus.Metric {
	if s.numeric {
		return prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "sysctl", s.metricName()),
				fmt.Sprintf("sysctl %s", s.name),
				nil, nil),
			prometheus.UntypedValue,
			float64(v.([]int)[0]))
	}
	return prometheus.MustNewConstMetric(
		sysctlInfoDesc,
		prometheus.GaugeValue,
		1.0,
		s.name,
		v.([]string)[0],
		"0",
	)
}

func (s *sysctl) newIndexedMetrics(v interface{}) []prometheus.Metric {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sysctl", s.metricName()),
		fmt.Sprintf("sysctl %s", s.name),
		[]string{"index"}, nil,
	)
	switch values := v.(type) {
	case []int:
		metrics := make([]prometheus.Metric, len(values))
		for i, n := range values {
			metrics[i] = prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, float64(n), strconv.Itoa(i))
		}
		return metrics
	case []string:
		metrics := make([]prometheus.Metric, len(values))
		for i, str := range values {
			metrics[i] = prometheus.MustNewConstMetric(sysctlInfoDesc, prometheus.GaugeValue, 1.0, s.name, str, strconv.Itoa(i))
		}
		return metrics
	default:
		panic(fmt.Sprintf("unexpected type %T", values))
	}
}

func (s *sysctl) newMappedMetrics(v interface{}) ([]prometheus.Metric, error) {
	switch values := v.(type) {
	case []int:
		metrics := make([]prometheus.Metric, len(values))
		for i, n := range values {
			key := s.keys[i]
			desc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "sysctl", s.metricName()+"_"+key),
				fmt.Sprintf("sysctl %s, field %d", s.name, i),
				nil,
				nil,
			)
			metrics[i] = prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, float64(n))
		}
		return metrics, nil
	case []string:
		return nil, fmt.Errorf("mapped sysctl string values not supported")
	default:
		return nil, fmt.Errorf("unexpected type %T", values)
	}
}
