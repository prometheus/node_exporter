package collector

import (
	"os"
	"strconv"
	"testing"
)

func TestSockStats(t *testing.T) {
	file, err := os.Open("fixtures/sockstat")
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()

	sockStats, err := parseSockStats(file, fileName)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "229", sockStats["sockets"]["used"]; want != got {
		t.Errorf("want sockstat sockets used %s, got %s", want, got)
	}

	if want, got := "4", sockStats["TCP"]["tw"]; want != got {
		t.Errorf("want sockstat TCP tw %s, got %s", want, got)
	}

	if want, got := "17", sockStats["TCP"]["alloc"]; want != got {
		t.Errorf("want sockstat TCP alloc %s, got %s", want, got)
	}

	// The test file has 1 for TCP mem, which is one page.  So we should get the
	// page size in bytes back from sockstat_linux.  We get the page size from
	// os here because this value can change from system to system.  The value is
	// 4096 by default from linux 2.4 onward.
	if want, got := strconv.Itoa(os.Getpagesize()), sockStats["TCP"]["mem_bytes"]; want != got {
		t.Errorf("want sockstat TCP mem_bytes %s, got %s", want, got)
	}
}
