{
  new(config):: {
    cpu: (import './cpu.libsonnet').new(config),
    disk: (import './disk.libsonnet').new(config),
    fleet: (import './fleet.libsonnet').new(config),
    hardware: (import './hardware.libsonnet').new(config),
    memory: (import './memory.libsonnet').new(config),
    network: (import './network.libsonnet').new(config),
    system: (import './system.libsonnet').new(config),
    use: (import './use.libsonnet').new(config),
    useCluster: (import './useCluster.libsonnet').new(config),
    useClusterMulti: (import './useClusterMulti.libsonnet').new(config),
  },
}
