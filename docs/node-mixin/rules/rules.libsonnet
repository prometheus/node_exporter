{
  prometheusRules+:: {
    groups+: [
      {
        name: 'node-exporter.rules',
        rules: [
          {
            // This rule gives the number of CPUs per node.
            record: 'instance:node_num_cpu:sum',
            expr: |||
              count without (cpu) (
                count without (mode) (
                  node_cpu_seconds_total{%(nodeExporterSelector)s}
                )
              )
            ||| % $._config,
          },
          {
            // CPU utilisation is % CPU is not idle.
            record: 'instance:node_cpu_utilisation:rate1m',
            expr: |||
              1 - avg without (cpu, mode) (
                rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle"}[1m])
              )
            ||| % $._config,
          },
          {
            // One minute load average per CPU.
            // This metrics is taken from /proc/loadavg file.
            // It represents the number of jobs in the run queue (processes that are running; state R),
            // or waiting for disk I/O (state D) averaged over 1 minute.
            // Values over 1 don't have to be outright considered problematic,
            // for example there can be multiple processes trying to access the disk.
            record: 'instance:node_load1_per_cpu:ratio',
            expr: |||
              (
                node_load1{%(nodeExporterSelector)s}
              /
                instance:node_num_cpu:sum{%(nodeExporterSelector)s}
              )
            ||| % $._config,
          },
          {
            // Memory utilisation (ratio of used memory per instance).
            record: 'instance:node_memory_utilisation:ratio',
            expr: |||
              1 - (
                node_memory_MemAvailable_bytes{%(nodeExporterSelector)s}
              /
                node_memory_MemTotal_bytes{%(nodeExporterSelector)s}
              )
            ||| % $._config,
          },
          {
            record: 'instance:node_vmstat_pgmajfault:rate1m',
            expr: |||
              rate(node_vmstat_pgmajfault{%(nodeExporterSelector)s}[1m])
            ||| % $._config,
          },
          {
            // Disk utilisation (seconds spent, 1 second rate).
            record: 'instance_device:node_disk_io_time_seconds:rate1m',
            expr: |||
              rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, %(diskDeviceSelector)s}[1m])
            ||| % $._config,
          },
          {
            // Disk saturation (weighted seconds spent, 1 second rate).
            record: 'instance_device:node_disk_io_time_weighted_seconds:rate1m',
            expr: |||
              rate(node_disk_io_time_weighted_seconds_total{%(nodeExporterSelector)s, %(diskDeviceSelector)s}[1m])
            ||| % $._config,
          },
          {
            record: 'instance:node_network_receive_bytes_excluding_lo:rate1m',
            expr: |||
              sum without (device) (
                rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, device!="lo"}[1m])
              )
            ||| % $._config,
          },
          {
            record: 'instance:node_network_transmit_bytes_excluding_lo:rate1m',
            expr: |||
              sum without (device) (
                rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, device!="lo"}[1m])
              )
            ||| % $._config,
          },
          // TODO: Find out if those drops ever happen on modern switched networks.
          {
            record: 'instance:node_network_receive_drop_excluding_lo:rate1m',
            expr: |||
              sum without (device) (
                rate(node_network_receive_drop_total{%(nodeExporterSelector)s, device!="lo"}[1m])
              )
            ||| % $._config,
          },
          {
            record: 'instance:node_network_transmit_drop_excluding_lo:rate1m',
            expr: |||
              sum without (device) (
                rate(node_network_transmit_drop_total{%(nodeExporterSelector)s, device!="lo"}[1m])
              )
            ||| % $._config,
          },
        ],
      },
    ],
  },
}
