local dashboards = (import 'mixin.libsonnet').grafanaDashboards;

{
  [name]: dashboards[name] + { uid: std.md5(name) },
  for name in std.objectFields(dashboards)
}
