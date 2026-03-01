// Copyright The Prometheus Authors
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

//go:build !nonetvf

package collector

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jsimonetti/rtnetlink/v2"
	"github.com/prometheus/client_golang/prometheus"
)

const netvfSubsystem = "net_vf"

var (
	netvfDeviceInclude = kingpin.Flag("collector.netvf.device-include", "Regexp of PF devices to include (mutually exclusive to device-exclude).").String()
	netvfDeviceExclude = kingpin.Flag("collector.netvf.device-exclude", "Regexp of PF devices to exclude (mutually exclusive to device-include).").String()
)

func init() {
	registerCollector("netvf", defaultDisabled, NewNetVFCollector)
}

type netvfCollector struct {
	logger       *slog.Logger
	deviceFilter deviceFilter

	info            *prometheus.Desc
	receivePackets  *prometheus.Desc
	transmitPackets *prometheus.Desc
	receiveBytes    *prometheus.Desc
	transmitBytes   *prometheus.Desc
	broadcast       *prometheus.Desc
	multicast       *prometheus.Desc
	receiveDropped  *prometheus.Desc
	transmitDropped *prometheus.Desc
}

func NewNetVFCollector(logger *slog.Logger) (Collector, error) {
	if *netvfDeviceExclude != "" && *netvfDeviceInclude != "" {
		return nil, errors.New("device-exclude & device-include are mutually exclusive")
	}

	if *netvfDeviceExclude != "" {
		logger.Info("Parsed flag --collector.netvf.device-exclude", "flag", *netvfDeviceExclude)
	}

	if *netvfDeviceInclude != "" {
		logger.Info("Parsed flag --collector.netvf.device-include", "flag", *netvfDeviceInclude)
	}

	return &netvfCollector{
		logger:       logger,
		deviceFilter: newDeviceFilter(*netvfDeviceExclude, *netvfDeviceInclude),
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "info"),
			"Virtual Function configuration information.",
			[]string{"device", "vf", "mac", "vlan", "link_state", "spoof_check", "trust", "pci_address"}, nil,
		),
		receivePackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "receive_packets_total"),
			"Number of received packets by the VF.",
			[]string{"device", "vf", "pci_address"}, nil,
		),
		transmitPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "transmit_packets_total"),
			"Number of transmitted packets by the VF.",
			[]string{"device", "vf", "pci_address"}, nil,
		),
		receiveBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "receive_bytes_total"),
			"Number of received bytes by the VF.",
			[]string{"device", "vf", "pci_address"}, nil,
		),
		transmitBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "transmit_bytes_total"),
			"Number of transmitted bytes by the VF.",
			[]string{"device", "vf", "pci_address"}, nil,
		),
		broadcast: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "broadcast_packets_total"),
			"Number of broadcast packets received by the VF.",
			[]string{"device", "vf", "pci_address"}, nil,
		),
		multicast: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "multicast_packets_total"),
			"Number of multicast packets received by the VF.",
			[]string{"device", "vf", "pci_address"}, nil,
		),
		receiveDropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "receive_dropped_total"),
			"Number of dropped received packets by the VF.",
			[]string{"device", "vf", "pci_address"}, nil,
		),
		transmitDropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, netvfSubsystem, "transmit_dropped_total"),
			"Number of dropped transmitted packets by the VF.",
			[]string{"device", "vf", "pci_address"}, nil,
		),
	}, nil
}

func (c *netvfCollector) Update(ch chan<- prometheus.Metric) error {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return fmt.Errorf("failed to connect to rtnetlink: %w", err)
	}
	defer conn.Close()

	links, err := conn.Link.ListWithVFInfo()
	if err != nil {
		return fmt.Errorf("failed to list interfaces with VF info: %w", err)
	}

	vfCount := 0
	for _, link := range links {
		if link.Attributes == nil {
			continue
		}

		// Skip interfaces without VFs
		if link.Attributes.NumVF == nil || *link.Attributes.NumVF == 0 {
			continue
		}

		device := link.Attributes.Name

		// Apply device filter
		if c.deviceFilter.ignored(device) {
			c.logger.Debug("Ignoring device", "device", device)
			continue
		}

		for _, vf := range link.Attributes.VFInfoList {
			vfID := fmt.Sprintf("%d", vf.ID)

			// Emit info metric with VF configuration
			mac := ""
			if vf.MAC != nil {
				mac = vf.MAC.String()
			}
			vlan := fmt.Sprintf("%d", vf.Vlan)
			linkState := vfLinkStateString(vf.LinkState)
			spoofCheck := fmt.Sprintf("%t", vf.SpoofCheck)
			trust := fmt.Sprintf("%t", vf.Trust)
			pciAddress := resolveVFPCIAddress(sysFilePath("class"), device, vf.ID)

			ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1, device, vfID, mac, vlan, linkState, spoofCheck, trust, pciAddress)

			// Emit stats metrics if available
			if vf.Stats == nil {
				c.logger.Debug("VF has no stats", "device", device, "vf", vf.ID)
				vfCount++
				continue
			}

			stats := vf.Stats

			ch <- prometheus.MustNewConstMetric(c.receivePackets, prometheus.CounterValue, float64(stats.RxPackets), device, vfID, pciAddress)
			ch <- prometheus.MustNewConstMetric(c.transmitPackets, prometheus.CounterValue, float64(stats.TxPackets), device, vfID, pciAddress)
			ch <- prometheus.MustNewConstMetric(c.receiveBytes, prometheus.CounterValue, float64(stats.RxBytes), device, vfID, pciAddress)
			ch <- prometheus.MustNewConstMetric(c.transmitBytes, prometheus.CounterValue, float64(stats.TxBytes), device, vfID, pciAddress)
			ch <- prometheus.MustNewConstMetric(c.broadcast, prometheus.CounterValue, float64(stats.Broadcast), device, vfID, pciAddress)
			ch <- prometheus.MustNewConstMetric(c.multicast, prometheus.CounterValue, float64(stats.Multicast), device, vfID, pciAddress)
			ch <- prometheus.MustNewConstMetric(c.receiveDropped, prometheus.CounterValue, float64(stats.RxDropped), device, vfID, pciAddress)
			ch <- prometheus.MustNewConstMetric(c.transmitDropped, prometheus.CounterValue, float64(stats.TxDropped), device, vfID, pciAddress)

			vfCount++
		}
	}

	if vfCount == 0 {
		return ErrNoData
	}

	return nil
}

func vfLinkStateString(state rtnetlink.VFLinkState) string {
	switch state {
	case rtnetlink.VFLinkStateAuto:
		return "auto"
	case rtnetlink.VFLinkStateEnable:
		return "enable"
	case rtnetlink.VFLinkStateDisable:
		return "disable"
	default:
		return "unknown"
	}
}

// resolveVFPCIAddress resolves the PCI BDF address of a VF by reading the
// sysfs virtfn symlink: <sysClassPath>/net/<pfDevice>/device/virtfn<vfID>.
// Returns empty string if the symlink doesn't exist or can't be resolved.
func resolveVFPCIAddress(sysClassPath, pfDevice string, vfID uint32) string {
	virtfnPath := filepath.Join(sysClassPath, "net", pfDevice, "device", fmt.Sprintf("virtfn%d", vfID))
	resolved, err := os.Readlink(virtfnPath)
	if err != nil {
		return ""
	}
	return filepath.Base(resolved)
}

// vfMetrics holds parsed VF metrics for a single VF
type vfMetrics struct {
	Device     string
	VFID       uint32
	MAC        string
	Vlan       uint32
	LinkState  string
	SpoofCheck bool
	Trust      bool
	PCIAddress string
	Stats      *rtnetlink.VFStats
}

// parseVFInfo extracts VF information from link messages for testing.
// sysClassPath is the path to the sysfs class directory used to resolve VF PCI addresses.
func parseVFInfo(links []rtnetlink.LinkMessage, filter *deviceFilter, logger *slog.Logger, sysClassPath string) []vfMetrics {
	var result []vfMetrics

	for _, link := range links {
		if link.Attributes == nil {
			continue
		}

		// Skip interfaces without VFs
		if link.Attributes.NumVF == nil || *link.Attributes.NumVF == 0 {
			continue
		}

		device := link.Attributes.Name

		// Apply device filter
		if filter.ignored(device) {
			logger.Debug("Ignoring device", "device", device)
			continue
		}

		for _, vf := range link.Attributes.VFInfoList {
			mac := ""
			if vf.MAC != nil {
				mac = vf.MAC.String()
			}

			result = append(result, vfMetrics{
				Device:     device,
				VFID:       vf.ID,
				MAC:        mac,
				Vlan:       vf.Vlan,
				LinkState:  vfLinkStateString(vf.LinkState),
				SpoofCheck: vf.SpoofCheck,
				Trust:      vf.Trust,
				PCIAddress: resolveVFPCIAddress(sysClassPath, device, vf.ID),
				Stats:      vf.Stats,
			})
		}
	}

	return result
}
