---
permalink: /inspect/
---

# inspect

```jsonnet
local inspect = import "github.com/jsonnet-libs/xtd/inspect.libsonnet"
```

`inspect` implements helper functions for inspecting Jsonnet

## Index

* [`fn diff(input1, input2)`](#fn-diff)
* [`fn filterKubernetesObjects(object, kind='')`](#fn-filterkubernetesobjects)
* [`fn filterObjects(filter_func, x)`](#fn-filterobjects)
* [`fn inspect(object, maxDepth)`](#fn-inspect)

## Fields

### fn diff

```ts
diff(input1, input2)
```

`diff` returns a JSON object describing the differences between two inputs. It
attemps to show diffs in nested objects and arrays too.

Simple example:

```jsonnet
local input1 = {
  same: 'same',
  change: 'this',
  remove: 'removed',
};

local input2 = {
  same: 'same',
  change: 'changed',
  add: 'added',
};

diff(input1, input2),
```

Output:
```json
{
  "add +": "added",
  "change ~": "~[ this , changed ]",
  "remove -": "removed"
}
```


### fn filterKubernetesObjects

```ts
filterKubernetesObjects(object, kind='')
```

`filterKubernetesObjects` implements `filterObjects` to return all Kubernetes objects in
an array, assuming that Kubernetes object are characterized by having an
`apiVersion` and `kind` field.

The `object` argument can either be an object or an array, other types will be
ignored. The `kind` allows to filter out a specific kind, if unset all kinds will
be returned.


### fn filterObjects

```ts
filterObjects(filter_func, x)
```

`filterObjects` walks a JSON tree returning all matching objects in an array.

The `x` argument can either be an object or an array, other types will be
ignored.


### fn inspect

```ts
inspect(object, maxDepth)
```

`inspect` reports the structure of a Jsonnet object with a recursion depth of
`maxDepth` (default maxDepth=10).
