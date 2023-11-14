local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#': d.package.newSub('util', 'Helper functions that work well with Grafonnet.'),
  dashboard: (import './dashboard.libsonnet'),
  grid: (import './grid.libsonnet'),
  panel: (import './panel.libsonnet'),
  string: (import './string.libsonnet'),
}
