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
                predict_linear(node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[%(fsSpaceFillingUpPredictionWindow)s], %(nodeWarningWindowHours)s*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of space within the next %(nodeWarningWindowHours)s hours.' % $._config,
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left and is filling up.',
            },
          },
          {
            alert: 'NodeFilesystemSpaceFillingUp',
            expr: |||
              (
                node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_size_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < %(fsSpaceFillingUpCriticalThreshold)d
              and
                predict_linear(node_filesystem_avail_bytes{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], %(nodeCriticalWindowHours)s*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % $._config,
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of space within the next %(nodeCriticalWindowHours)s hours.' % $._config,
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left and is filling up fast.',
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
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left.',
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
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available space left.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 40
              and
                predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], %(nodeWarningWindowHours)s*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of inodes within the next %(nodeWarningWindowHours)s hours.' % $._config,
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left and is filling up.',
            },
          },
          {
            alert: 'NodeFilesystemFilesFillingUp',
            expr: |||
              (
                node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} / node_filesystem_files{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} * 100 < 20
              and
                predict_linear(node_filesystem_files_free{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s}[6h], %(nodeCriticalWindowHours)s*60*60) < 0
              and
                node_filesystem_readonly{%(nodeExporterSelector)s,%(fsSelector)s,%(fsMountpointSelector)s} == 0
              )
            ||| % $._config,
            'for': '1h',
            labels: {
              severity: '%(nodeCriticalSeverity)s' % $._config,
            },
            annotations: {
              summary: 'Filesystem is predicted to run out of inodes within the next %(nodeCriticalWindowHours)s hours.' % $._config,
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left and is filling up fast.',
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
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left.',
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
              description: 'Filesystem on {{ $labels.device }}, mounted on {{ $labels.mountpoint }}, at {{ $labels.instance }} has only {{ printf "%.2f" $value }}% available inodes left.',
            },
          },
          {
            alert: 'NodeNetworkReceiveErrs',
            expr: |||
              rate(node_network_receive_errs_total{%(nodeExporterSelector)s}[2m]) / rate(node_network_receive_packets_total{%(nodeExporterSelector)s}[2m]) > 0.01
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
              rate(node_network_transmit_errs_total{%(nodeExporterSelector)s}[2m]) / rate(node_network_transmit_packets_total{%(nodeExporterSelector)s}[2m]) > 0.01
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
              (node_nf_conntrack_entries{%(nodeExporterSelector)s} / node_nf_conntrack_entries_limit) > 0.75
            ||| % $._config,
            annotations: {
              summary: 'Number of conntrack are getting close to the limit.',
              description: '{{ $labels.instance }} {{ $value | humanizePercentage }} of conntrack entries are used.',
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
              description: 'Clock at {{ $labels.instance }} is out of sync by more than 0.05s. Ensure NTP is configured correctly on this host.',
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
              description: 'Clock at {{ $labels.instance }} is not synchronising. Ensure NTP is configured on this host.',
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
              summary: 'RAID Array is degraded.',
              description: "RAID array '{{ $labels.device }}' at {{ $labels.instance }} is in degraded state due to one or more disks failures. Number of spare drives is insufficient to fix issue automatically.",
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
              summary: 'Failed device in RAID array.',
              description: "At least one device in RAID array at {{ $labels.instance }} failed. Array '{{ $labels.device }}' needs attention and possibly a disk swap.",
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
          {
            alert: 'NodeCPUHighUsage',
            expr: |||
              sum without(mode) (avg without (cpu) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode!~"idle|iowait"}[2m]))) * 100 > %(cpuHighUsageThreshold)d
            ||| % $._config,
            'for': '15m',
            labels: {
              severity: 'info',
            },
            annotations: {
              summary: 'High CPU usage.',
              description: |||
                CPU usage at {{ $labels.instance }} has been above %(cpuHighUsageThreshold)d%% for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}%%.
              ||| % $._config,
            },
          },
          {
            alert: 'NodeSystemSaturation',
            expr: |||
              node_load1{%(nodeExporterSelector)s}
              / count without (cpu, mode) (node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle"}) > %(systemSaturationPerCoreThreshold)d
            ||| % $._config,
            'for': '15m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'System saturated, load per core is very high.',
              description: |||
                System load per core at {{ $labels.instance }} has been above %(systemSaturationPerCoreThreshold)d for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}.
                This might indicate this instance resources saturation and can cause it becoming unresponsive.
              ||| % $._config,
            },
          },
          {
            alert: 'NodeMemoryMajorPagesFaults',
            expr: |||
              rate(node_vmstat_pgmajfault{%(nodeExporterSelector)s}[5m]) > %(memoryMajorPagesFaultsThreshold)d
            ||| % $._config,
            'for': '15m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Memory major page faults are occurring at very high rate.',
              description: |||
                Memory major pages are occurring at very high rate at {{ $labels.instance }}, %(memoryMajorPagesFaultsThreshold)d major page faults per second for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}.
                Please check that there is enough memory available at this instance.
              ||| % $._config,
            },
          },
          {
            alert: 'NodeMemoryHighUtilization',
            expr: |||
              100 - (node_memory_MemAvailable_bytes{%(nodeExporterSelector)s} / node_memory_MemTotal_bytes{%(nodeExporterSelector)s} * 100) > %(memoryHighUtilizationThreshold)d
            ||| % $._config,
            'for': '15m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Host is running out of memory.',
              description: |||
                Memory is filling up at {{ $labels.instance }}, has been above %(memoryHighUtilizationThreshold)d%% for the last 15 minutes, is currently at {{ printf "%%.2f" $value }}%%.
              ||| % $._config,
            },
          },
          {
            alert: 'NodeDiskIOSaturation',
            expr: |||
              rate(node_disk_io_time_weighted_seconds_total{%(nodeExporterSelector)s, %(diskDeviceSelector)s}[5m]) > %(diskIOSaturationThreshold)d
            ||| % $._config,
            'for': '30m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Disk IO queue is high.',
              description: |||
                Disk IO queue (aqu-sq) is high on {{ $labels.device }} at {{ $labels.instance }}, has been above %(diskIOSaturationThreshold)d for the last 30 minutes, is currently at {{ printf "%%.2f" $value }}.
                This symptom might indicate disk saturation.
              ||| % $._config,
            },
          },
          {
            alert: 'NodeSystemdServiceFailed',
            expr: |||
              node_systemd_unit_state{%(nodeExporterSelector)s, state="failed"} == 1
            ||| % $._config,
            'for': '5m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Systemd service has entered failed state.',
              description: 'Systemd service {{ $labels.name }} has entered failed state at {{ $labels.instance }}',
            },
          },
          {
            alert: 'NodeSystemdServiceCrashlooping',
            expr: |||
              increase(node_systemd_service_restart_total{%(nodeExporterSelector)s}[5m]) > 2
            ||| % $._config,
            'for': '15m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Systemd service keeps restaring, possibly crash looping.',
              description: 'Systemd service {{ $labels.name }} has being restarted too many times at {{ $labels.instance }} for the last 15 minutes. Please check if service is crash looping.',
            },
          },
          {
            alert: 'NodeBondingDegraded',
            expr: |||
              (node_bonding_slaves{%(nodeExporterSelector)s} - node_bonding_active{%(nodeExporterSelector)s}) != 0
            ||| % $._config,
            'for': '5m',
            labels: {
              severity: 'warning',
            },
            annotations: {
              summary: 'Bonding interface is degraded.',
              description: 'Bonding interface {{ $labels.master }} on {{ $labels.instance }} is in degraded state due to one or more slave failures.',
            },
          },
        ],
      },
    ],
  },
}
