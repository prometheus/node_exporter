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

type systemdCollector struct {
	unitDesc                      *prometheus.Desc
	unitStartTimeDesc             *prometheus.Desc
	systemRunningDesc             *prometheus.Desc
	summaryDesc                   *prometheus.Desc
	nRestartsDesc                 *prometheus.Desc
	timerLastTriggerDesc          *prometheus.Desc
	socketAcceptedConnectionsDesc *prometheus.Desc
	socketCurrentConnectionsDesc  *prometheus.Desc
	socketRefusedConnectionsDesc  *prometheus.Desc
	unitWhitelistPattern          *regexp.Regexp
	unitBlacklistPattern          *regexp.Regexp
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
	unitStartTimeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_start_time_seconds"),
		"Start time of the unit since unix epoch in seconds.", []string{"name"}, nil,
	)
	systemRunningDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "system_running"),
		"Whether the system is operational (see 'systemctl is-system-running')",
		nil, nil,
	)
	summaryDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "units"),
		"Summary of systemd unit states", []string{"state"}, nil)
	nRestartsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "service_restart_total"),
		"Service unit count of Restart triggers", []string{"state"}, nil)
	timerLastTriggerDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "timer_last_trigger_seconds"),
		"Seconds since epoch of last trigger.", []string{"name"}, nil)
	socketAcceptedConnectionsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "socket_accepted_connections_total"),
		"Total number of accepted socket connections", []string{"name"}, nil)
	socketCurrentConnectionsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "socket_current_connections"),
		"Current number of socket connections", []string{"name"}, nil)
	socketRefusedConnectionsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "socket_refused_connections_total"),
		"Total number of refused socket connections", []string{"name"}, nil)
	unitWhitelistPattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitWhitelist))
	unitBlacklistPattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitBlacklist))

	return &systemdCollector{
		unitDesc:                      unitDesc,
		unitStartTimeDesc:             unitStartTimeDesc,
		systemRunningDesc:             systemRunningDesc,
		summaryDesc:                   summaryDesc,
		nRestartsDesc:                 nRestartsDesc,
		timerLastTriggerDesc:          timerLastTriggerDesc,
		socketAcceptedConnectionsDesc: socketAcceptedConnectionsDesc,
		socketCurrentConnectionsDesc:  socketCurrentConnectionsDesc,
		socketRefusedConnectionsDesc:  socketRefusedConnectionsDesc,
		unitWhitelistPattern:          unitWhitelistPattern,
		unitBlacklistPattern:          unitBlacklistPattern,
	}, nil
}

func (c *systemdCollector) Update(ch chan<- prometheus.Metric) error {
	allUnits, err := c.getAllUnits()
	if err != nil {
		return fmt.Errorf("couldn't get units: %s", err)
	}

	summary := summarizeUnits(allUnits)
	c.collectSummaryMetrics(ch, summary)

	units := filterUnits(allUnits, c.unitWhitelistPattern, c.unitBlacklistPattern)
	c.collectUnitStatusMetrics(ch, units)
	c.collectUnitStartTimeMetrics(ch, units)
	c.collectTimers(ch, units)
	c.collectSockets(ch, units)

	systemState, err := c.getSystemState()
	if err != nil {
		return fmt.Errorf("couldn't get system state: %s", err)
	}
	c.collectSystemState(ch, systemState)

	return nil
}

func (c *systemdCollector) collectUnitStatusMetrics(ch chan<- prometheus.Metric, units []unit) {
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
		if strings.HasSuffix(unit.Name, ".service") && unit.nRestarts != nil {
			ch <- prometheus.MustNewConstMetric(
				c.nRestartsDesc, prometheus.CounterValue,
				float64(*unit.nRestarts), unit.Name)
		}
	}
}

func (c *systemdCollector) collectSockets(ch chan<- prometheus.Metric, units []unit) {
	for _, unit := range units {
		if !strings.HasSuffix(unit.Name, ".socket") {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			c.socketAcceptedConnectionsDesc, prometheus.CounterValue,
			float64(unit.acceptedConnections), unit.Name)
		ch <- prometheus.MustNewConstMetric(
			c.socketCurrentConnectionsDesc, prometheus.GaugeValue,
			float64(unit.currentConnections), unit.Name)
		if unit.refusedConnections != nil {
			ch <- prometheus.MustNewConstMetric(
				c.socketRefusedConnectionsDesc, prometheus.GaugeValue,
				float64(*unit.refusedConnections), unit.Name)
		}
	}
}

func (c *systemdCollector) collectUnitStartTimeMetrics(ch chan<- prometheus.Metric, units []unit) {
	for _, unit := range units {
		ch <- prometheus.MustNewConstMetric(
			c.unitStartTimeDesc, prometheus.GaugeValue,
			float64(unit.startTimeUsec)/1e6, unit.Name)
	}
}

func (c *systemdCollector) collectTimers(ch chan<- prometheus.Metric, units []unit) {
	for _, unit := range units {
		if !strings.HasSuffix(unit.Name, ".timer") {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			c.timerLastTriggerDesc, prometheus.GaugeValue,
			float64(unit.lastTriggerUsec)/1e6, unit.Name)
	}
}

func (c *systemdCollector) collectSummaryMetrics(ch chan<- prometheus.Metric, summary map[string]float64) {
	for stateName, count := range summary {
		ch <- prometheus.MustNewConstMetric(
			c.summaryDesc, prometheus.GaugeValue, count, stateName)
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

type unit struct {
	dbus.UnitStatus
	lastTriggerUsec     uint64
	startTimeUsec       uint64
	nRestarts           *uint32
	acceptedConnections uint32
	currentConnections  uint32
	refusedConnections  *uint32
}

func (c *systemdCollector) getAllUnits() ([]unit, error) {
	conn, err := c.newDbus()
	if err != nil {
		return nil, fmt.Errorf("couldn't get dbus connection: %s", err)
	}
	defer conn.Close()

	// Filter out any units that are not installed and are pulled in only as dependencies.
	allUnits, err := conn.ListUnits()

	if err != nil {
		return nil, err
	}

	result := make([]unit, 0, len(allUnits))
	for _, status := range allUnits {
		unit := unit{
			UnitStatus: status,
		}

		if strings.HasSuffix(unit.Name, ".timer") {
			lastTriggerValue, err := conn.GetUnitTypeProperty(unit.Name, "Timer", "LastTriggerUSec")
			if err != nil {
				log.Debugf("couldn't get unit '%s' LastTriggerUSec: %s", unit.Name, err)
				continue
			}

			unit.lastTriggerUsec = lastTriggerValue.Value.Value().(uint64)
		}
		if strings.HasSuffix(unit.Name, ".service") {
			// NRestarts wasn't added until systemd 235.
			restartsCount, err := conn.GetUnitTypeProperty(unit.Name, "Service", "NRestarts")
			if err != nil {
				log.Debugf("couldn't get unit '%s' NRestarts: %s", unit.Name, err)
			} else {
				nRestarts := restartsCount.Value.Value().(uint32)
				unit.nRestarts = &nRestarts
			}
		}

		if strings.HasSuffix(unit.Name, ".socket") {
			acceptedConnectionCount, err := conn.GetUnitTypeProperty(unit.Name, "Socket", "NAccepted")
			if err != nil {
				log.Debugf("couldn't get unit '%s' NAccepted: %s", unit.Name, err)
				continue
			}

			unit.acceptedConnections = acceptedConnectionCount.Value.Value().(uint32)

			currentConnectionCount, err := conn.GetUnitTypeProperty(unit.Name, "Socket", "NConnections")
			if err != nil {
				log.Debugf("couldn't get unit '%s' NConnections: %s", unit.Name, err)
				continue
			}
			unit.currentConnections = currentConnectionCount.Value.Value().(uint32)

			// NRefused wasn't added until systemd 239.
			refusedConnectionCount, err := conn.GetUnitTypeProperty(unit.Name, "Socket", "NRefused")
			if err != nil {
				log.Debugf("couldn't get unit '%s' NRefused: %s", unit.Name, err)
			} else {
				nRefused := refusedConnectionCount.Value.Value().(uint32)
				unit.refusedConnections = &nRefused
			}
		}

		if unit.ActiveState != "active" {
			unit.startTimeUsec = 0
		} else {
			timestampValue, err := conn.GetUnitProperty(unit.Name, "ActiveEnterTimestamp")
			if err != nil {
				log.Debugf("couldn't get unit '%s' StartTimeUsec: %s", unit.Name, err)
				continue
			}

			unit.startTimeUsec = timestampValue.Value.Value().(uint64)
		}

		result = append(result, unit)
	}

	return result, nil
}

func summarizeUnits(units []unit) map[string]float64 {
	summarized := make(map[string]float64)

	for _, unitStateName := range unitStatesName {
		summarized[unitStateName] = 0.0
	}

	for _, unit := range units {
		summarized[unit.ActiveState] += 1.0
	}

	return summarized
}

func filterUnits(units []unit, whitelistPattern, blacklistPattern *regexp.Regexp) []unit {
	filtered := make([]unit, 0, len(units))
	for _, unit := range units {
		if whitelistPattern.MatchString(unit.Name) && !blacklistPattern.MatchString(unit.Name) && unit.LoadState == "loaded" {
			log.Debugf("Adding unit: %s", unit.Name)
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
