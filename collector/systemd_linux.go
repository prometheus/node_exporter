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

// +build !nosystemd

package collector

import (
	"fmt"

	"github.com/coreos/go-systemd/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

type systemdCollector struct {
	unitDesc *prometheus.Desc
}

var unitStatesName = []string{"active", "activating", "deactivating", "inactive", "failed"}

func init() {
	Factories["systemd"] = NewSystemdCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// systemd statistics.
func NewSystemdCollector() (Collector, error) {
	unitDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "systemd", "unit_state"),
		"Systemd unit", []string{"name", "state"}, nil,
	)

	return &systemdCollector{
		unitDesc: unitDesc,
	}, nil
}

func (c *systemdCollector) Update(ch chan<- prometheus.Metric) (err error) {
	units, err := c.listUnits()
	if err != nil {
		return fmt.Errorf("couldn't get units states: %s", err)
	}

	c.collectMetrics(ch, units)

	return nil
}

func (c *systemdCollector) collectMetrics(ch chan<- prometheus.Metric, units []dbus.UnitStatus) {
	for _, unit := range units {
		for _, stateName := range unitStatesName {
			isActive := 0.0
			if stateName == unit.ActiveState {
				isActive = 1.0
			}
			ch <- prometheus.MustNewConstMetric(
				c.unitDesc, prometheus.GaugeValue, isActive,
				unit.Name, stateName)
		}
	}
}

func (c *systemdCollector) listUnits() ([]dbus.UnitStatus, error) {
	conn, err := dbus.New()
	if err != nil {
		return nil, fmt.Errorf("couldn't get dbus connection: %s", err)
	}
	units, err := conn.ListUnits()
	conn.Close()
	return units, err
}
