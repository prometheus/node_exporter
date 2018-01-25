package collector

import (
	"testing"
)

func TestReadProcessStatus(t *testing.T) {
	want := 4
	threads, err := getAllocatedThreads()
	if err != nil {
		t.Fatal(err)
	}
	if threads < want {
		t.Fatalf("Current threads: %d Shouldn't be less than wanted %d",threads,want)
	}
	max, err := getMaxThreads()
	if err != nil {
		t.Fatal(err)
	}
	if max <= 0 {
		t.Fatalf("Maximum allowed amount of threads in the system %d which sould be" +
			"greated than 0",max)
	}

}