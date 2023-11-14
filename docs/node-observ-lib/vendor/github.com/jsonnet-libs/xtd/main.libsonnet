local d = import 'doc-util/main.libsonnet';

{
  '#': d.pkg(
    name='xtd',
    url='github.com/jsonnet-libs/xtd/main.libsonnet',
    help=|||
      `xtd` aims to collect useful functions not included in the Jsonnet standard library (`std`).

      This package serves as a test field for functions intended to be contributed to `std`
      in the future, but also provides a place for less general, yet useful utilities.
    |||,
  ),

  aggregate: (import './aggregate.libsonnet'),
  array: (import './array.libsonnet'),
  ascii: (import './ascii.libsonnet'),
  camelcase: (import './camelcase.libsonnet'),
  date: (import './date.libsonnet'),
  inspect: (import './inspect.libsonnet'),
  jsonpath: (import './jsonpath.libsonnet'),
  string: (import './string.libsonnet'),
  url: (import './url.libsonnet'),
}
