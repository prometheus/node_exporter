package collector

import (
	"os"
	"testing"
)

func TestFileFDStats(t *testing.T) {
	file, err := os.Open("fixtures/file-nr")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileFDStats, err := parseFileFDStats(file, fileName)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := "1024", fileFDStats["allocated"]; want != got {
		t.Errorf("want filefd allocated %s, got %s", want, got)
	}

	if want, got := "1631329", fileFDStats["maximum"]; want != got {
		t.Errorf("want filefd maximum %s, got %s", want, got)
	}
}
