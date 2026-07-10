// Copyright 2024 The Prometheus Authors
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

//go:build linux && !noebsnvme
// +build linux,!noebsnvme

package collector

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

func TestMapDevices(t *testing.T) {
	// lsblk -nd --json -o NAME,SERIAL,MOUNTPOINT style output, including:
	//   - a mounted EBS volume (serial without dash)
	//   - an unmounted EBS volume (mountpoint null) -> NotMounted
	//   - an empty-string mountpoint -> NotMounted
	//   - a non-EBS device (serial does not start with "vol")
	in := []byte(`{
		"blockdevices": [
			{"name": "nvme1n1", "serial": "vol0e9d0ec0c3c24e723", "mountpoint": "/mysql/undo"},
			{"name": "nvme0n1", "serial": "vol08a138e8d1cfb88df", "mountpoint": null},
			{"name": "nvme2n1", "serial": "vol0cd972338c0a5f21b", "mountpoint": ""},
			{"name": "nvme9n1", "serial": "AWS-LOCAL-1234", "mountpoint": "/scratch"}
		]
	}`)

	got, err := mapDevices(in)
	if err != nil {
		t.Fatalf("mapDevices returned error: %v", err)
	}

	want := map[string]ebsDeviceInfo{
		"/dev/nvme1n1": {volumeID: "vol-0e9d0ec0c3c24e723", mountPath: "/mysql/undo"},
		"/dev/nvme0n1": {volumeID: "vol-08a138e8d1cfb88df", mountPath: notMounted},
		"/dev/nvme2n1": {volumeID: "vol-0cd972338c0a5f21b", mountPath: notMounted},
		"/dev/nvme9n1": {volumeID: "AWS-LOCAL-1234", mountPath: "/scratch"},
	}

	if len(got) != len(want) {
		t.Fatalf("got %d devices, want %d", len(got), len(want))
	}
	for path, w := range want {
		g, ok := got[path]
		if !ok {
			t.Errorf("missing device %s", path)
			continue
		}
		if g != w {
			t.Errorf("device %s: got %+v, want %+v", path, g, w)
		}
	}
}

func TestMapDevicesInvalidJSON(t *testing.T) {
	if _, err := mapDevices([]byte("not json")); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestParseEBSLogPageInvalidMagic(t *testing.T) {
	data := make([]byte, nvmeLogPageSize)
	binary.LittleEndian.PutUint64(data[0:8], 0xDEADBEEF)

	_, err := parseEBSLogPage(data)
	if !errors.Is(err, errInvalidEBSMagic) {
		t.Fatalf("expected errInvalidEBSMagic, got %v", err)
	}
}

func TestParseEBSLogPageValid(t *testing.T) {
	var m ebsMetrics
	m.EBSMagic = ebsMagic
	m.ReadOps = 100
	m.WriteOps = 200
	m.ReadBytes = 4096 * 100
	m.WriteBytes = 16384 * 200
	m.TotalReadTime = 5_000_000  // 5s in microseconds
	m.TotalWriteTime = 2_000_000 // 2s in microseconds
	m.QueueLength = 3

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, &m); err != nil {
		t.Fatalf("failed to encode test log page: %v", err)
	}
	// The EBS metrics struct is smaller than the full NVMe log page; the device
	// returns nvmeLogPageSize bytes with the trailing space unused. Pad the test
	// buffer to mirror that layout.
	page := make([]byte, nvmeLogPageSize)
	copy(page, buf.Bytes())

	got, err := parseEBSLogPage(page)
	if err != nil {
		t.Fatalf("parseEBSLogPage returned error: %v", err)
	}
	if got.ReadOps != 100 || got.WriteOps != 200 || got.QueueLength != 3 {
		t.Errorf("unexpected parsed metrics: %+v", got)
	}
}

func TestConvertEBSHistogram(t *testing.T) {
	var h ebsHistogram
	h.BinCount = 3
	// Upper bounds in microseconds; counts are per-bin (non-cumulative).
	h.Bins[0] = ebsHistogramBin{Lower: 0, Upper: 100, Count: 10}
	h.Bins[1] = ebsHistogramBin{Lower: 100, Upper: 500, Count: 5}
	h.Bins[2] = ebsHistogramBin{Lower: 500, Upper: 1000, Count: 2}

	count, buckets := convertEBSHistogram(h)

	if count != 17 {
		t.Errorf("total count = %d, want 17", count)
	}
	// Buckets are cumulative, keyed by upper bound converted to seconds.
	wantCumulative := map[float64]uint64{
		100.0 / microsecondsInSeconds:  10,
		500.0 / microsecondsInSeconds:  15,
		1000.0 / microsecondsInSeconds: 17,
	}
	for ub, wc := range wantCumulative {
		if buckets[ub] != wc {
			t.Errorf("bucket le=%g: got %d, want %d", ub, buckets[ub], wc)
		}
	}
}
