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
      nodeTimeseries.new(
        'Memory Pages In / Out',
        description=|||
          Page-In - Return of pages to physical memory. This is a common and normal event.

          Page-Out - process of writing pages to disk. Unlike page-in, page-outs can indicate trouble.
          When the kernel detects low memory, it attempts to free memory by paging out. 
          While occasional page-outs are normal, excessive and frequent page-outs can lead to thrashing.
          Thrashing is a state in which the kernel spends more time managing paging activity than running applications, resulting in poor system performance.
        |||
      )
      .withNegativeYByRegex('out')
      .withAxisLabel('out(-) | in(+)')
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pgpgin{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Page-In'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pgpgout{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Page-Out'
      )),
    local memoryPagesSwapInOut =
      nodeTimeseries.new(
        'Memory Pages Swapping In / Out',
        description=|||
          Compared to the speed of the CPU and main memory, writing pages out to disk is relatively slow.
          Nonetheless, it is a preferable option to crashing or killing off processes.

          The process of writing pages out to disk to free memory is known as swapping-out.
          If a page fault occurs because the page is on disk, in the swap area, rather than in memory,
          the kernel will read the page back in from the disk to satisfy the page fault. 
          This is known as swapping-in.
        |||
      )
      .withNegativeYByRegex('out')
      .withAxisLabel('out(-) | in(+)')
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pswpin{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pages swapped in'
      ))
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pswpout{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Pages swapped out'
      )),

    local memoryPagesFaults =
      nodeTimeseries.new(
        'Memory Page Faults',
        description=|||
          A page fault is an exception raised by the memory when a process accesses a memory page without the necessary preparations,
          requiring a mapping to be added to the process's virtual address space. The page contents may also need to be loaded from a backing store such as a disk.
          While the MMU detects the page fault, the operating system's kernel handles the exception by either making the required page accessible in physical memory or denying an illegal memory access.
          Valid page faults are common and necessary to increase memory availability in any operating system that uses virtual memory, including Windows, macOS, and the Linux kernel.
        |||
      )
      .addTarget(commonPromTarget(
        expr='irate(node_vmstat_pgmajfault{%(nodeQuerySelector)s}[$__rate_interval])' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Major page fault operations'
      ))
      .addTarget(commonPromTarget(
        expr=
        |||
          irate(node_vmstat_pgfault{%(nodeQuerySelector)s}[$__rate_interval])
          -
          irate(node_vmstat_pgmajfault{%(nodeQuerySelector)s}[$__rate_interval])
        ||| % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Minor page fault operations'
      )),

    local memoryOOMkiller =
      nodeTimeseries.new(
        'OOM Killer',
        description=|||
          Out Of Memory Killer is a process used by the Linux kernel when the system is running critically low on memory.
          This can happen when the kernel has allocated more memory than is available for its processes.
        |||
      )
      .addTarget(commonPromTarget(
        expr='increase(node_vmstat_oom_kill{%(nodeQuerySelector)s}[$__interval] offset -$__interval)' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='OOM killer invocations'
      )),

    local memoryActiveInactive =
      nodeTimeseries.new(
        'Memory Active / Inactive',
        description=|||
          Inactive: Memory which has been less recently used.  It is more eligible to be reclaimed for other purposes.
          Active: Memory that has been used more recently and usually not reclaimed unless absolutely necessary.
        |||)
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Inactive_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Inactive',
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Active_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Active',
      )),

    local memoryActiveInactiveDetail =
      nodeTimeseries.new(
        'Memory Active / Inactive Details',
        description=|||
          Inactive_file: File-backed memory on inactive LRU list.
          Inactive_anon: Anonymous and swap cache on inactive LRU list, including tmpfs (shmem).
          Active_file: File-backed memory on active LRU list.
          Active_anon: Anonymous and swap cache on active least-recently-used (LRU) list, including tmpfs.
        |||
      )
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Inactive_file_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Inactive_file',
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Inactive_anon_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Inactive_anon',
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Active_file_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Active_file',
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Active_anon_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Active_anon',
      )),

    local memoryCommited =
      nodeTimeseries.new(
        'Memory Commited',
        description=|||
          Committed_AS - Amount of memory presently allocated on the system.
          CommitLimit - Amount of memory currently available to be allocated on the system.
        |||
      )
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Committed_AS_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Committed_AS'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_CommitLimit_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='CommitLimit'
      )),
    local memorySharedAndMapped =
      nodeTimeseries.new(
        'Memory Shared and Mapped',
        description=|||
          Mapped: This refers to the memory used in mapped page files that have been memory mapped, such as libraries.
          Shmem: This is the memory used by shared memory, which is shared between multiple processes, including RAM disks.
          ShmemHugePages: This is the memory used by shared memory and tmpfs allocated with huge pages.
          ShmemPmdMapped: This is the amount of shared memory (shmem/tmpfs) backed by huge pages.
        |||
      )
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Mapped_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Mapped'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Shmem_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Shmem'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_ShmemHugePages_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ShmemHugePages'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_ShmemPmdMapped_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='ShmemPmdMapped'
      )),

    local memoryWriteAndDirty =
      nodeTimeseries.new(
        'Memory Writeback and Dirty',
        description=|||
          Writeback: This refers to the memory that is currently being actively written back to the disk.
          WritebackTmp: This is the memory used by FUSE for temporary writeback buffers.
          Dirty: This type of memory is waiting to be written back to the disk.
        |||
      )
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Writeback_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Writeback'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_WritebackTmp_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='WritebackTmp'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Dirty_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Dirty'
      )),

    local memoryVmalloc =
      nodeTimeseries.new(
        'Memory Vmalloc',
        description=|||
          Virtual Memory Allocation is a type of memory allocation in Linux that allows a process to request a contiguous block of memory larger than the amount of physically available memory. This is achieved by mapping the requested memory to virtual addresses that are backed by a combination of physical memory and swap space on disk.

          VmallocChunk: Largest contiguous block of vmalloc area which is free.
          VmallocTotal: Total size of vmalloc memory area.
          VmallocUsed: Amount of vmalloc area which is used.
        |||
      )
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_VmallocChunk_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='VmallocChunk'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_VmallocTotal_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='VmallocTotal'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_VmallocUsed_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='VmallocUsed'
      )),

    local memorySlab =
      nodeTimeseries.new('Memory Slab',
                         description=|||
                           Slab Allocation is a type of memory allocation in Linux that allows the kernel to efficiently manage the allocation and deallocation of small and frequently used data structures, such as network packets, file system objects, and process descriptors.

                           The Slab Allocator maintains a cache of pre-allocated objects of a fixed size and type, called slabs. When an application requests an object of a particular size and type, the Slab Allocator checks if a pre-allocated object of that size and type is available in the cache. If an object is available, it is returned to the application; if not, a new slab of objects is allocated and added to the cache.

                           SUnreclaim: Part of Slab, that cannot be reclaimed on memory pressure.
                           SReclaimable: Part of Slab, that might be reclaimed, such as caches.
                         |||)
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_SUnreclaim_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='SUnreclaim'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_SReclaimable_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='SReclaimable'
      )),

    local memoryAnonymous =
      nodeTimeseries.new(
        'Memory Anonymous',
        description=|||
          Memory Anonymous refers to the portion of the virtual memory that is used by a process for dynamically allocated memory that is not backed by any file or device.

          This type of memory is commonly used for heap memory allocation, which is used by programs to allocate and free memory dynamically during runtime.

          Memory Anonymous is different from Memory Mapped files, which refer to portions of the virtual memory space that are backed by a file or device,
          and from Memory Shared with other processes,
          which refers to memory regions that can be accessed and modified by multiple processes.

          AnonHugePages: Memory in anonymous huge pages.
          AnonPages: Memory in user pages not backed by files.
        |||
      )
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_AnonHugePages_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='AnonHugePages'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_AnonPages_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='AnonPages'
      )),

    local memoryHugePagesCounter =
      nodeTimeseries.new(
        'Memory HugePages Counter',
        description=|||
          Huge Pages are a feature that allows for the allocation of larger memory pages than the standard 4KB page size. By using larger page sizes, the kernel can reduce the overhead associated with managing a large number of smaller pages, which can improve system performance for certain workloads.

          HugePages_Free: Huge pages in the pool that are not yet allocated.
          HugePages_Rsvd: Huge pages for which a commitment to allocate from the pool has been made, but no allocation has yet been made.
          HugePages_Surp: Huge pages in the pool above the value in /proc/sys/vm/nr_hugepages.
        |||
      )
      .addTarget(commonPromTarget(
        expr='node_memory_HugePages_Free{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='HugePages_Free'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_HugePages_Rsvd{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='HugePages_Rsvd'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_HugePages_Surp{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='HugePages_Surp'
      )),
    local memoryHugePagesSize =
      nodeTimeseries.new(
        'Memory HugePages Size',
        description=|||
          Huge Pages are a feature that allows for the allocation of larger memory pages than the standard 4KB page size. By using larger page sizes, the kernel can reduce the overhead associated with managing a large number of smaller pages, which can improve system performance for certain workloads.
        |||
      )
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_HugePages_Total{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Huge pages total size'
      ))
      .addTarget(commonPromTarget(
        expr='node_memory_Hugepagesize_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Huge page size'
      )),
    local memoryDirectMap =
      nodeTimeseries.new(
        'Memory Direct Map',
        description=|||
          Direct Map memory refers to the portion of the kernel's virtual address space that is directly mapped to physical memory. This mapping is set up by the kernel during boot time and is used to provide fast access to certain critical kernel data structures, such as page tables and interrupt descriptor tables.
        |||
      )
      .withUnits('decbytes')
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
      nodeTimeseries.new(
        'Memory Bounce',
        description=|||
          Memory bounce is a technique used in the Linux kernel to handle situations where direct memory access (DMA) is required but the physical memory being accessed is not contiguous. This can happen when a device, such as a network interface card or a disk controller, requires access to a large amount of memory that is not available as a single contiguous block.

          To handle this situation, the kernel uses a technique called memory bouncing. In memory bouncing, the kernel sets up a temporary buffer in physical memory that is large enough to hold the entire data block being transferred by the device. The data is then copied from the non-contiguous source memory to the temporary buffer, which is physically contiguous.

          Bounce: Memory used for block device bounce buffers.
        |||
      )
      .withUnits('decbytes')
      .addTarget(commonPromTarget(
        expr='node_memory_Bounce_bytes{%(nodeQuerySelector)s}' % config { nodeQuerySelector: c.nodeQuerySelector },
        legendFormat='Bounce'
      )),
    local panelsGrid =
      [
        c.panelsWithTargets.memoryGauge { gridPos: { x: 0, w: 6, h: 6, y: 0 } },
        c.panelsWithTargets.memoryGraph { gridPos: { x: 6, w: 18, h: 6, y: 0 } },
        { type: 'row', title: 'Vmstat', gridPos: { y: 25 } },
        memoryPagesInOut { gridPos: { x: 0, w: 12, h: 8, y: 25 } },
        memoryPagesSwapInOut { gridPos: { x: 12, w: 12, h: 8, y: 25 } },
        memoryPagesFaults { gridPos: { x: 0, w: 12, h: 8, y: 25 } },
        memoryOOMkiller { gridPos: { x: 12, w: 12, h: 8, y: 25 } },
        { type: 'row', title: 'Memstat', gridPos: { y: 50 } },
        memoryActiveInactive { gridPos: { x: 0, w: 12, h: 8, y: 50 } },
        memoryActiveInactiveDetail { gridPos: { x: 12, w: 12, h: 8, y: 50 } },
        memoryCommited { gridPos: { x: 0, w: 12, h: 8, y: 50 } },
        memorySharedAndMapped { gridPos: { x: 12, w: 12, h: 8, y: 50 } },
        memoryWriteAndDirty { gridPos: { x: 0, w: 12, h: 8, y: 50 } },
        memoryVmalloc { gridPos: { x: 12, w: 12, h: 8, y: 50 } },
        memorySlab { gridPos: { x: 0, w: 12, h: 8, y: 50 } },
        memoryAnonymous { gridPos: { x: 12, w: 12, h: 8, y: 50 } },
        memoryHugePagesCounter { gridPos: { x: 0, w: 12, h: 8, y: 50 } },
        memoryHugePagesSize { gridPos: { x: 12, w: 12, h: 8, y: 50 } },
        memoryDirectMap { gridPos: { x: 0, w: 12, h: 8, y: 50 } },
        memoryBounce { gridPos: { x: 12, w: 12, h: 8, y: 50 } },
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Memory' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid=config.grafanaDashboardIDs['nodes-memory.json'],
      )
      .addLink(c.links.fleetDash)
      .addLink(c.links.nodeDash)
      .addLink(c.links.otherDashes)
      .addAnnotations(c.annotations)
      .addTemplates(templates)
      .addPanels(panelsGrid)
    else if platform == 'Darwin' then {},
  },
}
