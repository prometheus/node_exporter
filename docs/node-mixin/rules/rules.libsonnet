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
                sum without (mode) (
                  node_cpu_seconds_total{%(nodeExporterSelector)s}
                )
              )
            ||| % $._config,
          },
          {
            // CPU utilisation is % CPU is not idle.
            record: 'instance:node_cpu_utilisation:avg_rate1m',
            expr: |||
              1 - avg without (cpu, mode) (
                rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode="idle"}[1m])
              )
            ||| % $._config,
          },
          {
            // This is CPU saturation: 1min avg run queue length / number of CPUs.
            // Can go over 1. >1 is bad.
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
            // Memory utilisation per node, normalized by per-node memory
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
            record: 'instance:node_memory_swap_io_pages:rate1m',
            expr: |||
              (
                rate(node_vmstat_pgpgin{%(nodeExporterSelector)s}[1m])
              +
                rate(node_vmstat_pgpgout{%(nodeExporterSelector)s}[1m])
              )
            ||| % $._config,
          },
          {
            // Disk utilisation (seconds spent, 1 second rate)
            record: 'instance:node_disk_io_time:sum_rate1m',
            expr: |||
              sum without (device) (
                rate(node_disk_io_time_seconds_total{%(nodeExporterSelector)s, %(diskDeviceSelector)s}[1m])
              )
            ||| % $._config,
          },
          {
            // Disk saturation (weighted seconds spent, 1 second rate)
            record: 'instance:node_disk_io_time_weighted:sum_rate1m',
            expr: |||
              sum without (device) (
                rate(node_disk_io_time_weighted_seconds_total{%(nodeExporterSelector)s, %(diskDeviceSelector)s}[1m])
              )
            ||| % $._config,
          },
          // TODO: For the following rules, consider configurable filtering to exclude more network
          // device names than just "lo".
          {
            record: 'instance:node_network_receive_bytes:sum_rate1m',
            expr: |||
              sum without (device) (
                rate(node_network_receive_bytes_total{%(nodeExporterSelector)s, device!="lo"}[1m])
              )
            ||| % $._config,
          },
          {
            record: 'instance:node_network_transmit_bytes:sum_rate1m',
            expr: |||
              sum without (device) (
                rate(node_network_transmit_bytes_total{%(nodeExporterSelector)s, device!="lo"}[1m])
              )
            ||| % $._config,
          },
          {
            record: 'instance:node_network_receive_drop:sum_rate1m',
            expr: |||
              sum without (device) (
                rate(node_network_receive_drop_total{%(nodeExporterSelector)s, device!="lo"}[1m])
              )
            ||| % $._config,
          },
          {
            record: 'instance:node_network_transmit_drop:sum_rate1m',
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
