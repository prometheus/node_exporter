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
	"encoding/xml"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type domain struct {
	XMLName      xml.Name     `xml:"domain"`
	Name         string       `xml:"name"`
	NovaInstance novaInstance `xml:"http://openstack.org/xmlns/libvirt/nova/1.0 instance"`
	CPU          domainCPU    `xml:"cpu"`
	NumaTune     numaTune     `xml:"numatune"`
}

type novaInstance struct {
	Name string `xml:"http://openstack.org/xmlns/libvirt/nova/1.0 name"`
}

type domainCPU struct {
	NUMA domainNUMA `xml:"numa"`
}

type domainNUMA struct {
	Cells []numaCell `xml:"cell"`
}

type numaCell struct {
	ID     string `xml:"id,attr"`
	CPUs   string `xml:"cpus,attr"`
	Memory string `xml:"memory,attr"`
	Unit   string `xml:"unit,attr"`
}

type numaTune struct {
	MemNodes []memNode `xml:"memnode"`
}

type memNode struct {
	CellID  string `xml:"cellid,attr"`
	Mode    string `xml:"mode,attr"`
	NodeSet string `xml:"nodeset,attr"`
}

type domstatusWrapper struct {
	XMLName xml.Name `xml:"domstatus"`
	Domain  domain   `xml:"domain"`
}

// ParseVirshXML parses libvirt domain XML in either format:
//   - <domain> root (output of "virsh dumpxml")
//   - <domstatus> root (files in /run/libvirt/qemu/*.xml)
//
// Returns (nil, nil) when the domain has no explicit NUMA pinning
// (no <memnode> elements in <numatune>).
// Returns (nil, err) on XML parse errors.
func ParseVirshXML(xmlStr string) (*VirshResult, error) {
	var d domain
	if err := xml.Unmarshal([]byte(xmlStr), &d); err != nil {
		var ds domstatusWrapper
		if err2 := xml.Unmarshal([]byte(xmlStr), &ds); err2 != nil {
			return nil, fmt.Errorf("parsing domain XML: %w", err2)
		}
		d = ds.Domain
	}

	if len(d.NumaTune.MemNodes) == 0 {
		return nil, nil
	}

	vmName := strings.TrimSpace(d.NovaInstance.Name)
	if vmName == "" {
		vmName = strings.TrimSpace(d.Name)
	}

	cellToHostNode := make(map[int]int, len(d.NumaTune.MemNodes))
	for _, mn := range d.NumaTune.MemNodes {
		cellID, err := strconv.Atoi(strings.TrimSpace(mn.CellID))
		if err != nil {
			continue
		}
		hostNode, err := strconv.Atoi(strings.TrimSpace(mn.NodeSet))
		if err != nil {
			// Multiple-node sets cannot be represented by VirshResult without
			// incorrectly attributing all of the cell's resources to one node.
			continue
		}
		cellToHostNode[cellID] = hostNode
	}

	hostNUMAUsage := make(map[int][2]int64)
	for _, cell := range d.CPU.NUMA.Cells {
		cellID, err := strconv.Atoi(strings.TrimSpace(cell.ID))
		if err != nil {
			continue
		}
		hostNode, ok := cellToHostNode[cellID]
		if !ok {
			continue
		}
		cpuCount := int64(CountCPUList(cell.CPUs))
		memBytes, err := memoryBytes(cell.Memory, cell.Unit)
		if err != nil {
			return nil, fmt.Errorf("parsing memory for NUMA cell %q: %w", cell.ID, err)
		}

		prev := hostNUMAUsage[hostNode]
		hostNUMAUsage[hostNode] = [2]int64{prev[0] + cpuCount, prev[1] + memBytes}
	}

	return &VirshResult{
		VMName:        vmName,
		HostNUMAUsage: hostNUMAUsage,
	}, nil
}

func memoryBytes(value, unit string) (int64, error) {
	memory, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid value %q: %w", value, err)
	}
	if memory < 0 {
		return 0, fmt.Errorf("invalid negative value %q", value)
	}

	var multiplier int64
	switch strings.TrimSpace(unit) {
	case "bytes":
		multiplier = 1
	case "KB":
		multiplier = 1000
	case "", "KiB":
		multiplier = 1024
	case "MB":
		multiplier = 1000 * 1000
	case "MiB":
		multiplier = 1024 * 1024
	case "GB":
		multiplier = 1000 * 1000 * 1000
	case "GiB":
		multiplier = 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unsupported unit %q", unit)
	}

	if memory > math.MaxInt64/multiplier {
		return 0, fmt.Errorf("value %q %s overflows bytes", value, unit)
	}
	return memory * multiplier, nil
}
