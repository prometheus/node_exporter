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
              count without (cpu, mode) (
                node_cpu_seconds_total{%(nodeExporterSelector)s,mode="idle"}
              )
            ||| % $._config,
          },
          {
            // CPU utilisation is % CPU without {idle,iowait,steal}.
            record: 'instance:node_cpu_utilisation:rate%(rateInterval)s' % $._config,
            expr: |||
              1 - avg without (cpu) (
                sum without (mode) (rate(node_cpu_seconds_total{%(nodeExporterSelector)s, mode=~"idle|iowait|steal"}[%(rateInterval)s]))
              )
            ||| % $._config,
          },
          {
            // This is CPU saturation: 1min avg run queue length / number of CPUs.
            // Can go over 1.
            // TODO: There are situation where a run queue >1/core is just normal and fine.
            //       We need to clarify how to read this metric and if its usage is helpful at all.
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
                (
                  node_memory_MemAvailable_bytes{%(nodeExporterSelector)s}
                  or
                  (
                    node_memory_Buffers_bytes{%(nodeExporterSelector)s}
                    +
                    node_memory_Cached_bytes{%(nodeExporterSelector)s}
                    +
                    node_memory_MemFree_bytes{%(nodeExporterSelector)s}
                    +
                    node_memory_Slab_bytes{%(nodeExporterSelector)s}
                  )
                )
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
