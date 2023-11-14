# pieChart

grafonnet.panel.pieChart

## Subpackages

* [panelOptions.link](panelOptions/link.md)
* [queryOptions.transformation](queryOptions/transformation.md)
* [standardOptions.mapping](standardOptions/mapping.md)
* [standardOptions.override](standardOptions/override.md)
* [standardOptions.threshold.step](standardOptions/threshold/step.md)

## Index

* [`fn new(title)`](#fn-new)
* [`obj fieldConfig`](#obj-fieldconfig)
  * [`obj defaults`](#obj-fieldconfigdefaults)
    * [`obj custom`](#obj-fieldconfigdefaultscustom)
      * [`fn withHideFrom(value)`](#fn-fieldconfigdefaultscustomwithhidefrom)
      * [`fn withHideFromMixin(value)`](#fn-fieldconfigdefaultscustomwithhidefrommixin)
      * [`obj hideFrom`](#obj-fieldconfigdefaultscustomhidefrom)
        * [`fn withLegend(value=true)`](#fn-fieldconfigdefaultscustomhidefromwithlegend)
        * [`fn withTooltip(value=true)`](#fn-fieldconfigdefaultscustomhidefromwithtooltip)
        * [`fn withViz(value=true)`](#fn-fieldconfigdefaultscustomhidefromwithviz)
* [`obj libraryPanel`](#obj-librarypanel)
  * [`fn withName(value)`](#fn-librarypanelwithname)
  * [`fn withUid(value)`](#fn-librarypanelwithuid)
* [`obj options`](#obj-options)
  * [`fn withDisplayLabels(value)`](#fn-optionswithdisplaylabels)
  * [`fn withDisplayLabelsMixin(value)`](#fn-optionswithdisplaylabelsmixin)
  * [`fn withLegend(value)`](#fn-optionswithlegend)
  * [`fn withLegendMixin(value)`](#fn-optionswithlegendmixin)
  * [`fn withOrientation(value)`](#fn-optionswithorientation)
  * [`fn withPieType(value)`](#fn-optionswithpietype)
  * [`fn withReduceOptions(value)`](#fn-optionswithreduceoptions)
  * [`fn withReduceOptionsMixin(value)`](#fn-optionswithreduceoptionsmixin)
  * [`fn withText(value)`](#fn-optionswithtext)
  * [`fn withTextMixin(value)`](#fn-optionswithtextmixin)
  * [`fn withTooltip(value)`](#fn-optionswithtooltip)
  * [`fn withTooltipMixin(value)`](#fn-optionswithtooltipmixin)
  * [`obj legend`](#obj-optionslegend)
    * [`fn withAsTable(value=true)`](#fn-optionslegendwithastable)
    * [`fn withCalcs(value)`](#fn-optionslegendwithcalcs)
    * [`fn withCalcsMixin(value)`](#fn-optionslegendwithcalcsmixin)
    * [`fn withDisplayMode(value)`](#fn-optionslegendwithdisplaymode)
    * [`fn withIsVisible(value=true)`](#fn-optionslegendwithisvisible)
    * [`fn withPlacement(value)`](#fn-optionslegendwithplacement)
    * [`fn withShowLegend(value=true)`](#fn-optionslegendwithshowlegend)
    * [`fn withSortBy(value)`](#fn-optionslegendwithsortby)
    * [`fn withSortDesc(value=true)`](#fn-optionslegendwithsortdesc)
    * [`fn withValues(value)`](#fn-optionslegendwithvalues)
    * [`fn withValuesMixin(value)`](#fn-optionslegendwithvaluesmixin)
    * [`fn withWidth(value)`](#fn-optionslegendwithwidth)
  * [`obj reduceOptions`](#obj-optionsreduceoptions)
    * [`fn withCalcs(value)`](#fn-optionsreduceoptionswithcalcs)
    * [`fn withCalcsMixin(value)`](#fn-optionsreduceoptionswithcalcsmixin)
    * [`fn withFields(value)`](#fn-optionsreduceoptionswithfields)
    * [`fn withLimit(value)`](#fn-optionsreduceoptionswithlimit)
    * [`fn withValues(value=true)`](#fn-optionsreduceoptionswithvalues)
  * [`obj text`](#obj-optionstext)
    * [`fn withTitleSize(value)`](#fn-optionstextwithtitlesize)
    * [`fn withValueSize(value)`](#fn-optionstextwithvaluesize)
  * [`obj tooltip`](#obj-optionstooltip)
    * [`fn withMode(value)`](#fn-optionstooltipwithmode)
    * [`fn withSort(value)`](#fn-optionstooltipwithsort)
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

Creates a new pieChart panel with a title.
### obj fieldConfig


#### obj fieldConfig.defaults


##### obj fieldConfig.defaults.custom


###### fn fieldConfig.defaults.custom.withHideFrom

```jsonnet
fieldConfig.defaults.custom.withHideFrom(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
###### fn fieldConfig.defaults.custom.withHideFromMixin

```jsonnet
fieldConfig.defaults.custom.withHideFromMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
###### obj fieldConfig.defaults.custom.hideFrom


####### fn fieldConfig.defaults.custom.hideFrom.withLegend

```jsonnet
fieldConfig.defaults.custom.hideFrom.withLegend(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


####### fn fieldConfig.defaults.custom.hideFrom.withTooltip

```jsonnet
fieldConfig.defaults.custom.hideFrom.withTooltip(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


####### fn fieldConfig.defaults.custom.hideFrom.withViz

```jsonnet
fieldConfig.defaults.custom.hideFrom.withViz(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


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


#### fn options.withDisplayLabels

```jsonnet
options.withDisplayLabels(value)
```

PARAMETERS:

* **value** (`array`)


#### fn options.withDisplayLabelsMixin

```jsonnet
options.withDisplayLabelsMixin(value)
```

PARAMETERS:

* **value** (`array`)


#### fn options.withLegend

```jsonnet
options.withLegend(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withLegendMixin

```jsonnet
options.withLegendMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withOrientation

```jsonnet
options.withOrientation(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"auto"`, `"vertical"`, `"horizontal"`

TODO docs
#### fn options.withPieType

```jsonnet
options.withPieType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"pie"`, `"donut"`

Select the pie chart display style.
#### fn options.withReduceOptions

```jsonnet
options.withReduceOptions(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
#### fn options.withReduceOptionsMixin

```jsonnet
options.withReduceOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
#### fn options.withText

```jsonnet
options.withText(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
#### fn options.withTextMixin

```jsonnet
options.withTextMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
#### fn options.withTooltip

```jsonnet
options.withTooltip(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
#### fn options.withTooltipMixin

```jsonnet
options.withTooltipMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
#### obj options.legend


##### fn options.legend.withAsTable

```jsonnet
options.legend.withAsTable(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.legend.withCalcs

```jsonnet
options.legend.withCalcs(value)
```

PARAMETERS:

* **value** (`array`)


##### fn options.legend.withCalcsMixin

```jsonnet
options.legend.withCalcsMixin(value)
```

PARAMETERS:

* **value** (`array`)


##### fn options.legend.withDisplayMode

```jsonnet
options.legend.withDisplayMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"list"`, `"table"`, `"hidden"`

TODO docs
Note: "hidden" needs to remain as an option for plugins compatibility
##### fn options.legend.withIsVisible

```jsonnet
options.legend.withIsVisible(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.legend.withPlacement

```jsonnet
options.legend.withPlacement(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"bottom"`, `"right"`

TODO docs
##### fn options.legend.withShowLegend

```jsonnet
options.legend.withShowLegend(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.legend.withSortBy

```jsonnet
options.legend.withSortBy(value)
```

PARAMETERS:

* **value** (`string`)


##### fn options.legend.withSortDesc

```jsonnet
options.legend.withSortDesc(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.legend.withValues

```jsonnet
options.legend.withValues(value)
```

PARAMETERS:

* **value** (`array`)


##### fn options.legend.withValuesMixin

```jsonnet
options.legend.withValuesMixin(value)
```

PARAMETERS:

* **value** (`array`)


##### fn options.legend.withWidth

```jsonnet
options.legend.withWidth(value)
```

PARAMETERS:

* **value** (`number`)


#### obj options.reduceOptions


##### fn options.reduceOptions.withCalcs

```jsonnet
options.reduceOptions.withCalcs(value)
```

PARAMETERS:

* **value** (`array`)

When !values, pick one value for the whole field
##### fn options.reduceOptions.withCalcsMixin

```jsonnet
options.reduceOptions.withCalcsMixin(value)
```

PARAMETERS:

* **value** (`array`)

When !values, pick one value for the whole field
##### fn options.reduceOptions.withFields

```jsonnet
options.reduceOptions.withFields(value)
```

PARAMETERS:

* **value** (`string`)

Which fields to show.  By default this is only numeric fields
##### fn options.reduceOptions.withLimit

```jsonnet
options.reduceOptions.withLimit(value)
```

PARAMETERS:

* **value** (`number`)

if showing all values limit
##### fn options.reduceOptions.withValues

```jsonnet
options.reduceOptions.withValues(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

If true show each row value
#### obj options.text


##### fn options.text.withTitleSize

```jsonnet
options.text.withTitleSize(value)
```

PARAMETERS:

* **value** (`number`)

Explicit title text size
##### fn options.text.withValueSize

```jsonnet
options.text.withValueSize(value)
```

PARAMETERS:

* **value** (`number`)

Explicit value text size
#### obj options.tooltip


##### fn options.tooltip.withMode

```jsonnet
options.tooltip.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"single"`, `"multi"`, `"none"`

TODO docs
##### fn options.tooltip.withSort

```jsonnet
options.tooltip.withSort(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"asc"`, `"desc"`, `"none"`

TODO docs
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