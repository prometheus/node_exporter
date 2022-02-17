local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local promgrafonnet = import 'github.com/kubernetes-monitoring/kubernetes-mixin/lib/promgrafonnet/promgrafonnet.libsonnet';
local gauge = promgrafonnet.gauge;
local loki = grafana.loki;
local logPanel = grafana.logPanel;

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
              (1 - sum without (mode) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode=~"idle|iowait|steal", instance="$instance"}[$__rate_interval])))
            / ignoring(cpu) group_left
              count without (cpu, mode) (node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle", instance="$instance"})
            )
          ||| % $._config,
          legendFormat='{{cpu}}',
          intervalFactor=5,
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
          ||| % $._config,
          legendFormat='memory used'
        ))
        .addTarget(prometheus.target('node_memory_Buffers_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory buffers'))
        .addTarget(prometheus.target('node_memory_Cached_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory cached'))
        .addTarget(prometheus.target('node_memory_MemFree_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, legendFormat='memory free'));

      // TODO: It would be nicer to have a gauge that gets a 0-1 range and displays it as a percentage 0%-100%.
      // This needs to be added upstream in the promgrafonnet library and then changed here.
      // NOTE: avg() is used to circumvent a label change caused by a node_exporter rollout.
      local memoryGauge = gauge.new(
        'Memory Usage',
        |||
          100 -
          (
            avg(node_memory_MemAvailable_bytes{%(nodeExporterSelector)s, instance="$instance"})
          /
            avg(node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"})
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
          'rate(node_disk_read_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % $._config,
          legendFormat='{{device}} read',
        ))
        .addTarget(prometheus.target(
          'rate(node_disk_written_bytes_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % $._config,
          legendFormat='{{device}} written',
        ))
        .addTarget(prometheus.target(
          'rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, instance="$instance", %(diskDeviceSelector)s}[$__rate_interval])' % $._config,
          legendFormat='{{device}} io time',
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
          'rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[$__rate_interval])' % $._config,
          legendFormat='{{device}}',
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
          'rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, instance="$instance", device!="lo"}[$__rate_interval])' % $._config,
          legendFormat='{{device}}',
        ));

      local syslog = 
        logPanel.new(
          'syslog Errors',
          datasource='$logs_datasource',
        )
        .addTarget(
          loki.target('{filename=~"/var/log/syslog*|/var/log/messages*", %(nodeExporterSelector)s, instance=~"$instance"} |~".+(?i)error(?-i).+"' % $._config)
        );

      local authlog = 
        logPanel.new(
          'authlog',
          datasource='$logs_datasource',
        )
        .addTarget(
          loki.target('{filename=~"/var/log/auth.log|/var/log/secure", %(nodeExporterSelector)s, instance=~"$instance"}' % $._config)
        );

      local kernellog = 
        logPanel.new(
          'Kernel logs',
          datasource='$logs_datasource',
        )
        .addTarget(
          loki.target('{filename=~"/var/log/kern.log*", %(nodeExporterSelector)s, instance=~"$instance"}' % $._config)
        );
        
      local journalsyslog = 
        logPanel.new(
          'Journal syslogs',
          datasource='$logs_datasource',
        )
        .addTarget(
          loki.target('{transport="syslog", %(nodeExporterSelector)s, instance=~"$instance"}' % $._config)
        );
        
      local journalkernel = 
        logPanel.new(
          'Journal Kernel logs',
          datasource='$logs_datasource',
        )
        .addTarget(
          loki.target('{transport="kernel", %(nodeExporterSelector)s, instance=~"$instance"}' % $._config)
        );
        
      local journalstdout = 
        logPanel.new(
          'Journal stdout Errors',
          datasource='$logs_datasource',
        )
        .addTarget(
          loki.target('{transport="stdout", %(nodeExporterSelector)s, instance=~"$instance", unit=~"$unit"} |~".+(?i)error(?-i).+"' % $._config)
        );
        
      dashboard.new(
        '%sNodes' % $._config.dashboardNamePrefix,
        time_from='now-1h',
        tags=($._config.dashboardTags),
        timezone='utc',
        refresh='30s',
        graphTooltip='shared_crosshair'
      )      
      .addTemplate(
        {
          current: {
            text: 'default',
            value: 'default',
          },
          hide: 0,
          label: 'Data Source',
          name: 'datasource',
          options: [],
          query: 'prometheus',
          refresh: 1,
          regex: '',
          type: 'datasource',
        },
      )
      .addTemplate(
        {
          current: {
            text: 'Loki',
            value: 'Loki',
          },
          hide: 0,
          label: 'Loki Data Source',
          name: 'logs_datasource',
          options: [],
          query: 'loki',
          refresh: 1,
          regex: '',
          type: 'datasource',
        },
      )
      .addTemplate(
        template.new(
          'job',
          '$datasource',
          'label_values(node_exporter_build_info, job)',
          refresh='time',
        )
      )
      .addTemplate(
        template.new(
          'instance',
          '$datasource',
          'label_values(node_exporter_build_info{%(nodeExporterSelector)s}, instance)',
          refresh='time',
        )
      )
      .addTemplate(
        template.new(
          'unit',
          '$logs_datasource',
          'label_values(unit)',
          refresh='time',
          includeAll=true,
          multi=true,
          allValues='.+',
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
      )
      .addRow(
        row.new(
          'Direct Log Scrapes'
        )
        .addPanel(syslog)
        .addPanel(authlog)
        .addPanel(kernellog)
      )
      .addRow(
        row.new(
          'Journal Log Scrapes'
        )
        .addPanel(journalsyslog)
        .addPanel(journalkernel)
        .addPanel(journalstdout)
      ),
  },  
}
