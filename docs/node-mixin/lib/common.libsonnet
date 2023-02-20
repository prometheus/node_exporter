local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;

{

  new(config=null, platform=null):: {

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

    local instanceTemplatePrototype =
      template.new(
        'instance',
        '$datasource',
        '',
        sort=1,
        refresh='time',
        label='Instance',
      ),
    local instanceTemplate =
      if platform == 'Darwin' then
        instanceTemplatePrototype
        { query: 'label_values(node_uname_info{%(nodeExporterSelector)s, sysname="Darwin"}, instance)' % config }
      else
        instanceTemplatePrototype
        { query: 'label_values(node_uname_info{%(nodeExporterSelector)s, sysname!="Darwin"}, instance)' % config },

    // return common templates
    templates: [
      prometheusDatasourceTemplate,
      instanceTemplate,
    ],

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
    },
    // return common queries that could be used in multiple dashboards
    queries:: {
      systemLoad1:: 'avg by (instance) (node_load1{%(nodeExporterSelector)s, instance=~"$instance"})' % config,
      systemLoad5:: 'avg by (instance) (node_load5{%(nodeExporterSelector)s, instance=~"$instance"})' % config,
      systemLoad15:: 'avg by (instance) (node_load15{%(nodeExporterSelector)s, instance=~"$instance"})' % config,
      uptime:: 'time() - node_boot_time_seconds{' + config.nodeExporterSelector + ', instance=~"$instance"}',
      cpuCount:: 'count by (instance) (node_cpu_seconds_total{%(nodeExporterSelector)s, instance=~"$instance", mode="idle"})' % config,
      cpuUsage::
        |||
          (((count by (instance) (count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance=~"$instance"}) by (cpu, instance))) 
          - 
          avg by (instance) (sum by (instance, mode)(irate(node_cpu_seconds_total{mode='idle',%(nodeExporterSelector)s, instance=~"$instance"}[5m])))) * 100) 
          / 
          count by(instance) (count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance=~"$instance"}) by (cpu, instance))
        ||| % config,
      cpuUsagePerCore::
        |||
          (
            (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode=~"idle|iowait|steal", instance=~"$instance"}[$__rate_interval])))
          / ignoring(cpu) group_left
            count without (cpu, mode) (node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance=~"$instance"})
          )
        ||| % config,
      memoryTotal:: 'node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance=~"$instance"}' % config,
      memorySwapTotal:: 'node_memory_SwapTotal_bytes{%(nodeExporterSelector)s, instance=~"$instance"}' % config,
      memoryUsage::
        |||
          100 -
          (
            avg by (instance) (node_memory_MemAvailable_bytes{%(nodeExporterSelector)s, instance=~"$instance"}) /
            avg by (instance) (node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance=~"$instance"})
          * 100
          )
        ||| % config,
      fsSizeTotalRoot:: 'node_filesystem_size_bytes{%(nodeExporterSelector)s, instance=~"$instance", mountpoint="/",fstype!="rootfs"}' % config,
      osInfo:: 'node_os_info{%(nodeExporterSelector)s, instance=~"$instance"}' % config,
      nodeInfo:: 'node_uname_info{%(nodeExporterSelector)s, instance=~"$instance"}' % config,
      diskReadTime:: 'rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance=~"$instance", %(diskDeviceSelector)s}[$__rate_interval])' % config,
      diskWriteTime:: 'rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance=~"$instance", %(diskDeviceSelector)s}[$__rate_interval])' % config,
      diskIoTime:: 'rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, instance=~"$instance", %(diskDeviceSelector)s}[$__rate_interval])' % config,
      diskSpaceUsage::
        |||
          sort_desc(1 -
            (
            max by (job, instance, fstype, device) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance=~"$instance", %(fsSelector)s, %(fsMountpointSelector)s})
            /
            max by (job, instance, fstype, device) (node_filesystem_size_bytes{%(nodeExporterSelector)s, instance=~"$instance", %(fsSelector)s, %(fsMountpointSelector)s})
            ) != 0
          )
        ||| % config,
      networkReceiveBitsPerSec:: 'irate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance=~"$instance"}[$__rate_interval])*8' % config,
      networkTransmitBitsPerSec:: 'irate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance=~"$instance"}[$__rate_interval])*8' % config,
      networkReceiveErrorsPerSec:: 'irate(node_network_receive_errs_total{%(nodeExporterSelector)s, instance=~"$instance",}[$__rate_interval])' % config,
      networkTransmitErrorsPerSec:: 'irate(node_network_transmit_errs_total{%(nodeExporterSelector)s, instance=~"$instance",}[$__rate_interval])' % config,
      networkReceiveDropsPerSec:: 'irate(node_network_receive_drop_total{%(nodeExporterSelector)s, instance=~"$instance",}[$__rate_interval])' % config,
      networkTransmitDropsPerSec:: 'irate(node_network_transmit_drop_total{%(nodeExporterSelector)s, instance=~"$instance",}[$__rate_interval])' % config,
    },
  },

}
