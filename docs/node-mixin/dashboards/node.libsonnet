{
  local nodemixin = import '../lib/prom-mixin.libsonnet',
  grafanaDashboards+:: {
    'nodes.json': nodemixin.new(config=$._config, platform='Linux').dashboard,
    'nodes-darwin.json': nodemixin.new(config=$._config, platform='Darwin').dashboard,
  },
}
