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
        refresh='time',
        label='Instance',
      ) {"sort": 1},
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
      datasource="$datasource",
      legendFormat=null,
      format="timeseries",
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

    // return common queries that could be used in multiple dashboards
    queries:: {
      uptime:: 'time() - node_boot_time_seconds{' + config.nodeExporterSelector + ', instance="$instance"}',
      cpuCount:: 'count(count by (cpu)(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance"}))' % config,
      cpuUsage::
        |||
          (((count(count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance"}) by (cpu))) 
          - 
          avg(sum by (mode)(irate(node_cpu_seconds_total{mode='idle',%(nodeExporterSelector)s, instance="$instance"}[5m])))) * 100) 
          / 
          count(count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance"}) by (cpu))
        ||| % config,
      cpuUsagePerCore::
        |||
          (
            (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
          / ignoring(cpu) group_left
            count without (cpu, mode) (node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance"})
          )
        ||| % config,
      memoryTotal:: 'node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config,
      memorySwapTotal:: 'node_memory_SwapTotal_bytes{%(nodeExporterSelector)s, instance="$instance"}' % config,
      fsSizeTotalRoot:: 'node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", mountpoint="/",fstype!="rootfs"}' % config,
      osInfo:: 'node_os_info{%(nodeExporterSelector)s, instance="$instance"}' % config,
      nodeInfo:: 'node_uname_info{%(nodeExporterSelector)s, instance="$instance"}' % config,
      networkReceiveBitsPerSec:: 'irate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance"}[$__rate_interval])*8' % config,
      networkTransmitBitsPerSec:: 'irate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance"}[$__rate_interval])*8' % config,
      networkReceiveErrorsPerSec:: 'irate(node_network_receive_errs_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
      networkTransmitErrorsPerSec:: 'irate(node_network_transmit_errs_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
      networkReceiveDropsPerSec:: 'irate(node_network_receive_drop_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
      networkTransmitDropsPerSec:: 'irate(node_network_transmit_drop_total{%(nodeExporterSelector)s, instance="$instance",}[$__rate_interval])' % config,
    },
  },

}
