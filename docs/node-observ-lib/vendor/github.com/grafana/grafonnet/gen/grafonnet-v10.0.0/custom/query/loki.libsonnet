local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#new':: d.func.new(
    'Creates a new loki query target for panels.',
    args=[
      d.arg('datasource', d.T.string),
      d.arg('expr', d.T.string),
    ]
  ),
  new(datasource, expr):
    self.withDatasource(datasource)
    + self.withExpr(expr),

  '#withDatasource':: d.func.new(
    'Set the datasource for this query.',
    args=[
      d.arg('value', d.T.string),
    ]
  ),
  withDatasource(value): {
    datasource+: {
      type: 'loki',
      uid: value,
    },
  },
}
