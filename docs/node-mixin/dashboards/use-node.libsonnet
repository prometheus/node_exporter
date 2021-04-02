local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;

{
  grafanaDashboards+:: {
    'node-rsrc-use.json':

      local CPUUtilisation =
        graphPanel.new(
          'CPU Utilisation',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10,
          legend_show=false,
        )
        .addTarget(prometheus.target('instance:node_cpu_utilisation:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='Utilisation'));

      local CPUSaturation =
        // TODO: Is this a useful panel? At least there should be some explanation how load
        // average relates to the "CPU saturation" in the title.
        graphPanel.new(
          'CPU Saturation (Load1 per CPU)',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10,
          legend_show=false,
        )
        .addTarget(prometheus.target('instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='Saturation'));


      local memoryUtilisation =
        graphPanel.new(
          'Memory Utilisation',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10,
          legend_show=false,
        )
        .addTarget(prometheus.target('instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='Utilisation'));

      local memorySaturation =
        graphPanel.new(
          'Memory Saturation (Major Page Faults)',
          datasource='$datasource',
          span=6,
          format='rds',
          stack=true,
          fill=10,
          legend_show=false,
        )
        .addTarget(prometheus.target('instance:node_vmstat_pgmajfault:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='Major page Faults'));

      local networkUtilisation =
        graphPanel.new(
          'Network Utilisation (Bytes Receive/Transmit)',
          datasource='$datasource',
          span=6,
          format='Bps',
          stack=true,
          fill=10,
          legend_show=false,
        )
        .addSeriesOverride({ alias: '/Receive/', stack: 'A' })
        .addSeriesOverride({ alias: '/Transmit/', stack: 'B', transform: 'negative-Y' })
        .addTarget(prometheus.target('instance:node_network_receive_bytes_excluding_lo:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='Receive'))
        .addTarget(prometheus.target('instance:node_network_transmit_bytes_excluding_lo:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='Transmit'));

      local networkSaturation =
        graphPanel.new(
          'Network Saturation (Drops Receive/Transmit)',
          datasource='$datasource',
          span=6,
          format='Bps',
          stack=true,
          fill=10,
          legend_show=false,
        )
        .addSeriesOverride({ alias: '/ Receive/', stack: 'A' })
        .addSeriesOverride({ alias: '/ Transmit/', stack: 'B', transform: 'negative-Y' })
        .addTarget(prometheus.target('instance:node_network_receive_drop_excluding_lo:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='Receive'))
        .addTarget(prometheus.target('instance:node_network_transmit_drop_excluding_lo:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='Transmit'));

      local diskIOUtilisation =
        graphPanel.new(
          'Disk IO Utilisation',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10,
          legend_show=false,
        )
        .addTarget(prometheus.target('instance_device:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='{{device}}'));

      local diskIOSaturation =
        graphPanel.new(
          'Disk IO Saturation',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10,
          legend_show=false,
        )
        .addTarget(prometheus.target('instance_device:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='{{device}}'));

      local diskSpaceUtilisation =
        graphPanel.new(
          'Disk Space Utilisation',
          datasource='$datasource',
          span=12,
          format='percentunit',
          stack=false,
          fill=5,
          legend_show=true,
          legend_alignAsTable=true,
          legend_current=true,
          legend_avg=true,
          legend_rightSide=true,
          legend_sortDesc=true,
        )
        .addTarget(prometheus.target(
          |||
            sort_desc(1 -
              (
               max without (mountpoint, fstype) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, fstype!="", instance="$instance"})
               /
               max without (mountpoint, fstype) (node_filesystem_size_bytes{%(nodeExporterSelector)s, fstype!="", instance="$instance"})
              )
            )
          ||| % $._config, legendFormat='{{device}}'
        ));

      dashboard.new(
        '%sUSE Method / Node' % $._config.dashboardNamePrefix,
        time_from='now-1h',
        tags=($._config.dashboardTags),
        timezone='utc',
        refresh='30s',
        graphTooltip='shared_crosshair'
      )
      .addTemplate(
        {
          current: {
            text: 'Prometheus',
            value: 'Prometheus',
          },
          hide: 0,
          label: null,
          name: 'datasource',
          options: [],
          query: 'prometheus',
          refresh: 1,
          regex: '',
          type: 'datasource',
        },
      )
      .addTemplate(
        template.new(
          'instance',
          '$datasource',
          'label_values(node_exporter_build_info{%(nodeExporterSelector)s}, instance)' % $._config,
          refresh='time',
        )
      )
      .addRow(
        row.new('CPU')
        .addPanel(CPUUtilisation)
        .addPanel(CPUSaturation)
      )
      .addRow(
        row.new('Memory')
        .addPanel(memoryUtilisation)
        .addPanel(memorySaturation)
      )
      .addRow(
        row.new('Network')
        .addPanel(networkUtilisation)
        .addPanel(networkSaturation)
      )
      .addRow(
        row.new('Disk IO')
        .addPanel(diskIOUtilisation)
        .addPanel(diskIOSaturation)
      )
      .addRow(
        row.new('Disk Space')
        .addPanel(diskSpaceUtilisation)
      ),
  },
}
