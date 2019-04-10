// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !noprocesses

package collector

import (
	"testing"

	"github.com/prometheus/procfs"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func TestReadProcessStatus(t *testing.T) {
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "fixtures/proc"}); err != nil {
		t.Fatal(err)
	}
	want := 1
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		t.Errorf("failed to open procfs: %v", err)
	}
	c := processCollector{fs: fs}
	pids, states, threads, err := c.getAllocatedThreads()
	if err != nil {
		t.Fatalf("Cannot retrieve data from procfs getAllocatedThreads function: %v ", err)
	}
	if threads < want {
		t.Fatalf("Current threads: %d Shouldn't be less than wanted %d", threads, want)
	}
	if states == nil {

		t.Fatalf("Process states cannot be nil %v:", states)
	}
	maxPid, err := readUintFromFile(procFilePath("sys/kernel/pid_max"))
	if err != nil {
		t.Fatalf("Unable to retrieve limit number of maximum pids alloved %v\n", err)
	}
	if uint64(pids) > maxPid || pids == 0 {
		t.Fatalf("Total running pids cannot be greater than %d or equals to 0", maxPid)
	}
}
