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
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/user"
	"runtime"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"

	"github.com/prometheus/node_exporter/collector"
)

const (
	defaultIPHeaderXForwardedFor = "X-Forwarded-For"
	defaultIPHeaderXRealIP       = "X-Real-IP"
	defaultIPHeaderXForwarded    = "X-Forwarded"
)

// handler wraps an unfiltered http.Handler but uses a filtered handler,
// created on the fly, if filtering is requested. Create instances with
// newHandler.
type handler struct {
	unfilteredHandler http.Handler
	// enabledCollectors list is used for logging and filtering
	enabledCollectors []string
	// exporterMetricsRegistry is a separate registry for the metrics about
	// the exporter itself.
	exporterMetricsRegistry *prometheus.Registry
	includeExporterMetrics  bool
	maxRequests             int
	logger                  *slog.Logger
	allowedNetworks         []*net.IPNet
	ipHeaders               []string
}

func newHandler(includeExporterMetrics bool, maxRequests int, logger *slog.Logger, allowedNetworks []*net.IPNet, ipHeaders []string) *handler {
	h := &handler{
		exporterMetricsRegistry: prometheus.NewRegistry(),
		includeExporterMetrics:  includeExporterMetrics,
		maxRequests:             maxRequests,
		logger:                  logger,
		allowedNetworks:         allowedNetworks,
		ipHeaders:               ipHeaders,
	}
	if h.includeExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
			promcollectors.NewGoCollector(),
		)
	}
	if innerHandler, err := h.innerHandler(); err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	} else {
		h.unfilteredHandler = innerHandler
	}
	return h
}

func (h *handler) getClientIP(r *http.Request) net.IP {
	headers := h.ipHeaders
	if len(headers) == 0 {
		headers = []string{
			defaultIPHeaderXForwardedFor,
			defaultIPHeaderXRealIP,
			defaultIPHeaderXForwarded,
		}
	}

	for _, header := range headers {
		ipStr := r.Header.Get(header)
		if ipStr == "" {
			continue
		}

		if header == defaultIPHeaderXForwardedFor {
			ips := strings.Split(ipStr, ",")
			if len(ips) > 0 {
				ipStr = strings.TrimSpace(ips[0])
			}
		}

		if ip := net.ParseIP(ipStr); ip != nil {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return net.ParseIP(r.RemoteAddr)
	}
	return net.ParseIP(host)
}

func (h *handler) isIPAllowed(ip net.IP) bool {
	if len(h.allowedNetworks) == 0 {
		return true
	}
	for _, network := range h.allowedNetworks {
		if network.Contains(ip) {
			h.logger.Debug("IP allowed by network", "ip", ip.String(), "network", network.String())
			return true
		}
	}
	h.logger.Debug("IP not in any allowed network", "ip", ip.String())
	return false
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(h.allowedNetworks) > 0 {
		clientIP := h.getClientIP(r)
		if clientIP == nil {
			h.logger.Debug("Could not parse client IP address", "remote_addr", r.RemoteAddr)
			http.Error(w, "Access denied: could not parse client IP address", http.StatusForbidden)
			return
		}
		if !h.isIPAllowed(clientIP) {
			h.logger.Debug("Access denied for IP", "ip", clientIP.String(), "remote_addr", r.RemoteAddr)
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		h.logger.Debug("Access allowed for IP", "ip", clientIP.String(), "remote_addr", r.RemoteAddr)
	}

	collects := r.URL.Query()["collect[]"]
	h.logger.Debug("collect query:", "collects", collects)

	excludes := r.URL.Query()["exclude[]"]
	h.logger.Debug("exclude query:", "excludes", excludes)

	if len(collects) == 0 && len(excludes) == 0 {
		// No filters, use the prepared unfiltered handler.
		h.unfilteredHandler.ServeHTTP(w, r)
		return
	}

	if len(collects) > 0 && len(excludes) > 0 {
		h.logger.Debug("rejecting combined collect and exclude queries")
		http.Error(w, "Combined collect and exclude queries are not allowed.", http.StatusBadRequest)
		return
	}

	filters := &collects
	if len(excludes) > 0 {
		// In exclude mode, filtered collectors = enabled - excludeed.
		f := []string{}
		for _, c := range h.enabledCollectors {
			if (slices.Index(excludes, c)) == -1 {
				f = append(f, c)
			}
		}
		filters = &f
	}

	filteredHandler, err := h.innerHandler(*filters...)
	if err != nil {
		h.logger.Warn("Couldn't create filtered metrics handler", "err", err)
		http.Error(w, fmt.Sprintf("Couldn't create filtered metrics handler: %s", err), http.StatusBadRequest)
		return
	}
	filteredHandler.ServeHTTP(w, r)
}

// innerHandler is used to create both the one unfiltered http.Handler to be
// wrapped by the outer handler and also the filtered handlers created on the
// fly. The former is accomplished by calling innerHandler without any arguments
// (in which case it will log all the collectors enabled via command-line
// flags).
func (h *handler) innerHandler(filters ...string) (http.Handler, error) {
	nc, err := collector.NewNodeCollector(h.logger, filters...)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	// Only log the creation of an unfiltered handler, which should happen
	// only once upon startup.
	if len(filters) == 0 {
		h.logger.Info("Enabled collectors")
		for n := range nc.Collectors {
			h.enabledCollectors = append(h.enabledCollectors, n)
		}
		sort.Strings(h.enabledCollectors)
		for _, c := range h.enabledCollectors {
			h.logger.Info(c)
		}
	}

	r := prometheus.NewRegistry()
	r.MustRegister(versioncollector.NewCollector("node_exporter"))
	if err := r.Register(nc); err != nil {
		return nil, fmt.Errorf("couldn't register node collector: %s", err)
	}

	var handler http.Handler
	if h.includeExporterMetrics {
		handler = promhttp.HandlerFor(
			prometheus.Gatherers{h.exporterMetricsRegistry, r},
			promhttp.HandlerOpts{
				ErrorLog:            slog.NewLogLogger(h.logger.Handler(), slog.LevelError),
				ErrorHandling:       promhttp.ContinueOnError,
				MaxRequestsInFlight: h.maxRequests,
				Registry:            h.exporterMetricsRegistry,
			},
		)
		// Note that we have to use h.exporterMetricsRegistry here to
		// use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			h.exporterMetricsRegistry, handler,
		)
	} else {
		handler = promhttp.HandlerFor(
			r,
			promhttp.HandlerOpts{
				ErrorLog:            slog.NewLogLogger(h.logger.Handler(), slog.LevelError),
				ErrorHandling:       promhttp.ContinueOnError,
				MaxRequestsInFlight: h.maxRequests,
			},
		)
	}

	return handler, nil
}

func parseAllowedNetworks(networkStrings []string) ([]*net.IPNet, error) {
	if len(networkStrings) == 0 {
		return nil, nil
	}

	networks := make([]*net.IPNet, 0, len(networkStrings))
	for _, networkStr := range networkStrings {
		networkStr = strings.TrimSpace(networkStr)
		if networkStr == "" {
			continue
		}

		if !strings.Contains(networkStr, "/") {
			ip := net.ParseIP(networkStr)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP address: %s", networkStr)
			}
			if ip.To4() != nil {
				networkStr = networkStr + "/32"
			} else {
				networkStr = networkStr + "/128"
			}
		}

		_, network, err := net.ParseCIDR(networkStr)
		if err != nil {
			return nil, fmt.Errorf("invalid network CIDR %s: %w", networkStr, err)
		}
		networks = append(networks, network)
	}

	return networks, nil
}

type whitelistConfig struct {
	AllowedNetworks []string `yaml:"allowed_networks"`
	IPHeaders       []string `yaml:"ip_headers"`
}

func loadWhitelistConfig(configPath string) (*whitelistConfig, error) {
	if configPath == "" {
		return nil, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config struct {
		Whitelist whitelistConfig `yaml:"whitelist"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if len(config.Whitelist.AllowedNetworks) == 0 && len(config.Whitelist.IPHeaders) == 0 {
		return nil, nil
	}

	return &config.Whitelist, nil
}

func loadWhitelistSettings(configPath string, networksFlag string, logger *slog.Logger) ([]*net.IPNet, []string, error) {
	var networks []*net.IPNet
	var ipHeaders []string

	config, err := loadWhitelistConfig(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load whitelist config: %w", err)
	}

	if config != nil {
		if len(config.AllowedNetworks) > 0 {
			networks, err = parseAllowedNetworks(config.AllowedNetworks)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse allowed networks from config: %w", err)
			}
		}
		if len(config.IPHeaders) > 0 {
			ipHeaders = config.IPHeaders
		}
	}

	if networksFlag != "" {
		networkStrings := strings.Split(networksFlag, ",")
		networks, err = parseAllowedNetworks(networkStrings)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse allowed networks from flag: %w", err)
		}
	}

	return networks, ipHeaders, nil
}

func main() {
	var (
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		disableExporterMetrics = kingpin.Flag(
			"web.disable-exporter-metrics",
			"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
		).Bool()
		maxRequests = kingpin.Flag(
			"web.max-requests",
			"Maximum number of parallel scrape requests. Use 0 to disable.",
		).Default("40").Int()
		allowedNetworks = kingpin.Flag(
			"web.allowed-networks",
			"Comma-separated list of allowed IP networks in CIDR notation (e.g., 192.168.1.0/24,10.0.0.0/8). Single IPs are also accepted and will be treated as /32 for IPv4 or /128 for IPv6.",
		).String()
		whitelistConfigPath = kingpin.Flag(
			"web.whitelist-config",
			"Path to YAML configuration file for IP whitelist settings.",
		).String()
		disableDefaultCollectors = kingpin.Flag(
			"collector.disable-defaults",
			"Set all collectors to disabled by default.",
		).Default("false").Bool()
		maxProcs = kingpin.Flag(
			"runtime.gomaxprocs", "The target number of CPUs Go will run on (GOMAXPROCS)",
		).Envar("GOMAXPROCS").Default("1").Int()
		toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":9100")
	)

	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.Version(version.Print("node_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promslog.New(promslogConfig)

	networks, ipHeaders, err := loadWhitelistSettings(*whitelistConfigPath, *allowedNetworks, logger)
	if err != nil {
		logger.Error("Failed to load whitelist settings", "error", err)
		os.Exit(1)
	}

	if len(networks) > 0 {
		logger.Info("IP whitelist enabled", "networks", len(networks))
		for _, network := range networks {
			logger.Info("Allowed network", "network", network.String())
		}
		if len(ipHeaders) > 0 {
			logger.Info("IP headers configured", "headers", strings.Join(ipHeaders, ", "))
		}
	}

	if *disableDefaultCollectors {
		collector.DisableDefaultCollectors()
	}
	logger.Info("Starting node_exporter", "version", version.Info())
	logger.Info("Build context", "build_context", version.BuildContext())
	if user, err := user.Current(); err == nil && user.Uid == "0" {
		logger.Warn("Node Exporter is running as root user. This exporter is designed to run as unprivileged user, root is not required.")
	}
	runtime.GOMAXPROCS(*maxProcs)
	logger.Debug("Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))

	metricsHandler := newHandler(!*disableExporterMetrics, *maxRequests, logger, networks, ipHeaders)
	http.Handle(*metricsPath, metricsHandler)

	if *metricsPath != "/" {
		landingConfig := web.LandingConfig{
			Name:        "Node Exporter",
			Description: "Prometheus Node Exporter",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		landingHandler := newHandler(false, 0, logger, networks, ipHeaders)
		landingHandler.unfilteredHandler = landingPage
		http.Handle("/", landingHandler)
	}

	server := &http.Server{}
	if err := web.ListenAndServe(server, toolkitFlags, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
