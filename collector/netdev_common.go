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

// +build !nonetdev
// +build linux freebsd openbsd dragonfly darwin

package collector

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
	"reflect"
)

var (
	netdevIgnoredDevices = kingpin.Flag("collector.netdev.ignored-devices", "Regexp of net devices to ignore for netdev collector.").Default("^$").String()
)

type netDevCollector struct {
	subsystem             string
	ignoredDevicesPattern *regexp.Regexp
	metricDescs           map[string]*prometheus.Desc
}

func init() {
	registerCollector("netdev", defaultEnabled, NewNetDevCollector)
}

// NewNetDevCollector returns a new Collector exposing network device stats.
func NewNetDevCollector() (Collector, error) {
	pattern := regexp.MustCompile(*netdevIgnoredDevices)
	return &netDevCollector{
		subsystem:             "network",
		ignoredDevicesPattern: pattern,
		metricDescs:           map[string]*prometheus.Desc{},
	}, nil
}

func (c *netDevCollector) Update(ch chan<- prometheus.Metric) error {
	netDev, err := getNetDevStats(c.ignoredDevicesPattern)
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %s", err)
	}
	for dev, devStats := range netDev {
		for key, value := range devStats {
			desc, ok := c.metricDescs[key]
			if !ok {
				desc = prometheus.NewDesc(
					prometheus.BuildFQName(namespace, c.subsystem, key+"_total"),
					fmt.Sprintf("Network device statistic %s.", key),
					[]string{"device"},
					nil,
				)
				c.metricDescs[key] = desc
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in netstats: %s", value, err)
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, v, dev)
		}
	}

	netClass, err := getNetClassInfo(c.ignoredDevicesPattern)
	if err != nil {
		return fmt.Errorf("could not get net class info: %s", err)
	}
	for _, ifaceInfo := range netClass {
		upDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "up"),
			"Valid operstate for interface.",
			[]string{"interface", "address", "broadcast", "duplex", "operstate", "ifalias"},
			nil,
		)
		upValue := 0.0
		if ifaceInfo.OperState == "up" {
			upValue = 1.0
		}
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, upValue, ifaceInfo.Name, ifaceInfo.Address, ifaceInfo.Broadcast, ifaceInfo.Duplex, ifaceInfo.OperState, ifaceInfo.IfAlias)

		fields := []string{
			"AddrAssignType",
			"Carrier",
			"CarrierChanges",
			"CarrierUpCount",
			"CarrierDownCount",
			"DevId",
			"Dormant",
			"Flags",
			"IfIndex",
			"IfLink",
			"LinkMode",
			"Mtu",
			"NameAssignType",
			"NetDevGroup",
			"Speed",
			"TxQueueLen",
			"Type",
		}
		interfaceElem := reflect.ValueOf(&ifaceInfo).Elem()
		interfaceType := reflect.TypeOf(ifaceInfo)

		for _, fieldName := range fields {
			fieldValue := interfaceElem.FieldByName(fieldName)
			fieldType, found := interfaceType.FieldByName(fieldName)
			if !found {
				continue
			}
			fieldDesc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, c.subsystem, fieldType.Tag.Get("fileName")),
				fmt.Sprintf("value of /sys/class/net/<iface>/%s.", fieldType.Tag.Get("fileName")),
				[]string{"interface"},
				nil,
			)

			ch <- prometheus.MustNewConstMetric(fieldDesc, prometheus.GaugeValue, float64(fieldValue.Int()), ifaceInfo.Name)
		}
	}

	return nil
}
