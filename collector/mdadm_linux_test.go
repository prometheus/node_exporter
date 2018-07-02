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

func TestMdadm(t *testing.T) {
	mdStates, err := parseMdstat("fixtures/proc/mdstat")
	if err != nil {
		t.Fatalf("parsing of reference-file failed entirely: %s", err)
	}

	refs := map[string]mdStatus{
		// { "<name>", <active?>, <numDisksActive>, <totalNumDisks>, <amountSynced>, <totalSize>}
		"md3":   {"md3", true, 8, 8, 5853468288, 5853468288},
		"md127": {"md127", true, 2, 2, 312319552, 312319552},
		"md0":   {"md0", true, 2, 2, 248896, 248896},
		"md4":   {"md4", false, 0, 2, 4883648, 4883648},
		"md6":   {"md6", true, 1, 2, 195310144, 16775552},
		"md8":   {"md8", true, 2, 2, 195310144, 16775552},
		"md7":   {"md7", true, 3, 4, 7813735424, 7813735424},
		"md9":   {"md9", true, 4, 4, 523968, 523968},
		"md10":  {"md10", true, 2, 2, 314159265, 314159265},
		"md11":  {"md11", true, 2, 2, 4190208, 4190208},
		"md12":  {"md12", true, 2, 2, 3886394368, 3886394368},
		"md120": {"md120", true, 2, 2, 2095104, 2095104},
		"md126": {"md126", true, 2, 2, 1855870976, 1855870976},
		"md219": {"md219", false, 0, 3, 7932, 7932},
		"md00":  {"md00", true, 1, 1, 4186624, 4186624},
	}

	for _, md := range mdStates {
		if md != refs[md.name] {
			t.Errorf("failed parsing md-device %s correctly: want %v, got %v", md.name, refs[md.name], md)
		}
	}

	if len(mdStates) != len(refs) {
		t.Errorf("expected number of parsed md-device to be %d, but was %d", len(refs), len(mdStates))
	}
}

func TestInvalidMdstat(t *testing.T) {
	_, err := parseMdstat("fixtures/proc/mdstat_invalid")
	if err == nil {
		t.Fatalf("parsing of invalid reference file did not find any errors")
	}
}
