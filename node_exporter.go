// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	slog "log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/prometheus/node_exporter/collector"
)

const (
	defaultCollectors = "conntrack,cpu,diskstats,entropy,filefd,filesystem,loadavg,mdadm,meminfo,netdev,netstat,sockstat,stat,textfile,time,uname,vmstat"
)

var (
	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: "exporter",
			Name:      "scrape_duration_seconds",
			Help:      "node_exporter: Duration of a scrape job.",
		},
		[]string{"collector", "result"},
	)
)

// NodeCollector implements the prometheus.Collector interface.
type NodeCollector struct {
	collectors map[string]collector.Collector
}

// Describe implements the prometheus.Collector interface.
func (n NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (n NodeCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.collectors))
	for name, c := range n.collectors {
		go func(name string, c collector.Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	scrapeDurations.Collect(ch)
}

func filterAvailableCollectors(collectors string) string {
	availableCollectors := make([]string, 0)
	for _, c := range strings.Split(collectors, ",") {
		_, ok := collector.Factories[c]
		if ok {
			availableCollectors = append(availableCollectors, c)
		}
	}
	return strings.Join(availableCollectors, ",")
}

func execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var result string

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		result = "error"
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		result = "success"
	}
	scrapeDurations.WithLabelValues(name, result).Observe(duration.Seconds())
}

func loadCollectors(list string) (map[string]collector.Collector, error) {
	collectors := map[string]collector.Collector{}
	for _, name := range strings.Split(list, ",") {
		fn, ok := collector.Factories[name]
		if !ok {
			return nil, fmt.Errorf("collector '%s' not available", name)
		}
		c, err := fn()
		if err != nil {
			return nil, err
		}
		collectors[name] = c
	}
	return collectors, nil
}

func init() {
	prometheus.MustRegister(version.NewCollector("node_exporter"))
}

func main() {
	var (
		showVersion       = flag.Bool("version", false, "Print version information.")
		listenAddress     = flag.String("web.listen-address", ":9100", "Address on which to expose metrics and web interface.")
		metricsPath       = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		tlsCert           = flag.String("web.tls-cert", "", "Path to PEM file that conains the certificate (and opionally also the private key in PEM format).\n" +
		                                "\tThis should include the whole certificate chain.\n" +
		                                "\tIf provided: The web socket will be a HTTPS socket.\n" +
		                                "\tIf not provided: Only HTTP.")
		tlsPrivateKey     = flag.String("web.tls-private-key", "", "Path to PEM file that conains the private key (if not contained in web.tls-cert file).")
		tlsClientCa       = flag.String("web.tls-client-ca", "", "Path to PEM file that conains the CAs that are trused for client connections.\n" +
		                                "\tIf provided: Connecting clients should present a certificate signed by one of this CAs.\n" +
		                                "\tIf not provided: Every client will be accepted.")
		enabledCollectors = flag.String("collectors.enabled", filterAvailableCollectors(defaultCollectors), "Comma-separated list of collectors to use.")
		printCollectors   = flag.Bool("collectors.print", false, "If true, print available collectors and exit.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("node_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting node_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	if *printCollectors {
		collectorNames := make(sort.StringSlice, 0, len(collector.Factories))
		for n := range collector.Factories {
			collectorNames = append(collectorNames, n)
		}
		collectorNames.Sort()
		fmt.Printf("Available collectors:\n")
		for _, n := range collectorNames {
			fmt.Printf(" - %s\n", n)
		}
		return
	}
	collectors, err := loadCollectors(*enabledCollectors)
	if err != nil {
		log.Fatalf("Couldn't load collectors: %s", err)
	}

	log.Infof("Enabled collectors:")
	for n := range collectors {
		log.Infof(" - %s", n)
	}

	nodeCollector := NodeCollector{collectors: collectors}
	prometheus.MustRegister(nodeCollector)

	server := &http.Server{
		Addr: *listenAddress,
		ErrorLog: createHttpServerLogWrapper(),
	}

	handler := prometheus.Handler()

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

	if (len(*tlsCert) > 0) {
		clientValidation := "no"
		if (len(*tlsClientCa) > 0 && len(*tlsCert) > 0) {
			certificates, err := loadCertificatesFrom(*tlsClientCa);
			if err != nil {
				log.Fatalf("Couldn't load client CAs from %s. Got: %s", *tlsClientCa, err)
			}
			server.TLSConfig = &tls.Config{
				ClientCAs: certificates,
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
			clientValidation = "yes"
		}
		targetTlsPrivateKey := *tlsPrivateKey
		if (len(targetTlsPrivateKey) <= 0) {
			targetTlsPrivateKey = *tlsCert;
		}
		log.Infof("Listening on %s (scheme=HTTPS, secured=TLS, clientValidation=%s)", server.Addr, clientValidation)
		err = server.ListenAndServeTLS(*tlsCert, targetTlsPrivateKey);
	} else {
		log.Infof("Listening on %s (scheme=HTTP, secured=no, clientValidation=no)", server.Addr)
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func loadCertificatesFrom(pemFile string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	certificates := x509.NewCertPool()
	certificates.AppendCertsFromPEM(caCert)
	return certificates, nil
}

type bufferedLogWriter struct {
	buf []byte
}

func (w *bufferedLogWriter) Write(p []byte) (n int, err error) {
	log.Debug(strings.TrimSpace(strings.Replace(string(p), "\n", " ", -1)))
	return len(p), nil
}

func createHttpServerLogWrapper() *slog.Logger {
	return slog.New(&bufferedLogWriter{}, "", 0)
}
