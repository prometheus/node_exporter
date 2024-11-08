local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
local xtd = import 'github.com/jsonnet-libs/xtd/main.libsonnet';
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = xtd.array.slice(this.config.instanceLabels, -1)[0],

      memoryTotalBytes: commonlib.panels.memory.stat.total.new(targets=[t.memory.memoryTotalBytes]),
      memorySwapTotalBytes:
        commonlib.panels.memory.stat.total.new(
          'Total swap',
          targets=[t.memory.memorySwapTotal],
          description=|||
            Total swap available.

            Swap is a space on a storage device (usually a dedicated swap partition or a swap file) 
            used as virtual memory when the physical RAM (random-access memory) is fully utilized.
            Swap space helps prevent memory-related performance issues by temporarily transferring less-used data from RAM to disk,
            freeing up physical memory for active processes and applications.
          |||
        ),
      memoryUsageStatPercent: commonlib.panels.memory.stat.usage.new(targets=[t.memory.memoryUsagePercent]),
      memotyUsageTopKPercent: commonlib.panels.generic.timeSeries.topkPercentage.new(
        title='Memory usage',
        target=t.memory.memoryUsagePercent,
        topk=25,
        instanceLabels=this.config.instanceLabels,
        drillDownDashboardUid=this.grafana.dashboards['nodes.json'].uid,
      ),
      memoryUsageTsBytes:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          targets=[
            t.memory.memoryUsedBytes,
            t.memory.memoryCachedBytes,
            t.memory.memoryAvailableBytes,
            t.memory.memoryBuffersBytes,
            t.memory.memoryFreeBytes,
            t.memory.memoryTotalBytes,
          ],
          description=
          |||
            - Used: The amount of physical memory currently in use by the system.
            - Cached: The amount of physical memory used for caching data from disk. The Linux kernel uses available memory to cache data that is read from or written to disk. This helps speed up disk access times.
            - Free: The amount of physical memory that is currently not in use.
            - Buffers: The amount of physical memory used for temporary storage of data being transferred between devices or applications.
            - Available: The amount of physical memory that is available for use by applications. This takes into account memory that is currently being used for caching but can be freed up if needed.
          |||
        )
        + g.panel.timeSeries.standardOptions.withOverridesMixin(
          {
            __systemRef: 'hideSeriesFrom',
            matcher: {
              id: 'byNames',
              options: {
                mode: 'exclude',
                names: [
                  t.memory.memoryTotalBytes.legendFormat,
                  t.memory.memoryUsedBytes.legendFormat,
                ],
                prefix: 'All except:',
                readOnly: true,
              },
            },
            properties: [
              {
                id: 'custom.hideFrom',
                value: {
                  viz: true,
                  legend: false,
                  tooltip: false,
                },
              },
            ],
          }
        ),

      memoryPagesInOut:
        commonlib.panels.memory.timeSeries.base.new(
          'Memory pages in / out',
          targets=[t.memory.memoryPagesIn, t.memory.memoryPagesOut],
          description=|||
            Page-In - Return of pages to physical memory. This is a common and normal event.

            Page-Out - process of writing pages to disk. Unlike page-in, page-outs can indicate trouble.
            When the kernel detects low memory, it attempts to free memory by paging out. 
            While occasional page-outs are normal, excessive and frequent page-outs can lead to thrashing.
            Thrashing is a state in which the kernel spends more time managing paging activity than running applications, resulting in poor system performance.
          |||
        )
        + commonlib.panels.network.timeSeries.base.withNegateOutPackets(),

      memoryPagesSwapInOut:
        commonlib.panels.memory.timeSeries.base.new(
          'Memory pages swapping in / out',
          targets=[t.memory.memoryPagesSwapIn, t.memory.memoryPagesSwapOut],
          description=|||
            Compared to the speed of the CPU and main memory, writing pages out to disk is relatively slow.
            Nonetheless, it is a preferable option to crashing or killing off processes.

            The process of writing pages out to disk to free memory is known as swapping-out.
            If a page fault occurs because the page is on disk, in the swap area, rather than in memory,
            the kernel will read the page back in from the disk to satisfy the page fault. 
            This is known as swapping-in.
          |||
        )
        + commonlib.panels.network.timeSeries.base.withNegateOutPackets(),

      memoryPagesFaults:
        commonlib.panels.memory.timeSeries.base.new(
          'Memory page faults',
          targets=[t.memory.memoryPageMajorFaults, t.memory.memoryPageMinorFaults],
          description=|||
            A page fault is an exception raised by the memory when a process accesses a memory page without the necessary preparations,
            requiring a mapping to be added to the process's virtual address space.

            The page contents may also need to be loaded from a backing store such as a disk.
            While the MMU detects the page fault, the operating system's kernel handles the exception by either making the required page accessible in physical memory or denying an illegal memory access.
            Valid page faults are common and necessary to increase memory availability in any operating system that uses virtual memory, including Windows, macOS, and the Linux kernel.
          |||,
        ),

      memoryOOMkiller:
        commonlib.panels.memory.timeSeries.base.new(
          'OOM Killer',
          targets=[t.events.memoryOOMkiller],
          description=|||
            Out Of Memory killer is a process used by the Linux kernel when the system is running critically low on memory.

            This can happen when the kernel has allocated more memory than is available for its processes.
          |||
        ),

      memoryActiveInactive:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory active / inactive',
          targets=[t.memory.memoryActiveBytes, t.memory.memoryInactiveBytes],
          description=|||
            - Inactive: Memory which has been less recently used. It is more eligible to be reclaimed for other purposes.
            - Active: Memory that has been used more recently and usually not reclaimed unless absolutely necessary.
          |||,
        ),

      memoryActiveInactiveDetail:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory active / inactive details',
          targets=[t.memory.memoryInactiveFile, t.memory.memoryInactiveAnon, t.memory.memoryActiveFile, t.memory.memoryActiveAnon],
          description=|||
            - Inactive_file: File-backed memory on inactive LRU list.
            - Inactive_anon: Anonymous and swap cache on inactive LRU list, including tmpfs (shmem).
            - Active_file: File-backed memory on active LRU list.
            - Active_anon: Anonymous and swap cache on active least-recently-used (LRU) list, including tmpfs.
          |||,
        ),

      memoryCommited:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory commited',
          targets=[t.memory.memoryCommitedAs, t.memory.memoryCommitedLimit],
          description=|||
            - Committed_AS - Amount of memory presently allocated on the system.
            - CommitLimit - Amount of memory currently available to be allocated on the system.
          |||
        ),

      memorySharedAndMapped:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory shared and mapped',
          targets=[t.memory.memoryMappedBytes, t.memory.memoryShmemBytes, t.memory.memoryShmemPmdMappedBytes, t.memory.memoryShmemHugePagesBytes],
          description=|||
            - Mapped: This refers to the memory used in mapped page files that have been memory mapped, such as libraries.
            - Shmem: This is the memory used by shared memory, which is shared between multiple processes, including RAM disks.
            - ShmemHugePages: This is the memory used by shared memory and tmpfs allocated with huge pages.
            - ShmemPmdMapped: This is the amount of shared memory (shmem/tmpfs) backed by huge pages.
          |||
        ),
      memoryWriteAndDirty:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory writeback and dirty',
          targets=[t.memory.memoryWriteback, t.memory.memoryWritebackTmp, t.memory.memoryDirty],
          description=|||
            - Writeback: This refers to the memory that is currently being actively written back to the disk.
            - WritebackTmp: This is the memory used by FUSE for temporary writeback buffers.
            - Dirty: This type of memory is waiting to be written back to the disk.
          |||
        ),
      memoryVmalloc:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory Vmalloc',
          targets=[t.memory.memoryVmallocChunk, t.memory.memoryVmallocTotal, t.memory.memoryVmallocUsed],
          description=|||
            Virtual Memory Allocation is a type of memory allocation in Linux that allows a process to request a contiguous block of memory larger than the amount of physically available memory. This is achieved by mapping the requested memory to virtual addresses that are backed by a combination of physical memory and swap space on disk.

            - VmallocChunk: Largest contiguous block of vmalloc area which is free.
            - VmallocTotal: Total size of vmalloc memory area.
            - VmallocUsed: Amount of vmalloc area which is used.
          |||
        ),
      memorySlab:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory slab',
          targets=[t.memory.memorySlabSUnreclaim, t.memory.memorySlabSReclaimable],
          description=|||
            Slab Allocation is a type of memory allocation in Linux that allows the kernel to efficiently manage the allocation and deallocation of small and frequently used data structures, such as network packets, file system objects, and process descriptors.

            The Slab Allocator maintains a cache of pre-allocated objects of a fixed size and type, called slabs. When an application requests an object of a particular size and type, the Slab Allocator checks if a pre-allocated object of that size and type is available in the cache. If an object is available, it is returned to the application; if not, a new slab of objects is allocated and added to the cache.

            - SUnreclaim: Part of Slab, that cannot be reclaimed on memory pressure.
            - SReclaimable: Part of Slab, that might be reclaimed, such as caches.
          |||
        ),
      memoryAnonymous:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory anonymous',
          targets=[t.memory.memoryAnonHugePages, t.memory.memoryAnonPages],
          description=|||
            Memory Anonymous refers to the portion of the virtual memory that is used by a process for dynamically allocated memory that is not backed by any file or device.

            This type of memory is commonly used for heap memory allocation, which is used by programs to allocate and free memory dynamically during runtime.

            Memory Anonymous is different from Memory Mapped files, which refer to portions of the virtual memory space that are backed by a file or device,
            and from Memory Shared with other processes,
            which refers to memory regions that can be accessed and modified by multiple processes.

            - AnonHugePages: Memory in anonymous huge pages.
            - AnonPages: Memory in user pages not backed by files.
          |||
        ),

      memoryHugePagesCounter:
        commonlib.panels.memory.timeSeries.base.new(
          'Memory HugePages counter',
          targets=[t.memory.memoryHugePages_Free, t.memory.memoryHugePages_Rsvd, t.memory.memoryHugePages_Surp],
          description=
          |||
            Huge Pages are a feature that allows for the allocation of larger memory pages than the standard 4KB page size. By using larger page sizes, the kernel can reduce the overhead associated with managing a large number of smaller pages, which can improve system performance for certain workloads.

            - HugePages_Free: Huge pages in the pool that are not yet allocated.
            - HugePages_Rsvd: Huge pages for which a commitment to allocate from the pool has been made, but no allocation has yet been made.
            - HugePages_Surp: Huge pages in the pool above the value in /proc/sys/vm/nr_hugepages.
          |||
        ),
      memoryHugePagesSize:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory HugePages size',
          targets=[t.memory.memoryHugePagesTotalSize, t.memory.memoryHugePagesSize],

          description=|||
            Huge Pages are a feature that allows for the allocation of larger memory pages than the standard 4KB page size. By using larger page sizes, the kernel can reduce the overhead associated with managing a large number of smaller pages, which can improve system performance for certain workloads.
          |||
        ),

      memoryDirectMap:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory direct map',
          targets=[t.memory.memoryDirectMap1G, t.memory.memoryDirectMap2M, t.memory.memoryDirectMap4k],

          description=|||
            Direct Map memory refers to the portion of the kernel's virtual address space that is directly mapped to physical memory. This mapping is set up by the kernel during boot time and is used to provide fast access to certain critical kernel data structures, such as page tables and interrupt descriptor tables.
          |||
        ),
      memoryBounce:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory bounce',
          targets=[t.memory.memoryBounce],
          description=|||
            Memory bounce is a technique used in the Linux kernel to handle situations where direct memory access (DMA) is required but the physical memory being accessed is not contiguous. This can happen when a device, such as a network interface card or a disk controller, requires access to a large amount of memory that is not available as a single contiguous block.

            To handle this situation, the kernel uses a technique called memory bouncing. In memory bouncing, the kernel sets up a temporary buffer in physical memory that is large enough to hold the entire data block being transferred by the device. The data is then copied from the non-contiguous source memory to the temporary buffer, which is physically contiguous.

            - Bounce: Memory used for block device bounce buffers.
          |||
        ),
    },
}
