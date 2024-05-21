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

//go:build !386
// +build !386

package collector

import (
	"github.com/go-kit/log"
	"golang.org/x/sys/unix"
	"math/big"
	"runtime"
	"testing"
)

func TestOverflowHandling(t *testing.T) {
	if runtime.GOARCH == "386" {
		t.Skip("Skipping on 32-bit architecture")
	}
	factor := uint64(24)
	testcases := []struct {
		fsStats  *unix.Statfs_t
		overflow bool
	}{
		{
			fsStats: &unix.Statfs_t{
				Bsize:  int64(float64Mantissa),
				Blocks: factor,
			},
			overflow: true,
		},
		{
			fsStats: &unix.Statfs_t{
				Bsize:  int64(float64Mantissa / factor),
				Blocks: factor,
			},
			overflow: false,
		},
	}
	for _, testcase := range testcases {
		var fsC = filesystemCollector{
			logger:   log.NewNopLogger(),
			_fsStats: testcase.fsStats,
		}
		fsStats := fsC.processStat(filesystemLabels{
			mountPoint: "/",
		})

		oldsize := float64(testcase.fsStats.Blocks) * float64(testcase.fsStats.Bsize)
		size := new(big.Float).SetFloat64(fsStats.size)
		actualFloat64, _ := new(big.Int).Mul(new(big.Int).SetUint64(testcase.fsStats.Blocks), new(big.Int).SetInt64(testcase.fsStats.Bsize)).Float64()
		actual := new(big.Float).SetFloat64(actualFloat64)
		if testcase.overflow {
			oldsizeActualDiff := new(big.Float).Abs(new(big.Float).Sub(actual, new(big.Float).SetFloat64(oldsize)))
			sizeActualDiff := new(big.Float).Abs(new(big.Float).Sub(actual, size))
			if sizeActualDiff.Cmp(oldsizeActualDiff) < 0 {
				t.Errorf("Expected size to be closer to %f than %f, got %f instead", actual, oldsize, size)
			}
		} else {
			if size.Cmp(actual) != 0 {
				t.Errorf("Expected size to be %f, got %f instead", actual, size)
			}
		}
	}
}
