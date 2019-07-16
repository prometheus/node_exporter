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
                predict_linear(node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s}[6h], 24*60*60) < 0
              and
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s} < 0.4
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} is predicted to run out of space within the next 24 hours.',
            },
          },
          {
            alert: 'NodeFilesystemSpaceFillingUp',
            expr: |||
              (
                predict_linear(node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s}[6h], 4*60*60) < 0
              and
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s} < 0.2
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'critical',
            },
            annotations: {
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} is predicted to run out of space within the next 4 hours.',
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
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ $value }}% available space left.',
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
              severity: 'critical',
            },
            annotations: {
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ $value }}% available space left.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s}[6h], 24*60*60) < 0
              and
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s} < 0.4
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} is predicted to run out of files within the next 24 hours.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s}[6h], 4*60*60) < 0
              and
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s} < 0.2
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'critical',
            },
            annotations: {
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} is predicted to run out of files within the next 4 hours.',
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
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ $value }}% available inodes left.',
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
              severity: 'critical',
            },
            annotations: {
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ $value }}% available space left.',
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
              message: '{{ $labels.instance }} interface {{ $labels.device }} shows errors while receiving packets ({{ $value }} errors in two minutes).',
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
              message: '{{ $labels.instance }} interface {{ $labels.device }} shows errors while transmitting packets ({{ $value }} errors in two minutes).',
            },
          },
        ],
      },
    ],
  },
}
