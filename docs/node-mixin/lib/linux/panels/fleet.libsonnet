local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
local xtd = import 'github.com/jsonnet-libs/xtd/main.libsonnet';
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = xtd.array.slice(this.config.instanceLabels, -1)[0],
      fleetOverviewTable:
        commonlib.panels.generic.table.base.new(
          'Fleet overview',
          targets=
          [
            t.system.osInfoCombined
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('OS Info'),
            t.system.uptime
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Uptime'),
            t.system.systemLoad1
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Load 1'),
            t.cpu.cpuCount
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Cores'),
            t.cpu.cpuUsage
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('CPU usage'),
            t.memory.memoryTotalBytes
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Memory total'),
            t.memory.memoryUsagePercent
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Memory usage'),
            t.disk.diskTotalRoot
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Root mount size'),
            t.disk.diskUsageRootPercent
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Root mount used'),
            t.alerts.alertsCritical
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('CRITICAL'),
            t.alerts.alertsWarning
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('WARNING'),
          ],
          description="All nodes' perfomance at a glance."
        )
        + g.panel.table.options.withFooter(
          value={
            reducer: ['sum'],
            show: true,
            fields: [
              'Value #Cores',
              'Value #Load 1',
              'Value #Memory total',
              'Value #Root mount size',
            ],
          }
        )
        + commonlib.panels.system.table.uptime.stylizeByName('Uptime')
        + table.standardOptions.withOverridesMixin([
          fieldOverride.byRegexp.new('Product|^Hostname$')
          + fieldOverride.byRegexp.withProperty('custom.filterable', true),
          fieldOverride.byName.new('Instance')
          + fieldOverride.byName.withProperty('custom.filterable', true)
          + fieldOverride.byName.withProperty('links', [
            {
              targetBlank: false,
              title: 'Drill down to ${__field.name} ${__value.text}',
              url: 'd/%s?var-%s=${__data.fields.%s}&${__url_time_range}&${datasource:queryparam}' % [this.grafana.dashboards['nodes.json'].uid, instanceLabel, instanceLabel],
            },
          ]),
          fieldOverride.byRegexp.new(std.join('|', std.map(utils.toSentenceCase, this.config.groupLabels)))
          + fieldOverride.byRegexp.withProperty('custom.filterable', true)
          + fieldOverride.byRegexp.withProperty('links', [
            {
              targetBlank: false,
              title: 'Filter by ${__field.name}',
              url: 'd/%s?var-${__field.name}=${__value.text}&${__url_time_range}&${datasource:queryparam}' % [this.grafana.dashboards['fleet.json'].uid],
            },
          ]),
          fieldOverride.byName.new('Cores')
          + fieldOverride.byName.withProperty('custom.width', '120'),
          fieldOverride.byName.new('CPU usage')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withProperty(
            'custom.cellOptions', {
              type: 'gauge',
              mode: 'basic',
              valueDisplayMode: 'text',
            }
          )
          + fieldOverride.byName.withPropertiesFromOptions(
            commonlib.panels.cpu.timeSeries.utilization.stylize()
          ),
          fieldOverride.byName.new('Memory total')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withPropertiesFromOptions(
            table.standardOptions.withUnit('bytes')
          ),
          fieldOverride.byName.new('Memory usage')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withProperty(
            'custom.cellOptions', {
              type: 'gauge',
              mode: 'basic',
              valueDisplayMode: 'text',
            }
          )
          + fieldOverride.byName.withPropertiesFromOptions(
            commonlib.panels.cpu.timeSeries.utilization.stylize()
          ),
          fieldOverride.byName.new('Root mount size')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withPropertiesFromOptions(
            table.standardOptions.withUnit('bytes')
          ),
          fieldOverride.byName.new('Root mount used')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withProperty(
            'custom.cellOptions', {
              type: 'gauge',
              mode: 'basic',
              valueDisplayMode: 'text',
            }
          )
          + fieldOverride.byName.withPropertiesFromOptions(
            table.standardOptions.withUnit('percent')
          )
          + fieldOverride.byName.withPropertiesFromOptions(
            commonlib.panels.cpu.timeSeries.utilization.stylize()
          ),
        ])
        + table.queryOptions.withTransformationsMixin(
          [
            {
              id: 'joinByField',
              options: {
                byField: instanceLabel,
                mode: 'outer',
              },
            },
            {
              id: 'filterFieldsByName',
              options: {
                include: {
                  //' 1' - would only match first occurence of group label, so no duplicates
                  pattern: instanceLabel + '|'
                           +
                           std.join(
                             '|',
                             std.map(
                               function(x) '%s 1' % x, this.config.instanceLabels
                             )
                           )
                           + '|' +
                           std.join(
                             '|',
                             std.map(
                               function(x) '%s 1' % x, this.config.groupLabels
                             )
                           )
                           + '|product|^hostname$|^nodename$|^pretty_name$|Value.+',
                },
              },
            },
            {
              id: 'organize',
              options: {
                excludeByName: {
                  'Value #OS Info': true,
                },
                indexByName:
                  {
                    [instanceLabel]: 0,
                    nodename: 1,
                    hostname: 1,
                    pretty_name: 2,
                    product: 2,
                  }
                  +
                  // group labels are named as 'job 1' and so on.
                  {
                    [label]: 3
                    for label in this.config.groupLabels
                  },
                renameByName:
                  {
                    [label + ' 1']: utils.toSentenceCase(label)
                    for label in this.config.instanceLabels
                  }
                  {
                    [instanceLabel]: utils.toSentenceCase(instanceLabel),
                    product: 'OS',  // windows
                    pretty_name: 'OS',  // linux
                    hostname: 'Hostname',  // windows
                    nodename: 'Hostname',  // Linux
                  }
                  +
                  // group labels are named as 'job 1' and so on.
                  {
                    [label + ' 1']: utils.toSentenceCase(label)
                    for label in this.config.groupLabels
                  },

              },
            },
            {
              id: 'renameByRegex',
              options: {
                regex: 'Value #(.*)',
                renamePattern: '$1',
              },
            },
          ]
        ),
    },
}
