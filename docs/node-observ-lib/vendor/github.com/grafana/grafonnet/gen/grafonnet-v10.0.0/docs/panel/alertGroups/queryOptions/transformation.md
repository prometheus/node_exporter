# transformation



## Index

* [`fn withDisabled(value=true)`](#fn-withdisabled)
* [`fn withFilter(value)`](#fn-withfilter)
* [`fn withFilterMixin(value)`](#fn-withfiltermixin)
* [`fn withId(value)`](#fn-withid)
* [`fn withOptions(value)`](#fn-withoptions)
* [`obj filter`](#obj-filter)
  * [`fn withId(value="")`](#fn-filterwithid)
  * [`fn withOptions(value)`](#fn-filterwithoptions)

## Fields

### fn withDisabled

```jsonnet
withDisabled(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Disabled transformations are skipped
### fn withFilter

```jsonnet
withFilter(value)
```

PARAMETERS:

* **value** (`object`)


### fn withFilterMixin

```jsonnet
withFilterMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withId

```jsonnet
withId(value)
```

PARAMETERS:

* **value** (`string`)

Unique identifier of transformer
### fn withOptions

```jsonnet
withOptions(value)
```

PARAMETERS:

* **value** (`string`)

Options to be passed to the transformer
Valid options depend on the transformer id
### obj filter


#### fn filter.withId

```jsonnet
filter.withId(value="")
```

PARAMETERS:

* **value** (`string`)
   - default value: `""`


#### fn filter.withOptions

```jsonnet
filter.withOptions(value)
```

PARAMETERS:

* **value** (`string`)

