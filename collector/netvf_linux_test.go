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
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/jsimonetti/rtnetlink/v2"
	"github.com/prometheus/procfs/sysfs"
)

func uint32Ptr(v uint32) *uint32 {
	return &v
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

		if link.Attributes.NumVF == nil || *link.Attributes.NumVF == 0 {
			continue
		}

		device := link.Attributes.Name

		if filter.ignored(device) {
			logger.Debug("Ignoring device", "device", device)
			continue
		}

		for _, vf := range link.Attributes.VFInfoList {
			mac := ""
			if vf.MAC != nil {
				mac = vf.MAC.String()
			}

			pciAddress := ""
			if sysClassPath != "" {
				if fs, err := sysfs.NewFS(sysClassPath); err == nil {
					if dev, err := fs.NetClassPCIDevice(device); err == nil {
						pciAddress, _ = fs.PciDeviceVFAddress(dev, vf.ID)
					}
				}
			}
			result = append(result, vfMetrics{
				Device:     device,
				VFID:       vf.ID,
				MAC:        mac,
				Vlan:       vf.Vlan,
				LinkState:  vfLinkStateString(vf.LinkState),
				SpoofCheck: vf.SpoofCheck,
				Trust:      vf.Trust,
				PCIAddress: pciAddress,
				Stats:      vf.Stats,
			})
		}
	}

	return result
}

var vfLinks = []rtnetlink.LinkMessage{
	{
		// Interface without VFs
		Attributes: &rtnetlink.LinkAttributes{
			Name: "eth0",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 1000,
				TXPackets: 2000,
			},
		},
	},
	{
		// Interface with NumVF = 0
		Attributes: &rtnetlink.LinkAttributes{
			Name:  "eth1",
			NumVF: uint32Ptr(0),
		},
	},
	{
		// PF with 2 VFs
		Attributes: &rtnetlink.LinkAttributes{
			Name:  "enp3s0f0",
			NumVF: uint32Ptr(2),
			VFInfoList: []rtnetlink.VFInfo{
				{
					ID:         0,
					MAC:        net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0x00},
					Vlan:       100,
					LinkState:  rtnetlink.VFLinkStateAuto,
					SpoofCheck: true,
					Trust:      false,
					Stats: &rtnetlink.VFStats{
						RxPackets: 1000,
						TxPackets: 2000,
						RxBytes:   100000,
						TxBytes:   200000,
						Broadcast: 10,
						Multicast: 20,
						RxDropped: 5,
						TxDropped: 3,
					},
				},
				{
					ID:         1,
					MAC:        net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0x01},
					Vlan:       200,
					LinkState:  rtnetlink.VFLinkStateEnable,
					SpoofCheck: false,
					Trust:      true,
					Stats: &rtnetlink.VFStats{
						RxPackets: 3000,
						TxPackets: 4000,
						RxBytes:   300000,
						TxBytes:   400000,
						Broadcast: 30,
						Multicast: 40,
						RxDropped: 7,
						TxDropped: 9,
					},
				},
			},
		},
	},
	{
		// Another PF with 1 VF (no stats)
		Attributes: &rtnetlink.LinkAttributes{
			Name:  "enp3s0f1",
			NumVF: uint32Ptr(1),
			VFInfoList: []rtnetlink.VFInfo{
				{
					ID:         0,
					MAC:        net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
					Vlan:       0,
					LinkState:  rtnetlink.VFLinkStateDisable,
					SpoofCheck: true,
					Trust:      false,
					Stats:      nil, // No stats available
				},
			},
		},
	},
	{
		// Nil attributes (should be skipped)
		Attributes: nil,
	},
}

func TestParseVFInfo(t *testing.T) {
	filter := newDeviceFilter("", "")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	vfs := parseVFInfo(vfLinks, &filter, logger, "")

	// Should have 3 VFs total (2 from enp3s0f0, 1 from enp3s0f1)
	if want, got := 3, len(vfs); want != got {
		t.Errorf("want %d VFs, got %d", want, got)
	}
}

func TestParseVFInfoDeviceFilter(t *testing.T) {
	// Exclude enp3s0f1
	filter := newDeviceFilter("enp3s0f1", "")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	vfs := parseVFInfo(vfLinks, &filter, logger, "")

	// Should have 2 VFs (only from enp3s0f0)
	if want, got := 2, len(vfs); want != got {
		t.Errorf("want %d VFs, got %d", want, got)
	}

	for _, vf := range vfs {
		if vf.Device == "enp3s0f1" {
			t.Error("enp3s0f1 should be filtered out")
		}
	}
}

func TestParseVFInfoDeviceInclude(t *testing.T) {
	// Only include enp3s0f1
	filter := newDeviceFilter("", "^enp3s0f1$")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	vfs := parseVFInfo(vfLinks, &filter, logger, "")

	// Should have 1 VF (only from enp3s0f1)
	if want, got := 1, len(vfs); want != got {
		t.Errorf("want %d VFs, got %d", want, got)
	}

	if len(vfs) > 0 && vfs[0].Device != "enp3s0f1" {
		t.Errorf("want device enp3s0f1, got %s", vfs[0].Device)
	}
}

func TestParseVFInfoStats(t *testing.T) {
	filter := newDeviceFilter("", "^enp3s0f0$")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	vfs := parseVFInfo(vfLinks, &filter, logger, "")

	if len(vfs) != 2 {
		t.Fatalf("expected 2 VFs, got %d", len(vfs))
	}

	// Check VF 0 stats
	vf0 := vfs[0]
	if vf0.VFID != 0 {
		t.Errorf("expected VF ID 0, got %d", vf0.VFID)
	}
	if vf0.Stats == nil {
		t.Fatal("expected stats for VF 0")
	}
	if want, got := uint64(1000), vf0.Stats.RxPackets; want != got {
		t.Errorf("want RxPackets %d, got %d", want, got)
	}
	if want, got := uint64(200000), vf0.Stats.TxBytes; want != got {
		t.Errorf("want TxBytes %d, got %d", want, got)
	}

	// Check VF 1 stats
	vf1 := vfs[1]
	if vf1.VFID != 1 {
		t.Errorf("expected VF ID 1, got %d", vf1.VFID)
	}
	if want, got := uint64(4000), vf1.Stats.TxPackets; want != got {
		t.Errorf("want TxPackets %d, got %d", want, got)
	}
}

func TestParseVFInfoMetadata(t *testing.T) {
	filter := newDeviceFilter("", "^enp3s0f0$")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	vfs := parseVFInfo(vfLinks, &filter, logger, "")

	if len(vfs) != 2 {
		t.Fatalf("expected 2 VFs, got %d", len(vfs))
	}

	// Check VF 0 metadata
	vf0 := vfs[0]
	if want, got := "aa:bb:cc:dd:ee:00", vf0.MAC; want != got {
		t.Errorf("want MAC %s, got %s", want, got)
	}
	if want, got := uint32(100), vf0.Vlan; want != got {
		t.Errorf("want VLAN %d, got %d", want, got)
	}
	if want, got := "auto", vf0.LinkState; want != got {
		t.Errorf("want LinkState %s, got %s", want, got)
	}
	if !vf0.SpoofCheck {
		t.Error("expected SpoofCheck to be true")
	}
	if vf0.Trust {
		t.Error("expected Trust to be false")
	}

	// Check VF 1 metadata
	vf1 := vfs[1]
	if want, got := "enable", vf1.LinkState; want != got {
		t.Errorf("want LinkState %s, got %s", want, got)
	}
	if vf1.SpoofCheck {
		t.Error("expected SpoofCheck to be false")
	}
	if !vf1.Trust {
		t.Error("expected Trust to be true")
	}
}

func TestParseVFInfoNoStats(t *testing.T) {
	filter := newDeviceFilter("", "^enp3s0f1$")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	vfs := parseVFInfo(vfLinks, &filter, logger, "")

	if len(vfs) != 1 {
		t.Fatalf("expected 1 VF, got %d", len(vfs))
	}

	vf := vfs[0]
	if vf.Stats != nil {
		t.Error("expected Stats to be nil for this VF")
	}
	if want, got := "disable", vf.LinkState; want != got {
		t.Errorf("want LinkState %s, got %s", want, got)
	}
}

func TestVFLinkStateString(t *testing.T) {
	tests := []struct {
		state    rtnetlink.VFLinkState
		expected string
	}{
		{rtnetlink.VFLinkStateAuto, "auto"},
		{rtnetlink.VFLinkStateEnable, "enable"},
		{rtnetlink.VFLinkStateDisable, "disable"},
		{rtnetlink.VFLinkState(99), "unknown"},
	}

	for _, tt := range tests {
		got := vfLinkStateString(tt.state)
		if got != tt.expected {
			t.Errorf("vfLinkStateString(%d) = %s, want %s", tt.state, got, tt.expected)
		}
	}
}

func TestResolveVFPCIAddress(t *testing.T) {
	// Create a fake sysfs tree with:
	// - class/net/enp3s0f0/device -> symlink to the PCI device
	// - bus/pci/devices/0000:00:01.0 -> symlink to the real device path
	// - devices/pci0000:00/0000:00:01.0/ with required PCI files and virtfn0
	tmp := t.TempDir()

	// Real device directory with required PCI files and virtfn0 symlink
	pciDevDir := filepath.Join(tmp, "devices", "pci0000:00", "0000:00:01.0")
	if err := os.MkdirAll(pciDevDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, f := range []struct{ name, val string }{
		{"class", "0x020000"}, {"vendor", "0x8086"}, {"device", "0x1572"},
		{"subsystem_vendor", "0x8086"}, {"subsystem_device", "0x0000"}, {"revision", "0x00"},
	} {
		if err := os.WriteFile(filepath.Join(pciDevDir, f.name), []byte(f.val), 0o444); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink("../0000:00:02.0", filepath.Join(pciDevDir, "virtfn0")); err != nil {
		t.Fatal(err)
	}

	// bus/pci/devices/0000:00:01.0 -> symlink to real device
	busDir := filepath.Join(tmp, "bus", "pci", "devices")
	if err := os.MkdirAll(busDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../../../devices/pci0000:00/0000:00:01.0", filepath.Join(busDir, "0000:00:01.0")); err != nil {
		t.Fatal(err)
	}

	// class/net/enp3s0f0/device -> symlink to the PCI device
	netDir := filepath.Join(tmp, "class", "net", "enp3s0f0")
	if err := os.MkdirAll(netDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../../../devices/pci0000:00/0000:00:01.0", filepath.Join(netDir, "device")); err != nil {
		t.Fatal(err)
	}

	fs, err := sysfs.NewFS(tmp)
	if err != nil {
		t.Fatal(err)
	}
	dev, err := fs.NetClassPCIDevice("enp3s0f0")
	if err != nil {
		t.Fatal(err)
	}
	got, err := fs.PciDeviceVFAddress(dev, 0)
	if err != nil {
		t.Fatal(err)
	}
	if want := "0000:00:02.0"; got != want {
		t.Errorf("PciDeviceVFAddress() = %q, want %q", got, want)
	}
}

func TestResolveVFPCIAddressMissing(t *testing.T) {
	tmp := t.TempDir()

	fs, err := sysfs.NewFS(tmp)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fs.NetClassPCIDevice("enp3s0f0")
	if err == nil {
		t.Error("expected error for missing interface, got nil")
	}
}
