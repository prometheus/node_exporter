// +build !nomegacli

package collector

import (
	"flag"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	testMegaCliAdapter = "fixtures/megacli_adapter.txt"
	testMegaCliDisks   = "fixtures/megacli_disks.txt"

	physicalDevicesExpected = "5"
	virtualDevicesDegraded  = "0"
)

func TestMegaCliAdapter(t *testing.T) {
	data, err := os.Open(testMegaCliAdapter)
	if err != nil {
		t.Fatal(err)
	}
	stats, err := parseMegaCliAdapter(data)
	if err != nil {
		t.Fatal(err)
	}

	if stats["Device Present"]["Physical Devices"] != physicalDevicesExpected {
		t.Fatalf("Unexpected device count: %d != %d", stats["Device Present"]["Physical Devices"], physicalDevicesExpected)
	}

	if stats["Device Present"]["Degraded"] != virtualDevicesDegraded {
		t.Fatal()
	}
}

func TestMegaCliDisks(t *testing.T) {
	data, err := os.Open(testMegaCliDisks)
	if err != nil {
		t.Fatal(err)
	}
	stats, err := parseMegaCliDisks(data)
	if err != nil {
		t.Fatal(err)
	}

	if stats[32][0]["Drive Temperature"] != "37C (98.60 F)" {
		t.Fatalf("Unexpected drive temperature: %s", stats[32][0]["Drive Temperature"])
	}

	if stats[32][1]["Drive Temperature"] != "N/A" {
		t.Fatalf("Unexpected drive temperature: %s", stats[32][2]["Drive Temperature"])
	}

	if stats[32][3]["Predictive Failure Count"] != "23" {
		t.Fatal()
	}
}

func TestMegaCliCollectorDoesntCrash(t *testing.T) {
	if err := flag.Set("collector.megacli.command", "./fixtures/megacli"); err != nil {
		t.Fatal(err)
	}
	collector, err := NewMegaCliCollector()
	if err != nil {
		t.Fatal(err)
	}
	sink := make(chan prometheus.Metric)
	go func() {
		for {
			<-sink
		}
	}()

	err = collector.Update(sink)
	if err != nil {
		t.Fatal(err)
	}
}
