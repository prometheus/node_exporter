module github.com/prometheus/node_exporter/otel_collector

go 1.25.0

require (
	github.com/prometheus/client_golang v1.23.2
	github.com/prometheus/exporter-toolkit v0.14.1
	github.com/prometheus/node_exporter v1.10.2
	go.opentelemetry.io/collector/component v1.44.0
	go.opentelemetry.io/collector/receiver v1.44.0
)

require (
	github.com/alecthomas/kingpin/v2 v2.4.0 // indirect
	github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b // indirect
	github.com/beevik/ntp v1.5.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.6.0 // indirect
	github.com/dennwc/btrfs v0.0.0-20241002142654-12ae127e0bf6 // indirect
	github.com/dennwc/ioctl v1.0.0 // indirect
	github.com/ema/qdisc v1.0.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hashicorp/go-envparse v0.1.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hodgesds/perf-utils v0.7.0 // indirect
	github.com/illumos/go-kstat v0.0.0-20210513183136-173c9b0a9973 // indirect
	github.com/jsimonetti/rtnetlink/v2 v2.0.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lufia/iostat v1.2.1 // indirect
	github.com/mattn/go-xmlrpc v0.0.3 // indirect
	github.com/mdlayher/ethtool v0.5.0 // indirect
	github.com/mdlayher/genetlink v1.3.2 // indirect
	github.com/mdlayher/netlink v1.8.0 // indirect
	github.com/mdlayher/socket v0.5.1 // indirect
	github.com/mdlayher/wifi v0.6.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/selinux v1.12.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus-community/go-runit v0.1.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.1 // indirect
	github.com/prometheus/procfs v0.19.0 // indirect
	github.com/safchain/ethtool v0.6.2 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector/consumer v1.44.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.44.0 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.138.0 // indirect
	go.opentelemetry.io/collector/pdata v1.44.0 // indirect
	go.opentelemetry.io/collector/pipeline v1.44.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.13.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/log v0.14.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/exp v0.0.0-20250911091902-df9299821621 // indirect
	golang.org/x/net v0.45.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250825161204-c5933d9347a5 // indirect
	google.golang.org/grpc v1.76.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	howett.net/plist v1.0.1 // indirect
)

replace github.com/prometheus/exporter-toolkit => ../../exporter-toolkit

replace github.com/prometheus/node_exporter => ../
