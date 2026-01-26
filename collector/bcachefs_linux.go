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
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("bcachefs", defaultEnabled, NewBcachefsCollector)
}

// bcachefsCollector collects metrics from bcachefs filesystems.
type bcachefsCollector struct {
	logger *slog.Logger
}

// NewBcachefsCollector returns a new Collector exposing bcachefs statistics.
func NewBcachefsCollector(logger *slog.Logger) (Collector, error) {
	return &bcachefsCollector{
		logger: logger,
	}, nil
}

// Update retrieves and exports bcachefs statistics.
func (c *bcachefsCollector) Update(ch chan<- prometheus.Metric) error {
	const subsystem = "bcachefs"

	bcachefsPath := sysFilePath("fs/bcachefs")

	entries, err := os.ReadDir(bcachefsPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.logger.Debug("bcachefs sysfs path does not exist", "path", bcachefsPath)
			return ErrNoData
		}
		return fmt.Errorf("failed to read bcachefs sysfs: %w", err)
	}

	if len(entries) == 0 {
		return ErrNoData
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		uuid := entry.Name()
		fsPath := filepath.Join(bcachefsPath, uuid)

		// Emit info metric
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

		// Parse btree_cache_size
		if btreeCacheSize, err := c.parseBtreeCacheSize(fsPath); err == nil {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "btree_cache_size_bytes"),
					"Btree cache memory usage in bytes.",
					[]string{"uuid"},
					nil,
				),
				prometheus.GaugeValue,
				btreeCacheSize,
				uuid,
			)
		} else {
			c.logger.Debug("failed to parse btree_cache_size", "uuid", uuid, "err", err)
		}

		// Parse compression_stats
		if err := c.parseCompressionStats(ch, fsPath, uuid); err != nil {
			c.logger.Debug("failed to parse compression_stats", "uuid", uuid, "err", err)
		}

		// Parse errors
		if err := c.parseErrors(ch, fsPath, uuid); err != nil {
			c.logger.Debug("failed to parse errors", "uuid", uuid, "err", err)
		}

		// Parse counters
		if err := c.parseCounters(ch, fsPath, uuid); err != nil {
			c.logger.Debug("failed to parse counters", "uuid", uuid, "err", err)
		}

		// Parse btree_write_stats
		if err := c.parseBtreeWriteStats(ch, fsPath, uuid); err != nil {
			c.logger.Debug("failed to parse btree_write_stats", "uuid", uuid, "err", err)
		}

		// Parse device stats
		if err := c.parseDevices(ch, fsPath, uuid); err != nil {
			c.logger.Debug("failed to parse devices", "uuid", uuid, "err", err)
		}
	}

	return nil
}

// parseBtreeCacheSize parses the btree_cache_size file which contains a value with unit suffix (e.g., "524M").
func (c *bcachefsCollector) parseBtreeCacheSize(fsPath string) (float64, error) {
	data, err := os.ReadFile(filepath.Join(fsPath, "btree_cache_size"))
	if err != nil {
		return 0, err
	}
	return parseHumanReadableBytes(strings.TrimSpace(string(data)))
}

// parseHumanReadableBytes converts human-readable byte sizes (e.g., "524M", "2.00G", "1024k") to bytes.
func parseHumanReadableBytes(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	// Handle potential unit suffixes
	multiplier := float64(1)
	lastChar := s[len(s)-1]

	switch lastChar {
	case 'k', 'K':
		multiplier = 1024
		s = s[:len(s)-1]
	case 'm', 'M':
		multiplier = 1024 * 1024
		s = s[:len(s)-1]
	case 'g', 'G':
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	case 't', 'T':
		multiplier = 1024 * 1024 * 1024 * 1024
		s = s[:len(s)-1]
	case 'p', 'P':
		multiplier = 1024 * 1024 * 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, err
	}

	return value * multiplier, nil
}

// parseCompressionStats parses the compression_stats file.
// Format:
// type                    compressed  uncompressed  average extent size
// lz4:                         19.3G         67.0G                 112k
// incompressible:               5.5G          5.5G                  22k
func (c *bcachefsCollector) parseCompressionStats(ch chan<- prometheus.Metric, fsPath, uuid string) error {
	const subsystem = "bcachefs"

	file, err := os.Open(filepath.Join(fsPath, "compression_stats"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip header line
		if lineNum == 1 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		algorithm := strings.TrimSuffix(fields[0], ":")
		compressed, err := parseHumanReadableBytes(fields[1])
		if err != nil {
			c.logger.Debug("failed to parse compressed size", "algorithm", algorithm, "err", err)
			continue
		}
		uncompressed, err := parseHumanReadableBytes(fields[2])
		if err != nil {
			c.logger.Debug("failed to parse uncompressed size", "algorithm", algorithm, "err", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "compression_compressed_bytes"),
				"Compressed size by algorithm.",
				[]string{"uuid", "algorithm"},
				nil,
			),
			prometheus.GaugeValue,
			compressed,
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
			uncompressed,
			uuid, algorithm,
		)
	}

	return scanner.Err()
}

// parseErrors parses the errors file.
// Format: space-separated: error_type count timestamp
// Example: btree_node_read_err 5 1234567890
func (c *bcachefsCollector) parseErrors(ch chan<- prometheus.Metric, fsPath, uuid string) error {
	const subsystem = "bcachefs"

	file, err := os.Open(filepath.Join(fsPath, "errors"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		errorType := fields[0]
		count, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			c.logger.Debug("failed to parse error count", "error_type", errorType, "err", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "errors_total"),
				"Error count by error type.",
				[]string{"uuid", "error_type"},
				nil,
			),
			prometheus.CounterValue,
			count,
			uuid, errorType,
		)
	}

	return scanner.Err()
}

// parseCounters parses the counters directory.
// Each file has two lines:
// since mount:                  12345
// since filesystem creation:    67890
// We use the "since filesystem creation" value.
func (c *bcachefsCollector) parseCounters(ch chan<- prometheus.Metric, fsPath, uuid string) error {
	const subsystem = "bcachefs"

	countersPath := filepath.Join(fsPath, "counters")
	entries, err := os.ReadDir(countersPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		counterName := entry.Name()
		counterPath := filepath.Join(countersPath, counterName)

		value, err := c.parseCounterFile(counterPath)
		if err != nil {
			c.logger.Debug("failed to parse counter", "counter", counterName, "err", err)
			continue
		}

		// Convert counter name to Prometheus-friendly format (e.g., btree_node_read -> btree_node_read_total)
		metricName := sanitizeMetricName(counterName) + "_total"

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, metricName),
				fmt.Sprintf("Bcachefs counter %s since filesystem creation.", counterName),
				[]string{"uuid"},
				nil,
			),
			prometheus.CounterValue,
			value,
			uuid,
		)
	}

	return nil
}

// parseCounterFile parses a single counter file and returns the "since filesystem creation" value.
func (c *bcachefsCollector) parseCounterFile(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "since filesystem creation:") {
			valueStr := strings.TrimPrefix(line, "since filesystem creation:")
			valueStr = strings.TrimSpace(valueStr)
			return strconv.ParseFloat(valueStr, 64)
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("since filesystem creation line not found")
}

// sanitizeMetricName converts a string to a valid Prometheus metric name component.
func sanitizeMetricName(name string) string {
	// Replace hyphens and other invalid characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return re.ReplaceAllString(name, "_")
}

// parseDevices parses device information from dev-* directories.
func (c *bcachefsCollector) parseDevices(ch chan<- prometheus.Metric, fsPath, uuid string) error {
	const subsystem = "bcachefs"

	entries, err := os.ReadDir(fsPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "dev-") {
			continue
		}

		device := strings.TrimPrefix(entry.Name(), "dev-")
		devPath := filepath.Join(fsPath, entry.Name())

		// Parse device info (label, state)
		label := c.readSysfsFile(filepath.Join(devPath, "label"))
		state := c.readSysfsFile(filepath.Join(devPath, "state"))

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "device_info"),
				"Device information.",
				[]string{"uuid", "device", "label", "state"},
				nil,
			),
			prometheus.GaugeValue,
			1,
			uuid, device, label, state,
		)

		// Parse bucket_size
		if bucketSize, err := c.parseDeviceBucketSize(devPath); err == nil {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "device_bucket_size_bytes"),
					"Bucket size in bytes.",
					[]string{"uuid", "device"},
					nil,
				),
				prometheus.GaugeValue,
				bucketSize,
				uuid, device,
			)
		} else {
			c.logger.Debug("failed to parse bucket_size", "device", device, "err", err)
		}

		// Parse nbuckets
		if nbuckets, err := c.parseDeviceUint(devPath, "nbuckets"); err == nil {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "device_buckets"),
					"Total number of buckets.",
					[]string{"uuid", "device"},
					nil,
				),
				prometheus.GaugeValue,
				float64(nbuckets),
				uuid, device,
			)
		} else {
			c.logger.Debug("failed to parse nbuckets", "device", device, "err", err)
		}

		// Parse durability
		if durability, err := c.parseDeviceUint(devPath, "durability"); err == nil {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "device_durability"),
					"Device durability setting.",
					[]string{"uuid", "device"},
					nil,
				),
				prometheus.GaugeValue,
				float64(durability),
				uuid, device,
			)
		} else {
			c.logger.Debug("failed to parse durability", "device", device, "err", err)
		}

		// Parse io_done
		if err := c.parseDeviceIODone(ch, devPath, uuid, device); err != nil {
			c.logger.Debug("failed to parse io_done", "device", device, "err", err)
		}

		// Parse io_errors
		if err := c.parseDeviceIOErrors(ch, devPath, uuid, device); err != nil {
			c.logger.Debug("failed to parse io_errors", "device", device, "err", err)
		}
	}

	return nil
}

// readSysfsFile reads a sysfs file and returns its content as a trimmed string.
func (c *bcachefsCollector) readSysfsFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// parseDeviceBucketSize parses the bucket_size file.
func (c *bcachefsCollector) parseDeviceBucketSize(devPath string) (float64, error) {
	data, err := os.ReadFile(filepath.Join(devPath, "bucket_size"))
	if err != nil {
		return 0, err
	}
	return parseHumanReadableBytes(strings.TrimSpace(string(data)))
}

// parseDeviceUint parses a sysfs file containing a single unsigned integer.
func (c *bcachefsCollector) parseDeviceUint(devPath, filename string) (uint64, error) {
	data, err := os.ReadFile(filepath.Join(devPath, filename))
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

// parseDeviceIODone parses the io_done file.
// Format:
// read:
//
//	sb          :     3989504
//	btree       :  4411097088
//	user        :5768222552064
//
// write:
//
//	sb          :    31417344
//	btree       :     1171456
//	user        : 39196815360
func (c *bcachefsCollector) parseDeviceIODone(ch chan<- prometheus.Metric, devPath, uuid, device string) error {
	const subsystem = "bcachefs"

	file, err := os.Open(filepath.Join(devPath, "io_done"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentOp string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "read:" || trimmed == "write:" {
			currentOp = strings.TrimSuffix(trimmed, ":")
			continue
		}

		if currentOp == "" || trimmed == "" {
			continue
		}

		// Parse "dataType: value" format (value is in bytes as integer)
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}

		dataType := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			c.logger.Debug("failed to parse io_done value", "dataType", dataType, "err", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "device_io_done_bytes_total"),
				"IO bytes by operation type and data type.",
				[]string{"uuid", "device", "operation", "data_type"},
				nil,
			),
			prometheus.CounterValue,
			value,
			uuid, device, currentOp, dataType,
		)
	}

	return scanner.Err()
}

// parseDeviceIOErrors parses the io_errors file.
// Format:
// IO errors since filesystem creation
//
//	read:    197346
//	write:   0
//	checksum:0
//
// IO errors since 8 y ago
//
//	read:    197346
//	...
//
// We only parse the "IO errors since filesystem creation" section.
func (c *bcachefsCollector) parseDeviceIOErrors(ch chan<- prometheus.Metric, devPath, uuid, device string) error {
	const subsystem = "bcachefs"

	file, err := os.Open(filepath.Join(devPath, "io_errors"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inCreationSection := false

	for scanner.Scan() {
		line := scanner.Text()

		// Check for section headers
		if strings.HasPrefix(line, "IO errors since filesystem creation") {
			inCreationSection = true
			continue
		}
		if strings.HasPrefix(line, "IO errors since ") {
			// This is another section (e.g., "IO errors since 8 y ago"), stop parsing
			break
		}

		if !inCreationSection {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Parse "type: value" format
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}

		errorType := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			c.logger.Debug("failed to parse io_errors value", "errorType", errorType, "err", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "device_io_errors_total"),
				"IO errors by error type.",
				[]string{"uuid", "device", "type"},
				nil,
			),
			prometheus.CounterValue,
			value,
			uuid, device, errorType,
		)
	}

	return scanner.Err()
}

// parseBtreeWriteStats parses the btree_write_stats file.
// Format:
//
//	                   nr        size
//	initial:           19088     108k
//	init_next_bset:    6647      23.5k
//	cache_reclaim:     0         0
//	journal_reclaim:   541080    405
//	interior:          16788     354
func (c *bcachefsCollector) parseBtreeWriteStats(ch chan<- prometheus.Metric, fsPath, uuid string) error {
	const subsystem = "bcachefs"

	file, err := os.Open(filepath.Join(fsPath, "btree_write_stats"))
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip header line
		if lineNum == 1 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		writeType := strings.TrimSuffix(fields[0], ":")
		nrStr := fields[1]
		sizeStr := fields[2]

		nr, err := strconv.ParseFloat(nrStr, 64)
		if err != nil {
			c.logger.Debug("failed to parse btree_write_stats nr", "type", writeType, "err", err)
			continue
		}

		size, err := parseHumanReadableBytes(sizeStr)
		if err != nil {
			c.logger.Debug("failed to parse btree_write_stats size", "type", writeType, "err", err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "btree_writes_total"),
				"Number of btree writes by type.",
				[]string{"uuid", "type"},
				nil,
			),
			prometheus.CounterValue,
			nr,
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
			size,
			uuid, writeType,
		)
	}

	return scanner.Err()
}
