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
	"flag"
	"fmt"
	"regexp"

	"github.com/coreos/go-systemd/dbus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	unitWhitelist = flag.String("collector.systemd.unit-whitelist", ".+", "Regexp of systemd units to whitelist. Units must both match whitelist and not match blacklist to be included.")
	unitBlacklist = flag.String("collector.systemd.unit-blacklist", "", "Regexp of systemd units to blacklist. Units must both match whitelist and not match blacklist to be included.")
)

type systemdCollector struct {
	unitDesc             *prometheus.Desc
	systemRunningDesc    *prometheus.Desc
	unitWhitelistPattern *regexp.Regexp
	unitBlacklistPattern *regexp.Regexp
}

var unitStatesName = []string{"active", "activating", "deactivating", "inactive", "failed"}

var (
	systemdPrivate = flag.Bool(
		"collector.systemd.private",
		false,
		"Establish a private, direct connection to systemd without dbus.",
	)
)

func init() {
	Factories["systemd"] = NewSystemdCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// systemd statistics.
func NewSystemdCollector() (Collector, error) {
	const subsystem = "systemd"

	unitDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "unit_state"),
		"Systemd unit", []string{"name", "state"}, nil,
	)
	systemRunningDesc := prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, "system_running"),
		"Whether the system is operational (see 'systemctl is-system-running')",
		nil, nil,
	)
	unitWhitelistPattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitWhitelist))
	unitBlacklistPattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitBlacklist))

	return &systemdCollector{
		unitDesc:             unitDesc,
		systemRunningDesc:    systemRunningDesc,
		unitWhitelistPattern: unitWhitelistPattern,
		unitBlacklistPattern: unitBlacklistPattern,
	}, nil
}

func (c *systemdCollector) Update(ch chan<- prometheus.Metric) (err error) {
	units, err := c.listUnits()
	if err != nil {
		return fmt.Errorf("couldn't get units states: %s", err)
	}
	c.collectUnitStatusMetrics(ch, units)

	systemState, err := c.getSystemState()
	if err != nil {
		return fmt.Errorf("couldn't get system state: %s", err)
	}
	c.collectSystemState(ch, systemState)

	return nil
}

func (c *systemdCollector) collectUnitStatusMetrics(ch chan<- prometheus.Metric, units []dbus.UnitStatus) {
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

func (c *systemdCollector) collectSystemState(ch chan<- prometheus.Metric, systemState string) {
	isSystemRunning := 0.0
	if systemState == `"running"` {
		isSystemRunning = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.systemRunningDesc, prometheus.GaugeValue, isSystemRunning)
}

func (c *systemdCollector) newDbus() (*dbus.Conn, error) {
	if *systemdPrivate {
		return dbus.NewSystemdConnection()
	}
	return dbus.New()
}

func (c *systemdCollector) listUnits() ([]dbus.UnitStatus, error) {
	conn, err := c.newDbus()
	if err != nil {
		return nil, fmt.Errorf("couldn't get dbus connection: %s", err)
	}
	allUnits, err := conn.ListUnits()
	conn.Close()

	if err != nil {
		return []dbus.UnitStatus{}, err
	}

	units := filterUnits(allUnits, c.unitWhitelistPattern, c.unitBlacklistPattern)
	return units, nil
}

func filterUnits(units []dbus.UnitStatus, whitelistPattern, blacklistPattern *regexp.Regexp) []dbus.UnitStatus {
	filtered := make([]dbus.UnitStatus, 0, len(units))
	for _, unit := range units {
		if whitelistPattern.MatchString(unit.Name) && !blacklistPattern.MatchString(unit.Name) {
			filtered = append(filtered, unit)
		} else {
			log.Debugf("Ignoring unit: %s", unit.Name)
		}
	}

	return filtered
}

func (c *systemdCollector) getSystemState() (state string, err error) {
	conn, err := c.newDbus()
	if err != nil {
		return "", fmt.Errorf("couldn't get dbus connection: %s", err)
	}
	state, err = conn.GetManagerProperty("SystemState")
	conn.Close()
	return state, err
}
