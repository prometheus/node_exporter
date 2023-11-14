local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';
local xtd = import 'github.com/jsonnet-libs/xtd/main.libsonnet';

{
  '#slugify':: d.func.new(
    |||
      `slugify` will create a simple slug from `string`, keeping only alphanumeric
      characters and replacing spaces with dashes.
    |||,
    args=[d.arg('string', d.T.string)]
  ),
  slugify(string):
    std.strReplace(
      std.asciiLower(
        std.join('', [
          string[i]
          for i in std.range(0, std.length(string) - 1)
          if xtd.ascii.isUpper(string[i])
                 || xtd.ascii.isLower(string[i])
                 || xtd.ascii.isNumber(string[i])
                 || string[i] == ' '
        ])
      ),
      ' ',
      '-',
    ),
}
