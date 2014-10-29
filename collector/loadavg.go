// +build !noloadavg

package collector

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	procLoad = "/proc/loadavg"
)

var (
	load1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "load1",
		Help:      "1m load average.",
	})
)

type loadavgCollector struct {
	config Config
}

func init() {
	Factories["loadavg"] = NewLoadavgCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// load, seconds since last login and a list of tags as specified by config.
func NewLoadavgCollector(config Config) (Collector, error) {
	c := loadavgCollector{
		config: config,
	}
	return &c, nil
}

func (c *loadavgCollector) Update(ch chan<- prometheus.Metric) (err error) {
	load, err := getLoad1()
	if err != nil {
		return fmt.Errorf("Couldn't get load: %s", err)
	}
	glog.V(1).Infof("Set node_load: %f", load)
	load1.Set(load)
	load1.Collect(ch)
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
