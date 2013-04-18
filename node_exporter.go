package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang"
	"github.com/prometheus/client_golang/metrics"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	proto    = "tcp"
	procLoad = "/proc/loadavg"
)

var (
	verbose          = flag.Bool("verbose", false, "Verbose output.")
	listeningAddress = flag.String("listeningAddress", ":8080", "Address on which to expose JSON metrics.")
	metricsEndpoint  = flag.String("metricsEndpoint", "/metrics.json", "Path under which to expose JSON metrics.")
	configFile       = flag.String("config", "node_exporter.conf", "Config file.")
	scrapeInterval   = flag.Int("interval", 60, "Scrape interval.")

	loadAvg          = metrics.NewGauge()
	attributes       = metrics.NewGauge()
	lastSeen         = metrics.NewGauge()
)

type config struct {
	Attributes map[string]string `json:"attributes"`
}

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Couldn't get hostname: %s", err)
	}

	registry.DefaultRegistry.Register(
		"node_load",
		"node_exporter: system load.",
		map[string]string{"hostname": hostname},
		loadAvg,
	)

	registry.DefaultRegistry.Register(
		"node_last_login_seconds",
		"node_exporter: seconds since last login.",
		map[string]string{"hostname": hostname},
		lastSeen,
	)

	registry.DefaultRegistry.Register(
		"node_attributes",
		"node_exporter: system attributes.",
		map[string]string{"hostname": hostname},
		attributes,
	)
}

func debug(format string, a ...interface{}) {
	if *verbose {
		log.Printf(format, a...)
	}
}

func newConfig(filename string) (conf config, err error) {
	log.Printf("Reading config %s", filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &conf)
	return
}

func serveStatus() {
	exporter := registry.DefaultRegistry.Handler()

	http.Handle(*metricsEndpoint, exporter)
	http.ListenAndServe(*listeningAddress, nil)
}

// Takes a string, splits it, converts each element to int and returns them as new list.
// It will return an error in case any element isn't an int.
func splitToInts(str string, sep string) (ints []int, err error) {
	for _, part := range strings.Split(str, sep) {
		i, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("Could not split '%s' because %s is no int: %s", str, part, err)
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func main() {
	flag.Parse()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)
	configChan := make(chan config)
	go func() {
		for _ = range sig {
			config, err := newConfig(*configFile)
			if err != nil {
				log.Printf("Couldn't reload config: %s", err)
				continue
			}
			configChan <- config
		}
	}()

	conf, err := newConfig(*configFile)
	if err != nil {
		log.Fatalf("Couldn't read config: %s", err)
	}

	go serveStatus()

	tick := time.Tick(time.Duration(*scrapeInterval) * time.Second)
	for {
		select {
		case conf = <-configChan:
			log.Printf("Got new config")
		case <-tick:
			log.Printf("Starting new scrape interval")

			last, err := getSecondsSinceLastLogin()
			if err != nil {
				log.Printf("Couldn't get last seen: %s", err)
			} else {
				debug("last: %f", last)
				lastSeen.Set(nil, last)
			}

			load, err := getLoad()
			if err != nil {
				log.Printf("Couldn't get load: %s", err)
			} else {
				debug("load: %f", load)
				loadAvg.Set(nil, load)
			}

			debug("attributes: %s", conf.Attributes)
			attributes.Set(conf.Attributes, 1)

		}
	}
}

func getLoad() (float64, error) {
	data, err := ioutil.ReadFile(procLoad)
	if err != nil {
		return 0, err
	}
	parts := strings.Fields(string(data))
	load, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse load '%s': %s", parts[0], err)
	}
	return load, nil
}

func getSecondsSinceLastLogin() (float64, error) {
	who := exec.Command("who", "/var/log/wtmp", "-l", "-u", "-s")

	output, err := who.StdoutPipe()
	if err != nil {
		return 0, err
	}

	err = who.Start()
	if err != nil {
		return 0, err
	}

	reader := bufio.NewReader(output)
	var last time.Time
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if isPrefix {
			return 0, fmt.Errorf("line to long: %s(...)", line)
		}

		fields := strings.Fields(string(line))
		lastDate := fields[2]
		lastTime := fields[3]

		dateParts, err := splitToInts(lastDate, "-") // 2013-04-16
		if err != nil {
			return 0, err
		}

		timeParts, err := splitToInts(lastTime, ":") // 11:33
		if err != nil {
			return 0, err
		}

		last_t := time.Date(dateParts[0], time.Month(dateParts[1]), dateParts[2], timeParts[0], timeParts[1], 0, 0, time.UTC)
		last = last_t
	}
	err = who.Wait()
	if err != nil {
		return 0, err
	}

	return float64(time.Now().Sub(last).Seconds()), nil
}
