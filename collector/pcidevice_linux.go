// Copyright 2017-2019 The Prometheus Authors
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

//go:build !nopcidevice
// +build !nopcidevice

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type pciDeviceCollector struct {
	fs          sysfs.FS
	metricDescs map[string]*prometheus.Desc
	logger      *slog.Logger
	subsystem   string
}

func init() {
	registerCollector("pcidevice", defaultDisabled, NewPciDeviceCollector)
}

// NewPciDeviceCollector returns a new Collector exposing PCI devices stats.
func NewPciDeviceCollector(logger *slog.Logger) (Collector, error) {
	var i pciDeviceCollector
	var err error

	i.fs, err = sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	i.logger = logger

	// Detailed description for all metrics.
	descriptions := map[string]string{
		"max_link_transfers_per_second":     "Value of maximum link's transfers per second (T/s)",
		"max_link_width":                    "Value of maximum link's width (number of lanes)",
		"current_link_transfers_per_second": "Value of current link's transfers per second (T/s)",
		"current_link_width":                "Value of current link's width (number of lanes)",
	}

	i.metricDescs = make(map[string]*prometheus.Desc)
	i.subsystem = "pcidevice"

	for metricName, description := range descriptions {
		i.metricDescs[metricName] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, i.subsystem, metricName),
			description,
			[]string{"segment", "bus", "device", "function"},
			nil,
		)
	}

	return &i, nil
}

func (c *pciDeviceCollector) pushMetric(ch chan<- prometheus.Metric, name string, value *float64, location sysfs.PciDeviceLocation, valueType prometheus.ValueType) {
	if value != nil {
		ch <- prometheus.MustNewConstMetric(c.metricDescs[name], valueType, *value, location.Strings()...)
	}
}

func (c *pciDeviceCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.PciDevices()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("PCI device not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining PCI device info: %w", err)
	}

	for _, device := range devices {
		// The format follows the definition in drivers/pci/pci-sysfs.c
		infos := [][]string{
			{"class_id", fmt.Sprintf("0x%06x", device.Class)},
			{"vendor_id", fmt.Sprintf("0x%04x", device.Device)},
			{"subsystem_vendor_id", fmt.Sprintf("0x%04x", device.SubsystemVendor)},
			{"subsystem_device_id", fmt.Sprintf("0x%04x", device.SubsystemDevice)},
			{"revision", fmt.Sprintf("0x%02x", device.Revision)},
		}

		labels := []string{}
		values := []string{}
		for i := range infos {
			labels = append(labels, infos[i][0])
			values = append(values, infos[i][1])
		}

		// The device location is represented in separated format.
		labels = append(labels, []string{"segment", "bus", "device", "function"}...)
		values = append(values, device.Location.Strings()...)

		labels = append(labels, []string{"parent_segment", "parent_bus", "parent_device", "parent_function"}...)
		if device.ParentLocation != nil {
			values = append(values, device.ParentLocation.Strings()...)
		} else {
			// TODO: is this ok?
			values = append(values, []string{"*", "*", "*", "*"}...)
		}

		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data from /sys/bus/pci/devices/<location>, value is always 1.",
			labels,
			nil,
		)
		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, 1.0, values...)

		// MaxLinkSpeed and CurrentLinkSpeed are represnted in GT/s
		maxLinkSpeedTS := float64(int64(*device.MaxLinkSpeed * math.Pow10(9)))
		currentLinkSpeedTS := float64(int64(*device.CurrentLinkSpeed * math.Pow10(9)))

		c.pushMetric(ch, "max_link_transfers_per_second", &maxLinkSpeedTS, device.Location, prometheus.GaugeValue)
		c.pushMetric(ch, "max_link_width", device.MaxLinkWidth, device.Location, prometheus.GaugeValue)
		c.pushMetric(ch, "current_link_transfers_per_second", &currentLinkSpeedTS, device.Location, prometheus.GaugeValue)
		c.pushMetric(ch, "current_link_width", device.CurrentLinkWidth, device.Location, prometheus.GaugeValue)
	}

	return nil
}
