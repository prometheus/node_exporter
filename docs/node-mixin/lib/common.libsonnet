local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;

{

  new(config=null, platform=null):: {


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
          // do not use .* will get series without such label at all when ALL is selected.
          // do not use .+ will ignore nodeExporterSelector results
          // use null for group values
          allValues=null,
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
      iconColor: 'orange',
      tagKeys: config.instanceLabels,
      textFormat: '',
      titleFormat: 'Reboot',
      useValueForTime: 'on',
    },
    // return common annotations
    annotations: [rebootAnnotation],

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
        url='d/node-fleet'
      ) { type: 'link', icon: 'dashboard' },
      nodeDash:: grafana.link.dashboards(
        asDropdown=false,
        title='Back to Node Overview',
        tags=[],
        includeVars=true,
        keepTime=true,
        url='d/nodes'
      ) { type: 'link', icon: 'dashboard' },
      otherDashes:: grafana.link.dashboards(
        asDropdown=true,
        title='Other Node dashboards',
        includeVars=true,
        keepTime=true,
        tags=(config.dashboardTags),
      ),
      // used in fleet table
      instanceDataLinkForTable:: {
        title: 'Drill down to instance ${__data.fields.%s}' % std.split(config.instanceLabels, ',')[0],
        url: 'd/nodes?' + labelsToURLvars(std.split(config.instanceLabels, ','), prefix='__data.fields.') + '&${__url_time_range}',
      },
      // used in ts panels
      instanceDataLink:: {
        title: 'Drill down to instance ${__field.labels.%s}' % std.split(config.instanceLabels, ',')[0],
        url: 'd/nodes?' + labelsToURLvars(std.split(config.instanceLabels, ','), prefix='__field.labels.') + '&${__url_time_range}',
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
          avg by (%(instanceLabels)s) (sum by (%(instanceLabels)s, mode)(irate(node_cpu_seconds_total{mode='idle',%(nodeExporterSelector)s}[5m])))) * 100) 
          / 
          count by(%(instanceLabels)s) (count(node_cpu_seconds_total{%(nodeQuerySelector)s}) by (cpu, %(instanceLabels)s))
        ||| % config { nodeQuerySelector: nodeQuerySelector },
      cpuUsagePerCore::
        |||
          (
            (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeQuerySelector)s, mode=~"idle|iowait|steal"}[$__rate_interval])))
          / ignoring(cpu) group_left
            count without (cpu, mode) (node_cpu_seconds_total{%(nodeQuerySelector)s, mode="idle"})
          )
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
      fsSizeTotalRoot:: 'node_filesystem_size_bytes{%(nodeQuerySelector)s, mountpoint="/",fstype!="rootfs"}' % config { nodeQuerySelector: nodeQuerySelector },
      osInfo:: 'node_os_info{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      nodeInfo:: 'node_uname_info{%(nodeQuerySelector)s}' % config { nodeQuerySelector: nodeQuerySelector },
      diskReadTime:: 'rate(node_disk_read_bytes_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      diskWriteTime:: 'rate(node_disk_written_bytes_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      diskIoTime:: 'rate(node_disk_io_time_seconds_total{%(nodeQuerySelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
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
      networkReceiveBitsPerSec:: 'irate(node_network_receive_bytes_total{%(nodeQuerySelector)s}[$__rate_interval])*8' % config { nodeQuerySelector: nodeQuerySelector },
      networkTransmitBitsPerSec:: 'irate(node_network_transmit_bytes_total{%(nodeQuerySelector)s}[$__rate_interval])*8' % config { nodeQuerySelector: nodeQuerySelector },
      networkReceiveErrorsPerSec:: 'irate(node_network_receive_errs_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      networkTransmitErrorsPerSec:: 'irate(node_network_transmit_errs_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      networkReceiveDropsPerSec:: 'irate(node_network_receive_drop_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
      networkTransmitDropsPerSec:: 'irate(node_network_transmit_drop_total{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: nodeQuerySelector },
    },
  },

}
