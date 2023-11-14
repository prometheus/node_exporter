local xtd = import './main.libsonnet';
local d = import 'doc-util/main.libsonnet';

{
  '#': d.pkg(
    name='camelcase',
    url='github.com/jsonnet-libs/xtd/camelcase.libsonnet',
    help='`camelcase` can split camelCase words into an array of words.',
  ),

  '#split':: d.fn(
    |||
      `split` splits a camelcase word and returns an array  of words. It also supports
      digits. Both lower camel case and upper camel case are supported. It only supports
      ASCII characters.
      For more info please check: http://en.wikipedia.org/wiki/CamelCase
      Based on https://github.com/fatih/camelcase/
    |||,
    [d.arg('src', d.T.string)]
  ),
  split(src):
    if src == ''
    then ['']
    else
      local runes = std.foldl(
        function(acc, r)
          acc {
            local class =
              if xtd.ascii.isNumber(r)
              then 1
              else if xtd.ascii.isLower(r)
              then 2
              else if xtd.ascii.isUpper(r)
              then 3
              else 4,

            lastClass:: class,

            runes:
              if class == super.lastClass
              then super.runes[:std.length(super.runes) - 1]
                   + [super.runes[std.length(super.runes) - 1] + r]
              else super.runes + [r],
          },
        [src[i] for i in std.range(0, std.length(src) - 1)],
        { lastClass:: 0, runes: [] }
      ).runes;

      local fixRunes =
        std.foldl(
          function(runes, i)
            if xtd.ascii.isUpper(runes[i][0])
               && xtd.ascii.isLower(runes[i + 1][0])
               && !xtd.ascii.isNumber(runes[i + 1][0])
               && runes[i][0] != ' '
               && runes[i + 1][0] != ' '
            then
              std.mapWithIndex(
                function(index, r)
                  if index == i + 1
                  then runes[i][std.length(runes[i]) - 1:] + r
                  else
                    if index == i
                    then r[:std.length(r) - 1]
                    else r
                , runes
              )
            else runes
          ,
          [i for i in std.range(0, std.length(runes) - 2)],
          runes
        );

      [
        r
        for r in fixRunes
        if r != ''
      ],

  '#toCamelCase':: d.fn(
    |||
      `toCamelCase` transforms a string to camelCase format, splitting words by the `-`, `_` or spaces.
      For example: `hello_world` becomes `helloWorld`.
      For more info please check: http://en.wikipedia.org/wiki/CamelCase
    |||,
    [d.arg('str', d.T.string)]
  ),
  toCamelCase(str)::
    local separators = std.set(std.findSubstr('_', str) + std.findSubstr('-', str) + std.findSubstr(' ', str));
    local n = std.join('', [
      if std.setMember(i - 1, separators)
      then std.asciiUpper(str[i])
      else str[i]
      for i in std.range(0, std.length(str) - 1)
      if !std.setMember(i, separators)
    ]);
    if std.length(n) == 0
    then n
    else std.asciiLower(n[0]) + n[1:],
}
