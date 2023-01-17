// Copyright 2022 The Prometheus Authors
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

//go:build !nonetclass && linux
// +build !nonetclass,linux

package collector

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jsimonetti/rtnetlink"
	"github.com/mdlayher/ethtool"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	netclassRTNLIgnoredDevices = kingpin.Flag("collector.netclass_rtnl.ignored-devices", "Regexp of net devices to ignore for netclass_rtnl collector.").Default("^$").String()
	netclassRTNLWithStats      = kingpin.Flag("collector.netclass_rtnl.with-stats", "Expose the statistics for each network device, replacing netdev collector.").Bool()
	operstateStr               = []string{
		"unknown", "notpresent", "down", "lowerlayerdown", "testing",
		"dormant", "up",
	}
)

type netClassRTNLCollector struct {
	subsystem             string
	ignoredDevicesPattern *regexp.Regexp
	metricDescs           map[string]*prometheus.Desc
	logger                log.Logger
}

func init() {
	registerCollector("netclass_rtnl", defaultDisabled, NewNetClassRTNLCollector)
}

// NewNetClassCollector returns a new Collector exposing network class stats.
func NewNetClassRTNLCollector(logger log.Logger) (Collector, error) {
	pattern := regexp.MustCompile(*netclassRTNLIgnoredDevices)
	return &netClassRTNLCollector{
		subsystem:             "network",
		ignoredDevicesPattern: pattern,
		metricDescs:           map[string]*prometheus.Desc{},
		logger:                logger,
	}, nil
}

func (c *netClassRTNLCollector) Update(ch chan<- prometheus.Metric) error {
	linkModes := make(map[string]*ethtool.LinkMode)
	lms, err := c.getLinkModes()
	if err != nil {
		if !errors.Is(errors.Unwrap(err), fs.ErrNotExist) {
			return fmt.Errorf("could not get link modes: %w", err)
		}
		level.Info(c.logger).Log("msg", "ETHTOOL netlink interface unavailable, duplex and linkspeed are not scraped.")
	} else {
		for _, lm := range lms {
			if c.ignoredDevicesPattern.MatchString(lm.Interface.Name) {
				continue
			}
			if lm.SpeedMegabits >= 0 {
				speedBytes := uint64(lm.SpeedMegabits * 1000 * 1000 / 8)
				pushMetric(ch, c.getFieldDesc("speed_bytes"), "speed_bytes", speedBytes, prometheus.GaugeValue, lm.Interface.Name)
			}
			linkModes[lm.Interface.Name] = lm
		}
	}

	lMsgs, err := c.getNetClassInfoRTNL()
	if err != nil {
		return fmt.Errorf("could not get net class info: %w", err)
	}
	for _, msg := range lMsgs {
		if c.ignoredDevicesPattern.MatchString(msg.Attributes.Name) {
			continue
		}
		upDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "up"),
			"Value is 1 if operstate is 'up', 0 otherwise.",
			[]string{"device"},
			nil,
		)
		upValue := 0.0
		if msg.Attributes.OperationalState == rtnetlink.OperStateUp {
			upValue = 1.0
		}
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, upValue, msg.Attributes.Name)

		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data of <iface>, value is always 1.",
			[]string{"device", "address", "broadcast", "duplex", "operstate", "ifalias"},
			nil,
		)
		infoValue := 1.0

		var ifalias = ""
		if msg.Attributes.Alias != nil {
			ifalias = *msg.Attributes.Alias
		}

		duplex := ""
		lm, lmExists := linkModes[msg.Attributes.Name]
		if lmExists {
			duplex = lm.Duplex.String()
		}

		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue, msg.Attributes.Name, msg.Attributes.Address.String(), msg.Attributes.Broadcast.String(), duplex, operstateStr[int(msg.Attributes.OperationalState)], ifalias)

		pushMetric(ch, c.getFieldDesc("carrier"), "carrier", msg.Attributes.Carrier, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("carrier_changes_total"), "carrier_changes_total", msg.Attributes.CarrierChanges, prometheus.CounterValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("carrier_up_changes_total"), "carrier_up_changes_total", msg.Attributes.CarrierUpCount, prometheus.CounterValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("carrier_down_changes_total"), "carrier_down_changes_total", msg.Attributes.CarrierDownCount, prometheus.CounterValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("flags"), "flags", msg.Flags, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("iface_id"), "iface_id", msg.Index, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("iface_link_mode"), "iface_link_mode", msg.Attributes.LinkMode, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("dormant"), "dormant", msg.Attributes.LinkMode, prometheus.GaugeValue, msg.Attributes.Name)

		// kernel logic: IFLA_LINK attribute will be ignore when ifindex is the same as iflink
		// (dev->ifindex != dev_get_iflink(dev) && nla_put_u32(skb, IFLA_LINK, dev_get_iflink(dev)))
		// As interface ID is never 0, we assume msg.Attributes.Type 0  means iflink is omitted in RTM_GETLINK response.
		if msg.Attributes.Type > 0 {
			pushMetric(ch, c.getFieldDesc("iface_link"), "iface_link", msg.Attributes.Type, prometheus.GaugeValue, msg.Attributes.Name)
		} else {
			pushMetric(ch, c.getFieldDesc("iface_link"), "iface_link", msg.Index, prometheus.GaugeValue, msg.Attributes.Name)
		}

		pushMetric(ch, c.getFieldDesc("mtu_bytes"), "mtu_bytes", msg.Attributes.MTU, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("net_dev_group"), "net_dev_group", msg.Attributes.NetDevGroup, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_queue_length"), "transmit_queue_length", msg.Attributes.TxQueueLen, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("protocol_type"), "protocol_type", msg.Type, prometheus.GaugeValue, msg.Attributes.Name)

		// skip statistics if argument collector.netclass_rtnl.with-stats is false or statistics are unavailable.
		if netclassRTNLWithStats == nil || !*netclassRTNLWithStats || msg.Attributes.Stats64 == nil {
			continue
		}

		pushMetric(ch, c.getFieldDesc("receive_packets_total"), "receive_packets_total", msg.Attributes.Stats64.RXPackets, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_packets_total"), "transmit_packets_total", msg.Attributes.Stats64.TXPackets, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_bytes_total"), "receive_bytes_total", msg.Attributes.Stats64.RXBytes, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_bytes_total"), "transmit_bytes_total", msg.Attributes.Stats64.TXBytes, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_errors_total"), "receive_errors_total", msg.Attributes.Stats64.RXErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_errors_total"), "transmit_errors_total", msg.Attributes.Stats64.TXErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_dropped_total"), "receive_dropped_total", msg.Attributes.Stats64.RXDropped, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_dropped_total"), "transmit_dropped_total", msg.Attributes.Stats64.TXDropped, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("multicast_total"), "multicast_total", msg.Attributes.Stats64.Multicast, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("collisions_total"), "collisions_total", msg.Attributes.Stats64.Collisions, prometheus.GaugeValue, msg.Attributes.Name)

		// detailed rx_errors
		pushMetric(ch, c.getFieldDesc("receive_length_errors_total"), "receive_length_errors_total", msg.Attributes.Stats64.RXLengthErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_over_errors_total"), "receive_over_errors_total", msg.Attributes.Stats64.RXOverErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_crc_errors_total"), "receive_crc_errors_total", msg.Attributes.Stats64.RXCRCErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_frame_errors_total"), "receive_frame_errors_total", msg.Attributes.Stats64.RXFrameErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_fifo_errors_total"), "receive_fifo_errors_total", msg.Attributes.Stats64.RXFIFOErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_missed_errors_total"), "receive_missed_errors_total", msg.Attributes.Stats64.RXMissedErrors, prometheus.GaugeValue, msg.Attributes.Name)

		// detailed tx_errors
		pushMetric(ch, c.getFieldDesc("transmit_aborted_errors_total"), "transmit_aborted_errors_total", msg.Attributes.Stats64.TXAbortedErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_carrier_errors_total"), "transmit_carrier_errors_total", msg.Attributes.Stats64.TXCarrierErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_fifo_errors_total"), "transmit_fifo_errors_total", msg.Attributes.Stats64.TXFIFOErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_heartbeat_errors_total"), "transmit_heartbeat_errors_total", msg.Attributes.Stats64.TXHeartbeatErrors, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_window_errors_total"), "transmit_window_errors_total", msg.Attributes.Stats64.TXWindowErrors, prometheus.GaugeValue, msg.Attributes.Name)

		// for cslip etc
		pushMetric(ch, c.getFieldDesc("receive_compressed_total"), "receive_compressed_total", msg.Attributes.Stats64.RXCompressed, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("transmit_compressed_total"), "transmit_compressed_total", msg.Attributes.Stats64.TXCompressed, prometheus.GaugeValue, msg.Attributes.Name)
		pushMetric(ch, c.getFieldDesc("receive_nohandler_total"), "receive_nohandler_total", msg.Attributes.Stats64.RXNoHandler, prometheus.GaugeValue, msg.Attributes.Name)

	}

	return nil
}

func (c *netClassRTNLCollector) getFieldDesc(name string) *prometheus.Desc {
	fieldDesc, exists := c.metricDescs[name]

	if !exists {
		fieldDesc = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, name),
			fmt.Sprintf("Network device property %s.", name),
			[]string{"device"},
			nil,
		)
		c.metricDescs[name] = fieldDesc
	}

	return fieldDesc
}

func (c *netClassRTNLCollector) getNetClassInfoRTNL() ([]rtnetlink.LinkMessage, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	lMsgs, err := conn.Link.List()

	return lMsgs, err

}

func (c *netClassRTNLCollector) getLinkModes() ([]*ethtool.LinkMode, error) {
	conn, err := ethtool.New()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	lms, err := conn.LinkModes()

	return lms, err
}
