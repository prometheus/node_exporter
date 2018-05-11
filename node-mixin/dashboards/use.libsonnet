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
          g.queryPanel('instance:node_cpu_utilisation:avg1m * instance:node_num_cpu:sum / scalar(sum(instance:node_num_cpu:sum))', '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          g.panel('CPU Saturation (Load1)') +
          g.queryPanel(|||
            instance:node_cpu_saturation_load1: / scalar(sum(up{%(nodeExporterSelector)s}))
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
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
          // 1 sec per second doing I/O, normalize by node count for stacked charts
          g.queryPanel(|||
            instance:node_disk_utilisation:sum_irate / scalar(sum(up{%(nodeExporterSelector)s}))
          ||| % $._config, '{{instance}}', legendLink) +
          g.stack +
          { yaxes: g.yaxes({ format: 'percentunit', max: 1 }) },
        )
        .addPanel(
          g.panel('Disk IO Saturation') +
          g.queryPanel(|||
            instance:node_disk_saturation:sum_irate / scalar(sum(up{%(nodeExporterSelector)s}))
          ||| % $._config, '{{instance}}', legendLink) +
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
          g.queryPanel('sum(max(node_filesystem_size{fstype=~"ext[24]"} - node_filesystem_free{fstype=~"ext[24]"}) by (device,instance,namespace)) by (instance,namespace) / scalar(sum(max(node_filesystem_size{fstype=~"ext[24]"}) by (device,instance,namespace)))', '{{instance}}', legendLink) +
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
          g.panel('Memory Saturation (Swap I/O)') +
          g.queryPanel('instance:node_memory_swap_io_bytes:sum_rate{instance="$instance"}', 'Swap IO') +
          { yaxes: g.yaxes('Bps') },
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
          g.queryPanel('1 - sum(max by (device, node) (node_filesystem_free{fstype=~"ext[24]"})) / sum(max by (device, node) (node_filesystem_size{fstype=~"ext[24]"}))', 'Disk') +
          { yaxes: g.yaxes('percentunit') },
        ),
      ),
  },
}
