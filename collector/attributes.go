// +build !noattributes

package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

type attributesCollector struct {
	config Config
	metric *prometheus.GaugeVec
}

func init() {
	Factories["attributes"] = NewAttributesCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// labels from the config.
func NewAttributesCollector(config Config) (Collector, error) {
	labelNames := []string{}
	for l := range config.Attributes {
		labelNames = append(labelNames, l)
	}

	return &attributesCollector{
		config: config,
		metric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "attributes",
				Help:      "The node_exporter attributes.",
			},
			labelNames,
		),
	}, nil
}

func (c *attributesCollector) Update(ch chan<- prometheus.Metric) (err error) {
	glog.V(1).Info("Set node_attributes{%v}: 1", c.config.Attributes)
	c.metric.Reset()
	c.metric.With(c.config.Attributes).Set(1)
	c.metric.Collect(ch)
	return err
}
