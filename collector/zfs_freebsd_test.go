package collector

import (
	"os"
	"testing"
)

func TestArcstatsParsing(t *testing.T) {

	arcstatsOutput, err := os.Open("fixtures/sysctl/freebsd/kstat.zfs.misc.arcstats.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer arcstatsOutput.Close()

	c := zfsCollector{}
	if err != nil {
		t.Fatal(err)
	}

	handlerCalled := false
	err = c.parseArcstatsSysctlOutput(arcstatsOutput, func(s zfsSysctl, v zfsMetricValue) {

		if s != zfsSysctl("kstat.zfs.misc.arcstats.hits") {
			return
		}

		handlerCalled = true

		if v != zfsMetricValue(63068289) {
			t.Fatalf("Incorrect value parsed from sysctl output")
		}

	})

	if err != nil {
		t.Fatal(err)
	}

	if !handlerCalled {
		t.Fatal("Arcstats parsing handler was not called for some expected sysctls")
	}

}
