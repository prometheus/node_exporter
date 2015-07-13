package collector

import (
	"os"
	"regexp"
	"testing"
)

func TestNetDevStats(t *testing.T) {
	file, err := os.Open("fixtures/net-dev")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	netStats, err := parseNetDevStats(file, regexp.MustCompile("^veth"))
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "10437182923", netStats["receive"]["wlan0"]["bytes"]; want != got {
		t.Errorf("want netstat wlan0 bytes %s, got %s", want, got)
	}

	if want, got := "68210035552", netStats["receive"]["eth0"]["bytes"]; want != got {
		t.Errorf("want netstat eth0 bytes %s, got %s", want, got)
	}

	if want, got := "934", netStats["transmit"]["tun0"]["packets"]; want != got {
		t.Errorf("want netstat tun0 packets %s, got %s", want, got)
	}

	if want, got := 6, len(netStats["transmit"]); want != got {
		t.Errorf("want count of devices to be %d, got %d", want, got)
	}

	if _, ok := netStats["transmit"]["veth4B09XN"]; ok != false {
		t.Error("want fixture interface veth4B09XN to not exist, but it does")
	}
}
