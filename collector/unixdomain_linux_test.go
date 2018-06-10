// Copyright 2018 The Prometheus Authors
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

	"github.com/prometheus/client_golang/prometheus"
)

func TestUnixDomain(t *testing.T) {

	f, err := os.Open("fixtures/proc/net/unix")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	collector, err := NewUnixDomainCollector()
	if err != nil {
		t.Fatal(err)
	}
	resultCh := make(chan prometheus.Metric, 200)

	udCollector := collector.(*unixDomainCollector)
	err = udCollector.process(f, resultCh)
	if err != nil {
		t.Fatal(err)
	}
	// total 4 combinations of the tags in the fixtures
	expectedSize := 4 * 2
	gotSize := len(resultCh)
	if gotSize != expectedSize {
		t.Errorf("expected Unix domain results count %d, got %d", expectedSize, len(resultCh))
	}
	for i := 0; i < gotSize; i++ {
		<-resultCh
	}

	// open another file which contains less combinations than the previous one
	// it should have same number of results
	f, err = os.Open("fixtures/proc/net/unix2")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = udCollector.process(f, resultCh)
	if err != nil {
		t.Fatal(err)
	}
	gotSize = len(resultCh)
	if gotSize != expectedSize {
		t.Errorf("expected Unix domain results count %d, got %d", expectedSize, len(resultCh))
	}
}

func TestUnixDomainParsing(t *testing.T) {
	collector, err := NewUnixDomainCollector()
	if err != nil {
		t.Fatal(err)
	}
	udCollector := collector.(*unixDomainCollector)

	testDatas := []struct {
		line          string
		expectedType  string
		expectedFlags string
		expectedState string
		expectedUsers int64
	}{
		{
			"0000000000000000: 00000002 00000000 00010000 0001 01 3990228 @/containerd-shim/moby/16b1cd05b9b6e6268fc6cc3271fc09a6384dac585f1d0ba419ab1afebe99ae95/shim.sock@",
			"stream", "accepton", "unconnected", 2,
		},
		{
			"0000000000000000: 00000003 00000000 00000000 0001 03 4787297 /var/run/postgresql/.s.PGSQL.5432",
			"stream", "default", "connected", 3,
		},
		{
			"0000000000000000: 0000000f 00000000 00000000 0002 02 12392 /dev/log",
			"dgram", "default", "connecting", 15,
		},
	}
	for _, td := range testDatas {
		entry, err := udCollector.parseItem(td.line)
		if err != nil {
			t.Errorf("failed to parse unix domain entry %s: %s", td.line, err)
			continue
		}
		if entry.labelsComb.typ != td.expectedType {
			t.Errorf("unix domain entry %s, expected type %s, got %s", td.line, td.expectedType, entry.labelsComb.typ)
		}
		if entry.labelsComb.flags != td.expectedFlags {
			t.Errorf("unix domain entry %s, expected flags %s, got %s", td.line, td.expectedFlags, entry.labelsComb.flags)
		}
		if entry.labelsComb.state != td.expectedState {
			t.Errorf("unix domain entry %s, expected state %s, got %s", td.line, td.expectedState, entry.labelsComb.state)
		}
	}
}
