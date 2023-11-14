---
permalink: /camelcase/
---

# camelcase

```jsonnet
local camelcase = import "github.com/jsonnet-libs/xtd/camelcase.libsonnet"
```

`camelcase` can split camelCase words into an array of words.

## Index

* [`fn split(src)`](#fn-split)
* [`fn toCamelCase(str)`](#fn-tocamelcase)

## Fields

### fn split

```ts
split(src)
```

`split` splits a camelcase word and returns an array  of words. It also supports
digits. Both lower camel case and upper camel case are supported. It only supports
ASCII characters.
For more info please check: http://en.wikipedia.org/wiki/CamelCase
Based on https://github.com/fatih/camelcase/


### fn toCamelCase

```ts
toCamelCase(str)
```

`toCamelCase` transforms a string to camelCase format, splitting words by the `-`, `_` or spaces.
For example: `hello_world` becomes `helloWorld`.
For more info please check: http://en.wikipedia.org/wiki/CamelCase
