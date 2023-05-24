// Copyright 2023 The Prometheus Authors
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

func TestCPUTimes(t *testing.T) {
	times, err := getCPUTimes()
	if err != nil {
		t.Fatalf("getCPUTimes returned error: %v", err)
	}

	if len(times) == 0 {
		t.Fatalf("no CPU times found")
	}

	if got, want := len(times), runtime.NumCPU(); got != want {
		t.Fatalf("unexpected # of CPU times; got %d want %d", got, want)
	}
}

func TestCPUTemperatures(t *testing.T) {
	_, err := getCPUTemperatures()
	if err != nil {
		t.Fatalf("getCPUTemperatures returned error: %v", err)
	}
}
