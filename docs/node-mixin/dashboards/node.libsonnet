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
          'CPU Usage',
          datasource='$datasource',
          span=6,
          format='percentunit',
          max=1,
          min=0,
          stack=true,
        )
        .addTarget(prometheus.target(
          |||
            (
              (1 - rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance"}[$__interval]))
            / ignoring(cpu) group_left
              count without (cpu)( node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance"})
            )
          ||| % $._config,
          legendFormat='{{cpu}}',
          intervalFactor=5,
          interval='1m',
        ));

      local systemLoad =
        graphPanel.new(
          'Load Average',
          datasource='$datasource',
          span=6,
          format='short',
          min=0,
          fill=0,
        )
        .addTarget(prometheus.target('node_load1{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='1m load average'))
        .addTarget(prometheus.target('node_load5{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='5m load average'))
        .addTarget(prometheus.target('node_load15{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='15m load average'))
        .addTarget(prometheus.target('count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance", mode="idle"})' % $._config, legendFormat='logical cores'));

      local memoryGraph =
        graphPanel.new(
          'Memory Usage',
          datasource='$datasource',
          span=9,
          format='bytes',
          stack=true,
          min=0,
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
          span=6,
          min=0,
          fill=0,
        )
        // TODO: Does it make sense to have those three in the same panel?
        .addTarget(prometheus.target(
          'rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__interval])' % $._config,
          legendFormat='{{device}} read',
          interval='1m',
        ))
        .addTarget(prometheus.target(
          'rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__interval])' % $._config,
          legendFormat='{{device}} written',
          interval='1m',
        ))
        .addTarget(prometheus.target(
          'rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__interval])' % $._config,
          legendFormat='{{device}} io time',
          interval='1m',
        )) +
        {
          seriesOverrides: [
            {
              alias: '/ read| written/',
              yaxis: 1,
            },
            {
              alias: '/ io time/',
              yaxis: 2,
            },
          ],
          yaxes: [
            self.yaxe(format='bytes'),
            self.yaxe(format='s'),
          ],
        };

      // TODO: Somehow partition this by device while excluding read-only devices.
      local diskSpaceUsage =
        graphPanel.new(
          'Disk Space Usage',
          datasource='$datasource',
          span=6,
          format='bytes',
          min=0,
          fill=1,
          stack=true,
        )
        .addTarget(prometheus.target(
          |||
            sum(
              max by (device) (
                node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s}
              -
                node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s}
              )
            )
          ||| % $._config,
          legendFormat='used',
        ))
        .addTarget(prometheus.target(
          |||
            sum(
              max by (device) (
                node_filesystem_avail_bytes{%(nodeExporterSelector)s, instance="$instance", %(fsSelector)s}
              )
            )
          ||| % $._config,
          legendFormat='available',
        )) +
        {
          seriesOverrides: [
            {
              alias: 'used',
              color: '#E0B400',
            },
            {
              alias: 'available',
              color: '#73BF69',
            },
          ],
        };

      local networkReceived =
        graphPanel.new(
          'Network Received',
          datasource='$datasource',
          span=6,
          format='bytes',
          min=0,
          fill=0,
        )
        .addTarget(prometheus.target(
          'rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[$__interval])' % $._config,
          legendFormat='{{device}}',
          interval='1m',
        ));

      local networkTransmitted =
        graphPanel.new(
          'Network Transmitted',
          datasource='$datasource',
          span=6,
          format='bytes',
          min=0,
          fill=0,
        )
        .addTarget(prometheus.target(
          'rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[$__interval])' % $._config,
          legendFormat='{{device}}',
          interval='1m',
        ));

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
