# barChart

grafonnet.panel.barChart

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
      * [`fn withAxisCenteredZero(value=true)`](#fn-fieldconfigdefaultscustomwithaxiscenteredzero)
      * [`fn withAxisColorMode(value)`](#fn-fieldconfigdefaultscustomwithaxiscolormode)
      * [`fn withAxisGridShow(value=true)`](#fn-fieldconfigdefaultscustomwithaxisgridshow)
      * [`fn withAxisLabel(value)`](#fn-fieldconfigdefaultscustomwithaxislabel)
      * [`fn withAxisPlacement(value)`](#fn-fieldconfigdefaultscustomwithaxisplacement)
      * [`fn withAxisSoftMax(value)`](#fn-fieldconfigdefaultscustomwithaxissoftmax)
      * [`fn withAxisSoftMin(value)`](#fn-fieldconfigdefaultscustomwithaxissoftmin)
      * [`fn withAxisWidth(value)`](#fn-fieldconfigdefaultscustomwithaxiswidth)
      * [`fn withFillOpacity(value=80)`](#fn-fieldconfigdefaultscustomwithfillopacity)
      * [`fn withGradientMode(value)`](#fn-fieldconfigdefaultscustomwithgradientmode)
      * [`fn withHideFrom(value)`](#fn-fieldconfigdefaultscustomwithhidefrom)
      * [`fn withHideFromMixin(value)`](#fn-fieldconfigdefaultscustomwithhidefrommixin)
      * [`fn withLineWidth(value=1)`](#fn-fieldconfigdefaultscustomwithlinewidth)
      * [`fn withScaleDistribution(value)`](#fn-fieldconfigdefaultscustomwithscaledistribution)
      * [`fn withScaleDistributionMixin(value)`](#fn-fieldconfigdefaultscustomwithscaledistributionmixin)
      * [`fn withThresholdsStyle(value)`](#fn-fieldconfigdefaultscustomwiththresholdsstyle)
      * [`fn withThresholdsStyleMixin(value)`](#fn-fieldconfigdefaultscustomwiththresholdsstylemixin)
      * [`obj hideFrom`](#obj-fieldconfigdefaultscustomhidefrom)
        * [`fn withLegend(value=true)`](#fn-fieldconfigdefaultscustomhidefromwithlegend)
        * [`fn withTooltip(value=true)`](#fn-fieldconfigdefaultscustomhidefromwithtooltip)
        * [`fn withViz(value=true)`](#fn-fieldconfigdefaultscustomhidefromwithviz)
      * [`obj scaleDistribution`](#obj-fieldconfigdefaultscustomscaledistribution)
        * [`fn withLinearThreshold(value)`](#fn-fieldconfigdefaultscustomscaledistributionwithlinearthreshold)
        * [`fn withLog(value)`](#fn-fieldconfigdefaultscustomscaledistributionwithlog)
        * [`fn withType(value)`](#fn-fieldconfigdefaultscustomscaledistributionwithtype)
      * [`obj thresholdsStyle`](#obj-fieldconfigdefaultscustomthresholdsstyle)
        * [`fn withMode(value)`](#fn-fieldconfigdefaultscustomthresholdsstylewithmode)
* [`obj libraryPanel`](#obj-librarypanel)
  * [`fn withName(value)`](#fn-librarypanelwithname)
  * [`fn withUid(value)`](#fn-librarypanelwithuid)
* [`obj options`](#obj-options)
  * [`fn withBarRadius(value=0)`](#fn-optionswithbarradius)
  * [`fn withBarWidth(value=0.97)`](#fn-optionswithbarwidth)
  * [`fn withColorByField(value)`](#fn-optionswithcolorbyfield)
  * [`fn withFullHighlight(value=true)`](#fn-optionswithfullhighlight)
  * [`fn withGroupWidth(value=0.7)`](#fn-optionswithgroupwidth)
  * [`fn withLegend(value)`](#fn-optionswithlegend)
  * [`fn withLegendMixin(value)`](#fn-optionswithlegendmixin)
  * [`fn withOrientation(value)`](#fn-optionswithorientation)
  * [`fn withShowValue(value)`](#fn-optionswithshowvalue)
  * [`fn withStacking(value)`](#fn-optionswithstacking)
  * [`fn withText(value)`](#fn-optionswithtext)
  * [`fn withTextMixin(value)`](#fn-optionswithtextmixin)
  * [`fn withTooltip(value)`](#fn-optionswithtooltip)
  * [`fn withTooltipMixin(value)`](#fn-optionswithtooltipmixin)
  * [`fn withXField(value)`](#fn-optionswithxfield)
  * [`fn withXTickLabelMaxLength(value)`](#fn-optionswithxticklabelmaxlength)
  * [`fn withXTickLabelRotation(value=0)`](#fn-optionswithxticklabelrotation)
  * [`fn withXTickLabelSpacing(value=0)`](#fn-optionswithxticklabelspacing)
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
    * [`fn withWidth(value)`](#fn-optionslegendwithwidth)
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

Creates a new barChart panel with a title.
### obj fieldConfig


#### obj fieldConfig.defaults


##### obj fieldConfig.defaults.custom


###### fn fieldConfig.defaults.custom.withAxisCenteredZero

```jsonnet
fieldConfig.defaults.custom.withAxisCenteredZero(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn fieldConfig.defaults.custom.withAxisColorMode

```jsonnet
fieldConfig.defaults.custom.withAxisColorMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"text"`, `"series"`

TODO docs
###### fn fieldConfig.defaults.custom.withAxisGridShow

```jsonnet
fieldConfig.defaults.custom.withAxisGridShow(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn fieldConfig.defaults.custom.withAxisLabel

```jsonnet
fieldConfig.defaults.custom.withAxisLabel(value)
```

PARAMETERS:

* **value** (`string`)


###### fn fieldConfig.defaults.custom.withAxisPlacement

```jsonnet
fieldConfig.defaults.custom.withAxisPlacement(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"auto"`, `"top"`, `"right"`, `"bottom"`, `"left"`, `"hidden"`

TODO docs
###### fn fieldConfig.defaults.custom.withAxisSoftMax

```jsonnet
fieldConfig.defaults.custom.withAxisSoftMax(value)
```

PARAMETERS:

* **value** (`number`)


###### fn fieldConfig.defaults.custom.withAxisSoftMin

```jsonnet
fieldConfig.defaults.custom.withAxisSoftMin(value)
```

PARAMETERS:

* **value** (`number`)


###### fn fieldConfig.defaults.custom.withAxisWidth

```jsonnet
fieldConfig.defaults.custom.withAxisWidth(value)
```

PARAMETERS:

* **value** (`number`)


###### fn fieldConfig.defaults.custom.withFillOpacity

```jsonnet
fieldConfig.defaults.custom.withFillOpacity(value=80)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `80`

Controls the fill opacity of the bars.
###### fn fieldConfig.defaults.custom.withGradientMode

```jsonnet
fieldConfig.defaults.custom.withGradientMode(value)
```

PARAMETERS:

* **value** (`string`)

Set the mode of the gradient fill. Fill gradient is based on the line color. To change the color, use the standard color scheme field option.
Gradient appearance is influenced by the Fill opacity setting.
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
###### fn fieldConfig.defaults.custom.withLineWidth

```jsonnet
fieldConfig.defaults.custom.withLineWidth(value=1)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `1`

Controls line width of the bars.
###### fn fieldConfig.defaults.custom.withScaleDistribution

```jsonnet
fieldConfig.defaults.custom.withScaleDistribution(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
###### fn fieldConfig.defaults.custom.withScaleDistributionMixin

```jsonnet
fieldConfig.defaults.custom.withScaleDistributionMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
###### fn fieldConfig.defaults.custom.withThresholdsStyle

```jsonnet
fieldConfig.defaults.custom.withThresholdsStyle(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
###### fn fieldConfig.defaults.custom.withThresholdsStyleMixin

```jsonnet
fieldConfig.defaults.custom.withThresholdsStyleMixin(value)
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


###### obj fieldConfig.defaults.custom.scaleDistribution


####### fn fieldConfig.defaults.custom.scaleDistribution.withLinearThreshold

```jsonnet
fieldConfig.defaults.custom.scaleDistribution.withLinearThreshold(value)
```

PARAMETERS:

* **value** (`number`)


####### fn fieldConfig.defaults.custom.scaleDistribution.withLog

```jsonnet
fieldConfig.defaults.custom.scaleDistribution.withLog(value)
```

PARAMETERS:

* **value** (`number`)


####### fn fieldConfig.defaults.custom.scaleDistribution.withType

```jsonnet
fieldConfig.defaults.custom.scaleDistribution.withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"linear"`, `"log"`, `"ordinal"`, `"symlog"`

TODO docs
###### obj fieldConfig.defaults.custom.thresholdsStyle


####### fn fieldConfig.defaults.custom.thresholdsStyle.withMode

```jsonnet
fieldConfig.defaults.custom.thresholdsStyle.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"off"`, `"line"`, `"dashed"`, `"area"`, `"line+area"`, `"dashed+area"`, `"series"`

TODO docs
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


#### fn options.withBarRadius

```jsonnet
options.withBarRadius(value=0)
```

PARAMETERS:

* **value** (`number`)
   - default value: `0`

Controls the radius of each bar.
#### fn options.withBarWidth

```jsonnet
options.withBarWidth(value=0.97)
```

PARAMETERS:

* **value** (`number`)
   - default value: `0.97`

Controls the width of bars. 1 = Max width, 0 = Min width.
#### fn options.withColorByField

```jsonnet
options.withColorByField(value)
```

PARAMETERS:

* **value** (`string`)

Use the color value for a sibling field to color each bar value.
#### fn options.withFullHighlight

```jsonnet
options.withFullHighlight(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Enables mode which highlights the entire bar area and shows tooltip when cursor
hovers over highlighted area
#### fn options.withGroupWidth

```jsonnet
options.withGroupWidth(value=0.7)
```

PARAMETERS:

* **value** (`number`)
   - default value: `0.7`

Controls the width of groups. 1 = max with, 0 = min width.
#### fn options.withLegend

```jsonnet
options.withLegend(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
#### fn options.withLegendMixin

```jsonnet
options.withLegendMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
#### fn options.withOrientation

```jsonnet
options.withOrientation(value)
```

PARAMETERS:

* **value** (`string`)

Controls the orientation of the bar chart, either vertical or horizontal.
#### fn options.withShowValue

```jsonnet
options.withShowValue(value)
```

PARAMETERS:

* **value** (`string`)

This controls whether values are shown on top or to the left of bars.
#### fn options.withStacking

```jsonnet
options.withStacking(value)
```

PARAMETERS:

* **value** (`string`)

Controls whether bars are stacked or not, either normally or in percent mode.
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
#### fn options.withXField

```jsonnet
options.withXField(value)
```

PARAMETERS:

* **value** (`string`)

Manually select which field from the dataset to represent the x field.
#### fn options.withXTickLabelMaxLength

```jsonnet
options.withXTickLabelMaxLength(value)
```

PARAMETERS:

* **value** (`integer`)

Sets the max length that a label can have before it is truncated.
#### fn options.withXTickLabelRotation

```jsonnet
options.withXTickLabelRotation(value=0)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `0`

Controls the rotation of the x axis labels.
#### fn options.withXTickLabelSpacing

```jsonnet
options.withXTickLabelSpacing(value=0)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `0`

Controls the spacing between x axis labels.
negative values indicate backwards skipping behavior
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


##### fn options.legend.withWidth

```jsonnet
options.legend.withWidth(value)
```

PARAMETERS:

* **value** (`number`)


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