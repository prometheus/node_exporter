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

//go:build !noswap
// +build !noswap

package collector

import (
	"io"
	"log/slog"
	"testing"
)

func TestSwap(t *testing.T) {
	*procPath = "fixtures/proc"
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	collector, err := NewSwapCollector(logger)
	if err != nil {
		panic(err)
	}

	swapInfo, err := collector.(*swapCollector).getSwapInfo()
	if err != nil {
		panic(err)
	}

	if want, got := "/dev/zram0", swapInfo[0].Device; want != got {
		t.Errorf("want swap device %s, got %s", want, got)
	}

	if want, got := "partition", swapInfo[0].Type; want != got {
		t.Errorf("want swap type %s, got %s", want, got)
	}

	if want, got := 100, swapInfo[0].Priority; want != got {
		t.Errorf("want swap priority %d, got %d", want, got)
	}

	if want, got := 8388604, swapInfo[0].Size; want != got {
		t.Errorf("want swap size %d, got %d", want, got)
	}

	if want, got := 76, swapInfo[0].Used; want != got {
		t.Errorf("want swpa used %d, got %d", want, got)
	}
}
