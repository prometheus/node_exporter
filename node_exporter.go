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
func addFlags(app *kingpin.Application) func(*slog.Logger) (config.Config, error) {
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
	arpDeviceInclude := app.Flag("collector.arp.device-include", "Regexp of arp devices to include (mutually exclusive to device-exclude).").String()
	arpDeviceExclude := app.Flag("collector.arp.device-exclude", "Regexp of arp devices to exclude (mutually exclusive to device-include).").String()
	arpNetlink := app.Flag("collector.arp.netlink", "Use netlink to gather stats instead of /proc/net/arp.").Default(strconv.FormatBool(config.DefaultARPNetlink)).Bool()
	ethtoolDeviceInclude := app.Flag("collector.ethtool.device-include", "Regexp of ethtool devices to include (mutually exclusive to device-exclude).").String()
	ethtoolDeviceExclude := app.Flag("collector.ethtool.device-exclude", "Regexp of ethtool devices to exclude (mutually exclusive to device-include).").String()
	ethtoolMetricsInclude := app.Flag("collector.ethtool.metrics-include", "Regexp of ethtool stats to include.").Default(config.DefaultEthtoolMetricsInclude).String()
	infinibandDeviceInclude := app.Flag("collector.infiniband.device-include", "Regexp of infiniband devices to include (mutually exclusive to device-exclude).").String()
	infinibandDeviceExclude := app.Flag("collector.infiniband.device-exclude", "Regexp of infiniband devices to exclude (mutually exclusive to device-include).").String()
	netclassIgnoredDevices := app.Flag("collector.netclass.ignored-devices", "Regexp of net devices to ignore for netclass collector.").Default(config.DefaultNetClassIgnoredDevices).String()
	netclassInvalidSpeed := app.Flag("collector.netclass.ignore-invalid-speed", "Ignore devices where the speed is invalid. This will be the default behavior in 2.x.").Bool()
	netclassNetlink := app.Flag("collector.netclass.netlink", "Use netlink to gather stats instead of /proc/net/dev.").Default("false").Bool()
	netclassRTNLWithStats := app.Flag("collector.netclass_rtnl.with-stats", "Expose the statistics for each network device, replacing netdev collector.").Bool()
	netdevDeviceInclude := app.Flag("collector.netdev.device-include", "Regexp of net devices to include (mutually exclusive to device-exclude).").String()
	oldNetdevDeviceInclude := app.Flag("collector.netdev.device-whitelist", "DEPRECATED: Use collector.netdev.device-include").Hidden().String()
	netdevDeviceExclude := app.Flag("collector.netdev.device-exclude", "Regexp of net devices to exclude (mutually exclusive to device-include).").String()
	oldNetdevDeviceExclude := app.Flag("collector.netdev.device-blacklist", "DEPRECATED: Use collector.netdev.device-exclude").Hidden().String()
	netdevAddressInfo := app.Flag("collector.netdev.address-info", "Collect address-info for every device").Bool()
	netdevDetailedMetrics := app.Flag("collector.netdev.enable-detailed-metrics", "Use (incompatible) metric names that provide more detailed stats on Linux").Bool()
	netdevNetlink := app.Flag("collector.netdev.netlink", "Use netlink to gather stats instead of /proc/net/dev.").Default(strconv.FormatBool(config.DefaultNetDevNetlink)).Bool()
	netdevLabelIfAlias := app.Flag("collector.netdev.label-ifalias", "Add ifAlias label").Default("false").Bool()
	qdiscFixtures := app.Flag("collector.qdisc.fixtures", "test fixtures to use for qdisc collector end-to-end testing").Default("").String()
	qdiscDeviceInclude := app.Flag("collector.qdisc.device-include", "Regexp of qdisc devices to include (mutually exclusive to device-exclude).").String()
	oldQdiscDeviceInclude := app.Flag("collector.qdisk.device-include", "DEPRECATED: Use collector.qdisc.device-include").Hidden().String()
	qdiscDeviceExclude := app.Flag("collector.qdisc.device-exclude", "Regexp of qdisc devices to exclude (mutually exclusive to device-include).").String()
	oldQdiscDeviceExclude := app.Flag("collector.qdisk.device-exclude", "DEPRECATED: Use collector.qdisc.device-exclude").Hidden().String()
	ipvsBackendLabels := app.Flag("collector.ipvs.backend-labels", "Comma separated list for IPVS backend stats labels.").Default(config.DefaultIPVSBackendLabels).String()
	var diskstatsExcludeSet, mountPointsExcludeSet, fsTypesExcludeSet bool
	diskstatsDeviceExclude := app.Flag("collector.diskstats.device-exclude", "Regexp of diskstats devices to exclude (mutually exclusive to device-include).").Default(config.DefaultDiskstatsDeviceExclude).IsSetByUser(&diskstatsExcludeSet).String()
	oldDiskstatsDeviceExclude := app.Flag("collector.diskstats.ignored-devices", "DEPRECATED: Use collector.diskstats.device-exclude").Hidden().String()
	diskstatsDeviceInclude := app.Flag("collector.diskstats.device-include", "Regexp of diskstats devices to include (mutually exclusive to device-exclude).").String()
	mountPointsExclude := app.Flag("collector.filesystem.mount-points-exclude", "Regexp of mount points to exclude for filesystem collector. (mutually exclusive to mount-points-include)").Default(config.DefaultFilesystemMountPointsExclude).IsSetByUser(&mountPointsExcludeSet).String()
	oldMountPointsExclude := app.Flag("collector.filesystem.ignored-mount-points", "Regexp of mount points to ignore for filesystem collector.").Hidden().String()
	mountPointsInclude := app.Flag("collector.filesystem.mount-points-include", "Regexp of mount points to include for filesystem collector. (mutually exclusive to mount-points-exclude)").String()
	fsTypesExclude := app.Flag("collector.filesystem.fs-types-exclude", "Regexp of filesystem types to exclude for filesystem collector. (mutually exclusive to fs-types-include)").Default(config.DefaultFilesystemFSTypesExclude).IsSetByUser(&fsTypesExcludeSet).String()
	oldFSTypesExclude := app.Flag("collector.filesystem.ignored-fs-types", "Regexp of filesystem types to ignore for filesystem collector.").Hidden().String()
	fsTypesInclude := app.Flag("collector.filesystem.fs-types-include", "Regexp of filesystem types to include for filesystem collector. (mutually exclusive to fs-types-exclude)").String()
	mountTimeout := app.Flag("collector.filesystem.mount-timeout", "how long to wait for a mount to respond before marking it as stale").Hidden().Default(config.DefaultFilesystemMountTimeout.String()).Duration()
	statWorkerCount := app.Flag("collector.filesystem.stat-workers", "how many stat calls to process simultaneously").Hidden().Default(strconv.Itoa(config.DefaultFilesystemStatWorkerCount)).Int()
	var systemdIncludeSet, systemdExcludeSet bool
	systemdUnitInclude := app.Flag("collector.systemd.unit-include", "Regexp of systemd units to include. Units must both match include and not match exclude to be included.").Default(config.DefaultSystemdUnitInclude).IsSetByUser(&systemdIncludeSet).String()
	oldSystemdUnitInclude := app.Flag("collector.systemd.unit-whitelist", "DEPRECATED: Use --collector.systemd.unit-include").Hidden().String()
	systemdUnitExclude := app.Flag("collector.systemd.unit-exclude", "Regexp of systemd units to exclude. Units must both match include and not match exclude to be included.").Default(config.DefaultSystemdUnitExclude).IsSetByUser(&systemdExcludeSet).String()
	oldSystemdUnitExclude := app.Flag("collector.systemd.unit-blacklist", "DEPRECATED: Use collector.systemd.unit-exclude").Hidden().String()
	systemdPrivate := app.Flag("collector.systemd.private", "Establish a private, direct connection to systemd without dbus (Strongly discouraged since it requires root. For testing purposes only).").Hidden().Bool()
	systemdTaskMetrics := app.Flag("collector.systemd.enable-task-metrics", "Enables service unit tasks metrics unit_tasks_current and unit_tasks_max").Bool()
	systemdRestartsMetrics := app.Flag("collector.systemd.enable-restarts-metrics", "Enables service unit metric service_restart_total").Bool()
	systemdStartTimeMetrics := app.Flag("collector.systemd.enable-start-time-metrics", "Enables service unit metric unit_start_time_seconds").Bool()
	perfCPUs := app.Flag("collector.perf.cpus", "List of CPUs from which perf metrics should be collected").Default("").String()
	perfTracepoints := app.Flag("collector.perf.tracepoint", "perf tracepoint that should be collected").Strings()
	perfNoHwProfiler := app.Flag("collector.perf.disable-hardware-profilers", "disable perf hardware profilers").Default("false").Bool()
	perfHwProfilers := app.Flag("collector.perf.hardware-profilers", "perf hardware profilers that should be collected").Strings()
	perfNoSwProfiler := app.Flag("collector.perf.disable-software-profilers", "disable perf software profilers").Default("false").Bool()
	perfSwProfilers := app.Flag("collector.perf.software-profilers", "perf software profilers that should be collected").Strings()
	perfNoCaProfiler := app.Flag("collector.perf.disable-cache-profilers", "disable perf cache profilers").Default("false").Bool()
	perfCaProfilers := app.Flag("collector.perf.cache-profilers", "perf cache profilers that should be collected").Strings()
	sysctlInclude := app.Flag("collector.sysctl.include", "Select sysctl metrics to include").Strings()
	sysctlIncludeInfo := app.Flag("collector.sysctl.include-info", "Select sysctl metrics to include as info metrics").Strings()
	hwmonChipInclude := app.Flag("collector.hwmon.chip-include", "Regexp of hwmon chip to include (mutually exclusive to device-exclude).").String()
	hwmonChipExclude := app.Flag("collector.hwmon.chip-exclude", "Regexp of hwmon chip to exclude (mutually exclusive to device-include).").String()
	hwmonSensorInclude := app.Flag("collector.hwmon.sensor-include", "Regexp of hwmon sensor to include (mutually exclusive to sensor-exclude).").String()
	hwmonSensorExclude := app.Flag("collector.hwmon.sensor-exclude", "Regexp of hwmon sensor to exclude (mutually exclusive to sensor-include).").String()
	wifiFixtures := app.Flag("collector.wifi.fixtures", "test fixtures to use for wifi collector metrics").Default("").String()

	return func(logger *slog.Logger) (config.Config, error) {
		netdevInclude, err := deprecatedFlag(logger, "collector.netdev.device-whitelist", *oldNetdevDeviceInclude, "collector.netdev.device-include", *netdevDeviceInclude, *netdevDeviceInclude != "")
		if err != nil {
			return config.Config{}, err
		}
		netdevExclude, err := deprecatedFlag(logger, "collector.netdev.device-blacklist", *oldNetdevDeviceExclude, "collector.netdev.device-exclude", *netdevDeviceExclude, *netdevDeviceExclude != "")
		if err != nil {
			return config.Config{}, err
		}
		qdiscInclude, err := deprecatedFlag(logger, "collector.qdisk.device-include", *oldQdiscDeviceInclude, "collector.qdisc.device-include", *qdiscDeviceInclude, *qdiscDeviceInclude != "")
		if err != nil {
			return config.Config{}, err
		}
		qdiscExclude, err := deprecatedFlag(logger, "collector.qdisk.device-exclude", *oldQdiscDeviceExclude, "collector.qdisc.device-exclude", *qdiscDeviceExclude, *qdiscDeviceExclude != "")
		if err != nil {
			return config.Config{}, err
		}

		diskstatsExclude, err := deprecatedFlag(logger, "collector.diskstats.ignored-devices", *oldDiskstatsDeviceExclude, "collector.diskstats.device-exclude", *diskstatsDeviceExclude, diskstatsExcludeSet)
		if err != nil {
			return config.Config{}, err
		}
		mountPointsExcl, err := deprecatedFlag(logger, "collector.filesystem.ignored-mount-points", *oldMountPointsExclude, "collector.filesystem.mount-points-exclude", *mountPointsExclude, mountPointsExcludeSet)
		if err != nil {
			return config.Config{}, err
		}
		fsTypesExcl, err := deprecatedFlag(logger, "collector.filesystem.ignored-fs-types", *oldFSTypesExclude, "collector.filesystem.fs-types-exclude", *fsTypesExclude, fsTypesExcludeSet)
		if err != nil {
			return config.Config{}, err
		}

		systemdExclude, err := deprecatedFlag(logger, "collector.systemd.unit-blacklist", *oldSystemdUnitExclude, "collector.systemd.unit-exclude", *systemdUnitExclude, systemdExcludeSet)
		if err != nil {
			return config.Config{}, err
		}
		systemdInclude, err := deprecatedFlag(logger, "collector.systemd.unit-whitelist", *oldSystemdUnitInclude, "collector.systemd.unit-include", *systemdUnitInclude, systemdIncludeSet)
		if err != nil {
			return config.Config{}, err
		}

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
		cfg.ARP = config.ARPConfig{
			DeviceInclude: *arpDeviceInclude,
			DeviceExclude: *arpDeviceExclude,
			Netlink:       *arpNetlink,
		}
		cfg.Ethtool = config.EthtoolConfig{
			DeviceInclude:  *ethtoolDeviceInclude,
			DeviceExclude:  *ethtoolDeviceExclude,
			MetricsInclude: *ethtoolMetricsInclude,
		}
		cfg.Infiniband = config.InfinibandConfig{
			DeviceInclude: *infinibandDeviceInclude,
			DeviceExclude: *infinibandDeviceExclude,
		}
		cfg.NetClass = config.NetClassConfig{
			IgnoredDevices:     *netclassIgnoredDevices,
			IgnoreInvalidSpeed: *netclassInvalidSpeed,
			Netlink:            *netclassNetlink,
			RTNLWithStats:      *netclassRTNLWithStats,
		}
		cfg.NetDev = config.NetDevConfig{
			DeviceInclude:   netdevInclude,
			DeviceExclude:   netdevExclude,
			AddressInfo:     *netdevAddressInfo,
			DetailedMetrics: *netdevDetailedMetrics,
			Netlink:         *netdevNetlink,
			LabelIfAlias:    *netdevLabelIfAlias,
		}
		cfg.Qdisc = config.QdiscConfig{
			Fixtures:      *qdiscFixtures,
			DeviceInclude: qdiscInclude,
			DeviceExclude: qdiscExclude,
		}
		cfg.IPVS = config.IPVSConfig{BackendLabels: *ipvsBackendLabels}
		cfg.Wifi = config.WifiConfig{Fixtures: *wifiFixtures}
		cfg.Systemd = config.SystemdConfig{
			UnitInclude:            systemdInclude,
			UnitExclude:            systemdExclude,
			Private:                *systemdPrivate,
			EnableTaskMetrics:      *systemdTaskMetrics,
			EnableRestartsMetrics:  *systemdRestartsMetrics,
			EnableStartTimeMetrics: *systemdStartTimeMetrics,
		}
		cfg.Perf = config.PerfConfig{
			CPUs:                     *perfCPUs,
			Tracepoints:              *perfTracepoints,
			DisableHardwareProfilers: *perfNoHwProfiler,
			HardwareProfilers:        *perfHwProfilers,
			DisableSoftwareProfilers: *perfNoSwProfiler,
			SoftwareProfilers:        *perfSwProfilers,
			DisableCacheProfilers:    *perfNoCaProfiler,
			CacheProfilers:           *perfCaProfilers,
		}
		cfg.Sysctl = config.SysctlConfig{Include: *sysctlInclude, IncludeInfo: *sysctlIncludeInfo}
		cfg.Hwmon = config.HwmonConfig{
			ChipInclude:   *hwmonChipInclude,
			ChipExclude:   *hwmonChipExclude,
			SensorInclude: *hwmonSensorInclude,
			SensorExclude: *hwmonSensorExclude,
		}
		cfg.Diskstats = config.DiskstatsConfig{
			DeviceExclude: excludeFlag(diskstatsExclude, diskstatsExcludeSet, *diskstatsDeviceInclude),
			DeviceInclude: *diskstatsDeviceInclude,
		}
		cfg.Filesystem = config.FilesystemConfig{
			MountPointsExclude: excludeFlag(mountPointsExcl, mountPointsExcludeSet, *mountPointsInclude),
			MountPointsInclude: *mountPointsInclude,
			FSTypesExclude:     excludeFlag(fsTypesExcl, fsTypesExcludeSet, *fsTypesInclude),
			FSTypesInclude:     *fsTypesInclude,
			MountTimeout:       *mountTimeout,
			StatWorkerCount:    *statWorkerCount,
		}
		return cfg, nil
	}
}

// excludeFlag returns the value of an exclude flag that has a default. An
// include flag replaces that default, unless the command line also sets the
// exclude flag.
func excludeFlag(exclude string, set bool, include string) string {
	if include != "" && !set {
		return ""
	}
	return exclude
}

// deprecatedFlag returns the value of the flag newName, or the value of the
// deprecated flag oldName when the command line sets only that one.
func deprecatedFlag(logger *slog.Logger, oldName, oldValue, newName, value string, set bool) (string, error) {
	if oldValue == "" {
		return value, nil
	}
	if set {
		return "", fmt.Errorf("--%s and --%s are mutually exclusive", oldName, newName)
	}
	logger.Warn(fmt.Sprintf("--%s is DEPRECATED and will be removed in 2.0.0, use --%s", oldName, newName))
	return oldValue, nil
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
			cfg, err := configFromFlags(b.Logger)
			if err != nil {
				return nil, err
			}

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
