// Copyright 2018 The Prometheus Authors
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
	"regexp"
)

type netDevFilter struct {
	ignorePattern *regexp.Regexp
	acceptPattern *regexp.Regexp
}

func newNetDevFilter(ignoredPattern, acceptPattern string) (f netDevFilter) {
	if ignoredPattern != "" {
		f.ignorePattern = regexp.MustCompile(ignoredPattern)
	}

	if acceptPattern != "" {
		f.acceptPattern = regexp.MustCompile(acceptPattern)
	}

	return
}

// ignores returns whether the device should be ignored
func (f *netDevFilter) ignored(name string) bool {
	return ((f.ignorePattern != nil && f.ignorePattern.MatchString(name)) ||
		(f.acceptPattern != nil && !f.acceptPattern.MatchString(name)))
}
