---
permalink: /array/
---

# array

```jsonnet
local array = import "github.com/jsonnet-libs/xtd/array.libsonnet"
```

`array` implements helper functions for processing arrays.

## Index

* [`fn slice(indexable, index, end='null', step=1)`](#fn-slice)

## Fields

### fn slice

```ts
slice(indexable, index, end='null', step=1)
```

`slice` works the same as `std.slice` but with support for negative index/end.