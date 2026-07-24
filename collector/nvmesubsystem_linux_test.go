// Copyright The Prometheus Authors
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

//go:build !nonvmesubsystem

package collector

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type testNVMeSubsystemCollector struct {
	mc Collector
}

func (c testNVMeSubsystemCollector) Collect(ch chan<- prometheus.Metric) {
	c.mc.Update(ch)
}

func (c testNVMeSubsystemCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func newTestNVMeSubsystemCollector(logger *slog.Logger) (prometheus.Collector, error) {
	mc, err := NewNVMeSubsystemCollector(logger)
	if err != nil {
		return testNVMeSubsystemCollector{}, err
	}
	return &testNVMeSubsystemCollector{mc}, nil
}

func TestNVMeSubsystemMetrics(t *testing.T) {
	*sysPath = "fixtures/sys"

	testcase := `# HELP node_nvmesubsystem_info Non-numeric information about an NVMe subsystem.
# TYPE node_nvmesubsystem_info gauge
node_nvmesubsystem_info{iopolicy="round-robin",model="Dell PowerStore",nqn="nqn.2014-08.org.nvmexpress:uuid:a34c4f3a-0d6f-5cec-dead-beefcafebabe",serial="SN12345678",subsystem="nvme-subsys0"} 1
# HELP node_nvmesubsystem_namespace_info Maps an NVMe namespace block device to its subsystem.
# TYPE node_nvmesubsystem_namespace_info gauge
node_nvmesubsystem_namespace_info{device="nvme0n1",subsystem="nvme-subsys0"} 1
# HELP node_nvmesubsystem_path_state Current NVMe controller path state (1 for the current state, 0 for all others).
# TYPE node_nvmesubsystem_path_state gauge
node_nvmesubsystem_path_state{controller="nvme0",state="connecting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme0",state="dead",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme0",state="deleting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme0",state="live",subsystem="nvme-subsys0",transport="fc"} 1
node_nvmesubsystem_path_state{controller="nvme0",state="new",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme0",state="resetting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme0",state="unknown",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme1",state="connecting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme1",state="dead",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme1",state="deleting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme1",state="live",subsystem="nvme-subsys0",transport="fc"} 1
node_nvmesubsystem_path_state{controller="nvme1",state="new",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme1",state="resetting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme1",state="unknown",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme2",state="connecting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme2",state="dead",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme2",state="deleting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme2",state="live",subsystem="nvme-subsys0",transport="fc"} 1
node_nvmesubsystem_path_state{controller="nvme2",state="new",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme2",state="resetting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme2",state="unknown",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme3",state="connecting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme3",state="dead",subsystem="nvme-subsys0",transport="fc"} 1
node_nvmesubsystem_path_state{controller="nvme3",state="deleting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme3",state="live",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme3",state="new",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme3",state="resetting",subsystem="nvme-subsys0",transport="fc"} 0
node_nvmesubsystem_path_state{controller="nvme3",state="unknown",subsystem="nvme-subsys0",transport="fc"} 0
# HELP node_nvmesubsystem_paths Number of controller paths for an NVMe subsystem.
# TYPE node_nvmesubsystem_paths gauge
node_nvmesubsystem_paths{subsystem="nvme-subsys0"} 4
# HELP node_nvmesubsystem_paths_live Number of controller paths in live state for an NVMe subsystem.
# TYPE node_nvmesubsystem_paths_live gauge
node_nvmesubsystem_paths_live{subsystem="nvme-subsys0"} 3
`
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
	}))
	c, err := newTestNVMeSubsystemCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	err = testutil.GatherAndCompare(reg, strings.NewReader(testcase))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNVMeSubsystemNoDevices(t *testing.T) {
	*sysPath = t.TempDir()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewNVMeSubsystemCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	c := coll.(*nvmeSubsystemCollector)

	ch := make(chan prometheus.Metric, 200)
	err = c.Update(ch)
	close(ch)

	if err != ErrNoData {
		t.Fatalf("expected ErrNoData, got %v", err)
	}
}

func TestNormalizeControllerState(t *testing.T) {
	tests := []struct {
		raw      string
		expected string
	}{
		{"live", "live"},
		{"connecting", "connecting"},
		{"resetting", "resetting"},
		{"dead", "dead"},
		{"deleting", "deleting"},
		{"deleting (no IO)", "deleting"},
		{"new", "new"},
		{"", "unknown"},
		{"something-else", "unknown"},
	}
	for _, tc := range tests {
		got := normalizeControllerState(tc.raw)
		if got != tc.expected {
			t.Errorf("normalizeControllerState(%q) = %q, want %q", tc.raw, got, tc.expected)
		}
	}
}
