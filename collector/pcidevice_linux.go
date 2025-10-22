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
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

const (
	pcideviceSubsystem = "pcidevice"
)

var (
	pcideviceLabelNames = []string{"segment", "bus", "device", "function"}

	pcideviceMaxLinkTSDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, pcideviceSubsystem, "max_link_transfers_per_second"),
		"Value of maximum link's transfers per second (T/s)",
		pcideviceLabelNames, nil,
	)
	pcideviceMaxLinkWidthDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, pcideviceSubsystem, "max_link_width"),
		"Value of maximum link's width (number of lanes)",
		pcideviceLabelNames, nil,
	)

	pcideviceCurrentLinkTSDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, pcideviceSubsystem, "current_link_transfers_per_second"),
		"Value of current link's transfers per second (T/s)",
		pcideviceLabelNames, nil,
	)
	pcideviceCurrentLinkWidthDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, pcideviceSubsystem, "current_link_width"),
		"Value of current link's width (number of lanes)",
		pcideviceLabelNames, nil,
	)
)

type pcideviceCollector struct {
	fs       sysfs.FS
	infoDesc typedDesc
	descs    []typedFactorDesc
	logger   *slog.Logger
}

func init() {
	registerCollector("pcidevice", defaultDisabled, NewPcideviceCollector)
}

// NewPcideviceCollector returns a new Collector exposing PCI devices stats.
func NewPcideviceCollector(logger *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	c := pcideviceCollector{
		fs:     fs,
		logger: logger,
		infoDesc: typedDesc{
			desc: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, pcideviceSubsystem, "info"),
				"Non-numeric data from /sys/bus/pci/devices/<location>, value is always 1.",
				append(pcideviceLabelNames,
					[]string{"parent_segment", "parent_bus", "parent_device", "parent_function",
						"class_id", "vendor_id", "subsystem_vendor_id", "subsystem_device_id", "revision"}...),
				nil,
			),
			valueType: prometheus.GaugeValue,
		},
		descs: []typedFactorDesc{
			{desc: pcideviceMaxLinkTSDesc, valueType: prometheus.GaugeValue},
			{desc: pcideviceMaxLinkWidthDesc, valueType: prometheus.GaugeValue},
			{desc: pcideviceCurrentLinkTSDesc, valueType: prometheus.GaugeValue},
			{desc: pcideviceCurrentLinkWidthDesc, valueType: prometheus.GaugeValue},
		},
	}

	return &c, nil
}

func (c *pcideviceCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.PciDevices()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.logger.Debug("PCI device not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining PCI device info: %w", err)
	}

	for _, device := range devices {
		// The device location is represented in separated format.
		values := device.Location.Strings()
		if device.ParentLocation != nil {
			values = append(values, device.ParentLocation.Strings()...)
		} else {
			values = append(values, []string{"*", "*", "*", "*"}...)
		}
		values = append(values, fmt.Sprintf("0x%06x", device.Class))
		values = append(values, fmt.Sprintf("0x%04x", device.Device))
		values = append(values, fmt.Sprintf("0x%04x", device.SubsystemVendor))
		values = append(values, fmt.Sprintf("0x%04x", device.SubsystemDevice))
		values = append(values, fmt.Sprintf("0x%02x", device.Revision))

		ch <- c.infoDesc.mustNewConstMetric(1.0, values...)

		// MaxLinkSpeed and CurrentLinkSpeed are represented in GT/s
		maxLinkSpeedTS := float64(int64(*device.MaxLinkSpeed * 1e9))
		currentLinkSpeedTS := float64(int64(*device.CurrentLinkSpeed * 1e9))

		for i, val := range []float64{
			maxLinkSpeedTS,
			float64(*device.MaxLinkWidth),
			currentLinkSpeedTS,
			float64(*device.CurrentLinkWidth),
		} {
			ch <- c.descs[i].mustNewConstMetric(val, device.Location.Strings()...)
		}
	}

	return nil
}
