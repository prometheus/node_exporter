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

//go:build !nomeminfo
// +build !nomeminfo

package collector

import (
	"io"
	"log/slog"
	"testing"
)

func TestMemInfo(t *testing.T) {
	*procPath = "fixtures/proc"
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	collector, err := NewMeminfoCollector(logger)
	if err != nil {
		panic(err)
	}

	memInfo, err := collector.(*meminfoCollector).getMemInfo()
	if err != nil {
		panic(err)
	}

	if want, got := 3831959552.0, memInfo["MemTotal_bytes"]; want != got {
		t.Errorf("want memory total %f, got %f", want, got)
	}

	if want, got := 3787456512.0, memInfo["DirectMap2M_bytes"]; want != got {
		t.Errorf("want memory directMap2M %f, got %f", want, got)
	}
}
