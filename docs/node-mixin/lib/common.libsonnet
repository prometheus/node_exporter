local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local nodePanels = import '../lib/panels/panels.libsonnet';
local commonPanels = import '../lib/panels/common/panels.libsonnet';
local nodeTimeseries = nodePanels.timeseries;
{

  new(config=null, platform=null):: {

    local c = self,

    local labelsToRegexSelector(labels) =
      std.join(',', ['%s=~"$%s"' % [label, label] for label in labels]),
    local labelsToLegend(labels) =
      std.join('/', ['{{%s}}' % [label] for label in labels]),

    local labelsToURLvars(labels, prefix) =
      std.join('&', ['var-%s=${%s%s}' % [label, prefix, label] for label in labels]),
    // export
    labelsToLegend:: labelsToLegend,
    labelsToURLvars:: labelsToURLvars,
    // add to all queries but not templates
    local nodeQuerySelector = labelsToRegexSelector(std.split(config.groupLabels + ',' + config.instanceLabels, ',')),
    nodeQuerySelector:: nodeQuerySelector,

    // common templates
    local prometheusDatasourceTemplate = {
      current: {
        text: 'default',
        value: 'default',
      },
      hide: 0,
      label: 'Data Source',
      name: 'datasource',
      options: [],
      query: 'prometheus',
      refresh: 1,
      regex: '',
      type: 'datasource',
    },

    local chainLabelsfold(prev, label) = {
      chain:
        if std.length(prev) > 0
        then
          [[label] + prev.chain[0]] + prev.chain
        else
          [[label]],
    },

    local chainLabels(labels) =
      [
        {
          label: l[0:1][0],
          chainSelector: labelsToRegexSelector(std.reverse(l[1:])),
        }
        for l in std.reverse(std.foldl(chainLabelsfold, labels, init={}).chain)
      ],

    local groupTemplates =
      [
        template.new(
          name=label.label,
          label=label.label,
          datasource='$datasource',
          query='',
          current='',
          refresh=2,
          includeAll=true,
          // do not use .*, will get series without such label at all when ALL is selected, ignoring nodeExporterSelector results
          allValues=null,
          multi=true,
          sort=1
        )
        {
          query: if platform == 'Darwin' then 'label_values(node_uname_info{sysname="Darwin", %(nodeExporterSelector)s, %(chainSelector)s}, %(label)s)' % config { label: label.label, chainSelector: label.chainSelector }
          else 'label_values(node_uname_info{sysname!="Darwin", %(nodeExporterSelector)s, %(chainSelector)s}, %(label)s)' % config { label: label.label, chainSelector: label.chainSelector },
        }
        for label in chainLabels(std.split(config.groupLabels, ','))
      ],

    local instanceTemplates =
      [
        template.new(
          label.label,
          '$datasource',
          'label_values(node_uname_info{%(nodeExporterSelector)s, %(chainSelector)s}, %(label)s)' % config { label: label.label, chainSelector: labelsToRegexSelector(std.split(config.groupLabels, ',')) + ',' + label.chainSelector },
          sort=1,
          refresh='time',
          label=label.label,
        )
        for label in chainLabels(std.split(config.instanceLabels, ','))
      ],

    // return common templates
    templates: [prometheusDatasourceTemplate] + groupTemplates + instanceTemplates,
    // return templates where instance select is not required
    groupDashboardTemplates: [prometheusDatasourceTemplate] + groupTemplates,

    local rebootAnnotation = {
      datasource: {
        type: 'prometheus',
        uid: '$datasource',
      },
      enable: true,
      hide: true,
      expr: 'node_boot_time_seconds{%(nodeQuerySelector)s}*1000 > $__from < $__to' % config { nodeQuerySelector: nodeQuerySelector },
      name: 'Reboot',
      iconColor: 'light-orange',
      tagKeys: config.instanceLabels,
      textFormat: '',
      titleFormat: 'Reboot',
      useValueForTime: 'on',
    },
    local memoryOOMkillerAnnotation = {
      datasource: {
        type: 'prometheus',
        uid: '$datasource',
      },
      enable: true,
      hide: true,
      expr: 'increase(node_vmstat_oom_kill{%(nodeQuerySelector)s}[$__interval])' % config { nodeQuerySelector: nodeQuerySelector },
      name: 'OOMkill',
      iconColor: 'light-purple',
      tagKeys: config.instanceLabels,
      textFormat: '',
      titleFormat: 'OOMkill',
    },
    local newKernelAnnotation = {
      datasource: {
        type: 'prometheus',
        uid: '$datasource',
      },
      enable: true,
      hide: true,
      expr: |||
        changes(
        sum by (%(instanceLabels)s) (
            group by (%(instanceLabels)s,release) (node_uname_info{%(nodeQuerySelector)s})
            )
        [$__interval:1m] offset -$__interval) > 1
      ||| % config { nodeQuerySelector: nodeQuerySelector },
      name: 'Kernel update',
      iconColor: 'light-blue',
      tagKeys: config.instanceLabels,
      textFormat: '',
      titleFormat: 'Kernel update',
      step: '5m',  // must be larger than possible scrape periods
    },
    // return common annotations
    annotations: [rebootAnnotation, memoryOOMkillerAnnotation, newKernelAnnotation],

    // return common prometheus target (with project defaults)
    commonPromTarget(
      expr=null,
      intervalFactor=1,
      datasource='$datasource',
      legendFormat=null,
      format='timeseries',
      instant=null,
      hide=null,
      interval=null,
    )::
      prometheus.target(
        expr=expr,
        intervalFactor=intervalFactor,
        datasource=datasource,
        legendFormat=legendFormat,
        format=format,
        instant=instant,
        hide=hide,
        interval=interval
      ),
    // link to fleet panel
    links:: {
      fleetDash:: grafana.link.dashboards(
        asDropdown=false,
        title='Back to Node Fleet Overview',
        tags=[],
        includeVars=false,
        keepTime=true,
        url='d/' + config.grafanaDashboardIDs['nodes-fleet.json']
      ) { type: 'link', icon: 'dashboard' },
      nodeDash:: grafana.link.dashboards(
        asDropdown=false,
        title='Back to Node Overview',
        tags=[],
        includeVars=true,
        keepTime=true,
        url='d/' + config.grafanaDashboardIDs['nodes.json']
      ) { type: 'link', icon: 'dashboard' },
      otherDashes:: grafana.link.dashboards(
        asDropdown=true,
        title='Other Node Dashboards',
        includeVars=true,
        keepTime=true,
        tags=(config.dashboardTags),
      ),
      // used in fleet table
      instanceDataLinkForTable:: {
        title: 'Drill down to instance ${__data.fields.%s}' % std.split(config.instanceLabels, ',')[0],
        url: 'd/' + config.grafanaDashboardIDs['nodes.json'] + '?' + labelsToURLvars(std.split(config.instanceLabels, ','), prefix='__data.fields.') + '&${__url_time_range}&var-datasource=${datasource}',
      },
      // used in ts panels
      instanceDataLink:: {
        title: 'Drill down to instance ${__field.labels.%s}' % std.split(config.instanceLabels, ',')[0],
        url: 'd/' + config.grafanaDashboardIDs['nodes.json'] + '?' + labelsToURLvars(std.split(config.instanceLabels, ','), prefix='__field.labels.') + '&${__url_time_range}&var-datasource=${datasource}',
      },
    },
    // return common queries that could be used in multiple dashboards
    queries:: {
      systemLoad1:: 'avg by (%(instanceLabels)s) (node_load1{%(nodeQuerySelector)s})' % config { nodeQuerySelector: nodeQuerySelector },
      systemLoad5:: 'avg by (%(instanceLabels)s) (node_load5{%(nodeQuerySelector)s})' % config { nodeQuerySelector: nodeQuerySelector },
      systemLoad15:: 'avg by (%(instanceLabels)s) (node_load15{%(nodeQuerySelector)s})' % config { nodeQuerySelector: nodeQuerySelector },
      uptime:: 'time() - node_boot_time_seconds{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      cpuCount:: 'count by (%(instanceLabels)s) (node_cpu_seconds_total{%(nodeQuerySelector)s, mode="idle"})' % config { nodeQuerySelector: nodeQuerySelector },
      cpuUsage::
        |||
          (((count by (%(instanceLabels)s) (count(node_cpu_seconds_total{%(nodeQuerySelector)s}) by (cpu, %(instanceLabels)s))) 
          - 
          avg by (%(instanceLabels)s) (sum by (%(instanceLabels)s, mode)(irate(node_cpu_seconds_total{mode='idle',%(nodeQuerySelector)s}[$__rate_interval])))) * 100) 
          / 
          count by(%(instanceLabels)s) (count(node_cpu_seconds_total{%(nodeQuerySelector)s}) by (cpu, %(instanceLabels)s))
        ||| % config { nodeQuerySelector: nodeQuerySelector },
      cpuUsageModes::
        |||
          sum by(%(instanceLabels)s, mode) (irate(node_cpu_seconds_total{%(nodeQuerySelector)s}[$__rate_interval])) 
          / on(%(instanceLabels)s) 
          group_left sum by (%(instanceLabels)s)((irate(node_cpu_seconds_total{%(nodeQuerySelector)s}[$__rate_interval]))) * 100
        ||| % config { nodeQuerySelector: nodeQuerySelector },
      cpuUsagePerCore::
        |||
          (
            (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeQuerySelector)s, mode=~"idle|iowait|steal"}[$__rate_interval])))
          / ignoring(cpu) group_left
            count without (cpu, mode) (node_cpu_seconds_total{%(nodeQuerySelector)s, mode="idle"})
          ) * 100
        ||| % config { nodeQuerySelector: nodeQuerySelector },
      memoryTotal:: 'node_memory_MemTotal_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      memorySwapTotal:: 'node_memory_SwapTotal_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      memoryUsage::
        |||
          100 -
          (
            avg by (%(instanceLabels)s) (node_memory_MemAvailable_bytes{%(nodeQuerySelector)s}) /
            avg by (%(instanceLabels)s) (node_memory_MemTotal_bytes{%(nodeQuerySelector)s})
          * 100
          )
        ||| % config { nodeQuerySelector: nodeQuerySelector },

      process_max_fds:: 'process_max_fds{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      process_open_fds:: 'process_open_fds{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },

      fsSizeTotalRoot:: 'node_filesystem_size_bytes{%(nodeQuerySelector)s, mountpoint="/",fstype!="rootfs"}' % config { nodeQuerySelector: nodeQuerySelector },
      osInfo:: 'node_os_info{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      nodeInfo:: 'node_uname_info{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_disk_reads_completed_total:: 'irate(node_disk_reads_completed_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      node_disk_writes_completed_total:: 'irate(node_disk_writes_completed_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      diskReadTime:: 'rate(node_disk_read_bytes_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      diskWriteTime:: 'rate(node_disk_written_bytes_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      diskIoTime:: 'rate(node_disk_io_time_seconds_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      diskWaitReadTime::
        |||
          irate(node_disk_read_time_seconds_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])
          /
          irate(node_disk_reads_completed_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])
        ||| % config { nodeQuerySelector: nodeQuerySelector },
      diskWaitWriteTime::
        |||
          irate(node_disk_write_time_seconds_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])
          /
          irate(node_disk_writes_completed_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])
        ||| % config { nodeQuerySelector: nodeQuerySelector },
      diskAvgQueueSize:: 'irate(node_disk_io_time_weighted_seconds_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      diskSpaceUsage::
        |||
          sort_desc(1 -
            (
            max by (job, %(instanceLabels)s, fstype, device, mountpoint) (node_filesystem_avail_bytes{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s})
            /
            max by (job, %(instanceLabels)s, fstype, device, mountpoint) (node_filesystem_size_bytes{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s})
            ) != 0
          )
        ||| % config { nodeQuerySelector: nodeQuerySelector },
      node_filesystem_avail_bytes:: 'node_filesystem_avail_bytes{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_filesystem_files_free:: 'node_filesystem_files_free{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_filesystem_files:: 'node_filesystem_files{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_filesystem_readonly:: 'node_filesystem_readonly{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_filesystem_device_error:: 'node_filesystem_device_error{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      networkReceiveBitsPerSec:: 'irate(node_network_receive_bytes_total{%(nodeQuerySelector)s}[$__rate_interval])*8' % config { nodeQuerySelector: nodeQuerySelector },
      networkTransmitBitsPerSec:: 'irate(node_network_transmit_bytes_total{%(nodeQuerySelector)s}[$__rate_interval])*8' % config { nodeQuerySelector: nodeQuerySelector },
      networkReceiveErrorsPerSec:: 'irate(node_network_receive_errs_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      networkTransmitErrorsPerSec:: 'irate(node_network_transmit_errs_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      networkReceiveDropsPerSec:: 'irate(node_network_receive_drop_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      networkTransmitDropsPerSec:: 'irate(node_network_transmit_drop_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },

      systemContextSwitches:: 'irate(node_context_switches_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      systemInterrupts:: 'irate(node_intr_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },

      //time
      node_timex_estimated_error_seconds:: 'node_timex_estimated_error_seconds{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_timex_offset_seconds:: 'node_timex_offset_seconds{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_timex_maxerror_seconds:: 'node_timex_maxerror_seconds{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },

      node_timex_sync_status:: 'node_timex_sync_status{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_time_zone_offset_seconds:: 'node_time_zone_offset_seconds{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      node_systemd_units:: 'node_systemd_units{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },


    },
    // share across dashboards
    panelsWithTargets:: {
      // cpu
      idleCPU::
        nodePanels.timeseries.new(
          'CPU Usage',
          description='Total CPU utilisation percent.'
        )
        .withUnits('percent')
        .withStacking('normal')
        .withMin(0)
        .withMax(100)
        .addTarget(c.commonPromTarget(
          expr=c.queries.cpuUsagePerCore,
          legendFormat='cpu {{cpu}}',
        )),

      systemLoad::
        nodePanels.timeseries.new(
          'Load Average',
          description='System load average over the previous 1, 5, and 15 minute ranges. A measurement of how many processes are waiting for CPU cycles. The maximum number is the number of CPU cores for the node.',
        )
        .withUnits('short')
        .withMin(0)
        .withFillOpacity(0)
        .addTarget(c.commonPromTarget(c.queries.systemLoad1, legendFormat='1m load average'))
        .addTarget(c.commonPromTarget(c.queries.systemLoad5, legendFormat='5m load average'))
        .addTarget(c.commonPromTarget(c.queries.systemLoad15, legendFormat='15m load average'))
        .addTarget(c.commonPromTarget(c.queries.cpuCount, legendFormat='logical cores'))
        .addOverride(
          matcher={
            id: 'byName',
            options: 'logical cores',
          },
          properties=[
            {
              id: 'custom.lineStyle',
              value: {
                fill: 'dash',
                dash: [
                  10,
                  10,
                ],
              },
            },
          ]
        ),
      cpuStatPanel::
        commonPanels.percentUsageStat.new(
          'CPU Usage',
          description='Total CPU utilisation percent.'
        )
        .addTarget(c.commonPromTarget(
          expr=c.queries.cpuUsage
        )),
      systemContextSwitches::
        nodePanels.timeseries.new(
          'Context Switches / Interrupts',
          description=|||
            Context switches occur when the operating system switches from running one process to another.
            Interrupts are signals sent to the CPU by external devices to request its attention.

            A high number of context switches or interrupts can indicate that the system is overloaded or that there are problems with specific devices or processes.
          |||
        )
        .addTarget(c.commonPromTarget(c.queries.systemContextSwitches, legendFormat='Context Switches'))
        .addTarget(c.commonPromTarget(c.queries.systemInterrupts, legendFormat='Interrupts')),

      diskSpaceUsage::
        nodePanels.table.new(
          title='Disk Space Usage',
          description='Disk utilisation in percent, by mountpoint. Some duplication can occur if the same filesystem is mounted in multiple locations.',
        )
        .setFieldConfig(unit='decbytes')
        //.addThresholdStep(color='light-green', value=null)
        .addThresholdStep(color='light-blue', value=null)
        .addThresholdStep(color='light-yellow', value=0.8)
        .addThresholdStep(color='light-red', value=0.9)
        .addTarget(c.commonPromTarget(
          |||
            max by (mountpoint) (node_filesystem_size_bytes{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s})
          ||| % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat='',
          instant=true,
          format='table'
        ))
        .addTarget(c.commonPromTarget(
          |||
            max by (mountpoint) (node_filesystem_avail_bytes{%(nodeQuerySelector)s, %(fsSelector)s, %(fsMountpointSelector)s})
          ||| % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat='',
          instant=true,
          format='table',
        ))
        .addOverride(
          matcher={
            id: 'byName',
            options: 'Mounted on',
          },
          properties=[
            {
              id: 'custom.width',
              value: 260,
            },
          ],
        )
        .addOverride(
          matcher={
            id: 'byName',
            options: 'Size',
          },
          properties=[

            {
              id: 'custom.width',
              value: 93,
            },

          ],
        )
        .addOverride(
          matcher={
            id: 'byName',
            options: 'Used',
          },
          properties=[
            {
              id: 'custom.width',
              value: 72,
            },
          ],
        )
        .addOverride(
          matcher={
            id: 'byName',
            options: 'Available',
          },
          properties=[
            {
              id: 'custom.width',
              value: 88,
            },
          ],
        )

        .addOverride(
          matcher={
            id: 'byName',
            options: 'Used, %',
          },
          properties=[
            {
              id: 'unit',
              value: 'percentunit',
            },
            {
              id: 'custom.displayMode',
              value: 'basic',
            },
            {
              id: 'max',
              value: 1,
            },
            {
              id: 'min',
              value: 0,
            },
          ]
        )
        .sortBy('Mounted on')
        + {
          transformations+: [
            {
              id: 'groupBy',
              options: {
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
              },
            },
            {
              id: 'merge',
              options: {},
            },
            {
              id: 'calculateField',
              options: {
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
              },
            },
            {
              id: 'calculateField',
              options: {
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
              },
            },
            {
              id: 'organize',
              options: {
                excludeByName: {},
                indexByName: {},
                renameByName: {
                  'Value #A (lastNotNull)': 'Size',
                  'Value #B (lastNotNull)': 'Available',
                  mountpoint: 'Mounted on',
                },
              },
            },
          ],
        },
      memoryGraphPanelPrototype::
        nodePanels.timeseries.new(
          'Memory Usage',
          description='Memory usage by category, measured in bytes.',
        )
        .withMin(0)
        .withUnits('bytes'),
      memoryGraph::
        if platform == 'Linux' then
          self.memoryGraphPanelPrototype
          {
            description: |||
              Used: The amount of physical memory currently in use by the system.
              Cached: The amount of physical memory used for caching data from disk. The Linux kernel uses available memory to cache data that is read from or written to disk. This helps speed up disk access times.
              Free: The amount of physical memory that is currently not in use.
              Buffers: The amount of physical memory used for temporary storage of data being transferred between devices or applications.
              Available: The amount of physical memory that is available for use by applications. This takes into account memory that is currently being used for caching but can be freed up if needed.
            |||,
          }
          { stack: true }
          .addTarget(c.commonPromTarget(
            |||
              (
                node_memory_MemTotal_bytes{%(nodeQuerySelector)s}
              -
                node_memory_MemFree_bytes{%(nodeQuerySelector)s}
              -
                node_memory_Buffers_bytes{%(nodeQuerySelector)s}
              -
                node_memory_Cached_bytes{%(nodeQuerySelector)s}
              )
            ||| % config { nodeQuerySelector: c.nodeQuerySelector },
            legendFormat='Memory used'
          ))
          .addTarget(c.commonPromTarget('node_memory_Buffers_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Memory buffers'))
          .addTarget(c.commonPromTarget('node_memory_Cached_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Memory cached'))
          .addTarget(c.commonPromTarget('node_memory_MemFree_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Memory free'))
          .addTarget(c.commonPromTarget('node_memory_MemAvailable_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Memory available'))
          .addTarget(c.commonPromTarget('node_memory_MemTotal_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Memory total'))
        else if platform == 'Darwin' then
          // not useful to stack
          self.memoryGraphPanelPrototype { stack: false }
          .addTarget(c.commonPromTarget('node_memory_total_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Physical Memory'))
          .addTarget(c.commonPromTarget(
            |||
              (
                  node_memory_internal_bytes{%(nodeQuerySelector)s} -
                  node_memory_purgeable_bytes{%(nodeQuerySelector)s} +
                  node_memory_wired_bytes{%(nodeQuerySelector)s} +
                  node_memory_compressed_bytes{%(nodeQuerySelector)s}
              )
            ||| % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Memory Used'
          ))
          .addTarget(c.commonPromTarget(
            |||
              (
                  node_memory_internal_bytes{%(nodeQuerySelector)s} -
                  node_memory_purgeable_bytes{%(nodeQuerySelector)s}
              )
            ||| % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='App Memory'
          ))
          .addTarget(c.commonPromTarget('node_memory_wired_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Wired Memory'))
          .addTarget(c.commonPromTarget('node_memory_compressed_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector }, legendFormat='Compressed')),

      // NOTE: avg() is used to circumvent a label change caused by a node_exporter rollout.
      memoryGaugePanelPrototype::
        commonPanels.percentUsageStat.new(
          'Memory Usage',
          description='Total memory utilisation.',
        ),

      memoryGauge::
        if platform == 'Linux' then
          self.memoryGaugePanelPrototype

          .addTarget(c.commonPromTarget(c.queries.memoryUsage))

        else if platform == 'Darwin' then
          self.memoryGaugePanelPrototype
          .addTarget(c.commonPromTarget(
            |||
              (
                  (
                    avg(node_memory_internal_bytes{%(nodeQuerySelector)s}) -
                    avg(node_memory_purgeable_bytes{%(nodeQuerySelector)s}) +
                    avg(node_memory_wired_bytes{%(nodeQuerySelector)s}) +
                    avg(node_memory_compressed_bytes{%(nodeQuerySelector)s})
                  ) /
                  avg(node_memory_total_bytes{%(nodeQuerySelector)s})
              )
              *
              100
            ||| % config { nodeQuerySelector: c.nodeQuerySelector }
          )),
      diskIO::
        nodePanels.timeseries.new(
          'Disk I/O',
          description='Disk read/writes in bytes, and total IO seconds.'
        )
        .withFillOpacity(0)
        .withMin(0)
        .addTarget(c.commonPromTarget(
          c.queries.diskReadTime,
          legendFormat='{{device}} read',
        ))
        .addTarget(c.commonPromTarget(
          c.queries.diskWriteTime,
          legendFormat='{{device}} written',
        ))
        .addTarget(c.commonPromTarget(
          c.queries.diskIoTime,
          legendFormat='{{device}} io time',
        ))
        .addOverride(
          matcher={
            id: 'byRegexp',
            options: '/ read| written/',
          },
          properties=[
            {
              id: 'unit',
              value: 'bps',
            },
          ]
        )
        .addOverride(
          matcher={
            id: 'byRegexp',
            options: '/ io time/',
          },
          properties=[
            {
              id: 'unit',
              value: 'percentunit',
            },
            {
              id: 'custom.axisSoftMax',
              value: 1,
            },
            {
              id: 'custom.drawStyle',
              value: 'points',
            },
          ]
        ),
    },
  },

}
