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

//go:build !nobonding
// +build !nobonding

package collector

import (
	"testing"
)

func TestBonding(t *testing.T) {
	bondingStats, err := readBondingStats("fixtures/sys/class/net")
	if err != nil {
		t.Fatal(err)
	}
	if bondingStats[0].name != "bond0" || bondingStats[0].slaves != 0 || bondingStats[0].active != 0 || bondingStats[0].miimon != 100 {
		t.Fatal("bond0 in unexpected state")
	}

	if bondingStats[1].name != "dmz" || bondingStats[1].slaves != 2 || bondingStats[1].active != 2 || bondingStats[1].miimon != 0 {
		t.Fatal("dmz in unexpected state")
	}

	if bondingStats[2].name != "int" || bondingStats[2].slaves != 2 || bondingStats[2].active != 1 || bondingStats[2].miimon != 200 {
		t.Fatal("int in unexpected state")
	}
}
