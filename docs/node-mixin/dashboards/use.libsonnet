local g = import 'grafana-builder/grafana.libsonnet';

{
  grafanaDashboards+:: {
    'node-cluster-rsrc-use.json':
      local legendLink = '%s/dashboard/file/k8s-node-rsrc-use.json' % $._config.grafana_prefix;

      g.dashboard('USE Method / Cluster')
      .addRow(
        g.row('CPU')
        .addPanel(
          g.panel('CPU Utilisation') +
          g.queryPanel(|||
            (
              instance:node_cpu_utilisation:rate1m{%(nodeExporterSelector)s}
            *
              instance:node_num_cpu:sum{%(nodeExporterSelector)s}
            / ignoring (instance) group_left
              sum without (instance) (instance:node_num_cpu:sum{%(nodeExporterSelector)s})
            )
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          // TODO: Is this a useful panel?
          g.panel('CPU Saturation (load1 per CPU)') +
          g.queryPanel(|||
            (
              instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s}
            / ignoring (instance) group_left
              count without (instance) (instance:node_load1_per_cpu:ratio{%(nodeExporterSelector)s})
            )
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
          g.queryPanel('instance:node_memory_utilisation:ratio{%(nodeExporterSelector)s}' % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          g.panel('Memory Saturation (Swapped Pages)') +
          g.queryPanel('instance:node_memory_swap_io_pages:rate{%(nodeExporterSelector)s}' % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes('rps') },
        )
      )
      .addRow(
        g.row('Disk')
        .addPanel(
          g.panel('Disk IO Utilisation') +
          // Full utilisation would be all disks on each node spending an average of
          // 1 second per second doing I/O, normalize by metric cardinality for stacked charts.
          g.queryPanel(|||
            (
              instance:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s}
            / ignoring (instance) group_left
              count without (instance) (instance:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s})
            )
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          g.panel('Disk IO Saturation') +
          g.queryPanel(|||
            (
              instance:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s}
            / ignoring (instance) group_left
              count without (instance) (instance:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s})
            )
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
      )
      .addRow(
        g.row('Network')
        .addPanel(
          g.panel('Net Utilisation (Bytes Receive/Transmit)') +
          g.queryPanel(
            [
              'instance:node_network_receive_bytes:rate1m{%(nodeExporterSelector)s}' % $._config,
              '-instance:node_network_transmit_bytes:rate1m{%(nodeExporterSelector)s}' % $._config,
            ],
            ['{{instance}} Receive', '{{instance}} Transmit'],
            legendLink,
          ) +
          g.stack +
          { yaxes: g.yaxes('Bps') },
        )
        .addPanel(
          g.panel('Net Saturation (Drops Receive/Transmit)') +
          g.queryPanel(
            [
              'instance:node_network_receive_drop:rate1m{%(nodeExporterSelector)s}' % $._config,
              '-instance:node_network_transmit_drop:rate1m{%(nodeExporterSelector)s}' % $._config,
            ],
            ['{{instance}} Receive', '{{instance}} Transmit'],
            legendLink,
          ) +
          g.stack +
          { yaxes: g.yaxes('rps') },
        )
      )
      .addRow(
        g.row('Storage')
        .addPanel(
          g.panel('Disk Capacity') +
          g.queryPanel(|||
            (
              sum without (device) (
                max without (fstype, mountpoint) (
                  node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s} - node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s}
                )
              ) 
            / ignoring (instance) group_left
              sum without (instance, device) (
                max without (fstype, mountpoint) (
                  node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s}
                )
              )
            )  
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        ),
      ),

    'node-rsrc-use.json':
      g.dashboard('USE Method / Node')
      .addTemplate('instance', 'up{%(nodeExporterSelector)s}' % $._config, 'instance')
      .addRow(
        g.row('CPU')
        .addPanel(
          g.panel('CPU Utilisation') +
          g.queryPanel('instance:node_cpu_utilisation:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Utilisation') +
          { yaxes: g.yaxes('percentunit') },
        )
        .addPanel(
          g.panel('CPU Saturation (Load1)') +
          g.queryPanel('instance:node_cpu_saturation_load1:{%(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Saturation') +
          { yaxes: g.yaxes('percentunit') },
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
          g.panel('Memory Saturation (pages swapped per second)') +
          g.queryPanel('instance:node_memory_swap_io_pages:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Swap IO') +
          { yaxes: g.yaxes('short') },
        )
      )
      .addRow(
        g.row('Disk')
        .addPanel(
          g.panel('Disk IO Utilisation') +
          g.queryPanel('instance:node_disk_io_time_seconds:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Utilisation') +
          { yaxes: g.yaxes('percentunit') },
        )
        .addPanel(
          g.panel('Disk IO Saturation') +
          g.queryPanel('instance:node_disk_io_time_weighted_seconds:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config, 'Saturation') +
          { yaxes: g.yaxes('percentunit') },
        )
      )
      .addRow(
        g.row('Net')
        .addPanel(
          g.panel('Net Utilisation (Bytes Receive/Transmit)') +
          g.queryPanel(
            [
              'instance:node_network_receive_bytes:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config,
              '-instance:node_network_transmit_bytes:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config,
            ],
            ['Receive', 'Transmit'],
          ) +
          { yaxes: g.yaxes('Bps') },
        )
        .addPanel(
          g.panel('Net Saturation (Drops Receive/Transmit)') +
          g.queryPanel(
            [
              'instance:node_network_receive_drop:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config,
              '-instance:node_network_transmit_drop:rate1m{%(nodeExporterSelector)s, instance="$instance"}' % $._config,
            ],
            ['Receive drops', 'Transmit drops'],
          ) +
          { yaxes: g.yaxes('rps') },
        )
      )
      .addRow(
        g.row('Disk')
        .addPanel(
          g.panel('Disk Utilisation') +
          g.queryPanel(|||
            1 -
            (
              sum(max without (mountpoint, fstype) (node_filesystem_avail_bytes{%(nodeExporterSelector)s, %(fsSelector)s}))
            /
              sum(max without (mountpoint, fstype) (node_filesystem_size_bytes{%(nodeExporterSelector)s, %(fsSelector)s}))
            )
          ||| % $._config, 'Disk') +
          { yaxes: g.yaxes('percentunit') },
        ),
      ),
  },
}
