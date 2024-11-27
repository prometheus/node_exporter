local g = import '../../g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables.main,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',
    uptimeQuery:: 'node_boot_time_seconds',

    diskTotal:
      prometheusQuery.new(
        prometheusDatasource,
        'node_filesystem_size_bytes{%(fsSelector)s, %(fsMountpointSelector)s, %(queriesSelector)s}' % variables { fsMountpointSelector: config.fsMountpointSelector, fsSelector: config.fsSelector }
      ),
    diskTotalRoot:
      prometheusQuery.new(
        prometheusDatasource,
        'node_filesystem_size_bytes{%(queriesSelector)s, mountpoint="/", fstype!="rootfs"}' % variables,
      ),
    diskUsageRoot:
      prometheusQuery.new(
        prometheusDatasource,
        'node_filesystem_avail_bytes{%(queriesSelector)s, mountpoint="/",fstype!="rootfs"}' % variables
      ),
    diskUsageRootPercent:
      prometheusQuery.new(
        prometheusDatasource,
        '100 - node_filesystem_avail_bytes{mountpoint="/",fstype!="rootfs", %(queriesSelector)s}/node_filesystem_size_bytes{mountpoint="/",fstype!="rootfs", %(queriesSelector)s}*100' % variables
      ),
    diskFree:
      prometheusQuery.new(
        prometheusDatasource,
        'node_filesystem_avail_bytes{%(fsSelector)s, %(fsMountpointSelector)s, %(queriesSelector)s}' % variables { fsMountpointSelector: config.fsMountpointSelector, fsSelector: config.fsSelector }
      )
      + prometheusQuery.withLegendFormat('{{ mountpoint }} free'),
    diskUsagePercent:
      prometheusQuery.new(
        prometheusDatasource,
        '100 - node_filesystem_avail_bytes{%(fsSelector)s, %(fsMountpointSelector)s, %(queriesSelector)s}/node_filesystem_size_bytes{%(fsSelector)s, %(fsMountpointSelector)s, %(queriesSelector)s}*100' % variables { fsMountpointSelector: config.fsMountpointSelector, fsSelector: config.fsSelector }
      )
      + prometheusQuery.withLegendFormat('{{ mountpoint }} used, %'),

    diskInodesFree:
      prometheusQuery.new(
        prometheusDatasource,
        'node_filesystem_files_free{%(queriesSelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % variables { fsMountpointSelector: config.fsMountpointSelector, fsSelector: config.fsSelector },
      )
      + prometheusQuery.withLegendFormat('{{ mountpoint }} inodes free'),
    diskInodesTotal:
      prometheusQuery.new(
        prometheusDatasource,
        'node_filesystem_files{%(queriesSelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % variables { fsMountpointSelector: config.fsMountpointSelector, fsSelector: config.fsSelector }
      ) + prometheusQuery.withLegendFormat('{{ mountpoint }} inodes total'),
    diskReadOnly:
      prometheusQuery.new(
        prometheusDatasource,
        'node_filesystem_readonly{%(queriesSelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % variables { fsMountpointSelector: config.fsMountpointSelector, fsSelector: config.fsSelector }
      )
      + prometheusQuery.withLegendFormat('{{ mountpoint }} read-only'),
    diskDeviceError:
      prometheusQuery.new(
        prometheusDatasource,
        'node_filesystem_device_error{%(queriesSelector)s, %(fsSelector)s, %(fsMountpointSelector)s}' % variables { fsMountpointSelector: config.fsMountpointSelector, fsSelector: config.fsSelector }
      )
      + prometheusQuery.withLegendFormat('{{ mountpoint }} device error'),
    // descriptors
    processMaxFds:
      prometheusQuery.new(
        prometheusDatasource,
        'process_max_fds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Maximum open file descriptors'),
    processOpenFds:
      prometheusQuery.new(
        prometheusDatasource,
        'process_open_fds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Open file descriptors'),

    // disk(device)
    diskIOreadBytesPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_disk_read_bytes_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % variables { diskDeviceSelector: config.diskDeviceSelector },
      )
      + prometheusQuery.withLegendFormat('{{ device }} read'),
    diskIOwriteBytesPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_disk_written_bytes_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % variables { diskDeviceSelector: config.diskDeviceSelector },
      )
      + prometheusQuery.withLegendFormat('{{ device }} written'),
    diskIOutilization:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_disk_io_time_seconds_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % variables { diskDeviceSelector: config.diskDeviceSelector },
      )
      + prometheusQuery.withLegendFormat('{{ device }} io util'),
    diskAvgQueueSize:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_disk_io_time_weighted_seconds_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % variables { diskDeviceSelector: config.diskDeviceSelector },
      )
      + prometheusQuery.withLegendFormat('{{ device }} avg queue'),

    diskIOWaitWriteTime:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          irate(node_disk_write_time_seconds_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])
          /
          irate(node_disk_writes_completed_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])
        ||| % variables { diskDeviceSelector: config.diskDeviceSelector }
      )
      + prometheusQuery.withLegendFormat('{{ device }} avg write time'),
    diskIOWaitReadTime:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          irate(node_disk_read_time_seconds_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])
          /
          irate(node_disk_reads_completed_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])
        ||| % variables { diskDeviceSelector: config.diskDeviceSelector }
      )
      + prometheusQuery.withLegendFormat('{{ device }} avg read time'),
    diskIOReads:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          irate(node_disk_reads_completed_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])
        ||| % variables { diskDeviceSelector: config.diskDeviceSelector }
      )
      + prometheusQuery.withLegendFormat('{{ device }} reads'),
    diskIOWrites:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          irate(node_disk_writes_completed_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])
        ||| % variables { diskDeviceSelector: config.diskDeviceSelector }
      )
      + prometheusQuery.withLegendFormat('{{ device }} writes'),

  },
}
