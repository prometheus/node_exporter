// Copyright 2019 The Prometheus Authors
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

//go:build !nopowersupplyclass && linux
// +build !nopowersupplyclass,linux

package collector

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	powerSupplyClassIgnoredPowerSupplies = kingpin.Flag("collector.powersupply.ignored-supplies", "Regexp of power supplies to ignore for powersupplyclass collector.").Default("^$").String()
)

type powerSupplyClassCollector struct {
	subsystem      string
	ignoredPattern *regexp.Regexp
	metricDescs    map[string]*prometheus.Desc
	logger         log.Logger
}

func init() {
	registerCollector("powersupplyclass", defaultEnabled, NewPowerSupplyClassCollector)
}

func NewPowerSupplyClassCollector(logger log.Logger) (Collector, error) {
	pattern := regexp.MustCompile(*powerSupplyClassIgnoredPowerSupplies)
	return &powerSupplyClassCollector{
		subsystem:      "power_supply",
		ignoredPattern: pattern,
		metricDescs:    map[string]*prometheus.Desc{},
		logger:         logger,
	}, nil
}

func (c *powerSupplyClassCollector) Update(ch chan<- prometheus.Metric) error {
	powerSupplyClass, err := getPowerSupplyClassInfo(c.ignoredPattern)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoData
		}
		return fmt.Errorf("could not get power_supply class info: %w", err)
	}
	for _, powerSupply := range powerSupplyClass {

		for name, value := range map[string]*int64{
			"authentic":             powerSupply.Authentic,
			"calibrate":             powerSupply.Calibrate,
			"capacity":              powerSupply.Capacity,
			"capacity_alert_max":    powerSupply.CapacityAlertMax,
			"capacity_alert_min":    powerSupply.CapacityAlertMin,
			"cyclecount":            powerSupply.CycleCount,
			"online":                powerSupply.Online,
			"present":               powerSupply.Present,
			"time_to_empty_seconds": powerSupply.TimeToEmptyNow,
			"time_to_full_seconds":  powerSupply.TimeToFullNow,
		} {
			if value != nil {
				pushPowerSupplyMetric(ch, c.subsystem, name, float64(*value), powerSupply.Name, prometheus.GaugeValue)
			}
		}

		for name, value := range map[string]*int64{
			"current_boot":                powerSupply.CurrentBoot,
			"current_max":                 powerSupply.CurrentMax,
			"current_ampere":              powerSupply.CurrentNow,
			"energy_empty":                powerSupply.EnergyEmpty,
			"energy_empty_design":         powerSupply.EnergyEmptyDesign,
			"energy_full":                 powerSupply.EnergyFull,
			"energy_full_design":          powerSupply.EnergyFullDesign,
			"energy_watthour":             powerSupply.EnergyNow,
			"voltage_boot":                powerSupply.VoltageBoot,
			"voltage_max":                 powerSupply.VoltageMax,
			"voltage_max_design":          powerSupply.VoltageMaxDesign,
			"voltage_min":                 powerSupply.VoltageMin,
			"voltage_min_design":          powerSupply.VoltageMinDesign,
			"voltage_volt":                powerSupply.VoltageNow,
			"voltage_ocv":                 powerSupply.VoltageOCV,
			"charge_control_limit":        powerSupply.ChargeControlLimit,
			"charge_control_limit_max":    powerSupply.ChargeControlLimitMax,
			"charge_counter":              powerSupply.ChargeCounter,
			"charge_empty":                powerSupply.ChargeEmpty,
			"charge_empty_design":         powerSupply.ChargeEmptyDesign,
			"charge_full":                 powerSupply.ChargeFull,
			"charge_full_design":          powerSupply.ChargeFullDesign,
			"charge_ampere":               powerSupply.ChargeNow,
			"charge_term_current":         powerSupply.ChargeTermCurrent,
			"constant_charge_current":     powerSupply.ConstantChargeCurrent,
			"constant_charge_current_max": powerSupply.ConstantChargeCurrentMax,
			"constant_charge_voltage":     powerSupply.ConstantChargeVoltage,
			"constant_charge_voltage_max": powerSupply.ConstantChargeVoltageMax,
			"precharge_current":           powerSupply.PrechargeCurrent,
			"input_current_limit":         powerSupply.InputCurrentLimit,
			"power_watt":                  powerSupply.PowerNow,
		} {
			if value != nil {
				pushPowerSupplyMetric(ch, c.subsystem, name, float64(*value)/1e6, powerSupply.Name, prometheus.GaugeValue)
			}
		}

		for name, value := range map[string]*int64{
			"temp_celsius":             powerSupply.Temp,
			"temp_alert_max_celsius":   powerSupply.TempAlertMax,
			"temp_alert_min_celsius":   powerSupply.TempAlertMin,
			"temp_ambient_celsius":     powerSupply.TempAmbient,
			"temp_ambient_max_celsius": powerSupply.TempAmbientMax,
			"temp_ambient_min_celsius": powerSupply.TempAmbientMin,
			"temp_max_celsius":         powerSupply.TempMax,
			"temp_min_celsius":         powerSupply.TempMin,
		} {
			if value != nil {
				pushPowerSupplyMetric(ch, c.subsystem, name, float64(*value)/10.0, powerSupply.Name, prometheus.GaugeValue)
			}
		}

		var (
			keys   []string
			values []string
		)
		for name, value := range map[string]string{
			"power_supply":   powerSupply.Name,
			"capacity_level": powerSupply.CapacityLevel,
			"charge_type":    powerSupply.ChargeType,
			"health":         powerSupply.Health,
			"manufacturer":   powerSupply.Manufacturer,
			"model_name":     powerSupply.ModelName,
			"serial_number":  powerSupply.SerialNumber,
			"status":         powerSupply.Status,
			"technology":     powerSupply.Technology,
			"type":           powerSupply.Type,
			"usb_type":       powerSupply.UsbType,
			"scope":          powerSupply.Scope,
		} {
			if value != "" {
				keys = append(keys, name)
				values = append(values, strings.ToValidUTF8(value, "�"))
			}
		}

		fieldDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"info of /sys/class/power_supply/<power_supply>.",
			keys,
			nil,
		)
		ch <- prometheus.MustNewConstMetric(fieldDesc, prometheus.GaugeValue, 1.0, values...)

	}

	return nil
}

func pushPowerSupplyMetric(ch chan<- prometheus.Metric, subsystem string, name string, value float64, powerSupplyName string, valueType prometheus.ValueType) {
	fieldDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		fmt.Sprintf("%s value of /sys/class/power_supply/<power_supply>.", name),
		[]string{"power_supply"},
		nil,
	)

	ch <- prometheus.MustNewConstMetric(fieldDesc, valueType, value, powerSupplyName)
}

func getPowerSupplyClassInfo(ignore *regexp.Regexp) (sysfs.PowerSupplyClass, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, err
	}
	powerSupplyClass, err := fs.PowerSupplyClass()

	if err != nil {
		return powerSupplyClass, fmt.Errorf("error obtaining power_supply class info: %w", err)
	}

	for device := range powerSupplyClass {
		if ignore.MatchString(device) {
			delete(powerSupplyClass, device)
		}
	}

	return powerSupplyClass, nil
}
