local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#withTitle': { ignore: true },
  '#withName': super['#withTitle'],
  withName: super.withTitle,
  rule+: {
    '#':: d.package.newSub('rule', ''),
    '#withTitle': { ignore: true },
    '#withName': super['#withTitle'],
    withName: super.withTitle,
  },
}
