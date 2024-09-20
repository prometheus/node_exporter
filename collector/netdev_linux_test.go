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

//go:build !nonetdev
// +build !nonetdev

package collector

import (
	"io"
	"log/slog"
	"testing"

	"github.com/jsimonetti/rtnetlink/v2"
)

var links = []rtnetlink.LinkMessage{
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "tun0",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 24,
				TXPackets: 934,
				RXBytes:   1888,
				TXBytes:   67120,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "veth4B09XN",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 8,
				TXPackets: 10640,
				RXBytes:   648,
				TXBytes:   1943284,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "lo",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 1832522,
				TXPackets: 1832522,
				RXBytes:   435303245,
				TXBytes:   435303245,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "eth0",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 520993275,
				TXPackets: 43451486,
				RXBytes:   68210035552,
				TXBytes:   9315587528,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "lxcbr0",
			Stats64: &rtnetlink.LinkStats64{
				TXPackets: 28339,
				TXBytes:   2630299,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "wlan0",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 13899359,
				TXPackets: 11726200,
				RXBytes:   10437182923,
				TXBytes:   2851649360,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "docker0",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 1065585,
				TXPackets: 1929779,
				RXBytes:   64910168,
				TXBytes:   2681662018,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name:    "ibr10:30",
			Stats64: &rtnetlink.LinkStats64{},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "flannel.1",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 228499337,
				TXPackets: 258369223,
				RXBytes:   18144009813,
				TXBytes:   20758990068,
				TXDropped: 64,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "💩0",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets: 105557,
				TXPackets: 304261,
				RXBytes:   57750104,
				TXBytes:   404570255,
				Multicast: 72,
			},
		},
	},
	{
		Attributes: &rtnetlink.LinkAttributes{
			Name: "enp0s0f0",
			Stats64: &rtnetlink.LinkStats64{
				RXPackets:         226,
				TXPackets:         803,
				RXBytes:           231424,
				TXBytes:           822272,
				RXErrors:          14,
				TXErrors:          2,
				RXDropped:         10,
				TXDropped:         17,
				Multicast:         285,
				Collisions:        30,
				RXLengthErrors:    5,
				RXOverErrors:      3,
				RXCRCErrors:       1,
				RXFrameErrors:     4,
				RXFIFOErrors:      6,
				RXMissedErrors:    21,
				TXAbortedErrors:   22,
				TXCarrierErrors:   7,
				TXFIFOErrors:      24,
				TXHeartbeatErrors: 9,
				TXWindowErrors:    19,
				RXCompressed:      23,
				TXCompressed:      20,
				RXNoHandler:       62,
			},
		},
	},
}

func TestNetDevStatsIgnore(t *testing.T) {
	filter := newDeviceFilter("^veth", "")

	netStats := parseNetlinkStats(links, &filter, slog.New(slog.NewTextHandler(io.Discard, nil)))

	if want, got := uint64(10437182923), netStats["wlan0"]["receive_bytes"]; want != got {
		t.Errorf("want netstat wlan0 bytes %v, got %v", want, got)
	}

	if want, got := uint64(68210035552), netStats["eth0"]["receive_bytes"]; want != got {
		t.Errorf("want netstat eth0 bytes %v, got %v", want, got)
	}

	if want, got := uint64(934), netStats["tun0"]["transmit_packets"]; want != got {
		t.Errorf("want netstat tun0 packets %v, got %v", want, got)
	}

	if want, got := 10, len(netStats); want != got {
		t.Errorf("want count of devices to be %d, got %d", want, got)
	}

	if _, ok := netStats["veth4B09XN"]["transmit_bytes"]; ok {
		t.Error("want fixture interface veth4B09XN to not exist, but it does")
	}

	if want, got := uint64(0), netStats["ibr10:30"]["receive_fifo"]; want != got {
		t.Error("want fixture interface ibr10:30 to exist, but it does not")
	}

	if want, got := uint64(72), netStats["💩0"]["multicast"]; want != got {
		t.Error("want fixture interface 💩0 to exist, but it does not")
	}
}

func TestNetDevStatsAccept(t *testing.T) {
	filter := newDeviceFilter("", "^💩0$")
	netStats := parseNetlinkStats(links, &filter, slog.New(slog.NewTextHandler(io.Discard, nil)))

	if want, got := 1, len(netStats); want != got {
		t.Errorf("want count of devices to be %d, got %d", want, got)
	}
	if want, got := uint64(72), netStats["💩0"]["multicast"]; want != got {
		t.Error("want fixture interface 💩0 to exist, but it does not")
	}
}

func TestNetDevLegacyMetricNames(t *testing.T) {
	expected := []string{
		"receive_packets",
		"transmit_packets",
		"receive_bytes",
		"transmit_bytes",
		"receive_errs",
		"transmit_errs",
		"receive_drop",
		"transmit_drop",
		"receive_multicast",
		"transmit_colls",
		"receive_frame",
		"receive_fifo",
		"transmit_carrier",
		"transmit_fifo",
		"receive_compressed",
		"transmit_compressed",
	}

	filter := newDeviceFilter("", "")
	netStats := parseNetlinkStats(links, &filter, slog.New(slog.NewTextHandler(io.Discard, nil)))

	for dev, devStats := range netStats {
		legacy(devStats)
		for _, name := range expected {
			if _, ok := devStats[name]; !ok {
				t.Errorf("metric %s should be defined on interface %s", name, dev)
			}
		}
	}
}

func TestNetDevLegacyMetricValues(t *testing.T) {
	expected := map[string]uint64{
		"receive_packets":     226,
		"transmit_packets":    803,
		"receive_bytes":       231424,
		"transmit_bytes":      822272,
		"receive_errs":        14,
		"transmit_errs":       2,
		"receive_drop":        10 + 21,
		"transmit_drop":       17,
		"receive_multicast":   285,
		"transmit_colls":      30,
		"receive_frame":       5 + 3 + 1 + 4,
		"receive_fifo":        6,
		"transmit_carrier":    22 + 7 + 9 + 19,
		"transmit_fifo":       24,
		"receive_compressed":  23,
		"transmit_compressed": 20,
	}

	filter := newDeviceFilter("", "^enp0s0f0$")
	netStats := parseNetlinkStats(links, &filter, slog.New(slog.NewTextHandler(io.Discard, nil)))
	metrics, ok := netStats["enp0s0f0"]
	if !ok {
		t.Error("expected stats for interface enp0s0f0")
	}

	legacy(metrics)

	for name, want := range expected {
		got, ok := metrics[name]
		if !ok {
			t.Errorf("metric %s should be defined on interface enp0s0f0", name)
			continue
		}
		if want != got {
			t.Errorf("want %s %d, got %d", name, want, got)
		}
	}
}

func TestNetDevMetricValues(t *testing.T) {
	filter := newDeviceFilter("", "")
	netStats := parseNetlinkStats(links, &filter, slog.New(slog.NewTextHandler(io.Discard, nil)))

	for _, msg := range links {
		device := msg.Attributes.Name
		stats := msg.Attributes.Stats64

		expected := map[string]uint64{
			"receive_packets":  stats.RXPackets,
			"transmit_packets": stats.TXPackets,
			"receive_bytes":    stats.RXBytes,
			"transmit_bytes":   stats.TXBytes,
			"receive_errors":   stats.RXErrors,
			"transmit_errors":  stats.TXErrors,
			"receive_dropped":  stats.RXDropped,
			"transmit_dropped": stats.TXDropped,
			"multicast":        stats.Multicast,
			"collisions":       stats.Collisions,

			// detailed rx_errors
			"receive_length_errors": stats.RXLengthErrors,
			"receive_over_errors":   stats.RXOverErrors,
			"receive_crc_errors":    stats.RXCRCErrors,
			"receive_frame_errors":  stats.RXFrameErrors,
			"receive_fifo_errors":   stats.RXFIFOErrors,
			"receive_missed_errors": stats.RXMissedErrors,

			// detailed tx_errors
			"transmit_aborted_errors":   stats.TXAbortedErrors,
			"transmit_carrier_errors":   stats.TXCarrierErrors,
			"transmit_fifo_errors":      stats.TXFIFOErrors,
			"transmit_heartbeat_errors": stats.TXHeartbeatErrors,
			"transmit_window_errors":    stats.TXWindowErrors,

			// for cslip etc
			"receive_compressed":  stats.RXCompressed,
			"transmit_compressed": stats.TXCompressed,
			"receive_nohandler":   stats.RXNoHandler,
		}

		for name, want := range expected {
			devStats, ok := netStats[device]
			if !ok {
				t.Errorf("expected stats for interface %s", device)
			}
			got, ok := devStats[name]
			if !ok {
				t.Errorf("metric %s should be defined on interface %s", name, device)
			}
			if want != got {
				t.Errorf("want %s %d, got %d", name, want, got)
			}
		}
	}
}
