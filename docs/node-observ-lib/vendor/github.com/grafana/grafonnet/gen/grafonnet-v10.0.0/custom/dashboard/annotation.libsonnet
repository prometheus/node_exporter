local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';

{
  '#annotation':: {},

  '#withAnnotations':
    d.func.new(
      |||
        `withAnnotations` adds an array of annotations to a dashboard.

        This function appends passed data to existing values
      |||,
      args=[d.arg('value', d.T.array)]
    ),
  withAnnotations(value): super.annotation.withList(value),

  '#withAnnotationsMixin':
    d.func.new(
      |||
        `withAnnotationsMixin` adds an array of annotations to a dashboard.

        This function appends passed data to existing values
      |||,
      args=[d.arg('value', d.T.array)]
    ),
  withAnnotationsMixin(value): super.annotation.withListMixin(value),

  annotation:
    super.annotation.list
    + {
      '#':: d.package.newSub(
        'annotation',
        '',
      ),
    },
}
