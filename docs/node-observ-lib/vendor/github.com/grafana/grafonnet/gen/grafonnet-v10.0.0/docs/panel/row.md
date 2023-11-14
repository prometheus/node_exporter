# row

grafonnet.panel.row

## Index

* [`fn new(title)`](#fn-new)
* [`fn withCollapsed(value=true)`](#fn-withcollapsed)
* [`fn withDatasource(value)`](#fn-withdatasource)
* [`fn withDatasourceMixin(value)`](#fn-withdatasourcemixin)
* [`fn withGridPos(value)`](#fn-withgridpos)
* [`fn withGridPosMixin(value)`](#fn-withgridposmixin)
* [`fn withId(value)`](#fn-withid)
* [`fn withPanels(value)`](#fn-withpanels)
* [`fn withPanelsMixin(value)`](#fn-withpanelsmixin)
* [`fn withRepeat(value)`](#fn-withrepeat)
* [`fn withTitle(value)`](#fn-withtitle)
* [`fn withType()`](#fn-withtype)
* [`obj datasource`](#obj-datasource)
  * [`fn withType(value)`](#fn-datasourcewithtype)
  * [`fn withUid(value)`](#fn-datasourcewithuid)
* [`obj gridPos`](#obj-gridpos)
  * [`fn withH(value=9)`](#fn-gridposwithh)
  * [`fn withStatic(value=true)`](#fn-gridposwithstatic)
  * [`fn withW(value=12)`](#fn-gridposwithw)
  * [`fn withX(value=0)`](#fn-gridposwithx)
  * [`fn withY(value=0)`](#fn-gridposwithy)

## Fields

### fn new

```jsonnet
new(title)
```

PARAMETERS:

* **title** (`string`)

Creates a new row panel with a title.
### fn withCollapsed

```jsonnet
withCollapsed(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


### fn withDatasource

```jsonnet
withDatasource(value)
```

PARAMETERS:

* **value** (`object`)

Name of default datasource.
### fn withDatasourceMixin

```jsonnet
withDatasourceMixin(value)
```

PARAMETERS:

* **value** (`object`)

Name of default datasource.
### fn withGridPos

```jsonnet
withGridPos(value)
```

PARAMETERS:

* **value** (`object`)


### fn withGridPosMixin

```jsonnet
withGridPosMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withId

```jsonnet
withId(value)
```

PARAMETERS:

* **value** (`integer`)


### fn withPanels

```jsonnet
withPanels(value)
```

PARAMETERS:

* **value** (`array`)


### fn withPanelsMixin

```jsonnet
withPanelsMixin(value)
```

PARAMETERS:

* **value** (`array`)


### fn withRepeat

```jsonnet
withRepeat(value)
```

PARAMETERS:

* **value** (`string`)

Name of template variable to repeat for.
### fn withTitle

```jsonnet
withTitle(value)
```

PARAMETERS:

* **value** (`string`)


### fn withType

```jsonnet
withType()
```



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


### obj gridPos


#### fn gridPos.withH

```jsonnet
gridPos.withH(value=9)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `9`

Panel
#### fn gridPos.withStatic

```jsonnet
gridPos.withStatic(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

true if fixed
#### fn gridPos.withW

```jsonnet
gridPos.withW(value=12)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `12`

Panel
#### fn gridPos.withX

```jsonnet
gridPos.withX(value=0)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `0`

Panel x
#### fn gridPos.withY

```jsonnet
gridPos.withY(value=0)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `0`

Panel y