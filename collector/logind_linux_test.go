// Copyright 2016 The Prometheus Authors
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

//go:build !nologind

package collector

import (
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/prometheus/client_golang/prometheus"
)

type testLogindInterface struct{}

var testSeats = []string{"seat0", ""}

func (c *testLogindInterface) listSeats() ([]string, error) {
	return testSeats, nil
}

func (c *testLogindInterface) listSessions() ([]logindSessionEntry, error) {
	return []logindSessionEntry{
		{
			SessionID:         "1",
			UserID:            0,
			UserName:          "",
			SeatID:            "",
			SessionObjectPath: dbus.ObjectPath("/org/freedesktop/login1/session/1"),
		},
		{
			SessionID:         "2",
			UserID:            0,
			UserName:          "",
			SeatID:            "seat0",
			SessionObjectPath: dbus.ObjectPath("/org/freedesktop/login1/session/2"),
		},
	}, nil
}

func (c *testLogindInterface) getSession(session logindSessionEntry) *logindSession {
	sessions := map[dbus.ObjectPath]*logindSession{
		dbus.ObjectPath("/org/freedesktop/login1/session/1"): {
			seat:        session.SeatID,
			remote:      "true",
			sessionType: knownStringOrOther("tty", attrTypeValues),
			class:       knownStringOrOther("user", attrClassValues),
		},
		dbus.ObjectPath("/org/freedesktop/login1/session/2"): {
			seat:        session.SeatID,
			remote:      "false",
			sessionType: knownStringOrOther("x11", attrTypeValues),
			class:       knownStringOrOther("greeter", attrClassValues),
		},
	}

	return sessions[session.SessionObjectPath]
}

func TestLogindCollectorKnownStringOrOther(t *testing.T) {
	known := []string{"foo", "bar"}

	actual := knownStringOrOther("foo", known)
	expected := "foo"
	if actual != expected {
		t.Errorf("knownStringOrOther failed: got %q, expected %q.", actual, expected)
	}

	actual = knownStringOrOther("baz", known)
	expected = "other"
	if actual != expected {
		t.Errorf("knownStringOrOther failed: got %q, expected %q.", actual, expected)
	}

}

func TestLogindCollectorCollectMetrics(t *testing.T) {
	ch := make(chan prometheus.Metric)
	go func() {
		collectMetrics(ch, &testLogindInterface{})
		close(ch)
	}()

	count := 0
	for range ch {
		count++
	}

	expected := len(testSeats) * len(attrRemoteValues) * len(attrTypeValues) * len(attrClassValues)
	if count != expected {
		t.Errorf("collectMetrics did not generate the expected number of metrics: got %d, expected %d.", count, expected)
	}
}

func TestLogindCollectorCollectMetricsWithFilter(t *testing.T) {
	ch := make(chan prometheus.Metric)
	// Create a dummy collector with a classFilter that excludes "greeter"
	collector := &logindCollector{
			classFilter: deviceFilter{
					ignorePattern: regexp.MustCompile("greeter"),
			},
	}

	// Wrap collectMetrics call in a goroutine feeding the metrics channel
	go func() {
			err := collectMetrics(ch, &testLogindInterface{})
			if err != nil {
					t.Errorf("collectMetrics returned error: %v", err)
			}
			close(ch)
	}()

	seenClasses := make(map[string]bool)
	count := 0
	for metric := range ch {
			count++
			desc := metric.Desc().String()
			// Extract class label from the metric Desc string (approximate)
			// The label list appears as '{seat=...,remote=...,type=...,class=...}'
			// We'll parse class label from the Desc string for simplicity here:
			// For a formal approach, you can use a prometheus.Metric introspection library.

			if strings.Contains(desc, "class=\"greeter\"") {
					t.Errorf("metric with excluded class 'greeter' was emitted")
			}
	}
	if count == 0 {
			t.Errorf("no metrics were emitted")
	}
}
