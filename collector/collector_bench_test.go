//go:build linux

// Copyright 2026 The Prometheus Authors
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
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs"
)

func BenchmarkParseFilesystemLabels(b *testing.B) {
	mountInfo := benchmarkMountInfo(256)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filesystems, err := parseFilesystemLabels(mountInfo)
		if err != nil {
			b.Fatal(err)
		}
		if len(filesystems) != len(mountInfo) {
			b.Fatalf("got %d filesystems, want %d", len(filesystems), len(mountInfo))
		}
	}
}

func BenchmarkTextfileConvertMetricFamily(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	metricFamily := benchmarkMetricFamily(256, 8)
	ch := make(chan prometheus.Metric, len(metricFamily.Metric)+1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		convertMetricFamily(metricFamily, ch, logger)
		drainMetrics(ch)
	}
}

func benchmarkMountInfo(count int) []*procfs.MountInfo {
	mountInfo := make([]*procfs.MountInfo, 0, count)
	for i := 0; i < count; i++ {
		mountInfo = append(mountInfo, &procfs.MountInfo{
			MajorMinorVer: fmt.Sprintf("%d:%d", i/16+8, i%16),
			Source:        fmt.Sprintf("/dev/vd%c", 'a'+rune(i%26)),
			MountPoint:    fmt.Sprintf("/var/lib/containers/storage/overlay/%d", i),
			FSType:        "ext4",
			Options: map[string]string{
				"rw":       "",
				"relatime": "",
				"discard":  "",
			},
			SuperOptions: map[string]string{
				"rw":     "",
				"errors": "remount-ro",
			},
		})
	}

	return mountInfo
}

func benchmarkMetricFamily(metricCount, labelCount int) *dto.MetricFamily {
	metrics := make([]*dto.Metric, 0, metricCount)
	for i := 0; i < metricCount; i++ {
		labels := make([]*dto.LabelPair, 0, labelCount)
		for j := 0; j < labelCount; j++ {
			// Leave some labels out on each metric so the benchmark exercises
			// label union normalization as well as desc reuse.
			if (i+j)%3 == 0 {
				continue
			}
			name := fmt.Sprintf("label_%d", j)
			value := fmt.Sprintf("value_%d_%d", i, j)
			labels = append(labels, &dto.LabelPair{
				Name:  stringPtr(name),
				Value: stringPtr(value),
			})
		}

		value := float64(i)
		metrics = append(metrics, &dto.Metric{
			Label: labels,
			Gauge: &dto.Gauge{Value: &value},
		})
	}

	metricType := dto.MetricType_GAUGE

	return &dto.MetricFamily{
		Name:   stringPtr("node_benchmark_textfile_metric"),
		Help:   stringPtr("Benchmark metric family for textfile collector."),
		Type:   &metricType,
		Metric: metrics,
	}
}

func drainMetrics(ch chan prometheus.Metric) {
	for len(ch) > 0 {
		<-ch
	}
}

func stringPtr(s string) *string {
	return &s
}
