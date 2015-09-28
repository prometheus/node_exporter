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
)

func TestTCPStat(t *testing.T) {
	file, err := os.Open("fixtures/proc/net/tcpstat")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	tcpStats, err := parseTCPStats(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 1, int(tcpStats[TCP_ESTABLISHED]); want != got {
		t.Errorf("want tcpstat number of established state %d, got %d", want, got)
	}

	if want, got := 1, int(tcpStats[TCP_LISTEN]); want != got {
		t.Errorf("want tcpstat number of listen state %d, got %d", want, got)
	}
}
