package collector

import (
	"os"
	"testing"
)

func TestTCPStat(t *testing.T) {
	file, err := os.Open("fixtures/tcpstat")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	tcpStats, err := parseTCPStats(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 1, int(tcpStats[TCP_ESTABLISHED]); want != got {
		t.Errorf("want tcpstat number of established state %d, got %d", want, got)
	}

	if want, got := 1, int(tcpStats[TCP_LISTEN]); want != got {
		t.Errorf("want tcpstat number of listen state %d, got %d", want, got)
	}
}
