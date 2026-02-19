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

//go:build !nocpu
// +build !nocpu

package collector

// Exported metrics, from /sys/class/sas_phy/<name>/*:
//
// - invalid_dword_count
// - loss_of_dword_sync_count
// - negotiated_linkrate
// - phy_reset_problem_count
// - running_disparity_error_count
//
// Four of these are counters, one is a gauge.

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
)

type sasCollector struct {
	fs sysfs.FS

	sasInvalidDword          *prometheus.Desc
	sasLossOfDwordSync       *prometheus.Desc
	sasNegotiatedLinkrate    *prometheus.Desc
	sasPhyResetProblem       *prometheus.Desc
	sasRunningDisparityError *prometheus.Desc

	logger log.Logger
}

const sasCollectorSubsystem = "sas"

func init() {
	registerCollector("sas", defaultEnabled, NewSASCollector)
}

// NewSASCollector returns a new Collector exposing SAS storage statistics.
func NewSASCollector(logger log.Logger) (Collector, error) {
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}
	c := &sasCollector{
		fs: fs,
		sasInvalidDword: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sasCollectorSubsystem, "phy_invalid_dword"),
			"SAS PHY count of invalid dwords.",
			[]string{"phy", "port", "host", "expander", "block_device", "sas_address"}, nil,
		),
		sasLossOfDwordSync: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sasCollectorSubsystem, "phy_loss_of_dword_sync"),
			"SAS PHY count of lost dword sync.",
			[]string{"phy", "port", "host", "expander", "block_device", "sas_address"}, nil,
		),
		sasNegotiatedLinkrate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sasCollectorSubsystem, "phy_negotiated_linkrate"),
			"SAS PHY negotiated link rate in Gbps.",
			[]string{"phy", "port", "host", "expander", "block_device", "sas_address"}, nil,
		),
		sasPhyResetProblem: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sasCollectorSubsystem, "phy_reset_problem"),
			"SAS PHY count of reset problems.",
			[]string{"phy", "port", "host", "expander", "block_device", "sas_address"}, nil,
		),
		sasRunningDisparityError: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, sasCollectorSubsystem, "phy_running_disparity_error"),
			"SAS PHY count of running disparity errors.",
			[]string{"phy", "port", "host", "expander", "block_device", "sas_address"}, nil,
		),
		logger: logger,
	}

	return c, nil
}

// Update implements Collector and exposes cpu related metrics from /proc/stat and /sys/.../cpu/.
func (s *sasCollector) Update(ch chan<- prometheus.Metric) error {
	sasEndDevices, err := s.fs.SASEndDeviceClass()
	if err != nil {
		return err
	}

	sasExpanders, err := s.fs.SASExpanderClass()
	if err != nil {
		return err
	}

	sasHosts, err := s.fs.SASHostClass()
	if err != nil {
		return err
	}

	sasPhys, err := s.fs.SASPhyClass()
	if err != nil {
		return err
	}

	sasPorts, err := s.fs.SASPortClass()
	if err != nil {
		return err
	}

	for _, phy := range sasPhys {
		// One per SAS link in the system.
		phyName := phy.Name
		portName := phy.SASPort
		sasAddress := phy.SASAddress
		hostName := ""
		expanderName := ""
		blockDeviceName := ""

		// Check to see if this Phy is connected directly to a SAS host (~controller card)
		host := sasHosts.GetByPhy(phyName)
		if host != nil {
			hostName = host.Name
		}

		expander := sasExpanders.GetByPhy(phyName)
		if expander != nil {
			expanderName = expander.Name

			// If we didn't find a host name before, then maybe it's connected
			// host -> expander -> Phy.  So try checking the expander's Ports?
			// TODO: Do this recursively for nested expanders.
			if hostName == "" {
				port := sasPorts.GetByExpander(expanderName)
				if port != nil {
					host := sasHosts.GetByPort(port.Name)
					if host != nil {
						hostName = host.Name
					}
				}
			}

		}

		// This is the best mapping that I've found for phy->blockDevice.  Yes, it's horrible.
		port := sasPorts.GetByName(portName)
		if port != nil {
			for _, portEndDevice := range port.EndDevices {
				endDevice := sasEndDevices.GetByName(portEndDevice)
				if endDevice != nil {
					if len(endDevice.BlockDevices) > 0 {
						blockDeviceName = endDevice.BlockDevices[0]
					}
				}
			}
		}

		ch <- prometheus.MustNewConstMetric(s.sasInvalidDword,
			prometheus.CounterValue,
			float64(phy.InvalidDwordCount),
			phyName,
			portName,
			hostName,
			expanderName,
			blockDeviceName,
			sasAddress)
		ch <- prometheus.MustNewConstMetric(s.sasLossOfDwordSync,
			prometheus.CounterValue,
			float64(phy.LossOfDwordSyncCount),
			phyName,
			portName,
			hostName,
			expanderName,
			blockDeviceName,
			sasAddress)
		ch <- prometheus.MustNewConstMetric(s.sasNegotiatedLinkrate,
			prometheus.GaugeValue,
			float64(phy.NegotiatedLinkrate),
			phyName,
			portName,
			hostName,
			expanderName,
			blockDeviceName,
			sasAddress)
		ch <- prometheus.MustNewConstMetric(s.sasPhyResetProblem,
			prometheus.CounterValue,
			float64(phy.PhyResetProblemCount),
			phyName,
			portName,
			hostName,
			expanderName,
			blockDeviceName,
			sasAddress)
		ch <- prometheus.MustNewConstMetric(s.sasRunningDisparityError,
			prometheus.CounterValue,
			float64(phy.RunningDisparityErrorCount),
			phyName,
			portName,
			hostName,
			expanderName,
			blockDeviceName,
			sasAddress)
	}
	return nil
}
