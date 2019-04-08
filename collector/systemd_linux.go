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
	"math"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-systemd/dbus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	unitWhitelist              = kingpin.Flag("collector.systemd.unit-whitelist", "Regexp of systemd units to whitelist. Units must both match whitelist and not match blacklist to be included.").Default(".+").String()
	unitBlacklist              = kingpin.Flag("collector.systemd.unit-blacklist", "Regexp of systemd units to blacklist. Units must both match whitelist and not match blacklist to be included.").Default(".+\\.(automount|device|mount|scope|slice)").String()
	systemdPrivate             = kingpin.Flag("collector.systemd.private", "Establish a private, direct connection to systemd without dbus.").Bool()
	enableResourceUsageMetrics = kingpin.Flag("collector.systemd.enable-resource-usage-metrics", "Enables service unit resource usage metrics (CPU, memory, tasks, IP)").Bool()
	enableRestartsMetrics      = kingpin.Flag("collector.systemd.enable-restarts-metrics", "Enables service unit metric service_restart_total").Bool()
	enableStartTimeMetrics     = kingpin.Flag("collector.systemd.enable-start-time-metrics", "Enables service unit metric unit_start_time_seconds").Bool()
)

type systemdCollector struct {
	unitDesc                      *prometheus.Desc
	unitStartTimeDesc             *prometheus.Desc
	unitCPUUsageDesc              *prometheus.Desc
	unitCPUWeightDesc             *prometheus.Desc
	unitStartupCPUWeightDesc      *prometheus.Desc
	unitCPUQuotaDesc              *prometheus.Desc
	unitMemoryCurrentDesc         *prometheus.Desc
	unitMemoryMinDesc             *prometheus.Desc
	unitMemoryLowDesc             *prometheus.Desc
	unitMemoryHighDesc            *prometheus.Desc
	unitMemoryMaxDesc             *prometheus.Desc
	unitMemorySwapMaxDesc         *prometheus.Desc
	unitTasksCurrentDesc          *prometheus.Desc
	unitTasksMaxDesc              *prometheus.Desc
	unitIPIngressBytesDesc        *prometheus.Desc
	unitIPIngressPacketsDesc      *prometheus.Desc
	unitIPEgressBytesDesc         *prometheus.Desc
	unitIPEgressPacketsDesc       *prometheus.Desc
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
		"Systemd unit", []string{"name", "state", "type"}, nil,
	)
	unitStartTimeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_start_time_seconds"),
		"Start time of the unit since unix epoch in seconds.", []string{"name"}, nil,
	)
	unitCPUUsageDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_cpu_usage_seconds"),
		"CPU usage per Systemd unit", []string{"name"}, nil,
	)
	unitCPUWeightDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_cpu_weight"),
		"CPU time weight (between 1 and 10000) per Systemd unit", []string{"name"}, nil,
	)
	unitStartupCPUWeightDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_startup_cpu_weight"),
		"Startup CPU time weight (between 1 and 10000) per Systemd unit", []string{"name"}, nil,
	)
	unitCPUQuotaDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_cpu_quota"),
		"CPU time quota per Systemd unit", []string{"name"}, nil,
	)
	unitMemoryCurrentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_memory_current_bytes"),
		"Current memory usage per Systemd unit", []string{"name"}, nil,
	)
	unitMemoryMinDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_memory_min_bytes"),
		"Memory usage protection per Systemd unit", []string{"name"}, nil,
	)
	unitMemoryLowDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_memory_low_bytes"),
		"Best-effort memory usage protection per Systemd unit", []string{"name"}, nil,
	)
	unitMemoryHighDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_memory_high_bytes"),
		"High memory limit per Systemd unit", []string{"name"}, nil,
	)
	unitMemoryMaxDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_memory_max_bytes"),
		"Absolute limit on memory usage per Systemd unit", []string{"name"}, nil,
	)
	unitMemorySwapMaxDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_memory_swap_max_bytes"),
		"Absolute limit on swap usage per Systemd unit", []string{"name"}, nil,
	)
	unitTasksCurrentDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_tasks_current"),
		"Current number of tasks per Systemd unit", []string{"name"}, nil,
	)
	unitTasksMaxDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_tasks_max"),
		"Maximum number of tasks per Systemd unit", []string{"name"}, nil,
	)
	unitIPIngressBytesDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_ip_ingress_bytes"),
		"IP ingress bytes per Systemd unit", []string{"name"}, nil,
	)
	unitIPIngressPacketsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_ip_ingress_packets"),
		"IP ingress packets per Systemd unit", []string{"name"}, nil,
	)
	unitIPEgressBytesDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_ip_egress_bytes"),
		"IP egress bytes per Systemd unit", []string{"name"}, nil,
	)
	unitIPEgressPacketsDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "unit_ip_egress_packets"),
		"IP egress packets per Systemd unit", []string{"name"}, nil,
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
		unitCPUUsageDesc:              unitCPUUsageDesc,
		unitCPUWeightDesc:             unitCPUWeightDesc,
		unitStartupCPUWeightDesc:      unitStartupCPUWeightDesc,
		unitCPUQuotaDesc:              unitCPUQuotaDesc,
		unitMemoryCurrentDesc:         unitMemoryCurrentDesc,
		unitMemoryMinDesc:             unitMemoryMinDesc,
		unitMemoryLowDesc:             unitMemoryLowDesc,
		unitMemoryHighDesc:            unitMemoryHighDesc,
		unitMemoryMaxDesc:             unitMemoryMaxDesc,
		unitMemorySwapMaxDesc:         unitMemorySwapMaxDesc,
		unitTasksCurrentDesc:          unitTasksCurrentDesc,
		unitTasksMaxDesc:              unitTasksMaxDesc,
		unitIPIngressBytesDesc:        unitIPIngressBytesDesc,
		unitIPIngressPacketsDesc:      unitIPIngressPacketsDesc,
		unitIPEgressBytesDesc:         unitIPEgressBytesDesc,
		unitIPEgressPacketsDesc:       unitIPEgressPacketsDesc,
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

// Update gathers metrics from systemd.  Dbus collection is done in parallel
// to reduce wait time for responses.
func (c *systemdCollector) Update(ch chan<- prometheus.Metric) error {
	begin := time.Now()
	conn, err := c.newDbus()
	if err != nil {
		return fmt.Errorf("couldn't get dbus connection: %s", err)
	}
	defer conn.Close()

	allUnits, err := c.getAllUnits(conn)
	if err != nil {
		return fmt.Errorf("couldn't get units: %s", err)
	}
	log.Debugf("systemd getAllUnits took %f", time.Since(begin).Seconds())

	begin = time.Now()
	summary := summarizeUnits(allUnits)
	c.collectSummaryMetrics(ch, summary)
	log.Debugf("systemd collectSummaryMetrics took %f", time.Since(begin).Seconds())

	begin = time.Now()
	units := filterUnits(allUnits, c.unitWhitelistPattern, c.unitBlacklistPattern)
	log.Debugf("systemd filterUnits took %f", time.Since(begin).Seconds())

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		begin = time.Now()
		c.collectUnitStatusMetrics(conn, ch, units)
		log.Debugf("systemd collectUnitStatusMetrics took %f", time.Since(begin).Seconds())
	}()

	if *enableStartTimeMetrics {
		wg.Add(1)
		go func() {
			defer wg.Done()
			begin = time.Now()
			c.collectUnitStartTimeMetrics(conn, ch, units)
			log.Debugf("systemd collectUnitStartTimeMetrics took %f", time.Since(begin).Seconds())
		}()
	}

	if *enableResourceUsageMetrics {
		wg.Add(1)
		go func() {
			defer wg.Done()
			begin = time.Now()
			c.collectUnitResourceUsageMetrics(conn, ch, units)
			log.Debugf("systemd collectUnitResourceUsageMetrics took %f", time.Since(begin).Seconds())
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		begin = time.Now()
		c.collectTimers(conn, ch, units)
		log.Debugf("systemd collectTimers took %f", time.Since(begin).Seconds())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		begin = time.Now()
		c.collectSockets(conn, ch, units)
		log.Debugf("systemd collectSockets took %f", time.Since(begin).Seconds())
	}()

	begin = time.Now()
	err = c.collectSystemState(conn, ch)
	log.Debugf("systemd collectSystemState took %f", time.Since(begin).Seconds())
	return err
}

func (c *systemdCollector) collectUnitStatusMetrics(conn *dbus.Conn, ch chan<- prometheus.Metric, units []unit) {
	for _, unit := range units {
		serviceType := ""
		if strings.HasSuffix(unit.Name, ".service") {
			serviceTypeProperty, err := conn.GetUnitTypeProperty(unit.Name, "Service", "Type")
			if err != nil {
				log.Debugf("couldn't get unit '%s' Type: %s", unit.Name, err)
			} else {
				serviceType = serviceTypeProperty.Value.Value().(string)
			}
		} else if strings.HasSuffix(unit.Name, ".mount") {
			serviceTypeProperty, err := conn.GetUnitTypeProperty(unit.Name, "Mount", "Type")
			if err != nil {
				log.Debugf("couldn't get unit '%s' Type: %s", unit.Name, err)
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
				log.Debugf("couldn't get unit '%s' NRestarts: %s", unit.Name, err)
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
			log.Debugf("couldn't get unit '%s' NAccepted: %s", unit.Name, err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.socketAcceptedConnectionsDesc, prometheus.CounterValue,
			float64(acceptedConnectionCount.Value.Value().(uint32)), unit.Name)

		currentConnectionCount, err := conn.GetUnitTypeProperty(unit.Name, "Socket", "NConnections")
		if err != nil {
			log.Debugf("couldn't get unit '%s' NConnections: %s", unit.Name, err)
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
				log.Debugf("couldn't get unit '%s' StartTimeUsec: %s", unit.Name, err)
				continue
			}
			startTimeUsec = timestampValue.Value.Value().(uint64)
		}

		ch <- prometheus.MustNewConstMetric(
			c.unitStartTimeDesc, prometheus.GaugeValue,
			float64(startTimeUsec)/1e6, unit.Name)
	}
}

func getServicePropertyUint64(conn *dbus.Conn, unitName, propertyName string, ch chan<- prometheus.Metric, desc *prometheus.Desc) bool {
	if property, err := conn.GetUnitTypeProperty(unitName, "Service", propertyName); err != nil {
		log.Debugf("couldn't get unit '%s' %s: %s", unitName, propertyName, err)
		return false
	} else if val := property.Value.Value().(uint64); val != math.MaxUint64 {
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(val), unitName)
	}
	return true
}

func (c *systemdCollector) collectUnitResourceUsageMetrics(conn *dbus.Conn, ch chan<- prometheus.Metric, units []unit) {
	for _, unit := range units {
		if strings.HasSuffix(unit.Name, ".service") {
			// CPU
			getServicePropertyUint64(conn, unit.Name, "CPUUsageNSec", ch, c.unitCPUUsageDesc)
			// CPUWeight replaces CPUShares
			if !getServicePropertyUint64(conn, unit.Name, "CPUWeight", ch, c.unitCPUWeightDesc) {
				getServicePropertyUint64(conn, unit.Name, "CPUShares", ch, c.unitCPUWeightDesc)
			}
			// StartupCPUWeight replaces StartupCPUShares
			if !getServicePropertyUint64(conn, unit.Name, "StartupCPUWeight", ch, c.unitStartupCPUWeightDesc) {
				getServicePropertyUint64(conn, unit.Name, "StartupCPUShares", ch, c.unitStartupCPUWeightDesc)
			}
			getServicePropertyUint64(conn, unit.Name, "CPUQuota", ch, c.unitCPUQuotaDesc)

			// Memory
			getServicePropertyUint64(conn, unit.Name, "MemoryCurrent", ch, c.unitMemoryCurrentDesc)
			getServicePropertyUint64(conn, unit.Name, "MemoryMin", ch, c.unitMemoryMinDesc)
			getServicePropertyUint64(conn, unit.Name, "MemoryLow", ch, c.unitMemoryLowDesc)
			getServicePropertyUint64(conn, unit.Name, "MemoryHigh", ch, c.unitMemoryHighDesc)
			// MemoryMax replaces MemoryLimit
			if !getServicePropertyUint64(conn, unit.Name, "MemoryMax", ch, c.unitMemoryHighDesc) {
				getServicePropertyUint64(conn, unit.Name, "MemoryLimit", ch, c.unitMemoryHighDesc)
			}
			getServicePropertyUint64(conn, unit.Name, "MemorySwapMax", ch, c.unitMemorySwapMaxDesc)

			// Tasks
			getServicePropertyUint64(conn, unit.Name, "TasksCurrent", ch, c.unitTasksCurrentDesc)
			getServicePropertyUint64(conn, unit.Name, "TasksMax", ch, c.unitTasksMaxDesc)

			// IP
			getServicePropertyUint64(conn, unit.Name, "IPIngressBytes", ch, c.unitIPIngressBytesDesc)
			getServicePropertyUint64(conn, unit.Name, "IPIngressPackets", ch, c.unitIPIngressPacketsDesc)
			getServicePropertyUint64(conn, unit.Name, "IPEgressBytes", ch, c.unitIPEgressBytesDesc)
			getServicePropertyUint64(conn, unit.Name, "IPEgressPackets", ch, c.unitIPEgressPacketsDesc)
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
			log.Debugf("couldn't get unit '%s' LastTriggerUSec: %s", unit.Name, err)
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
		return fmt.Errorf("couldn't get system state: %s", err)
	}
	isSystemRunning := 0.0
	if systemState == `"running"` {
		isSystemRunning = 1.0
	}
	ch <- prometheus.MustNewConstMetric(c.systemRunningDesc, prometheus.GaugeValue, isSystemRunning)
	return nil
}

func (c *systemdCollector) newDbus() (*dbus.Conn, error) {
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
