# annotation



## Index

* [`fn withDatasource(value)`](#fn-withdatasource)
* [`fn withDatasourceMixin(value)`](#fn-withdatasourcemixin)
* [`fn withEnable(value=true)`](#fn-withenable)
* [`fn withFilter(value)`](#fn-withfilter)
* [`fn withFilterMixin(value)`](#fn-withfiltermixin)
* [`fn withHide(value=true)`](#fn-withhide)
* [`fn withIconColor(value)`](#fn-withiconcolor)
* [`fn withName(value)`](#fn-withname)
* [`fn withTarget(value)`](#fn-withtarget)
* [`fn withTargetMixin(value)`](#fn-withtargetmixin)
* [`fn withType(value)`](#fn-withtype)
* [`obj datasource`](#obj-datasource)
  * [`fn withType(value)`](#fn-datasourcewithtype)
  * [`fn withUid(value)`](#fn-datasourcewithuid)
* [`obj filter`](#obj-filter)
  * [`fn withExclude(value=true)`](#fn-filterwithexclude)
  * [`fn withIds(value)`](#fn-filterwithids)
  * [`fn withIdsMixin(value)`](#fn-filterwithidsmixin)
* [`obj target`](#obj-target)
  * [`fn withLimit(value)`](#fn-targetwithlimit)
  * [`fn withMatchAny(value=true)`](#fn-targetwithmatchany)
  * [`fn withTags(value)`](#fn-targetwithtags)
  * [`fn withTagsMixin(value)`](#fn-targetwithtagsmixin)
  * [`fn withType(value)`](#fn-targetwithtype)

## Fields

### fn withDatasource

```jsonnet
withDatasource(value)
```

PARAMETERS:

* **value** (`object`)

TODO: Should be DataSourceRef
### fn withDatasourceMixin

```jsonnet
withDatasourceMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO: Should be DataSourceRef
### fn withEnable

```jsonnet
withEnable(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

When enabled the annotation query is issued with every dashboard refresh
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


### fn withHide

```jsonnet
withHide(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Annotation queries can be toggled on or off at the top of the dashboard.
When hide is true, the toggle is not shown in the dashboard.
### fn withIconColor

```jsonnet
withIconColor(value)
```

PARAMETERS:

* **value** (`string`)

Color to use for the annotation event markers
### fn withName

```jsonnet
withName(value)
```

PARAMETERS:

* **value** (`string`)

Name of annotation.
### fn withTarget

```jsonnet
withTarget(value)
```

PARAMETERS:

* **value** (`object`)

TODO: this should be a regular DataQuery that depends on the selected dashboard
these match the properties of the "grafana" datasouce that is default in most dashboards
### fn withTargetMixin

```jsonnet
withTargetMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO: this should be a regular DataQuery that depends on the selected dashboard
these match the properties of the "grafana" datasouce that is default in most dashboards
### fn withType

```jsonnet
withType(value)
```

PARAMETERS:

* **value** (`string`)

TODO -- this should not exist here, it is based on the --grafana-- datasource
### obj datasource


#### fn datasource.withType

```jsonnet
datasource.withType(value)
```

PARAMETERS:

* **value** (`string`)


#### fn datasource.withUid

```jsonnet
datasource.withUid(value)
```

PARAMETERS:

* **value** (`string`)


### obj filter


#### fn filter.withExclude

```jsonnet
filter.withExclude(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Should the specified panels be included or excluded
#### fn filter.withIds

```jsonnet
filter.withIds(value)
```

PARAMETERS:

* **value** (`array`)

Panel IDs that should be included or excluded
#### fn filter.withIdsMixin

```jsonnet
filter.withIdsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Panel IDs that should be included or excluded
### obj target


#### fn target.withLimit

```jsonnet
target.withLimit(value)
```

PARAMETERS:

* **value** (`integer`)

Only required/valid for the grafana datasource...
but code+tests is already depending on it so hard to change
#### fn target.withMatchAny

```jsonnet
target.withMatchAny(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Only required/valid for the grafana datasource...
but code+tests is already depending on it so hard to change
#### fn target.withTags

```jsonnet
target.withTags(value)
```

PARAMETERS:

* **value** (`array`)

Only required/valid for the grafana datasource...
but code+tests is already depending on it so hard to change
#### fn target.withTagsMixin

```jsonnet
target.withTagsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Only required/valid for the grafana datasource...
but code+tests is already depending on it so hard to change
#### fn target.withType

```jsonnet
target.withType(value)
```

PARAMETERS:

* **value** (`string`)

Only required/valid for the grafana datasource...
but code+tests is already depending on it so hard to change