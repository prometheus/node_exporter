package collector

import (
	"testing"
)

func TestReadProcessStatus(t *testing.T) {
	want := 4
	states, threads, err := getAllocatedThreads()
	if err != nil {
		t.Fatalf("Cannot retrieve data from procfs getAllocatedThreads function: %v ", err)
	}
	if threads < want {
		t.Fatalf("Current threads: %d Shouldn't be less than wanted %d", threads, want)
	}
	if states == nil {

		t.Fatalf("Prcess states cannot be nil %v:", states)
	}
}
