package collector

import (
	"os"
	"testing"
)

func TestArcstatsParsing(t *testing.T) {

	arcstatsFile, err := os.Open("fixtures/proc/spl/kstat/zfs/arcstats")
	if err != nil {
		t.Fatal(err)
	}
	defer arcstatsFile.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseArcstatsProcfsFile(arcstatsFile, func(s zfsSysctl, v zfsMetricValue) {

		if s != zfsSysctl("kstat.zfs.misc.arcstats.hits") {
			return
		}

		handlerCalled = true

		if v != zfsMetricValue(8772612) {
			t.Fatalf("Incorrect value parsed from procfs data")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("Arcstats parsing handler was not called for some expected sysctls")
	}

}
