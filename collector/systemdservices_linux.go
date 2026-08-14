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
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("systemdservices", defaultDisabled, NewSystemdServicesCollector)
	registerCollector("systemdinfo", defaultDisabled, NewSystemdInfoCollector)
}

var errSystemdConnClosed = errors.New("systemd dbus connection closed")

const systemdInfoUnitDeadline = 10 * time.Second

// systemdConn serializes get/close of a long-lived dbus connection.
type systemdConn struct {
	mu     sync.Mutex
	conn   *dbus.Conn
	closed bool
	dial   func() (*dbus.Conn, error)
}

func (s *systemdConn) get() (*dbus.Conn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, errSystemdConnClosed
	}
	if s.conn != nil {
		return s.conn, nil
	}
	dial := s.dial
	if dial == nil {
		dial = newSystemdDbusConn
	}
	conn, err := dial()
	if err != nil {
		return nil, fmt.Errorf("couldn't get dbus connection: %w", err)
	}
	s.conn = conn
	return conn, nil
}

func (s *systemdConn) drop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return
	}
	s.conn.Close()
	s.conn = nil
}

func (s *systemdConn) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	if s.conn == nil {
		return nil
	}
	s.conn.Close()
	s.conn = nil
	return nil
}

func (s *systemdConn) connAndUnits() (*dbus.Conn, []dbus.UnitStatus, error) {
	conn, err := s.get()
	if err != nil {
		return nil, nil, err
	}
	units, err := listSystemdUnits(conn)
	if err == nil {
		return conn, units, nil
	}
	s.drop()
	conn, err2 := s.get()
	if err2 != nil {
		return nil, nil, err2
	}
	units, err = listSystemdUnits(conn)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get units: %w", err)
	}
	return conn, units, nil
}

func listSystemdUnits(conn *dbus.Conn) ([]dbus.UnitStatus, error) {
	if conn == nil {
		return nil, errSystemdConnClosed
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return conn.ListUnitsContext(ctx)
}

func isServiceUnit(name string) bool {
	return strings.HasSuffix(name, ".service")
}

// systemdServicesCollector is health: ActiveState only (failed at 15s).
type systemdServicesCollector struct {
	serviceState *prometheus.Desc
	logger       *slog.Logger
	sc           systemdConn
}

func NewSystemdServicesCollector(logger *slog.Logger) (Collector, error) {
	return &systemdServicesCollector{
		serviceState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "systemd_service", "state"),
			"Systemd service state: 0 = unknown, 1 = active, 2 = reloading, 3 = inactive, 4 = failed, 5 = activating, 6 = deactivating.",
			[]string{"name"},
			nil,
		),
		logger: logger,
	}, nil
}

func (c *systemdServicesCollector) Update(ch chan<- prometheus.Metric) error {
	_, units, err := c.sc.connAndUnits()
	if err != nil {
		return err
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
	return c.sc.Close()
}

// systemdInfoCollector is analysis: Type, load/sub state, and NRestarts.
type systemdInfoCollector struct {
	serviceInfo         *prometheus.Desc
	serviceSubState     *prometheus.Desc
	serviceLoadState    *prometheus.Desc
	serviceRestartTotal *prometheus.Desc
	logger              *slog.Logger
	sc                  systemdConn
}

func NewSystemdInfoCollector(logger *slog.Logger) (Collector, error) {
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
	}, nil
}

func (c *systemdInfoCollector) Update(ch chan<- prometheus.Metric) error {
	conn, units, err := c.sc.connAndUnits()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), systemdInfoUnitDeadline)
	defer cancel()
	for _, unit := range units {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !isServiceUnit(unit.Name) {
			continue
		}
		c.collectAnalysisMetrics(ctx, ch, conn, unit)
	}
	return nil
}

func (c *systemdInfoCollector) collectAnalysisMetrics(ctx context.Context, ch chan<- prometheus.Metric, conn *dbus.Conn, unit dbus.UnitStatus) {
	serviceType := "unknown"
	props, err := conn.GetUnitTypePropertiesContext(ctx, unit.Name, "Service")
	if err != nil {
		if ctx.Err() == nil {
			c.logger.Debug("couldn't get unit service properties", "unit", unit.Name, "err", err)
		}
	} else if t := serviceTypeFromProps(props); t != "" {
		serviceType = t
	}
	ch <- prometheus.MustNewConstMetric(c.serviceInfo, prometheus.GaugeValue, 1, unit.Name, serviceType)
	ch <- prometheus.MustNewConstMetric(c.serviceSubState, prometheus.GaugeValue, parseSystemdSubState(unit.SubState), unit.Name)
	ch <- prometheus.MustNewConstMetric(c.serviceLoadState, prometheus.GaugeValue, parseSystemdLoadState(unit.LoadState), unit.Name)

	if err != nil {
		return
	}
	// NRestarts wasn't added until systemd 235; absence is expected on older versions.
	fv, ok := nRestartsFromProps(props)
	if !ok {
		if raw, exists := props["NRestarts"]; exists {
			c.logger.Debug("unexpected NRestarts value type", "unit", unit.Name, "type", fmt.Sprintf("%T", raw))
		}
		return
	}
	ch <- prometheus.MustNewConstMetric(c.serviceRestartTotal, prometheus.CounterValue, fv, unit.Name)
}

func serviceTypeFromProps(props map[string]any) string {
	v, ok := props["Type"].(string)
	if !ok {
		return ""
	}
	return v
}

func nRestartsFromProps(props map[string]any) (float64, bool) {
	raw, ok := props["NRestarts"]
	if !ok {
		return 0, false
	}
	return dbusNumericToFloat64(raw)
}

func (c *systemdInfoCollector) Close() error {
	return c.sc.Close()
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
