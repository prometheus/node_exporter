{
  new(this): {
    groups: [
      {
        name: if this.config.uid == 'node' then 'node-exporter-filesystem' else this.config.uid + '-linux-filesystem-alerts',
        rules: [
          {
            alert: 'NodeFilesystemSpaceFillingUp',
            expr: |||
              (
                node_filesystem_avail_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceFillingUpWarningThreshold)d
              and
                predict_linear(node_filesystem_avail_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], 24*60*60) < 0
              and
                node_filesystem_readonly{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % this.config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of space within the next 24 hours.',
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left and is filling up.',
            },
          },
          {
            alert: 'NodeFilesystemSpaceFillingUp',
            expr: |||
              (
                node_filesystem_avail_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceFillingUpCriticalThreshold)d
              and
                predict_linear(node_filesystem_avail_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], 4*60*60) < 0
              and
                node_filesystem_readonly{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % this.config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % this.config,
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of space within the next 4 hours.',
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left and is filling up fast.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfSpace',
            expr: |||
              (
                node_filesystem_avail_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceAvailableWarningThreshold)d
              and
                node_filesystem_readonly{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % this.config,
            'for': '30m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem has less than %(fsSpaceAvailableWarningThreshold)d%% space left.' % this.config,
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfSpace',
            expr: |||
              (
                node_filesystem_avail_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceAvailableCriticalThreshold)d
              and
                node_filesystem_readonly{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % this.config,
            'for': '30m',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % this.config,
            },
            annotations: {
              summary: 'Filesystem has less than %(fsSpaceAvailableCriticalThreshold)d%% space left.' % this.config,
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                node_filesystem_files_free{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 40
              and
                predict_linear(node_filesystem_files_free{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], 24*60*60) < 0
              and
                node_filesystem_readonly{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % this.config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of inodes within the next 24 hours.',
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left and is filling up.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                node_filesystem_files_free{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 20
              and
                predict_linear(node_filesystem_files_free{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], 4*60*60) < 0
              and
                node_filesystem_readonly{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % this.config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % this.config,
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of inodes within the next 4 hours.',
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left and is filling up fast.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfFiles',
            expr: |||
              (
                node_filesystem_files_free{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 5
              and
                node_filesystem_readonly{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % this.config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem has less than 5% inodes left.',
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left.',
            },
          },
          {
            alert: 'NodeFilesystemAlmostOutOfFiles',
            expr: |||
              (
                node_filesystem_files_free{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 3
              and
                node_filesystem_readonly{%(filteringSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % this.config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % this.config,
            },
            annotations: {
              summary: 'Filesystem has less than 3% inodes left.',
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left.',
            },
          },
        ],
      },
      {
        // defaults to 'node-exporter for backward compatibility with old node-mixin
        name: if this.config.uid == 'node' then 'node-exporter' else this.config.uid + '-linux-alerts',
        rules: [
          {
            alert: 'NodeNetworkReceiveErrs',
            expr: |||
              rate(node_network_receive_errs_total{%(filteringSelector)s}[2m]) / rate(node_network_receive_packets_total{%(filteringSelector)s}[2m]) > 0.01
            ||| % this.config,
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
              rate(node_network_transmit_errs_total{%(filteringSelector)s}[2m]) / rate(node_network_transmit_packets_total{%(filteringSelector)s}[2m]) > 0.01
            ||| % this.config,
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
              (node_nf_conntrack_entries{%(filteringSelector)s} / node_nf_conntrack_entries_limit) > 0.75
            ||| % this.config,
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
              node_textfile_scrape_error{%(filteringSelector)s} == 1
            ||| % this.config,
            annotations: {
              summary: 'Node Exporter text file collector failed to scrape.',
              description: 'Node Exporter text file collector on {{ $labels.instance }} failed to scrape.',
            },
            labels: {
              severity: 'warning',
            },
          },
          {
            alert: 'NodeClockSkewDetected',
            expr: |||
              (
                node_timex_offset_seconds{%(filteringSelector)s} > 0.05
              and
                deriv(node_timex_offset_seconds{%(filteringSelector)s}[5m]) >= 0
              )
              or
              (
                node_timex_offset_seconds{%(filteringSelector)s} < -0.05
              and
                deriv(node_timex_offset_seconds{%(filteringSelector)s}[5m]) <= 0
              )
            ||| % this.config,
            'for': '10m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Clock skew detected.',
              description: 'Clock at {{ $labels.instance }} is out of sync by more than 0.05s. Ensure NTP is configured correctly on this host.',
            },
          },
          {
            alert: 'NodeClockNotSynchronising',
            expr: |||
              min_over_time(node_timex_sync_status{%(filteringSelector)s}[5m]) == 0
              and
              node_timex_maxerror_seconds{%(filteringSelector)s} >= 16
            ||| % this.config,
            'for': '10m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Clock not synchronising.',
              description: 'Clock at {{ $labels.instance }} is not synchronising. Ensure NTP is configured on this host.',
            },
          },
          {
            alert: 'NodeRAIDDegraded',
            expr: |||
              node_md_disks_required{%(filteringSelector)s,%(diskDeviceSelector)s} - ignoring (state) (node_md_disks{state="active",%(filteringSelector)s,%(diskDeviceSelector)s}) > 0
            ||| % this.config,
            'for': '15m',
            labels: {
              severity: 'critical',
            },
            annotations: {
              summary: 'RAID Array is degraded.',
              description: "RAID array '{{ $labels.device }}' at {{ $labels.instance }} is in degraded state due to one or more disks failures. Number of spare drives is insufficient to fix issue automatically.",
            },
          },
          {
            alert: 'NodeRAIDDiskFailure',
            expr: |||
              node_md_disks{state="failed",%(filteringSelector)s,%(diskDeviceSelector)s} > 0
            ||| % this.config,
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Failed device in RAID array.',
              description: "At least one device in RAID array at {{ $labels.instance }} failed. Array '{{ $labels.device }}' needs attention and possibly a disk swap.",
            },
          },
          {
            alert: 'NodeFileDescriptorLimit',
            expr: |||
              (
                node_filefd_allocated{%(filteringSelector)s} * 100 / node_filefd_maximum{%(filteringSelector)s} > 70
              )
            ||| % this.config,
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
                node_filefd_allocated{%(filteringSelector)s} * 100 / node_filefd_maximum{%(filteringSelector)s} > 90
              )
            ||| % this.config,
            'for': '15m',
            labels: {
              severity: 'critical',
            },
            annotations: {
              summary: 'Kernel is predicted to exhaust file descriptors limit soon.',
              description: 'File descriptors limit at {{ $labels.instance }} is currently at {{ printf "%.2f" $value }}%.',
            },
          },
          {
            alert: 'NodeCPUHighUsage',
            expr: |||
              sum without(mode) (avg without (cpu) (rate(node_cpu_seconds_total{%(filteringSelector)s, mode!="idle"}[2m]))) * 100 > %(cpuHighUsageThreshold)d
            ||| % this.config,
            'for': '15m',
            labels: {
              severity: 'info',
            },
            annotations: {
              summary: 'High CPU usage.',
              description: |||
                CPU usage at {{ $labels.instance }} has been above %(cpuHighUsageThreshold)d%% for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}%%.
              ||| % this.config,
            },
          },
          {
            alert: 'NodeSystemSaturation',
            expr: |||
              node_load1{%(filteringSelector)s}
              / count without (cpu, mode) (node_cpu_seconds_total{%(filteringSelector)s, mode="idle"}) > %(systemSaturationPerCoreThreshold)d
            ||| % this.config,
            'for': '15m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'System saturated, load per core is very high.',
              description: |||
                System load per core at {{ $labels.instance }} has been above %(systemSaturationPerCoreThreshold)d for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}.
                This might indicate this instance resources saturation and can cause it becoming unresponsive.
              ||| % this.config,
            },
          },
          {
            alert: 'NodeMemoryMajorPagesFaults',
            expr: |||
              rate(node_vmstat_pgmajfault{%(filteringSelector)s}[5m]) > %(memoryMajorPagesFaultsThreshold)d
            ||| % this.config,
            'for': '15m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Memory major page faults are occurring at very high rate.',
              description: |||
                Memory major pages are occurring at very high rate at {{ $labels.instance }}, %(memoryMajorPagesFaultsThreshold)d major page faults per second for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}.
                Please check that there is enough memory available at this instance.
              ||| % this.config,
            },
          },
          {
            alert: 'NodeMemoryHighUtilization',
            expr: |||
              100 - (node_memory_MemAvailable_bytes{%(filteringSelector)s} / node_memory_MemTotal_bytes{%(filteringSelector)s} * 100) > %(memoryHighUtilizationThreshold)d
            ||| % this.config,
            'for': '15m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Host is running out of memory.',
              description: |||
                Memory is filling up at {{ $labels.instance }}, has been above %(memoryHighUtilizationThreshold)d%% for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}%%.
              ||| % this.config,
            },
          },
          {
            alert: 'NodeDiskIOSaturation',
            expr: |||
              rate(node_disk_io_time_weighted_seconds_total{%(filteringSelector)s, %(diskDeviceSelector)s}[5m]) > %(diskIOSaturationThreshold)d
            ||| % this.config,
            'for': '30m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Disk IO queue is high.',
              description: |||
                Disk IO queue (aqu-sq) is high on {{ $labels.device }} at {{ $labels.instance }}, has been above %(diskIOSaturationThreshold)d for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}.
                This symptom might indicate disk saturation.
              ||| % this.config,
            },
          },
          {
            alert: 'NodeSystemdServiceFailed',
            expr: |||
              node_systemd_unit_state{%(filteringSelector)s, state="failed"} == 1
            ||| % this.config,
            'for': '5m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Systemd service has entered failed state.',
              description: 'Systemd service {{ $labels.name }} has entered failed state at {{ $labels.instance }}',
            },
          },
        ],
      },
    ],
  },
}
