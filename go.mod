module github.com/prometheus/node_exporter

require (
	github.com/beevik/ntp v0.3.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/ema/qdisc v0.0.0-20200603082823-62d0308e3e00
	github.com/go-kit/log v0.1.0
	github.com/godbus/dbus v0.0.0-20190402143921-271e53dc4968
	github.com/hodgesds/perf-utils v0.2.5
	github.com/illumos/go-kstat v0.0.0-20210513183136-173c9b0a9973
	github.com/jsimonetti/rtnetlink v0.0.0-20210713125558-2bfdf1dbdbd6
	github.com/lufia/iostat v1.1.0
	github.com/mattn/go-xmlrpc v0.0.3
	github.com/mdlayher/wifi v0.0.0-20200527114002-84f0b9457fdd
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.29.0
	github.com/prometheus/exporter-toolkit v0.6.0
	github.com/prometheus/procfs v0.7.2
	github.com/safchain/ethtool v0.0.0-20201023143004-874930cb3ce0
	github.com/siebenmann/go-kstat v0.0.0-20210513183136-173c9b0a9973 // indirect
	github.com/soundcloud/go-runit v0.0.0-20150630195641-06ad41a06c4a
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

go 1.14
