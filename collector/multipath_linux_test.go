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

//go:build !nomultipath

package collector

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestMultipathJSONParsing(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("fixtures", "multipathd", "show_maps_json.json"))
	if err != nil {
		t.Fatal(err)
	}

	var topo multipathTopology
	if err := json.Unmarshal(data, &topo); err != nil {
		t.Fatal(err)
	}

	if topo.MajorVersion != 0 {
		t.Errorf("expected major_version 0, got %d", topo.MajorVersion)
	}
	if len(topo.Maps) != 2 {
		t.Fatalf("expected 2 maps, got %d", len(topo.Maps))
	}

	m := topo.Maps[0]
	if m.Name != "mpathA" {
		t.Errorf("expected name mpathA, got %s", m.Name)
	}
	if m.Paths != 4 {
		t.Errorf("expected 4 paths, got %d", m.Paths)
	}
	if m.PathFaults != 2 {
		t.Errorf("expected 2 path_faults, got %d", m.PathFaults)
	}
	if len(m.PathGroups) != 2 {
		t.Fatalf("expected 2 path groups, got %d", len(m.PathGroups))
	}
	if len(m.PathGroups[0].Paths) != 2 {
		t.Errorf("expected 2 paths in PG1, got %d", len(m.PathGroups[0].Paths))
	}
	if m.PathGroups[1].Paths[1].ChkSt != "faulty" {
		t.Errorf("expected faulty chk_st for sdd, got %s", m.PathGroups[1].Paths[1].ChkSt)
	}
}

func TestNormalizeCheckerState(t *testing.T) {
	tests := []struct {
		raw      string
		expected string
	}{
		{"ready", "ready"},
		{"faulty", "faulty"},
		{"ghost", "ghost"},
		{"shaky", "shaky"},
		{"delayed", "delayed"},
		{"disconnected", "disconnected"},
		{"i/o pending", "pending"},
		{"i/o timeout", "timeout"},
		{"undef", "unknown"},
		{"", "unknown"},
		{"something-else", "unknown"},
	}
	for _, tc := range tests {
		got := normalizeCheckerState(tc.raw)
		if got != tc.expected {
			t.Errorf("normalizeCheckerState(%q) = %q, want %q", tc.raw, got, tc.expected)
		}
	}
}

func TestMultipathCollectorUpdate(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("fixtures", "multipathd", "show_maps_json.json"))
	if err != nil {
		t.Fatal(err)
	}

	var topo multipathTopology
	if err := json.Unmarshal(data, &topo); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewMultipathCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	mc := coll.(*multipathCollector)
	mc.queryTopology = func() (*multipathTopology, error) {
		return &topo, nil
	}
	mc.readDeviceSize = func(sysfsName string) (uint64, error) {
		sizes := map[string]uint64{
			"dm-0": 53687091200,  // 50 GiB
			"dm-1": 107374182400, // 100 GiB
		}
		if s, ok := sizes[sysfsName]; ok {
			return s, nil
		}
		return 0, fmt.Errorf("device %s not found", sysfsName)
	}
	mc.scanNVMeSubsystems = func() ([]nvmeSubsystem, error) {
		return nil, nil
	}

	ch := make(chan prometheus.Metric, 200)
	if err := mc.Update(ch); err != nil {
		t.Fatal(err)
	}
	close(ch)

	metrics := make(map[string][]*dto.Metric)
	for m := range ch {
		d := &dto.Metric{}
		if err := m.Write(d); err != nil {
			t.Fatal(err)
		}
		desc := m.Desc().String()
		metrics[desc] = append(metrics[desc], d)
	}

	assertGaugeValue(t, metrics, "daemon_up", nil, 1)
	assertGaugeValue(t, metrics, "device_active", labelMap{"device": "mpathA"}, 1)
	assertGaugeValue(t, metrics, "device_active", labelMap{"device": "mpathB"}, 1)
	assertGaugeValue(t, metrics, "device_size_bytes", labelMap{"device": "mpathA"}, 53687091200)
	assertGaugeValue(t, metrics, "device_size_bytes", labelMap{"device": "mpathB"}, 107374182400)
	assertGaugeValue(t, metrics, "device_paths_total", labelMap{"device": "mpathA"}, 4)
	assertGaugeValue(t, metrics, "device_paths_total", labelMap{"device": "mpathB"}, 2)

	// mpathA: sda(active+ready) + sdb(active+ready) = 2 active paths.
	// sdc is active but ghost (not "ready"), sdd is failed+faulty.
	assertGaugeValue(t, metrics, "device_paths_active", labelMap{"device": "mpathA"}, 2)
	// mpathA: sdd has dm_st=failed and chk_st=faulty → 1 failed path.
	assertGaugeValue(t, metrics, "device_paths_failed", labelMap{"device": "mpathA"}, 1)
	// mpathB: all paths active+ready.
	assertGaugeValue(t, metrics, "device_paths_active", labelMap{"device": "mpathB"}, 2)
	assertGaugeValue(t, metrics, "device_paths_failed", labelMap{"device": "mpathB"}, 0)

	assertCounterValue(t, metrics, "device_path_faults_total", labelMap{"device": "mpathA"}, 2)
	assertCounterValue(t, metrics, "device_path_faults_total", labelMap{"device": "mpathB"}, 0)

	assertGaugeValue(t, metrics, "path_active", labelMap{"device": "mpathA", "path": "sda", "path_group": "1"}, 1)
	assertGaugeValue(t, metrics, "path_active", labelMap{"device": "mpathA", "path": "sdd", "path_group": "2"}, 0)

	assertGaugeValue(t, metrics, "path_checker_state",
		labelMap{"device": "mpathA", "path": "sda", "path_group": "1", "state": "ready"}, 1)
	assertGaugeValue(t, metrics, "path_checker_state",
		labelMap{"device": "mpathA", "path": "sda", "path_group": "1", "state": "faulty"}, 0)
	assertGaugeValue(t, metrics, "path_checker_state",
		labelMap{"device": "mpathA", "path": "sdd", "path_group": "2", "state": "faulty"}, 1)
	assertGaugeValue(t, metrics, "path_checker_state",
		labelMap{"device": "mpathA", "path": "sdd", "path_group": "2", "state": "ready"}, 0)
	assertGaugeValue(t, metrics, "path_checker_state",
		labelMap{"device": "mpathA", "path": "sdc", "path_group": "2", "state": "ghost"}, 1)
}

func TestMultipathCollectorDaemonDown(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	coll, err := NewMultipathCollector(logger)
	if err != nil {
		t.Fatal(err)
	}

	mc := coll.(*multipathCollector)
	mc.queryTopology = func() (*multipathTopology, error) {
		return nil, fmt.Errorf("connection refused")
	}
	mc.scanNVMeSubsystems = func() ([]nvmeSubsystem, error) {
		return nil, fmt.Errorf("no nvme-subsystem")
	}

	ch := make(chan prometheus.Metric, 200)
	updateErr := mc.Update(ch)
	close(ch)

	if updateErr == nil {
		t.Fatal("expected error from Update when daemon is down and no NVMe")
	}

	for m := range ch {
		d := &dto.Metric{}
		if err := m.Write(d); err != nil {
			t.Fatal(err)
		}
		if d.GetGauge() != nil && d.GetGauge().GetValue() == 0 {
			return
		}
	}
	t.Error("expected daemon_up=0 metric")
}

func TestMultipathSocketProtocol(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("fixtures", "multipathd", "show_maps_json.json"))
	if err != nil {
		t.Fatal(err)
	}

	socketDir := t.TempDir()
	socketPath := filepath.Join(socketDir, "test.sock")

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		lenBuf := make([]byte, sizeOfSizeT)
		if _, err := io.ReadFull(conn, lenBuf); err != nil {
			return
		}
		var cmdLen uint64
		switch sizeOfSizeT {
		case 8:
			cmdLen = binary.NativeEndian.Uint64(lenBuf)
		case 4:
			cmdLen = uint64(binary.NativeEndian.Uint32(lenBuf))
		}
		cmdBuf := make([]byte, cmdLen)
		if _, err := io.ReadFull(conn, cmdBuf); err != nil {
			return
		}

		reply := append(data, 0)
		replyLen := uint64(len(reply))
		switch sizeOfSizeT {
		case 8:
			binary.NativeEndian.PutUint64(lenBuf, replyLen)
		case 4:
			binary.NativeEndian.PutUint32(lenBuf, uint32(replyLen))
		}
		conn.Write(lenBuf)
		conn.Write(reply)
	}()

	topo, err := queryMultipathd(socketPath, 5e9)
	if err != nil {
		t.Fatal(err)
	}

	if len(topo.Maps) != 2 {
		t.Errorf("expected 2 maps, got %d", len(topo.Maps))
	}
	if topo.Maps[0].Name != "mpathA" {
		t.Errorf("expected mpathA, got %s", topo.Maps[0].Name)
	}
}

type labelMap map[string]string

func assertGaugeValue(t *testing.T, metrics map[string][]*dto.Metric, metricSubstring string, labels labelMap, expected float64) {
	t.Helper()
	for desc, ms := range metrics {
		if !strings.Contains(desc, metricSubstring) {
			continue
		}
		for _, m := range ms {
			if matchLabels(m.GetLabel(), labels) {
				got := m.GetGauge().GetValue()
				if got != expected {
					t.Errorf("%s%v: got %v, want %v", metricSubstring, labels, got, expected)
				}
				return
			}
		}
	}
	t.Errorf("metric %s%v not found", metricSubstring, labels)
}

func assertCounterValue(t *testing.T, metrics map[string][]*dto.Metric, metricSubstring string, labels labelMap, expected float64) {
	t.Helper()
	for desc, ms := range metrics {
		if !strings.Contains(desc, metricSubstring) {
			continue
		}
		for _, m := range ms {
			if matchLabels(m.GetLabel(), labels) {
				got := m.GetCounter().GetValue()
				if got != expected {
					t.Errorf("%s%v: got %v, want %v", metricSubstring, labels, got, expected)
				}
				return
			}
		}
	}
	t.Errorf("metric %s%v not found", metricSubstring, labels)
}

func matchLabels(pairs []*dto.LabelPair, want labelMap) bool {
	if want == nil {
		return len(pairs) == 0
	}
	found := 0
	for _, lp := range pairs {
		if v, ok := want[lp.GetName()]; ok && v == lp.GetValue() {
			found++
		}
	}
	return found == len(want)
}
