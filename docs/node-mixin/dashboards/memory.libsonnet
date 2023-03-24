local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local nodePanels = import '../lib/panels/panels.libsonnet';
local commonPanels = import '../lib/panels/common/panels.libsonnet';
local nodeTimeseries = nodePanels.timeseries;
local common = import '../lib/common.libsonnet';

{

  new(config=null, platform=null):: {
    local c = common.new(config=config, platform=platform),
    local commonPromTarget = c.commonPromTarget,
    local templates = c.templates,
    local q = c.queries,

    local memoryPagesInOut = 
      nodeTimeseries.new('Memory Pages In / Out')
      .withNegativeYByRegex('out')
      .withAxisLabel('out(-) / in(+)')
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pgpgin{%(nodeQuerySelector)s}[$__rate_interval])'  % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pagesin - Page in operations'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pgpgin{%(nodeQuerySelector)s}[$__rate_interval])'  % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pagesout - Page out operations'
      )),
    local memoryPagesSwapInOut = 
      nodeTimeseries.new('Memory Pages Swap In / Out')
      .withNegativeYByRegex('out')
      .withAxisLabel('out(-) / in(+)')
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pswpin{%(nodeQuerySelector)s}[$__rate_interval])'  % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pswpin - Pages swapped in'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pswpout{%(nodeQuerySelector)s}[$__rate_interval])'  % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pswpout - Pages swapped out'
      )),

    local memoryPagesFaults = 
      nodeTimeseries.new('Memory Page Faults')
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pgfault{%(nodeQuerySelector)s}[$__rate_interval])'  % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pgfault - Page major and minor fault operations'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pgmajfault{%(nodeQuerySelector)s}[$__rate_interval])'  % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pgmajfault - Major page fault operations'
      ))
      .addTarget(commonPromTarget(
        expr=
          |||
            irate(node_vmstat_pgfault{%(nodeQuerySelector)s}[$__rate_interval])
            -
            irate(node_vmstat_pgmajfault{%(nodeQuerySelector)s}[$__rate_interval])
          ||| % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pgminfault - Minor page fault operations'
      )),

      local memoryOOMkiller =
      nodeTimeseries.new('OOM Killer')
        .addTarget(commonPromTarget(
          expr="increase(node_vmstat_oom_kill{%(nodeQuerySelector)s}[$__rate_interval])" % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat="OOM killer invocations"
        )),
      
      local memoryActiveInactive = 
        nodeTimeseries.new('Memory Active / Inactive')
        .withUnits("decbytes")
        .addTarget(commonPromTarget(
          expr='node_memory_Inactive_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat='Inactive - Memory which has been less recently used.  It is more eligible to be reclaimed for other purposes',
        ))
        .addTarget(commonPromTarget(
          expr='node_memory_Active_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat='Active - Memory that has been used more recently and usually not reclaimed unless absolutely necessary',
        )),

      local memoryActiveInactiveDetail = 
        nodeTimeseries.new('Memory Active / Inactive Details')
        .withUnits("decbytes")
        .addTarget(commonPromTarget(
          expr='node_memory_Inactive_file_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat='Inactive_file - File-backed memory on inactive LRU list',
        ))
        .addTarget(commonPromTarget(
          expr='node_memory_Inactive_anon_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat='Inactive_anon - Anonymous and swap cache on inactive LRU list, including tmpfs (shmem)',
        ))
        .addTarget(commonPromTarget(
          expr='node_memory_Active_file_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat='Active_file - File-backed memory on active LRU list',
        ))
        .addTarget(commonPromTarget(
          expr='node_memory_Active_anon_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
          legendFormat='Active_anon - Anonymous and swap cache on active least-recently-used (LRU) list, including tmpfs',
        )),
      
    local memoryCommited =
      nodeTimeseries.new('Memory Commited')
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Committed_AS_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Committed_AS - Amount of memory presently allocated on the system'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_CommitLimit_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='CommitLimit - Amount of memory currently available to be allocated on the system'
      )),
    local memorySharedAndMapped =
      nodeTimeseries.new('Memory Shared and Mapped')
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Mapped_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Mapped - Used memory in mapped pages files which have been mmaped, such as libraries'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Shmem_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Shmem - Used shared memory (shared between several processes, thus including RAM disks)'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_ShmemHugePages_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ShmemHugePages - Memory used by shared memory (shmem) and tmpfs allocated  with huge pages'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_ShmemPmdMapped_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ShmemPmdMapped - Amount of shared (shmem/tmpfs) memory backed by huge pages'
      )),

    local memoryWriteAndDirty =
      nodeTimeseries.new('Memory Writeback and Dirty')
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Writeback_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Writeback - Memory which is actively being written back to disk'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_WritebackTmp_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='WritebackTmp - Memory used by FUSE for temporary writeback buffers'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Dirty_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Dirty - Memory which is waiting to get written back to the disk'
      )),

    local memoryVmalloc =
      nodeTimeseries.new('Memory Vmalloc')
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_VmallocChunk_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='VmallocChunk - Largest contiguous block of vmalloc area which is free'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_VmallocTotal_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='VmallocTotal - Total size of vmalloc memory area'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_VmallocUsed_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='VmallocUsed - Amount of vmalloc area which is used'
      )),

    local memorySlab =
      nodeTimeseries.new('Memory Slab')
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_SUnreclaim_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='SUnreclaim - Part of Slab, that cannot be reclaimed on memory pressure'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_SReclaimable_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='SReclaimable - Part of Slab, that might be reclaimed, such as caches'
      )),

    local memoryAnonymous =
      nodeTimeseries.new('Memory Anonymous')
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_AnonHugePages_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='AnonHugePages - Memory in anonymous huge pages'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_AnonPages_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='AnonPages - Memory in user pages not backed by files'
      )),

    local memoryHugePagesCounter =
      nodeTimeseries.new('Memory HugePages Counter')
      .addTarget(commonPromTarget(
        expr='node_memory_HugePages_Free{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='HugePages_Free - Huge pages in the pool that are not yet allocated'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_HugePages_Rsvd{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='HugePages_Rsvd - Huge pages for which a commitment to allocate from the pool has been made, but no allocation has yet been made'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_HugePages_Surp{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='HugePages_Surp - Huge pages in the pool above the value in /proc/sys/vm/nr_hugepages'
      )),
    local memoryHugePagesSize =
      nodeTimeseries.new('Memory HugePages Size')
      .withUnits("decbytes")
      .addTarget(commonPromTarget(
        expr='node_memory_HugePages_Total{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='HugePages - Total size of the pool of huge pages'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Hugepagesize_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Hugepagesize - Huge Page size'
      )),
    local memoryDirectMap =
      nodeTimeseries.new('Memory Direct Map')
      .withUnits("decbytes")
      .addTarget(commonPromTarget(
        expr='node_memory_DirectMap1G_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='DirectMap1G'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_DirectMap2M_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='DirectMap2M'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_DirectMap4k_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='DirectMap4k'
      )),
      
    local memoryBounce =
      nodeTimeseries.new('Memory Bounce')
      .withUnits("decbytes")
      .addTarget(commonPromTarget(
        expr='node_memory_Bounce_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Bounce - Memory used for block device bounce buffers'
      )),
    local panelsGrid =
      [
        //TODO add generic graphs
        
        {type:"row", title: "Vmstat" , gridPos: {y: 25}},
        memoryPagesInOut {gridPos: {x:0, w:12, h: 8, y: 25}},
        memoryPagesSwapInOut {gridPos: {x:12, w:12, h: 8, y: 25}},
        memoryPagesFaults {gridPos: {x:0, w:12, h: 8, y: 25}},
        memoryOOMkiller {gridPos: {x:12, w:12, h: 8, y: 25}},
        {type:"row", title: "Memstat", gridPos: {y: 50}},
        memoryActiveInactive {gridPos: {x:0, w:12, h: 8, y: 50}},
        memoryActiveInactiveDetail {gridPos: {x:12, w:12, h: 8, y: 50}},
        memoryCommited {gridPos: {x:0, w:12, h: 8, y: 50}},
        memorySharedAndMapped {gridPos: {x:12, w:12, h: 8, y: 50}},
        memoryWriteAndDirty {gridPos: {x:0, w:12, h: 8, y: 50}},
        memoryVmalloc {gridPos: {x:12, w:12, h: 8, y: 50}},
        memorySlab {gridPos: {x:0, w:12, h: 8, y: 50}},
        memoryAnonymous {gridPos: {x:12, w:12, h: 8, y: 50}},
        memoryHugePagesCounter {gridPos: {x:0, w:12, h: 8, y: 50}},
        memoryHugePagesSize {gridPos: {x:12, w:12, h: 8, y: 50}},
        memoryDirectMap {gridPos: {x:0, w:12, h: 8, y: 50}},
        memoryBounce {gridPos: {x:12, w:12, h: 8, y: 50}},
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Memory' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='node-memory'
      ) { editable: true }
      .addLink(c.links.fleetDash)
      .addLink(c.links.nodeDash)
      .addLink(c.links.otherDashes)
      .addAnnotations(c.annotations)
      .addTemplates(templates)
      .addPanels(panelsGrid)
    else if platform == 'Darwin' then {},
  },
}
