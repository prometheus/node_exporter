package collector

import (
	"os"
	"testing"
)

func TestInterrupts(t *testing.T) {
	file, err := os.Open("fixtures/interrupts")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	interrupts, err := parseInterrupts(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "5031", interrupts["NMI"].values[1]; want != got {
		t.Errorf("want interrupts %s, got %s", want, got)
	}
}
