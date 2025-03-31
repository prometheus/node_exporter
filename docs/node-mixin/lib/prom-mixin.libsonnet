local grafana = import 'github.com/grafana/grafonnet/gen/grafonnet-latest/main.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.panel.row;
local prometheus = grafana.query.prometheus;
local variable = dashboard.variable;

local timeSeriesPanel = grafana.panel.timeSeries;
local tsOptions = timeSeriesPanel.options;
local tsStandardOptions = timeSeriesPanel.standardOptions;
local tsQueryOptions = timeSeriesPanel.queryOptions;
local tsCustom = timeSeriesPanel.fieldConfig.defaults.custom;

local gaugePanel = grafana.panel.gauge;
local gaugeStep = gaugePanel.standardOptions.threshold.step;

local table = grafana.panel.table;
local tableStep = table.standardOptions.threshold.step;
local tableOverride = table.standardOptions.override;
local tableTransformation = table.queryOptions.transformation;

{

  new(config=null, platform=null, uid=null):: {

    local prometheusDatasourceVariable = variable.datasource.new(
      'datasource', 'prometheus'
    ),

    local clusterVariablePrototype =
      variable.query.new('cluster')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + (if config.showMultiCluster then variable.query.generalOptions.showOnDashboard.withLabelAndValue() else variable.query.generalOptions.showOnDashboard.withNothing())
      + variable.query.refresh.onTime()
      + variable.query.generalOptions.withLabel('Cluster'),

    local clusterVariable =
      if platform == 'Darwin' then
        clusterVariablePrototype
        + variable.query.queryTypes.withLabelValues(
          ' %(clusterLabel)s' % config,
          'node_uname_info{%(nodeExporterSelector)s, sysname="Darwin"}' % config,
        )
      else
        clusterVariablePrototype
        + variable.query.queryTypes.withLabelValues(
          '%(clusterLabel)s' % config,
          'node_uname_info{%(nodeExporterSelector)s, sysname!="Darwin"}' % config,
        ),

    local instanceVariablePrototype =
      variable.query.new('instance')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + variable.query.refresh.onTime()
      + variable.query.generalOptions.withLabel('Instance'),

    local instanceVariable =
      if platform == 'Darwin' then
        instanceVariablePrototype
        + variable.query.queryTypes.withLabelValues(
          'instance',
          'node_uname_info{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster", sysname="Darwin"}' % config,
        )
      else
        instanceVariablePrototype
        + variable.query.queryTypes.withLabelValues(
          'instance',
          'node_uname_info{%(nodeExporterSelector)s, %(clusterLabel)s="$cluster", sysname!="Darwin"}' % config,
        ),

    local idleCPU =
      timeSeriesPanel.new('CPU Usage')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + tsStandardOptions.withUnit('percentunit')
      + tsCustom.stacking.withMode('normal')
      + tsStandardOptions.withMax(1)
      + tsStandardOptions.withMin(0)
      + tsOptions.tooltip.withMode('multi')
      + tsCustom.withFillOpacity(10)
      + tsCustom.withShowPoints('never')
      + tsQueryOptions.withTargets([
        prometheus.new(
          '$datasource',
          |||
            (
              (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode=~"idle|iowait|steal", instance="$instance", %(clusterLabel)s="$cluster"}[$__rate_interval])))
            / ignoring(cpu) group_left
              count without (cpu, mode) (node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance", %(clusterLabel)s="$cluster"})
            )
          ||| % config,
        )
        + prometheus.withLegendFormat('{{cpu}}')
        + prometheus.withIntervalFactor(5),
      ]),

    local systemLoad =
      timeSeriesPanel.new('Load Average')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + tsStandardOptions.withUnit('short')
      + tsStandardOptions.withMin(0)
      + tsCustom.withFillOpacity(0)
      + tsCustom.withShowPoints('never')
      + tsOptions.tooltip.withMode('multi')
      + tsQueryOptions.withTargets([
        prometheus.new('$datasource', 'node_load1{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('1m load average'),
        prometheus.new('$datasource', 'node_load5{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('5m load average'),
        prometheus.new('$datasource', 'node_load15{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('15m load average'),
        prometheus.new('$datasource', 'count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", mode="idle"})' % config) + prometheus.withLegendFormat('logical cores'),
      ]),

    local memoryGraphPanelPrototype =
      timeSeriesPanel.new('Memory Usage')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + tsStandardOptions.withUnit('bytes')
      + tsStandardOptions.withMin(0)
      + tsOptions.tooltip.withMode('multi')
      + tsCustom.withFillOpacity(10)
      + tsCustom.withShowPoints('never'),

    local memoryGraph =
      if platform == 'Linux' then
        memoryGraphPanelPrototype
        + tsCustom.stacking.withMode('normal')
        + tsQueryOptions.withTargets([
          prometheus.new(
            '$datasource',
            |||
              (
                node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
              -
                node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
              -
                node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
              -
                node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
              )
            ||| % config,
          ) + prometheus.withLegendFormat('memory used'),
          prometheus.new('$datasource', 'node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('memory buffers'),
          prometheus.new('$datasource', 'node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('memory cached'),
          prometheus.new('$datasource', 'node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('memory free'),
        ])
      else if platform == 'Darwin' then
        // not useful to stack
        memoryGraphPanelPrototype
        + tsCustom.stacking.withMode('none')
        + tsQueryOptions.withTargets([
          prometheus.new('$datasource', 'node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('Physical Memory'),
          prometheus.new(
            '$datasource',
            |||
              (
                  node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} -
                  node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} +
                  node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} +
                  node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
              )
            ||| % config
          ) + prometheus.withLegendFormat(
            'Memory Used'
          ),
          prometheus.new(
            '$datasource',
            |||
              (
                  node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} -
                  node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
              )
            ||| % config
          ) + prometheus.withLegendFormat(
            'App Memory'
          ),
          prometheus.new('$datasource', 'node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('Wired Memory'),
          prometheus.new('$datasource', 'node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('Compressed'),
        ])

      else if platform == 'AIX' then
        memoryGraphPanelPrototype
        + tsCustom.stacking.withMode('none')
        + tsQueryOptions.withTargets([
          prometheus.new('$datasource', 'node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}' % config) + prometheus.withLegendFormat('Physical Memory'),
          prometheus.new(
            '$datasource',
            |||
              (
                  node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"} -
                  node_memory_available_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}
              )
            ||| % config
          ) + prometheus.withLegendFormat('Memory Used'),
        ]),


    // NOTE: avg() is used to circumvent a label change caused by a node_exporter rollout.
    local memoryGaugePanelPrototype =
      gaugePanel.new('Memory Usage')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + gaugePanel.standardOptions.thresholds.withSteps([
        gaugeStep.withColor('rgba(50, 172, 45, 0.97)'),
        gaugeStep.withColor('rgba(237, 129, 40, 0.89)') + gaugeStep.withValue(80),
        gaugeStep.withColor('rgba(245, 54, 54, 0.9)') + gaugeStep.withValue(90),
      ])
      + gaugePanel.standardOptions.withMax(100)
      + gaugePanel.standardOptions.withMin(0)
      + gaugePanel.standardOptions.withUnit('percent'),

    local memoryGauge =
      if platform == 'Linux' then
        memoryGaugePanelPrototype
        + gaugePanel.queryOptions.withTargets([
          prometheus.new(
            '$datasource',
            |||
              100 -
              (
                avg(node_memory_MemAvailable_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) /
                avg(node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"})
              * 100
              )
            ||| % config,
          ),
        ])

      else if platform == 'Darwin' then
        memoryGaugePanelPrototype
        + gaugePanel.queryOptions.withTargets([
          prometheus.new(
            '$datasource',
            |||
              (
                  (
                    avg(node_memory_internal_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) -
                    avg(node_memory_purgeable_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) +
                    avg(node_memory_wired_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) +
                    avg(node_memory_compressed_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"})
                  ) /
                  avg(node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"})
              )
              *
              100
            ||| % config
          ),
        ])

      else if platform == 'AIX' then
        memoryGaugePanelPrototype
        + gaugePanel.queryOptions.withTargets([
          prometheus.new(
            '$datasource',
            |||
              100 -
              (
                avg(node_memory_available_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"}) /
                avg(node_memory_total_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster"})
                * 100
              )
            ||| % config
          ),
        ]),


    local diskIO =
      timeSeriesPanel.new('Disk I/O')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + tsStandardOptions.withMin(0)
      + tsCustom.withFillOpacity(0)
      + tsCustom.withShowPoints('never')
      + tsOptions.tooltip.withMode('multi')
      + tsQueryOptions.withTargets([
        // TODO: Does it make sense to have those three in the same panel?
        prometheus.new('$datasource', 'rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(diskDeviceSelector)s}[$__rate_interval])' % config)
        + prometheus.withLegendFormat('{{device}} read')
        + prometheus.withIntervalFactor(1),
        prometheus.new('$datasource', 'rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(diskDeviceSelector)s}[$__rate_interval])' % config)
        + prometheus.withLegendFormat('{{device}} written')
        + prometheus.withIntervalFactor(1),
        prometheus.new('$datasource', 'rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(diskDeviceSelector)s}[$__rate_interval])' % config)
        + prometheus.withLegendFormat('{{device}} io time')
        + prometheus.withIntervalFactor(1),
      ])
      + tsStandardOptions.withOverrides(
        [
          tsStandardOptions.override.byRegexp.new('/ read| written/')
          + tsStandardOptions.override.byRegexp.withPropertiesFromOptions(
            tsStandardOptions.withUnit('Bps')
          ),
          tsStandardOptions.override.byRegexp.new('/ io time/')
          + tsStandardOptions.override.byRegexp.withPropertiesFromOptions(tsStandardOptions.withUnit('percentunit')),
        ]
      ),

    local diskSpaceUsage =
      table.new('Disk Space Usage')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + table.standardOptions.withUnit('decbytes')
      + table.standardOptions.thresholds.withSteps(
        [
          tableStep.withColor('green'),
          tableStep.withColor('yellow') + gaugeStep.withValue(0.8),
          tableStep.withColor('red') + gaugeStep.withValue(0.9),
        ]
      )
      + table.queryOptions.withTargets([
        prometheus.new(
          '$datasource',
          |||
            max by (mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(fsSelector)s, %(fsMountpointSelector)s})
          ||| % config
        )
        + prometheus.withLegendFormat('')
        + prometheus.withInstant()
        + prometheus.withFormat('table'),
        prometheus.new(
          '$datasource',
          |||
            max by (mountpoint) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", %(fsSelector)s, %(fsMountpointSelector)s})
          ||| % config
        )
        + prometheus.withLegendFormat('')
        + prometheus.withInstant()
        + prometheus.withFormat('table'),
      ])
      + table.standardOptions.withOverrides([
        tableOverride.byName.new('Mounted on')
        + tableOverride.byName.withProperty('custom.width', 260),
        tableOverride.byName.new('Size')
        + tableOverride.byName.withProperty('custom.width', 93),
        tableOverride.byName.new('Used')
        + tableOverride.byName.withProperty('custom.width', 72),
        tableOverride.byName.new('Available')
        + tableOverride.byName.withProperty('custom.width', 88),
        tableOverride.byName.new('Used, %')
        + tableOverride.byName.withProperty('unit', 'percentunit')
        + tableOverride.byName.withPropertiesFromOptions(
          table.fieldConfig.defaults.custom.withCellOptions(
            { type: 'gauge' },
          )
        )
        + tableOverride.byName.withProperty('max', 1)
        + tableOverride.byName.withProperty('min', 0),
      ])
      + table.queryOptions.withTransformations([
        tableTransformation.withId('groupBy')
        + tableTransformation.withOptions(
          {
            fields: {
              'Value #A': {
                aggregations: [
                  'lastNotNull',
                ],
                operation: 'aggregate',
              },
              'Value #B': {
                aggregations: [
                  'lastNotNull',
                ],
                operation: 'aggregate',
              },
              mountpoint: {
                aggregations: [],
                operation: 'groupby',
              },
            },
          }
        ),
        tableTransformation.withId('merge'),
        tableTransformation.withId('calculateField')
        + tableTransformation.withOptions(
          {
            alias: 'Used',
            binary: {
              left: 'Value #A (lastNotNull)',
              operator: '-',
              reducer: 'sum',
              right: 'Value #B (lastNotNull)',
            },
            mode: 'binary',
            reduce: {
              reducer: 'sum',
            },
          }
        ),
        tableTransformation.withId('calculateField')
        + tableTransformation.withOptions(
          {
            alias: 'Used, %',
            binary: {
              left: 'Used',
              operator: '/',
              reducer: 'sum',
              right: 'Value #A (lastNotNull)',
            },
            mode: 'binary',
            reduce: {
              reducer: 'sum',
            },
          }
        ),
        tableTransformation.withId('organize')
        + tableTransformation.withOptions(
          {
            excludeByName: {},
            indexByName: {},
            renameByName: {
              'Value #A (lastNotNull)': 'Size',
              'Value #B (lastNotNull)': 'Available',
              mountpoint: 'Mounted on',
            },
          }
        ),
        tableTransformation.withId('sortBy')
        + tableTransformation.withOptions(
          {
            fields: {},
            sort: [
              {
                field: 'Mounted on',
              },
            ],
          }
        ),

      ]),

    local networkReceived =
      timeSeriesPanel.new('Network Received')
      + timeSeriesPanel.panelOptions.withDescription('Network received (bits/s)')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + tsStandardOptions.withUnit('bps')
      + tsStandardOptions.withMin(0)
      + tsCustom.withFillOpacity(0)
      + tsCustom.withShowPoints('never')
      + tsOptions.tooltip.withMode('multi')
      + tsQueryOptions.withTargets([
        prometheus.new('$datasource', 'rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", device!="lo"}[$__rate_interval]) * 8' % config)
        + prometheus.withLegendFormat('{{device}}')
        + prometheus.withIntervalFactor(1),
      ]),

    local networkTransmitted =
      timeSeriesPanel.new('Network Transmitted')
      + timeSeriesPanel.panelOptions.withDescription('Network transmitted (bits/s)')
      + variable.query.withDatasourceFromVariable(prometheusDatasourceVariable)
      + tsStandardOptions.withUnit('bps')
      + tsStandardOptions.withMin(0)
      + tsCustom.withFillOpacity(0)
      + tsOptions.tooltip.withMode('multi')
      + tsQueryOptions.withTargets([
        prometheus.new('$datasource', 'rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(clusterLabel)s="$cluster", device!="lo"}[$__rate_interval]) * 8' % config)
        + prometheus.withLegendFormat('{{device}}')
        + prometheus.withIntervalFactor(1),
      ]),

    local cpuRow =
      row.new('CPU')
      + row.withPanels([
        idleCPU,
        systemLoad,
      ]),

    local memoryRow = [
      row.new('Memory') + row.gridPos.withY(8),
      memoryGraph + row.gridPos.withX(0) + row.gridPos.withY(9) + row.gridPos.withH(7) + row.gridPos.withW(18),
      memoryGauge + row.gridPos.withX(18) + row.gridPos.withY(9) + row.gridPos.withH(7) + row.gridPos.withW(6),
    ],

    local diskRow =
      row.new('Disk')
      + row.withPanels([
        diskIO,
        diskSpaceUsage,
      ]),

    local networkRow =
      row.new('Network')
      + row.withPanels([
        networkReceived,
        networkTransmitted,
      ]),

    local panels =
      grafana.util.grid.makeGrid([
        cpuRow,
      ], panelWidth=12, panelHeight=7)
      + memoryRow
      + grafana.util.grid.makeGrid([
        diskRow,
        networkRow,
      ], panelWidth=12, panelHeight=7, startY=18),

    local variables =
      [
        prometheusDatasourceVariable,
        clusterVariable,
        instanceVariable,
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNodes' % config.dashboardNamePrefix,
      )
      + dashboard.time.withFrom('now-1h')
      + dashboard.withTags(config.dashboardTags)
      + dashboard.withTimezone('utc')
      + dashboard.withRefresh('30s')
      + dashboard.withUid(std.md5(uid))
      + dashboard.graphTooltip.withSharedCrosshair()
      + dashboard.withVariables(variables)
      + dashboard.withPanels(panels)
    else if platform == 'Darwin' then
      dashboard.new(
        '%sMacOS' % config.dashboardNamePrefix,
      )
      + dashboard.time.withFrom('now-1h')
      + dashboard.withTags(config.dashboardTags)
      + dashboard.withTimezone('utc')
      + dashboard.withRefresh('30s')
      + dashboard.withUid(std.md5(uid))
      + dashboard.graphTooltip.withSharedCrosshair()
      + dashboard.withVariables(variables)
      + dashboard.withPanels(panels)
    else if platform == 'AIX' then
      dashboard.new(
        '%sAIX' % config.dashboardNamePrefix,
      )
      + dashboard.time.withFrom('now-1h')
      + dashboard.withTags(config.dashboardTags)
      + dashboard.withTimezone('utc')
      + dashboard.withRefresh('30s')
      + dashboard.withUid(std.md5(uid))
      + dashboard.graphTooltip.withSharedCrosshair()
      + dashboard.withVariables(variables)
      + dashboard.withPanels(panels),

  },
}
