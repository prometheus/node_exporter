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

// +build linux
// +build !noinfiniband

package collector

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
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
func NewInfiniBandCollector() (Collector, error) {
	var i infinibandCollector

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

	// Deprecated counters for some older versions of InfiniBand drivers.
	i.legacyCounters = map[string]infinibandMetric{
		"legacy_multicast_packets_received_total":    {"port_multicast_rcv_packets", "Number of multicast packets received"},
		"legacy_multicast_packets_transmitted_total": {"port_multicast_xmit_packets", "Number of multicast packets transmitted"},
		"legacy_data_received_bytes_total":           {"port_rcv_data_64", "Number of data octets received on all links"},
		"legacy_packets_received_total":              {"port_rcv_packets_64", "Number of data packets received on all links"},
		"legacy_unicast_packets_received_total":      {"port_unicast_rcv_packets", "Number of unicast packets received"},
		"legacy_unicast_packets_transmitted_total":   {"port_unicast_xmit_packets", "Number of unicast packets transmitted"},
		"legacy_data_transmitted_bytes_total":        {"port_xmit_data_64", "Number of data octets transmitted on all links"},
		"legacy_packets_transmitted_total":           {"port_xmit_packets_64", "Number of data packets received on all links"},
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
	devices, err := infinibandDevices(sysFilePath(infinibandPath))

	// If no devices are found or another error is raised while attempting to find devices,
	// InfiniBand is likely not installed and the collector should be skipped.
	switch err {
	case nil:
	case errInfinibandNoDevicesFound:
		return nil
	default:
		return err
	}

	for _, device := range devices {
		ports, err := infinibandPorts(sysFilePath(infinibandPath), device)

		// If no ports are found for the specified device, skip to the next device.
		switch err {
		case nil:
		case errInfinibandNoPortsFound:
			continue
		default:
			return err
		}

		for _, port := range ports {
			portFiles := sysFilePath(filepath.Join(infinibandPath, device, "ports", port))

			// Add metrics for the InfiniBand counters.
			for metricName, infinibandMetric := range c.counters {
				if _, err := os.Stat(filepath.Join(portFiles, "counters", infinibandMetric.File)); os.IsNotExist(err) {
					continue
				}
				metric, err := readMetric(filepath.Join(portFiles, "counters"), infinibandMetric.File)
				if err != nil {
					return err
				}

				ch <- prometheus.MustNewConstMetric(
					c.metricDescs[metricName],
					prometheus.CounterValue,
					float64(metric),
					device,
					port,
				)
			}

			// Add metrics for the legacy InfiniBand counters.
			for metricName, infinibandMetric := range c.legacyCounters {
				if _, err := os.Stat(filepath.Join(portFiles, "counters_ext", infinibandMetric.File)); os.IsNotExist(err) {
					continue
				}
				metric, err := readMetric(filepath.Join(portFiles, "counters_ext"), infinibandMetric.File)
				if err != nil {
					return err
				}

				ch <- prometheus.MustNewConstMetric(
					c.metricDescs[metricName],
					prometheus.CounterValue,
					float64(metric),
					device,
					port,
				)
			}
		}
	}

	return nil
}
