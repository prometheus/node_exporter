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
              instance:node_cpu_utilisation:avg1m
            *
              instance:node_num_cpu:sum
            / ignoring (instance) group_left
              sum without (instance) (instance:node_num_cpu:sum)
            )
          |||, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          // TODO: Is this a useful panel?
          g.panel('CPU Saturation (load1 per CPU)') +
          g.queryPanel(|||
            (
              instance:node_load1_per_cpu:ratio
            / ignoring (instance) group_left
              count without (instance) (instance:node_load1_per_cpu:ratio)
            )
          |||, '{{instance}}', legendLink) +
          g.stack +
          // TODO: Does `max: 1` make sense? The stack can go over 1 in high-load scenarios.
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
      )
      .addRow(
        g.row('Memory')
        .addPanel(
          g.panel('Memory Utilisation') +
          g.queryPanel('instance:node_memory_utilisation:ratio', '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          g.panel('Memory Saturation (Swap I/O)') +
          g.queryPanel('instance:node_memory_swap_io_bytes:sum_rate', '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes('Bps') },
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
              instance:node_disk_utilisation:sum_irate
            / ignoring (instance) group_left
              count without (instance) (instance:node_disk_utilisation:sum_irate)
            )
          |||, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          g.panel('Disk IO Saturation') +
          g.queryPanel(|||
            (
              instance:node_disk_saturation:sum_irate
            / ignoring (instance) group_left
              count without (instance) (instance:node_disk_saturation:sum_irate)
            )
          |||, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
      )
      .addRow(
        g.row('Network')
        .addPanel(
          g.panel('Net Utilisation (Transmitted)') +
          g.queryPanel('instance:node_net_utilisation:sum_irate', '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes('Bps') },
        )
        .addPanel(
          g.panel('Net Saturation (Dropped)') +
          g.queryPanel('instance:node_net_saturation:sum_irate', '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes('Bps') },
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
                  node_filesystem_size_bytes{fstype=~"ext[24]"} - node_filesystem_avail_bytes{fstype=~"ext[24]"}
                )
              ) 
            / ignoring (instance) group_left
              sum without (instance, device) (
                max without (fstype, mountpoint) (
                  node_filesystem_size_bytes{fstype=~"ext[24]"}
                )
              )
            )  
          |||, '{{instance}}', legendLink) +
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
          g.queryPanel('instance:node_cpu_utilisation:avg1m{instance="$instance"}', 'Utilisation') +
          { yaxes: g.yaxes('percentunit') },
        )
        .addPanel(
          g.panel('CPU Saturation (Load1)') +
          g.queryPanel('instance:node_cpu_saturation_load1:{instance="$instance"}', 'Saturation') +
          { yaxes: g.yaxes('percentunit') },
        )
      )
      .addRow(
        g.row('Memory')
        .addPanel(
          g.panel('Memory Utilisation') +
          g.queryPanel('instance:node_memory_utilisation:ratio{instance="$instance"}', 'Memory') +
          { yaxes: g.yaxes('percentunit') },
        )
        .addPanel(
          g.panel('Memory Saturation (pages swapped per second)') +
          g.queryPanel('instance:node_memory_swap_io_pages:sum_rate{instance="$instance"}', 'Swap IO') +
          { yaxes: g.yaxes('short') },
        )
      )
      .addRow(
        g.row('Disk')
        .addPanel(
          g.panel('Disk IO Utilisation') +
          g.queryPanel('instance:node_disk_utilisation:sum_irate{instance="$instance"}', 'Utilisation') +
          { yaxes: g.yaxes('percentunit') },
        )
        .addPanel(
          g.panel('Disk IO Saturation') +
          g.queryPanel('instance:node_disk_saturation:sum_irate{instance="$instance"}', 'Saturation') +
          { yaxes: g.yaxes('percentunit') },
        )
      )
      .addRow(
        g.row('Net')
        .addPanel(
          g.panel('Net Utilisation (Transmitted)') +
          g.queryPanel('instance:node_net_utilisation:sum_irate{instance="$instance"}', 'Utilisation') +
          { yaxes: g.yaxes('Bps') },
        )
        .addPanel(
          g.panel('Net Saturation (Dropped)') +
          g.queryPanel('instance:node_net_saturation:sum_irate{instance="$instance"}', 'Saturation') +
          { yaxes: g.yaxes('Bps') },
        )
      )
      .addRow(
        g.row('Disk')
        .addPanel(
          g.panel('Disk Utilisation') +
          g.queryPanel(|||
            1 -
            (
              sum(max without (mountpoint, fstype) (node_filesystem_avail_bytes{fstype=~"ext[24]"}))
            /
              sum(max without (mountpoint, fstype) (node_filesystem_size_bytes{fstype=~"ext[24]"}))
            )
          |||, 'Disk') +
          { yaxes: g.yaxes('percentunit') },
        ),
      ),
  },
}
