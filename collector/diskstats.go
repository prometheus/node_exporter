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
	diskSubsystem = "disk"
)

var (
	ignoredDevices = flag.String("diskstatsIgnoredDevices", "^(ram|loop|(h|s|xv)d[a-z])\\d+$", "Regexp of devices to ignore for diskstats.")

	diskLabelNames = []string{"device"}

	// Docs from https://www.kernel.org/doc/Documentation/iostats.txt
	diskStatsMetrics = []prometheus.Collector{
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "reads_completed",
				Help:      "The total number of reads completed successfully.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "reads_merged",
				Help:      "The number of reads merged. See https://www.kernel.org/doc/Documentation/iostats.txt.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "sectors_read",
				Help:      "The total number of sectors read successfully.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "read_time_ms",
				Help:      "The total number of milliseconds spent by all reads.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "writes_completed",
				Help:      "The total number of writes completed successfully.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "writes_merged",
				Help:      "The number of writes merged. See https://www.kernel.org/doc/Documentation/iostats.txt.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "sectors_written",
				Help:      "The total number of sectors written successfully.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "write_time_ms",
				Help:      "This is the total number of milliseconds spent by all writes.",
			},
			diskLabelNames,
		),
		prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "io_now",
				Help:      "The number of I/Os currently in progress.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "io_time_ms",
				Help:      "Milliseconds spent doing I/Os.",
			},
			diskLabelNames,
		),
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: Namespace,
				Subsystem: diskSubsystem,
				Name:      "io_time_weighted",
				Help:      "The weighted # of milliseconds spent doing I/Os. See https://www.kernel.org/doc/Documentation/iostats.txt.",
			},
			diskLabelNames,
		),
	}
)

type diskstatsCollector struct {
	config                Config
	ignoredDevicesPattern *regexp.Regexp
}

func init() {
	Factories["diskstats"] = NewDiskstatsCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// disk device stats.
func NewDiskstatsCollector(config Config) (Collector, error) {
	c := diskstatsCollector{
		config:                config,
		ignoredDevicesPattern: regexp.MustCompile(*ignoredDevices),
	}

	for _, c := range diskStatsMetrics {
		if _, err := prometheus.RegisterOrGet(c); err != nil {
			return nil, err
		}
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
			counter, ok := diskStatsMetrics[k].(*prometheus.CounterVec)
			if ok {
				counter.WithLabelValues(dev).Set(v)
			} else {
				var gauge = diskStatsMetrics[k].(*prometheus.GaugeVec)
				gauge.WithLabelValues(dev).Set(v)
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
