// +build !nofilefd

package collector

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	procFileFDStat      = "/proc/sys/fs/file-nr"
	fileFDStatSubsystem = "filefd"
)

type fileFDStatCollector struct {
	metrics map[string]prometheus.Gauge
}

func init() {
	Factories[fileFDStatSubsystem] = NewFileFDStatCollector
}

// NewFileFDStatCollector returns a new Collector exposing file-nr stats.
func NewFileFDStatCollector() (Collector, error) {
	return &fileFDStatCollector{
		metrics: map[string]prometheus.Gauge{},
	}, nil
}

func (c *fileFDStatCollector) Update(ch chan<- prometheus.Metric) (err error) {
	fileFDStat, err := getFileFDStats(procFileFDStat)
	if err != nil {
		return fmt.Errorf("couldn't get file-nr: %s", err)
	}
	for name, value := range fileFDStat {
		if _, ok := c.metrics[name]; !ok {
			c.metrics[name] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: Namespace,
					Subsystem: fileFDStatSubsystem,
					Name:      name,
					Help:      fmt.Sprintf("filefd %s from %s.", name, procFileFDStat),
				},
			)
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in file-nr: %s", value, err)
			}
			c.metrics[name].Set(v)
		}
	}
	for _, m := range c.metrics {
		m.Collect(ch)
	}
	return err
}

func getFileFDStats(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseFileFDStats(file, fileName)
}

func parseFileFDStats(r io.Reader, fileName string) (map[string]string, error) {
	var scanner = bufio.NewScanner(r)
	scanner.Scan()
	// The file-nr proc file is separated by tabs, not spaces.
	line := strings.Split(string(scanner.Text()), "\u0009")
	var fileFDStat = map[string]string{}
	// The file-nr proc is only 1 line with 3 values.
	fileFDStat["allocated"] = line[0]
	// The second value is skipped as it will alwasy be zero in linux 2.6.
	fileFDStat["maximum"] = line[2]

	return fileFDStat, nil
}
