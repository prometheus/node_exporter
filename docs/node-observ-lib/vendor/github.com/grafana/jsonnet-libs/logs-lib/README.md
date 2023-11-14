# Logs lib

This logs lib can be used to generate logs dashboard using [grafonnet](https://github.com/grafana/grafonnet).

## Import

```sh
jb init
jb install https://github.com/grafana/jsonnet-libs/logs-lib
```

## Examples

### Generate kubernetes logs dashboard

```jsonnet
local logslib = import 'github.com/grafana/jsonnet-libs/logs-lib/logs/main.libsonnet';

//Additional selector to add to all variable queries and alerts(if any)
local kubeFilterSelector = 'namespace!=""';
// Array of labels to compose chained grafana variables (order matters)
local kubeLabels = ['cluster', 'namespace', 'app', 'pod', 'container'];

// pick one of Loki's parsers to use: i.e. logfmt, json.
// | __error__=`` is appended automatically
// https://grafana.com/docs/loki/latest/logql/log_queries/#parser-expression
// set null or do not provide at all if parsing is not required.
local formatParser = 'logfmt';

//group by 'app' label instead of 'level':
local logsVolumeGroupBy = 'app';

//extra filters to do advanced line_format:
local extraFilters = |||
  | label_format timestamp="{{__timestamp__}}"
  | line_format `{{ if eq "[[pod]]" ".*" }}{{.pod | trunc 20}}:{{else}}{{.container}}:{{end}} {{__line__}}`
|||;

(
  logsDashboard.new('Kubernetes apps logs',
                    datasourceRegex='',
                    filterSelector=kubeFilterSelector,
                    labels=kubeLabels,
                    formatParser=formatParser,
                    logsVolumeGroupBy=logsVolumeGroupBy,
                    extraFilters=extraFilters)
).dashboards.logs
```

![image](https://github.com/grafana/jsonnet-libs/assets/14870891/7b246cc9-5de1-42f5-b3cd-bb9f89302405)

### Generate systemd logs dashboard and modify panels and variables

This lib exposes `variables`, `targets`, `panels`, and `dashboards`.

Because of that, you can override options of those objects before exporting the dashboard.

Again, use [Grafonnet](https://grafana.github.io/grafonnet/API/panel/index.html) for this:

```jsonnet
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';
local logslib = import 'github.com/grafana/jsonnet-libs/logs-lib/logs/main.libsonnet';


local linuxFilterSelector = 'unit!=""';
local linuxLabels = ['job', 'instance', 'unit', 'level'];

// pick one of Loki's parsers to use: i.e. logfmt, json.
// | __error__=`` is appended automatically
// https://grafana.com/docs/loki/latest/logql/log_queries/#parser-expression
// set null or do not provide at all if parsing is not required.
local formatParser = 'unpack';

// 2. create and export systemd logs dashboard
local systemdLogs =
  logslib.new('Linux systemd logs',
                    datasourceRegex='',
                    filterSelector=linuxFilterSelector,
                    labels=linuxLabels,
                    formatParser=formatParser,
                    showLogsVolume=true)
  // override panels or variables using grafonnet
  {
    panels+:
      {
        logs+:
          g.panel.logs.options.withEnableLogDetails(false),
      },
    variables+:
      {
        regex_search+:
          g.dashboard.variable.textbox.new('regex_search', default='error'),
      },
  };
// export logs dashboard
systemdLogs.dashboards.logs

```

![image](https://github.com/grafana/jsonnet-libs/assets/14870891/5e6313fd-9135-446a-b7bf-cf124b436970)

### Generate docker logs dashboard

```jsonnet

local logslib = import 'github.com/grafana/jsonnet-libs/logs-lib/logs/main.libsonnet';

// Array of labels to compose chained grafana variables
local dockerFilterSelector = 'container_name!=""';
local dockerLabels = ['job', 'instance', 'container_name'];

// pick one of Loki's parsers to use: i.e. logfmt, json.
// | __error__=`` is appended automatically
// https://grafana.com/docs/loki/latest/logql/log_queries/#parser-expression
// set null or do not provide at all if parsing is not required.
local formatParser = 'logfmt';

(
  logslib.new('Docker logs',
                    datasourceRegex='',
                    filterSelector=dockerFilterSelector,
                    labels=dockerLabels,
                    formatParser=formatParser)
).dashboards.logs

```
