local g = import 'grafana-builder/grafana.libsonnet';

{
  grafanaDashboards+:: {
    'node-cluster-rsrc-use.json':
      local legendLink = '%s/dashboard/file/node-rsrc-use.json' % $._config.grafana_prefix;

      g.dashboard('USE Method / Cluster')
      .addRow(
        g.row('')
        .addPanel(
          g.textPanel('About', |||
            This dashboard is inspired by Brendan Gregg's [USE Method](http://www.brendangregg.com/usemethod.html) - 
            showing **U**tilisation, **S**aturation and **E**rrors for the various resources for each node in your cluster.
          |||)
        ) + {
          height: '100px',
          showTitle: false,
        },
      )
      .addRow(
        g.row('CPU')
        .addPanel(
          g.panel('CPU Utilisation') +
          g.queryPanel(|||
            (
              instance:node_cpu_utilisation:rate1m{%(nodeExporterSelector)s}
            *
              instance:node_num_cpu:sum{%(nodeExporterSelector)s}
            )
            / scalar(sum(instance:node_num_cpu:sum{%(nodeExporterSelector)s}))
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          // TODO: Is this a useful panel? At least there should be some explanation how load
          // average relates to the "CPU saturation" in the title.
          g.panel('CPU Saturation (load1 per CPU)') +
          g.queryPanel(|||
            instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}
            / scalar(count(instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}))
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          // TODO: Does `max: 1` make sense? The stack can go over 1 in high-load scenarios.
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
      )
      .addRow(
        g.row('Memory')
        .addPanel(
          g.panel('Memory Utilisation') +
          g.queryPanel(|||
            instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}
            / scalar(count(instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}))
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          g.panel('Memory Saturation (Major Page Faults)') +
          g.queryPanel('instance:node_vmstat_pgmajfault:rate1m{%(nodeExporterSelector)s}' % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes('rps') },
        )
      )
      .addRow(
        g.row('Network')
        .addPanel(
          g.panel('Net Utilisation (Bytes Receive/Transmit)') +
          g.queryPanel(
            [
              'instance:node_network_receive_bytes_excluding_lo:rate1m{%(nodeExporterSelector)s}' % $._config,
              'instance:node_network_transmit_bytes_excluding_lo:rate1m{%(nodeExporterSelector)s}' % $._config,
            ],
            ['{{instance}} Receive', '{{instance}} Transmit'],
            legendLink,
          ) +
          g.stack +
          {
            yaxes: g.yaxes({ format: 'Bps', min: null }),
            seriesOverrides: [
              {
                alias: '/ Receive/',
                stack: 'A',
              },
              {
                alias: '/ Transmit/',
                stack: 'B',
                transform: 'negative-Y',
              },
            ],
          },
        )
        .addPanel(
          g.panel('Net Saturation (Drops Receive/Transmit)') +
          g.queryPanel(
            [
              'instance:node_network_receive_drop_excluding_lo:rate1m{%(nodeExporterSelector)s}' % $._config,
              'instance:node_network_transmit_drop_excluding_lo:rate1m{%(nodeExporterSelector)s}' % $._config,
            ],
            ['{{instance}} Receive', '{{instance}} Transmit'],
            legendLink,
          ) +
          g.stack +
          {
            yaxes: g.yaxes({ format: 'rps', min: null }),
            seriesOverrides: [
              {
                alias: '/ Receive/',
                stack: 'A',
              },
              {
                alias: '/ Transmit/',
                stack: 'B',
                transform: 'negative-Y',
              },
            ],
          },
        )
      )
      .addRow(
        g.row('Disk IO')
        .addPanel(
          g.panel('Disk IO Utilisation') +
          // Full utilisation would be all disks on each node spending an average of
          // 1 second per second doing I/O, normalize by metric cardinality for stacked charts.
          // TODO: Does the partition by device make sense? Using the most utilized device per
          // instance might make more sense.
          g.queryPanel(|||
            instance_device:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s}
            / scalar(count(instance_device:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s}))
          ||| % $._config, '{{instance}} {{device}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          g.panel('Disk IO Saturation') +
          g.queryPanel(|||
            instance_device:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s}
            / scalar(count(instance_device:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s}))
          ||| % $._config, '{{instance}} {{device}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
      )
      .addRow(
        g.row('Disk Space')
        .addPanel(
          g.panel('Disk Space Utilisation') +
          g.queryPanel(|||
            sum without (device) (
              max without (fstype, mountpoint) (
                node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s} - node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s}
              )
            ) 
            / scalar(sum(max without (fstype, mountpoint) (node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s})))
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        ),
      ),

    'node-rsrc-use.json':
      g.dashboard('USE Method / Node')
      .addTemplate('instance', 'up{%(nodeExporterSelector)s}' % $._config, 'instance')
      .addRow(
        g.row('')
        .addPanel(
          g.textPanel('About', |||
            This dashboard is inspired by Brendan Gregg's [USE Method](http://www.brendangregg.com/usemethod.html) - 
            showing **U**tilisation, **S**aturation and **E**rrors for the various resources for a specific node.
          |||)
        ) + {
          height: '100px',
          showTitle: false,
        },
      )
      .addRow(
        (g.row('Highlights') +
         {
           height: '100px',
           showTitle: false,
         })
        .addPanel(
          g.panel('CPU Cores') +
          g.statPanel('count(node_cpu_seconds_total{%(nodeExporterSelector)s, instance="$instance"})' % $._config, format='short')
        )
        .addPanel(
          g.panel('Total Memory') +
          g.statPanel('node_memory_MemTotal_bytes{%(nodeExporterSelector)s, instance="$instance"}' % $._config, format='decbytes')
        )
        .addPanel(
          g.panel('Disks') +
          g.statPanel('count(node_disk_io_now{%(nodeExporterSelector)s, instance="$instance"})' % $._config, format='short')
        )
        .addPanel(
          g.panel('Disks Space') +
          g.statPanel('sum(node_filesystem_size_bytes{%(nodeExporterSelector)s, instance="$instance"})' % $._config, format='decbytes')
        )
        .addPanel(
          g.panel('Network Interfaces') +
          g.statPanel('count(node_network_device_id{%(nodeExporterSelector)s, instance="$instance", device!~"(veth.+)|(docker[0-9]+)|(cbr[0-9]+)|lo"})' % $._config, format='short')
        )
        .addPanel(
          g.panel('Uptime') +
          g.statPanel('time() - node_boot_time_seconds{%(nodeExporterSelector)s, instance="$instance"}' % $._config, format='dtdhms')
        )
      )
      .addRow(
        g.row('CPU')
        .addPanel(
          g.panel('CPU Utilisation') +
          g.queryPanel('instance:node_cpu_utilisation:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Utilisation') +
          {
            yaxes: g.yaxes('percentunit'),
            legend+: { show: false },
          },
        )
        .addPanel(
          // TODO: Is this a useful panel? At least there should be some explanation how load
          // average relates to the "CPU saturation" in the title.
          g.panel('CPU Saturation (Load1 per CPU)') +
          g.queryPanel('instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Saturation') +
          {
            yaxes: g.yaxes('percentunit'),
            legend+: { show: false },
          },
        )
      )
      .addRow(
        g.row('Memory')
        .addPanel(
          g.panel('Memory Utilisation') +
          g.queryPanel('instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s, %(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Memory') +
          { yaxes: g.yaxes('percentunit') },
        )
        .addPanel(
          g.panel('Memory Saturation (Major Page Faults)') +
          g.queryPanel('instance:node_vmstat_pgmajfault:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Major page faults') +
          {
            yaxes: g.yaxes('short'),
            legend+: { show: false },
          },
        )
      )
      .addRow(
        g.row('Net')
        .addPanel(
          g.panel('Net Utilisation (Bytes Receive/Transmit)') +
          g.queryPanel(
            [
              'instance:node_network_receive_bytes_excluding_lo:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config,
              'instance:node_network_transmit_bytes_excluding_lo:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config,
            ],
            ['Receive', 'Transmit'],
          ) +
          {
            yaxes: g.yaxes({ format: 'Bps', min: null }),
            seriesOverrides: [
              {
                alias: '/Receive/',
                stack: 'A',
              },
              {
                alias: '/Transmit/',
                stack: 'B',
                transform: 'negative-Y',
              },
            ],
          },
        )
        .addPanel(
          g.panel('Net Saturation (Drops Receive/Transmit)') +
          g.queryPanel(
            [
              'instance:node_network_receive_drop_excluding_lo:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config,
              'instance:node_network_transmit_drop_excluding_lo:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config,
            ],
            ['Receive drops', 'Transmit drops'],
          ) +
          {
            yaxes: g.yaxes({ format: 'rps', min: null }),
            seriesOverrides: [
              {
                alias: '/Receive/',
                stack: 'A',
              },
              {
                alias: '/Transmit/',
                stack: 'B',
                transform: 'negative-Y',
              },
            ],
          },
        )
      )
      .addRow(
        g.row('Disk IO')
        .addPanel(
          g.panel('Disk IO Utilisation') +
          g.queryPanel('instance_device:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, '{{device}}') +
          { yaxes: g.yaxes('percentunit') },
        )
        .addPanel(
          g.panel('Disk IO Saturation') +
          g.queryPanel('instance_device:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, '{{device}}') +
          { yaxes: g.yaxes('percentunit') },
        )
      )
      .addRow(
        g.row('Disk Space')
        .addPanel(
          g.panel('Disk Space Utilisation') +
          g.queryPanel(|||
            1 -
            (
              max without (mountpoint, fstype) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s, instance="$instance"})
            /
              max without (mountpoint, fstype) (node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s, instance="$instance"})
            )
          ||| % $._config, '{{device}}') +
          {
            yaxes: g.yaxes('percentunit'),
            legend+: { show: false },
          },
        ),
      ),
  },
}
