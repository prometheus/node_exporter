local g = import '../../g.libsonnet';
local prometheusQuery = g.query.prometheus;
local lokiQuery = g.query.loki;

{
  new(this): {
    local variables = this.grafana.variables.main,
    local config = this.config,
    local prometheusDatasource = '${' + variables.datasources.prometheus.name + '}',
    local lokiDatasource = '${' + variables.datasources.loki.name + '}',

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
    memoryCachedBytes:
      prometheusQuery.new(
        prometheusDatasource,
        'node_memory_Cached_bytes{%(queriesSelector)s}' % variables
      )
      + prometheusQuery.withLegendFormat('Memory cached'),
    memoryBuffersBytes:
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

  },
}
