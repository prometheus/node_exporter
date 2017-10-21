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
	"reflect"
	"regexp"
	"strings"

	"github.com/coreos/go-systemd/dbus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	unitWhitelist  = kingpin.Flag("collector.systemd.unit-whitelist", "Regexp of systemd units to whitelist. Units must both match whitelist and not match blacklist to be included.").Default(".+").String()
	unitBlacklist  = kingpin.Flag("collector.systemd.unit-blacklist", "Regexp of systemd units to blacklist. Units must both match whitelist and not match blacklist to be included.").Default(".+\\.scope").String()
	systemdPrivate = kingpin.Flag("collector.systemd.private", "Establish a private, direct connection to systemd without dbus.").Bool()
	accounting     = kingpin.Flag("collector.systemd.accounting", "Whether to expose systemd accounting metrics. These can usually also be scraped using something like cAdvisor, as it's just info from CGroups").Bool()
)

type metric struct {
	desc      string
	valueType prometheus.ValueType
	kind      reflect.Kind
}

type metricsMap map[string]*metric

type systemdCollector struct {
	unitDesc             *prometheus.Desc
	systemRunningDesc    *prometheus.Desc
	unitPropsMetrics     metricsMap
	unitWhitelistPattern *regexp.Regexp
	unitBlacklistPattern *regexp.Regexp
}

var unitStatesName = []string{"active", "activating", "deactivating", "inactive", "failed"}

const subsystem = "systemd"

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(camel string) string {
	snake := matchFirstCap.ReplaceAllString(camel, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func init() {
	registerCollector("systemd", defaultDisabled, NewSystemdCollector)
}

// NewSystemdCollector returns a new Collector exposing systemd statistics.
func NewSystemdCollector() (Collector, error) {

	unitDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_state"),
		"Systemd unit", []string{"name", "state"}, nil,
	)
	systemRunningDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "system_running"),
		"Whether the system is operational (see 'systemctl is-system-running')",
		nil, nil,
	)

	unitMetrics := metricsMap{
		"NRestarts":                       &metric{desc: "Total number of service restart", valueType: prometheus.CounterValue, kind: reflect.Uint32},
		"NAccepted":                       &metric{desc: "Number of accepted connections", valueType: prometheus.CounterValue, kind: reflect.Uint32},
		"NConnections":                    &metric{desc: "Number of open connections", valueType: prometheus.CounterValue, kind: reflect.Uint32},
		"ActiveEnterTimestampMonotonic":   &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ActiveEnterTimestamp":            &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ActiveExitTimestampMonotonic":    &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ActiveExitTimestamp":             &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"AssertTimestampMonotonic":        &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"AssertTimestamp":                 &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ConditionTimestampMonotonic":     &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ConditionTimestamp":              &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ExecMainExitTimestampMonotonic":  &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ExecMainExitTimestamp":           &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ExecMainStartTimestampMonotonic": &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"ExecMainStartTimestamp":          &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"InactiveEnterTimestampMonotonic": &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"InactiveEnterTimestamp":          &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"InactiveExitTimestampMonotonic":  &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"InactiveExitTimestamp":           &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"LastTriggerUSecMonotonic":        &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"LastTriggerUSec":                 &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"NextElapseUSecMonotonic":         &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"NextElapseUSecRealtime":          &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"StateChangeTimestampMonotonic":   &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"StateChangeTimestamp":            &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"WatchdogTimestampMonotonic":      &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		"WatchdogTimestamp":               &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
	}
	if *accounting {
		accountingMetrics := metricsMap{
			"CPUQuotaPerSecUSec": &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"CPUShares":          &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"CPUUsageNSec":       &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"CPUWeight":          &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"IPEgressBytes":      &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"IPEgressPackets":    &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"IPIngressBytes":     &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"IPIngressPackets":   &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"MemoryCurrent":      &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"MemoryHigh":         &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"MemoryLimit":        &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"MemoryLow":          &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"MemoryMax":          &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"MemorySwapMax":      &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"TasksCurrent":       &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
			"TasksMax":           &metric{desc: "", valueType: prometheus.CounterValue, kind: reflect.Uint64},
		}
		for k, v := range accountingMetrics {
			unitMetrics[k] = v
		}
	}
	unitWhitelistPattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitWhitelist))
	unitBlacklistPattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitBlacklist))

	return &systemdCollector{
		unitDesc:             unitDesc,
		systemRunningDesc:    systemRunningDesc,
		unitPropsMetrics:     unitMetrics,
		unitWhitelistPattern: unitWhitelistPattern,
		unitBlacklistPattern: unitBlacklistPattern,
	}, nil
}

func (c *systemdCollector) Update(ch chan<- prometheus.Metric) error {
	units, err := c.listUnits()
	if err != nil {
		return fmt.Errorf("couldn't get units states: %s", err)
	}
	c.collectUnitStatusMetrics(ch, units)
	c.collectUnitProperiesMetrics(ch)

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

func (c *systemdCollector) collectUnitProperiesMetrics(ch chan<- prometheus.Metric) error {
	conn, err := c.newDbus()
	if err != nil {
		return fmt.Errorf("couldn't get dbus connection: %s", err)
	}

	defer conn.Close()

	units, err := conn.ListUnits()
	if err != nil {
		return err
	}

	units = filterUnits(units, c.unitWhitelistPattern, c.unitBlacklistPattern)

	for _, unit := range units {
		splitted := strings.Split(unit.Name, ".")
		unitType := strings.Title(splitted[len(splitted)-1])

		props, err := conn.GetUnitTypeProperties(unit.Name, "Unit")
		if err != nil {
			return err
		}
		unitSpecificProps, err := conn.GetUnitTypeProperties(unit.Name, unitType)

		if err != nil {
			return err
		}
		for k, v := range unitSpecificProps {
			props[k] = v
		}

		for prop, value := range props {
			metric := c.unitPropsMetrics[prop]
			if metric == nil {
				continue
			}
			desc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, toSnakeCase(prop)),
				metric.desc,
				[]string{"name"},
				nil,
			)
			var out float64
			var isOk bool

			// unset values are representedas 0xFFFF... in DBus
			switch metric.kind {
			case reflect.Uint32:
				casted, ok := value.(uint32)
				ok = ok && casted != ^uint32(0)
				isOk = ok
				out = float64(casted)
			case reflect.Uint64:
				casted, ok := value.(uint64)
				ok = ok && casted != ^uint64(0)
				isOk = ok
				out = float64(casted)

			}
			if isOk {
				ch <- prometheus.MustNewConstMetric(desc, metric.valueType, out, unit.Name)
			}

		}
	}
	return nil
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
