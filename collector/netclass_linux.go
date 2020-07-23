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
	"fmt"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	netclassIgnoredDevices = kingpin.Flag("collector.netclass.ignored-devices", "Regexp of net devices to ignore for netclass collector.").Default("^$").String()
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
		return fmt.Errorf("could not get net class info: %w", err)
	}
	for _, ifaceInfo := range netClass {
		upDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "up"),
			"Value is 1 if operstate is 'up', 0 otherwise.",
			[]string{"device"},
			nil,
		)
		upValue := 0.0
		if ifaceInfo.OperState == "up" {
			upValue = 1.0
		}

		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, upValue, ifaceInfo.Name)

		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data from /sys/class/net/<iface>, value is always 1.",
			[]string{"device", "address", "broadcast", "duplex", "operstate", "ifalias"},
			nil,
		)
		infoValue := 1.0

		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue, ifaceInfo.Name, ifaceInfo.Address, ifaceInfo.Broadcast, ifaceInfo.Duplex, ifaceInfo.OperState, ifaceInfo.IfAlias)

		if ifaceInfo.AddrAssignType != nil {
			pushMetric(ch, c.subsystem, "address_assign_type", *ifaceInfo.AddrAssignType, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.Carrier != nil {
			pushMetric(ch, c.subsystem, "carrier", *ifaceInfo.Carrier, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.CarrierChanges != nil {
			pushMetric(ch, c.subsystem, "carrier_changes_total", *ifaceInfo.CarrierChanges, ifaceInfo.Name, prometheus.CounterValue)
		}

		if ifaceInfo.CarrierUpCount != nil {
			pushMetric(ch, c.subsystem, "carrier_up_changes_total", *ifaceInfo.CarrierUpCount, ifaceInfo.Name, prometheus.CounterValue)
		}

		if ifaceInfo.CarrierDownCount != nil {
			pushMetric(ch, c.subsystem, "carrier_down_changes_total", *ifaceInfo.CarrierDownCount, ifaceInfo.Name, prometheus.CounterValue)
		}

		if ifaceInfo.DevID != nil {
			pushMetric(ch, c.subsystem, "device_id", *ifaceInfo.DevID, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.Dormant != nil {
			pushMetric(ch, c.subsystem, "dormant", *ifaceInfo.Dormant, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.Flags != nil {
			pushMetric(ch, c.subsystem, "flags", *ifaceInfo.Flags, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.IfIndex != nil {
			pushMetric(ch, c.subsystem, "iface_id", *ifaceInfo.IfIndex, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.IfLink != nil {
			pushMetric(ch, c.subsystem, "iface_link", *ifaceInfo.IfLink, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.LinkMode != nil {
			pushMetric(ch, c.subsystem, "iface_link_mode", *ifaceInfo.LinkMode, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.MTU != nil {
			pushMetric(ch, c.subsystem, "mtu_bytes", *ifaceInfo.MTU, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.NameAssignType != nil {
			pushMetric(ch, c.subsystem, "name_assign_type", *ifaceInfo.NameAssignType, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.NetDevGroup != nil {
			pushMetric(ch, c.subsystem, "net_dev_group", *ifaceInfo.NetDevGroup, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.Speed != nil {
			speedBytes := int64(*ifaceInfo.Speed * 1000 * 1000 / 8)
			pushMetric(ch, c.subsystem, "speed_bytes", speedBytes, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.TxQueueLen != nil {
			pushMetric(ch, c.subsystem, "transmit_queue_length", *ifaceInfo.TxQueueLen, ifaceInfo.Name, prometheus.GaugeValue)
		}

		if ifaceInfo.Type != nil {
			pushMetric(ch, c.subsystem, "protocol_type", *ifaceInfo.Type, ifaceInfo.Name, prometheus.GaugeValue)
		}
	}

	return nil
}

func pushMetric(ch chan<- prometheus.Metric, subsystem string, name string, value int64, ifaceName string, valueType prometheus.ValueType) {
	fieldDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		fmt.Sprintf("%s value of /sys/class/net/<iface>.", name),
		[]string{"device"},
		nil,
	)

	ch <- prometheus.MustNewConstMetric(fieldDesc, valueType, float64(value), ifaceName)
}

func (c *netClassCollector) getNetClassInfo() (sysfs.NetClass, error) {
	netClass, err := c.fs.NetClass()

	if err != nil {
		return netClass, fmt.Errorf("error obtaining net class info: %w", err)
	}

	for device := range netClass {
		if c.ignoredDevicesPattern.MatchString(device) {
			delete(netClass, device)
		}
	}

	return netClass, nil
}
