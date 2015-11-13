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

func TestMemInfoNuma(t *testing.T) {
	file, err := os.Open("fixtures/sys/devices/system/node/node0/meminfo")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	memInfo, err := parseMemInfoNuma(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 707915776.0, memInfo[meminfoKey{"Active_anon", "0"}]; want != got {
		t.Errorf("want memory Active(anon) %f, got %f", want, got)
	}

	if want, got := 150994944.0, memInfo[meminfoKey{"AnonHugePages", "0"}]; want != got {
		t.Errorf("want memory AnonHugePages %f, got %f", want, got)
	}

	file, err = os.Open("fixtures/sys/devices/system/node/node1/meminfo")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	memInfo, err = parseMemInfoNuma(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 291930112.0, memInfo[meminfoKey{"Inactive_anon", "1"}]; want != got {
		t.Errorf("want memory Inactive(anon) %f, got %f", want, got)
	}

	if want, got := 85585088512.0, memInfo[meminfoKey{"FilePages", "1"}]; want != got {
		t.Errorf("want memory FilePages %f, got %f", want, got)
	}
}
