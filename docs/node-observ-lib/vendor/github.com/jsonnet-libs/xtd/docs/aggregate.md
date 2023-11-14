---
permalink: /aggregate/
---

# aggregate

```jsonnet
local aggregate = import "github.com/jsonnet-libs/xtd/aggregate.libsonnet"
```

`aggregate` implements helper functions to aggregate arrays of objects into objects with arrays.

Example:

```jsonnet
local apps = [
  {
    appid: 'id1',
    name: 'yo',
    id: i,
  }
  for i in std.range(0, 10)
];

aggregate.byKeys(apps, ['appid', 'name']);
```

Output:

```json
{
   "id1": {
      "yo": [
         {
            "appid": "id1",
            "id": 0,
            "name": "yo"
         },
         {
            "appid": "id1",
            "id": 1,
            "name": "yo"
         },
         ...
      ]
   }
}
```


## Index

* [`fn byKey(arr, key)`](#fn-bykey)
* [`fn byKeys(arr, keys)`](#fn-bykeys)

## Fields

### fn byKey

```ts
byKey(arr, key)
```

`byKey` aggregates an array by the value of `key`


### fn byKeys

```ts
byKeys(arr, keys)
```

`byKey` aggregates an array by iterating over `keys`, each item in `keys` nests the
aggregate one layer deeper.
