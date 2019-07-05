local grafana = import 'grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local promgrafonnet = import '../lib/promgrafonnet/promgrafonnet.libsonnet';
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
          |||
            1 - avg by (cpu) (irate(node_cpu{%(nodeExporterSelector)s, mode="idle", instance="$instance"}[1m]))
          ||| % $._config,
          legendFormat='{{cpu}}',
          intervalFactor=10,
        ));

      local systemLoad =
        graphPanel.new(
          'System load',
          datasource='$datasource',
          span=6,
          format='percentunit',
        )
        .addTarget(prometheus.target('node_load1{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='load 1m'))
        .addTarget(prometheus.target('node_load5{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='load 5m'))
        .addTarget(prometheus.target('node_load15{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='load 15m'));

      local memoryGraph =
        graphPanel.new(
          'Memory Usage',
          datasource='$datasource',
          span=9,
          format='bytes',
        )
        .addTarget(prometheus.target(
          |||
            node_memory_MemTotal{%(nodeExporterSelector)s, instance="$instance"}
            - node_memory_MemFree{%(nodeExporterSelector)s, instance="$instance"}
            - node_memory_Buffers{%(nodeExporterSelector)s, instance="$instance"}
            - node_memory_Cached{%(nodeExporterSelector)s, instance="$instance"}
          ||| % $._config, legendFormat='memory used'
        ))
        .addTarget(prometheus.target('node_memory_Buffers{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory buffers'))
        .addTarget(prometheus.target('node_memory_Cached{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory cached'))
        .addTarget(prometheus.target('node_memory_MemFree{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory free'));

      local memoryGauge = gauge.new(
        'Memory Usage',
        |||
          node_memory_MemAvailable{%(nodeExporterSelector)s, instance="$instance"}
          /
          node_memory_MemTotal{%(nodeExporterSelector)s, instance="$instance"}
        ||| % $._config,
      ).withLowerBeingBetter();

      local diskIO =
        graphPanel.new(
          'Disk I/O',
          datasource='$datasource',
          span=9,
        )
        .addTarget(prometheus.target('sum by (instance) (irate(node_disk_bytes_read_total{%(nodeExporterSelector)s, instance="$instance"}[1m]))' % $._config, legendFormat='read'))
        .addTarget(prometheus.target('sum by (instance) (irate(node_disk_bytes_written_total{%(nodeExporterSelector)s, instance="$instance"}[1m]))' % $._config, legendFormat='written'))
        .addTarget(prometheus.target('sum by (instance) (irate(node_disk_io_time_ms{%(nodeExporterSelector)s,  instance="$instance"}[1m]))' % $._config, legendFormat='io time')) +
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
            self.yaxe(format='ms'),
          ],
        };

      local diskSpaceUsage = gauge.new(
        'Disk Space Usage',
        |||
          1 - (
              sum(node_filesystem_free{%(nodeExporterSelector)s, device!="rootfs", instance="$instance"}
            /
              sum(node_filesystem_size{%(nodeExporterSelector)s, device!="rootfs", instance="$instance"}
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
        .addTarget(prometheus.target('irate(node_network_receive_bytes{%(nodeExporterSelector)s, instance="$instance", device!~"lo"}[1m])' % $._config, legendFormat='{{device}}'));

      local networkTransmitted =
        graphPanel.new(
          'Network Transmitted',
          datasource='$datasource',
          span=6,
          format='bytes',
        )
        .addTarget(prometheus.target('irate(node_network_transmit_bytes{%(nodeExporterSelector)s, instance="$instance", device!~"lo"}[1m])' % $._config, legendFormat='{{device}}'));

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
          'label_values(node_boot_time{%(nodeExporterSelector)s}, instance)' % $._config,
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
