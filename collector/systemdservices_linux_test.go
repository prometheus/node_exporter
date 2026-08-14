// Copyright 2015 The Prometheus Authors
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

package collector

import (
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewSystemdCollectorsSkipDial(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if _, err := NewSystemdServicesCollector(logger); err != nil {
		t.Fatalf("NewSystemdServicesCollector: %v", err)
	}
	if _, err := NewSystemdInfoCollector(logger); err != nil {
		t.Fatalf("NewSystemdInfoCollector: %v", err)
	}
}

func TestSystemdConnGetReusesHealthy(t *testing.T) {
	dummy := &dbus.Conn{}
	s := &systemdConn{conn: dummy}
	got, err := s.get()
	if err != nil {
		t.Fatal(err)
	}
	if got != dummy {
		t.Fatal("expected stored connection to be reused")
	}
}

func TestSystemdConnGetAfterClose(t *testing.T) {
	s := &systemdConn{}
	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if _, err := s.get(); !errors.Is(err, errSystemdConnClosed) {
		t.Fatalf("get after close: got %v want %v", err, errSystemdConnClosed)
	}
}

func TestListSystemdUnitsNilConn(t *testing.T) {
	if _, err := listSystemdUnits(nil); !errors.Is(err, errSystemdConnClosed) {
		t.Fatalf("listSystemdUnits(nil): got %v want %v", err, errSystemdConnClosed)
	}
}

func TestSystemdCollectorsUpdateClosed(t *testing.T) {
	ch := make(chan prometheus.Metric, 1)
	svc := &systemdServicesCollector{
		serviceState: prometheus.NewDesc("node_systemd_service_state", "x", []string{"name"}, nil),
	}
	if err := svc.Close(); err != nil {
		t.Fatalf("services Close: %v", err)
	}
	if err := svc.Update(ch); !errors.Is(err, errSystemdConnClosed) {
		t.Fatalf("services Update after Close: got %v want %v", err, errSystemdConnClosed)
	}

	info := &systemdInfoCollector{}
	if err := info.Close(); err != nil {
		t.Fatalf("info Close: %v", err)
	}
	if err := info.Update(ch); !errors.Is(err, errSystemdConnClosed) {
		t.Fatalf("info Update after Close: got %v want %v", err, errSystemdConnClosed)
	}
}

func TestSystemdCollectorsUpdateDialFailure(t *testing.T) {
	dialErr := errors.New("dbus unavailable")
	dial := func() (*dbus.Conn, error) { return nil, dialErr }
	ch := make(chan prometheus.Metric, 1)

	svc := &systemdServicesCollector{
		serviceState: prometheus.NewDesc("node_systemd_service_state", "x", []string{"name"}, nil),
		sc:           systemdConn{dial: dial},
	}
	err := svc.Update(ch)
	if err == nil || !strings.Contains(err.Error(), "couldn't get dbus connection") || !errors.Is(err, dialErr) {
		t.Fatalf("services Update dial failure: %v", err)
	}

	info := &systemdInfoCollector{sc: systemdConn{dial: dial}}
	err = info.Update(ch)
	if err == nil || !strings.Contains(err.Error(), "couldn't get dbus connection") || !errors.Is(err, dialErr) {
		t.Fatalf("info Update dial failure: %v", err)
	}
}

func TestServiceTypeFromProps(t *testing.T) {
	if got := serviceTypeFromProps(nil); got != "" {
		t.Fatalf("nil props: %q", got)
	}
	if got := serviceTypeFromProps(map[string]any{"Type": "oneshot"}); got != "oneshot" {
		t.Fatalf("Type: %q", got)
	}
	if got := serviceTypeFromProps(map[string]any{"Type": ""}); got != "" {
		t.Fatalf("empty Type: %q", got)
	}
}

func TestNRestartsFromProps(t *testing.T) {
	if _, ok := nRestartsFromProps(nil); ok {
		t.Fatal("missing NRestarts should be ok on older systemd")
	}
	if _, ok := nRestartsFromProps(map[string]any{"Type": "simple"}); ok {
		t.Fatal("expected missing NRestarts")
	}
	got, ok := nRestartsFromProps(map[string]any{"NRestarts": uint32(3)})
	if !ok || got != 3 {
		t.Fatalf("NRestarts: got %v ok=%v", got, ok)
	}
	if _, ok := nRestartsFromProps(map[string]any{"NRestarts": "nope"}); ok {
		t.Fatal("unexpected type should not parse")
	}
}
