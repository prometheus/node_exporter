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

package collector

import (
	"testing"

	"github.com/go-kit/log"

	"github.com/jsimonetti/rtnetlink"
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
}

func TestNetDevStatsIgnore(t *testing.T) {
	filter := newDeviceFilter("^veth", "")

	netStats := netlinkStats(links, &filter, log.NewNopLogger())

	if want, got := uint64(10437182923), netStats["wlan0"]["receive_bytes"]; want != got {
		t.Errorf("want netstat wlan0 bytes %v, got %v", want, got)
	}

	if want, got := uint64(68210035552), netStats["eth0"]["receive_bytes"]; want != got {
		t.Errorf("want netstat eth0 bytes %v, got %v", want, got)
	}

	if want, got := uint64(934), netStats["tun0"]["transmit_packets"]; want != got {
		t.Errorf("want netstat tun0 packets %v, got %v", want, got)
	}

	if want, got := 9, len(netStats); want != got {
		t.Errorf("want count of devices to be %d, got %d", want, got)
	}

	if _, ok := netStats["veth4B09XN"]["transmit_bytes"]; ok {
		t.Error("want fixture interface veth4B09XN to not exist, but it does")
	}

	if want, got := uint64(0), netStats["ibr10:30"]["receive_fifo"]; want != got {
		t.Error("want fixture interface ibr10:30 to exist, but it does not")
	}

	if want, got := uint64(72), netStats["💩0"]["receive_multicast"]; want != got {
		t.Error("want fixture interface 💩0 to exist, but it does not")
	}
}

func TestNetDevStatsAccept(t *testing.T) {
	filter := newDeviceFilter("", "^💩0$")
	netStats := netlinkStats(links, &filter, log.NewNopLogger())

	if want, got := 1, len(netStats); want != got {
		t.Errorf("want count of devices to be %d, got %d", want, got)
	}
	if want, got := uint64(72), netStats["💩0"]["receive_multicast"]; want != got {
		t.Error("want fixture interface 💩0 to exist, but it does not")
	}
}
