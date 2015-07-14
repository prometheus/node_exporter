// +build !nomeminfo

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	procMemInfo      = "/proc/meminfo"
	memInfoSubsystem = "memory"
)

type meminfoCollector struct {
	metrics map[string]prometheus.Gauge
}

func init() {
	Factories["meminfo"] = NewMeminfoCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// memory stats.
func NewMeminfoCollector() (Collector, error) {
	return &meminfoCollector{
		metrics: map[string]prometheus.Gauge{},
	}, nil
}

func (c *meminfoCollector) Update(ch chan<- prometheus.Metric) (err error) {
	memInfo, err := getMemInfo()
	if err != nil {
		return fmt.Errorf("Couldn't get meminfo: %s", err)
	}
	log.Debugf("Set node_mem: %#v", memInfo)
	for k, v := range memInfo {
		if _, ok := c.metrics[k]; !ok {
			c.metrics[k] = prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: memInfoSubsystem,
				Name:      k,
				Help:      k + " from /proc/meminfo.",
			})
		}
		c.metrics[k].Set(v)
		c.metrics[k].Collect(ch)
	}
	return err
}

func getMemInfo() (map[string]float64, error) {
	file, err := os.Open(procMemInfo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseMemInfo(file)
}

func parseMemInfo(r io.Reader) (map[string]float64, error) {
	var (
		memInfo = map[string]float64{}
		scanner = bufio.NewScanner(r)
		re      = regexp.MustCompile("\\((.*)\\)")
	)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(string(line))
		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid value in meminfo: %s", err)
		}
		switch len(parts) {
		case 2: // no unit
		case 3: // has unit, we presume kB
			fv *= 1024
		default:
			return nil, fmt.Errorf("Invalid line in %s: %s", procMemInfo, line)
		}
		key := parts[0][:len(parts[0])-1] // remove trailing : from key
		// Active(anon) -> Active_anon
		key = re.ReplaceAllString(key, "_${1}")
		memInfo[key] = fv
	}

	return memInfo, nil
}
