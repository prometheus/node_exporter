// +build !nonative

package collector

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	procDiskStats = "/proc/diskstats"
)

type diskStat struct {
	name          string
	metric        prometheus.Metric
	documentation string
}

var (
	ignoredDevices = flag.String("diskstatsIgnoredDevices", "^(ram|loop|[hs]d[a-z])\\d+$", "Regexp of devices to ignore for diskstats.")

	// Docs from https://www.kernel.org/doc/Documentation/iostats.txt
	diskStatsMetrics = []diskStat{
		{"reads_completed", prometheus.NewCounter(), "The total number of reads completed successfully."},
		{"reads_merged", prometheus.NewCounter(), "The number of reads merged. See https://www.kernel.org/doc/Documentation/iostats.txt"},
		{"sectors_read", prometheus.NewCounter(), "The total number of sectors read successfully."},
		{"read_time_ms", prometheus.NewCounter(), "the total number of milliseconds spent by all reads."},
		{"writes_completed", prometheus.NewCounter(), "The total number of writes completed successfully."},
		{"writes_merged", prometheus.NewCounter(), "The number of writes merged. See https://www.kernel.org/doc/Documentation/iostats.txt"},
		{"sectors_written", prometheus.NewCounter(), "The total number of sectors written successfully."},
		{"write_time_ms", prometheus.NewCounter(), "This is the total number of milliseconds spent by all writes."},
		{"io_now", prometheus.NewGauge(), "The number of I/Os currently in progress."},
		{"io_time_ms", prometheus.NewCounter(), "Milliseconds spent doing I/Os."},
		{"io_time_weighted", prometheus.NewCounter(), "The weighted # of milliseconds spent doing I/Os. See https://www.kernel.org/doc/Documentation/iostats.txt"},
	}
)

type diskstatsCollector struct {
	registry              prometheus.Registry
	config                Config
	ignoredDevicesPattern *regexp.Regexp
}

func init() {
	Factories["diskstats"] = NewDiskstatsCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// disk device stats.
func NewDiskstatsCollector(config Config, registry prometheus.Registry) (Collector, error) {
	c := diskstatsCollector{
		config:                config,
		registry:              registry,
		ignoredDevicesPattern: regexp.MustCompile(*ignoredDevices),
	}

	for _, v := range diskStatsMetrics {
		registry.Register(
			"node_disk_"+v.name,
			v.documentation,
			prometheus.NilLabels,
			v.metric,
		)
	}
	return &c, nil
}

func (c *diskstatsCollector) Update() (updates int, err error) {
	diskStats, err := getDiskStats()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get diskstats: %s", err)
	}
	for dev, stats := range diskStats {
		if c.ignoredDevicesPattern.MatchString(dev) {
			glog.V(1).Infof("Ignoring device: %s", dev)
			continue
		}
		for k, value := range stats {
			updates++
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return updates, fmt.Errorf("Invalid value %s in diskstats: %s", value, err)
			}
			labels := map[string]string{"device": dev}
			counter, ok := diskStatsMetrics[k].metric.(prometheus.Counter)
			if ok {
				counter.Set(labels, v)
			} else {
				var gauge = diskStatsMetrics[k].metric.(prometheus.Gauge)
				gauge.Set(labels, v)
			}
		}
	}
	return updates, err
}

func getDiskStats() (map[string]map[int]string, error) {
	file, err := os.Open(procDiskStats)
	if err != nil {
		return nil, err
	}
	return parseDiskStats(file)
}

func parseDiskStats(r io.ReadCloser) (map[string]map[int]string, error) {
	defer r.Close()
	diskStats := map[string]map[int]string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Fields(string(scanner.Text()))
		if len(parts) != len(diskStatsMetrics)+3 { // we strip major, minor and dev
			return nil, fmt.Errorf("Invalid line in %s: %s", procDiskStats, scanner.Text())
		}
		dev := parts[2]
		diskStats[dev] = map[int]string{}
		for i, v := range parts[3:] {
			diskStats[dev][i] = v
		}
	}
	return diskStats, nil
}
