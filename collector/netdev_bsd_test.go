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

// +build !nonetdev
// +build freebsd dragonfly

package collector

import "testing"

type uintToStringTest struct {
	in  uint64
	out string
}

var uinttostringtests = []uintToStringTest{
	// Copied base10 values from strconv's tests:
	{0, "0"},
	{1, "1"},
	{12345678, "12345678"},
	{1<<31 - 1, "2147483647"},
	{1 << 31, "2147483648"},
	{1<<31 + 1, "2147483649"},
	{1<<32 - 1, "4294967295"},
	{1 << 32, "4294967296"},
	{1<<32 + 1, "4294967297"},
	{1 << 50, "1125899906842624"},
	{1<<63 - 1, "9223372036854775807"},

	// Some values that convert correctly on amd64, but not on i386.
	{0x1bf0c640a, "7500227594"},
	{0xbee5df75, "3202735989"},
}

func TestUintToString(t *testing.T) {
	for _, test := range uinttostringtests {
		is := convertFreeBSDCPUTime(test.in)
		if is != test.out {
			t.Errorf("convertFreeBSDCPUTime(%v) = %v want %v",
				test.in, is, test.out)
		}
	}
}
