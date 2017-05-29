// Copyright 2017 The Prometheus Authors
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

// +build !nobcache

package collector

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	// https://godoc.org/github.com/prometheus/client_golang/prometheus
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	bcacheUuidRE = regexp.MustCompile(`.*bcache/(.*-.*)`)
	bdevNoRE     = regexp.MustCompile(`.*(bdev[0-9]*)`)
	cacheNoRE    = regexp.MustCompile(`.*(cache[0-9]*)`)
)

func init() {
	Factories["bcache"] = NewBcacheCollector
}

type bcacheCollector struct {
	descs map[string]typedDesc
}

// NewBcacheCollector returns a newly allocated bcacheCollector.
// It exposes a number of Linux bcache statistics.
func NewBcacheCollector() (Collector, error) {
	return &bcacheCollector{
		// /sys/fs/bcache/<uuid>/
		descs: map[string]typedDesc{
			"average_key_size": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "average_key_size_sectors"),
					"Average data per key in the btree (sectors).",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			"btree_cache_size": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "btree_cache_size_bytes"),
					"Amount of memory currently used by the btree cache.",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			"cache_available_percent": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "cache_available_percent"),
					"Percentage of cache device without dirty data, useable for writeback (may contain clean cached data).",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			"congested": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "congested"),
					"Congestion.",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			"root_usage_percent": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "root_usage_percent"),
					"Percentage of the root btree node in use (tree depth increases if too high).",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			"tree_depth": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "tree_depth"),
					"Depth of the btree.",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			// /sys/fs/bcache/<uuid>/internal/
			"active_journal_entries": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "active_journal_entries"),
					"Number of journal entries that are newer than the index.",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			"btree_nodes": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "btree_nodes"),
					"Depth of the btree.",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			"btree_read_average_duration_us": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "btree_read_average_duration_us"),
					"Average btree read duration.",
					[]string{"uuid"}, nil,
				), prometheus.GaugeValue,
			},
			"cache_read_races": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "cache_read_races"),
					"Counts instances where while data was being read from the cache, the bucket was reused and invalidated - i.e. where the pointer was stale after the read completed.",
					[]string{"uuid"}, nil,
				), prometheus.CounterValue,
			},
			// /sys/fs/bcache/{,<uuid>}/stats_five_minute/
			"bypassed": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "bypassed_bytes_5min"),
					"Amount of IO (both reads and writes) that has bypassed the cache in 5 minutes.",
					[]string{"uuid", "bdev_no"}, nil,
				), prometheus.GaugeValue,
			},
			"cache_hits": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "cache_hits_5min"),
					"Hits counted per individual IO as bcache sees them in 5 minutes.",
					[]string{"uuid", "bdev_no"}, nil,
				), prometheus.GaugeValue,
			},
			"cache_misses": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "cache_misses_5min"),
					"Misses counted per individual IO as bcache sees them in 5 minutes.",
					[]string{"uuid", "bdev_no"}, nil,
				), prometheus.GaugeValue,
			},
			"cache_bypass_hits": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "cache_bypass_hits_5min"),
					"Hits for IO intended to skip the cache in 5 minutes.",
					[]string{"uuid", "bdev_no"}, nil,
				), prometheus.GaugeValue,
			},
			"cache_bypass_misses": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "cache_bypass_misses_5min"),
					"Misses for IO intended to skip the cache in 5 minutes.",
					[]string{"uuid", "bdev_no"}, nil,
				), prometheus.GaugeValue,
			},
			"cache_miss_collisions": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "cache_miss_collisions_5min"),
					"Instances where data insertion from cache miss raced with write (data already present) in 5 minutes.",
					[]string{"uuid", "bdev_no"}, nil,
				), prometheus.GaugeValue,
			},
			"cache_readaheads": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "cache_readaheads_5min"),
					"Count of times readahead occurred in 5 minutes.",
					[]string{"uuid", "bdev_no"}, nil,
				), prometheus.GaugeValue,
			},
			// /sys/fs/bcache/<uuid>/<bdev_num>/
			"dirty_data": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "dirty_data_bytes"),
					"Amount of dirty data for this backing device in the cache.",
					[]string{"uuid", "bdev_no"}, nil,
				), prometheus.GaugeValue,
			},
			// /sys/fs/bcache/<uuid>/<cache_num>/
			"io_errors": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "io_errors"),
					"Number of errors that have occurred, decayed by io_error_halflife.",
					[]string{"uuid", "cache_num"}, nil,
				), prometheus.GaugeValue,
			},
			"metadata_written": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "metadata_written_bytes_total"),
					"Sum of all non data writes (btree writes and all other metadata).",
					[]string{"uuid", "cache_num"}, nil,
				), prometheus.CounterValue,
			},
			"written": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "written_bytes_total"),
					"Sum of all data that has been written to the cache.",
					[]string{"uuid", "cache_num"}, nil,
				), prometheus.CounterValue,
			},
			// /sys/fs/bcache/<uuid>/<cache_num>/priority_stats
			"priority_stats_unused_percent": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "priority_stats_unused_percent"),
					"The percentage of the cache that doesn't contain any data.",
					[]string{"uuid", "cache_num"}, nil,
				), prometheus.GaugeValue,
			},
			"priority_stats_metadata_percent": {
				prometheus.NewDesc(
					prometheus.BuildFQName(Namespace, "bcache", "priority_stats_metadate_percent"),
					"Bcache's metadata overhead.",
					[]string{"uuid", "cache_num"}, nil,
				), prometheus.GaugeValue,
			},
		},
	}, nil
}

// Update reads and exposes bcache stats.
// It implements the Collector interface.
func (c *bcacheCollector) Update(ch chan<- prometheus.Metric) error {

	bcacheUuidPaths, err := filepath.Glob(sysFilePath("fs/bcache/*-*"))
	if err != nil {
		return err
	}
	if len(bcacheUuidPaths) == 0 {
		log.Debugf("No bcache UUIDs found. Skipping.")
		return nil
	}

	for _, uuidPath := range bcacheUuidPaths {
		uuidMatch := bcacheUuidRE.FindStringSubmatch(uuidPath)
		if uuidMatch == nil {
			return fmt.Errorf("no UUID in %s", uuidPath)
		}
		bcacheUuid := uuidMatch[1]

		// stats
		files := []string{
			"average_key_size", "btree_cache_size",
			"cache_available_percent", "congested",
			"root_usage_percent", "tree_depth",
		}
		stats, err := getStats(uuidPath, files)
		if err != nil {
			return err
		}
		for metric, val := range stats {
			ch <- prometheus.MustNewConstMetric(c.descs[metric].desc, c.descs[metric].valueType, val, bcacheUuid)
		}

		// internal stats
		files = []string{
			"active_journal_entries", "btree_nodes",
			"btree_read_average_duration_us", "cache_read_races",
		}
		internalPath := path.Join(uuidPath, "internal")
		stats, err = getStats(internalPath, files)
		if err != nil {
			return err
		}
		for metric, val := range stats {
			ch <- prometheus.MustNewConstMetric(c.descs[metric].desc, c.descs[metric].valueType, val, bcacheUuid)
		}

		// bdev stats
		files = []string{
			"cache_hits", "cache_misses", "cache_bypass_hits",
			"cache_bypass_misses", "cache_miss_collisions",
			"cache_readaheads", "bypassed",
		}
		reg := path.Join(uuidPath, "bdev[0-9]*")
		bdevDirs, err := filepath.Glob(reg)
		if err != nil {
			return err
		}
		bdevDirs = append(bdevDirs, uuidPath)
		for _, bdevDir := range bdevDirs {
			// Label bdev_all for assumed path
			// /sys/fs/bcache/<uuid>/stats_five_minute/
			var bdevLabel = "bdev_all"

			// Conditionally update label value for
			// /sys/fs/bcache/<uuid>/<bdev_num>/stats_five_minutes/
			bdevMatch := bdevNoRE.FindStringSubmatch(bdevDir)
			if bdevMatch != nil {
				// Actually, it is
				bdevLabel = bdevMatch[1]
			}

			subDir := path.Join(bdevDir, "stats_five_minute")

			stats, err := getStats(subDir, files)
			if err != nil {
				return err
			}
			for metric, val := range stats {
				ch <- prometheus.MustNewConstMetric(c.descs[metric].desc, c.descs[metric].valueType, val, bcacheUuid, bdevLabel)
			}
		}

		// cache stats
		reg = path.Join(uuidPath, "cache[0-9]*")
		cacheDirs, err := filepath.Glob(reg)
		if err != nil {
			return err
		}
		for _, cacheDir := range cacheDirs {
			cacheMatch := cacheNoRE.FindStringSubmatch(cacheDir)
			if cacheMatch == nil {
				return fmt.Errorf("invalid path: %s", cacheDir)
			}
			cacheLabel := cacheMatch[1]
			stats, err := getCacheStats(cacheDir)
			if err != nil {
				return err
			}
			for metric, val := range stats {
				ch <- prometheus.MustNewConstMetric(c.descs[metric].desc, c.descs[metric].valueType, val, bcacheUuid, cacheLabel)

			}
		}

	}
	return nil
}

// getStats collects data from sysfs files
func getStats(dir string, files []string) (map[string]float64, error) {
	var (
		stats = map[string]float64{}
	)

	for _, metric := range files {
		sysPath := path.Join(dir, metric)
		byt, err := ioutil.ReadFile(sysPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read: %s", sysPath)
		}
		// Remove trailing newline
		byt = byt[:len(byt)-1]
		stats[metric] = dehumanize(byt)
	}
	return stats, nil
}

// getCacheStats collects data from the <cache_num> directories
func getCacheStats(dir string) (map[string]float64, error) {
	var (
		stats = map[string]float64{}
		files = []string{
			"io_errors", "metadata_written", "priority_stats",
			"written",
		}
	)

	for _, fi := range files {
		sysPath := path.Join(dir, fi)
		if fi == "priority_stats" {
			file, err := os.Open(sysPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read: %s", sysPath)

			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				key, val, err := parsePriorityStats(scanner.Text())
				if err != nil {
					return nil, fmt.Errorf("failed to parse: %s (%s)", sysPath, err)
				}
				if len(key) > 0 {
					stats[key] = val
				}
			}
			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("failed to parse: %s (%s)", sysPath, err)
			}
		} else {
			byt, err := ioutil.ReadFile(sysPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read: %s", sysPath)
			}
			// Remove trailing newline
			byt = byt[:len(byt)-1]
			stats[fi] = dehumanize(byt)
		}
	}

	return stats, nil
}

const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB
	ZiB
	YiB
)

// dehumanize converts human-readable byte slice into float64
func dehumanize(hbytes []byte) float64 {
	lastByte := hbytes[len(hbytes)-1]
	mul := float64(1)
	if lastByte > 57 {
		// beyond range of ASCII digits, must be a multiplier
		hbytes = hbytes[:len(hbytes)-1]
		multipliers := map[rune]float64{
			// Source for conversion rules:
			// linux-kernel/drivers/md/bcache/util.c:bch_hprint()
			'k': KiB,
			'M': MiB,
			'G': GiB,
			'T': TiB,
			'P': PiB,
			'E': EiB,
			'Z': ZiB,
			'Y': YiB,
		}
		mul = float64(multipliers[rune(lastByte)])
	}
	base, _ := strconv.ParseFloat(string(hbytes), 64)
	res := base * mul
	return res
}

// parsePriorityStats parses lines from the priority_stats file
func parsePriorityStats(line string) (key string, value float64, err error) {
	switch {
	case strings.HasPrefix(line, "Unused:"):
		key = "priority_stats_unused_percent"
		fields := strings.Fields(line)
		rawValue := fields[len(fields)-1]
		valueStr := strings.TrimSuffix(rawValue, "%")
		value, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return "", -1, err
		}
	case strings.HasPrefix(line, "Metadata:"):
		key = "priority_stats_metadata_percent"
		fields := strings.Fields(line)
		rawValue := fields[len(fields)-1]
		valueStr := strings.TrimSuffix(rawValue, "%")
		value, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return "", -1, err
		}
	default:
		key = ""
		value = 0
	}
	return key, value, nil
}
