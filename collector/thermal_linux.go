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

// +build !nonetstat

package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const (
	thermalZoneSubsystem = "thermal"
)

type thermalCollector struct {
	fs   sysfs.FS
	temp *prometheus.Desc
}

func init() {
	registerCollector("thermal", defaultEnabled, NewThermalCollector)
}

// NewThermalCollector takes and returns
// a new Collector exposing thermal stats.
func NewThermalCollector() (Collector, error) {

	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %v", err)
	}

	return &thermalCollector{
		fs: fs,
		temp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, thermalZoneSubsystem, "zone_temp_celsius"),
			"Current thermal zones temperature in celsius degrees.",
			[]string{"zone"}, nil,
		),
	}, nil
}

func (c *thermalCollector) Update(ch chan<- prometheus.Metric) error {
	tzStats, err := c.fs.ClassThermalZoneStats()
	if err != nil {
		return err
	}
	for _, stat := range tzStats {
		ch <- prometheus.MustNewConstMetric(
			c.temp,
			prometheus.GaugeValue,
			float64(stat.Temp)/1000.0,
			stat.Name,
		)
	}
	return nil

}
