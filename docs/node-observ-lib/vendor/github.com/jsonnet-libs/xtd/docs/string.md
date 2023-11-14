---
permalink: /string/
---

# string

```jsonnet
local string = import "github.com/jsonnet-libs/xtd/string.libsonnet"
```

`string` implements helper functions for processing strings.

## Index

* [`fn splitEscape(str, c, escape='\\')`](#fn-splitescape)

## Fields

### fn splitEscape

```ts
splitEscape(str, c, escape='\\')
```

`split` works the same as `std.split` but with support for escaping the dividing
string `c`.
