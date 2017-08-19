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
		// State bit:
		//  MD_SB_CLEAN          0 - value 1
		//  MD_SB_ERRORS         1 - value 2
		//  MD_SB_CLUSTERED      5 - value 32  /* MD is clustered */
		//  MD_SB_BITMAP_PRESENT 8 - value 256 /* bitmap may be present nearby */

		// { "<name>", state, disksActive, disksFailed, disksMissing, disksSpare, bytesTotal, bytesSynced }
		// TODO: What should disks[A-Z]* be?
		// md6: raid1 recovery for 1 disk: 1, 0, 1, 0?
		"md3":   {"md3", 1, 8, 0, 0, 0, 5853468288, 5853468288},
		"md127": {"md127", 1, 2, 0, 0, 0, 312319552, 312319552},
		"md0":   {"md0", 1, 2, 0, 0, 0, 248896, 248896},
		"md4":   {"md4", 2, 2, 0, 0, 0, 4883648, 4883648},
		"md6":   {"md6", 1, 1, 0, 1, 0, 16775552, 195310144},
		"md8":   {"md8", 1, 2, 0, 1, 1, 16775552, 195310144},
		"md7":   {"md7", 1, 3, 1, 0, 0, 7813735424, 7813735424},
		"md9":   {"md9", 1, 4, 0, 0, 0, 523968, 523968},
		"md10":  {"md10", 1, 2, 0, 0, 0, 314159265, 314159265},
		"md11":  {"md11", 1, 2, 0, 0, 0, 4190208, 4190208},
		"md12":  {"md12", 1, 2, 0, 0, 0, 3886394368, 3886394368},
		"md126": {"md126", 1, 2, 0, 0, 0, 1855870976, 1855870976},
		"md219": {"md219", 2, 0, 0, 0, 3, 7932, 7932},
		"md00":  {"md00", 1, 1, 0, 0, 0, 4186624, 4186624},
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
