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
	"strconv"
	"testing"
)

func TestSockStats(t *testing.T) {
	testSockStats(t, "fixtures/proc/net/sockstat")
	testSockStats(t, "fixtures/proc/net/sockstat_rhe4")
}

func testSockStats(t *testing.T, fixture string) {
	file, err := os.Open(fixture)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	sockStats, err := parseSockStats(file, fixture)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "229", sockStats["sockets"]["used"]; want != got {
		t.Errorf("want sockstat sockets used %s, got %s", want, got)
	}

	if want, got := "4", sockStats["TCP"]["tw"]; want != got {
		t.Errorf("want sockstat TCP tw %s, got %s", want, got)
	}

	if want, got := "17", sockStats["TCP"]["alloc"]; want != got {
		t.Errorf("want sockstat TCP alloc %s, got %s", want, got)
	}

	// The test file has 1 for TCP mem, which is one page.  So we should get the
	// page size in bytes back from sockstat_linux.  We get the page size from
	// os here because this value can change from system to system.  The value is
	// 4096 by default from linux 2.4 onward.
	if want, got := strconv.Itoa(os.Getpagesize()), sockStats["TCP"]["mem_bytes"]; want != got {
		t.Errorf("want sockstat TCP mem_bytes %s, got %s", want, got)
	}
}
