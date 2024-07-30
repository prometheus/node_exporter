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

//go:build !nomeminfo
// +build !nomeminfo

package collector

import (
	"log/slog"

	"github.com/power-devops/perfstat"
)

type meminfoCollector struct {
	logger *slog.Logger
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewMeminfoCollector(logger *slog.Logger) (Collector, error) {
	return &meminfoCollector{
		logger: logger,
	}, nil
}

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	stats, err := perfstat.MemoryTotalStat()
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"total_bytes":     float64(stats.RealTotal * 4096),
		"free_bytes":      float64(stats.RealFree * 4096),
		"available_bytes": float64(stats.RealAvailable * 4096),
	}, nil
}
