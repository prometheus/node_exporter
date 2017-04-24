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
)

func TestBonding(t *testing.T) {
	bondingStats, err := readBondingStats("fixtures/sys/class/net")
	if err != nil {
		t.Fatal(err)
	}
	if bondingStats["bond0"][0] != 0 || bondingStats["bond0"][1] != 0 {
		t.Fatal("bond0 in unexpected state")
	}

	if bondingStats["int"][0] != 2 || bondingStats["int"][1] != 1 {
		t.Fatal("int in unexpected state")
	}

	if bondingStats["dmz"][0] != 2 || bondingStats["dmz"][1] != 2 {
		t.Fatal("dmz in unexpected state")
	}
}
