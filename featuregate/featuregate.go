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
	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"golang.org/x/mod/semver"
	"io"
	"os"
)

const versionFilePath = "../VERSION"

var (
	enabledDefaultFeatureGates = map[string]bool{}
	gotFeatureGates            = kingpin.Flag("feature-gates", "A set of key=value pairs that describe enabled feature gates for experimental features.").Default("").StringMap()
)

type (
	FeatureGate struct {
		Name                 string
		Desc                 string
		InitialAddingVersion string
		FinalRetiringVersion string
		logger               log.Logger
	}
)

// NewFeatureGate creates a new feature gate, but errors if the feature gate is invalid.
func NewFeatureGate(name, desc, initialAddingVersion, finalRetiringVersion string) *FeatureGate {
	return &FeatureGate{
		Name:                 name,
		Desc:                 desc,
		InitialAddingVersion: semver.Canonical(initialAddingVersion),
		FinalRetiringVersion: semver.Canonical(finalRetiringVersion),
		logger:               log.NewLogfmtLogger(os.Stdout),
	}
}

func (fg *FeatureGate) String() string {
	return fg.Name
}

// IsEnabled returns true if the feature gate is enabled, false if it is disabled.
func (fg *FeatureGate) IsEnabled() (bool, error) {
	_, err := fg.isValid()
	if err != nil {
		return false, err
	}

	v, err := currentVersion()
	if err != nil {
		return false, err
	}

	// Check if the feature gate is within the specified version range.
	// Note: Do not return immediately to allow overriding the specified range enablement.
	isFeatureGateWithinSpecifiedVersionRange := semver.Compare(fg.InitialAddingVersion, v) != 1 &&
		semver.Compare(fg.FinalRetiringVersion, v) != -1

	// Check if the feature gate is enabled or disabled by default.
	var isFeatureGateEnabledOrDisabledByDefault *bool
	_, isFeatureGateInDefaultFeatureGates := enabledDefaultFeatureGates[fg.Name]
	if isFeatureGateInDefaultFeatureGates {
		isFeatureGateEnabledOrDisabledByDefault = ptrTo(enabledDefaultFeatureGates[fg.Name])
	}

	// Check if the feature gate is enabled or disabled by the user.
	var isFeatureGateEnabledOrDisabledByFlag *bool
	_, isFeatureGateInFlag := (*gotFeatureGates)[fg.Name]
	if isFeatureGateInFlag {
		t := (*gotFeatureGates)[fg.Name]
		// Check for exact values since none of these is the default value (nil).
		if t == "true" {
			isFeatureGateEnabledOrDisabledByFlag = ptrTo(true)
		} else if t == "false" {
			isFeatureGateEnabledOrDisabledByFlag = ptrTo(false)
		}
	}

	// Check for overrides by the user. This taken precedence over the default states and version range.
	featureGateState := isFeatureGateWithinSpecifiedVersionRange
	if isFeatureGateEnabledOrDisabledByDefault != nil {
		featureGateState = *isFeatureGateEnabledOrDisabledByDefault
	}
	if isFeatureGateEnabledOrDisabledByFlag != nil {
		if isFeatureGateEnabledOrDisabledByDefault != nil &&
			*isFeatureGateEnabledOrDisabledByFlag != *isFeatureGateEnabledOrDisabledByDefault {
			fg.logger.Log(fg.Name, "default feature gate state overridden by user, this is not recommended")
		}
		featureGateState = *isFeatureGateEnabledOrDisabledByFlag
	}
	return featureGateState, nil
}

// isValid checks if the feature gate is valid.
// Eventually, the current version will surpass the final retiring version, so we don't check for that.
func (fg *FeatureGate) isValid() (*bool, error) {
	v, err := currentVersion()
	if err != nil {
		return nil, err
	}
	if !semver.IsValid(fg.InitialAddingVersion) {
		return ptrTo(false), fmt.Errorf("invalid adding version %q", fg.InitialAddingVersion)
	}
	if !semver.IsValid(fg.FinalRetiringVersion) {
		return ptrTo(false), fmt.Errorf("invalid retiring version %q", fg.FinalRetiringVersion)
	}
	if semver.Compare(fg.InitialAddingVersion, fg.FinalRetiringVersion) != -1 {
		return ptrTo(false), fmt.Errorf("adding version %q is not before the retiring version %q", fg.InitialAddingVersion, fg.FinalRetiringVersion)
	}
	if semver.Compare(fg.InitialAddingVersion, v) == 1 {
		return ptrTo(false), fmt.Errorf("adding version %q is greater than the current version %q", fg.InitialAddingVersion, v)
	}
	return ptrTo(true), nil
}

// currentVersion reads the current version from the VERSION file.
func currentVersion() (string, error) {
	file, err := os.Open(versionFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	version, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return "v" + string(version[:len(version)-1]), nil
}
