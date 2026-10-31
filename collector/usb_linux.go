// Copyright 2026 The Prometheus Authors
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

//go:build !nousb

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"unicode"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

// usbDeviceLabels are the labels shared by all per-device metrics. devpath is
// the sysfs entry name (e.g. "1-3"), which embeds the bus number and is
// therefore unique on a host. serial is a join key to other collectors, e.g.
// node_disk_info; devnum is deliberately not exported as it changes whenever a
// device is re-enumerated.
var usbDeviceLabels = []string{"devpath", "serial"}

type usbCollector struct {
	fs sysfs.FS

	infoDesc *prometheus.Desc

	speedBytesDesc         *prometheus.Desc
	removableDesc          *prometheus.Desc
	maxCurrentDesc         *prometheus.Desc
	powerInfoDesc          *prometheus.Desc
	autosuspendDelayDesc   *prometheus.Desc
	urbnumDesc             *prometheus.Desc
	numInterfacesDesc      *prometheus.Desc
	numConfigurationsDesc  *prometheus.Desc
	interfaceInfoDesc      *prometheus.Desc
	controllerInfoDesc     *prometheus.Desc
	controllerSpeedDesc    *prometheus.Desc
	controllerMaxChildDesc *prometheus.Desc
}

func init() {
	registerCollector("usb", defaultDisabled, NewUSBCollector)
}

// NewUSBCollector returns a new Collector exposing USB device stats.
func NewUSBCollector(_ *slog.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	c := &usbCollector{
		fs: fs,
		infoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "info"),
			"Info about USB devices.",
			[]string{"devpath", "serial", "id_vendor", "id_product", "manufacturer", "product", "device_class", "device_subclass", "device_protocol", "bcd_device", "version"},
			nil,
		),
		speedBytesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "speed_bytes"),
			"USB device speed in bytes per second.",
			usbDeviceLabels,
			nil,
		),
		removableDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "removable"),
			"USB device removable state (1 = removable).",
			usbDeviceLabels,
			nil,
		),
		maxCurrentDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "max_current_milliamps"),
			"USB device maximum current in milliamps.",
			usbDeviceLabels,
			nil,
		),
		powerInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "power_info"),
			"USB device power information.",
			[]string{"devpath", "serial", "control", "runtime_status"},
			nil,
		),
		autosuspendDelayDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "autosuspend_delay_seconds"),
			"USB device autosuspend delay in seconds.",
			usbDeviceLabels,
			nil,
		),
		urbnumDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "urbnum_total"),
			"Total number of USB request blocks.",
			usbDeviceLabels,
			nil,
		),
		numInterfacesDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "num_interfaces"),
			"Number of USB interfaces.",
			usbDeviceLabels,
			nil,
		),
		numConfigurationsDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "num_configurations"),
			"Number of USB configurations.",
			usbDeviceLabels,
			nil,
		),
		interfaceInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "interface_info"),
			"Info about USB interfaces.",
			[]string{"devpath", "serial", "configuration", "interface", "binterface_class", "binterface_subclass", "binterface_protocol", "driver"},
			nil,
		),
		controllerInfoDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "controller_info"),
			"Info about USB controllers.",
			[]string{"bus", "version"},
			nil,
		),
		controllerSpeedDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "controller_speed_bytes"),
			"USB controller speed in bytes per second.",
			[]string{"bus"},
			nil,
		),
		controllerMaxChildDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "usb", "controller_maxchild"),
			"USB controller maximum number of child ports.",
			[]string{"bus"},
			nil,
		),
	}

	return c, nil
}

// sanitizeLabel makes s safe to use as a Prometheus label value: invalid UTF-8
// sequences are replaced, and control characters such as the newlines some
// devices put into their descriptors are removed.
func sanitizeLabel(s string) string {
	s = strings.ToValidUTF8(s, "\uFFFD")
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

func (c *usbCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.USBDevices()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoData
		}
		return fmt.Errorf("couldn't get USB devices: %w", err)
	}

	controllers, err := c.fs.USBControllers()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNoData
		}
		return fmt.Errorf("couldn't get USB controllers: %w", err)
	}

	for _, controller := range controllers {
		ch <- prometheus.MustNewConstMetric(c.controllerInfoDesc, prometheus.GaugeValue, 1.0,
			controller.BusNum,
			controller.BCDUSB,
		)
		if controller.Speed > 0 {
			speedBytes := float64(controller.Speed) * 1000000.0 / 8.0
			ch <- prometheus.MustNewConstMetric(c.controllerSpeedDesc, prometheus.GaugeValue, speedBytes,
				controller.BusNum,
			)
		}
		ch <- prometheus.MustNewConstMetric(c.controllerMaxChildDesc, prometheus.GaugeValue, float64(controller.MaxChild),
			controller.BusNum,
		)
	}

	for _, device := range devices {
		devpath := sanitizeLabel(device.Name)
		serial := sanitizeLabel(device.Serial)

		ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1.0,
			devpath,
			serial,
			device.IDVendor,
			device.IDProduct,
			sanitizeLabel(device.Manufacturer),
			sanitizeLabel(device.Product),
			device.DeviceClass,
			device.DeviceSubClass,
			device.DeviceProtocol,
			device.BCDDevice,
			sanitizeLabel(device.BCDUSB),
		)

		if device.Speed > 0 {
			speedBytes := float64(device.Speed) * 1000000.0 / 8.0
			ch <- prometheus.MustNewConstMetric(c.speedBytesDesc, prometheus.GaugeValue, speedBytes,
				devpath, serial)
		}

		var removable float64
		if device.Removable != nil && *device.Removable == sysfs.USBRemovableRemovable {
			removable = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.removableDesc, prometheus.GaugeValue, removable,
			devpath, serial)

		if device.BMaxPower > 0 {
			ch <- prometheus.MustNewConstMetric(c.maxCurrentDesc, prometheus.GaugeValue, float64(device.BMaxPower),
				devpath, serial)
		}

		control := sanitizeLabel(device.PowerControl)
		runtimeStatus := sanitizeLabel(device.RuntimeStatus)
		if control != "" || runtimeStatus != "" {
			ch <- prometheus.MustNewConstMetric(c.powerInfoDesc, prometheus.GaugeValue, 1.0,
				devpath, serial, control, runtimeStatus)
		}

		if device.AutosuspendDelayMillis != nil {
			delaySeconds := float64(*device.AutosuspendDelayMillis) / 1000.0
			ch <- prometheus.MustNewConstMetric(c.autosuspendDelayDesc, prometheus.GaugeValue, delaySeconds,
				devpath, serial)
		}

		if device.URBNum > 0 {
			ch <- prometheus.MustNewConstMetric(c.urbnumDesc, prometheus.CounterValue, float64(device.URBNum),
				devpath, serial)
		}

		if device.NumInterfaces > 0 {
			ch <- prometheus.MustNewConstMetric(c.numInterfacesDesc, prometheus.GaugeValue, float64(device.NumInterfaces),
				devpath, serial)
		}

		if device.NumConfigurations > 0 {
			ch <- prometheus.MustNewConstMetric(c.numConfigurationsDesc, prometheus.GaugeValue, float64(device.NumConfigurations),
				devpath, serial)
		}

		for _, iface := range device.Interfaces {
			ch <- prometheus.MustNewConstMetric(c.interfaceInfoDesc, prometheus.GaugeValue, 1.0,
				devpath,
				serial,
				iface.Configuration,
				iface.Interface,
				iface.Class,
				iface.SubClass,
				iface.Protocol,
				sanitizeLabel(iface.Driver),
			)
		}
	}

	return nil
}
