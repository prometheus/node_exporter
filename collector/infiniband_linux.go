// Copyright 2017 The Prometheus Authors
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

//go:build linux && !noinfiniband
// +build linux,!noinfiniband

package collector

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const infinibandPath = "class/infiniband"

var (
	errInfinibandNoDevicesFound = errors.New("no InfiniBand devices detected")
	errInfinibandNoPortsFound   = errors.New("no InfiniBand ports detected")
)

type infinibandCollector struct {
	metricDescs    map[string]*prometheus.Desc
	counters       map[string]infinibandMetric
	legacyCounters map[string]infinibandMetric
}

type infinibandMetric struct {
	File string
	Help string
}

func init() {
	registerCollector("infiniband", defaultEnabled, NewInfiniBandCollector)
}

// NewInfiniBandCollector returns a new Collector exposing InfiniBand stats.
func NewInfiniBandCollector(logger log.Logger) (Collector, error) {
	var i infinibandCollector
	var err error

	i.fs, err = sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}
	i.logger = logger

	// Filenames of all InfiniBand counter metrics including a detailed description.
	i.counters = map[string]infinibandMetric{
		"link_downed_total":                   {"link_downed", "Number of times the link failed to recover from an error state and went down"},
		"link_error_recovery_total":           {"link_error_recovery", "Number of times the link successfully recovered from an error state"},
		"multicast_packets_received_total":    {"multicast_rcv_packets", "Number of multicast packets received (including errors)"},
		"multicast_packets_transmitted_total": {"multicast_xmit_packets", "Number of multicast packets transmitted (including errors)"},
		"port_data_received_bytes_total":      {"port_rcv_data", "Number of data octets received on all links"},
		"port_data_transmitted_bytes_total":   {"port_xmit_data", "Number of data octets transmitted on all links"},
		"unicast_packets_received_total":      {"unicast_rcv_packets", "Number of unicast packets received (including errors)"},
		"unicast_packets_transmitted_total":   {"unicast_xmit_packets", "Number of unicast packets transmitted (including errors)"},
	}

	// Detailed description for all metrics.
	descriptions := map[string]string{
		"legacy_multicast_packets_received_total":    "Number of multicast packets received",
		"legacy_multicast_packets_transmitted_total": "Number of multicast packets transmitted",
		"legacy_data_received_bytes_total":           "Number of data octets received on all links",
		"legacy_packets_received_total":              "Number of data packets received on all links",
		"legacy_unicast_packets_received_total":      "Number of unicast packets received",
		"legacy_unicast_packets_transmitted_total":   "Number of unicast packets transmitted",
		"legacy_data_transmitted_bytes_total":        "Number of data octets transmitted on all links",
		"legacy_packets_transmitted_total":           "Number of data packets received on all links",
		"excessive_buffer_overrun_errors_total":      "Number of times that OverrunErrors consecutive flow control update periods occurred, each having at least one overrun error.",
		"link_downed_total":                          "Number of times the link failed to recover from an error state and went down",
		"link_error_recovery_total":                  "Number of times the link successfully recovered from an error state",
		"local_link_integrity_errors_total":          "Number of times that the count of local physical errors exceeded the threshold specified by LocalPhyErrors.",
		"multicast_packets_received_total":           "Number of multicast packets received (including errors)",
		"multicast_packets_transmitted_total":        "Number of multicast packets transmitted (including errors)",
		"physical_state_id":                          "Physical state of the InfiniBand port (0: no change, 1: sleep, 2: polling, 3: disable, 4: shift, 5: link up, 6: link error recover, 7: phytest)",
		"port_constraint_errors_received_total":      "Number of packets received on the switch physical port that are discarded",
		"port_constraint_errors_transmitted_total":   "Number of packets not transmitted from the switch physical port",
		"port_data_received_bytes_total":             "Number of data octets received on all links",
		"port_data_transmitted_bytes_total":          "Number of data octets transmitted on all links",
		"port_discards_received_total":               "Number of inbound packets discarded by the port because the port is down or congested",
		"port_discards_transmitted_total":            "Number of outbound packets discarded by the port because the port is down or congested",
		"port_errors_received_total":                 "Number of packets containing an error that were received on this port",
		"port_packets_received_total":                "Number of packets received on all VLs by this port (including errors)",
		"port_packets_transmitted_total":             "Number of packets transmitted on all VLs from this port (including errors)",
		"port_transmit_wait_total":                   "Number of ticks during which the port had data to transmit but no data was sent during the entire tick",
		"rate_bytes_per_second":                      "Maximum signal transfer rate",
		"state_id":                                   "State of the InfiniBand port (0: no change, 1: down, 2: init, 3: armed, 4: active, 5: act defer)",
		"unicast_packets_received_total":             "Number of unicast packets received (including errors)",
		"unicast_packets_transmitted_total":          "Number of unicast packets transmitted (including errors)",
		"port_receive_remote_physical_errors_total":  "Number of packets marked with the EBP (End of Bad Packet) delimiter received on the port.",
		"port_receive_switch_relay_errors_total":     "Number of packets that could not be forwarded by the switch.",
		"symbol_error_total":                         "Number of minor link errors detected on one or more physical lanes.",
		"vl15_dropped_total":                         "Number of incoming VL15 packets dropped due to resource limitations.",
	}

	subsystem := "infiniband"
	i.metricDescs = make(map[string]*prometheus.Desc)

	for metricName, infinibandMetric := range i.counters {
		i.metricDescs[metricName] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, metricName),
			infinibandMetric.Help,
			[]string{"device", "port"},
			nil,
		)
	}

	for metricName, infinibandMetric := range i.legacyCounters {
		i.metricDescs[metricName] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, metricName),
			infinibandMetric.Help,
			[]string{"device", "port"},
			nil,
		)
	}

	return &i, nil
}

// infinibandDevices retrieves a list of InfiniBand devices.
func infinibandDevices(infinibandPath string) ([]string, error) {
	devices, err := filepath.Glob(filepath.Join(infinibandPath, "/*"))
	if err != nil {
		return nil, err
	}

	if len(devices) < 1 {
		log.Debugf("Unable to detect InfiniBand devices")
		err = errInfinibandNoDevicesFound
		return nil, err
	}

	// Extract just the filenames which equate to the device names.
	for i, device := range devices {
		devices[i] = filepath.Base(device)
	}

	return devices, nil
}

// Retrieve a list of ports for the InfiniBand device.
func infinibandPorts(infinibandPath, device string) ([]string, error) {
	ports, err := filepath.Glob(filepath.Join(infinibandPath, device, "ports/*"))
	if err != nil {
		return nil, err
	}

	if len(ports) < 1 {
		log.Debugf("Unable to detect ports for %s", device)
		err = errInfinibandNoPortsFound
		return nil, err
	}

	// Extract just the filenames which equates to the port numbers.
	for i, port := range ports {
		ports[i] = filepath.Base(port)
	}

	return ports, nil
}

func readMetric(directory, metricFile string) (uint64, error) {
	metric, err := readUintFromFile(filepath.Join(directory, metricFile))
	if err != nil {
		// Ugly workaround for handling #966, when counters are
		// `N/A (not available)`.
		// This was already patched and submitted, see
		// https://www.spinics.net/lists/linux-rdma/msg68596.html
		// Remove this as soon as the fix lands in the enterprise distros.
		if strings.Contains(err.Error(), "N/A (no PMA)") {
			log.Debugf("%q value is N/A", metricFile)
			return 0, nil
		}
		log.Debugf("Error reading %q file", metricFile)
		return 0, err
	}

	// According to Mellanox, the following metrics "are divided by 4 unconditionally"
	// as they represent the amount of data being transmitted and received per lane.
	// Mellanox cards have 4 lanes per port, so all values must be multiplied by 4
	// to get the expected value.
	switch metricFile {
	case "port_rcv_data", "port_xmit_data", "port_rcv_data_64", "port_xmit_data_64":
		metric *= 4
	}

	return metric, nil
}

func (c *infinibandCollector) Update(ch chan<- prometheus.Metric) error {
	devices, err := c.fs.InfiniBandClass()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			level.Debug(c.logger).Log("msg", "infiniband statistics not found, skipping")
			return ErrNoData
		}
		return fmt.Errorf("error obtaining InfiniBand class info: %w", err)
	}

	for _, device := range devices {
		infoDesc := prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, "info"),
			"Non-numeric data from /sys/class/infiniband/<device>, value is always 1.",
			[]string{"device", "board_id", "firmware_version", "hca_type"},
			nil,
		)
		infoValue := 1.0
		ch <- prometheus.MustNewConstMetric(infoDesc, prometheus.GaugeValue, infoValue, device.Name, device.BoardID, device.FirmwareVersion, device.HCAType)

		for _, port := range device.Ports {
			portStr := strconv.FormatUint(uint64(port.Port), 10)

			c.pushMetric(ch, "state_id", uint64(port.StateID), port.Name, portStr, prometheus.GaugeValue)
			c.pushMetric(ch, "physical_state_id", uint64(port.PhysStateID), port.Name, portStr, prometheus.GaugeValue)
			c.pushMetric(ch, "rate_bytes_per_second", port.Rate, port.Name, portStr, prometheus.GaugeValue)

			c.pushCounter(ch, "legacy_multicast_packets_received_total", port.Counters.LegacyPortMulticastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_multicast_packets_transmitted_total", port.Counters.LegacyPortMulticastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_data_received_bytes_total", port.Counters.LegacyPortRcvData64, port.Name, portStr)
			c.pushCounter(ch, "legacy_packets_received_total", port.Counters.LegacyPortRcvPackets64, port.Name, portStr)
			c.pushCounter(ch, "legacy_unicast_packets_received_total", port.Counters.LegacyPortUnicastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_unicast_packets_transmitted_total", port.Counters.LegacyPortUnicastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "legacy_data_transmitted_bytes_total", port.Counters.LegacyPortXmitData64, port.Name, portStr)
			c.pushCounter(ch, "legacy_packets_transmitted_total", port.Counters.LegacyPortXmitPackets64, port.Name, portStr)
			c.pushCounter(ch, "excessive_buffer_overrun_errors_total", port.Counters.ExcessiveBufferOverrunErrors, port.Name, portStr)
			c.pushCounter(ch, "link_downed_total", port.Counters.LinkDowned, port.Name, portStr)
			c.pushCounter(ch, "link_error_recovery_total", port.Counters.LinkErrorRecovery, port.Name, portStr)
			c.pushCounter(ch, "local_link_integrity_errors_total", port.Counters.LocalLinkIntegrityErrors, port.Name, portStr)
			c.pushCounter(ch, "multicast_packets_received_total", port.Counters.MulticastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "multicast_packets_transmitted_total", port.Counters.MulticastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_constraint_errors_received_total", port.Counters.PortRcvConstraintErrors, port.Name, portStr)
			c.pushCounter(ch, "port_constraint_errors_transmitted_total", port.Counters.PortXmitConstraintErrors, port.Name, portStr)
			c.pushCounter(ch, "port_data_received_bytes_total", port.Counters.PortRcvData, port.Name, portStr)
			c.pushCounter(ch, "port_data_transmitted_bytes_total", port.Counters.PortXmitData, port.Name, portStr)
			c.pushCounter(ch, "port_discards_received_total", port.Counters.PortRcvDiscards, port.Name, portStr)
			c.pushCounter(ch, "port_discards_transmitted_total", port.Counters.PortXmitDiscards, port.Name, portStr)
			c.pushCounter(ch, "port_errors_received_total", port.Counters.PortRcvErrors, port.Name, portStr)
			c.pushCounter(ch, "port_packets_received_total", port.Counters.PortRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "port_packets_transmitted_total", port.Counters.PortXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_transmit_wait_total", port.Counters.PortXmitWait, port.Name, portStr)
			c.pushCounter(ch, "unicast_packets_received_total", port.Counters.UnicastRcvPackets, port.Name, portStr)
			c.pushCounter(ch, "unicast_packets_transmitted_total", port.Counters.UnicastXmitPackets, port.Name, portStr)
			c.pushCounter(ch, "port_receive_remote_physical_errors_total", port.Counters.PortRcvRemotePhysicalErrors, port.Name, portStr)
			c.pushCounter(ch, "port_receive_switch_relay_errors_total", port.Counters.PortRcvSwitchRelayErrors, port.Name, portStr)
			c.pushCounter(ch, "symbol_error_total", port.Counters.SymbolError, port.Name, portStr)
			c.pushCounter(ch, "vl15_dropped_total", port.Counters.VL15Dropped, port.Name, portStr)
		}
	}

	return nil
}