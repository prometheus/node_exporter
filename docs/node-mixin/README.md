# Node Mixin

_This is work in progress. We aim for it to become a good role model for alerts
and dashboards eventually, but it is not quite there yet._

The Node Mixin is a set of configurable, reusable, and extensible alerts and
dashboards based on the metrics exported by the Node Exporter. The mixin create
recording and alerting rules for Prometheus and suitable dashboard descriptions
for Grafana.

To use them, you need to have `jsonnet` (v0.13+) and `jb` installed. If you
have a working Go development environment, it's easiest to run the following:
```bash
$ go get github.com/google/go-jsonnet/cmd/jsonnet
$ go get github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb
```

_Note: The make targets `lint` and `fmt` need the `jsonnetfmt` binary, which is
currently not included in the Go implementation of `jsonnet`. For the time
being, you have to install the [C++ version of
jsonnetfmt](https://github.com/google/jsonnet) if you want to use `make lint`
or `make fmt`._

Next, install the dependencies by running the following command in this
directory:
```bash
$ jb install
```

You can then build the Prometheus rules files `node_alerts.yaml` and
`node_rules.yaml`:
```bash
$ make node_alerts.yaml node_rules.yaml
```

You can also build a directory `dashboard_out` with the JSON dashboard files
for Grafana:
```bash
$ make dashboards_out
```

For more advanced uses of mixins, see
https://github.com/monitoring-mixins/docs.

