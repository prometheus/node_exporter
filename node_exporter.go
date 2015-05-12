package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sort"
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
	// set at build time
	Version = "0.0.0.dev"

	configFile        = flag.String("config.file", "", "Path to config file.")
	memProfile        = flag.String("debug.memprofile-file", "", "Write memory profile to this file upon receipt of SIGUSR1.")
	listenAddress     = flag.String("web.listen-address", ":9100", "Address on which to expose metrics and web interface.")
	metricsPath       = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	enabledCollectors = flag.String("collectors.enabled", "attributes,diskstats,filesystem,loadavg,meminfo,stat,textfile,time,netdev,netstat", "Comma-separated list of collectors to use.")
	printCollectors   = flag.Bool("collectors.print", false, "If true, print available collectors and exit.")
	authUser          = flag.String("auth.user", "", "Username for basic auth.")
	authPass          = flag.String("auth.pass", "", "Password for basic auth.")

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
)

// Implements Collector.
type NodeCollector struct {
	collectors map[string]collector.Collector
}

// Implements Collector.
func (n NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Implements Collector.
func (n NodeCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.collectors))
	for name, c := range n.collectors {
		go func(name string, c collector.Collector) {
			Execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	scrapeDurations.Collect(ch)
}

type basicAuthHandler struct {
	handler  http.HandlerFunc
	user     string
	password string
}

func (h *basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if !ok || password != h.password || user != h.user {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"metrics\"")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	h.handler(w, r)
	return
}

func Execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
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

func loadCollectors(file string) (map[string]collector.Collector, error) {
	collectors := map[string]collector.Collector{}
	config := &collector.Config{}
	if file != "" {
		var err error
		config, err = getConfig(file)
		if err != nil {
			return nil, fmt.Errorf("couldn't read config %s: %s", file, err)
		}
	}
	for _, name := range strings.Split(*enabledCollectors, ",") {
		fn, ok := collector.Factories[name]
		if !ok {
			return nil, fmt.Errorf("collector '%s' not available", name)
		}
		c, err := fn(*config)
		if err != nil {
			return nil, err
		}
		collectors[name] = c
	}
	return collectors, nil
}

func main() {
	flag.Parse()

	if *printCollectors {
		collectorNames := make(sort.StringSlice, 0, len(collector.Factories))
		for n, _ := range collector.Factories {
			collectorNames = append(collectorNames, n)
		}
		collectorNames.Sort()
		fmt.Printf("Available collectors:\n")
		for _, n := range collectorNames {
			fmt.Printf(" - %s\n", n)
		}
		return
	}
	collectors, err := loadCollectors(*configFile)
	if err != nil {
		glog.Fatalf("Couldn't load config and collectors: %s", err)
	}

	glog.Infof("Enabled collectors:")
	for n, _ := range collectors {
		glog.Infof(" - %s", n)
	}

	nodeCollector := NodeCollector{collectors: collectors}
	prometheus.MustRegister(nodeCollector)

	sigUsr1 := make(chan os.Signal)
	signal.Notify(sigUsr1, syscall.SIGUSR1)

	handler := prometheus.Handler()
	if *authUser != "" || *authPass != "" {
		if *authUser == "" || *authPass == "" {
			glog.Fatal("You need to specify -auth.user and -auth.pass to enable basic auth")
		}
		handler = &basicAuthHandler{
			handler:  prometheus.Handler().ServeHTTP,
			user:     *authUser,
			password: *authPass,
		}
	}

	http.Handle(*metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Node Exporter</title></head>
			<body>
			<h1>Node Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	glog.Infof("Starting node_exporter v%s at %s", Version, *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		glog.Fatal(err)
	}
}
