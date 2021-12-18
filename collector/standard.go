package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("standard.go", defaultDisabled, NewStandardGoCollector)
	registerCollector("standard.process", defaultDisabled, NewStandardProcessCollector)
}

type standardGoCollector struct {
	origin prometheus.Collector
	logger log.Logger
}

// NewStandardGoCollector creates standard go collector.
func NewStandardGoCollector(logger log.Logger) (Collector, error) {
	c := prometheus.NewGoCollector()
	return &standardGoCollector{origin: c}, nil
}

func (c *standardGoCollector) Update(ch chan<- prometheus.Metric) error {
	c.origin.Collect(ch)
	return nil
}

type standardProcessCollector struct {
	origin prometheus.Collector
}

// NewStandardProcessCollector creates standard process collector.
func NewStandardProcessCollector(logger log.Logger) (Collector, error) {
	c := prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{})
	return &standardProcessCollector{origin: c}, nil
}

func (c *standardProcessCollector) Update(ch chan<- prometheus.Metric) error {
	c.origin.Collect(ch)
	return nil
}
