local d = import 'doc-util/main.libsonnet';

{
  '#': d.pkg(
    name='array',
    url='github.com/jsonnet-libs/xtd/array.libsonnet',
    help='`array` implements helper functions for processing arrays.',
  ),

  '#slice':: d.fn(
    '`slice` works the same as `std.slice` but with support for negative index/end.',
    [
      d.arg('indexable', d.T.array),
      d.arg('index', d.T.number),
      d.arg('end', d.T.number, default='null'),
      d.arg('step', d.T.number, default=1),
    ]
  ),
  slice(indexable, index, end=null, step=1):
    local invar = {
      index:
        if index != null
        then
          if index < 0
          then std.length(indexable) + index
          else index
        else 0,
      end:
        if end != null
        then
          if end < 0
          then std.length(indexable) + end
          else end
        else std.length(indexable),
    };
    indexable[invar.index:invar.end:step],
}
