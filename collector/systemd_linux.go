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

//go:build !nosystemd
// +build !nosystemd

package collector

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/coreos/go-systemd/dbus"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	unitIncludeSet bool
	unitInclude    = kingpin.Flag("collector.systemd.unit-include", "Regexp of systemd units to include. Units must both match include and not match exclude to be included.").Default(".+").PreAction(func(c *kingpin.ParseContext) error {
		unitIncludeSet = true
		return nil
	}).String()
	oldUnitInclude = kingpin.Flag("collector.systemd.unit-whitelist", "DEPRECATED: Use --collector.systemd.unit-include").Hidden().String()
	unitExcludeSet bool
	unitExclude    = kingpin.Flag("collector.systemd.unit-exclude", "Regexp of systemd units to exclude. Units must both match include and not match exclude to be included.").Default(".+\\.(automount|device|mount|scope|slice)").PreAction(func(c *kingpin.ParseContext) error {
		unitExcludeSet = true
		return nil
	}).String()
	oldUnitExclude         = kingpin.Flag("collector.systemd.unit-blacklist", "DEPRECATED: Use collector.systemd.unit-exclude").Hidden().String()
	systemdPrivate         = kingpin.Flag("collector.systemd.private", "Establish a private, direct connection to systemd without dbus (Strongly discouraged since it requires root. For testing purposes only).").Hidden().Bool()
	enableTaskMetrics      = kingpin.Flag("collector.systemd.enable-task-metrics", "Enables service unit tasks metrics unit_tasks_current and unit_tasks_max").Bool()
	enableRestartsMetrics  = kingpin.Flag("collector.systemd.enable-restarts-metrics", "Enables service unit metric service_restart_total").Bool()
	enableStartTimeMetrics = kingpin.Flag("collector.systemd.enable-start-time-metrics", "Enables service unit metric unit_start_time_seconds").Bool()
)

type systemdCollector struct {
	unitDesc                      *prometheus.Desc
	unitStartTimeDesc             *prometheus.Desc
	unitTasksCurrentDesc          *prometheus.Desc
	unitTasksMaxDesc              *prometheus.Desc
	systemRunningDesc             *prometheus.Desc
	summaryDesc                   *prometheus.Desc
	nRestartsDesc                 *prometheus.Desc
	timerLastTriggerDesc          *prometheus.Desc
	socketAcceptedConnectionsDesc *prometheus.Desc
	socketCurrentConnectionsDesc  *prometheus.Desc
	socketRefusedConnectionsDesc  *prometheus.Desc
	systemdVersionDesc            *prometheus.Desc
	systemdVersion                int
	unitIncludePattern            *regexp.Regexp
	unitExcludePattern            *regexp.Regexp
	logger                        log.Logger
}

var unitStatesName = []string{"active", "activating", "deactivating", "inactive", "failed"}

func init() {
	registerCollector("systemd", defaultDisabled, NewSystemdCollector)
}

// NewSystemdCollector returns a new Collector exposing systemd statistics.
func NewSystemdCollector(logger log.Logger) (Collector, error) {
	const subsystem = "systemd"

	unitDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_state"),
		"Systemd unit", []string{"name", "state", "type"}, nil,
	)
	unitStartTimeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_start_time_seconds"),
		"Start time of the unit since unix epoch in seconds.", []string{"name"}, nil,
	)
	unitTasksCurrentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_tasks_current"),
		"Current number of tasks per Systemd unit", []string{"name"}, nil,
	)
	unitTasksMaxDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_tasks_max"),
		"Maximum number of tasks per Systemd unit", []string{"name"}, nil,
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
		"Service unit count of Restart triggers", []string{"name"}, nil)
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

	if *oldUnitExclude != "" {
		if !unitExcludeSet {
			level.Warn(logger).Log("msg", "--collector.systemd.unit-blacklist is DEPRECATED and will be removed in 2.0.0, use --collector.systemd.unit-exclude")
			*unitExclude = *oldUnitExclude
		} else {
			return nil, errors.New("--collector.systemd.unit-blacklist and --collector.systemd.unit-exclude are mutually exclusive")
		}
	}
	if *oldUnitInclude != "" {
		if !unitIncludeSet {
			level.Warn(logger).Log("msg", "--collector.systemd.unit-whitelist is DEPRECATED and will be removed in 2.0.0, use --collector.systemd.unit-include")
			*unitInclude = *oldUnitInclude
		} else {
			return nil, errors.New("--collector.systemd.unit-whitelist and --collector.systemd.unit-include are mutually exclusive")
		}
	}
	level.Info(logger).Log("msg", "Parsed flag --collector.systemd.unit-include", "flag", *unitInclude)
	unitIncludePattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitInclude))
	level.Info(logger).Log("msg", "Parsed flag --collector.systemd.unit-exclude", "flag", *unitExclude)
	unitExcludePattern := regexp.MustCompile(fmt.Sprintf("^(?:%s)$", *unitExclude))

	systemdVersion := getSystemdVersion(logger)
	if systemdVersion < minSystemdVersionSystemState {
		level.Warn(logger).Log("msg", "Detected systemd version is lower than minimum", "current", systemdVersion, "minimum", minSystemdVersionSystemState)
		level.Warn(logger).Log("msg", "Some systemd state and timer metrics will not be available")
	}

	return &systemdCollector{
		unitDesc:                      unitDesc,
		unitStartTimeDesc:             unitStartTimeDesc,
		unitTasksCurrentDesc:          unitTasksCurrentDesc,
		unitTasksMaxDesc:              unitTasksMaxDesc,
		systemRunningDesc:             systemRunningDesc,
		summaryDesc:                   summaryDesc,
		nRestartsDesc:                 nRestartsDesc,
		timerLastTriggerDesc:          timerLastTriggerDesc,
		socketAcceptedConnectionsDesc: socketAcceptedConnectionsDesc,
		socketCurrentConnectionsDesc:  socketCurrentConnectionsDesc,
		socketRefusedConnectionsDesc:  socketRefusedConnectionsDesc,
		systemdVersionDesc:            systemdVersionDesc,
		systemdVersion:                systemdVersion,
		unitIncludePattern:            unitIncludePattern,
		unitExcludePattern:            unitExcludePattern,
		logger:                        logger,
	}, nil
}

// Update gathers metrics from systemd.  Dbus collection is done in parallel
// to reduce wait time for responses.
func (c *systemdCollector) Update(ch chan<- prometheus.Metric) error {
	begin := time.Now()
	conn, err := newSystemdDbusConn()
	if err != nil {
		return fmt.Errorf("couldn't get dbus connection: %w", err)
	}
	defer conn.Close()

	allUnits, err := c.getAllUnits(conn)
	if err != nil {
		return fmt.Errorf("couldn't get units: %w", err)
	}
	level.Debug(c.logger).Log("msg", "getAllUnits took", "duration_seconds", time.Since(begin).Seconds())

	begin = time.Now()
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

func (c *systemdCollector) collectUnitStatusMetrics(conn *dbus.Conn, ch chan<- prometheus.Metric, units []unit) {
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

func (c *systemdCollector) collectSystemState(conn *dbus.Conn, ch chan<- prometheus.Metric) error {
	systemState, err := conn.GetManagerProperty("SystemState")
	if err != nil {
		return fmt.Errorf("couldn't get system state: %w", err)
	}
	isSystemRunning := 0.0
	if systemState == `"running"` {
		isSystemRunning = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.systemRunningDesc, prometheus.GaugeValue, isSystemRunning)
	return nil
}

func newSystemdDbusConn() (*dbus.Conn, error) {
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

func (c *systemdCollector) getAllUnits(conn *dbus.Conn) ([]unit, error) {
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

func filterUnits(units []unit, includePattern, excludePattern *regexp.Regexp, logger log.Logger) []unit {
	filtered := make([]unit, 0, len(units))
	for _, unit := range units {
		if includePattern.MatchString(unit.Name) && !excludePattern.MatchString(unit.Name) && unit.LoadState == "loaded" {
			level.Debug(logger).Log("msg", "Adding unit", "unit", unit.Name)
			filtered = append(filtered, unit)
		} else {
			level.Debug(logger).Log("msg", "Ignoring unit", "unit", unit.Name)
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
