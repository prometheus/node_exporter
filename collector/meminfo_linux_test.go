package collector

import (
	"os"
	"testing"
)

func TestMemInfo(t *testing.T) {
	file, err := os.Open("fixtures/meminfo")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	memInfo, err := parseMemInfo(file)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := 3831959552.0, memInfo["MemTotal"]; want != got {
		t.Errorf("want memory total %f, got %f", want, got)
	}

	if want, got := 3787456512.0, memInfo["DirectMap2M"]; want != got {
		t.Errorf("want memory directMap2M %f, got %f", want, got)
	}
}
