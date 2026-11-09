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
	"log/slog"
	"maps"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/user"
	"runtime"
	"slices"
	"strconv"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/bootstrap"

	"github.com/prometheus/node_exporter/collector"
	"github.com/prometheus/node_exporter/config"
)

// handler wraps an unfiltered http.Handler but uses a filtered handler,
// created on the fly, if filtering is requested. Create instances with
// newHandler.
type handler struct {
	unfilteredHandler http.Handler
	runtime           *collector.Runtime
	// exporterMetricsRegistry is a separate registry for the metrics about
	// the exporter itself.
	exporterMetricsRegistry *prometheus.Registry
	includeExporterMetrics  bool
	maxRequests             int
	logger                  *slog.Logger
}

func newHandler(cfg config.Config, includeExporterMetrics bool, maxRequests int, logger *slog.Logger) (*handler, error) {
	runtime, err := collector.NewRuntime(cfg, logger)
	if err != nil {
		return nil, err
	}

	h := &handler{
		runtime:                 runtime,
		exporterMetricsRegistry: prometheus.NewRegistry(),
		includeExporterMetrics:  includeExporterMetrics,
		maxRequests:             maxRequests,
		logger:                  logger,
	}
	if h.includeExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
			promcollectors.NewGoCollector(),
		)
	}
	innerHandler, err := h.innerHandler()
	if err != nil {
		return nil, fmt.Errorf("couldn't create metrics handler: %w", err)
	}
	h.unfilteredHandler = innerHandler
	return h, nil
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Combined collect and exclude queries are not allowed."))
		return
	}

	filters := &collects
	if len(excludes) > 0 {
		// In exclude mode, filtered collectors = enabled - excludeed.
		f := []string{}
		for _, c := range h.runtime.EnabledCollectors() {
			if slices.Index(excludes, c) == -1 {
				f = append(f, c)
			}
		}
		filters = &f
	}

	// To serve filtered metrics, we create a filtering handler on the fly.
	filteredHandler, err := h.innerHandler(*filters...)
	if err != nil {
		h.logger.Warn("Couldn't create filtered metrics handler:", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Couldn't create filtered metrics handler: %s", err)
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
	runtime := h.runtime
	if len(filters) == 0 {
		// Only the unfiltered handler logs, and it is built once at startup.
		h.logger.Info("Enabled collectors")
		for _, c := range runtime.EnabledCollectors() {
			h.logger.Info(c)
		}
	} else {
		r, err := h.runtime.Filtered(filters...)
		if err != nil {
			return nil, fmt.Errorf("couldn't create collector runtime: %w", err)
		}
		runtime = r
	}

	r, err := runtime.Registry()
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector registry: %w", err)
	}
	r.MustRegister(versioncollector.NewCollector("node_exporter"))

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

// addFlags adds the collector flags to app. The returned function builds a
// Config from them after app parses the command line.
func addFlags(app *kingpin.Application) func() config.Config {
	disableDefaults := app.Flag("collector.disable-defaults", "Set all collectors to disabled by default.").Default("false").Bool()

	type collectorFlag struct {
		enabled, setByUser bool
	}
	collectors := collector.AvailableCollectors()
	collectorFlags := make(map[string]*collectorFlag, len(collectors))
	for _, name := range slices.Sorted(maps.Keys(collectors)) {
		state := "disabled"
		if collectors[name] {
			state = "enabled"
		}
		f := &collectorFlag{}
		app.Flag("collector."+name, fmt.Sprintf("Enable the %s collector (default: %s).", name, state)).
			Default(strconv.FormatBool(collectors[name])).IsSetByUser(&f.setByUser).BoolVar(&f.enabled)
		collectorFlags[name] = f
	}

	cpuGuest := app.Flag("collector.cpu.guest", "Enables metric node_cpu_guest_seconds_total").Default(strconv.FormatBool(config.DefaultCPUEnableGuest)).Bool()
	cpuInfo := app.Flag("collector.cpu.info", "Enables metric cpu_info").Bool()
	cpuFlagsInclude := app.Flag("collector.cpu.info.flags-include", "Filter the `flags` field in cpuInfo with a value that must be a regular expression").String()
	cpuBugsInclude := app.Flag("collector.cpu.info.bugs-include", "Filter the `bugs` field in cpuInfo with a value that must be a regular expression").String()
	textfileDirectories := app.Flag("collector.textfile.directory", "Directory to read text files with metrics from, supports glob matching. (repeatable)").Strings()
	ntpServer := app.Flag("collector.ntp.server", "NTP server to use for ntp collector").Default(config.DefaultNTPServer).String()
	ntpServerPort := app.Flag("collector.ntp.server-port", "UDP port number to connect to on NTP server").Default(strconv.Itoa(config.DefaultNTPServerPort)).Int()
	ntpProtocolVersion := app.Flag("collector.ntp.protocol-version", "NTP protocol version").Default(strconv.Itoa(config.DefaultNTPProtocolVersion)).Int()
	ntpServerIsLocal := app.Flag("collector.ntp.server-is-local", "Certify that collector.ntp.server address is not a public ntp server").Default("false").Bool()
	ntpIPTTL := app.Flag("collector.ntp.ip-ttl", "IP TTL to use while sending NTP query").Default(strconv.Itoa(config.DefaultNTPIPTTL)).Int()
	ntpMaxDistance := app.Flag("collector.ntp.max-distance", "Max accumulated distance to the root").Default(config.DefaultNTPMaxDistance.String()).Duration()
	ntpOffsetTolerance := app.Flag("collector.ntp.local-offset-tolerance", "Offset between local clock and local ntpd time to tolerate").Default(config.DefaultNTPLocalOffsetTolerance.String()).Duration()
	statSoftirq := app.Flag("collector.stat.softirq", "Export softirq calls per vector").Default("false").Bool()
	vmstatFields := app.Flag("collector.vmstat.fields", "Regexp of fields to return for vmstat collector.").Default(config.DefaultVMStatFields).String()
	netstatFields := app.Flag("collector.netstat.fields", "Regexp of fields to return for netstat collector.").Default(config.DefaultNetstatFields).String()
	bcachePriorityStats := app.Flag("collector.bcache.priorityStats", "Expose expensive priority stats.").Bool()
	nvmeNamespaceUsage := app.Flag("collector.nvme.namespace-usage", "Expose namespace used bytes (reads nuse, which makes the kernel query the controller).").Default(strconv.FormatBool(config.DefaultNVMeNamespaceUsage)).Bool()
	raplZoneLabel := app.Flag("collector.rapl.enable-zone-label", "Enables service unit metric unit_start_time_seconds").Bool()

	return func() config.Config {
		cfg := collector.ConfigFromFlags(config.NewConfigWithDefaults())
		cfg.CollectorDisableDefaults = *disableDefaults
		for name, f := range collectorFlags {
			if f.setByUser {
				if cfg.Collectors == nil {
					cfg.Collectors = map[string]bool{}
				}
				cfg.Collectors[name] = f.enabled
			}
		}
		cfg.CPU = config.CPUConfig{
			EnableGuest:  *cpuGuest,
			EnableInfo:   *cpuInfo,
			FlagsInclude: *cpuFlagsInclude,
			BugsInclude:  *cpuBugsInclude,
		}
		cfg.Textfile = config.TextfileConfig{Directories: *textfileDirectories}
		cfg.NTP = config.NTPConfig{
			Server:               *ntpServer,
			ServerPort:           *ntpServerPort,
			ProtocolVersion:      *ntpProtocolVersion,
			ServerIsLocal:        *ntpServerIsLocal,
			IPTTL:                *ntpIPTTL,
			MaxDistance:          *ntpMaxDistance,
			LocalOffsetTolerance: *ntpOffsetTolerance,
		}
		cfg.Stat = config.StatConfig{Softirq: *statSoftirq}
		cfg.VMStat = config.VMStatConfig{Fields: *vmstatFields}
		cfg.Netstat = config.NetstatConfig{Fields: *netstatFields}
		cfg.Bcache = config.BcacheConfig{PriorityStats: *bcachePriorityStats}
		cfg.NVMe = config.NVMeConfig{NamespaceUsage: *nvmeNamespaceUsage}
		cfg.RAPL = config.RAPLConfig{EnableZoneLabel: *raplZoneLabel}
		return cfg
	}
}

func main() {
	var (
		configFromFlags = addFlags(kingpin.CommandLine)
		maxProcs        = kingpin.Flag(
			"runtime.gomaxprocs", "The target number of CPUs Go will run on (GOMAXPROCS)",
		).Envar("GOMAXPROCS").Default("1").Int()
	)

	kingpin.CommandLine.UsageWriter(os.Stdout)
	runner := bootstrap.New(bootstrap.Config{
		App:            kingpin.CommandLine,
		Name:           "node_exporter",
		Description:    "Prometheus Node Exporter",
		DefaultAddress: ":9100",
		MetricsHandlerFactory: func(b *bootstrap.Bootstrap) (http.Handler, error) {
			cfg := configFromFlags()

			if user, err := user.Current(); err == nil && user.Uid == "0" {
				b.Logger.Warn("Node Exporter is running as root user. This exporter is designed to run as unprivileged user, root is not required.")
			}
			runtime.GOMAXPROCS(*maxProcs)
			b.Logger.Debug("Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))
			b.Handle("/debug/pprof/", http.DefaultServeMux)
			return newHandler(cfg, !b.DisableExporterMetrics, b.MaxRequests, b.Logger)
		},
	})

	if err := runner.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
