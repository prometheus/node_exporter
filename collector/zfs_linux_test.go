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

	p := NewZFSMetricProvider()
	err = p.parseArcstatsProcfsFile(arcstatsFile)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Parsed values mapped to sysctls:\n%v", p.values)

	value, err := p.Value(zfsSysctl("kstat.zfs.misc.arcstats.hits"))

	if err != nil {
		t.Fatal(err)
	}

	if value != zfsMetricValue(8772612) {
		t.Fatalf("Incorrect value parsed from procfs data")
	}

}
