local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#new':: d.func.new(
    'Creates a new row panel with a title.',
    args=[d.arg('title', d.T.string)]
  ),
  new(title):
    self.withTitle(title)
    + self.withType(),
}
