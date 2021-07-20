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
            record: 'instance:node_cpu_utilisation:rate%(rateInterval)s' % $._config,
            expr: |||
              1 - avg without (cpu, mode) (
                rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle"}[%(rateInterval)s])
              )
            ||| % $._config,
          },
          {
            // This is CPU saturation, that is 1 minute load average per CPU.
            // I.e. in the Linux kernel it is the number of processes in "running" state,
            // plus processes in "uninterruptible sleep" state (which are usually processes waiting for I/O),
            // averaged over one minute, and per CPU.
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
            record: 'instance:node_vmstat_pgmajfault:rate%(rateInterval)s' % $._config,
            expr: |||
              rate(node_vmstat_pgmajfault{%(nodeExporterSelector)s}[%(rateInterval)s])
            ||| % $._config,
          },
          {
            // Disk utilisation (seconds spent, 1 second rate).
            record: 'instance_device:node_disk_io_time_seconds:rate%(rateInterval)s' % $._config,
            expr: |||
              rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, %(diskDeviceSelector)s}[%(rateInterval)s])
            ||| % $._config,
          },
          {
            // Disk saturation (weighted seconds spent, 1 second rate).
            record: 'instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s' % $._config,
            expr: |||
              rate(node_disk_io_time_weighted_seconds_total{%(nodeExporterSelector)s, %(diskDeviceSelector)s}[%(rateInterval)s])
            ||| % $._config,
          },
          {
            record: 'instance:node_network_receive_bytes_excluding_lo:rate%(rateInterval)s' % $._config,
            expr: |||
              sum without (device) (
                rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, device!="lo"}[%(rateInterval)s])
              )
            ||| % $._config,
          },
          {
            record: 'instance:node_network_transmit_bytes_excluding_lo:rate%(rateInterval)s' % $._config,
            expr: |||
              sum without (device) (
                rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, device!="lo"}[%(rateInterval)s])
              )
            ||| % $._config,
          },
          // TODO: Find out if those drops ever happen on modern switched networks.
          {
            record: 'instance:node_network_receive_drop_excluding_lo:rate%(rateInterval)s' % $._config,
            expr: |||
              sum without (device) (
                rate(node_network_receive_drop_total{%(nodeExporterSelector)s, device!="lo"}[%(rateInterval)s])
              )
            ||| % $._config,
          },
          {
            record: 'instance:node_network_transmit_drop_excluding_lo:rate%(rateInterval)s' % $._config,
            expr: |||
              sum without (device) (
                rate(node_network_transmit_drop_total{%(nodeExporterSelector)s, device!="lo"}[%(rateInterval)s])
              )
            ||| % $._config,
          },
        ],
      },
    ],
  },
}
