package collector

import (
	"os"
	"testing"
)

var fileName = "fixtures/netstat"

func TestNetStats(t *testing.T) {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	netStats, err := parseNetStats(file, fileName)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "102471", netStats["TcpExt"]["DelayedACKs"]; want != got {
		t.Errorf("want netstat TCP DelayedACKs %s, got %s", want, got)
	}

	if want, got := "2786264347", netStats["IpExt"]["OutOctets"]; want != got {
		t.Errorf("want netstat IP OutOctets %s, got %s", want, got)
	}
}
