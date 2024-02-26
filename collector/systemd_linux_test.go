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

//go:build !nosystemd
// +build !nosystemd

package collector

import (
	"regexp"
	"testing"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/go-kit/log"
)

// Creates mock UnitLists
func getUnitListFixtures() [][]unit {
	fixture1 := []unit{
		{
			UnitStatus: dbus.UnitStatus{
				Name:        "foo",
				Description: "foo desc",
				LoadState:   "loaded",
				ActiveState: "active",
				SubState:    "running",
				Followed:    "",
				Path:        "/org/freedesktop/systemd1/unit/foo",
				JobId:       0,
				JobType:     "",
				JobPath:     "/",
			},
		},
		{
			UnitStatus: dbus.UnitStatus{
				Name:        "bar",
				Description: "bar desc",
				LoadState:   "not-found",
				ActiveState: "inactive",
				SubState:    "dead",
				Followed:    "",
				Path:        "/org/freedesktop/systemd1/unit/bar",
				JobId:       0,
				JobType:     "",
				JobPath:     "/",
			},
		},
		{
			UnitStatus: dbus.UnitStatus{
				Name:        "foobar",
				Description: "bar desc",
				LoadState:   "not-found",
				ActiveState: "inactive",
				SubState:    "dead",
				Followed:    "",
				Path:        "/org/freedesktop/systemd1/unit/bar",
				JobId:       0,
				JobType:     "",
				JobPath:     "/",
			},
		},
		{
			UnitStatus: dbus.UnitStatus{
				Name:        "baz",
				Description: "bar desc",
				LoadState:   "not-found",
				ActiveState: "inactive",
				SubState:    "dead",
				Followed:    "",
				Path:        "/org/freedesktop/systemd1/unit/bar",
				JobId:       0,
				JobType:     "",
				JobPath:     "/",
			},
		},
	}

	fixture2 := []unit{}

	return [][]unit{fixture1, fixture2}
}

func TestSystemdIgnoreFilter(t *testing.T) {
	fixtures := getUnitListFixtures()
	includePattern := regexp.MustCompile("^foo$")
	excludePattern := regexp.MustCompile("^bar$")
	filtered := filterUnits(fixtures[0], includePattern, excludePattern, log.NewNopLogger())
	for _, unit := range filtered {
		if excludePattern.MatchString(unit.Name) || !includePattern.MatchString(unit.Name) {
			t.Error(unit.Name, "should not be in the filtered list")
		}
	}
}
func TestSystemdIgnoreFilterDefaultKeepsAll(t *testing.T) {
	logger := log.NewNopLogger()
	c, err := NewSystemdCollector(logger)
	if err != nil {
		t.Fatal(err)
	}
	fixtures := getUnitListFixtures()
	collector := c.(*systemdCollector)
	filtered := filterUnits(fixtures[0], collector.systemdUnitIncludePattern, collector.systemdUnitExcludePattern, logger)
	// Adjust fixtures by 3 "not-found" units.
	if len(filtered) != len(fixtures[0])-3 {
		t.Error("Default filters removed units")
	}
}

func TestSystemdSummary(t *testing.T) {
	fixtures := getUnitListFixtures()
	summary := summarizeUnits(fixtures[0])

	for _, state := range unitStatesName {
		if state == "inactive" {
			testSummaryHelper(t, state, summary[state], 3.0)
		} else if state == "active" {
			testSummaryHelper(t, state, summary[state], 1.0)
		} else {
			testSummaryHelper(t, state, summary[state], 0.0)
		}
	}
}

func testSummaryHelper(t *testing.T, state string, actual float64, expected float64) {
	if actual != expected {
		t.Errorf("Summary mode didn't count %s jobs correctly. Actual: %f, expected: %f", state, actual, expected)
	}
}
