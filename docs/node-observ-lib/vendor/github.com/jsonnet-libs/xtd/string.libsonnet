local d = import 'doc-util/main.libsonnet';

{
  '#': d.pkg(
    name='string',
    url='github.com/jsonnet-libs/xtd/string.libsonnet',
    help='`string` implements helper functions for processing strings.',
  ),

  // BelRune is a string of the Ascii character BEL which made computers ring in ancient times.
  // We use it as "magic" char to temporarily replace an escaped string as it is a non printable
  // character and thereby will unlikely be in a valid key by accident. Only when we include it.
  local BelRune = std.char(7),

  '#splitEscape':: d.fn(
    |||
      `split` works the same as `std.split` but with support for escaping the dividing
      string `c`.
    |||,
    [
      d.arg('str', d.T.string),
      d.arg('c', d.T.string),
      d.arg('escape', d.T.string, default='\\'),
    ]
  ),
  splitEscape(str, c, escape='\\'):
    std.map(
      function(i)
        std.strReplace(i, BelRune, escape + c),
      std.split(
        std.strReplace(str, escape + c, BelRune),
        c,
      )
    ),
}
