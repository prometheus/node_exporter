{
  local nodemixin = import './prom-mixin.libsonnet',
  local cpu = import './cpu.libsonnet',
  local system = import './system.libsonnet',
  local memory = import './memory.libsonnet',
  local disk = import './disk.libsonnet',
  local network = import './network.libsonnet',

  grafanaDashboards+:: {
    'nodes.json': nodemixin.new(config=$._config, platform='Linux').dashboard,
    'nodes-darwin.json': nodemixin.new(config=$._config, platform='Darwin').dashboard,
    'nodes-system.json': system.new(config=$._config, platform='Linux').dashboard,
    'nodes-memory.json': memory.new(config=$._config, platform='Linux').dashboard,
    'nodes-network.json': network.new(config=$._config, platform='Linux').dashboard,
    'nodes-disk.json': disk.new(config=$._config, platform='Linux').dashboard,
    'nodes-fleet.json': (import './fleet.libsonnet').new(config=$._config, platform='Linux').dashboard,
  },
}
