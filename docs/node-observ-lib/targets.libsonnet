local g = import './g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',
    uptimeQuery:: 'node_boot_time_seconds',

    reboot:
      prometheusQuery.new(
        prometheusDatasource,
        self.uptimeQuery + '{%(queriesSelector)s}*1000 > $__from < $__to' % variables,
      ),

    serviceFailed:
      lokiQuery.new(
        lokiDatasource,
        '{%(queriesSelector)s, unit="init.scope"} |= "code=exited, status=1/FAILURE"' % variables
      ),
    // those events should be rare, so can be shown as annotations
    criticalEvents:
      lokiQuery.new(
        lokiDatasource,
        '{%(queriesSelector)s, transport="kernel", level="emerg"}' % variables
      ),
    memoryOOMkiller:
      prometheusQuery.new(
        prometheusDatasource,
        'increase(node_vmstat_oom_kill{%(queriesSelector)s}[$__interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('OOM killer invocations'),

    kernelUpdate:
      prometheusQuery.new(
        prometheusDatasource,
        expr=|||
          changes(
          sum by (%(instanceLabels)s) (
              group by (%(instanceLabels)s,release) (node_uname_info{%(queriesSelector)s})
              )
          [$__interval:1m] offset -$__interval) > 1
        ||| % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      ),

    // new interactive session in logs:
    sessionOpened:
      lokiQuery.new(
        lokiDatasource,
        '{%(queriesSelector)s, unit="systemd-logind.service"}|= "New session"' % variables
      ),
    sessionClosed:
      lokiQuery.new(
        lokiDatasource,
        '{%(queriesSelector)s, unit="systemd-logind.service"} |= "logged out"' % variables
      ),

    alertsCritical:
      prometheusQuery.new(
        prometheusDatasource,
        'count by (%(instanceLabels)s) (max_over_time(ALERTS{%(queriesSelector)s, alertstate="firing", severity="critical"}[1m])) * group by (%(instanceLabels)s) (node_uname_info{%(queriesSelector)s})' % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      ),
    alertsWarning:
      prometheusQuery.new(
        prometheusDatasource,
        'count by (%(instanceLabels)s) (max_over_time(ALERTS{%(queriesSelector)s, alertstate="firing", severity="warning"}[1m])) * group by (%(instanceLabels)s) (node_uname_info{%(queriesSelector)s})' % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      ),

    uptime:
      prometheusQuery.new(
        prometheusDatasource,
        'time() - ' + self.uptimeQuery + '{%(queriesSelector)s}' % variables
      ),
    cpuCount:
      prometheusQuery.new(
        prometheusDatasource,
        'count without (cpu) (node_cpu_seconds_total{%(queriesSelector)s, mode="idle"})' % variables
      )
      + prometheusQuery.withLegendFormat('Cores'),
    cpuUsage:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (((count by (%(instanceLabels)s) (count(node_cpu_seconds_total{%(queriesSelector)s}) by (cpu, %(instanceLabels)s))) 
          - 
          avg by (%(instanceLabels)s) (sum by (%(instanceLabels)s, mode)(irate(node_cpu_seconds_total{mode='idle',%(queriesSelector)s}[$__rate_interval])))) * 100) 
          / 
          count by(%(instanceLabels)s) (count(node_cpu_seconds_total{%(queriesSelector)s}) by (cpu, %(instanceLabels)s))
        ||| % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      )
      + prometheusQuery.withLegendFormat('CPU usage'),
    cpuUsagePerCore:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
            (1 - sum without (mode) (rate(node_cpu_seconds_total{%(queriesSelector)s, mode=~"idle|iowait|steal"}[$__rate_interval])))
          / ignoring(cpu) group_left
            count without (cpu, mode) (node_cpu_seconds_total{%(queriesSelector)s, mode="idle"})
          ) * 100
        ||| % variables,
      )
      + prometheusQuery.withLegendFormat('CPU {{cpu}}'),
    cpuUsageByMode:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          sum by(%(instanceLabels)s, mode) (irate(node_cpu_seconds_total{%(queriesSelector)s}[$__rate_interval])) 
          / on(%(instanceLabels)s) 
          group_left sum by (%(instanceLabels)s)((irate(node_cpu_seconds_total{%(queriesSelector)s}[$__rate_interval]))) * 100
        ||| % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      )
      + prometheusQuery.withLegendFormat('{{ mode }}'),
    memoryTotalBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_MemTotal_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory total'),
    memoryFreeBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_MemFree_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory free'),
    memoryAvailableBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_MemAvailable_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory available'),
    memoryBuffersBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Cached_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory buffers'),
    memoryCachedBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Buffers_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory buffers'),
    memoryUsedBytes:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          (
            node_memory_MemTotal_bytes{%(queriesSelector)s}
          -
            node_memory_MemFree_bytes{%(queriesSelector)s}
          -
            node_memory_Buffers_bytes{%(queriesSelector)s}
          -
            node_memory_Cached_bytes{%(queriesSelector)s}
          )
        ||| % variables
      )
      + prometheusQuery.withLegendFormat('Memory used'),
    memoryUsagePercent:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          100 -
          (
            avg by (%(instanceLabels)s) (node_memory_MemAvailable_bytes{%(queriesSelector)s}) /
            avg by (%(instanceLabels)s) (node_memory_MemTotal_bytes{%(queriesSelector)s})
          * 100
          )
        |||
        % variables { instanceLabels: std.join(',', this.config.instanceLabels) },
      ),
    memorySwapTotal:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_SwapTotal_bytes{%(queriesSelector)s}' % variables
      ),
    memoryPagesIn:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_vmstat_pgpgin{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Page-In'),
    memoryPagesOut:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_vmstat_pgpgout{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Page-Out'),

    memoryPagesSwapIn:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_vmstat_pswpin{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Pages swapped in'),
    memoryPagesSwapOut:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_vmstat_pswpout{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Pages swapped out'),

    memoryPageMajorFaults:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_vmstat_pgmajfault{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Major page fault operations'),
    memoryPageMinorFaults:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          irate(node_vmstat_pgfault{%(queriesSelector)s}[$__rate_interval])
          -
          irate(node_vmstat_pgmajfault{%(queriesSelector)s}[$__rate_interval])
        ||| % variables,
      )
      + prometheusQuery.withLegendFormat('Minor page fault operations'),

    memoryInactiveBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Inactive_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Inactive'),
    memoryActiveBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Active_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Active'),

    memoryInactiveFile:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Inactive_file_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Inactive_file'),

    memoryInactiveAnon:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Inactive_anon_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Inactive_anon'),

    memoryActiveFile:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Active_file_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Active_file'),

    memoryActiveAnon:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Active_anon_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Active_anon'),

    memoryCommitedAs:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Committed_AS_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Commited_AS'),
    memoryCommitedLimit:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_CommitLimit_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('CommitLimit'),

    memoryMappedBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Mapped_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Mapped'),
    memoryShmemBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Shmem_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Shmem'),
    memoryShmemHugePagesBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_ShmemHugePages_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('ShmemHugePages'),
    memoryShmemPmdMappedBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_ShmemPmdMapped_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('ShmemPmdMapped'),
    memoryWriteback:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Writeback_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Writeback'),
    memoryWritebackTmp:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_WritebackTmp_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('WritebackTmp'),
    memoryDirty:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Dirty_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Dirty'),

    memoryVmallocChunk:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_VmallocChunk_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('VmallocChunk'),
    memoryVmallocTotal:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_VmallocTotal_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('VmallocTotal'),
    memoryVmallocUsed:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_VmallocUsed_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('VmallocUsed'),
    memorySlabSUnreclaim:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_SUnreclaim_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('SUnreclaim'),
    memorySlabSReclaimable:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_SReclaimable_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('SReclaimable'),

    memoryAnonHugePages:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_AnonHugePages_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('AnonHugePages'),
    memoryAnonPages:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_AnonPages_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('AnonPages'),

    memoryHugePages_Free:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_HugePages_Free{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('HugePages_Free'),
    memoryHugePages_Rsvd:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_HugePages_Rsvd{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('HugePages_Rsvd'),
    memoryHugePages_Surp:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_HugePages_Surp{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('HugePages_Surp'),
    memoryHugePagesTotalSize:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_HugePages_Total{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Huge pages total size'),
    memoryHugePagesSize:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Hugepagesize_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Huge page size'),
    memoryDirectMap1G:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_DirectMap1G_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('DirectMap1G'),
    memoryDirectMap2M:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_DirectMap2M_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('DirectMap2M'),
    memoryDirectMap4k:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_DirectMap4k_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('DirectMap4k'),
    memoryBounce:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Bounce_bytes{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Bounce'),

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
        'irate(node_disk_reads_completed_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % variables { diskDeviceSelector: config.diskDeviceSelector },
      )
      + prometheusQuery.withLegendFormat('{{ device }} read'),
    diskIOwriteBytesPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_disk_writes_completed_total{%(queriesSelector)s, %(diskDeviceSelector)s}[$__rate_interval])' % variables { diskDeviceSelector: config.diskDeviceSelector },
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

    unameInfo:
      prometheusQuery.new(
        prometheusDatasource,
        'node_uname_info{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withFormat('table'),
    osInfo:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          node_os_info{%(queriesSelector)s}
        ||| % variables,
      )
      + prometheusQuery.withFormat('table'),
    osInfoCombined:
      prometheusQuery.new(
        prometheusDatasource,
        |||
          node_uname_info{%(queriesSelector)s} 
          * on (%(groupLabels)s,%(instanceLabels)s)
          group_left(pretty_name)
          node_os_info{%(queriesSelector)s}
        ||| % variables {
          instanceLabels: std.join(',', this.config.instanceLabels),
          groupLabels: std.join(',', this.config.groupLabels),
        },
      )
      + prometheusQuery.withFormat('table'),

    osTimezone:  //timezone label
      prometheusQuery.new(
        prometheusDatasource,
        'node_time_zone_offset_seconds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withFormat('table'),

    systemLoad1:
      prometheusQuery.new(
        prometheusDatasource,
        'node_load1{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('1m'),
    systemLoad5:
      prometheusQuery.new(
        prometheusDatasource,
        'node_load5{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('5m'),
    systemLoad15:
      prometheusQuery.new(
        prometheusDatasource,
        'node_load15{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('15m'),

    systemContextSwitches:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_context_switches_total{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Context switches'),

    systemInterrupts:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_intr_total{%(queriesSelector)s}[$__rate_interval])' % variables,
      )
      + prometheusQuery.withLegendFormat('Interrupts'),

    timeNtpStatus:
      prometheusQuery.new(
        prometheusDatasource,
        'node_timex_sync_status{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('NTP status'),

    timeOffset:
      prometheusQuery.new(
        prometheusDatasource,
        'node_timex_offset_seconds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Time offset'),

    timeEstimatedError:
      prometheusQuery.new(
        prometheusDatasource,
        'node_timex_estimated_error_seconds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Estimated error in seconds'),
    timeMaxError:
      prometheusQuery.new(
        prometheusDatasource,
        'node_timex_maxerror_seconds{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('Maximum error in seconds'),

    networkUp:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_up{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('{{device}}'),
    networkCarrier:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_carrier{%(queriesSelector)s}' % variables,
      )
      + prometheusQuery.withLegendFormat('{{device}}'),
    networkArpEntries:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_arp{%(queriesSelector)s}' % variables,
      ),
    networkMtuBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_mtu_bytes{%(queriesSelector)s}' % variables,
      ),
    networkSpeedBitsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_speed_bytes{%(queriesSelector)s} * 8' % variables,
      ),
    networkTransmitQueueLength:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_transmit_queue_length{%(queriesSelector)s}' % variables,
      ),
    networkInfo:
      prometheusQuery.new(
        prometheusDatasource,
        'node_network_info{%(queriesSelector)s}' % variables,
      ),

    networkOutBitPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_bytes_total{%(queriesSelector)s}[$__rate_interval])*8' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),
    networkInBitPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_bytes_total{%(queriesSelector)s}[$__rate_interval])*8' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkOutErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_errs_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} errors transmitted'),
    networkInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_errs_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} errors received'),
    networkOutDroppedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_drop_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted dropped'),
    networkInDroppedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_drop_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received dropped'),

    networkInPacketsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_packets_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkOutPacketsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_packets_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),

    networkInMulticastPacketsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_multicast_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkOutMulticastPacketsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_multicast_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),
    networkFifoInPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_fifo_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkFifoOutPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_fifo_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),

    networkCompressedInPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_receive_compressed_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} received'),
    networkCompressedOutPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_network_transmit_compressed_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('{{ device }} transmitted'),

    networkNFConntrackEntries:
      prometheusQuery.new(
        prometheusDatasource,
        'node_nf_conntrack_entries{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('NF conntrack entries'),
    networkNFConntrackLimits:
      prometheusQuery.new(
        prometheusDatasource,
        'node_nf_conntrack_entries_limit{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('NF conntrack limits'),

    networkSoftnetProcessedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_softnet_processed_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('CPU {{ cpu }} processed'),
    networkSoftnetDroppedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_softnet_dropped_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('CPU {{ cpu }} dropped'),
    networkSoftnetSqueezedPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_softnet_times_squeezed_total{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('CPU {{ cpu }} out of quota'),

    networkSocketsUsed:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_sockets_used{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 sockets in use'),
    networkSocketsTCPAllocated:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_alloc{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Allocated'),
    networkSocketsTCPIPv6:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 in use'),
    networkSocketsTCPIPv4:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 in use'),
    networkSocketsTCPOrphans:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_orphan{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Orphan sockets'),
    networkSocketsTCPTimeWait:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_tw{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Time wait'),

    networkSocketsUDPLiteInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDPLITE_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 UDPLITE in use'),
    networkSocketsUDPInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDP_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 UDP in use'),
    networkSocketsUDPLiteIPv6InUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDPLITE6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 UDPLITE in use'),
    networkSocketsUDPIPv6InUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDP6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 UDP in use'),

    networkSocketsFragInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_FRAG_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 Frag sockets in use'),
    networkSocketsFragIPv6InUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_FRAG6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 Frag sockets in use'),
    networkSocketsRawInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_RAW_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv4 Raw sockets in use'),
    networkSocketsIPv6RawInUse:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_RAW6_inuse{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('IPv6 Raw sockets in use'),

    networkSocketsTCPMemoryPages:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_mem{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory pages allocated for TCP sockets'),
    networkSocketsUDPMemoryPages:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDP_mem{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory pages allocated for UDP sockets'),

    networkSocketsTCPMemoryBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_TCP_mem_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory bytes allocated for TCP sockets'),
    networkSocketsUDPMemoryBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_sockstat_UDP_mem_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory bytes allocated for UDP sockets'),

    networkNetstatIPInOctetsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_IpExt_InOctets{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('Octets received'),
    networkNetstatIPOutOctetsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_IpExt_OutOctets{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('Octets transmitted'),

    networkNetstatTCPInSegmentsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_InSegs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP received'),
    networkNetstatTCPOutSegmentsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_OutSegs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP transmitted'),

    networkNetstatTCPOverflowPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_TcpExt_ListenOverflows{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP overflow'),

    networkNetstatTCPListenDropsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_TcpExt_ListenDrops{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP ListenDrops - SYNs to LISTEN sockets ignored'),

    networkNetstatTCPRetransPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_TcpExt_TCPSynRetrans{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP SYN rentransmits'),

    networkNetstatTCPRetransSegPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_RetransSegs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP retransmitted segments, containing one or more previously transmitted octets'),
    networkNetstatTCPInWithErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_InErrs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP received with errors'),

    networkNetstatTCPOutWithRstPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Tcp_OutRsts{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('TCP segments sent with RST flag'),

    networkNetstatIPInUDPPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_InDatagrams{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP received'),

    networkNetstatIPOutUDPPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_OutDatagrams{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP transmitted'),

    //UDP errors
    networkNetstatUDPLiteInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_UdpLite_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDPLite InErrors'),

    networkNetstatUDPInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP InErrors'),
    networkNetstatUDP6InErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 InErrors'),
    networkNetstatUDPNoPortsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_NoPorts{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP NoPorts'),
    networkNetstatUDP6NoPortsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_NoPorts{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 NoPorts'),
    networkNetstatUDPRcvBufErrsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_RcvbufErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP receive buffer errors'),
    networkNetstatUDP6RcvBufErrsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_RcvbufErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 receive buffer errors'),
    networkNetstatUDPSndBufErrsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp_SndbufErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP transmit buffer errors'),
    networkNetstatUDP6SndBufErrsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Udp6_SndbufErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('UDP6 transmit buffer errors'),

    //ICMP
    networkNetstatICMPInPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp_InMsgs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP received'),
    networkNetstatICMPOutPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp_OutMsgs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP transmitted'),
    networkNetstatICMP6InPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp6_InMsgs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP6 received'),
    networkNetstatICMP6OutPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp6_OutMsgs{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP6 transmitted'),

    networkNetstatICMPInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP6 errors'),
    networkNetstatICM6PInErrorsPerSec:
      prometheusQuery.new(
        prometheusDatasource,
        'irate(node_netstat_Icmp6_InErrors{%(queriesSelector)s}[$__rate_interval])' % variables
      )
      + prometheusQuery.withLegendFormat('ICMP6 errors'),
  },
}
