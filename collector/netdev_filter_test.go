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
	"testing"
)

func TestNetDevFilter(t *testing.T) {
	tests := []struct {
		ignore         string
		accept         string
		name           string
		expectedResult bool
	}{
		{"", "", "eth0", false},
		{"", "^💩0$", "💩0", false},
		{"", "^💩0$", "💩1", true},
		{"", "^💩0$", "veth0", true},
		{"^💩", "", "💩3", true},
		{"^💩", "", "veth0", false},
	}

	for _, test := range tests {
		filter := newNetDevFilter(test.ignore, test.accept)
		result := filter.ignored(test.name)

		if result != test.expectedResult {
			t.Errorf("ignorePattern=%v acceptPattern=%v ifname=%v expected=%v result=%v", test.ignore, test.accept, test.name, test.expectedResult, result)
		}
	}
}
