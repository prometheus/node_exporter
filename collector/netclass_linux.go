// Copyright 2018 The Prometheus Authors
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

// +build !nonetclass
// +build linux

package collector

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	netclassIgnoredDevices = kingpin.Flag("collector.netclass.ignored-devices", "Regexp of net devices to ignore for netclass collector.").Default("^$").String()
	netclassInvalidSpeed   = kingpin.Flag("collector.netclass.ignore-invalid-speed", "Ignore devices where the speed is invalid. This will be the default behavior in 2.x.").Bool()
)

type netClassCollector struct {
	fs                    sysfs.FS
	subsystem             string
	ignoredDevicesPattern *regexp.Regexp
	metricDescs           map[string]*prometheus.Desc
	logger                log.Logger
}

func init() {
	registerCollector("netclass", defaultEnabled, NewNetClassCollector)
}

// NewNetClassCollector returns a new Collector exposing network class stats.
func NewNetClassCollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	pattern := regexp.MustCompile(*netclassIgnoredDevices)
	return &netClassCollector{
		fs:                    fs,
		subsystem:             "network",
		ignoredDevicesPattern: pattern,
		metricDescs:           map[string]*prometheus.Desc{},
		logger:                logger,
	}, nil
}

func (c *netClassCollector) Update(ch chan<- prometheus.Metric) error {
	netClass, err := c.getNetClassInfo()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			level.Debug(c.logger).Log("msg", "Could not read netclass file", "err", err)
			return ErrNoData
		}
		return fmt.Errorf("could not get net class info: %w", err)
	}
	for _, ifaceInfo := range netClass {
		labels := getLabelsFromIfAlias(ifaceInfo.Name)
		labelKeys := append([]string{"device"}, labels.keys()...)
		labelValues := append([]string{ifaceInfo.Name}, labels.values()...)

		upDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "up"),
			"Value is 1 if operstate is 'up', 0 otherwise.",
			labelKeys,
			nil,
		)
		upValue := 0.0
		if ifaceInfo.OperState == "up" {
			upValue = 1.0
		}

		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, upValue, labelValues...)

		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data from /sys/class/net/<iface>, value is always 1.",
			append([]string{"device", "address", "broadcast", "duplex", "operstate", "ifalias"}, labels.keys()...),
			nil,
		)
		infoValue := 1.0

		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue, append([]string{ifaceInfo.Name, ifaceInfo.Address, ifaceInfo.Broadcast, ifaceInfo.Duplex, ifaceInfo.OperState, ifaceInfo.IfAlias}, labels.values()...)...)

		if ifaceInfo.AddrAssignType != nil {
			pushMetric(ch, c.subsystem, "address_assign_type", *ifaceInfo.AddrAssignType, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.Carrier != nil {
			pushMetric(ch, c.subsystem, "carrier", *ifaceInfo.Carrier, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.CarrierChanges != nil {
			pushMetric(ch, c.subsystem, "carrier_changes_total", *ifaceInfo.CarrierChanges, labelKeys, labelValues, prometheus.CounterValue)
		}

		if ifaceInfo.CarrierUpCount != nil {
			pushMetric(ch, c.subsystem, "carrier_up_changes_total", *ifaceInfo.CarrierUpCount, labelKeys, labelValues, prometheus.CounterValue)
		}

		if ifaceInfo.CarrierDownCount != nil {
			pushMetric(ch, c.subsystem, "carrier_down_changes_total", *ifaceInfo.CarrierDownCount, labelKeys, labelValues, prometheus.CounterValue)
		}

		if ifaceInfo.DevID != nil {
			pushMetric(ch, c.subsystem, "device_id", *ifaceInfo.DevID, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.Dormant != nil {
			pushMetric(ch, c.subsystem, "dormant", *ifaceInfo.Dormant, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.Flags != nil {
			pushMetric(ch, c.subsystem, "flags", *ifaceInfo.Flags, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.IfIndex != nil {
			pushMetric(ch, c.subsystem, "iface_id", *ifaceInfo.IfIndex, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.IfLink != nil {
			pushMetric(ch, c.subsystem, "iface_link", *ifaceInfo.IfLink, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.LinkMode != nil {
			pushMetric(ch, c.subsystem, "iface_link_mode", *ifaceInfo.LinkMode, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.MTU != nil {
			pushMetric(ch, c.subsystem, "mtu_bytes", *ifaceInfo.MTU, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.NameAssignType != nil {
			pushMetric(ch, c.subsystem, "name_assign_type", *ifaceInfo.NameAssignType, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.NetDevGroup != nil {
			pushMetric(ch, c.subsystem, "net_dev_group", *ifaceInfo.NetDevGroup, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.Speed != nil {
			// Some devices return -1 if the speed is unknown.
			if *ifaceInfo.Speed >= 0 || !*netclassInvalidSpeed {
				speedBytes := int64(*ifaceInfo.Speed * 1000 * 1000 / 8)
				pushMetric(ch, c.subsystem, "speed_bytes", speedBytes, labelKeys, labelValues, prometheus.GaugeValue)
			}
		}

		if ifaceInfo.TxQueueLen != nil {
			pushMetric(ch, c.subsystem, "transmit_queue_length", *ifaceInfo.TxQueueLen, labelKeys, labelValues, prometheus.GaugeValue)
		}

		if ifaceInfo.Type != nil {
			pushMetric(ch, c.subsystem, "protocol_type", *ifaceInfo.Type, labelKeys, labelValues, prometheus.GaugeValue)
		}
	}

	return nil
}

func pushMetric(ch chan<- prometheus.Metric, subsystem string, name string, value int64, labelKeys []string, labelValues []string, valueType prometheus.ValueType) {
	fieldDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		fmt.Sprintf("%s value of /sys/class/net/<iface>.", name),
		labelKeys,
		nil,
	)

	ch <- prometheus.MustNewConstMetric(fieldDesc, valueType, float64(value), labelValues...)
}

func (c *netClassCollector) getNetClassInfo() (sysfs.NetClass, error) {
	netClass := sysfs.NetClass{}
	netDevices, err := c.fs.NetClassDevices()
	if err != nil {
		return netClass, err
	}

	for _, device := range netDevices {
		if c.ignoredDevicesPattern.MatchString(device) {
			continue
		}
		interfaceClass, err := c.fs.NetClassByIface(device)
		if err != nil {
			return netClass, err
		}
		netClass[device] = *interfaceClass
	}

	return netClass, nil
}
