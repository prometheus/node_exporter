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

      diskTotalRoot:
        commonlib.panels.disk.stat.total.new(
          'Root mount size',
          targets=[t.disk.diskTotalRoot],
          description=|||
            Total capacity on the primary mount point /.
          |||
        ),
      diskUsage:
        commonlib.panels.disk.table.usage.new(
          totalTarget=
          (
            t.disk.diskTotal
            + g.query.prometheus.withFormat('table')
            + g.query.prometheus.withInstant(true)
          ),
          freeTarget=
          t.disk.diskFree
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
            t.disk.diskFree,
          ],
          description='Filesystem space utilisation in bytes, by mountpoint.'
        ),
      diskInodesFree:
        commonlib.panels.disk.timeSeries.base.new(
          'Free inodes',
          targets=[t.disk.diskInodesFree],
          description='The inode is a data structure in a Unix-style file system that describes a file-system object such as a file or a directory.'
        )
        + g.panel.timeSeries.standardOptions.withUnit('short'),
      diskInodesTotal:
        commonlib.panels.disk.timeSeries.base.new(
          'Total inodes',
          targets=[t.disk.diskInodesTotal],
          description='The inode is a data structure in a Unix-style file system that describes a file-system object such as a file or a directory.',
        )
        + g.panel.timeSeries.standardOptions.withUnit('short'),
      diskErrorsandRO:
        commonlib.panels.disk.timeSeries.base.new(
          'Filesystems with errors / read-only',
          targets=[
            t.disk.diskDeviceError,
            t.disk.diskReadOnly,
          ],
          description='',
        )
        + g.panel.timeSeries.standardOptions.withMax(1),
      fileDescriptors:
        commonlib.panels.disk.timeSeries.base.new(
          'File descriptors',
          targets=[
            t.disk.processMaxFds,
            t.disk.processOpenFds,
          ],
          description=|||
            File descriptor is a handle to an open file or input/output (I/O) resource, such as a network socket or a pipe.
            The operating system uses file descriptors to keep track of open files and I/O resources, and provides a way for programs to read from and write to them.
          |||
        ),
      diskUsagePercentTopK: commonlib.panels.generic.timeSeries.topkPercentage.new(
        title='Disk space usage',
        target=t.disk.diskUsagePercent,
        topk=25,
        instanceLabels=this.config.instanceLabels + ['volume'],
        drillDownDashboardUid=this.grafana.dashboards['nodes.json'].uid,
      ),
      diskIOBytesPerSec: commonlib.panels.disk.timeSeries.ioBytesPerSec.new(
        targets=[t.disk.diskIOreadBytesPerSec, t.disk.diskIOwriteBytesPerSec, t.disk.diskIOutilization]
      ),
      diskIOutilPercentTopK:
        commonlib.panels.generic.timeSeries.topkPercentage.new(
          title='Disk IO',
          target=t.disk.diskIOutilization,
          topk=25,
          instanceLabels=this.config.instanceLabels + ['volume'],
          drillDownDashboardUid=this.grafana.dashboards['nodes.json'].uid,
        ),
      diskIOps:
        commonlib.panels.disk.timeSeries.iops.new(
          targets=[
            t.disk.diskIOReads,
            t.disk.diskIOWrites,
          ]
        ),

      diskQueue:
        commonlib.panels.disk.timeSeries.ioQueue.new(
          'Disk average queue',
          targets=
          [
            t.disk.diskAvgQueueSize,
          ]
        ),
      diskIOWaitTime: commonlib.panels.disk.timeSeries.ioWaitTime.new(
        targets=[
          t.disk.diskIOWaitReadTime,
          t.disk.diskIOWaitWriteTime,
        ]
      ),
    },
}
