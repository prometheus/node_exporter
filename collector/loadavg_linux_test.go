package collector

import "testing"

func TestLoad(t *testing.T) {
	load, err := parseLoad("0.21 0.37 0.39 1/719 19737")
	if err != nil {
		t.Fatal(err)
	}

	if want := 0.21; want != load {
		t.Fatalf("want load %f, got %f", want, load)
	}
}
