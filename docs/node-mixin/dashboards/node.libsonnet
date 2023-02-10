{
  local nodemixin = import './prom-mixin.libsonnet',
  local network = import './network.libsonnet',
  grafanaDashboards+:: {
    'nodes.json': nodemixin.new(config=$._config, platform='Linux').dashboard,
    'nodes-darwin.json': nodemixin.new(config=$._config, platform='Darwin').dashboard,
    'nodes-network.json': network.new(config=$._config, platform='Linux').dashboard,
    'nodes-fleet.json': (import './fleet.libsonnet').new(config=$._config, platform='Linux').dashboard,
  },
}
