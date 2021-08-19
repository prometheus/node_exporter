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
            // CPU utilisation is % CPU is not idle. This represents CPU saturation.
            record: 'instance:node_cpu_utilisation:rate%(rateInterval)s' % $._config,
            expr: |||
              1 - avg without (cpu, mode) (
                rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle"}[%(rateInterval)s])
              )
            ||| % $._config,
          },
          {
            // CPU pressure represents over-saturation. This is the amount of CPU seconds
            // requested, but the kernel was not able to schedule.
            // NOTE: This is only availalbe on Linux >= 4.19 and `CONFIG_PSI` is enabled.
            // See also:
            // - https://www.kernel.org/doc/html/latest/accounting/psi.html
            // - https://facebookmicrosites.github.io/psi/docs/overview
            expr: |||
              rate(node_pressure_cpu_waiting_seconds_total{%(nodeExporterSelector)s}[%(rateInterval)s])
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
