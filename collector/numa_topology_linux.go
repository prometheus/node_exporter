// Copyright 2024 The Prometheus Authors
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

//go:build !nonumatopology

package collector

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/node_exporter/collector/numatopology"
)

const numaTopologySubsystem = "numatopology"

var (
	numaTopologyLibvirtXMLDir = kingpin.Flag(
		"collector.numatopology.libvirt-xml-dir",
		"Directory containing libvirt status XML files.",
	).Default("/run/libvirt/qemu").String()

	numaTopologyVMMetrics = kingpin.Flag(
		"collector.numatopology.vm-metrics",
		"Enable per-VM NUMA assignment metrics (node_numatopology_vm_cpu, node_numatopology_vm_memory_bytes). Cardinality scales with VM count per host.",
	).Default("true").Bool()

	numaNodeDirRE = regexp.MustCompile(`node(\d+)$`)
)

type numaTopologyCollector struct {
	logger        *slog.Logger
	libvirtXMLDir string
	emitVMMetrics bool

	cpuCapacity *prometheus.Desc
	memCapacity *prometheus.Desc
	cpuUsed     *prometheus.Desc
	memUsed     *prometheus.Desc
	vmCPU       *prometheus.Desc
	vmMemBytes  *prometheus.Desc
}

func init() {
	registerCollector(numaTopologySubsystem, defaultDisabled, NewNumaTopologyCollector)
}

// NewNumaTopologyCollector returns a new Collector exposing NUMA topology metrics.
func NewNumaTopologyCollector(logger *slog.Logger) (Collector, error) {
	return &numaTopologyCollector{
		logger:        logger,
		libvirtXMLDir: *numaTopologyLibvirtXMLDir,
		emitVMMetrics: *numaTopologyVMMetrics,
		cpuCapacity: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, numaTopologySubsystem, "cpu_capacity"),
			"Number of physical CPUs on this NUMA node.",
			[]string{"node"}, nil,
		),
		memCapacity: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, numaTopologySubsystem, "memory_capacity_bytes"),
			"Total memory on this NUMA node in bytes.",
			[]string{"node"}, nil,
		),
		cpuUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, numaTopologySubsystem, "cpu_used"),
			"vCPUs assigned to NUMA-pinned VMs on this node.",
			[]string{"node"}, nil,
		),
		memUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, numaTopologySubsystem, "memory_used_bytes"),
			"Memory assigned to NUMA-pinned VMs on this node in bytes.",
			[]string{"node"}, nil,
		),
		vmCPU: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, numaTopologySubsystem, "vm_cpu"),
			"vCPUs assigned to a specific VM on this NUMA node. Cardinality scales with VM count per host.",
			[]string{"node", "vm"}, nil,
		),
		vmMemBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, numaTopologySubsystem, "vm_memory_bytes"),
			"Memory assigned to a specific VM on this NUMA node in bytes. Cardinality scales with VM count per host.",
			[]string{"node", "vm"}, nil,
		),
	}, nil
}

// Update implements Collector.
func (c *numaTopologyCollector) Update(ch chan<- prometheus.Metric) error {
	nodeDirs, err := filepath.Glob(sysFilePath("devices/system/node/node[0-9]*"))
	if err != nil {
		return fmt.Errorf("globbing NUMA node sysfs: %w", err)
	}
	if len(nodeDirs) == 0 {
		return fmt.Errorf("no NUMA nodes found under %s", sysFilePath("devices/system/node"))
	}

	type nodeCapacity struct {
		cpus int
		mem  int64
	}
	capacities := make(map[string]nodeCapacity, len(nodeDirs))

	for _, dir := range nodeDirs {
		m := numaNodeDirRE.FindStringSubmatch(dir)
		if m == nil {
			continue
		}
		nodeID := m[1]

		cpulistBytes, err := os.ReadFile(filepath.Join(dir, "cpulist"))
		if err != nil {
			c.logger.Warn("reading cpulist", "node", nodeID, "err", err)
			continue
		}

		meminfoBytes, err := os.ReadFile(filepath.Join(dir, "meminfo"))
		if err != nil {
			c.logger.Warn("reading meminfo", "node", nodeID, "err", err)
			continue
		}

		cpuCount := numatopology.CountCPUList(strings.TrimSpace(string(cpulistBytes)))
		memBytes, err := numatopology.ParseMeminfo(string(meminfoBytes))
		if err != nil {
			c.logger.Warn("parsing meminfo", "node", nodeID, "err", err)
			continue
		}

		capacities[nodeID] = nodeCapacity{cpus: cpuCount, mem: memBytes}
	}

	if len(capacities) == 0 {
		return fmt.Errorf("no NUMA nodes with readable sysfs data")
	}

	// vmUsage: node ID string → vm name → [vCPUs, memBytes]
	vmUsage := make(map[string]map[string][2]int64)

	xmlFiles, err := filepath.Glob(filepath.Join(c.libvirtXMLDir, "*.xml"))
	if err != nil {
		c.logger.Warn("globbing libvirt XML dir", "dir", c.libvirtXMLDir, "err", err)
	}
	for _, xmlFile := range xmlFiles {
		data, err := os.ReadFile(xmlFile)
		if err != nil {
			c.logger.Warn("reading libvirt XML", "file", xmlFile, "err", err)
			continue
		}
		res, err := numatopology.ParseVirshXML(string(data))
		if err != nil {
			c.logger.Warn("parsing libvirt XML", "file", xmlFile, "err", err)
			continue
		}
		if res == nil {
			continue // not NUMA-pinned
		}
		for hostNode, usage := range res.HostNUMAUsage {
			nodeStr := strconv.Itoa(hostNode)
			if vmUsage[nodeStr] == nil {
				vmUsage[nodeStr] = make(map[string][2]int64)
			}
			prev := vmUsage[nodeStr][res.VMName]
			vmUsage[nodeStr][res.VMName] = [2]int64{prev[0] + usage[0], prev[1] + usage[1]}
		}
	}

	for nodeID, cap := range capacities {
		ch <- prometheus.MustNewConstMetric(c.cpuCapacity, prometheus.GaugeValue, float64(cap.cpus), nodeID)
		ch <- prometheus.MustNewConstMetric(c.memCapacity, prometheus.GaugeValue, float64(cap.mem), nodeID)

		var totalCPU, totalMem int64
		for vmName, usage := range vmUsage[nodeID] {
			totalCPU += usage[0]
			totalMem += usage[1]
			if c.emitVMMetrics {
				ch <- prometheus.MustNewConstMetric(c.vmCPU, prometheus.GaugeValue, float64(usage[0]), nodeID, vmName)
				ch <- prometheus.MustNewConstMetric(c.vmMemBytes, prometheus.GaugeValue, float64(usage[1]), nodeID, vmName)
			}
		}
		ch <- prometheus.MustNewConstMetric(c.cpuUsed, prometheus.GaugeValue, float64(totalCPU), nodeID)
		ch <- prometheus.MustNewConstMetric(c.memUsed, prometheus.GaugeValue, float64(totalMem), nodeID)
	}

	return nil
}
