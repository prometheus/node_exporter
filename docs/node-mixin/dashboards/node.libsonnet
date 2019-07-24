local grafana = import 'grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local promgrafonnet = import 'promgrafonnet/promgrafonnet.libsonnet';
local gauge = promgrafonnet.gauge;

{
  grafanaDashboards+:: {
    'nodes.json':
      local idleCPU =
        graphPanel.new(
          'Idle CPU',
          datasource='$datasource',
          span=6,
          format='percentunit',
          max=100,
          min=0,
        )
        .addTarget(prometheus.target(
          // TODO: Consider using `${__interval}` as range and a 1m min step.
          |||
            1 - rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance"}[1m])
          ||| % $._config,
          legendFormat='{{cpu}}',
          intervalFactor=10,
        ));

      // TODO: Is this panel useful?
      local systemLoad =
        graphPanel.new(
          'Load Average',
          datasource='$datasource',
          span=6,
          format='short',
        )
        .addTarget(prometheus.target('node_load1{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='1m load average'))
        .addTarget(prometheus.target('node_load5{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='5m load average'))
        .addTarget(prometheus.target('node_load15{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='15m load average'));

      local memoryGraph =
        graphPanel.new(
          'Memory Usage',
          datasource='$datasource',
          span=9,
          format='bytes',
        )
        .addTarget(prometheus.target(
          |||
            (
              node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"}
            -
              node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance"}
            -
              node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance"}
            -
              node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance"}
            )
          ||| % $._config, legendFormat='memory used'
        ))
        .addTarget(prometheus.target('node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory buffers'))
        .addTarget(prometheus.target('node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory cached'))
        .addTarget(prometheus.target('node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory free'));

      // TODO: It would be nicer to have a gauge that gets a 0-1 range and displays it as a percentage 0%-100%.
      // This needs to be added upstream in the promgrafonnet library and then changed here.
      local memoryGauge = gauge.new(
        'Memory Usage',
        |||
          100 -
          (
            node_memory_MemAvailable_bytes{%(nodeExporterSelector)s, instance="$instance"}
          /
            node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"}
          * 100
          )
        ||| % $._config,
      ).withLowerBeingBetter();

      local diskIO =
        graphPanel.new(
          'Disk I/O',
          datasource='$datasource',
          span=9,
        )
        // TODO: Does it make sense to have those three in the same panel?
        // TODO: Consider using `${__interval}` as range and a 1m min step.
        .addTarget(prometheus.target('rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[1m])' % $._config, legendFormat='{{device}} read'))
        .addTarget(prometheus.target('rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[1m])' % $._config, legendFormat='{{device}} written'))
        .addTarget(prometheus.target('rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[1m])' % $._config, legendFormat='{{device}} io time')) +
        {
          seriesOverrides: [
            {
              alias: 'read',
              yaxis: 1,
            },
            {
              alias: 'io time',
              yaxis: 2,
            },
          ],
          yaxes: [
            self.yaxe(format='bytes'),
            self.yaxe(format='s'),
          ],
        };

      // TODO: It would be nicer to have a gauge that gets a 0-1 range and displays it as a percentage 0%-100%.
      // This needs to be added upstream in the promgrafonnet library and then changed here.
      // TODO: Should this be partitioned by mountpoint?
      local diskSpaceUsage = gauge.new(
        'Disk Space Usage',
        |||
          100 -
          (
            sum(node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s})
          /
            sum(node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s})
          * 100
          )
        ||| % $._config,
      ).withLowerBeingBetter();

      local networkReceived =
        graphPanel.new(
          'Network Received',
          datasource='$datasource',
          span=6,
          format='bytes',
        )
        // TODO: Consider using `${__interval}` as range and a 1m min step.
        .addTarget(prometheus.target('rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[1m])' % $._config, legendFormat='{{device}}'));

      local networkTransmitted =
        graphPanel.new(
          'Network Transmitted',
          datasource='$datasource',
          span=6,
          format='bytes',
        )
        // TODO: Consider using `${__interval}` as range and a 1m min step.
        .addTarget(prometheus.target('rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[1m])' % $._config, legendFormat='{{device}}'));

      dashboard.new('Nodes', time_from='now-1h')
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
        row.new()
        .addPanel(idleCPU)
        .addPanel(systemLoad)
      )
      .addRow(
        row.new()
        .addPanel(memoryGraph)
        .addPanel(memoryGauge)
      )
      .addRow(
        row.new()
        .addPanel(diskIO)
        .addPanel(diskSpaceUsage)
      )
      .addRow(
        row.new()
        .addPanel(networkReceived)
        .addPanel(networkTransmitted)
      ),
  },
}
