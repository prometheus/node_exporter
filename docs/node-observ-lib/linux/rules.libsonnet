{
  new(this): {
    groups+: [
      {
        name: if this.config.uid == 'node' then 'node-exporter.rules' else this.config.uid + '-linux-rules',
        rules: [
          {
            // This rule gives the number of CPUs per node.
            record: 'instance:node_num_cpu:sum',
            expr: |||
              count without (cpu, mode) (
                node_cpu_seconds_total{%(filteringSelector)s,mode="idle"}
              )
            ||| % this.config,
          },
          {
            // CPU utilisation is % CPU without {idle,iowait,steal}.
            record: 'instance:node_cpu_utilisation:rate%(rateInterval)s' % this.config,
            expr: |||
              1 - avg without (cpu) (
                sum without (mode) (rate(node_cpu_seconds_total{%(filteringSelector)s, mode=~"idle|iowait|steal"}[%(rateInterval)s]))
              )
            ||| % this.config,
          },
          {
            // This is CPU saturation: 1min avg run queue length / number of CPUs.
            // Can go over 1.
            // TODO: There are situation where a run queue >1/core is just normal and fine.
            //       We need to clarify how to read this metric and if its usage is helpful at all.
            record: 'instance:node_load1_per_cpu:ratio',
            expr: |||
              (
                node_load1{%(filteringSelector)s}
              /
                instance:node_num_cpu:sum{%(filteringSelector)s}
              )
            ||| % this.config,
          },
          {
            // Memory utilisation (ratio of used memory per instance).
            record: 'instance:node_memory_utilisation:ratio',
            expr: |||
              1 - (
                (
                  node_memory_MemAvailable_bytes{%(filteringSelector)s}
                  or
                  (
                    node_memory_Buffers_bytes{%(filteringSelector)s}
                    +
                    node_memory_Cached_bytes{%(filteringSelector)s}
                    +
                    node_memory_MemFree_bytes{%(filteringSelector)s}
                    +
                    node_memory_Slab_bytes{%(filteringSelector)s}
                  )
                )
              /
                node_memory_MemTotal_bytes{%(filteringSelector)s}
              )
            ||| % this.config,
          },
          {
            record: 'instance:node_vmstat_pgmajfault:rate%(rateInterval)s' % this.config,
            expr: |||
              rate(node_vmstat_pgmajfault{%(filteringSelector)s}[%(rateInterval)s])
            ||| % this.config,
          },
          {
            // Disk utilisation (seconds spent, 1 second rate).
            record: 'instance_device:node_disk_io_time_seconds:rate%(rateInterval)s' % this.config,
            expr: |||
              rate(node_disk_io_time_seconds_total{%(filteringSelector)s, %(diskDeviceSelector)s}[%(rateInterval)s])
            ||| % this.config,
          },
          {
            // Disk saturation (weighted seconds spent, 1 second rate).
            record: 'instance_device:node_disk_io_time_weighted_seconds:rate%(rateInterval)s' % this.config,
            expr: |||
              rate(node_disk_io_time_weighted_seconds_total{%(filteringSelector)s, %(diskDeviceSelector)s}[%(rateInterval)s])
            ||| % this.config,
          },
          {
            record: 'instance:node_network_receive_bytes_excluding_lo:rate%(rateInterval)s' % this.config,
            expr: |||
              sum without (device) (
                rate(node_network_receive_bytes_total{%(filteringSelector)s, device!="lo"}[%(rateInterval)s])
              )
            ||| % this.config,
          },
          {
            record: 'instance:node_network_transmit_bytes_excluding_lo:rate%(rateInterval)s' % this.config,
            expr: |||
              sum without (device) (
                rate(node_network_transmit_bytes_total{%(filteringSelector)s, device!="lo"}[%(rateInterval)s])
              )
            ||| % this.config,
          },
          // TODO: Find out if those drops ever happen on modern switched networks.
          {
            record: 'instance:node_network_receive_drop_excluding_lo:rate%(rateInterval)s' % this.config,
            expr: |||
              sum without (device) (
                rate(node_network_receive_drop_total{%(filteringSelector)s, device!="lo"}[%(rateInterval)s])
              )
            ||| % this.config,
          },
          {
            record: 'instance:node_network_transmit_drop_excluding_lo:rate%(rateInterval)s' % this.config,
            expr: |||
              sum without (device) (
                rate(node_network_transmit_drop_total{%(filteringSelector)s, device!="lo"}[%(rateInterval)s])
              )
            ||| % this.config,
          },
        ],
      },
    ],
  },
}
