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
      nodeTimeseries.new('Filesystem Space Available')
      .withUnits('decbytes')
      .withFillOpacity(5)
      .addTarget(commonPromTarget(
        expr=q.node_filesystem_avail_bytes,
        legendFormat='{{ mountpoint }}',
      )),

    local fsInodes =
      nodeTimeseries.new('Free inodes')
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
      nodeTimeseries.new('Total inodes')
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

    local panelsGrid =
      [
        { type: 'row', title: 'Filesystem', gridPos: { y: 0 } },
        fsAvailable { gridPos: { x: 0, w: 12, h: 8, y: 0 } },
        c.panelsWithTargets.fsSpaceUsage { gridPos: { x: 12, w: 12, h: 8, y: 0 } },
        fsInodes { gridPos: { x: 0, w: 12, h: 8, y: 0 } },
        fsInodesTotal { gridPos: { x: 12, w: 12, h: 8, y: 0 } },
        fsErrorsandRO { gridPos: { x: 0, w: 12, h: 8, y: 0 } },

        { type: 'row', title: 'Disk', gridPos: { y: 25 } },
      ],

    dashboard: if platform == 'Linux' then
      dashboard.new(
        '%sNode Filesystem and Disk' % config { nodeQuerySelector: c.nodeQuerySelector }.dashboardNamePrefix,
        time_from=config.dashboardInterval,
        tags=(config.dashboardTags),
        timezone=config.dashboardTimezone,
        refresh=config.dashboardRefresh,
        graphTooltip='shared_crosshair',
        uid='node-disk'
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
