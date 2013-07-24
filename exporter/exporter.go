// Exporter is a prometheus exporter using multiple collectors to collect and export system metrics.
package exporter

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/exp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
)

var verbose = flag.Bool("verbose", false, "Verbose output.")

// Interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update() (n int, err error)

	// Returns the name of the collector
	Name() string
}

type config struct {
	Attributes       map[string]string `json:"attributes"`
	ListeningAddress string            `json:"listeningAddress"`
	ScrapeInterval   int               `json:"scrapeInterval"`
	Collectors       []string          `json:"collectors"`
}

func (e *exporter) loadConfig() (err error) {
	log.Printf("Reading config %s", e.configFile)
	bytes, err := ioutil.ReadFile(e.configFile)
	if err != nil {
		return
	}

	return json.Unmarshal(bytes, &e.config) // Make sure this is safe
}

type exporter struct {
	configFile       string
	listeningAddress string
	scrapeInterval   time.Duration
	scrapeDurations  prometheus.Histogram
	metricsUpdated   prometheus.Gauge
	config           config
	registry         prometheus.Registry
	collectors       []Collector
	MemProfile       string
}

// New takes the path to a config file and returns an exporter instance
func New(configFile string) (e exporter, err error) {
	registry := prometheus.NewRegistry()
	e = exporter{
		configFile:       configFile,
		scrapeDurations:  prometheus.NewDefaultHistogram(),
		metricsUpdated:   prometheus.NewGauge(),
		listeningAddress: ":8080",
		scrapeInterval:   60 * time.Second,
		registry:         registry,
	}

	err = e.loadConfig()
	if err != nil {
		return e, fmt.Errorf("Couldn't read config: %s", err)
	}

	cn, err := NewNativeCollector(e.config, e.registry)
	if err != nil {
		log.Fatalf("Couldn't attach collector: %s", err)
	}

	cg, err := NewGmondCollector(e.config, e.registry)
	if err != nil {
		log.Fatalf("Couldn't attach collector: %s", err)
	}

	e.collectors = []Collector{&cn, &cg}

	if e.config.ListeningAddress != "" {
		e.listeningAddress = e.config.ListeningAddress
	}
	if e.config.ScrapeInterval != 0 {
		e.scrapeInterval = time.Duration(e.config.ScrapeInterval) * time.Second
	}

	registry.Register("node_exporter_scrape_duration_seconds", "node_exporter: Duration of a scrape job.", prometheus.NilLabels, e.scrapeDurations)
	registry.Register("node_exporter_metrics_updated", "node_exporter: Number of metrics updated.", prometheus.NilLabels, e.metricsUpdated)

	return e, nil
}

func (e *exporter) serveStatus() {
	exp.Handle(prometheus.ExpositionResource, e.registry.Handler())
	http.ListenAndServe(e.listeningAddress, exp.DefaultCoarseMux)
}

func (e *exporter) Execute(c Collector) {
	begin := time.Now()
	updates, err := c.Update()
	duration := time.Since(begin)

	label := map[string]string{
		"collector": c.Name(),
	}
	if err != nil {
		log.Printf("ERROR: %s failed after %fs: %s", c.Name(), duration.Seconds(), err)
		label["result"] = "error"
	} else {
		log.Printf("OK: %s success after %fs.", c.Name(), duration.Seconds())
		label["result"] = "success"
	}
	e.scrapeDurations.Add(label, duration.Seconds())
	e.metricsUpdated.Set(label, float64(updates))
}

func (e *exporter) Loop() {
	sigHup := make(chan os.Signal)
	sigUsr1 := make(chan os.Signal)
	signal.Notify(sigHup, syscall.SIGHUP)
	signal.Notify(sigUsr1, syscall.SIGUSR1)

	go e.serveStatus()

	tick := time.Tick(e.scrapeInterval)
	for {
		select {
		case <-sigHup:
			err := e.loadConfig()
			if err != nil {
				log.Printf("Couldn't reload config: %s", err)
				continue
			}
			log.Printf("Got new config")
			tick = time.Tick(e.scrapeInterval)

		case <-tick:
			log.Printf("Starting new scrape interval")
			wg := sync.WaitGroup{}
			wg.Add(len(e.collectors))
			for _, c := range e.collectors {
				go func(c Collector) {
					e.Execute(c)
					wg.Done()
				}(c)
			}
			wg.Wait()

		case <-sigUsr1:
			log.Printf("got signal")
			if e.MemProfile != "" {
				log.Printf("Writing memory profile to %s", e.MemProfile)
				f, err := os.Create(e.MemProfile)
				if err != nil {
					log.Fatal(err)
				}
				pprof.WriteHeapProfile(f)
				f.Close()
			}
		}
	}
}
