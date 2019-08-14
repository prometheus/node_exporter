{
  prometheusAlerts+:: {
    groups+: [
      {
        name: 'node-exporter',
        rules: [
          {
            alert: 'NodeFilesystemSpaceFillingUp',
            expr: |||
              (
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s} < 0.4
              and
                predict_linear(node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s}[6h], 24*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of space within the next 24 hours.',
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left and is filling up.',
            },
          },
          {
            alert: 'NodeFilesystemSpaceFillingUp',
            expr: |||
              (
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s} < 0.2
              and
                predict_linear(node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s}[6h], 4*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % $._config,
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of space within the next 4 hours.',
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left and is filling up fast.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfSpace',
            expr: |||
              (
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s} * 100 < 5
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem has less than 5% space left.',
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfSpace',
            expr: |||
              (
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s} * 100 < 3
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % $._config,
            },
            annotations: {
              summary: 'Filesystem has less than 3% space left.',
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s} < 0.4
              and
                predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s}[6h], 24*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of inodes within the next 24 hours.',
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left and is filling up.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s} < 0.2
              and
                predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s}[6h], 4*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % $._config,
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of inodes within the next 4 hours.',
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left and is filling up fast.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfFiles',
            expr: |||
              (
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s} * 100 < 5
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem has less than 5% inodes left.',
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfFiles',
            expr: |||
              (
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s} * 100 < 3
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % $._config,
            },
            annotations: {
              summary: 'Filesystem has less than 3% inodes left.',
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left.',
            },
          },
          {
            alert: 'NodeNetworkReceiveErrs',
            expr: |||
              increase(node_network_receive_errs_total[2m]) > 10
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Network interface is reporting many receive errors.',
              description: '{{ $labels.instance }} interface {{ $labels.device }} has encountered {{ printf "%.0f" $value }} receive errors in the last two minutes.',
            },
          },
          {
            alert: 'NodeNetworkTransmitErrs',
            expr: |||
              increase(node_network_transmit_errs_total[2m]) > 10
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Network interface is reporting many transmit errors.',
              description: '{{ $labels.instance }} interface {{ $labels.device }} has encountered {{ printf "%.0f" $value }} transmit errors in the last two minutes.',
            },
          },
        ],
      },
    ],
  },
}
