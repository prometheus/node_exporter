# Grafana integrations common lib

This common library can be used to quickly create dashboards' `panels` and `annotations`.

By using this common library we can 'enforce' common style choices across multiple dashboards and mixins.

## Import

```sh
jb init
jb install https://github.com/grafana/jsonnet-libs/common-lib
```

## Use

### Create new panel

```jsonnet

local commonlib = import 'github.com/grafana/jsonnet-libs/common-lib/common/main.libsonnet';
local cpuUsage = commonlib.panels.cpu.timeSeries.utilization.new(targets=[targets.cpuUsage]);

```

### Mutate exisiting panel with style options

```jsonnet

local commonlib = import 'github.com/grafana/jsonnet-libs/common-lib/common/main.libsonnet';
local cpuPanel = oldPanel + commonlib.panels.cpu.timeSeries.utilization.stylize();
```

See [windows-observ-lib](./windows-observ-lib) for full example.
