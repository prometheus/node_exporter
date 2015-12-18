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
		"md3":   mdStatus{"md3", true, 8, 8, 5853468288, 5853468288},
		"md127": mdStatus{"md127", true, 2, 2, 312319552, 312319552},
		"md0":   mdStatus{"md0", true, 2, 2, 248896, 248896},
		"md4":   mdStatus{"md4", false, 2, 2, 4883648, 4883648},
		"md6":   mdStatus{"md6", true, 1, 2, 195310144, 16775552},
		"md8":   mdStatus{"md8", true, 2, 2, 195310144, 16775552},
		"md7":   mdStatus{"md7", true, 3, 4, 7813735424, 7813735424},
		"md9":   mdStatus{"md9", true, 4, 4, 523968, 523968},
	}

	for _, md := range mdStates {
		if md != refs[md.mdName] {
			t.Errorf("failed parsing md-device %s correctly: want %v, got %v", md.mdName, refs[md.mdName], md)
		}
	}

	if len(mdStates) != len(refs) {
		t.Errorf("expected number of parsed md-device to be %d, but was %d", len(refs), len(mdStates))
	}
}
