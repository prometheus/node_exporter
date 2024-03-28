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
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceFillingUpWarningThreshold)d
              and
                predict_linear(node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], 24*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
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
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceFillingUpCriticalThreshold)d
              and
                predict_linear(node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], 4*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
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
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceAvailableWarningThreshold)d
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % $._config,
            'for': '30m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem has less than %(fsSpaceAvailableWarningThreshold)d%% space left.' % $._config,
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfSpace',
            expr: |||
              (
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceAvailableCriticalThreshold)d
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % $._config,
            'for': '30m',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % $._config,
            },
            annotations: {
              summary: 'Filesystem has less than %(fsSpaceAvailableCriticalThreshold)d%% space left.' % $._config,
              description: 'Filesystem on {{ $labels.device }} at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 40
              and
                predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], 24*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
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
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 20
              and
                predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], 4*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
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
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 5
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
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
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 3
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
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
              rate(node_network_receive_errs_total[2m]) / rate(node_network_receive_packets_total[2m]) > 0.01
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
              rate(node_network_transmit_errs_total[2m]) / rate(node_network_transmit_packets_total[2m]) > 0.01
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
          {
            alert: 'NodeHighNumberConntrackEntriesUsed',
            expr: |||
              (node_nf_conntrack_entries / node_nf_conntrack_entries_limit) > 0.75
            ||| % $._config,
            annotations: {
              summary: 'Number of conntrack are getting close to the limit.',
              description: '{{ $value | humanizePercentage }} of conntrack entries are used.',
            },
            labels: {
              severity: 'warning',
            },
          },
          {
            alert: 'NodeTextFileCollectorScrapeError',
            expr: |||
              node_textfile_scrape_error{%(nodeExporterSelector)s} == 1
            ||| % $._config,
            annotations: {
              summary: 'Node Exporter text file collector failed to scrape.',
              description: 'Node Exporter text file collector failed to scrape.',
            },
            labels: {
              severity: 'warning',
            },
          },
          {
            alert: 'NodeClockSkewDetected',
            expr: |||
              (
                node_timex_offset_seconds{%(nodeExporterSelector)s} > 0.05
              and
                deriv(node_timex_offset_seconds{%(nodeExporterSelector)s}[5m]) >= 0
              )
              or
              (
                node_timex_offset_seconds{%(nodeExporterSelector)s} < -0.05
              and
                deriv(node_timex_offset_seconds{%(nodeExporterSelector)s}[5m]) <= 0
              )
            ||| % $._config,
            'for': '10m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Clock skew detected.',
              description: 'Clock on {{ $labels.instance }} is out of sync by more than 0.05s. Ensure NTP is configured correctly on this host.',
            },
          },
          {
            alert: 'NodeClockNotSynchronising',
            expr: |||
              min_over_time(node_timex_sync_status{%(nodeExporterSelector)s}[5m]) == 0
              and
              node_timex_maxerror_seconds{%(nodeExporterSelector)s} >= 16
            ||| % $._config,
            'for': '10m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Clock not synchronising.',
              description: 'Clock on {{ $labels.instance }} is not synchronising. Ensure NTP is configured on this host.',
            },
          },
          {
            alert: 'NodeRAIDDegraded',
            expr: |||
              node_md_disks_required{%(nodeExporterSelector)s,%(diskDeviceSelector)s} - ignoring (state) (node_md_disks{state="active",%(nodeExporterSelector)s,%(diskDeviceSelector)s}) > 0
            ||| % $._config,
            'for': '15m',
            labels: {
              severity: 'critical',
            },
            annotations: {
              summary: 'RAID Array is degraded',
              description: "RAID array '{{ $labels.device }}' on {{ $labels.instance }} is in degraded state due to one or more disks failures. Number of spare drives is insufficient to fix issue automatically.",
            },
          },
          {
            alert: 'NodeRAIDDiskFailure',
            expr: |||
              node_md_disks{state="failed",%(nodeExporterSelector)s,%(diskDeviceSelector)s} > 0
            ||| % $._config,
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Failed device in RAID array',
              description: "At least one device in RAID array on {{ $labels.instance }} failed. Array '{{ $labels.device }}' needs attention and possibly a disk swap.",
            },
          },
          {
            alert: 'NodeFileDescriptorLimit',
            expr: |||
              (
                node_filefd_allocated{%(nodeExporterSelector)s} * 100 / node_filefd_maximum{%(nodeExporterSelector)s} > 70
              )
            ||| % $._config,
            'for': '15m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Kernel is predicted to exhaust file descriptors limit soon.',
              description: 'File descriptors limit at {{ $labels.instance }} is currently at {{ printf "%.2f" $value }}%.',
            },
          },
          {
            alert: 'NodeFileDescriptorLimit',
            expr: |||
              (
                node_filefd_allocated{%(nodeExporterSelector)s} * 100 / node_filefd_maximum{%(nodeExporterSelector)s} > 90
              )
            ||| % $._config,
            'for': '15m',
            labels: {
              severity: 'critical',
            },
            annotations: {
              summary: 'Kernel is predicted to exhaust file descriptors limit soon.',
              description: 'File descriptors limit at {{ $labels.instance }} is currently at {{ printf "%.2f" $value }}%.',
            },
          },
        ],
      },
    ],
  },
}
