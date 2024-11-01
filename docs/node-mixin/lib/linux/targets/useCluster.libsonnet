local g = import '../../g.libsonnet';
local commonlib = import 'github.com/grafana/jsonnet-libs/common-lib/common/main.libsonnet';
local xtd = import 'github.com/jsonnet-libs/xtd/main.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;


{
  new(this): {
    local variables = this.grafana.variables.useCluster,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',

    cpuUtilization:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          ((
            instance:node_cpu_utilisation:rate5m{%(queriesSelector)s}
            *
            instance:node_num_cpu:sum{%(queriesSelector)s}
          ) != 0 )
          / scalar(sum(instance:node_num_cpu:sum{%(queriesSelector)s})) * 100
        ||| % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: Utilization' % xtd.array.slice(this.config.instanceLabels, -1)),

    cpuSaturation:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
            instance:node_load1_per_cpu:ratio{%(queriesSelector)s}
            / scalar(count(instance:node_load1_per_cpu:ratio{%(queriesSelector)s}))
          ) * 100 != 0
        ||| % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: Saturation' % xtd.array.slice(this.config.instanceLabels, -1)),
    memoryUtilization:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
            instance:node_memory_utilisation:ratio{%(queriesSelector)s}
            / scalar(count(instance:node_memory_utilisation:ratio{%(queriesSelector)s}))
          ) * 100 != 0 
        |||
        % variables
      )
      + prometheusQuery.withLegendFormat('{{%s}}: Utilization' % xtd.array.slice(this.config.instanceLabels, -1)),

    memorySaturation:
      prometheusQuery.new(
        prometheusDatasource,
        'instance:node_vmstat_pgmajfault:rate5m{%(queriesSelector)s} != 0' % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: Major page fault operations' % xtd.array.slice(this.config.instanceLabels, -1)),
    networkUtilizationReceive:
      prometheusQuery.new(
        prometheusDatasource,
        'instance:node_network_receive_bytes_excluding_lo:rate5m{%(queriesSelector)s} != 0' % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: Receive' % xtd.array.slice(this.config.instanceLabels, -1)),
    networkUtilizationTransmit:
      prometheusQuery.new(
        prometheusDatasource,
        'instance:node_network_transmit_bytes_excluding_lo:rate5m{%(queriesSelector)s} != 0' % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: Transmit' % xtd.array.slice(this.config.instanceLabels, -1)),
    networkSaturationReceive:
      prometheusQuery.new(
        prometheusDatasource,
        'instance:node_network_receive_drop_excluding_lo:rate5m{%(queriesSelector)s} != 0' % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: Receive' % xtd.array.slice(this.config.instanceLabels, -1)),
    networkSaturationTransmit:
      prometheusQuery.new(
        prometheusDatasource,
        'instance:node_network_transmit_drop_excluding_lo:rate5m{%(queriesSelector)s} != 0' % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: Transmit' % xtd.array.slice(this.config.instanceLabels, -1)),

    diskUtilization:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
            instance_device:node_disk_io_time_seconds:rate5m{%(queriesSelector)s}
            / scalar(count(instance_device:node_disk_io_time_seconds:rate5m{%(queriesSelector)s}))
          ) * 100 != 0
        ||| % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: {{device}}' % xtd.array.slice(this.config.instanceLabels, -1)),
    diskSaturation:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
            instance_device:node_disk_io_time_weighted_seconds:rate5m{%(queriesSelector)s}
            / scalar(count(instance_device:node_disk_io_time_weighted_seconds:rate5m{%(queriesSelector)s}))
          ) * 100 != 0
        ||| % variables,
      )
      + prometheusQuery.withLegendFormat('{{%s}}: {{device}}' % xtd.array.slice(this.config.instanceLabels, -1)),

    filesystemUtilization:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          sum without (device) (
            max without (fstype, mountpoint) ((
              node_filesystem_size_bytes{%(queriesSelector)s, %(fsMountpointSelector)s, %(fsSelector)s}
              -
              node_filesystem_avail_bytes{%(queriesSelector)s, %(fsMountpointSelector)s, %(fsSelector)s}
            ) != 0)
          )
          / scalar(sum(max without (fstype, mountpoint) (node_filesystem_size_bytes{%(queriesSelector)s, %(fsMountpointSelector)s, %(fsSelector)s}))) * 100
        ||| % variables { fsMountpointSelector: this.config.fsMountpointSelector, fsSelector: this.config.fsSelector },
      )
      + prometheusQuery.withLegendFormat('{{%s}}: {{device}}' % xtd.array.slice(this.config.instanceLabels, -1)),


  },
}
