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

//go:build linux && !nonumatopology

package collector

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// testNumaTopologyCollector wraps the collector to satisfy prometheus.Collector.
type testNumaTopologyCollector struct {
	c Collector
}

func (tc testNumaTopologyCollector) Collect(ch chan<- prometheus.Metric) {
	tc.c.Update(ch)
}

func (tc testNumaTopologyCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(tc, ch)
}

// pinnedDomainXML is a <domstatus> file with VM "test-vm" pinned to both NUMA nodes.
// Cell 0 → host node 0: 2 vCPUs (0-1), 2097152 KiB = 2147483648 bytes
// Cell 1 → host node 1: 2 vCPUs (2-3), 2097152 KiB = 2147483648 bytes
const pinnedDomainXML = `<domstatus state='running' reason='booted' pid='12345'>
  <domain type='kvm' xmlns:nova='http://openstack.org/xmlns/libvirt/nova/1.0'>
    <name>instance-0000001a</name>
    <nova:instance>
      <nova:name>test-vm</nova:name>
    </nova:instance>
    <cpu>
      <numa>
        <cell id='0' cpus='0-1' memory='2097152' unit='KiB'/>
        <cell id='1' cpus='2-3' memory='2097152' unit='KiB'/>
      </numa>
    </cpu>
    <numatune>
      <memnode cellid='0' mode='strict' nodeset='0'/>
      <memnode cellid='1' mode='strict' nodeset='1'/>
    </numatune>
  </domain>
</domstatus>`

// unpinnedDomainXML has no <memnode> elements — must be skipped by the collector.
const unpinnedDomainXML = `<domstatus state='running' reason='booted' pid='99999'>
  <domain type='kvm'>
    <name>instance-0000002b</name>
    <cpu><topology sockets='1' cores='4' threads='2'/></cpu>
  </domain>
</domstatus>`

func TestNumaTopologyCollector(t *testing.T) {
	// Point sysfs at the existing test fixtures.
	*sysPath = "fixtures/sys"
	*numaTopologyVMMetrics = true

	// Create a temp dir with libvirt XML files.
	libvirtDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(libvirtDir, "instance-pinned.xml"), []byte(pinnedDomainXML), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libvirtDir, "instance-unpinned.xml"), []byte(unpinnedDomainXML), 0644); err != nil {
		t.Fatal(err)
	}
	*numaTopologyLibvirtXMLDir = libvirtDir

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewNumaTopologyCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(&testNumaTopologyCollector{c: c})

	// Metric values derived from sys.ttar fixtures and pinnedDomainXML above.
	// node0 MemTotal: 134182340 kB * 1024 = 137402716160 bytes
	// node1 MemTotal: 134217728 kB * 1024 = 137438953472 bytes
	// node2 MemTotal: 134217728 kB * 1024 = 137438953472 bytes (cpulist empty → 0 CPUs)
	// VM mem per node: 2097152 KiB * 1024 = 2147483648 bytes
	expected := `
		# HELP node_numatopology_cpu_capacity Number of physical CPUs on this NUMA node.
		# TYPE node_numatopology_cpu_capacity gauge
		node_numatopology_cpu_capacity{node="0"} 2
		node_numatopology_cpu_capacity{node="1"} 2
		node_numatopology_cpu_capacity{node="2"} 0
		# HELP node_numatopology_cpu_used vCPUs assigned to NUMA-pinned VMs on this node.
		# TYPE node_numatopology_cpu_used gauge
		node_numatopology_cpu_used{node="0"} 2
		node_numatopology_cpu_used{node="1"} 2
		node_numatopology_cpu_used{node="2"} 0
		# HELP node_numatopology_memory_capacity_bytes Total memory on this NUMA node in bytes.
		# TYPE node_numatopology_memory_capacity_bytes gauge
		node_numatopology_memory_capacity_bytes{node="0"} 1.3740271616e+11
		node_numatopology_memory_capacity_bytes{node="1"} 1.37438953472e+11
		node_numatopology_memory_capacity_bytes{node="2"} 1.37438953472e+11
		# HELP node_numatopology_memory_used_bytes Memory assigned to NUMA-pinned VMs on this node in bytes.
		# TYPE node_numatopology_memory_used_bytes gauge
		node_numatopology_memory_used_bytes{node="0"} 2.147483648e+09
		node_numatopology_memory_used_bytes{node="1"} 2.147483648e+09
		node_numatopology_memory_used_bytes{node="2"} 0
		# HELP node_numatopology_vm_cpu vCPUs assigned to a specific VM on this NUMA node. Cardinality scales with VM count per host.
		# TYPE node_numatopology_vm_cpu gauge
		node_numatopology_vm_cpu{node="0",vm="test-vm"} 2
		node_numatopology_vm_cpu{node="1",vm="test-vm"} 2
		# HELP node_numatopology_vm_memory_bytes Memory assigned to a specific VM on this NUMA node in bytes. Cardinality scales with VM count per host.
		# TYPE node_numatopology_vm_memory_bytes gauge
		node_numatopology_vm_memory_bytes{node="0",vm="test-vm"} 2.147483648e+09
		node_numatopology_vm_memory_bytes{node="1",vm="test-vm"} 2.147483648e+09
	`

	if err := testutil.GatherAndCompare(reg, strings.NewReader(expected)); err != nil {
		t.Fatal(err)
	}
}

func TestNumaTopologyCollectorVMMetricsDisabled(t *testing.T) {
	*sysPath = "fixtures/sys"

	libvirtDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(libvirtDir, "instance-pinned.xml"), []byte(pinnedDomainXML), 0644); err != nil {
		t.Fatal(err)
	}
	*numaTopologyLibvirtXMLDir = libvirtDir

	*numaTopologyVMMetrics = false
	t.Cleanup(func() { *numaTopologyVMMetrics = true })

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c, err := NewNumaTopologyCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(&testNumaTopologyCollector{c: c})

	// Per-node aggregates still emit; per-VM metrics must be absent.
	expected := `
		# HELP node_numatopology_cpu_capacity Number of physical CPUs on this NUMA node.
		# TYPE node_numatopology_cpu_capacity gauge
		node_numatopology_cpu_capacity{node="0"} 2
		node_numatopology_cpu_capacity{node="1"} 2
		node_numatopology_cpu_capacity{node="2"} 0
		# HELP node_numatopology_cpu_used vCPUs assigned to NUMA-pinned VMs on this node.
		# TYPE node_numatopology_cpu_used gauge
		node_numatopology_cpu_used{node="0"} 2
		node_numatopology_cpu_used{node="1"} 2
		node_numatopology_cpu_used{node="2"} 0
		# HELP node_numatopology_memory_capacity_bytes Total memory on this NUMA node in bytes.
		# TYPE node_numatopology_memory_capacity_bytes gauge
		node_numatopology_memory_capacity_bytes{node="0"} 1.3740271616e+11
		node_numatopology_memory_capacity_bytes{node="1"} 1.37438953472e+11
		node_numatopology_memory_capacity_bytes{node="2"} 1.37438953472e+11
		# HELP node_numatopology_memory_used_bytes Memory assigned to NUMA-pinned VMs on this node in bytes.
		# TYPE node_numatopology_memory_used_bytes gauge
		node_numatopology_memory_used_bytes{node="0"} 2.147483648e+09
		node_numatopology_memory_used_bytes{node="1"} 2.147483648e+09
		node_numatopology_memory_used_bytes{node="2"} 0
	`

	if err := testutil.GatherAndCompare(reg, strings.NewReader(expected),
		"node_numatopology_cpu_capacity",
		"node_numatopology_cpu_used",
		"node_numatopology_memory_capacity_bytes",
		"node_numatopology_memory_used_bytes",
		"node_numatopology_vm_cpu",
		"node_numatopology_vm_memory_bytes",
	); err != nil {
		t.Fatal(err)
	}
}
