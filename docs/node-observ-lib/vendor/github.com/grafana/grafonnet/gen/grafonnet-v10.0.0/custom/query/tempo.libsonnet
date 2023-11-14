local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#new':: d.func.new(
    'Creates a new tempo query target for panels.',
    args=[
      d.arg('datasource', d.T.string),
      d.arg('query', d.T.string),
      d.arg('filters', d.T.array),
    ]
  ),
  new(datasource, query, filters):
    self.withDatasource(datasource)
    + self.withQuery(query)
    + self.withFilters(filters),

  '#withDatasource':: d.func.new(
    'Set the datasource for this query.',
    args=[
      d.arg('value', d.T.string),
    ]
  ),
  withDatasource(value): {
    datasource+: {
      type: 'tempo',
      uid: value,
    },
  },
}
