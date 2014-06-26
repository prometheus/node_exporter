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
	"github.com/prometheus/node_exporter/collector"
)

const subsystem = "exporter"

var (
	configFile        = flag.String("config", "node_exporter.conf", "config file.")
	memProfile        = flag.String("memprofile", "", "write memory profile to this file")
	listeningAddress  = flag.String("listen", ":8080", "address to listen on")
	enabledCollectors = flag.String("enabledCollectors", "attributes,diskstats,filesystem,loadavg,meminfo,stat,netdev", "comma-seperated list of collectors to use")
	printCollectors   = flag.Bool("printCollectors", false, "If true, print available collectors and exit")
	interval          = flag.Duration("interval", 60*time.Second, "refresh interval")

	collectorLabelNames = []string{"collector", "result"}

	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: subsystem,
			Name:      "scrape_duration_seconds",
			Help:      "node_exporter: Duration of a scrape job.",
		},
		collectorLabelNames,
	)
	metricsUpdated = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: collector.Namespace,
			Subsystem: subsystem,
			Name:      "metrics_updated",
			Help:      "node_exporter: Number of metrics updated.",
		},
		collectorLabelNames,
	)
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
	collectors, err := loadCollectors(*configFile)
	if err != nil {
		log.Fatalf("Couldn't load config and collectors: %s", err)
	}

	prometheus.MustRegister(scrapeDurations)
	prometheus.MustRegister(metricsUpdated)

	glog.Infof("Enabled collectors:")
	for n, _ := range collectors {
		glog.Infof(" - %s", n)
	}

	sigHup := make(chan os.Signal)
	sigUsr1 := make(chan os.Signal)
	signal.Notify(sigHup, syscall.SIGHUP)
	signal.Notify(sigUsr1, syscall.SIGUSR1)

	go serveStatus()

	glog.Infof("Starting initial collection")
	collect(collectors)

	tick := time.Tick(*interval)
	for {
		select {
		case <-sigHup:
			collectors, err = loadCollectors(*configFile)
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

func loadCollectors(file string) (map[string]collector.Collector, error) {
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
		c, err := fn(*config)
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

func serveStatus() {
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(*listeningAddress, nil)
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
	var result string

	if err != nil {
		glog.Infof("ERROR: %s failed after %fs: %s", name, duration.Seconds(), err)
		result = "error"
	} else {
		glog.Infof("OK: %s success after %fs.", name, duration.Seconds())
		result = "success"
	}
	scrapeDurations.WithLabelValues(name, result).Observe(duration.Seconds())
	metricsUpdated.WithLabelValues(name, result).Set(float64(updates))
}
