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

//go:build !nocpu
// +build !nocpu

package collector

import (
	"runtime"
	"testing"
)

func TestCPU(t *testing.T) {
	var (
		fieldsCount = 5
		times, err  = getDragonFlyCPUTimes()
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(times) == 0 {
		t.Fatalf("no cputimes found")
	}

	want := runtime.NumCPU() * fieldsCount
	if len(times) != want {
		t.Fatalf("should have %d cpuTimes: got %d", want, len(times))
	}
}
