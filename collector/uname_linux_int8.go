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

// +build 386 amd64 arm64 mips64 mips64le mips mipsle
// +build linux
// +build !nouname

package collector

func unameToString(input [65]int8) string {
	var str string
	for _, a := range input {
		if a == 0 {
			break
		}
		str += string(a)
	}
	return str
}
