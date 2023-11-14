# link

Dashboard links are displayed at the top of the dashboard, these can either link to other dashboards or to external URLs.

The [docs](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/manage-dashboard-links/#dashboard-links) give a more comprehensive description.

Example:

```jsonnet
local g = import 'g.libsonnet';
local link = g.dashboard.link;

g.dashboard.new('Title dashboard')
+ g.dashboard.withLinks([
  link.link.new('My title', 'https://wikipedia.org/'),
])
```


## Index

* [`obj dashboards`](#obj-dashboards)
  * [`fn new(title, tags)`](#fn-dashboardsnew)
  * [`obj options`](#obj-dashboardsoptions)
    * [`fn withAsDropdown(value=true)`](#fn-dashboardsoptionswithasdropdown)
    * [`fn withIncludeVars(value=true)`](#fn-dashboardsoptionswithincludevars)
    * [`fn withKeepTime(value=true)`](#fn-dashboardsoptionswithkeeptime)
    * [`fn withTargetBlank(value=true)`](#fn-dashboardsoptionswithtargetblank)
* [`obj link`](#obj-link)
  * [`fn new(title, url)`](#fn-linknew)
  * [`fn withIcon(value)`](#fn-linkwithicon)
  * [`fn withTooltip(value)`](#fn-linkwithtooltip)
  * [`obj options`](#obj-linkoptions)
    * [`fn withAsDropdown(value=true)`](#fn-linkoptionswithasdropdown)
    * [`fn withIncludeVars(value=true)`](#fn-linkoptionswithincludevars)
    * [`fn withKeepTime(value=true)`](#fn-linkoptionswithkeeptime)
    * [`fn withTargetBlank(value=true)`](#fn-linkoptionswithtargetblank)

## Fields

### obj dashboards


#### fn dashboards.new

```jsonnet
dashboards.new(title, tags)
```

PARAMETERS:

* **title** (`string`)
* **tags** (`array`)

Create links to dashboards based on `tags`.

#### obj dashboards.options


##### fn dashboards.options.withAsDropdown

```jsonnet
dashboards.options.withAsDropdown(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn dashboards.options.withIncludeVars

```jsonnet
dashboards.options.withIncludeVars(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn dashboards.options.withKeepTime

```jsonnet
dashboards.options.withKeepTime(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn dashboards.options.withTargetBlank

```jsonnet
dashboards.options.withTargetBlank(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


### obj link


#### fn link.new

```jsonnet
link.new(title, url)
```

PARAMETERS:

* **title** (`string`)
* **url** (`string`)

Create link to an arbitrary URL.

#### fn link.withIcon

```jsonnet
link.withIcon(value)
```

PARAMETERS:

* **value** (`string`)


#### fn link.withTooltip

```jsonnet
link.withTooltip(value)
```

PARAMETERS:

* **value** (`string`)


#### obj link.options


##### fn link.options.withAsDropdown

```jsonnet
link.options.withAsDropdown(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn link.options.withIncludeVars

```jsonnet
link.options.withIncludeVars(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn link.options.withKeepTime

```jsonnet
link.options.withKeepTime(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn link.options.withTargetBlank

```jsonnet
link.options.withTargetBlank(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

