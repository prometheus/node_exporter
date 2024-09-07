{
  local nodemixin = import '../lib/prom-mixin.libsonnet',
  grafanaDashboards+:: {
    'nodes.json': nodemixin.new(config=$._config, platform='Linux', uid=std.md5('nodes.json')).dashboard,
    'nodes-darwin.json': nodemixin.new(config=$._config, platform='Darwin', uid=std.md5('nodes-darwin.json')).dashboard,
    'nodes-aix.json': nodemixin.new(config=$._config, platform='AIX', uid=std.md5('nodes-aix.json')).dashboard,
  },
}
