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

  // https://www.robustperception.io/filesystem-metrics-from-the-node-exporter/
  new(config=null, platform=null):: {
    local c = common.new(config=config, platform=platform),
    local commonPromTarget = c.commonPromTarget,
    local templates = c.templates,
    local q = c.queries,

    local fsAvailable =
      nodeTimeseries.new(
        'Filesystem Space Available',
        description=|||
          Filesystem space utilisation in bytes, by mountpoint.
        |||
      )
      .withUnits('decbytes')
      .withFillOpacity(5)
      .addTarget(commonPromTarget(
        expr=q.node_filesystem_avail_bytes,
        legendFormat='{{ mountpoint }}',
      )),

    local fsInodes =
      nodeTimeseries.new(
        'Free inodes',
        description='The inode is a data structure in a Unix-style file system that describes a file-system object such as a file or a directory.',
      )
      .withUnits('short')
      .addTarget(commonPromTarget(
        expr=q.node_filesystem_files_free,
        legendFormat='{{ mountpoint }}'
      ))
      .addTarget(commonPromTarget(
        expr=q.node_filesystem_files,
        legendFormat='{{ mountpoint }}'
      )),
    local fsInodesTotal =
      nodeTimeseries.new(
        'Total inodes',
        description='The inode is a data structure in a Unix-style file system that describes a file-system object such as a file or a directory.',
      )
      .withUnits('short')
      .addTarget(commonPromTarget(
        expr=q.node_filesystem_files,
        legendFormat='{{ mountpoint }}'
      )),
    local fsErrorsandRO =
      nodeTimeseries.new('Filesystems with errors / read-only')
      .withMax(1)
      .addTarget(commonPromTarget(
        expr=q.node_filesystem_readonly,
        legendFormat='{{ mountpoint }}'
      ))
      .addTarget(commonPromTarget(
        expr=q.node_filesystem_device_error,
        legendFormat='{{ mountpoint }}'
      )),
    local fileDescriptors =
      nodeTimeseries.new(
        'File Descriptors',
        description=|||
          File descriptor is a handle to an open file or input/output (I/O) resource, such as a network socket or a pipe.
          The operating system uses file descriptors to keep track of open files and I/O resources, and provides a way for programs to read from and write to them.
        |||
      )
      .addTarget(commonPromTarget(
        expr=q.process_max_fds,
        legendFormat='Maximum open file descriptors',
      ))
      .addTarget(commonPromTarget(
        expr=q.process_open_fds,
        legendFormat='Open file descriptors',
      )),

    local diskIOcompleted =
      nodeTimeseries.new(
        title='Disk IOps completed',
        description='The number (after merges) of I/O requests completed per second for the device'
      )
      .withUnits('iops')
      .withNegativeYByRegex('reads')
      .withAxisLabel('read(-) | write(+)')
      .addTarget(commonPromTarget(
        expr=q.node_disk_reads_completed_total,
        legendFormat='{{device}} reads completed',
      ))
      .addTarget(commonPromTarget(
        expr=q.node_disk_writes_completed_total,
        legendFormat='{{device}} writes completed',
      )),

    local diskAvgWaitTime =
      nodeTimeseries.new(
        title='Disk Average Wait Time',
        description='The average time for requests issued to the device to be served. This includes the time spent by the requests in queue and the time spent servicing them.'
      )
      .withUnits('s')
      .withNegativeYByRegex('read')
      .withAxisLabel('read(-) | write(+)')
      .addTarget(commonPromTarget(
        expr=q.diskWaitReadTime,
        legendFormat='{{device}} read wait time avg',
      ))
      .addTarget(commonPromTarget(
        expr=q.diskWaitWriteTime,
        legendFormat='{{device}} write wait time avg',
      )),

    local diskAvgQueueSize =
      nodeTimeseries.new(
        title='Average Queue Size (aqu-sz)',
        description='The average queue length of the requests that were issued to the device.'
      )
      .addTarget(commonPromTarget(
        expr=q.diskAvgQueueSize,
        legendFormat='{{device}}',
      )),

    local panelsGrid =
      [
        { type: 'row', title: 'Filesystem', gridPos: { y: 0 } },
        fsAvailable { gridPos: { x: 0, w: 12, h: 8, y: 0 } },
        c.panelsWithTargets.diskSpaceUsage { gridPos: { x: 12, w: 12, h: 8, y: 0 } },
        fsInodes { gridPos: { x: 0, w: 12, h: 8, y: 0 } },
        fsInodesTotal { gridPos: { x: 12, w: 12, h: 8, y: 0 } },
        fsErrorsandRO { gridPos: { x: 0, w: 12, h: 8, y: 0 } },
        fileDescriptors { gridPos: { x: 12, w: 12, h: 8, y: 0 } },
        { type: 'row', title: 'Disk', gridPos: { y: 25 } },
        c.panelsWithTargets.diskIO { gridPos: { x: 0, w: 12, h: 8, y: 25 } },
        diskIOcompleted { gridPos: { x: 12, w: 12, h: 8, y: 25 } },
        diskAvgWaitTime { gridPos: { x: 0, w: 12, h: 8, y: 25 } },
        diskAvgQueueSize { gridPos: { x: 12, w: 12, h: 8, y: 25 } },
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Filesystem and Disk' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid=config.grafanaDashboardIDs['nodes-disk.json']
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
