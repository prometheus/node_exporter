{
  prometheusAlerts+:: {
    groups+: [
      {
        name: 'node-exporter',
        rules: [
          {
            alert: 'NodeFilesystemSpaceFillingUp',
            expr: |||
              predict_linear(node_filesystem_avail{%(nodeExporterSelector)s,%(fsSelectors)s}[6h], 24*60*60) < 0
              and
              node_filesystem_avail{%(nodeExporterSelector)s,%(fsSelectors)s} / node_filesystem_size{%(nodeExporterSelector)s,%(fsSelectors)s} < 0.4
              and
              node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelectors)s} == 0
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
              predict_linear(node_filesystem_avail{%(nodeExporterSelector)s,%(fsSelectors)s}[6h], 4*60*60) < 0
              and
              node_filesystem_avail{%(nodeExporterSelector)s,%(fsSelectors)s} / node_filesystem_size{%(nodeExporterSelector)s,%(fsSelectors)s} < 0.2
              and
              node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelectors)s} == 0
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
            alert: 'NodeFilesystemOutOfSpace',
            expr: |||
              node_filesystem_avail{%(nodeExporterSelector)s,%(fsSelectors)s} / node_filesystem_size{%(nodeExporterSelector)s,%(fsSelectors)s} * 100 < 5
              and
              node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelectors)s} == 0
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
            alert: 'NodeFilesystemOutOfSpace',
            expr: |||
              node_filesystem_avail{%(nodeExporterSelector)s,%(fsSelectors)s} / node_filesystem_size{%(nodeExporterSelector)s,%(fsSelectors)s} * 100 < 3
              and
              node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelectors)s} == 0
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
              predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelectors)s}[6h], 24*60*60) < 0
              and
              node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelectors)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelectors)s} < 0.4
              and
              node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelectors)s} == 0
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
              predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelectors)s}[6h], 4*60*60) < 0
              and
              node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelectors)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelectors)s} < 0.2
              and
              node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelectors)s} == 0
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              message: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} is predicted to run out of files within the next 4 hours.',
            },
          },
          {
            alert: 'NodeFilesystemOutOfFiles',
            expr: |||
              node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelectors)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelectors)s} * 100 < 5
              and
              node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelectors)s} == 0
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
            alert: 'NodeFilesystemOutOfSpace',
            expr: |||
              node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelectors)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelectors)s} * 100 < 3
              and
              node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelectors)s} == 0
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
              increase(node_network_receive_errs[2m]) > 10
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'critical',
            },
            annotations: {
              message: '{{ $labels.instance }} interface {{ $labels.device }} shows errors while receiving packets ({{ $value }} errors in two minutes).',
            },
          },
          {
            alert: 'NodeNetworkTransmitErrs',
            expr: |||
              increase(node_network_transmit_errs[2m]) > 10
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'critical',
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
