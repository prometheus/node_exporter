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
)

type metric struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
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

func init() {
	registerCollector("systemd", defaultDisabled, NewSystemdCollector)
}

// NewSystemdCollector returns a new Collector exposing systemd statistics.
func NewSystemdCollector() (Collector, error) {
	const subsystem = "systemd"

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
		"CPUUsageNSec": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "cpu_usage_nanoseconds_total"),
				"Total CPU seconds of a unit",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"MemoryCurrent": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "memory_current"),
				"Amount of bytes",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.GaugeValue,
		},
		"TasksCurrent": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "tasks_current"),
				"amount of tasks. Includes both user processes and kernel threads.",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.GaugeValue,
		},
		"IPIngressBytes": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "ip_ingress_bytes_total"),
				"Ingress bytes total",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"IPIngressPackets": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "ip_ingress_packets_total"),
				"Ingress packets total",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"IPEgressBytes": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "ip_egress_bytes_total"),
				"Egress bytes total",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"IPEgressPackets": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "ip_egress_packets_total"),
				"Egress packets total",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"NRestarts": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "nrestarts"),
				"Total number of service restarts",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"AssertTimestampMonotonic": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "nrestarts"),
				"Total number of service restarts",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"ConditionTimestampMonotonic": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "nrestarts"),
				"Total number of service restarts",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"InactiveEnterTimestampMonotonic": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "nrestarts"),
				"Total number of service restarts",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"InactiveExitTimestampMonotonic": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "nrestarts"),
				"Total number of service restarts",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"ActiveEnterTimestampMonotonic": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "nrestarts"),
				"Total number of service restarts",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"ActiveExitTimestampMonotonic": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "nrestarts"),
				"Total number of service restarts",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
		"StateChangeTimestampMonotonic": &metric{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "nrestarts"),
				"Total number of service restarts",
				[]string{"name"},
				nil,
			),
			valueType: prometheus.CounterValue,
		},
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

	for _, unit := range units {
		splitted := strings.Split(unit.Name, ".")
		unitType := strings.Title(splitted[1])

		props, err := conn.GetUnitTypeProperties(unit.Name, unitType)

                if err != nil {
                  return err
                }

		for prop, value := range props {
                        metric := c.unitPropsMetrics[prop]
			if metric != nil {
				value, ok := value.(uint64)
				if ok {
					if value != ^uint64(0) {
						ch <- prometheus.MustNewConstMetric(metric.desc, metric.valueType, float64(value), unit.Name)
					}
				}
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
