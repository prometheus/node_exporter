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

    local networkInterfacesTable =
      nodePanels.table.new(
        title='Linux Nodes Overview'
      )
      // "Value #A"
      .addTarget(commonPromTarget(
        expr='avg by (instance) (node_load1{%(nodeExporterSelector)s, instance="$instance"})' % config,
        format='table',
        instant=true,
      ))
      // "Value #B" (number of cpu)
      .addTarget(commonPromTarget(
        expr='count by (instance)(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance"})' % config,
        format='table',
        instant=true,
      ))
      // "Value #C (memory usage)"
      .addTarget(commonPromTarget(
        expr=
        |||
          100 -
          (
            avg by (instance) (node_memory_MemAvailable_bytes{%(nodeExporterSelector)s, instance=~"$instance"}) /
            avg by (instance) (node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance=~"$instance"})
          * 100
          )
        ||| % config,
        format='table',
        instant=true,
      ))
      // "Value #D (mount / usage)"
      .addTarget(commonPromTarget(
        expr=
        |||
          100-(max by (instance) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance=~"$instance", fstype!="", mountpoint="/"})
          /
          max by (instance) (node_filesystem_size_bytes{%(nodeExporterSelector)s, instance=~"$instance", fstype!="", mountpoint="/"}) * 100)
        ||| % config,
        format='table',
        instant=true,
      ))
      .withTransform()
      .joinByField(field='instance')
      // .merge()
      .filterFieldsByName('instance|Value.+')
      .organize(
        excludeByName={
        },
        renameByName=
        {
          instance: 'Instance',
          'Value #A': 'Load average 1',
          'Value #B': 'Cores',
          'Value #C': 'Memory usage',
          'Value #D': 'Root / disk usage',
        }
      )
      .withFooter(reducer=['mean'], fields=['Value #C'])
      .addThresholdStep(color='light-green', value=null)
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
                title: 'Drill down to instance',
                url: 'd/nodes?var-instance=${__data.fields.instance}&${__url_time_range}',
              },
            ],
          },
        ]
      )
      .addOverride(
        matcher={
          id: 'byRegexp',
          options: 'Memory usage|Root / disk usage',
        },
        properties=[
          {
            id: 'unit',
            value: 'percent',
          },
          {
            id: 'custom.displayMode',
            value: 'gradient-gauge',
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
      ),
    local rows =
      [
        row.new('Overview')
        .addPanel(networkInterfacesTable { span: 12, height: '800px' }),
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
