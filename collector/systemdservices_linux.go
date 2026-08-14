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

package collector

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("systemdservices", defaultDisabled, NewSystemdServicesCollector)
	registerCollector("systemdinfo", defaultDisabled, NewSystemdInfoCollector)
}

func listSystemdUnits(conn *dbus.Conn) ([]dbus.UnitStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return conn.ListUnitsContext(ctx)
}

func closeSystemdConn(conn **dbus.Conn) error {
	if conn != nil && *conn != nil {
		(*conn).Close()
		*conn = nil
	}
	return nil
}

func isServiceUnit(name string) bool {
	return strings.HasSuffix(name, ".service")
}

// systemdServicesCollector is health: ActiveState only (failed at 15s).
type systemdServicesCollector struct {
	serviceState *prometheus.Desc
	logger       *slog.Logger
	conn         *dbus.Conn
}

func NewSystemdServicesCollector(logger *slog.Logger) (Collector, error) {
	conn, err := newSystemdDbusConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get dbus connection: %w", err)
	}

	return &systemdServicesCollector{
		serviceState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "systemd_service", "state"),
			"Systemd service state: 0 = unknown, 1 = active, 2 = reloading, 3 = inactive, 4 = failed, 5 = activating, 6 = deactivating.",
			[]string{"name"},
			nil,
		),
		logger: logger,
		conn:   conn,
	}, nil
}

func (c *systemdServicesCollector) Update(ch chan<- prometheus.Metric) error {
	units, err := listSystemdUnits(c.conn)
	if err != nil {
		return fmt.Errorf("couldn't get units: %w", err)
	}
	for _, unit := range units {
		if !isServiceUnit(unit.Name) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.serviceState,
			prometheus.GaugeValue,
			parseSystemdState(unit.ActiveState),
			unit.Name,
		)
	}
	return nil
}

func (c *systemdServicesCollector) Close() error {
	return closeSystemdConn(&c.conn)
}

// systemdInfoCollector is analysis: Type, load/sub state, and NRestarts.
type systemdInfoCollector struct {
	serviceInfo         *prometheus.Desc
	serviceSubState     *prometheus.Desc
	serviceLoadState    *prometheus.Desc
	serviceRestartTotal *prometheus.Desc
	logger              *slog.Logger
	conn                *dbus.Conn
}

func NewSystemdInfoCollector(logger *slog.Logger) (Collector, error) {
	conn, err := newSystemdDbusConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get dbus connection: %w", err)
	}
	return &systemdInfoCollector{
		serviceInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "systemd_service", "info"),
			"Static systemd service information via D-Bus API. Value is always 1.",
			[]string{"name", "type"},
			nil,
		),
		serviceSubState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "systemd_service", "sub_state"),
			"Systemd service sub-state: 0 = unknown, 1 = running, 2 = exited, 3 = failed, 4 = dead, 5 = start, 6 = stop, 7 = reload, 8 = auto-restart.",
			[]string{"name"},
			nil,
		),
		serviceLoadState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "systemd_service", "load_state"),
			"Systemd service load state: 0 = unknown, 1 = loaded, 2 = error, 3 = masked, 4 = not-found.",
			[]string{"name"},
			nil,
		),
		serviceRestartTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "systemd_service", "restart_total"),
			"Total number of restart triggers for the service unit (systemd Service NRestarts).",
			[]string{"name"},
			nil,
		),
		logger: logger,
		conn:   conn,
	}, nil
}

func (c *systemdInfoCollector) Update(ch chan<- prometheus.Metric) error {
	units, err := listSystemdUnits(c.conn)
	if err != nil {
		return fmt.Errorf("couldn't get units: %w", err)
	}
	for _, unit := range units {
		if !isServiceUnit(unit.Name) {
			continue
		}
		c.collectAnalysisMetrics(ch, unit)
	}
	return nil
}

func (c *systemdInfoCollector) collectAnalysisMetrics(ch chan<- prometheus.Metric, unit dbus.UnitStatus) {
	serviceType := "unknown"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	typeProperty, err := c.conn.GetUnitTypePropertyContext(ctx, unit.Name, "Service", "Type")
	cancel()
	if err == nil {
		if v, ok := typeProperty.Value.Value().(string); ok && v != "" {
			serviceType = v
		}
	}
	ch <- prometheus.MustNewConstMetric(c.serviceInfo, prometheus.GaugeValue, 1, unit.Name, serviceType)
	ch <- prometheus.MustNewConstMetric(c.serviceSubState, prometheus.GaugeValue, parseSystemdSubState(unit.SubState), unit.Name)
	ch <- prometheus.MustNewConstMetric(c.serviceLoadState, prometheus.GaugeValue, parseSystemdLoadState(unit.LoadState), unit.Name)

	// NRestarts wasn't added until systemd 235; older versions return an error (logged at Debug).
	restartCtx, restartCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer restartCancel()
	restartsCount, err := c.conn.GetUnitTypePropertyContext(restartCtx, unit.Name, "Service", "NRestarts")
	if err != nil {
		c.logger.Debug("couldn't get unit NRestarts", "unit", unit.Name, "err", err)
		return
	}
	raw := restartsCount.Value.Value()
	fv, ok := dbusNumericToFloat64(raw)
	if !ok {
		c.logger.Debug("unexpected NRestarts value type", "unit", unit.Name, "type", fmt.Sprintf("%T", raw))
		return
	}
	ch <- prometheus.MustNewConstMetric(c.serviceRestartTotal, prometheus.CounterValue, fv, unit.Name)
}

func (c *systemdInfoCollector) Close() error {
	return closeSystemdConn(&c.conn)
}

func parseSystemdState(state string) float64 {
	switch strings.ToLower(state) {
	case "active":
		return 1
	case "reloading":
		return 2
	case "inactive":
		return 3
	case "failed":
		return 4
	case "activating":
		return 5
	case "deactivating":
		return 6
	default:
		return 0
	}
}

func parseSystemdSubState(subState string) float64 {
	switch strings.ToLower(subState) {
	case "running":
		return 1
	case "exited":
		return 2
	case "failed":
		return 3
	case "dead":
		return 4
	case "start":
		return 5
	case "stop":
		return 6
	case "reload":
		return 7
	case "auto-restart":
		return 8
	default:
		return 0
	}
}

func dbusNumericToFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}

func parseSystemdLoadState(loadState string) float64 {
	switch strings.ToLower(loadState) {
	case "loaded":
		return 1
	case "error":
		return 2
	case "masked":
		return 3
	case "not-found":
		return 4
	default:
		return 0
	}
}
