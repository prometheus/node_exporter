package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/exp"
	"github.com/prometheus/node_exporter/collector"
)

var (
	configFile        = flag.String("config", "node_exporter.conf", "config file.")
	memProfile        = flag.String("memprofile", "", "write memory profile to this file")
	listeningAddress  = flag.String("listen", ":8080", "address to listen on")
	enabledCollectors = flag.String("enabledCollectors", "attributes,diskstats,filesystem,loadavg,meminfo,stat,netdev", "comma-seperated list of collectors to use")
	printCollectors   = flag.Bool("printCollectors", false, "If true, print available collectors and exit")
	interval          = flag.Duration("interval", 60*time.Second, "refresh interval")
	scrapeDurations   = prometheus.NewDefaultHistogram()
	metricsUpdated    = prometheus.NewGauge()
)

func main() {
	flag.Parse()
	if *printCollectors {
		fmt.Printf("Available collectors:\n")
		for n, _ := range collector.Factories {
			fmt.Printf(" - %s\n", n)
		}
		return
	}
	registry := prometheus.NewRegistry()
	collectors, err := loadCollectors(*configFile, registry)
	if err != nil {
		log.Fatalf("Couldn't load config and collectors: %s", err)
	}

	registry.Register("node_exporter_scrape_duration_seconds", "node_exporter: Duration of a scrape job.", prometheus.NilLabels, scrapeDurations)
	registry.Register("node_exporter_metrics_updated", "node_exporter: Number of metrics updated.", prometheus.NilLabels, metricsUpdated)

	glog.Infof("Enabled collectors:")
	for n, _ := range collectors {
		glog.Infof(" - %s", n)
	}

	sigHup := make(chan os.Signal)
	sigUsr1 := make(chan os.Signal)
	signal.Notify(sigHup, syscall.SIGHUP)
	signal.Notify(sigUsr1, syscall.SIGUSR1)

	go serveStatus(registry)

	glog.Infof("Starting initial collection")
	collect(collectors)

	tick := time.Tick(*interval)
	for {
		select {
		case <-sigHup:
			collectors, err = loadCollectors(*configFile, registry)
			if err != nil {
				log.Fatalf("Couldn't load config and collectors: %s", err)
			}
			glog.Infof("Reloaded collectors and config")
			tick = time.Tick(*interval)

		case <-tick:
			glog.Infof("Starting new interval")
			collect(collectors)

		case <-sigUsr1:
			glog.Infof("got signal")
			if *memProfile != "" {
				glog.Infof("Writing memory profile to %s", *memProfile)
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

func loadCollectors(file string, registry prometheus.Registry) (map[string]collector.Collector, error) {
	collectors := map[string]collector.Collector{}
	config, err := getConfig(file)
	if err != nil {
		log.Fatalf("Couldn't read config %s: %s", file, err)
	}
	for _, name := range strings.Split(*enabledCollectors, ",") {
		fn, ok := collector.Factories[name]
		if !ok {
			log.Fatalf("Collector '%s' not available", name)
		}
		c, err := fn(*config, registry)
		if err != nil {
			return nil, err
		}
		collectors[name] = c
	}
	return collectors, nil
}

func getConfig(file string) (*collector.Config, error) {
	config := &collector.Config{}
	glog.Infof("Reading config %s", *configFile)
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

func collect(collectors map[string]collector.Collector) {
	wg := sync.WaitGroup{}
	wg.Add(len(collectors))
	for n, c := range collectors {
		go func(n string, c collector.Collector) {
			Execute(n, c)
			wg.Done()
		}(n, c)
	}
	wg.Wait()
}

func Execute(name string, c collector.Collector) {
	begin := time.Now()
	updates, err := c.Update()
	duration := time.Since(begin)

	label := map[string]string{
		"collector": name,
	}
	if err != nil {
		glog.Infof("ERROR: %s failed after %fs: %s", name, duration.Seconds(), err)
		label["result"] = "error"
	} else {
		glog.Infof("OK: %s success after %fs.", name, duration.Seconds())
		label["result"] = "success"
	}
	scrapeDurations.Add(label, duration.Seconds())
	metricsUpdated.Set(label, float64(updates))
}
