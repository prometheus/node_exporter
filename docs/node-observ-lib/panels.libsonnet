local g = import './g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
local utils = commonlib.utils;
{
  new(this):
    {
      local t = this.grafana.targets,
      local table = g.panel.table,
      local fieldOverride = g.panel.table.fieldOverride,
      local instanceLabel = this.config.instanceLabels[0],
      fleetOverviewTable:
        commonlib.panels.generic.table.base.new(
          'Fleet overview',
          targets=
          [
            t.osInfoCombined
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('OS Info'),
            t.uptime
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Uptime'),
            t.systemLoad1
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Load 1'),
            t.cpuCount
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Cores'),
            t.cpuUsage
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('CPU usage'),
            t.memoryTotalBytes
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Memory total'),
            t.memoryUsagePercent
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Memory usage'),
            t.diskTotalRoot
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Root mount size'),
            t.diskUsageRootPercent
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Root mount used'),
            t.alertsCritical
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('CRITICAL'),
            t.alertsWarning
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('WARNING'),
          ],
          description="All nodes' perfomance at a glance."
        )
        + g.panel.table.options.withFooter(
          value={
            reducer: ['sum'],
            show: true,
            fields: [
              'Value #Cores',
              'Value #Load 1',
              'Value #Memory total',
              'Value #Root mount size',
            ],
          }
        )
        + commonlib.panels.system.table.uptime.stylizeByName('Uptime')
        + table.standardOptions.withOverridesMixin([
          fieldOverride.byRegexp.new('Product|^Hostname$')
          + fieldOverride.byRegexp.withProperty('custom.filterable', true),
          fieldOverride.byName.new('Instance')
          + fieldOverride.byName.withProperty('custom.filterable', true)
          + fieldOverride.byName.withProperty('links', [
            {
              targetBlank: false,
              title: 'Drill down to ${__field.name} ${__value.text}',
              url: 'd/%s?var-%s=${__data.fields.%s}&${__url_time_range}' % [this.grafana.dashboards.overview.uid, instanceLabel, instanceLabel],
            },
          ]),
          fieldOverride.byRegexp.new(std.join('|', std.map(utils.toSentenceCase, this.config.groupLabels)))
          + fieldOverride.byRegexp.withProperty('custom.filterable', true)
          + fieldOverride.byRegexp.withProperty('links', [
            {
              targetBlank: false,
              title: 'Filter by ${__field.name}',
              url: 'd/%s?var-${__field.name}=${__value.text}&${__url_time_range}' % [this.grafana.dashboards.fleet.uid],
            },
          ]),
          fieldOverride.byName.new('Cores')
          + fieldOverride.byName.withProperty('custom.width', '120'),
          fieldOverride.byName.new('CPU usage')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withProperty('custom.displayMode', 'basic')
          + fieldOverride.byName.withPropertiesFromOptions(
            commonlib.panels.cpu.timeSeries.utilization.stylize()
          ),
          fieldOverride.byName.new('Memory total')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withPropertiesFromOptions(
            table.standardOptions.withUnit('bytes')
          ),
          fieldOverride.byName.new('Memory usage')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withProperty('custom.displayMode', 'basic')
          + fieldOverride.byName.withPropertiesFromOptions(
            commonlib.panels.cpu.timeSeries.utilization.stylize()
          ),
          fieldOverride.byName.new('Root mount size')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withPropertiesFromOptions(
            table.standardOptions.withUnit('bytes')
          ),
          fieldOverride.byName.new('Root mount used')
          + fieldOverride.byName.withProperty('custom.width', '120')
          + fieldOverride.byName.withProperty('custom.displayMode', 'basic')
          + fieldOverride.byName.withPropertiesFromOptions(
            table.standardOptions.withUnit('percent')
          )
          + fieldOverride.byName.withPropertiesFromOptions(
            commonlib.panels.cpu.timeSeries.utilization.stylize()
          ),
        ])
        + table.queryOptions.withTransformationsMixin(
          [
            {
              id: 'joinByField',
              options: {
                byField: instanceLabel,
                mode: 'outer',
              },
            },
            {
              id: 'filterFieldsByName',
              options: {
                include: {
                  //' 1' - would only match first occurence of group label, so no duplicates
                  pattern: instanceLabel + '|'
                           +
                           std.join(
                             '|',
                             std.map(
                               function(x) '%s 1' % x, this.config.instanceLabels
                             )
                           )
                           + '|' +
                           std.join(
                             '|',
                             std.map(
                               function(x) '%s 1' % x, this.config.groupLabels
                             )
                           )
                           + '|product|^hostname$|^nodename$|^pretty_name$|Value.+',
                },
              },
            },
            {
              id: 'organize',
              options: {
                excludeByName: {
                  'Value #OS Info': true,
                },
                indexByName:
                  {
                    [instanceLabel]: 0,
                    nodename: 1,
                    hostname: 1,
                    pretty_name: 2,
                    product: 2,
                  }
                  +
                  // group labels are named as 'job 1' and so on.
                  {
                    [label]: 3
                    for label in this.config.groupLabels
                  },
                renameByName:
                  {
                    [label + ' 1']: utils.toSentenceCase(label)
                    for label in this.config.instanceLabels
                  }
                  {
                    [instanceLabel]: utils.toSentenceCase(instanceLabel),
                    product: 'OS',  // windows
                    pretty_name: 'OS',  // linux
                    hostname: 'Hostname',  // windows
                    nodename: 'Hostname',  // Linux
                  }
                  +
                  // group labels are named as 'job 1' and so on.
                  {
                    [label + ' 1']: utils.toSentenceCase(label)
                    for label in this.config.groupLabels
                  },

              },
            },
            {
              id: 'renameByRegex',
              options: {
                regex: 'Value #(.*)',
                renamePattern: '$1',
              },
            },
          ]
        ),
      uptime: commonlib.panels.system.stat.uptime.new(targets=[t.uptime]),

      systemLoad:
        commonlib.panels.system.timeSeries.loadAverage.new(
          loadTargets=[t.systemLoad1, t.systemLoad5, t.systemLoad15],
          cpuCountTarget=t.cpuCount,
        ),

      systemContextSwitchesAndInterrupts:
        commonlib.panels.generic.timeSeries.base.new(
          'Context switches/Interrupts',
          targets=[
            t.systemContextSwitches,
            t.systemInterrupts,
          ],
          description=|||
            Context switches occur when the operating system switches from running one process to another. Interrupts are signals sent to the CPU by external devices to request its attention.

            A high number of context switches or interrupts can indicate that the system is overloaded or that there are problems with specific devices or processes.
          |||
        ),

      timeNtpStatus:
        commonlib.panels.system.statusHistory.ntp.new(
          'NTP status',
          targets=[t.timeNtpStatus],
          description='Status of time synchronization.'
        )
        + g.panel.timeSeries.standardOptions.withNoValue('No data.')
        + g.panel.statusHistory.options.withLegend(false),
      timeSyncDrift:
        commonlib.panels.generic.timeSeries.base.new(
          'Time synchronized drift',
          targets=[
            t.timeEstimatedError,
            t.timeOffset,
            t.timeMaxError,
          ],
          description=|||
            Time synchronization is essential to ensure accurate timekeeping, which is critical for many system operations such as logging, authentication, and network communication, as well as distributed systems or clusters where data consistency is important.
          |||
        )
        + g.panel.timeSeries.standardOptions.withUnit('seconds')
        + g.panel.timeSeries.standardOptions.withNoValue('No data.'),
      cpuCount: commonlib.panels.cpu.stat.count.new(targets=[t.cpuCount]),
      cpuUsageTsPerCore: commonlib.panels.cpu.timeSeries.utilization.new(targets=[t.cpuUsagePerCore])
                         + g.panel.timeSeries.fieldConfig.defaults.custom.withStacking({ mode: 'normal' }),

      cpuUsageTopk: commonlib.panels.generic.timeSeries.topkPercentage.new(
        title='CPU usage',
        target=t.cpuUsage,
        topk=25,
        instanceLabels=this.config.instanceLabels,
        drillDownDashboardUid=this.grafana.dashboards.overview.uid,
      ),
      cpuUsageStat: commonlib.panels.cpu.stat.usage.new(targets=[t.cpuUsage]),
      cpuUsageByMode: commonlib.panels.cpu.timeSeries.utilizationByMode.new(
        targets=[t.cpuUsageByMode],
        description=|||
          - System: Processes executing in kernel mode.
          - User: Normal processes executing in user mode.
          - Nice: Niced processes executing in user mode.
          - Idle: Waiting for something to happen.
          - Iowait: Waiting for I/O to complete.
          - Irq: Servicing interrupts.
          - Softirq: Servicing softirqs.
          - Steal: Time spent in other operating systems when running in a virtualized environment.
        |||
      ),

      memoryTotalBytes: commonlib.panels.memory.stat.total.new(targets=[t.memoryTotalBytes]),
      memorySwapTotalBytes:
        commonlib.panels.memory.stat.total.new(
          'Total swap',
          targets=[t.memorySwapTotal],
          description=|||
            Total swap available.

            Swap is a space on a storage device (usually a dedicated swap partition or a swap file) 
            used as virtual memory when the physical RAM (random-access memory) is fully utilized.
            Swap space helps prevent memory-related performance issues by temporarily transferring less-used data from RAM to disk,
            freeing up physical memory for active processes and applications.
          |||
        ),
      memoryUsageStatPercent: commonlib.panels.memory.stat.usage.new(targets=[t.memoryUsagePercent]),
      memotyUsageTopKPercent: commonlib.panels.generic.timeSeries.topkPercentage.new(
        title='Memory usage',
        target=t.memoryUsagePercent,
        topk=25,
        instanceLabels=this.config.instanceLabels,
        drillDownDashboardUid=this.grafana.dashboards.overview.uid,
      ),
      memoryUsageTsBytes:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          targets=[
            t.memoryUsedBytes,
            t.memoryCachedBytes,
            t.memoryAvailableBytes,
            t.memoryBuffersBytes,
            t.memoryFreeBytes,
            t.memoryTotalBytes,
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
                  t.memoryTotalBytes.legendFormat,
                  t.memoryUsedBytes.legendFormat,
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
          targets=[t.memoryPagesIn, t.memoryPagesOut],
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
          targets=[t.memoryPagesSwapIn, t.memoryPagesSwapOut],
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
          targets=[t.memoryPageMajorFaults, t.memoryPageMinorFaults],
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
          targets=[t.memoryOOMkiller],
          description=|||
            Out Of Memory killer is a process used by the Linux kernel when the system is running critically low on memory.

            This can happen when the kernel has allocated more memory than is available for its processes.
          |||
        ),

      memoryActiveInactive:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory active / inactive',
          targets=[t.memoryActiveBytes, t.memoryInactiveBytes],
          description=|||
            - Inactive: Memory which has been less recently used. It is more eligible to be reclaimed for other purposes.
            - Active: Memory that has been used more recently and usually not reclaimed unless absolutely necessary.
          |||,
        ),

      memoryActiveInactiveDetail:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory active / inactive details',
          targets=[t.memoryInactiveFile, t.memoryInactiveAnon, t.memoryActiveFile, t.memoryActiveAnon],
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
          targets=[t.memoryCommitedAs, t.memoryCommitedLimit],
          description=|||
            - Committed_AS - Amount of memory presently allocated on the system.
            - CommitLimit - Amount of memory currently available to be allocated on the system.
          |||
        ),

      memorySharedAndMapped:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory shared and mapped',
          targets=[t.memoryMappedBytes, t.memoryShmemBytes, t.memoryShmemBytes, t.memoryShmemHugePagesBytes],
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
          targets=[t.memoryWriteback, t.memoryWritebackTmp, t.memoryDirty],
          description=|||
            - Writeback: This refers to the memory that is currently being actively written back to the disk.
            - WritebackTmp: This is the memory used by FUSE for temporary writeback buffers.
            - Dirty: This type of memory is waiting to be written back to the disk.
          |||
        ),
      memoryVmalloc:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory Vmalloc',
          targets=[t.memoryVmallocChunk, t.memoryVmallocTotal, t.memoryVmallocUsed],
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
          targets=[t.memorySlabSUnreclaim, t.memorySlabSReclaimable],
          description=|||
            Slab Allocation is a type of memory allocation in Linux that allows the kernel to efficiently manage the allocation and deallocation of small and frequently used data structures, such as network packets, file system objects, and process descriptors.

            The Slab Allocator maintains a cache of pre-allocated objects of a fixed size and type, called slabs. When an application requests an object of a particular size and type, the Slab Allocator checks if a pre-allocated object of that size and type is available in the cache. If an object is available, it is returned to the application; if not, a new slab of objects is allocated and added to the cache.

            - SUnreclaim: Part of Slab, that cannot be reclaimed on memory pressure.
            - SReclaimable: Part of Slab, that might be reclaimed, such as caches.
          |||
        ),
      memoryAnonymous:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory slab',
          targets=[t.memoryAnonHugePages, t.memoryAnonPages],
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
          targets=[t.memoryHugePages_Free, t.memoryHugePages_Rsvd, t.memoryHugePages_Surp],
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
          targets=[t.memoryHugePagesTotalSize, t.memoryHugePagesSize],

          description=|||
            Huge Pages are a feature that allows for the allocation of larger memory pages than the standard 4KB page size. By using larger page sizes, the kernel can reduce the overhead associated with managing a large number of smaller pages, which can improve system performance for certain workloads.
          |||
        ),

      memoryDirectMap:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory direct map',
          targets=[t.memoryDirectMap1G, t.memoryDirectMap2M, t.memoryDirectMap4k],

          description=|||
            Direct Map memory refers to the portion of the kernel's virtual address space that is directly mapped to physical memory. This mapping is set up by the kernel during boot time and is used to provide fast access to certain critical kernel data structures, such as page tables and interrupt descriptor tables.
          |||
        ),
      memoryBounce:
        commonlib.panels.memory.timeSeries.usageBytes.new(
          'Memory bounce',
          targets=[t.memoryBounce],
          description=|||
            Memory bounce is a technique used in the Linux kernel to handle situations where direct memory access (DMA) is required but the physical memory being accessed is not contiguous. This can happen when a device, such as a network interface card or a disk controller, requires access to a large amount of memory that is not available as a single contiguous block.

            To handle this situation, the kernel uses a technique called memory bouncing. In memory bouncing, the kernel sets up a temporary buffer in physical memory that is large enough to hold the entire data block being transferred by the device. The data is then copied from the non-contiguous source memory to the temporary buffer, which is physically contiguous.

            - Bounce: Memory used for block device bounce buffers.
          |||
        ),
      diskTotalRoot:
        commonlib.panels.disk.stat.total.new(
          'Root mount size',
          targets=[t.diskTotalRoot],
          description=|||
            Total capacity on the primary mount point /.
          |||
        ),
      diskUsage:
        commonlib.panels.disk.table.usage.new(
          totalTarget=
          (
            t.diskTotal
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
          ),
          freeTarget=
          t.diskFree
          + g.query.prometheus.withFormat('table')
          + g.query.prometheus.withInstant(true),
          groupLabel='mountpoint'
          ,
          description='Disk utilisation in percent, by mountpoint. Some duplication can occur if the same filesystem is mounted in multiple locations.'
        ),
      diskFreeTs:
        commonlib.panels.disk.timeSeries.available.new(
          'Filesystem space availabe',
          targets=[
            t.diskFree,
          ],
          description='Filesystem space utilisation in bytes, by mountpoint.'
        ),
      diskInodesFree:
        commonlib.panels.disk.timeSeries.base.new(
          'Free inodes',
          targets=[t.diskInodesFree],
          description='The inode is a data structure in a Unix-style file system that describes a file-system object such as a file or a directory.'
        )
        + g.panel.timeSeries.standardOptions.withUnit('short'),
      diskInodesTotal:
        commonlib.panels.disk.timeSeries.base.new(
          'Total inodes',
          targets=[t.diskInodesTotal],
          description='The inode is a data structure in a Unix-style file system that describes a file-system object such as a file or a directory.',
        )
        + g.panel.timeSeries.standardOptions.withUnit('short'),
      diskErrorsandRO:
        commonlib.panels.disk.timeSeries.base.new(
          'Filesystems with errors / read-only',
          targets=[
            t.diskDeviceError,
            t.diskReadOnly,
          ],
          description='',
        )
        + g.panel.timeSeries.standardOptions.withMax(1),
      fileDescriptors:
        commonlib.panels.disk.timeSeries.base.new(
          'File descriptors',
          targets=[
            t.processMaxFds,
            t.processOpenFds,
          ],
          description=|||
            File descriptor is a handle to an open file or input/output (I/O) resource, such as a network socket or a pipe.
            The operating system uses file descriptors to keep track of open files and I/O resources, and provides a way for programs to read from and write to them.
          |||
        ),
      diskUsagePercentTopK: commonlib.panels.generic.timeSeries.topkPercentage.new(
        title='Disk space usage',
        target=t.diskUsagePercent,
        topk=25,
        instanceLabels=this.config.instanceLabels + ['volume'],
        drillDownDashboardUid=this.grafana.dashboards.overview.uid,
      ),
      diskIOBytesPerSec: commonlib.panels.disk.timeSeries.ioBytesPerSec.new(
        targets=[t.diskIOreadBytesPerSec, t.diskIOwriteBytesPerSec, t.diskIOutilization]
      ),
      diskIOutilPercentTopK:
        commonlib.panels.generic.timeSeries.topkPercentage.new(
          title='Disk IO',
          target=t.diskIOutilization,
          topk=25,
          instanceLabels=this.config.instanceLabels + ['volume'],
          drillDownDashboardUid=this.grafana.dashboards.overview.uid,
        ),
      diskIOps:
        commonlib.panels.disk.timeSeries.iops.new(
          targets=[
            t.diskIOReads,
            t.diskIOWrites,
          ]
        ),

      diskQueue:
        commonlib.panels.disk.timeSeries.ioQueue.new(
          'Disk average queue',
          targets=
          [
            t.diskAvgQueueSize,
          ]
        ),
      diskIOWaitTime: commonlib.panels.disk.timeSeries.ioWaitTime.new(
        targets=[
          t.diskIOWaitReadTime,
          t.diskIOWaitWriteTime,
        ]
      ),
      osInfo: commonlib.panels.generic.stat.info.new(
        'OS',
        targets=[t.osInfo],
        description='Operating system'
      )
              { options+: { reduceOptions+: { fields: '/^pretty_name$/' } } },
      kernelVersion:
        commonlib.panels.generic.stat.info.new('Kernel version',
                                               targets=[t.unameInfo],
                                               description='Kernel version of linux host.')
        { options+: { reduceOptions+: { fields: '/^release$/' } } },
      osTimezone:
        commonlib.panels.generic.stat.info.new(
          'Timezone', targets=[t.osTimezone], description='Current system timezone.'
        )
        { options+: { reduceOptions+: { fields: '/^time_zone$/' } } },
      hostname:
        commonlib.panels.generic.stat.info.new(
          'Hostname',
          targets=[t.unameInfo],
          description="System's hostname."
        )
        { options+: { reduceOptions+: { fields: '/^nodename$/' } } },
      networkErrorsAndDroppedPerSec:
        commonlib.panels.network.timeSeries.errors.new(
          'Network errors and dropped packets',
          targets=[
            t.networkOutErrorsPerSec,
            t.networkInErrorsPerSec,
            t.networkOutDroppedPerSec,
            t.networkInDroppedPerSec,
          ],
          description=|||
            **Network errors**:

            Network errors refer to issues that occur during the transmission of data across a network. 

            These errors can result from various factors, including physical issues, jitter, collisions, noise and interference.

            Monitoring network errors is essential for diagnosing and resolving issues, as they can indicate problems with network hardware or environmental factors affecting network quality.

            **Dropped packets**:

            Dropped packets occur when data packets traveling through a network are intentionally discarded or lost due to congestion, resource limitations, or network configuration issues. 

            Common causes include network congestion, buffer overflows, QoS settings, and network errors, as corrupted or incomplete packets may be discarded by receiving devices.

            Dropped packets can impact network performance and lead to issues such as degraded voice or video quality in real-time applications.
          |||
        )
        + commonlib.panels.network.timeSeries.errors.withNegateOutPackets(),
      networkErrorsAndDroppedPerSecTopK:
        commonlib.panels.network.timeSeries.errors.new(
          'Network errors and dropped packets',
          targets=std.map(
            function(t) t
                        {
              expr: 'topk(25, ' + t.expr + ')>0.5',
              legendFormat: '{{' + this.config.instanceLabels[0] + '}}: ' + std.get(t, 'legendFormat', '{{ nic }}'),
            },
            [
              t.networkOutErrorsPerSec,
              t.networkInErrorsPerSec,
              t.networkOutDroppedPerSec,
              t.networkInDroppedPerSec,
            ]
          ),
          description=|||
            Top 25.

            **Network errors**:

            Network errors refer to issues that occur during the transmission of data across a network. 

            These errors can result from various factors, including physical issues, jitter, collisions, noise and interference.

            Monitoring network errors is essential for diagnosing and resolving issues, as they can indicate problems with network hardware or environmental factors affecting network quality.

            **Dropped packets**:

            Dropped packets occur when data packets traveling through a network are intentionally discarded or lost due to congestion, resource limitations, or network configuration issues. 

            Common causes include network congestion, buffer overflows, QoS settings, and network errors, as corrupted or incomplete packets may be discarded by receiving devices.

            Dropped packets can impact network performance and lead to issues such as degraded voice or video quality in real-time applications.
          |||
        )
        + g.panel.timeSeries.fieldConfig.defaults.custom.withDrawStyle('points')
        + g.panel.timeSeries.fieldConfig.defaults.custom.withPointSize(5),

      networkErrorsPerSec:
        commonlib.panels.network.timeSeries.errors.new(
          'Network errors',
          targets=[t.networkInErrorsPerSec, t.networkOutErrorsPerSec]
        )
        + commonlib.panels.network.timeSeries.errors.withNegateOutPackets(),
      networkDroppedPerSec:
        commonlib.panels.network.timeSeries.dropped.new(
          targets=[t.networkInDroppedPerSec, t.networkOutDroppedPerSec]
        )
        + commonlib.panels.network.timeSeries.errors.withNegateOutPackets(),
      networkUsagePerSec:
        commonlib.panels.network.timeSeries.traffic.new(
          targets=[t.networkInBitPerSec, t.networkOutBitPerSec]
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),
      networkPacketsPerSec:
        commonlib.panels.network.timeSeries.packets.new(
          targets=[t.networkInPacketsPerSec, t.networkOutPacketsPerSec]
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),
      networkMulticastPerSec:
        commonlib.panels.network.timeSeries.multicast.new(
          'Multicast packets',
          targets=[t.networkInMulticastPacketsPerSec, t.networkOutMulticastPacketsPerSec],
          description='Multicast packets received and transmitted.'
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),

      networkFifo:
        commonlib.panels.network.timeSeries.packets.new(
          'Network FIFO',
          targets=[t.networkFifoInPerSec, t.networkFifoOutPerSec],
          description=|||
            Network FIFO (First-In, First-Out) refers to a buffer used by the network stack to store packets in a queue.
            It is a mechanism used to manage network traffic and ensure that packets are delivered to their destination in the order they were received.
            Packets are stored in the FIFO buffer until they can be transmitted or processed further.
          |||
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),
      networkCompressedPerSec:
        commonlib.panels.network.timeSeries.packets.new(
          'Compressed packets',
          targets=[t.networkCompressedInPerSec, t.networkCompressedOutPerSec],
          description=|||
            - Compressed received: 
            Number of correctly received compressed packets. This counters is only meaningful for interfaces which support packet compression (e.g. CSLIP, PPP).

            - Compressed transmitted:
            Number of transmitted compressed packets. This counters is only meaningful for interfaces which support packet compression (e.g. CSLIP, PPP).

            https://docs.kernel.org/networking/statistics.html
          |||,
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets(),
      networkNFConntrack:
        commonlib.panels.generic.timeSeries.base.new(
          'NF conntrack',
          targets=[t.networkNFConntrackEntries, t.networkNFConntrackLimits],
          description=|||
            NF Conntrack is a component of the Linux kernel's netfilter framework that provides stateful packet inspection to track and manage network connections,
            enforce firewall rules, perform NAT, and manage network address/port translation.
          |||
        )
        + g.panel.timeSeries.fieldConfig.defaults.custom.withFillOpacity(0),

      networkSoftnet:
        commonlib.panels.network.timeSeries.packets.new(
          'Softnet packets',
          targets=[t.networkSoftnetProcessedPerSec, t.networkSoftnetDroppedPerSec],
          description=|||
            Softnet packets are received by the network and queued for processing by the kernel's networking stack.
            Softnet packets are usually generated by network traffic that is directed to the local host, and they are typically processed by the kernel's networking subsystem before being passed on to the relevant application. 
          |||
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets('/dropped/')
        + g.panel.timeSeries.fieldConfig.defaults.custom.withAxisLabel('Dropped(-) | Processed(+)'),
      networkSoftnetSqueeze:
        commonlib.panels.network.timeSeries.packets.new(
          'Softnet out of quota',
          targets=[t.networkSoftnetSqueezedPerSec],
          description=|||
            "Softnet out of quota" is a network-related metric in Linux that measures the number of times the kernel's softirq processing was unable to handle incoming network traffic due to insufficient softirq processing capacity.
            This means that the kernel has reached its processing capacity limit for incoming packets, and any additional packets will be dropped or deferred.
          |||
        ),
      networkOperStatus:
        commonlib.panels.network.statusHistory.interfaceStatus.new(
          'Network interfaces carrier status',
          targets=[t.networkCarrier],
          description='Network interfaces carrier status',
        ),
      networkOverviewTable:
        commonlib.panels.generic.table.base.new(
          'Network interfaces overview',
          targets=
          [
            t.networkUp
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Up'),
            t.networkCarrier
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Carrier'),
            t.networkOutBitPerSec
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(false)
            + g.query.prometheus.withRefId('Transmitted'),
            t.networkInBitPerSec
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(false)
            + g.query.prometheus.withRefId('Received'),
            t.networkArpEntries
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('ARP entries'),
            t.networkMtuBytes
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('MTU'),
            t.networkSpeedBitsPerSec
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Speed'),
            t.networkTransmitQueueLength
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Queue length'),
            t.networkInfo
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
            + g.query.prometheus.withRefId('Info'),
          ],
          description='Network interfaces overview.'
        )
        + g.panel.table.standardOptions.withOverridesMixin([
          fieldOverride.byName.new('Speed')
          + fieldOverride.byName.withPropertiesFromOptions(
            table.standardOptions.withUnit('bps')
          ),
        ])
        + g.panel.table.standardOptions.withOverridesMixin([
          fieldOverride.byRegexp.new('Transmitted|Received')
          + fieldOverride.byRegexp.withProperty('custom.displayMode', 'gradient-gauge')
          + fieldOverride.byRegexp.withPropertiesFromOptions(
            table.standardOptions.withUnit('bps')
            + table.standardOptions.color.withMode('continuous-BlYlRd')
            + table.standardOptions.withMax(1000 * 1000 * 100)
          ),
        ])
        + g.panel.table.standardOptions.withOverridesMixin([
          fieldOverride.byRegexp.new('Carrier|Up')
          + fieldOverride.byRegexp.withProperty('custom.displayMode', 'color-text')
          + fieldOverride.byRegexp.withPropertiesFromOptions(
            table.standardOptions.withMappings(
              {
                type: 'value',
                options: {
                  '0': {
                    text: 'Down',
                    color: 'light-red',
                    index: 0,
                  },
                  '1': {
                    text: 'Up',
                    color: 'light-green',
                    index: 1,
                  },
                },
              }
            ),
          ),
        ])
        + table.queryOptions.withTransformationsMixin(
          [
            {
              id: 'joinByField',
              options: {
                byField: 'device',
                mode: 'outer',
              },
            },
            {
              id: 'filterFieldsByName',
              options: {
                include: {
                  pattern: 'device|duplex|address|Value.+',
                },
              },
            },
            {
              id: 'renameByRegex',
              options: {
                regex: '(Value) #(.*)',
                renamePattern: '$2',
              },
            },
            {
              id: 'organize',
              options: {
                excludeByName: {
                  Info: true,
                },
                renameByName:
                  {
                    device: 'Interface',
                    duplex: 'Duplex',
                    address: 'Address',
                  },
              },
            },
            {
              id: 'organize',
              options: {
                indexByName: {
                  Interface: 0,
                  Up: 1,
                  Carrier: 2,
                  Received: 3,
                  Transmitted: 4,
                },
              },
            },
          ]
        ),
      networkSockstatAll:
        commonlib.panels.generic.timeSeries.base.new(
          'Sockets in use',
          targets=[t.networkSocketsUsed],
          description='Number of sockets currently in use.',
        ),

      networkSockstatTCP:
        commonlib.panels.generic.timeSeries.base.new(
          'Sockets TCP',
          targets=[t.networkSocketsTCPAllocated, t.networkSocketsTCPIPv4, t.networkSocketsTCPIPv6, t.networkSocketsTCPOrphans, t.networkSocketsTCPTimeWait],
          description=|||
            TCP sockets are used for establishing and managing network connections between two endpoints over the TCP/IP protocol.

            Orphan sockets: If a process terminates unexpectedly or is terminated without closing its sockets properly, the sockets may become orphaned.
          |||
        ),
      networkSockstatUDP:
        commonlib.panels.generic.timeSeries.base.new(
          'Sockets UDP',
          targets=[t.networkSocketsUDPLiteInUse, t.networkSocketsUDPInUse, t.networkSocketsUDPLiteIPv6InUse, t.networkSocketsUDPIPv6InUse],
          description=|||
            UDP (User Datagram Protocol) and UDPlite (UDP-Lite) sockets are used for transmitting and receiving data over the UDP and UDPlite protocols, respectively.
            Both UDP and UDPlite are connectionless protocols that do not provide a reliable data delivery mechanism.
          |||
        ),
      networkSockstatOther:
        commonlib.panels.generic.timeSeries.base.new(
          'Sockets other',
          targets=[t.networkSocketsFragInUse, t.networkSocketsFragIPv6InUse, t.networkSocketsRawInUse, t.networkSocketsIPv6RawInUse],
          description=|||
            FRAG (IP fragment) sockets: Used to receive and process fragmented IP packets. FRAG sockets are useful in network monitoring and analysis.

            RAW sockets: Allow applications to send and receive raw IP packets directly without the need for a transport protocol like TCP or UDP.
          |||
        ),
      networkSockstatMemory:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.generic.timeSeries.base.new(
          title='Sockets memory',
          targets=[t.networkSocketsTCPMemoryPages, t.networkSocketsUDPMemoryPages, t.networkSocketsTCPMemoryBytes, t.networkSocketsUDPMemoryBytes],
          description=|||
            Memory currently in use for sockets.
          |||,
        )
        + panel.queryOptions.withMaxDataPoints(100)
        + panel.fieldConfig.defaults.custom.withAxisLabel('Pages')
        + panel.standardOptions.withOverridesMixin(
          panel.standardOptions.override.byRegexp.new('/bytes/')
          + override.byType.withPropertiesFromOptions(
            panel.standardOptions.withDecimals(2)
            + panel.standardOptions.withUnit('bytes')
            + panel.fieldConfig.defaults.custom.withDrawStyle('bars')
            + panel.fieldConfig.defaults.custom.withStacking(value={ mode: 'normal', group: 'A' })
            + panel.fieldConfig.defaults.custom.withAxisLabel('Bytes')
          )
        ),

      networkNetstatIP:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.packets.new(
          'IP octets',
          targets=[t.networkNetstatIPInOctetsPerSec, t.networkNetstatIPOutOctetsPerSec],
          description='Rate of IP octets received and transmitted.'
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets()
        + panel.standardOptions.withUnit('oct/s'),

      networkNetstatTCP:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.packets.new(
          'TCP segments',
          targets=[t.networkNetstatTCPInSegmentsPerSec, t.networkNetstatTCPOutSegmentsPerSec],
          description='Rate of TCP segments received and transmitted.'
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets()
        + panel.standardOptions.withUnit('seg/s'),

      networkNetstatTCPerrors:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.errors.new(
          title='TCP errors rate',
          targets=[
            t.networkNetstatTCPOverflowPerSec,
            t.networkNetstatTCPListenDropsPerSec,
            t.networkNetstatTCPRetransPerSec,
            t.networkNetstatTCPRetransSegPerSec,
            t.networkNetstatTCPInWithErrorsPerSec,
            t.networkNetstatTCPOutWithRstPerSec,
          ],
          description='Rate of TCP errors.'
        )
        + panel.standardOptions.withUnit('err/s'),

      networkNetstatUDP:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.packets.new(
          'UDP datagrams',
          targets=[t.networkNetstatIPInUDPPerSec, t.networkNetstatIPOutUDPPerSec],
          description='Rate of UDP datagrams received and transmitted.'
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets()
        + panel.standardOptions.withUnit('dat/s'),

      networkNetstatUDPerrors:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.errors.new(
          title='UDP errors rate',
          targets=[
            t.networkNetstatUDPLiteInErrorsPerSec,
            t.networkNetstatUDPInErrorsPerSec,
            t.networkNetstatUDP6InErrorsPerSec,
            t.networkNetstatUDPNoPortsPerSec,
            t.networkNetstatUDP6NoPortsPerSec,
            t.networkNetstatUDPRcvBufErrsPerSec,
            t.networkNetstatUDP6RcvBufErrsPerSec,
            t.networkNetstatUDPSndBufErrsPerSec,
            t.networkNetstatUDP6SndBufErrsPerSec,
          ],
          description='Rate of UDP errors.'
        )
        + panel.standardOptions.withUnit('err/s'),

      networkNetstatICMP:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.packets.new(
          'ICMP messages',
          targets=[
            t.networkNetstatICMPInPerSec,
            t.networkNetstatICMPOutPerSec,
            t.networkNetstatICMP6InPerSec,
            t.networkNetstatICMP6OutPerSec,
          ],
          description="Rate of ICMP messages, like 'ping', received and transmitted."
        )
        + commonlib.panels.network.timeSeries.traffic.withNegateOutPackets()
        + panel.standardOptions.withUnit('msg/s'),

      networkNetstatICMPerrors:
        local panel = g.panel.timeSeries;
        local override = g.panel.timeSeries.standardOptions.override;
        commonlib.panels.network.timeSeries.errors.new(
          title='ICMP errors rate',
          targets=[
            t.networkNetstatICMPInErrorsPerSec,
            t.networkNetstatICM6PInErrorsPerSec,
          ],
          description='Rate of ICMP messages received and transmitted with errors.'
        )
        + panel.standardOptions.withUnit('err/s'),
    },
}
