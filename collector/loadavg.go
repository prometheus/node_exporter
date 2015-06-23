// +build !noloadavg

package collector

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

// #include <stdlib.h>
import "C"

type loadavgCollector struct {
	metric prometheus.Gauge
}

func init() {
	Factories["loadavg"] = NewLoadavgCollector
}

// load1 stat.
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
	var loadavg [1]C.double
	samples := C.getloadavg(&loadavg[0], 1)
	if samples > 0 {
		return float64(loadavg[0]), nil
	} else {
		return 0, errors.New("Failed to get load average!")
	}

}
