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

package featuregate

import (
	"fmt"
	"golang.org/x/mod/semver"
	"os"
	"testing"
)

func TestFeatureGate(t *testing.T) {
	// Set a valid version range.
	cleanupFn, err := changeVersionTo("v1.8.1")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = cleanupFn()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Default enabled or disabled feature gates.
	enabledDefaultFeatureGates = map[string]bool{
		"TestFoo":  true,
		"TestBar":  false,
		"TestBaz":  false,
		"TestNorf": true,
	}

	// User enabled or disabled feature gates.
	gotFeatureGates = &map[string]string{
		"TestFoo": "false",
		"TestBar": "true",
		"TestQue": "false",
	}

	// Check feature gate enablement under various parameters.
	testcases := []struct {
		name        string
		fg          *FeatureGate
		enabled     bool
		shouldError bool
	}{
		{
			name: "feature gate 'TestFoo', should be disabled due to user override",
			fg:   NewFeatureGate("TestFoo", "", "v1.8.0", "v1.8.2"),
		},
		{
			name:    "feature gate 'TestBar', should be enabled due to user override",
			fg:      NewFeatureGate("TestBar", "", "v1.8.0", "v1.8.2"),
			enabled: true,
		},
		{
			name: "feature gate 'TestBaz', should be disabled by default",
			fg:   NewFeatureGate("TestBaz", "", "v1.8.0", "v1.8.2"),
		},
		{
			name:    "feature gate 'TestNorf', should be enabled by default",
			fg:      NewFeatureGate("TestNorf", "", "v1.8.0", "v1.8.2"),
			enabled: true,
		},
		{
			name: "feature gate 'TestQue', should be disabled due to user override",
			fg:   NewFeatureGate("TestQue", "", "v1.8.0", "v1.8.2"),
		},
		{
			name:    "feature gate 'TestQux', should be enabled due to valid version range",
			fg:      NewFeatureGate("TestQux", "", "v1.8.0", "v1.8.2"),
			enabled: true,
		},
		{
			// This will happen with all feature gates eventually if they are not removed, and expected.
			// We do not expect an error, or an enabled state.
			name: "feature gate 'TestQux', should be disabled due to current version exceeding final retiring version",
			fg:   NewFeatureGate("TestQux", "", "v1", "v1.8.0"),
		},
		{
			name:        "feature gate 'TestQux', should be disabled due to invalid version range (current version < initial adding version)",
			fg:          NewFeatureGate("TestQux", "", "v1.8.2", "v2"),
			shouldError: true,
		},
		{
			name:        "feature gate 'TestQux', should be disabled due to invalid version range (final retiring version < initial adding version)",
			fg:          NewFeatureGate("TestQux", "", "v1.8.2", "v1.8.0"),
			shouldError: true,
		},
		{
			name:        "feature gate 'TestQux', should be disabled due to invalid initial adding version (missing 'v' prefix)",
			fg:          NewFeatureGate("TestQux", "", "1.8.0", "v1.8.2"),
			shouldError: true,
		},
		{
			name:        "feature gate 'qux', should be disabled due to invalid final retiring version (missing 'v' prefix)",
			fg:          NewFeatureGate("TestQux", "", "v1.8.0", "1.8.2"),
			shouldError: true,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			got, err := testcase.fg.IsEnabled()
			if err != nil {
				if !testcase.shouldError {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if testcase.shouldError {
					t.Fatalf("expected error, got nil")
				}
			}
			if got != testcase.enabled {
				t.Fatalf("expected %v, got %v", testcase.enabled, got)
			}
		})
	}
}

func changeVersionTo(version string) (func() error, error) {
	if isValid := semver.IsValid(version); !isValid {
		return nil, fmt.Errorf("invalid version: %v", version)
	}
	v, err := currentVersion()
	if err != nil {
		return nil, err
	}
	if isValid := semver.IsValid(v); !isValid {
		return nil, fmt.Errorf("invalid version: %v", v)
	}
	cleanupFn := func() error {
		err := os.WriteFile(versionFilePath, []byte(v + "\n")[1:], 0644)
		if err != nil {
			return err
		}
		return nil
	}
	err = os.WriteFile(versionFilePath, []byte(version + "\n")[1:], 0644)
	if err != nil {
		return nil, err
	}
	return cleanupFn, nil
}
