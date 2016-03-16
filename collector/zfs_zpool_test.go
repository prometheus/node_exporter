package collector

import (
	"os"
	"testing"
)

func TestZpoolParsing(t *testing.T) {

	zpoolOutput, err := os.Open("fixtures/zfs/zpool_stats_stdout.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer zpoolOutput.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	pools := make([]string, 2)
	troutSize := float64(-1)
	troutDedupratio := float64(-1)
	zrootCapacity := float64(-1)

	err = c.parseZpoolOutput(zpoolOutput, func(pool, name string, value float64) {
		pools = append(pools, pool)
		if pool == "trout" && name == "size" {
			troutSize = value
		}
		if pool == "trout" && name == "dedupratio" {
			troutDedupratio = value
		}
		if pool == "zroot" && name == "capacity" {
			zrootCapacity = value
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if pools[0] == "trout" && pools[1] == "zroot" {
		t.Fatal("Did not parse all pools in fixture")
	}

	if troutSize != float64(4294967296) {
		t.Fatal("Unexpected value for pool 'trout's size value")
	}

	if troutDedupratio != float64(1.0) {
		t.Fatal("Unexpected value for pool 'trout's dedupratio value")
	}

	if zrootCapacity != float64(0.5) {
		t.Fatal("Unexpected value for pool 'zroot's capacity value")
	}

}
