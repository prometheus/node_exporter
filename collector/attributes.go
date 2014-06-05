// +build !noattributes

package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	attributes = prometheus.NewGauge()
)

type attributesCollector struct {
	registry prometheus.Registry
	config   Config
}

func init() {
	Factories["attributes"] = NewAttributesCollector
}

// Takes a config struct and prometheus registry and returns a new Collector exposing
// labels from the config.
func NewAttributesCollector(config Config, registry prometheus.Registry) (Collector, error) {
	c := attributesCollector{
		config:   config,
		registry: registry,
	}
	registry.Register(
		"node_attributes",
		"node_exporter attributes",
		prometheus.NilLabels,
		attributes,
	)
	return &c, nil
}

func (c *attributesCollector) Update() (updates int, err error) {
	glog.V(1).Info("Set node_attributes{%v}: 1", c.config.Attributes)
	attributes.Set(c.config.Attributes, 1)
	return updates, err
}
