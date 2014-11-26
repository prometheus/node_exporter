package collector

import (
	"bytes"
	"testing"
)

func TestConsulAgentInput(t *testing.T) {
	var (
		output = `
agent:
        check_monitors = 0
        check_ttls = 0
        checks = 0
        services = 8
build:
        prerelease =
        revision = 461c1e18
        version = 0.4.2.soundcloud4
consul:
        known_servers = 3
        server = false
runtime:
        arch = amd64
        cpu_count = 1
        goroutines = 36
        max_procs = 16
        os = linux
        version = go1.3
serf_lan:
        event_queue = 0
        event_time = 55
        failed = 0
        intent_queue = 0
        left = 0
        member_time = 116
        members = 11
        query_queue = 0
        query_time = 1
`

		r = bytes.NewBufferString(output)
	)

	stats, err := parseConsulStats(r)
	if err != nil {
		t.Fatalf("parse failed: %s", err)
	}

	if want, got := 18, len(stats); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestConsulServerInput(t *testing.T) {
	var (
		output = `
agent:
        check_monitors = 0
        check_ttls = 0
        checks = 0
        services = 4
build:
        prerelease =
        revision = 461c1e18
        version = 0.4.2.soundcloud4
consul:
        bootstrap = false
        known_datacenters = 1
        leader = true
        server = true
raft:
        applied_index = 1712470
        commit_index = 1712470
        fsm_pending = 0
        last_contact = 0
        last_log_index = 1712470
        last_log_term = 131
        last_snapshot_index = 1708783
        last_snapshot_term = 131
        num_peers = 2
        state = Leader
        term = 131
runtime:
        arch = amd64
        cpu_count = 1
        goroutines = 81
        max_procs = 16
        os = linux
        version = go1.3
serf_lan:
        event_queue = 0
        event_time = 55
        failed = 0
        intent_queue = 0
        left = 0
        member_time = 116
        members = 11
        query_queue = 0
        query_time = 1
serf_wan:
        event_queue = 0
        event_time = 1
        failed = 0
        intent_queue = 0
        left = 0
        member_time = 1
        members = 1
        query_queue = 0
        query_time = 1
`
		r = bytes.NewBufferString(output)
	)

	stats, err := parseConsulStats(r)
	if err != nil {
		t.Fatalf("parse failed: %s", err)
	}

	if want, got := 39, len(stats); want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
