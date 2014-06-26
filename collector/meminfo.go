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

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	procMemInfo      = "/proc/meminfo"
	memInfoSubsystem = "memory"
)

var (
	memInfoMetrics = map[string]prometheus.Gauge{}
)

type meminfoCollector struct {
	config Config
}

func init() {
	Factories["meminfo"] = NewMeminfoCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// memory stats.
func NewMeminfoCollector(config Config) (Collector, error) {
	c := meminfoCollector{
		config: config,
	}
	return &c, nil
}

func (c *meminfoCollector) Update() (updates int, err error) {
	memInfo, err := getMemInfo()
	if err != nil {
		return updates, fmt.Errorf("Couldn't get meminfo: %s", err)
	}
	glog.V(1).Infof("Set node_mem: %#v", memInfo)
	for k, v := range memInfo {
		if _, ok := memInfoMetrics[k]; !ok {
			gauge := prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Subsystem: memInfoSubsystem,
				Name:      k,
				Help:      k + " from /proc/meminfo.",
			})
			memInfoMetrics[k] = prometheus.MustRegisterOrGet(gauge).(prometheus.Gauge)
		}
		updates++
		memInfoMetrics[k].Set(v)
	}
	return updates, err
}

func getMemInfo() (map[string]float64, error) {
	file, err := os.Open(procMemInfo)
	if err != nil {
		return nil, err
	}
	return parseMemInfo(file)
}

func parseMemInfo(r io.ReadCloser) (map[string]float64, error) {
	defer r.Close()
	memInfo := map[string]float64{}
	scanner := bufio.NewScanner(r)
	re := regexp.MustCompile("\\((.*)\\)")
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
