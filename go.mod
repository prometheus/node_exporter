module github.com/prometheus/node_exporter

go 1.23.0

require (
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/beevik/ntp v1.4.3
	github.com/coreos/go-systemd/v22 v22.5.0
	github.com/dennwc/btrfs v0.0.0-20240418142341-0167142bde7a
	github.com/ema/qdisc v1.0.0
	github.com/godbus/dbus/v5 v5.1.0
	github.com/hashicorp/go-envparse v0.1.0
	github.com/hodgesds/perf-utils v0.7.0
	github.com/illumos/go-kstat v0.0.0-20210513183136-173c9b0a9973
	github.com/josharian/native v1.1.0
	github.com/jsimonetti/rtnetlink/v2 v2.0.2
	github.com/lufia/iostat v1.2.1
	github.com/mattn/go-xmlrpc v0.0.3
	github.com/mdlayher/ethtool v0.2.0
	github.com/mdlayher/netlink v1.7.2
	github.com/mdlayher/wifi v0.3.1
	github.com/opencontainers/selinux v1.11.1
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55
	github.com/prometheus-community/go-runit v0.1.0
	github.com/prometheus/client_golang v1.20.5
	github.com/prometheus/client_model v0.6.1
	github.com/prometheus/common v0.62.0
	github.com/prometheus/exporter-toolkit v0.14.0
	github.com/prometheus/procfs v0.15.2-0.20240603130017-1754b780536b // == v0.15.1 + https://github.com/prometheus/procfs/commit/1754b780536bb81082baa913e04cc4fff4d2baea
	github.com/safchain/ethtool v0.5.10
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0
	golang.org/x/sys v0.31.0
	howett.net/plist v1.0.1
)

require (
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dennwc/ioctl v1.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mdlayher/genetlink v1.3.2 // indirect
	github.com/mdlayher/socket v0.4.1 // indirect
	github.com/mdlayher/vsock v1.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/siebenmann/go-kstat v0.0.0-20210513183136-173c9b0a9973 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/oauth2 v0.28.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

//replace github.com/prometheus/procfs => github.com/rexagod/procfs v0.0.0-20241124020414-857c5b813f1b
