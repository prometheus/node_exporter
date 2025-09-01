package collector

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

type systemdServicesCollector struct {
	serviceMetrics *prometheus.Desc
	logger         *slog.Logger
	conn           *dbus.Conn
}

func init() {
	registerCollector("systemdservices", defaultDisabled, NewSystemdServicesCollector)
}

func NewSystemdServicesCollector(logger *slog.Logger) (Collector, error) {
	conn, err := newSystemdDbusConn()
	if err != nil {
		return nil, fmt.Errorf("couldn't get dbus connection: %w", err)
	}

	return &systemdServicesCollector{
		serviceMetrics: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "systemd_services", "info"),
			"Systemd service status information via D-Bus API. Value is 1 if service is active, 0 otherwise. "+
				"States include: active, reloading, inactive, failed, activating, deactivating. "+
				"Types can be: simple, forking, oneshot, dbus, notify, idle. "+
				"Load states: loaded, error, masked, not-found. "+
				"Sub states: running, exited, failed, dead, start, stop, reload, auto-restart.",
			[]string{"name", "state", "type", "load_state", "sub_state"},
			nil,
		),
		logger: logger,
		conn:   conn,
	}, nil
}

func (c *systemdServicesCollector) Update(ch chan<- prometheus.Metric) error {
	units, err := c.getAllUnits(c.conn)
	if err != nil {
		return fmt.Errorf("couldn't get units: %w", err)
	}

	for _, unit := range units {
		if !strings.HasSuffix(unit.Name, ".service") {
			continue
		}

		if err := c.collectServiceMetrics(c.conn, ch, unit); err != nil {
			c.logger.Debug("failed to collect metrics for unit", "unit", unit.Name, "error", err)
			continue
		}
	}

	return nil
}

func (c *systemdServicesCollector) getAllUnits(conn *dbus.Conn) ([]dbus.UnitStatus, error) {
	units, err := conn.ListUnitsContext(context.TODO())
	if err != nil {
		return nil, err
	}
	return units, nil
}

func (c *systemdServicesCollector) collectServiceMetrics(conn *dbus.Conn, ch chan<- prometheus.Metric, unit dbus.UnitStatus) error {
	serviceType := "unknown"
	typeProperty, err := conn.GetUnitTypePropertyContext(context.TODO(), unit.Name, "Service", "Type")
	if err == nil {
		serviceType = typeProperty.Value.Value().(string)
	}

	ch <- prometheus.MustNewConstMetric(
		c.serviceMetrics,
		prometheus.GaugeValue,
		1,
		unit.Name,
		unit.ActiveState,
		serviceType,
		unit.LoadState,
		unit.SubState,
	)
	return nil
}

func (c *systemdServicesCollector) Close() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	return nil
}
