// +build !noconsul

package collector

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	consulSubsystem = "consul"
)

var (
	consulBinary = flag.String("consul.binary", "consul", "Location of the consul binary")

	consulIgnoredFields = map[string]struct{}{
		"arch":    struct{}{},
		"os":      struct{}{},
		"state":   struct{}{},
		"version": struct{}{},
	}
	consulIgnoredKeys = map[string]struct{}{
		"build": struct{}{},
	}
)

func init() {
	Factories["consul"] = NewConsulCollector
}

type consulCollector struct {
	metrics map[string]prometheus.Gauge
}

// NewConsulCollector returns the collector implementation to export Consul
// internals.
func NewConsulCollector(config Config) (Collector, error) {
	return &consulCollector{
		metrics: map[string]prometheus.Gauge{},
	}, nil
}

func (c *consulCollector) Update(metricc chan<- prometheus.Metric) error {
	stats, err := getConsulStats()
	if err != nil {
		return fmt.Errorf("consul info failed: %s", err)
	}

	for name, value := range stats {
		if _, ok := c.metrics[name]; !ok {
			c.metrics[name] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: consulSubsystem,
					Name:      name,
					Help:      fmt.Sprintf("%s from consul info", name),
				},
			)
		}

		c.metrics[name].Set(float64(value))
	}

	for _, m := range c.metrics {
		m.Collect(metricc)
	}

	return nil
}

type consulStats map[string]int64

func getConsulStats() (consulStats, error) {
	info := exec.Command(*consulBinary, "info")

	outPipe, err := info.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := info.Start(); err != nil {
		return nil, err
	}

	stats, err := parseConsulStats(outPipe)
	if err != nil {
		return nil, err
	}

	if err := info.Wait(); err != nil {
		return nil, err
	}

	return stats, nil
}

func parseConsulStats(r io.Reader) (consulStats, error) {
	var (
		s     = bufio.NewScanner(r)
		stats = consulStats{}

		key string
	)

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if strings.Contains(line, ":") {
			key = strings.TrimSuffix(line, ":")
		}

		if _, ok := consulIgnoredKeys[key]; ok {
			continue
		}

		if strings.Contains(line, "=") {
			var (
				parts        = strings.SplitN(line, "=", 2)
				field, value = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
				name         = strings.Join([]string{key, field}, "_")
			)

			if _, ok := consulIgnoredFields[field]; ok {
				continue
			}

			switch value {
			case "true":
				stats[name] = 1
			case "false":
				stats[name] = 0
			default:
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, err
				}
				stats[name] = i
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
