// Copyright 2026 The Prometheus Authors
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

//go:build !nondisc

package collector

import (
	"net"
	"testing"

	"github.com/jsimonetti/rtnetlink/v2/rtnl"
	"golang.org/x/sys/unix"
)

func TestGetTotalNdiscEntries(t *testing.T) {
	t.Parallel()

	neighbors := []*rtnl.Neigh{
		{Interface: &net.Interface{Name: "eth0"}, State: unix.NUD_REACHABLE},
		{Interface: &net.Interface{Name: "eth0"}, State: unix.NUD_STALE},
		{Interface: &net.Interface{Name: "eth1"}, State: unix.NUD_DELAY},
		{Interface: &net.Interface{Name: "eth1"}, State: unix.NUD_NOARP},
	}

	entries := getTotalNdiscEntries(neighbors)

	if got, want := entries["eth0"], uint32(2); got != want {
		t.Fatalf("unexpected entry count for eth0: got %d, want %d", got, want)
	}
	if got, want := entries["eth1"], uint32(1); got != want {
		t.Fatalf("unexpected entry count for eth1: got %d, want %d", got, want)
	}
}
