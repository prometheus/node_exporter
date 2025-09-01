package collector

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type kernelModulesCollector struct {
	procModules *prometheus.Desc
	logger      *slog.Logger
}

func init() {
	registerCollector("kernelmodules", defaultDisabled, NewKernelModulesCollector)
}

// NewKernelModulesCollector returns a new Collector exposing kernel module information.
func NewKernelModulesCollector(logger *slog.Logger) (Collector, error) {
	return &kernelModulesCollector{
		procModules: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "kernel_modules", "info"),
			"Information about loaded kernel modules from /proc/modules. "+
				"State values indicate module status: "+
				"'Live' (module is fully loaded and functioning), "+
				"'Loading' (module load in progress), "+
				"'Unloading' (module removal in progress). "+
				"Size shows memory usage in bytes. "+
				"References count shows number of other modules or processes using this module.",
			[]string{
				"module",     // Module name
				"size",       // Memory usage in bytes
				"references", // Usage count
				"state",      // Module state
			},
			nil,
		),
		logger: logger,
	}, nil
}

// Update implements Collector and exposes kernel module metrics from /proc/modules.
func (c *kernelModulesCollector) Update(ch chan<- prometheus.Metric) error {
	file, err := os.Open("/proc/modules")
	if err != nil {
		return fmt.Errorf("failed to read /proc/modules: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 5 {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			c.procModules,
			prometheus.GaugeValue,
			1,
			parts[0], // module
			parts[1], // size
			parts[2], // references
			parts[4], // state
		)
	}

	return scanner.Err()
}
