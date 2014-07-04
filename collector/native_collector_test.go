package collector

import (
	"io/ioutil"
	"os"
	"testing"
)

const (
	loadExpected = 0.21

	memTotalExpected       = 3831959552
	memDirectMap2MExpected = 3787456512

	interruptsNmi1Expected = "5031"

	netReceiveWlan0Bytes    = "10437182923"
	netTransmitTun0Packages = "934"

	diskSda4ReadsCompleted = "25353629"
	diskMmcIoTimeWeighted  = "68"

	testProcLoad       = "fixtures/loadavg"
	testProcMemInfo    = "fixtures/meminfo"
	testProcInterrupts = "fixtures/interrupts"
	testProcNetDev     = "fixtures/net-dev"
	testProcDiskStats  = "fixtures/diskstats"
)

func TestLoad(t *testing.T) {
	data, err := ioutil.ReadFile(testProcLoad)
	if err != nil {
		t.Fatal(err)
	}
	load, err := parseLoad(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if load != loadExpected {
		t.Fatalf("Unexpected load: %f != %f", load, loadExpected)
	}
}

func TestMemInfo(t *testing.T) {
	file, err := os.Open(testProcMemInfo)
	if err != nil {
		t.Fatal(err)
	}

	memInfo, err := parseMemInfo(file)
	if err != nil {
		t.Fatal(err)
	}
	if memInfo["MemTotal"] != memTotalExpected {
		t.Fatalf("Unexpected memory: %s != %s", memInfo["MemTotal"], memTotalExpected)
	}
	if memInfo["DirectMap2M"] != memDirectMap2MExpected {
		t.Fatalf("Unexpected memory: %s != %s", memInfo["MemTotal"], memTotalExpected)
	}

}

func TestInterrupts(t *testing.T) {
	file, err := os.Open(testProcInterrupts)
	if err != nil {
		t.Fatal(err)
	}

	interrupts, err := parseInterrupts(file)
	if err != nil {
		t.Fatal(err)
	}
	if interrupts["NMI"].values[1] != interruptsNmi1Expected {
		t.Fatalf("Unexpected interrupts: %s != %s", interrupts["NMI"].values[1],
			interruptsNmi1Expected)
	}

}

func TestNetStats(t *testing.T) {
	file, err := os.Open(testProcNetDev)
	if err != nil {
		t.Fatal(err)
	}
	netStats, err := parseNetStats(file)
	if err != nil {
		t.Fatal(err)
	}
	if netStats["receive"]["wlan0"]["bytes"] != netReceiveWlan0Bytes {
		t.Fatalf("Unexpected netstats: %s != %s", netStats["receive"]["wlan0"]["bytes"],
			netReceiveWlan0Bytes)
	}
	if netStats["transmit"]["tun0"]["packets"] != netTransmitTun0Packages {
		t.Fatalf("Unexpected netstats: %s != %s", netStats["transmit"]["tun0"]["packets"],
			netTransmitTun0Packages)
	}
}

func TestDiskStats(t *testing.T) {
	file, err := os.Open(testProcDiskStats)
	if err != nil {
		t.Fatal(err)
	}
	diskStats, err := parseDiskStats(file)
	if err != nil {
		t.Fatal(err)
	}
	if diskStats["sda4"][0] != diskSda4ReadsCompleted {
		t.Fatalf("Unexpected diskstats: %s != %s", diskStats["sda4"][0],
			diskSda4ReadsCompleted)
	}

	if diskStats["mmcblk0p2"][10] != diskMmcIoTimeWeighted {
		t.Fatalf("Unexpected diskstats: %s != %s",
			diskStats["mmcblk0p2"][10], diskMmcIoTimeWeighted)
	}
}
