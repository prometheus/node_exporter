# node_exporter

Prometheus exporter with pluggable metric collectors.



## Available collectors

By default the build will only include the native collectors
that expose information from /proc.

To include other collectors, specify the build tags lile this:

    go build -tags 'ganglia runit' node_exporter.go


Which collectors are used is controlled by the --enabledCollectors flag.

### NativeCollector

Provides metrics for load, seconds since last login and a list of tags
read from `node_exporter.conf`.


### GmondCollector (tag: ganglia)

Talks to a local gmond and provide it's metrics.


### RunitCollector (tag: runit)

Provides metrics for each runit services like state and how long it
has been in that state.

