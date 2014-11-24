package collector

import (
	"os"
	"testing"
)

func TestNetDevStats(t *testing.T) {
	file, err := os.Open("fixtures/net-dev")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	netStats, err := parseNetDevStats(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "10437182923", netStats["receive"]["wlan0"]["bytes"]; want != got {
		t.Errorf("want netstat wlan0 bytes %s, got %s", want, got)
	}

	if want, got := "934", netStats["transmit"]["tun0"]["packets"]; want != got {
		t.Errorf("want netstat tun0 packets %s, got %s", want, got)
	}
}
