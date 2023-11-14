local xtd = import './main.libsonnet';
local d = import 'doc-util/main.libsonnet';

{
  '#': d.pkg(
    name='jsonpath',
    url='github.com/jsonnet-libs/xtd/jsonpath.libsonnet',
    help='`jsonpath` implements helper functions to use JSONPath expressions.',
  ),


  '#getJSONPath':: d.fn(
    |||
      `getJSONPath` gets the value at `path` from `source` where path is a JSONPath.

      This is a rudimentary implementation supporting the slice operator `[0:3:2]` and
      partially supporting filter expressions `?(@.attr==value)`.
    |||,
    [
      d.arg('source', d.T.any),
      d.arg('path', d.T.string,),
      d.arg('default', d.T.any, default='null'),
    ]
  ),
  getJSONPath(source, path, default=null):
    local _path = self.convertBracketToDot(path);
    std.foldl(
      function(acc, key)
        if acc == null
        then acc
        else get(acc, key, default),
      xtd.string.splitEscape(_path, '.'),
      source,
    ),

  '#convertBracketToDot':: d.fn(
    |||
      `convertBracketToDot` converts the bracket notation to dot notation.

      This function does not support  escaping brackets/quotes in path keys.
    |||,
    [
      d.arg('path', d.T.string,),
    ]
  ),
  convertBracketToDot(path):
    if std.length(std.findSubstr('[', path)) > 0
    then
      local split = std.split(path, '[');
      std.join('.', [
        local a = std.stripChars(i, "[]'");
        std.strReplace(a, '@.', '@\\.')
        for i in split
      ])
    else path,

  local get(source, key, default) =
    if key == ''
       || key == '$'
       || key == '*'
    then source
    else if std.isArray(source)
    then getFromArray(source, key)
    else std.get(source, key, default),

  local getFromArray(arr, key) =
    if std.startsWith(key, '?(@\\.')
    then
      std.filter(
        self.parseFilterExpr(std.stripChars(key, '?(@\\.)')),
        arr
      )
    else if std.length(std.findSubstr(':', key)) >= 1
    then
      local split = std.splitLimit(key, ':', 2);
      local step =
        if std.length(split) < 3
        then 1
        else parseIntOrNull(split[2]);
      xtd.array.slice(
        arr,
        parseIntOrNull(split[0]),
        parseIntOrNull(split[1]),
        step,
      )
    else
      arr[std.parseInt(key)],

  local parseIntOrNull(str) =
    if str == ''
    then null
    else std.parseInt(str),

  '#parseFilterExpr':: d.fn(
    |||
      `parseFilterExpr` returns a filter function `f(x)` for a filter expression `expr`.

       It supports comparisons (<, <=, >, >=) and equality checks (==, !=). If it doesn't
       have an operator, it will check if the `expr` value exists.
    |||,
    [
      d.arg('path', d.T.string,),
    ]
  ),
  parseFilterExpr(expr):
    local operandFunctions = {
      '=='(a, b): a == b,
      '!='(a, b): a != b,
      '<='(a, b): a <= b,
      '>='(a, b): a >= b,
      '<'(a, b): a < b,
      '>'(a, b): a > b,
    };

    local findOperands = std.filter(
      function(op) std.length(std.findSubstr(op, expr)) > 0,
      std.reverse(  // reverse to match '<=' before '<'
        std.objectFields(operandFunctions)
      )
    );

    if std.length(findOperands) > 0
    then
      local op = findOperands[0];
      local s = [
        std.stripChars(i, ' ')
        for i in std.splitLimit(expr, op, 1)
      ];
      function(x)
        if s[0] in x
        then
          local left = x[s[0]];
          local right =
            if std.isNumber(left)
            then std.parseInt(s[1])  // Only parse if comparing numbers
            else s[1];
          operandFunctions[op](left, right)
        else false
    else
      // Default to key matching
      function(x) (expr in x),
}
