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
	"regexp"
	"testing"
)

func TestNetDevStats(t *testing.T) {
	file, err := os.Open("fixtures/proc/net/dev")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	netStats, err := parseNetDevStats(file, regexp.MustCompile("^veth"))
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "10437182923", netStats["wlan0"]["receive_bytes"]; want != got {
		t.Errorf("want netstat wlan0 bytes %s, got %s", want, got)
	}

	if want, got := "68210035552", netStats["eth0"]["receive_bytes"]; want != got {
		t.Errorf("want netstat eth0 bytes %s, got %s", want, got)
	}

	if want, got := "934", netStats["tun0"]["transmit_packets"]; want != got {
		t.Errorf("want netstat tun0 packets %s, got %s", want, got)
	}

	if want, got := 6, len(netStats); want != got {
		t.Errorf("want count of devices to be %d, got %d", want, got)
	}

	if _, ok := netStats["veth4B09XN"]["transmit_bytes"]; ok {
		t.Error("want fixture interface veth4B09XN to not exist, but it does")
	}
}
