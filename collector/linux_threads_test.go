package collector

import (
	"testing"
)

func TestReadProcessStatus(t *testing.T) {
	want := 1
	states, threads, err := getAllocatedThreads()
	if err != nil {
		t.Fatalf("Cannot retrieve data from procfs getAllocatedThreads function: %v ", err)
	}
	if threads < want {
		t.Fatalf("Current threads: %d Shouldn't be less than wanted %d", threads, want)
	}
	if states == nil {

		t.Fatalf("Process states cannot be nil %v:", states)
	}
}
