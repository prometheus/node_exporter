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
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-systemd/dbus"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	// minSystemdVersionSystemState is the minimum SystemD version for availability of
	// the 'SystemState' manager property and the timer property 'LastTriggerUSec'
	// https://github.com/prometheus/node_exporter/issues/291
	minSystemdVersionSystemState = 212
)

var (
	unitInclude            = kingpin.Flag("collector.systemd.unit-include", "Regexp of systemd units to include. Units must both match include and not match exclude to be included.").Default(".+").String()
	oldUnitInclude         = kingpin.Flag("collector.systemd.unit-whitelist", "DEPRECATED: Use --collector.systemd.unit-include").Hidden().String()
	unitExclude            = kingpin.Flag("collector.systemd.unit-exclude", "Regexp of systemd units to exclude. Units must both match include and not match exclude to be included.").Default(".+\\.(automount|device|mount|scope|slice)").String()
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
	systemdVersionDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "version"),
		"Detected systemd version", []string{}, nil)

	if *oldUnitExclude != "" {
		if *unitExclude == "" {
			level.Warn(logger).Log("msg", "--collector.systemd.unit-blacklist is DEPRECATED and will be removed in 2.0.0, use --collector.systemd.unit-exclude")
			*unitExclude = *oldUnitExclude
		} else {
			return nil, errors.New("--collector.systemd.unit-blacklist and --collector.systemd.unit-exclude are mutually exclusive")
		}
	}
	if *oldUnitInclude != "" {
		if *unitInclude == "" {
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
	level.Debug(c.logger).Log("msg", "collectSummaryMetrics took", "duration_seconds", time.Since(begin).Seconds())

	begin = time.Now()
	units := filterUnits(allUnits, c.unitIncludePattern, c.unitExcludePattern, c.logger)
	level.Debug(c.logger).Log("msg", "filterUnits took", "duration_seconds", time.Since(begin).Seconds())

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		begin = time.Now()
		c.collectUnitStatusMetrics(conn, ch, units)
		level.Debug(c.logger).Log("msg", "collectUnitStatusMetrics took", "duration_seconds", time.Since(begin).Seconds())
	}()

	if *enableStartTimeMetrics {
		wg.Add(1)
		go func() {
			defer wg.Done()
			begin = time.Now()
			c.collectUnitStartTimeMetrics(conn, ch, units)
			level.Debug(c.logger).Log("msg", "collectUnitStartTimeMetrics took", "duration_seconds", time.Since(begin).Seconds())
		}()
	}

	if *enableTaskMetrics {
		wg.Add(1)
		go func() {
			defer wg.Done()
			begin = time.Now()
			c.collectUnitTasksMetrics(conn, ch, units)
			level.Debug(c.logger).Log("msg", "collectUnitTasksMetrics took", "duration_seconds", time.Since(begin).Seconds())
		}()
	}

	if c.systemdVersion >= minSystemdVersionSystemState {
		wg.Add(1)
		go func() {
			defer wg.Done()
			begin = time.Now()
			c.collectTimers(conn, ch, units)
			level.Debug(c.logger).Log("msg", "collectTimers took", "duration_seconds", time.Since(begin).Seconds())
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		begin = time.Now()
		c.collectSockets(conn, ch, units)
		level.Debug(c.logger).Log("msg", "collectSockets took", "duration_seconds", time.Since(begin).Seconds())
	}()

	if c.systemdVersion >= minSystemdVersionSystemState {
		begin = time.Now()
		err = c.collectSystemState(conn, ch)
		level.Debug(c.logger).Log("msg", "collectSystemState took", "duration_seconds", time.Since(begin).Seconds())
	}

	ch <- prometheus.MustNewConstMetric(
		c.systemdVersionDesc, prometheus.GaugeValue, float64(c.systemdVersion))

	return err
}

func (c *systemdCollector) collectUnitStatusMetrics(conn *dbus.Conn, ch chan<- prometheus.Metric, units []unit) {
	for _, unit := range units {
		serviceType := ""
		if strings.HasSuffix(unit.Name, ".service") {
			serviceTypeProperty, err := conn.GetUnitTypeProperty(unit.Name, "Service", "Type")
			if err != nil {
				level.Debug(c.logger).Log("msg", "couldn't get unit type", "unit", unit.Name, "err", err)
			} else {
				serviceType = serviceTypeProperty.Value.Value().(string)
			}
		} else if strings.HasSuffix(unit.Name, ".mount") {
			serviceTypeProperty, err := conn.GetUnitTypeProperty(unit.Name, "Mount", "Type")
			if err != nil {
				level.Debug(c.logger).Log("msg", "couldn't get unit type", "unit", unit.Name, "err", err)
			} else {
				serviceType = serviceTypeProperty.Value.Value().(string)
			}
		}
		for _, stateName := range unitStatesName {
			isActive := 0.0
			if stateName == unit.ActiveState {
				isActive = 1.0
			}
			ch <- prometheus.MustNewConstMetric(
				c.unitDesc, prometheus.GaugeValue, isActive,
				unit.Name, stateName, serviceType)
		}
		if *enableRestartsMetrics && strings.HasSuffix(unit.Name, ".service") {
			// NRestarts wasn't added until systemd 235.
			restartsCount, err := conn.GetUnitTypeProperty(unit.Name, "Service", "NRestarts")
			if err != nil {
				level.Debug(c.logger).Log("msg", "couldn't get unit NRestarts", "unit", unit.Name, "err", err)
			} else {
				ch <- prometheus.MustNewConstMetric(
					c.nRestartsDesc, prometheus.CounterValue,
					float64(restartsCount.Value.Value().(uint32)), unit.Name)
			}
		}
	}
}

func (c *systemdCollector) collectSockets(conn *dbus.Conn, ch chan<- prometheus.Metric, units []unit) {
	for _, unit := range units {
		if !strings.HasSuffix(unit.Name, ".socket") {
			continue
		}

		acceptedConnectionCount, err := conn.GetUnitTypeProperty(unit.Name, "Socket", "NAccepted")
		if err != nil {
			level.Debug(c.logger).Log("msg", "couldn't get unit NAccepted", "unit", unit.Name, "err", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.socketAcceptedConnectionsDesc, prometheus.CounterValue,
			float64(acceptedConnectionCount.Value.Value().(uint32)), unit.Name)

		currentConnectionCount, err := conn.GetUnitTypeProperty(unit.Name, "Socket", "NConnections")
		if err != nil {
			level.Debug(c.logger).Log("msg", "couldn't get unit NConnections", "unit", unit.Name, "err", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.socketCurrentConnectionsDesc, prometheus.GaugeValue,
			float64(currentConnectionCount.Value.Value().(uint32)), unit.Name)

		// NRefused wasn't added until systemd 239.
		refusedConnectionCount, err := conn.GetUnitTypeProperty(unit.Name, "Socket", "NRefused")
		if err != nil {
			//log.Debugf("couldn't get unit '%s' NRefused: %s", unit.Name, err)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.socketRefusedConnectionsDesc, prometheus.GaugeValue,
				float64(refusedConnectionCount.Value.Value().(uint32)), unit.Name)
		}
	}
}

func (c *systemdCollector) collectUnitStartTimeMetrics(conn *dbus.Conn, ch chan<- prometheus.Metric, units []unit) {
	var startTimeUsec uint64

	for _, unit := range units {
		if unit.ActiveState != "active" {
			startTimeUsec = 0
		} else {
			timestampValue, err := conn.GetUnitProperty(unit.Name, "ActiveEnterTimestamp")
			if err != nil {
				level.Debug(c.logger).Log("msg", "couldn't get unit StartTimeUsec", "unit", unit.Name, "err", err)
				continue
			}
			startTimeUsec = timestampValue.Value.Value().(uint64)
		}

		ch <- prometheus.MustNewConstMetric(
			c.unitStartTimeDesc, prometheus.GaugeValue,
			float64(startTimeUsec)/1e6, unit.Name)
	}
}

func (c *systemdCollector) collectUnitTasksMetrics(conn *dbus.Conn, ch chan<- prometheus.Metric, units []unit) {
	var val uint64
	for _, unit := range units {
		if strings.HasSuffix(unit.Name, ".service") {
			tasksCurrentCount, err := conn.GetUnitTypeProperty(unit.Name, "Service", "TasksCurrent")
			if err != nil {
				level.Debug(c.logger).Log("msg", "couldn't get unit TasksCurrent", "unit", unit.Name, "err", err)
			} else {
				val = tasksCurrentCount.Value.Value().(uint64)
				// Don't set if tasksCurrent if dbus reports MaxUint64.
				if val != math.MaxUint64 {
					ch <- prometheus.MustNewConstMetric(
						c.unitTasksCurrentDesc, prometheus.GaugeValue,
						float64(val), unit.Name)
				}
			}
			tasksMaxCount, err := conn.GetUnitTypeProperty(unit.Name, "Service", "TasksMax")
			if err != nil {
				level.Debug(c.logger).Log("msg", "couldn't get unit TasksMax", "unit", unit.Name, "err", err)
			} else {
				val = tasksMaxCount.Value.Value().(uint64)
				// Don't set if tasksMax if dbus reports MaxUint64.
				if val != math.MaxUint64 {
					ch <- prometheus.MustNewConstMetric(
						c.unitTasksMaxDesc, prometheus.GaugeValue,
						float64(val), unit.Name)
				}
			}
		}
	}
}

func (c *systemdCollector) collectTimers(conn *dbus.Conn, ch chan<- prometheus.Metric, units []unit) {
	for _, unit := range units {
		if !strings.HasSuffix(unit.Name, ".timer") {
			continue
		}

		lastTriggerValue, err := conn.GetUnitTypeProperty(unit.Name, "Timer", "LastTriggerUSec")
		if err != nil {
			level.Debug(c.logger).Log("msg", "couldn't get unit LastTriggerUSec", "unit", unit.Name, "err", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			c.timerLastTriggerDesc, prometheus.GaugeValue,
			float64(lastTriggerValue.Value.Value().(uint64))/1e6, unit.Name)
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

func getSystemdVersion(logger log.Logger) int {
	conn, err := newSystemdDbusConn()
	if err != nil {
		level.Warn(logger).Log("msg", "Unable to get systemd dbus connection, defaulting systemd version to 0", "err", err)
		return 0
	}
	defer conn.Close()
	version, err := conn.GetManagerProperty("Version")
	if err != nil {
		level.Warn(logger).Log("msg", "Unable to get systemd version property, defaulting to 0")
		return 0
	}
	re := regexp.MustCompile(`[0-9][0-9][0-9]`)
	version = re.FindString(version)
	v, err := strconv.Atoi(version)
	if err != nil {
		level.Warn(logger).Log("msg", "Got invalid systemd version", "version", version)
		return 0
	}
	return v
}
