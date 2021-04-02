local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;

{
  grafanaDashboards+:: {
    'node-cluster-rsrc-use.json':
      local CPUUtilisation =
        graphPanel.new(
          'CPU Utilisation',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10,
        )
        .addTarget(prometheus.target(|||
          (
            instance:node_cpu_utilisation:rate1m{%(nodeExporterSelector)s}
          *
            instance:node_num_cpu:sum{%(nodeExporterSelector)s}
          )
          / scalar(sum(instance:node_num_cpu:sum{%(nodeExporterSelector)s}))
        ||| % $._config, legendFormat='{{ instance }}'));

      local CPUSaturation =
        graphPanel.new(
          'CPU Saturation (load1 per CPU)',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10,
        )
        .addTarget(prometheus.target(
          |||
            instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}
            / scalar(count(instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}))
          ||| % $._config, legendFormat='{{instance}}',
        ));

      local memoryUtilisation =
        graphPanel.new(
          'Memory Utilisation',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10,
        )
        .addTarget(prometheus.target(
          |||
            instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}
            / scalar(count(instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}))
          ||| % $._config, legendFormat='{{instance}}',
        ));

      local memorySaturation =
        graphPanel.new(
          'Memory Saturation (Major Page Faults)',
          datasource='$datasource',
          span=6,
          format='rps',
          stack=true,
          fill=10
        )
        .addTarget(prometheus.target('instance:node_vmstat_pgmajfault:rate1m{%(nodeExporterSelector)s}' % $._config, legendFormat='{{instance}}'));


      local networkUtilisation =
        graphPanel.new(
          'Network Utilisation (Bytes Receive/Transmit)',
          datasource='$datasource',
          span=6,
          format='Bps',
          stack=true,
          fill=10
        )
        .addSeriesOverride({ alias: '/ Receive/', stack: 'A' })
        .addSeriesOverride({ alias: '/ Transmit/', stack: 'B', transform: 'negative-Y' })
        .addTarget(prometheus.target('instance:node_network_receive_bytes_excluding_lo:rate1m{%(nodeExporterSelector)s}' % $._config, legendFormat='{{instance}} Receive'))
        .addTarget(prometheus.target('instance:node_network_transmit_bytes_excluding_lo:rate1m{%(nodeExporterSelector)s}' % $._config, legendFormat='{{instance}} Transmit'));

      local networkSaturation =
        graphPanel.new(
          'Net Saturation (Drops Receive/Transmit)',
          datasource='$datasource',
          span=6,
          format='rps',
          stack=true,
          fill=10
        )
        .addSeriesOverride({ alias: '/ Receive/', stack: 'A' })
        .addSeriesOverride({ alias: '/ Transmit/', stack: 'B', transform: 'negative-Y' })
        .addTarget(prometheus.target('instance:node_network_receive_drop_excluding_lo:rate1m{%(nodeExporterSelector)s}' % $._config, legendFormat='{{instance}} Receive'))
        .addTarget(prometheus.target('instance:node_network_transmit_drop_excluding_lo:rate1m{%(nodeExporterSelector)s}' % $._config, legendFormat='{{instance}} Transmit'));

      local diskIOUtilisation =
        // Full utilisation would be all disks on each node spending an average of
        // 1 second per second doing I/O, normalize by metric cardinality for stacked charts.
        // TODO: Does the partition by device make sense? Using the most utilized device per
        // instance might make more sense.
        graphPanel.new(
          'Disk IO Utilisation',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10
        )
        .addTarget(prometheus.target(
          |||
            instance_device:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s}
            / scalar(count(instance_device:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s}))
          ||| % $._config, legendFormat='{{instance}} {{device}}'
        ));

      local diskIOSaturation =
        graphPanel.new(
          'Disk IO Saturation',
          datasource='$datasource',
          span=6,
          format='percentunit',
          stack=true,
          fill=10
        )
        .addTarget(prometheus.target(
          |||
            instance_device:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s}
            / scalar(count(instance_device:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s}))
          ||| % $._config, legendFormat='{{instance}} {{device}}'
        ));

      local diskSpaceUtilisation =
        graphPanel.new(
          'Disk Space Utilisation',
          datasource='$datasource',
          span=12,
          format='percentunit',
          stack=true,
          fill=10
        )
        .addTarget(prometheus.target(
          |||
            sum without (device) (
              max without (fstype, mountpoint) (
                node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s} - node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s}
              )
            )
            / scalar(sum(max without (fstype, mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s})))
          ||| % $._config, legendFormat='{{instance}}'
        ));

      dashboard.new(
        '%sUSE Method / Cluster' % $._config.dashboardNamePrefix,
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
