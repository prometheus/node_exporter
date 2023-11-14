# alertList

grafonnet.panel.alertList

## Subpackages

* [panelOptions.link](panelOptions/link.md)
* [queryOptions.transformation](queryOptions/transformation.md)
* [standardOptions.mapping](standardOptions/mapping.md)
* [standardOptions.override](standardOptions/override.md)
* [standardOptions.threshold.step](standardOptions/threshold/step.md)

## Index

* [`fn new(title)`](#fn-new)
* [`obj libraryPanel`](#obj-librarypanel)
  * [`fn withName(value)`](#fn-librarypanelwithname)
  * [`fn withUid(value)`](#fn-librarypanelwithuid)
* [`obj options`](#obj-options)
  * [`fn withAlertListOptions(value)`](#fn-optionswithalertlistoptions)
  * [`fn withAlertListOptionsMixin(value)`](#fn-optionswithalertlistoptionsmixin)
  * [`fn withUnifiedAlertListOptions(value)`](#fn-optionswithunifiedalertlistoptions)
  * [`fn withUnifiedAlertListOptionsMixin(value)`](#fn-optionswithunifiedalertlistoptionsmixin)
  * [`obj AlertListOptions`](#obj-optionsalertlistoptions)
    * [`fn withAlertName(value)`](#fn-optionsalertlistoptionswithalertname)
    * [`fn withDashboardAlerts(value=true)`](#fn-optionsalertlistoptionswithdashboardalerts)
    * [`fn withDashboardTitle(value)`](#fn-optionsalertlistoptionswithdashboardtitle)
    * [`fn withFolderId(value)`](#fn-optionsalertlistoptionswithfolderid)
    * [`fn withMaxItems(value)`](#fn-optionsalertlistoptionswithmaxitems)
    * [`fn withShowOptions(value)`](#fn-optionsalertlistoptionswithshowoptions)
    * [`fn withSortOrder(value)`](#fn-optionsalertlistoptionswithsortorder)
    * [`fn withStateFilter(value)`](#fn-optionsalertlistoptionswithstatefilter)
    * [`fn withStateFilterMixin(value)`](#fn-optionsalertlistoptionswithstatefiltermixin)
    * [`fn withTags(value)`](#fn-optionsalertlistoptionswithtags)
    * [`fn withTagsMixin(value)`](#fn-optionsalertlistoptionswithtagsmixin)
    * [`obj stateFilter`](#obj-optionsalertlistoptionsstatefilter)
      * [`fn withAlerting(value=true)`](#fn-optionsalertlistoptionsstatefilterwithalerting)
      * [`fn withExecutionError(value=true)`](#fn-optionsalertlistoptionsstatefilterwithexecutionerror)
      * [`fn withNoData(value=true)`](#fn-optionsalertlistoptionsstatefilterwithnodata)
      * [`fn withOk(value=true)`](#fn-optionsalertlistoptionsstatefilterwithok)
      * [`fn withPaused(value=true)`](#fn-optionsalertlistoptionsstatefilterwithpaused)
      * [`fn withPending(value=true)`](#fn-optionsalertlistoptionsstatefilterwithpending)
  * [`obj UnifiedAlertListOptions`](#obj-optionsunifiedalertlistoptions)
    * [`fn withAlertInstanceLabelFilter(value)`](#fn-optionsunifiedalertlistoptionswithalertinstancelabelfilter)
    * [`fn withAlertName(value)`](#fn-optionsunifiedalertlistoptionswithalertname)
    * [`fn withDashboardAlerts(value=true)`](#fn-optionsunifiedalertlistoptionswithdashboardalerts)
    * [`fn withDatasource(value)`](#fn-optionsunifiedalertlistoptionswithdatasource)
    * [`fn withFolder(value)`](#fn-optionsunifiedalertlistoptionswithfolder)
    * [`fn withFolderMixin(value)`](#fn-optionsunifiedalertlistoptionswithfoldermixin)
    * [`fn withGroupBy(value)`](#fn-optionsunifiedalertlistoptionswithgroupby)
    * [`fn withGroupByMixin(value)`](#fn-optionsunifiedalertlistoptionswithgroupbymixin)
    * [`fn withGroupMode(value)`](#fn-optionsunifiedalertlistoptionswithgroupmode)
    * [`fn withMaxItems(value)`](#fn-optionsunifiedalertlistoptionswithmaxitems)
    * [`fn withShowInstances(value=true)`](#fn-optionsunifiedalertlistoptionswithshowinstances)
    * [`fn withSortOrder(value)`](#fn-optionsunifiedalertlistoptionswithsortorder)
    * [`fn withStateFilter(value)`](#fn-optionsunifiedalertlistoptionswithstatefilter)
    * [`fn withStateFilterMixin(value)`](#fn-optionsunifiedalertlistoptionswithstatefiltermixin)
    * [`fn withViewMode(value)`](#fn-optionsunifiedalertlistoptionswithviewmode)
    * [`obj folder`](#obj-optionsunifiedalertlistoptionsfolder)
      * [`fn withId(value)`](#fn-optionsunifiedalertlistoptionsfolderwithid)
      * [`fn withTitle(value)`](#fn-optionsunifiedalertlistoptionsfolderwithtitle)
    * [`obj stateFilter`](#obj-optionsunifiedalertlistoptionsstatefilter)
      * [`fn withError(value=true)`](#fn-optionsunifiedalertlistoptionsstatefilterwitherror)
      * [`fn withFiring(value=true)`](#fn-optionsunifiedalertlistoptionsstatefilterwithfiring)
      * [`fn withInactive(value=true)`](#fn-optionsunifiedalertlistoptionsstatefilterwithinactive)
      * [`fn withNoData(value=true)`](#fn-optionsunifiedalertlistoptionsstatefilterwithnodata)
      * [`fn withNormal(value=true)`](#fn-optionsunifiedalertlistoptionsstatefilterwithnormal)
      * [`fn withPending(value=true)`](#fn-optionsunifiedalertlistoptionsstatefilterwithpending)
* [`obj panelOptions`](#obj-paneloptions)
  * [`fn withDescription(value)`](#fn-paneloptionswithdescription)
  * [`fn withGridPos(h="null", w="null", x="null", y="null")`](#fn-paneloptionswithgridpos)
  * [`fn withLinks(value)`](#fn-paneloptionswithlinks)
  * [`fn withLinksMixin(value)`](#fn-paneloptionswithlinksmixin)
  * [`fn withRepeat(value)`](#fn-paneloptionswithrepeat)
  * [`fn withRepeatDirection(value="h")`](#fn-paneloptionswithrepeatdirection)
  * [`fn withTitle(value)`](#fn-paneloptionswithtitle)
  * [`fn withTransparent(value=true)`](#fn-paneloptionswithtransparent)
* [`obj queryOptions`](#obj-queryoptions)
  * [`fn withDatasource(type, uid)`](#fn-queryoptionswithdatasource)
  * [`fn withDatasourceMixin(value)`](#fn-queryoptionswithdatasourcemixin)
  * [`fn withInterval(value)`](#fn-queryoptionswithinterval)
  * [`fn withMaxDataPoints(value)`](#fn-queryoptionswithmaxdatapoints)
  * [`fn withTargets(value)`](#fn-queryoptionswithtargets)
  * [`fn withTargetsMixin(value)`](#fn-queryoptionswithtargetsmixin)
  * [`fn withTimeFrom(value)`](#fn-queryoptionswithtimefrom)
  * [`fn withTimeShift(value)`](#fn-queryoptionswithtimeshift)
  * [`fn withTransformations(value)`](#fn-queryoptionswithtransformations)
  * [`fn withTransformationsMixin(value)`](#fn-queryoptionswithtransformationsmixin)
* [`obj standardOptions`](#obj-standardoptions)
  * [`fn withDecimals(value)`](#fn-standardoptionswithdecimals)
  * [`fn withDisplayName(value)`](#fn-standardoptionswithdisplayname)
  * [`fn withFilterable(value=true)`](#fn-standardoptionswithfilterable)
  * [`fn withLinks(value)`](#fn-standardoptionswithlinks)
  * [`fn withLinksMixin(value)`](#fn-standardoptionswithlinksmixin)
  * [`fn withMappings(value)`](#fn-standardoptionswithmappings)
  * [`fn withMappingsMixin(value)`](#fn-standardoptionswithmappingsmixin)
  * [`fn withMax(value)`](#fn-standardoptionswithmax)
  * [`fn withMin(value)`](#fn-standardoptionswithmin)
  * [`fn withNoValue(value)`](#fn-standardoptionswithnovalue)
  * [`fn withOverrides(value)`](#fn-standardoptionswithoverrides)
  * [`fn withOverridesMixin(value)`](#fn-standardoptionswithoverridesmixin)
  * [`fn withPath(value)`](#fn-standardoptionswithpath)
  * [`fn withUnit(value)`](#fn-standardoptionswithunit)
  * [`obj color`](#obj-standardoptionscolor)
    * [`fn withFixedColor(value)`](#fn-standardoptionscolorwithfixedcolor)
    * [`fn withMode(value)`](#fn-standardoptionscolorwithmode)
    * [`fn withSeriesBy(value)`](#fn-standardoptionscolorwithseriesby)
  * [`obj thresholds`](#obj-standardoptionsthresholds)
    * [`fn withMode(value)`](#fn-standardoptionsthresholdswithmode)
    * [`fn withSteps(value)`](#fn-standardoptionsthresholdswithsteps)
    * [`fn withStepsMixin(value)`](#fn-standardoptionsthresholdswithstepsmixin)

## Fields

### fn new

```jsonnet
new(title)
```

PARAMETERS:

* **title** (`string`)

Creates a new alertlist panel with a title.
### obj libraryPanel


#### fn libraryPanel.withName

```jsonnet
libraryPanel.withName(value)
```

PARAMETERS:

* **value** (`string`)


#### fn libraryPanel.withUid

```jsonnet
libraryPanel.withUid(value)
```

PARAMETERS:

* **value** (`string`)


### obj options


#### fn options.withAlertListOptions

```jsonnet
options.withAlertListOptions(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withAlertListOptionsMixin

```jsonnet
options.withAlertListOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withUnifiedAlertListOptions

```jsonnet
options.withUnifiedAlertListOptions(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withUnifiedAlertListOptionsMixin

```jsonnet
options.withUnifiedAlertListOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### obj options.AlertListOptions


##### fn options.AlertListOptions.withAlertName

```jsonnet
options.AlertListOptions.withAlertName(value)
```

PARAMETERS:

* **value** (`string`)


##### fn options.AlertListOptions.withDashboardAlerts

```jsonnet
options.AlertListOptions.withDashboardAlerts(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.AlertListOptions.withDashboardTitle

```jsonnet
options.AlertListOptions.withDashboardTitle(value)
```

PARAMETERS:

* **value** (`string`)


##### fn options.AlertListOptions.withFolderId

```jsonnet
options.AlertListOptions.withFolderId(value)
```

PARAMETERS:

* **value** (`number`)


##### fn options.AlertListOptions.withMaxItems

```jsonnet
options.AlertListOptions.withMaxItems(value)
```

PARAMETERS:

* **value** (`number`)


##### fn options.AlertListOptions.withShowOptions

```jsonnet
options.AlertListOptions.withShowOptions(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"current"`, `"changes"`


##### fn options.AlertListOptions.withSortOrder

```jsonnet
options.AlertListOptions.withSortOrder(value)
```

PARAMETERS:

* **value** (`number`)
   - valid values: `1`, `2`, `3`, `4`, `5`


##### fn options.AlertListOptions.withStateFilter

```jsonnet
options.AlertListOptions.withStateFilter(value)
```

PARAMETERS:

* **value** (`object`)


##### fn options.AlertListOptions.withStateFilterMixin

```jsonnet
options.AlertListOptions.withStateFilterMixin(value)
```

PARAMETERS:

* **value** (`object`)


##### fn options.AlertListOptions.withTags

```jsonnet
options.AlertListOptions.withTags(value)
```

PARAMETERS:

* **value** (`array`)


##### fn options.AlertListOptions.withTagsMixin

```jsonnet
options.AlertListOptions.withTagsMixin(value)
```

PARAMETERS:

* **value** (`array`)


##### obj options.AlertListOptions.stateFilter


###### fn options.AlertListOptions.stateFilter.withAlerting

```jsonnet
options.AlertListOptions.stateFilter.withAlerting(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.AlertListOptions.stateFilter.withExecutionError

```jsonnet
options.AlertListOptions.stateFilter.withExecutionError(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.AlertListOptions.stateFilter.withNoData

```jsonnet
options.AlertListOptions.stateFilter.withNoData(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.AlertListOptions.stateFilter.withOk

```jsonnet
options.AlertListOptions.stateFilter.withOk(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.AlertListOptions.stateFilter.withPaused

```jsonnet
options.AlertListOptions.stateFilter.withPaused(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.AlertListOptions.stateFilter.withPending

```jsonnet
options.AlertListOptions.stateFilter.withPending(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


#### obj options.UnifiedAlertListOptions


##### fn options.UnifiedAlertListOptions.withAlertInstanceLabelFilter

```jsonnet
options.UnifiedAlertListOptions.withAlertInstanceLabelFilter(value)
```

PARAMETERS:

* **value** (`string`)


##### fn options.UnifiedAlertListOptions.withAlertName

```jsonnet
options.UnifiedAlertListOptions.withAlertName(value)
```

PARAMETERS:

* **value** (`string`)


##### fn options.UnifiedAlertListOptions.withDashboardAlerts

```jsonnet
options.UnifiedAlertListOptions.withDashboardAlerts(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.UnifiedAlertListOptions.withDatasource

```jsonnet
options.UnifiedAlertListOptions.withDatasource(value)
```

PARAMETERS:

* **value** (`string`)


##### fn options.UnifiedAlertListOptions.withFolder

```jsonnet
options.UnifiedAlertListOptions.withFolder(value)
```

PARAMETERS:

* **value** (`object`)


##### fn options.UnifiedAlertListOptions.withFolderMixin

```jsonnet
options.UnifiedAlertListOptions.withFolderMixin(value)
```

PARAMETERS:

* **value** (`object`)


##### fn options.UnifiedAlertListOptions.withGroupBy

```jsonnet
options.UnifiedAlertListOptions.withGroupBy(value)
```

PARAMETERS:

* **value** (`array`)


##### fn options.UnifiedAlertListOptions.withGroupByMixin

```jsonnet
options.UnifiedAlertListOptions.withGroupByMixin(value)
```

PARAMETERS:

* **value** (`array`)


##### fn options.UnifiedAlertListOptions.withGroupMode

```jsonnet
options.UnifiedAlertListOptions.withGroupMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"default"`, `"custom"`


##### fn options.UnifiedAlertListOptions.withMaxItems

```jsonnet
options.UnifiedAlertListOptions.withMaxItems(value)
```

PARAMETERS:

* **value** (`number`)


##### fn options.UnifiedAlertListOptions.withShowInstances

```jsonnet
options.UnifiedAlertListOptions.withShowInstances(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.UnifiedAlertListOptions.withSortOrder

```jsonnet
options.UnifiedAlertListOptions.withSortOrder(value)
```

PARAMETERS:

* **value** (`number`)
   - valid values: `1`, `2`, `3`, `4`, `5`


##### fn options.UnifiedAlertListOptions.withStateFilter

```jsonnet
options.UnifiedAlertListOptions.withStateFilter(value)
```

PARAMETERS:

* **value** (`object`)


##### fn options.UnifiedAlertListOptions.withStateFilterMixin

```jsonnet
options.UnifiedAlertListOptions.withStateFilterMixin(value)
```

PARAMETERS:

* **value** (`object`)


##### fn options.UnifiedAlertListOptions.withViewMode

```jsonnet
options.UnifiedAlertListOptions.withViewMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"list"`, `"stat"`


##### obj options.UnifiedAlertListOptions.folder


###### fn options.UnifiedAlertListOptions.folder.withId

```jsonnet
options.UnifiedAlertListOptions.folder.withId(value)
```

PARAMETERS:

* **value** (`number`)


###### fn options.UnifiedAlertListOptions.folder.withTitle

```jsonnet
options.UnifiedAlertListOptions.folder.withTitle(value)
```

PARAMETERS:

* **value** (`string`)


##### obj options.UnifiedAlertListOptions.stateFilter


###### fn options.UnifiedAlertListOptions.stateFilter.withError

```jsonnet
options.UnifiedAlertListOptions.stateFilter.withError(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.UnifiedAlertListOptions.stateFilter.withFiring

```jsonnet
options.UnifiedAlertListOptions.stateFilter.withFiring(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.UnifiedAlertListOptions.stateFilter.withInactive

```jsonnet
options.UnifiedAlertListOptions.stateFilter.withInactive(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.UnifiedAlertListOptions.stateFilter.withNoData

```jsonnet
options.UnifiedAlertListOptions.stateFilter.withNoData(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.UnifiedAlertListOptions.stateFilter.withNormal

```jsonnet
options.UnifiedAlertListOptions.stateFilter.withNormal(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.UnifiedAlertListOptions.stateFilter.withPending

```jsonnet
options.UnifiedAlertListOptions.stateFilter.withPending(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


### obj panelOptions


#### fn panelOptions.withDescription

```jsonnet
panelOptions.withDescription(value)
```

PARAMETERS:

* **value** (`string`)

Description.
#### fn panelOptions.withGridPos

```jsonnet
panelOptions.withGridPos(h="null", w="null", x="null", y="null")
```

PARAMETERS:

* **h** (`number`)
   - default value: `"null"`
* **w** (`number`)
   - default value: `"null"`
* **x** (`number`)
   - default value: `"null"`
* **y** (`number`)
   - default value: `"null"`

`withGridPos` configures the height, width and xy coordinates of the panel. Also see `grafonnet.util.grid` for helper functions to calculate these fields.

All arguments default to `null`, which means they will remain unchanged or unset.

#### fn panelOptions.withLinks

```jsonnet
panelOptions.withLinks(value)
```

PARAMETERS:

* **value** (`array`)

Panel links.
TODO fill this out - seems there are a couple variants?
#### fn panelOptions.withLinksMixin

```jsonnet
panelOptions.withLinksMixin(value)
```

PARAMETERS:

* **value** (`array`)

Panel links.
TODO fill this out - seems there are a couple variants?
#### fn panelOptions.withRepeat

```jsonnet
panelOptions.withRepeat(value)
```

PARAMETERS:

* **value** (`string`)

Name of template variable to repeat for.
#### fn panelOptions.withRepeatDirection

```jsonnet
panelOptions.withRepeatDirection(value="h")
```

PARAMETERS:

* **value** (`string`)
   - default value: `"h"`
   - valid values: `"h"`, `"v"`

Direction to repeat in if 'repeat' is set.
"h" for horizontal, "v" for vertical.
TODO this is probably optional
#### fn panelOptions.withTitle

```jsonnet
panelOptions.withTitle(value)
```

PARAMETERS:

* **value** (`string`)

Panel title.
#### fn panelOptions.withTransparent

```jsonnet
panelOptions.withTransparent(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Whether to display the panel without a background.
### obj queryOptions


#### fn queryOptions.withDatasource

```jsonnet
queryOptions.withDatasource(type, uid)
```

PARAMETERS:

* **type** (`string`)
* **uid** (`string`)

`withDatasource` sets the datasource for all queries in a panel.

The default datasource for a panel is set to 'Mixed datasource' so panels can be datasource agnostic, which is a lot more interesting from a reusability standpoint. Note that this requires query targets to explicitly set datasource for the same reason.

#### fn queryOptions.withDatasourceMixin

```jsonnet
queryOptions.withDatasourceMixin(value)
```

PARAMETERS:

* **value** (`object`)

The datasource used in all targets.
#### fn queryOptions.withInterval

```jsonnet
queryOptions.withInterval(value)
```

PARAMETERS:

* **value** (`string`)

TODO docs
TODO tighter constraint
#### fn queryOptions.withMaxDataPoints

```jsonnet
queryOptions.withMaxDataPoints(value)
```

PARAMETERS:

* **value** (`number`)

TODO docs
#### fn queryOptions.withTargets

```jsonnet
queryOptions.withTargets(value)
```

PARAMETERS:

* **value** (`array`)

TODO docs
#### fn queryOptions.withTargetsMixin

```jsonnet
queryOptions.withTargetsMixin(value)
```

PARAMETERS:

* **value** (`array`)

TODO docs
#### fn queryOptions.withTimeFrom

```jsonnet
queryOptions.withTimeFrom(value)
```

PARAMETERS:

* **value** (`string`)

TODO docs
TODO tighter constraint
#### fn queryOptions.withTimeShift

```jsonnet
queryOptions.withTimeShift(value)
```

PARAMETERS:

* **value** (`string`)

TODO docs
TODO tighter constraint
#### fn queryOptions.withTransformations

```jsonnet
queryOptions.withTransformations(value)
```

PARAMETERS:

* **value** (`array`)


#### fn queryOptions.withTransformationsMixin

```jsonnet
queryOptions.withTransformationsMixin(value)
```

PARAMETERS:

* **value** (`array`)


### obj standardOptions


#### fn standardOptions.withDecimals

```jsonnet
standardOptions.withDecimals(value)
```

PARAMETERS:

* **value** (`number`)

Significant digits (for display)
#### fn standardOptions.withDisplayName

```jsonnet
standardOptions.withDisplayName(value)
```

PARAMETERS:

* **value** (`string`)

The display value for this field.  This supports template variables blank is auto
#### fn standardOptions.withFilterable

```jsonnet
standardOptions.withFilterable(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

True if data source field supports ad-hoc filters
#### fn standardOptions.withLinks

```jsonnet
standardOptions.withLinks(value)
```

PARAMETERS:

* **value** (`array`)

The behavior when clicking on a result
#### fn standardOptions.withLinksMixin

```jsonnet
standardOptions.withLinksMixin(value)
```

PARAMETERS:

* **value** (`array`)

The behavior when clicking on a result
#### fn standardOptions.withMappings

```jsonnet
standardOptions.withMappings(value)
```

PARAMETERS:

* **value** (`array`)

Convert input values into a display string
#### fn standardOptions.withMappingsMixin

```jsonnet
standardOptions.withMappingsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Convert input values into a display string
#### fn standardOptions.withMax

```jsonnet
standardOptions.withMax(value)
```

PARAMETERS:

* **value** (`number`)


#### fn standardOptions.withMin

```jsonnet
standardOptions.withMin(value)
```

PARAMETERS:

* **value** (`number`)


#### fn standardOptions.withNoValue

```jsonnet
standardOptions.withNoValue(value)
```

PARAMETERS:

* **value** (`string`)

Alternative to empty string
#### fn standardOptions.withOverrides

```jsonnet
standardOptions.withOverrides(value)
```

PARAMETERS:

* **value** (`array`)


#### fn standardOptions.withOverridesMixin

```jsonnet
standardOptions.withOverridesMixin(value)
```

PARAMETERS:

* **value** (`array`)


#### fn standardOptions.withPath

```jsonnet
standardOptions.withPath(value)
```

PARAMETERS:

* **value** (`string`)

An explicit path to the field in the datasource.  When the frame meta includes a path,
This will default to `${frame.meta.path}/${field.name}

When defined, this value can be used as an identifier within the datasource scope, and
may be used to update the results
#### fn standardOptions.withUnit

```jsonnet
standardOptions.withUnit(value)
```

PARAMETERS:

* **value** (`string`)

Numeric Options
#### obj standardOptions.color


##### fn standardOptions.color.withFixedColor

```jsonnet
standardOptions.color.withFixedColor(value)
```

PARAMETERS:

* **value** (`string`)

Stores the fixed color value if mode is fixed
##### fn standardOptions.color.withMode

```jsonnet
standardOptions.color.withMode(value)
```

PARAMETERS:

* **value** (`string`)

The main color scheme mode
##### fn standardOptions.color.withSeriesBy

```jsonnet
standardOptions.color.withSeriesBy(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"min"`, `"max"`, `"last"`

TODO docs
#### obj standardOptions.thresholds


##### fn standardOptions.thresholds.withMode

```jsonnet
standardOptions.thresholds.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"absolute"`, `"percentage"`


##### fn standardOptions.thresholds.withSteps

```jsonnet
standardOptions.thresholds.withSteps(value)
```

PARAMETERS:

* **value** (`array`)

Must be sorted by 'value', first value is always -Infinity
##### fn standardOptions.thresholds.withStepsMixin

```jsonnet
standardOptions.thresholds.withStepsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Must be sorted by 'value', first value is always -Infinity