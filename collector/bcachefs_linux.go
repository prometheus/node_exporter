// Copyright 2025 The Prometheus Authors
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

//go:build !nobcachefs

package collector

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/bcachefs"
)

func init() {
	registerCollector("bcachefs", defaultEnabled, NewBcachefsCollector)
}

// bcachefsCollector collects metrics from bcachefs filesystems.
type bcachefsCollector struct {
	fs     bcachefs.FS
	logger *slog.Logger
}

// NewBcachefsCollector returns a new Collector exposing bcachefs statistics.
func NewBcachefsCollector(logger *slog.Logger) (Collector, error) {
	fs, err := bcachefs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %w", err)
	}

	return &bcachefsCollector{
		fs:     fs,
		logger: logger,
	}, nil
}

// Update retrieves and exports bcachefs statistics.
func (c *bcachefsCollector) Update(ch chan<- prometheus.Metric) error {
	const subsystem = "bcachefs"

	stats, err := c.fs.Stats()
	if err != nil {
		if os.IsNotExist(err) {
			c.logger.Debug("bcachefs sysfs path does not exist", "path", sysFilePath("fs/bcachefs"))
			return ErrNoData
		}
		return fmt.Errorf("failed to retrieve bcachefs stats: %w", err)
	}

	if len(stats) == 0 {
		return ErrNoData
	}

	for _, s := range stats {
		uuid := s.UUID

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "info"),
				"Filesystem information.",
				[]string{"uuid"},
				nil,
			),
			prometheus.GaugeValue,
			1,
			uuid,
		)

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "btree_cache_size_bytes"),
				"Btree cache memory usage in bytes.",
				[]string{"uuid"},
				nil,
			),
			prometheus.GaugeValue,
			float64(s.BtreeCacheSizeBytes),
			uuid,
		)

		for algorithm, comp := range s.Compression {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "compression_compressed_bytes"),
					"Compressed size by algorithm.",
					[]string{"uuid", "algorithm"},
					nil,
				),
				prometheus.GaugeValue,
				float64(comp.CompressedBytes),
				uuid, algorithm,
			)
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "compression_uncompressed_bytes"),
					"Uncompressed size by algorithm.",
					[]string{"uuid", "algorithm"},
					nil,
				),
				prometheus.GaugeValue,
				float64(comp.UncompressedBytes),
				uuid, algorithm,
			)
		}

		for errorType, errStats := range s.Errors {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "errors_total"),
					"Error count by error type.",
					[]string{"uuid", "error_type"},
					nil,
				),
				prometheus.CounterValue,
				float64(errStats.Count),
				uuid, errorType,
			)
		}

		for counterName, counterStats := range s.Counters {
			metricName := sanitizeMetricName(counterName) + "_total"
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, metricName),
					fmt.Sprintf("Bcachefs counter %s since filesystem creation.", counterName),
					[]string{"uuid"},
					nil,
				),
				prometheus.CounterValue,
				float64(counterStats.SinceFilesystemCreation),
				uuid,
			)
		}

		for writeType, writeStats := range s.BtreeWrites {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "btree_writes_total"),
					"Number of btree writes by type.",
					[]string{"uuid", "type"},
					nil,
				),
				prometheus.CounterValue,
				float64(writeStats.Count),
				uuid, writeType,
			)
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "btree_write_average_size_bytes"),
					"Average btree write size by type.",
					[]string{"uuid", "type"},
					nil,
				),
				prometheus.GaugeValue,
				float64(writeStats.SizeBytes),
				uuid, writeType,
			)
		}

		for device, devStats := range s.Devices {
			if devStats == nil {
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "device_info"),
					"Device information.",
					[]string{"uuid", "device", "label", "state"},
					nil,
				),
				prometheus.GaugeValue,
				1,
				uuid, device, devStats.Label, devStats.State,
			)

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "device_bucket_size_bytes"),
					"Bucket size in bytes.",
					[]string{"uuid", "device"},
					nil,
				),
				prometheus.GaugeValue,
				float64(devStats.BucketSizeBytes),
				uuid, device,
			)

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "device_buckets"),
					"Total number of buckets.",
					[]string{"uuid", "device"},
					nil,
				),
				prometheus.GaugeValue,
				float64(devStats.Buckets),
				uuid, device,
			)

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "device_durability"),
					"Device durability setting.",
					[]string{"uuid", "device"},
					nil,
				),
				prometheus.GaugeValue,
				float64(devStats.Durability),
				uuid, device,
			)

			for op, dataTypes := range devStats.IODone {
				for dataType, value := range dataTypes {
					ch <- prometheus.MustNewConstMetric(
						prometheus.NewDesc(
							prometheus.BuildFQName(namespace, subsystem, "device_io_done_bytes_total"),
							"IO bytes by operation type and data type.",
							[]string{"uuid", "device", "operation", "data_type"},
							nil,
						),
						prometheus.CounterValue,
						float64(value),
						uuid, device, op, dataType,
					)
				}
			}

			for errorType, value := range devStats.IOErrors {
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						prometheus.BuildFQName(namespace, subsystem, "device_io_errors_total"),
						"IO errors by error type.",
						[]string{"uuid", "device", "type"},
						nil,
					),
					prometheus.CounterValue,
					float64(value),
					uuid, device, errorType,
				)
			}
		}
	}

	return nil
}

// sanitizeMetricName converts a string to a valid Prometheus metric name component.
func sanitizeMetricName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return re.ReplaceAllString(name, "_")
}
