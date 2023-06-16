// Copyright 2023 The Prometheus Authors
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

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const (
	cpuVulerabilitiesCollector = "cpu_vulnerabilities"
)

var (
	vulnerabilityDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuVulerabilitiesCollector, "info"),
		"Details of each CPU vulnerability reported by sysfs. The value of the series is an int encoded state of the vulnerability. The same state is stored as a string in the label",
		[]string{"codename", "state"},
		nil,
	)
	totalCountDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, cpuVulerabilitiesCollector, "total"),
		"Total number of reported vulnerabilities",
		nil,
		nil,
	)
)

type cpuVulnerabilitiesCollector struct {
	vulnerabilities map[string]*sysfs.Vulnerability
	total           float64
}

func init() {
	registerCollector(cpuVulerabilitiesCollector, defaultDisabled, NewVulnerabilitySysfsCollector)
}

func NewVulnerabilitySysfsCollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	vulnerabilities, err := fs.CPUVulnerabilities()
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities: %w", err)
	}

	return &cpuVulnerabilitiesCollector{
		vulnerabilities: vulnerabilities,
		total:           float64(len(vulnerabilities)), // this assumes the number of cpu vulnerabilities doesn't change during the lifetime of the node_exporter process
	}, nil
}

func (v *cpuVulnerabilitiesCollector) Update(ch chan<- prometheus.Metric) error {
	ch <- prometheus.MustNewConstMetric(totalCountDesc, prometheus.CounterValue, v.total)
	for _, vulnerability := range v.vulnerabilities {
		ch <- prometheus.MustNewConstMetric(
			vulnerabilityDesc,
			prometheus.GaugeValue,
			float64(vulnerability.State),
			vulnerability.CodeName,
			sysfs.VulnerabilityHumanEncoding[vulnerability.State],
		)
	}
	return nil
}
