// +build !nostat

package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#include <unistd.h>
*/
import "C"

const (
	procStat = "/proc/stat"
)

var (
	cpuMetrics         = prometheus.NewCounter()
	intrMetric         = prometheus.NewCounter()
	ctxtMetric         = prometheus.NewCounter()
	btimeMetric        = prometheus.NewGauge()
	forksMetric        = prometheus.NewCounter()
	procsRunningMetric = prometheus.NewGauge()
	procsBlockedMetric = prometheus.NewGauge()
)

type statCollector struct {
	registry prometheus.Registry
	config   Config
}

func init() {
	Factories["stat"] = NewStatCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// network device stats.
func NewStatCollector(config Config, registry prometheus.Registry) (Collector, error) {
	c := statCollector{
		config:   config,
		registry: registry,
	}
	registry.Register(
		"node_cpu",
		"Seconds the cpus spent in each mode.",
		prometheus.NilLabels,
		cpuMetrics,
	)
	registry.Register(
		"node_intr",
		"Total number of interrupts serviced",
		prometheus.NilLabels,
		intrMetric,
	)
	registry.Register(
		"node_context_switches",
		"Total number of context switches.",
		prometheus.NilLabels,
		ctxtMetric,
	)
	registry.Register(
		"node_forks",
		"Total number of forks.",
		prometheus.NilLabels,
		forksMetric,
	)
	registry.Register(
		"node_boot_time",
		"Node boot time, in unixtime.",
		prometheus.NilLabels,
		btimeMetric,
	)
	registry.Register(
		"node_procs_running",
		"Number of processes in runnable state.",
		prometheus.NilLabels,
		procsRunningMetric,
	)
	registry.Register(
		"node_procs_blocked",
		"Number of processes blocked waiting for I/O to complete.",
		prometheus.NilLabels,
		procsBlockedMetric,
	)
	return &c, nil
}

// Expose a variety of stats from /proc/stats.
func (c *statCollector) Update() (updates int, err error) {
	file, err := os.Open(procStat)
	if err != nil {
		return updates, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		switch {
		case strings.HasPrefix(parts[0], "cpu"):
			// Export only per-cpu stats, it can be aggregated up in prometheus.
			if parts[0] == "cpu" {
				break
			}
			// Only some of these may be present, depending on kernel version.
			cpuFields := []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest"}
			for i, v := range parts[1 : len(cpuFields)+1] {
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return updates, err
				}
                                // Convert from ticks to seconds
				value /= float64(C.sysconf(C._SC_CLK_TCK))
				cpuMetrics.Set(map[string]string{"cpu": parts[0], "mode": cpuFields[i]}, value)
			}
		case parts[0] == "intr":
			// Only expose the overall number, use the 'interrupts' collector for more detail.
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			intrMetric.Set(prometheus.NilLabels, value)
		case parts[0] == "ctxt":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			ctxtMetric.Set(prometheus.NilLabels, value)
		case parts[0] == "processes":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			forksMetric.Set(prometheus.NilLabels, value)
		case parts[0] == "btime":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			btimeMetric.Set(prometheus.NilLabels, value)
		case parts[0] == "procs_running":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			procsRunningMetric.Set(prometheus.NilLabels, value)
		case parts[0] == "procs_blocked":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			procsBlockedMetric.Set(prometheus.NilLabels, value)
		}
	}
	return updates, err
}
