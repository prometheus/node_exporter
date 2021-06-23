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
	"os"
	"testing"

	"github.com/go-kit/log"
)

func TestNetDevStatsIgnore(t *testing.T) {
	file, err := os.Open("fixtures/proc/net/dev")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	filter := newNetDevFilter("^veth", "")

	netStats, err := parseNetDevStats(file, &filter, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}

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
	file, err := os.Open("fixtures/proc/net/dev")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	filter := newNetDevFilter("", "^💩0$")
	netStats, err := parseNetDevStats(file, &filter, log.NewNopLogger())
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 1, len(netStats); want != got {
		t.Errorf("want count of devices to be %d, got %d", want, got)
	}
	if want, got := uint64(72), netStats["💩0"]["receive_multicast"]; want != got {
		t.Error("want fixture interface 💩0 to exist, but it does not")
	}
}
