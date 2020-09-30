// Copyright 2020 The Prometheus Authors
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
	"fmt"
	"testing"
)

type Logger struct {
	messages []string
}

func (l Logger) Log(keyvals ...interface{}) error {
	for _, value := range keyvals {
		l.messages = append(l.messages, fmt.Sprintf("%s", value))
	}
	return nil
}

func TestParseLineSimple(t *testing.T) {
	key, value, err := parseLine("some_metric 12", 2)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "some_metric", key; want != got {
		t.Errorf("want %s, got %s", want, got)
	}
	if want, got := 12, value; want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestParseLine(t *testing.T) {
	key, value, err := parseLine("   some  multi word metric  15", 5)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "some_multi_word_metric", key; want != got {
		t.Errorf("want %s, got %s", want, got)
	}
	if want, got := 15, value; want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestGetZoneInfo(t *testing.T) {
	*zoneinfoPath = "fixtures/proc/zoneinfo"
	logger := Logger{}
	zones, err := getZoneInfo(logger)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := 0, len(logger.messages); want != got {
		t.Errorf("Logger has %d messages, wanted %d", got, want)
	}
	if want, got := 5, len(zones); want != got {
		t.Errorf("got %d zones , wanted %d", got, want)
	}
	if want, got := uint64(95612), zones[0].PerNodeStats["nr_inactive_anon"]; want != got {
		t.Errorf("got nr_inactive_anon %d , wanted %d", got, want)
	}
	if want, got := uint64(528427), zones[1].Pages["free"]; want != got {
		t.Errorf("got pages free %d , wanted %d", got, want)
	}
	if want, got := uint64(357), zones[1].PageSets[0].Count; want != got {
		t.Errorf("got pageset cpu 0 %d , wanted %d", got, want)
	}
	if want, got := 8, len(zones[1].PageSets); want != got {
		t.Errorf("got pagesets %d , wanted %d", got, want)
	}
}
