local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#new':: d.func.new(
    'Creates a new prometheus query target for panels.',
    args=[
      d.arg('datasource', d.T.string),
      d.arg('expr', d.T.string),
    ]
  ),
  new(datasource, expr):
    self.withDatasource(datasource)
    + self.withExpr(expr),

  '#withIntervalFactor':: d.func.new(
    'Set the interval factor for this query.',
    args=[
      d.arg('value', d.T.string),
    ]
  ),
  withIntervalFactor(value): {
    intervalFactor: value,
  },

  '#withLegendFormat':: d.func.new(
    'Set the legend format for this query.',
    args=[
      d.arg('value', d.T.string),
    ]
  ),
  withLegendFormat(value): {
    legendFormat: value,
  },

  '#withDatasource':: d.func.new(
    'Set the datasource for this query.',
    args=[
      d.arg('value', d.T.string),
    ]
  ),
  withDatasource(value): {
    datasource+: {
      type: 'prometheus',
      uid: value,
    },
  },
}
