// Copyright The Prometheus Authors
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

//go:build !nopowermeter

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const powerMeterCollectorSubsystem = "power_meter"

var (
	powerMeterIgnoredMeters = kingpin.Flag(
		"collector.power-meter.ignored-meters",
		"Regexp of ACPI power meter names to ignore for power_meter collector.",
	).Default("^$").String()
)

type powerMeterCollector struct {
	ignoredPattern *regexp.Regexp
	logger         *slog.Logger
}

func init() {
	registerCollector(powerMeterCollectorSubsystem, defaultEnabled, NewPowerMeterCollector)
}

// NewPowerMeterCollector returns a new Collector exposing ACPI 4.0 power meter
// metrics read from /sys/bus/acpi/drivers/power_meter/ACPI000D:XX/.
func NewPowerMeterCollector(logger *slog.Logger) (Collector, error) {
	pattern, err := regexp.Compile(*powerMeterIgnoredMeters)
	if err != nil {
		return nil, err
	}
	return &powerMeterCollector{
		ignoredPattern: pattern,
		logger:         logger,
	}, nil
}

// pushPowerMeterMetric is a helper that builds a single gauge metric under the
// node_power_meter_<name> family and writes it to ch. The `meter` label is
// always set to the ACPI device name (e.g. "ACPI000D:00").
func pushPowerMeterMetric(ch chan<- prometheus.Metric, name string, value float64, meter string, valueType prometheus.ValueType) {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, powerMeterCollectorSubsystem, name),
		"ACPI power meter "+name+" value.",
		[]string{"meter"},
		nil,
	)
	ch <- prometheus.MustNewConstMetric(desc, valueType, value, meter)
}

// Update implements Collector and emits ACPI power meter metrics.
func (c *powerMeterCollector) Update(ch chan<- prometheus.Metric) error {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return fmt.Errorf("failed to open sysfs: %w", err)
	}

	meters, err := sysfs.GetPowerMeters(fs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("ACPI power meter bus path not present; collector idle")
			return ErrNoData
		}
		if errors.Is(err, os.ErrPermission) {
			c.logger.Debug("cannot access ACPI power meter sysfs path", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("failed to read ACPI power meters: %w", err)
	}

	for _, pm := range meters {
		if c.ignoredPattern.MatchString(pm.Name) {
			c.logger.Debug("ignoring ACPI power meter", "meter", pm.Name)
			continue
		}
		c.emitPowerMeterMetrics(ch, pm)
	}
	return nil
}

// emitPowerMeterMetrics writes every metric for a single PowerMeter.
func (c *powerMeterCollector) emitPowerMeterMetrics(ch chan<- prometheus.Metric, pm sysfs.PowerMeter) {
	meter := pm.Name

	// Power gauges (µW -> W).
	for name, v := range map[string]*int64{
		"average_watts":     pm.Average,
		"average_min_watts": pm.AverageMin,
		"average_max_watts": pm.AverageMax,
	} {
		if v != nil {
			pushPowerMeterMetric(ch, name, float64(*v)/1e6, meter, prometheus.GaugeValue)
		}
	}

	// Cap gauges (µW -> W).
	for name, v := range map[string]*int64{
		"cap_watts":            pm.Cap,
		"cap_min_watts":        pm.CapMin,
		"cap_max_watts":        pm.CapMax,
		"cap_hysteresis_watts": pm.CapHyst,
	} {
		if v != nil {
			pushPowerMeterMetric(ch, name, float64(*v)/1e6, meter, prometheus.GaugeValue)
		}
	}

	// Interval gauges (ms -> s).
	for name, v := range map[string]*int64{
		"average_interval_seconds":     pm.AverageInterval,
		"average_interval_min_seconds": pm.AverageIntervalMin,
		"average_interval_max_seconds": pm.AverageIntervalMax,
	} {
		if v != nil {
			pushPowerMeterMetric(ch, name, float64(*v)/1e3, meter, prometheus.GaugeValue)
		}
	}

	// State gauges (0/1, no unit conversion).
	for name, v := range map[string]*int64{
		"alarm":      pm.Alarm,
		"is_battery": pm.IsBattery,
	} {
		if v != nil {
			pushPowerMeterMetric(ch, name, float64(*v), meter, prometheus.GaugeValue)
		}
	}

	// Accuracy (percent string "1.50%" -> 1.50).
	if pm.Accuracy != "" {
		if percent, err := parseAccuracy(pm.Accuracy); err == nil {
			pushPowerMeterMetric(ch, "accuracy_percent", percent, meter, prometheus.GaugeValue)
		} else {
			c.logger.Debug("could not parse power meter accuracy", "meter", meter, "accuracy", pm.Accuracy, "err", err)
		}
	}

	// Info metric: meter metadata as labels, value = 1.
	infoLabels := []string{"meter"}
	infoValues := []string{meter}
	for _, kv := range []struct{ name, value string }{
		{"model_number", pm.ModelNumber},
		{"serial_number", pm.SerialNumber},
		{"oem_info", pm.OEMInfo},
	} {
		if kv.value != "" {
			infoLabels = append(infoLabels, kv.name)
			infoValues = append(infoValues, strings.ToValidUTF8(kv.value, ""))
		}
	}
	infoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, powerMeterCollectorSubsystem, "info"),
		"ACPI power meter metadata.",
		infoLabels,
		nil,
	)
	ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, 1.0, infoValues...)

	// Measures info: one series per device this meter measures.
	if len(pm.Measures) > 0 {
		measuresDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, powerMeterCollectorSubsystem, "measures_info"),
			"Device measured by this ACPI power meter.",
			[]string{"meter", "device"},
			nil,
		)
		for _, dev := range pm.Measures {
			ch <- prometheus.MustNewConstMetric(measuresDesc, prometheus.GaugeValue, 1.0, meter, dev)
		}
	}
}

// parseAccuracy converts a sysfs accuracy string like "1.50%" to a percentage
// value (1.50). Returns an error if the string cannot be parsed.
func parseAccuracy(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	return strconv.ParseFloat(s, 64)
}
