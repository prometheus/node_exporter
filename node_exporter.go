package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/exp"
	"github.com/prometheus/node_exporter/collector"
)

var (
	configFile       = flag.String("config", "node_exporter.conf", "config file.")
	memProfile       = flag.String("memprofile", "", "write memory profile to this file")
	listeningAddress = flag.String("listen", ":8080", "address to listen on")
	interval         = flag.Duration("interval", 60*time.Second, "refresh interval")
	scrapeDurations  = prometheus.NewDefaultHistogram()
	metricsUpdated   = prometheus.NewGauge()
)

func main() {
	flag.Parse()
	registry := prometheus.NewRegistry()
	collectors, err := loadCollectors(*configFile, registry)
	if err != nil {
		log.Fatalf("Couldn't load config and collectors: %s", err)
	}

	registry.Register("node_exporter_scrape_duration_seconds", "node_exporter: Duration of a scrape job.", prometheus.NilLabels, scrapeDurations)
	registry.Register("node_exporter_metrics_updated", "node_exporter: Number of metrics updated.", prometheus.NilLabels, metricsUpdated)

	log.Printf("Registered collectors:")
	for _, c := range collectors {
		log.Print(" - ", c.Name())
	}

	sigHup := make(chan os.Signal)
	sigUsr1 := make(chan os.Signal)
	signal.Notify(sigHup, syscall.SIGHUP)
	signal.Notify(sigUsr1, syscall.SIGUSR1)

	go serveStatus(registry)

	log.Printf("Starting initial collection")
	collect(collectors)

	tick := time.Tick(*interval)
	for {
		select {
		case <-sigHup:
			collectors, err = loadCollectors(*configFile, registry)
			if err != nil {
				log.Fatalf("Couldn't load config and collectors: %s", err)
			}
			log.Printf("Reloaded collectors and config")
			tick = time.Tick(*interval)

		case <-tick:
			log.Printf("Starting new interval")
			collect(collectors)

		case <-sigUsr1:
			log.Printf("got signal")
			if *memProfile != "" {
				log.Printf("Writing memory profile to %s", *memProfile)
				f, err := os.Create(*memProfile)
				if err != nil {
					log.Fatal(err)
				}
				pprof.WriteHeapProfile(f)
				f.Close()
			}
		}
	}

}

func loadCollectors(file string, registry prometheus.Registry) ([]collector.Collector, error) {
	collectors := []collector.Collector{}
	config, err := getConfig(file)
	if err != nil {
		log.Fatalf("Couldn't read config %s: %s", file, err)
	}
	for _, fn := range collector.Factories {
		c, err := fn(*config, registry)
		if err != nil {
			return nil, err
		}
		collectors = append(collectors, c)
	}
	return collectors, nil
}

func getConfig(file string) (*collector.Config, error) {
	config := &collector.Config{}
	log.Printf("Reading config %s", *configFile)
	bytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}
	return config, json.Unmarshal(bytes, &config)
}

func serveStatus(registry prometheus.Registry) {
	exp.Handle(prometheus.ExpositionResource, registry.Handler())
	http.ListenAndServe(*listeningAddress, exp.DefaultCoarseMux)
}

func collect(collectors []collector.Collector) {
	wg := sync.WaitGroup{}
	wg.Add(len(collectors))
	for _, c := range collectors {
		go func(c collector.Collector) {
			Execute(c)
			wg.Done()
		}(c)
	}
	wg.Wait()
}

func Execute(c collector.Collector) {
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
	scrapeDurations.Add(label, duration.Seconds())
	metricsUpdated.Set(label, float64(updates))
}
