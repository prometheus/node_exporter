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

package collector

import (
	"bytes"
	"strings"
)

func SafeDereference[T any](s ...*T) []T {
	var resolved []T
	for _, v := range s {
		if v != nil {
			resolved = append(resolved, *v)
		} else {
			var zeroValue T
			resolved = append(resolved, zeroValue)
		}
	}
	return resolved
}

// SafeBytesToString takes a slice of bytes and sanitizes it for Prometheus label
// values.
// * Terminate the string at the first null byte.
// * Convert any invalid UTF-8 to "�".
func SafeBytesToString(b []byte) string {
	var s string
	before, _, ok := bytes.Cut(b, []byte{0})
	if !ok {
		s = string(b)
	} else {
		s = string(before)
	}
	return strings.ToValidUTF8(s, "�")
}
