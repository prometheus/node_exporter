package collector

import (
	"testing"
)

func TestBonding(t *testing.T) {
	bondingStats, err := readBondingStats("fixtures/bonding")
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
