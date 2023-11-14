---
permalink: /jsonpath/
---

# jsonpath

```jsonnet
local jsonpath = import "github.com/jsonnet-libs/xtd/jsonpath.libsonnet"
```

`jsonpath` implements helper functions to use JSONPath expressions.

## Index

* [`fn convertBracketToDot(path)`](#fn-convertbrackettodot)
* [`fn getJSONPath(source, path, default='null')`](#fn-getjsonpath)
* [`fn parseFilterExpr(path)`](#fn-parsefilterexpr)

## Fields

### fn convertBracketToDot

```ts
convertBracketToDot(path)
```

`convertBracketToDot` converts the bracket notation to dot notation.

This function does not support  escaping brackets/quotes in path keys.


### fn getJSONPath

```ts
getJSONPath(source, path, default='null')
```

`getJSONPath` gets the value at `path` from `source` where path is a JSONPath.

This is a rudimentary implementation supporting the slice operator `[0:3:2]` and
partially supporting filter expressions `?(@.attr==value)`.


### fn parseFilterExpr

```ts
parseFilterExpr(path)
```

`parseFilterExpr` returns a filter function `f(x)` for a filter expression `expr`.

 It supports comparisons (<, <=, >, >=) and equality checks (==, !=). If it doesn't
 have an operator, it will check if the `expr` value exists.
