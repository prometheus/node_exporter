// +build !noloadavg

package collector

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	procLoad = "/proc/loadavg"
)

type loadavgCollector struct {
	metric prometheus.Gauge
}

func init() {
	Factories["loadavg"] = NewLoadavgCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewLoadavgCollector() (Collector, error) {
	return &loadavgCollector{
		metric: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "load1",
			Help:      "1m load average.",
		}),
	}, nil
}

func (c *loadavgCollector) Update(ch chan<- prometheus.Metric) (err error) {
	load, err := getLoad1()
	if err != nil {
		return fmt.Errorf("Couldn't get load: %s", err)
	}
	log.Debugf("Set node_load: %f", load)
	c.metric.Set(load)
	c.metric.Collect(ch)
	return err
}

func getLoad1() (float64, error) {
	data, err := ioutil.ReadFile(procLoad)
	if err != nil {
		return 0, err
	}
	return parseLoad(string(data))
}

func parseLoad(data string) (float64, error) {
	parts := strings.Fields(data)
	load, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse load '%s': %s", parts[0], err)
	}
	return load, nil
}
