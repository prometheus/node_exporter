// +build !nostat

package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// #include <unistd.h>
import "C"

const (
	procStat = "/proc/stat"
)

var (
	cpuMetrics = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "cpu",
			Help:      "Seconds the cpus spent in each mode.",
		},
		[]string{"cpu", "mode"},
	)
	intrMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "intr",
		Help:      "Total number of interrupts serviced.",
	})
	ctxtMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "context_switches",
		Help:      "Total number of context switches.",
	})
	forksMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "forks",
		Help:      "Total number of forks.",
	})
	btimeMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "boot_time",
		Help:      "Node boot time, in unixtime.",
	})
	procsRunningMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "procs_running",
		Help:      "Number of processes in runnable state.",
	})
	procsBlockedMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "procs_blocked",
		Help:      "Number of processes blocked waiting for I/O to complete.",
	})
)

type statCollector struct {
	config Config
}

func init() {
	Factories["stat"] = NewStatCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// network device stats.
func NewStatCollector(config Config) (Collector, error) {
	c := statCollector{
		config: config,
	}
	if _, err := prometheus.RegisterOrGet(cpuMetrics); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(intrMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(ctxtMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(forksMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(btimeMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(procsRunningMetric); err != nil {
		return nil, err
	}
	if _, err := prometheus.RegisterOrGet(procsBlockedMetric); err != nil {
		return nil, err
	}
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
				cpuMetrics.With(prometheus.Labels{"cpu": parts[0], "mode": cpuFields[i]}).Set(value)
			}
		case parts[0] == "intr":
			// Only expose the overall number, use the 'interrupts' collector for more detail.
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			intrMetric.Set(value)
		case parts[0] == "ctxt":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			ctxtMetric.Set(value)
		case parts[0] == "processes":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			forksMetric.Set(value)
		case parts[0] == "btime":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			btimeMetric.Set(value)
		case parts[0] == "procs_running":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			procsRunningMetric.Set(value)
		case parts[0] == "procs_blocked":
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return updates, err
			}
			procsBlockedMetric.Set(value)
		}
	}
	return updates, err
}
