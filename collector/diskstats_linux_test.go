package collector

import (
	"os"
	"testing"
)

func TestDiskStats(t *testing.T) {
	file, err := os.Open("fixtures/diskstats")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	diskStats, err := parseDiskStats(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "25353629", diskStats["sda4"][0]; want != got {
		t.Errorf("want diskstats sda4 %s, got %s", want, got)
	}

	if want, got := "68", diskStats["mmcblk0p2"][10]; want != got {
		t.Errorf("want diskstats mmcblk0p2 %s, got %s", want, got)
	}
}
