package collector

import (
	"testing"
)

func TestMdadm(t *testing.T) {
	mdStates, err := parseMdstat("fixtures/mdstat")

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
