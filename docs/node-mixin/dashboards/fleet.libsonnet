local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local nodePanels = import '../lib/panels/panels.libsonnet';
local commonPanels = import '../lib/panels/common/panels.libsonnet';
local nodeTimeseries = nodePanels.timeseries;
local common = import '../lib/common.libsonnet';

{

  new(config=null, platform=null):: {
    local c = common.new(config=config, platform=platform),
    local commonPromTarget = c.commonPromTarget,
    local templates = [
      if template.label == 'Instance'
      then
        template
        {
          allValue: '.+',
          includeAll: true,
          multi: true,
        }
      else template
      for template in c.templates
    ],

    local q = c.queries,

    local fleetTable =
      nodePanels.table.new(
        title='Linux Nodes Overview'
      )
      //ABCDEFGHIJKLMNOPQRSTUVWXYZ
      .addTarget(commonPromTarget(expr=q.osInfo, format='table', instant=true) {refId: "INFO"})
      .addTarget(commonPromTarget(expr=q.nodeInfo, format='table', instant=true) {refId: "OS"})
      .addTarget(commonPromTarget(expr=q.uptime, format='table', instant=true) {refId: "UPTIME"})
      .addTarget(commonPromTarget(expr=q.systemLoad1, format='table', instant=true) {refId: "LOAD1"})
      .addTarget(commonPromTarget(expr=q.systemLoad5, format='table', instant=true) {refId: "LOAD5"})
      .addTarget(commonPromTarget(expr=q.systemLoad15, format='table', instant=true) {refId: "LOAD15"})
      .addTarget(commonPromTarget(
        expr=q.cpuCount,
        format='table',
        instant=true,
      ) {refId: "CPUCOUNT"})
      .addTarget(commonPromTarget(
        expr=q.cpuUsage, format='table', instant=true,
      ) {refId: "CPUUSAGE"})
      .addTarget(commonPromTarget(expr=q.memoryTotal, format='table', instant=true) {refId: "MEMTOTAL"})
      .addTarget(commonPromTarget(expr=q.memoryUsage, format='table', instant=true) {refId: "MEMUSAGE"})
      .addTarget(commonPromTarget(expr=q.fsSizeTotalRoot, format='table', instant=true) {refId: "FSTOTAL"})
      .addTarget(commonPromTarget(
        expr=
        |||
          100-(max by (instance) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance=~"$instance", fstype!="", mountpoint="/"})
          /
          max by (instance) (node_filesystem_size_bytes{%(nodeExporterSelector)s, instance=~"$instance", fstype!="", mountpoint="/"}) * 100)
        ||| % config,
        format='table',
        instant=true,
      ) {refId: "FSUSAGE"})
      .addTarget(commonPromTarget(
        expr='count by (instance) (max_over_time(ALERTS{%(nodeExporterSelector)s, instance=~"$instance", alertstate="firing", severity="critical"}[1m]))' %config,
        format='table',
        instant=true
        ) {refId: "CRITICAL"})
      .addTarget(commonPromTarget(
        expr='count by (instance) (max_over_time(ALERTS{%(nodeExporterSelector)s, instance=~"$instance", alertstate="firing", severity="warning"}[1m]))' %config,
        format='table',
        instant=true
        ) {refId: "WARNING"})
      .withTransform()
        .joinByField(field='instance')
        //disable kernel and os:
        //.filterFieldsByName('instance|pretty_name|nodename|release|Value.+')
        .filterFieldsByName('instance|nodename|Value.+')
        .organize(
          excludeByName={
            'Value #OS': true,
            'Value #INFO': true,
            'Value #LOAD5': true,
            'Value #LOAD15': true,
          },
          renameByName={
            'instance': 'Instance',
            'pretty_name': 'OS',
            'nodename': 'Hostname',
            'release': 'Kernel version',
            'Value #LOAD1': 'Load 1m',
            'Value #LOAD5': 'Load 5m',
            'Value #LOAD15': 'Load 15m',
            'Value #CPUCOUNT': 'Cores',
            'Value #CPUUSAGE': 'CPU usage',
            'Value #MEMTOTAL': 'Memory total',
            'Value #MEMUSAGE': 'Memory usage',
            'Value #FSTOTAL': 'Root disk size',
            'Value #FSUSAGE': 'Root disk usage',
            'Value #UPTIME': 'Uptime',
            'Value #CRITICAL': 'Crit Alerts',
            'Value #WARNING': 'Warnings',
          }
        )
      .withFooter(reducer=['mean'], fields=[
        'Value #LOAD1',
        "Value #MEMUSAGE",
        "Value #CPUUSAGE"
      ])
      //.addThresholdStep(color='light-green', value=null)
      .addThresholdStep(color='text', value=null)
      .addThresholdStep(color='light-yellow', value=80)
      .addThresholdStep(color='light-red', value=90)
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Instance',
        },
        properties=[
          {
            id: 'links',
            value: [
              {
                targetBlank: true,
                title: 'Drill down to instance ${__data.fields.instance}',
                url: 'd/nodes?var-instance=${__data.fields.instance}&${__url_time_range}',
              },
            ],
          },
          {
            "id": "custom.filterable",
            "value": true
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'OS|Kernel version|Hostname',
        },
        properties=[
          {
            "id": "custom.filterable",
            "value": true
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'Memory total|Root disk size',
        },
        properties=[
          {
            id: 'unit',
            value: 'bytes',
          },
          {
            "id": "decimals",
            "value": 0
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Cores'
        },
        properties=[
          {
            "id": "custom.width",
            "value": 60
          }
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'Load.+'
        },
        properties=[
          {
            "id": "custom.width",
            "value": 60
          }
        ]
      )
      .addOverride(
        matcher={
          id: 'byName',
          options: 'Uptime',
        },
        properties=[
          {
            id: 'unit',
            value: 'dtdurations',
          },
          {
            "id": "custom.displayMode",
            "value": "color-text"
          },
          {
            "id": "thresholds",
            "value": {
              "mode": "absolute",
              "steps": [
                {
                  "color": "light-orange",
                  "value": null
                },
                {
                  "color": "text",
                  "value": 300
                }
              ]
            }
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'CPU usage|Memory usage|Root disk usage',
        },
        properties=[
          {
            id: 'unit',
            value: 'percent',
          },
          // {
          //   id: 'custom.displayMode',
          //   value: 'gradient-gauge',
          // },
          {
            "id": "custom.displayMode",
            "value": "basic"
          },
          {
            id: 'max',
            value: 100,
          },
          {
            id: 'min',
            value: 0,
          },
        ]
      )
      .sortBy('Instance')
      ,
    local rows =
      [
        row.new('Overview')
        .addPanel(fleetTable { span: 12, height: '800px' }),
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Fleet Overview' % config.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='node-fleet'
      ) { editable: true }
      .addLink(c.links.otherDashes { includeVars: false })
      .addTemplates(templates)
      .addRows(rows)
    else if platform == 'Darwin' then {},
  },
}
