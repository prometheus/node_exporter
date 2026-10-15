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

package numatopology

import (
	"strings"
	"testing"
)

// virshXMLNUMAPinned is <domain> root format (from "virsh dumpxml").
const virshXMLNUMAPinned = `<domain type='kvm' xmlns:nova='http://openstack.org/xmlns/libvirt/nova/1.0'>
  <name>instance-0000001a</name>
  <nova:instance>
    <nova:name>my-vm-01</nova:name>
  </nova:instance>
  <cpu>
    <numa>
      <cell id='0' cpus='0-3' memory='4194304' unit='KiB'/>
      <cell id='1' cpus='4-7' memory='4194304' unit='KiB'/>
    </numa>
  </cpu>
  <numatune>
    <memnode cellid='0' mode='strict' nodeset='0'/>
    <memnode cellid='1' mode='strict' nodeset='1'/>
  </numatune>
</domain>`

// virshXMLDomstatus mirrors the format of /run/libvirt/qemu/*.xml files.
const virshXMLDomstatus = `<domstatus state='running' reason='booted' pid='98765'>
  <domain type='kvm' xmlns:nova='http://openstack.org/xmlns/libvirt/nova/1.0'>
    <name>instance-0000001a</name>
    <nova:instance>
      <nova:name>my-vm-01</nova:name>
    </nova:instance>
    <cpu>
      <numa>
        <cell id='0' cpus='0-3' memory='4194304' unit='KiB'/>
      </numa>
    </cpu>
    <numatune>
      <memnode cellid='0' mode='strict' nodeset='0'/>
    </numatune>
  </domain>
</domstatus>`

// virshXMLNotPinned has no <numatune><memnode> — not NUMA-pinned.
const virshXMLNotPinned = `<domain type='kvm'>
  <name>instance-0000002b</name>
  <cpu><topology sockets='1' cores='4' threads='2'/></cpu>
</domain>`

// virshXMLFallbackName tests that libvirt domain name is used when nova:name is absent.
const virshXMLFallbackName = `<domain type='kvm'>
  <name>instance-0000003c</name>
  <cpu>
    <numa>
      <cell id='0' cpus='0-1' memory='2097152' unit='KiB'/>
    </numa>
  </cpu>
  <numatune>
    <memnode cellid='0' mode='strict' nodeset='0'/>
  </numatune>
</domain>`

func TestParseVirshXMLNUMAPinned(t *testing.T) {
	res, err := ParseVirshXML(virshXMLNUMAPinned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result for NUMA-pinned VM")
	}
	if res.VMName != "my-vm-01" {
		t.Errorf("VMName = %q, want %q", res.VMName, "my-vm-01")
	}

	n0, ok := res.HostNUMAUsage[0]
	if !ok {
		t.Fatal("host NUMA node 0 not found")
	}
	if n0[0] != 4 { // CountCPUList("0-3") = 4
		t.Errorf("node 0 vCPUs = %d, want 4", n0[0])
	}
	if n0[1] != 4194304*1024 { // 4194304 KiB in bytes
		t.Errorf("node 0 memBytes = %d, want %d", n0[1], int64(4194304*1024))
	}

	n1, ok := res.HostNUMAUsage[1]
	if !ok {
		t.Fatal("host NUMA node 1 not found")
	}
	if n1[0] != 4 {
		t.Errorf("node 1 vCPUs = %d, want 4", n1[0])
	}
	if n1[1] != 4194304*1024 {
		t.Errorf("node 1 memBytes = %d, want %d", n1[1], int64(4194304*1024))
	}
}

func TestParseVirshXMLDomstatus(t *testing.T) {
	res, err := ParseVirshXML(virshXMLDomstatus)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result for domstatus NUMA-pinned VM")
	}
	if res.VMName != "my-vm-01" {
		t.Errorf("VMName = %q, want %q", res.VMName, "my-vm-01")
	}
	n0, ok := res.HostNUMAUsage[0]
	if !ok {
		t.Fatal("host NUMA node 0 not found")
	}
	if n0[0] != 4 {
		t.Errorf("node 0 vCPUs = %d, want 4", n0[0])
	}
}

func TestParseVirshXMLNotPinned(t *testing.T) {
	res, err := ParseVirshXML(virshXMLNotPinned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Errorf("expected nil for non-NUMA-pinned VM, got %+v", res)
	}
}

func TestParseVirshXMLFallbackName(t *testing.T) {
	res, err := ParseVirshXML(virshXMLFallbackName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if res.VMName != "instance-0000003c" {
		t.Errorf("VMName = %q, want %q", res.VMName, "instance-0000003c")
	}
}

func TestParseVirshXMLInvalid(t *testing.T) {
	_, err := ParseVirshXML("<not valid xml")
	if err == nil {
		t.Error("expected error for invalid XML")
	}
}

func TestParseVirshXMLMemoryUnits(t *testing.T) {
	tests := []struct {
		unit string
		want int64
	}{
		{unit: "bytes", want: 2},
		{unit: "KB", want: 2000},
		{unit: "KiB", want: 2 * 1024},
		{unit: "MB", want: 2 * 1000 * 1000},
		{unit: "MiB", want: 2 * 1024 * 1024},
		{unit: "GB", want: 2 * 1000 * 1000 * 1000},
		{unit: "GiB", want: 2 * 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			xml := `<domain><name>vm</name><cpu><numa><cell id='0' cpus='0' memory='2' unit='` + tt.unit + `'/></numa></cpu><numatune><memnode cellid='0' mode='strict' nodeset='0'/></numatune></domain>`
			res, err := ParseVirshXML(xml)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := res.HostNUMAUsage[0][1]; got != tt.want {
				t.Errorf("memory = %d bytes, want %d", got, tt.want)
			}
		})
	}
}

func TestParseVirshXMLSkipsMultipleNodeSet(t *testing.T) {
	xml := `<domain><name>vm</name><cpu><numa><cell id='0' cpus='0-1' memory='2' unit='GiB'/></numa></cpu><numatune><memnode cellid='0' mode='interleave' nodeset='0-1'/></numatune></domain>`
	res, err := ParseVirshXML(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.HostNUMAUsage) != 0 {
		t.Errorf("expected no host attribution for multiple-node set, got %+v", res.HostNUMAUsage)
	}
}

func TestParseVirshXMLReturnsDomstatusError(t *testing.T) {
	_, err := ParseVirshXML("<unexpected/>")
	if err == nil {
		t.Fatal("expected error for unexpected root element")
	}
	if !strings.Contains(err.Error(), "domstatus") {
		t.Errorf("error = %q, want domstatus unmarshal error", err)
	}
}
