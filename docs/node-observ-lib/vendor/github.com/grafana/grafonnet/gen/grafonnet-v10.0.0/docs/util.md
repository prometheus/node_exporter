# util

Helper functions that work well with Grafonnet.

## Index

* [`obj dashboard`](#obj-dashboard)
  * [`fn getOptionsForCustomQuery(query)`](#fn-dashboardgetoptionsforcustomquery)
* [`obj grid`](#obj-grid)
  * [`fn makeGrid(panels, panelWidth, panelHeight, startY)`](#fn-gridmakegrid)
  * [`fn wrapPanels(panels, panelWidth, panelHeight, startY)`](#fn-gridwrappanels)
* [`obj panel`](#obj-panel)
  * [`fn setPanelIDs(panels)`](#fn-panelsetpanelids)
* [`obj string`](#obj-string)
  * [`fn slugify(string)`](#fn-stringslugify)

## Fields

### obj dashboard


#### fn dashboard.getOptionsForCustomQuery

```jsonnet
dashboard.getOptionsForCustomQuery(query)
```

PARAMETERS:

* **query** (`string`)

`getOptionsForCustomQuery` provides values for the `options` and `current` fields.
These are required for template variables of type 'custom'but do not automatically
get populated by Grafana when importing a dashboard from JSON.

This is a bit of a hack and should always be called on functions that set `type` on
a template variable. Ideally Grafana populates these fields from the `query` value
but this provides a backwards compatible solution.

### obj grid


#### fn grid.makeGrid

```jsonnet
grid.makeGrid(panels, panelWidth, panelHeight, startY)
```

PARAMETERS:

* **panels** (`array`)
* **panelWidth** (`number`)
* **panelHeight** (`number`)
* **startY** (`number`)

`makeGrid` returns an array of `panels` organized in a grid with equal `panelWidth`
and `panelHeight`. Row panels are used as "linebreaks", if a Row panel is collapsed,
then all panels below it will be folded into the row.

This function will use the full grid of 24 columns, setting `panelWidth` to a value
that can divide 24 into equal parts will fill up the page nicely. (1, 2, 3, 4, 6, 8, 12)
Other value for `panelWidth` will leave a gap on the far right.

Optional `startY` can be provided to place generated grid above or below existing panels.

#### fn grid.wrapPanels

```jsonnet
grid.wrapPanels(panels, panelWidth, panelHeight, startY)
```

PARAMETERS:

* **panels** (`array`)
* **panelWidth** (`number`)
* **panelHeight** (`number`)
* **startY** (`number`)

`wrapPanels` returns an array of `panels` organized in a grid, wrapping up to next 'row' if total width exceeds full grid of 24 columns.
'panelHeight' and 'panelWidth' are used unless panels already have height and width defined.

### obj panel


#### fn panel.setPanelIDs

```jsonnet
panel.setPanelIDs(panels)
```

PARAMETERS:

* **panels** (`array`)

`setPanelIDs` ensures that all `panels` have a unique ID, this functions is used in
`dashboard.withPanels` and `dashboard.withPanelsMixin` to provide a consistent
experience.

used in ../dashboard.libsonnet

### obj string


#### fn string.slugify

```jsonnet
string.slugify(string)
```

PARAMETERS:

* **string** (`string`)

`slugify` will create a simple slug from `string`, keeping only alphanumeric
characters and replacing spaces with dashes.
