// Exporter is a prometheus exporter using multiple Factories to collect and export system metrics.
package collector

const Namespace = "node"

var Factories = make(map[string]func(Config) (Collector, error))

// Interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update() (n int, err error)
}

type Config struct {
	Attributes map[string]string `json:"attributes"`
}
