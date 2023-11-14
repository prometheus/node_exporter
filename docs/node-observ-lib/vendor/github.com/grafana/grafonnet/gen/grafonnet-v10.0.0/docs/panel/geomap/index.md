# geomap

grafonnet.panel.geomap

## Subpackages

* [options.layers](options/layers.md)
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
  * [`fn withBasemap(value)`](#fn-optionswithbasemap)
  * [`fn withBasemapMixin(value)`](#fn-optionswithbasemapmixin)
  * [`fn withControls(value)`](#fn-optionswithcontrols)
  * [`fn withControlsMixin(value)`](#fn-optionswithcontrolsmixin)
  * [`fn withLayers(value)`](#fn-optionswithlayers)
  * [`fn withLayersMixin(value)`](#fn-optionswithlayersmixin)
  * [`fn withTooltip(value)`](#fn-optionswithtooltip)
  * [`fn withTooltipMixin(value)`](#fn-optionswithtooltipmixin)
  * [`fn withView(value)`](#fn-optionswithview)
  * [`fn withViewMixin(value)`](#fn-optionswithviewmixin)
  * [`obj basemap`](#obj-optionsbasemap)
    * [`fn withConfig(value)`](#fn-optionsbasemapwithconfig)
    * [`fn withFilterData(value)`](#fn-optionsbasemapwithfilterdata)
    * [`fn withLocation(value)`](#fn-optionsbasemapwithlocation)
    * [`fn withLocationMixin(value)`](#fn-optionsbasemapwithlocationmixin)
    * [`fn withName(value)`](#fn-optionsbasemapwithname)
    * [`fn withOpacity(value)`](#fn-optionsbasemapwithopacity)
    * [`fn withTooltip(value=true)`](#fn-optionsbasemapwithtooltip)
    * [`fn withType(value)`](#fn-optionsbasemapwithtype)
    * [`obj location`](#obj-optionsbasemaplocation)
      * [`fn withGazetteer(value)`](#fn-optionsbasemaplocationwithgazetteer)
      * [`fn withGeohash(value)`](#fn-optionsbasemaplocationwithgeohash)
      * [`fn withLatitude(value)`](#fn-optionsbasemaplocationwithlatitude)
      * [`fn withLongitude(value)`](#fn-optionsbasemaplocationwithlongitude)
      * [`fn withLookup(value)`](#fn-optionsbasemaplocationwithlookup)
      * [`fn withMode(value)`](#fn-optionsbasemaplocationwithmode)
      * [`fn withWkt(value)`](#fn-optionsbasemaplocationwithwkt)
  * [`obj controls`](#obj-optionscontrols)
    * [`fn withMouseWheelZoom(value=true)`](#fn-optionscontrolswithmousewheelzoom)
    * [`fn withShowAttribution(value=true)`](#fn-optionscontrolswithshowattribution)
    * [`fn withShowDebug(value=true)`](#fn-optionscontrolswithshowdebug)
    * [`fn withShowMeasure(value=true)`](#fn-optionscontrolswithshowmeasure)
    * [`fn withShowScale(value=true)`](#fn-optionscontrolswithshowscale)
    * [`fn withShowZoom(value=true)`](#fn-optionscontrolswithshowzoom)
  * [`obj tooltip`](#obj-optionstooltip)
    * [`fn withMode(value)`](#fn-optionstooltipwithmode)
  * [`obj view`](#obj-optionsview)
    * [`fn withAllLayers(value=true)`](#fn-optionsviewwithalllayers)
    * [`fn withId(value="zero")`](#fn-optionsviewwithid)
    * [`fn withLastOnly(value=true)`](#fn-optionsviewwithlastonly)
    * [`fn withLat(value=0)`](#fn-optionsviewwithlat)
    * [`fn withLayer(value)`](#fn-optionsviewwithlayer)
    * [`fn withLon(value=0)`](#fn-optionsviewwithlon)
    * [`fn withMaxZoom(value)`](#fn-optionsviewwithmaxzoom)
    * [`fn withMinZoom(value)`](#fn-optionsviewwithminzoom)
    * [`fn withPadding(value)`](#fn-optionsviewwithpadding)
    * [`fn withShared(value=true)`](#fn-optionsviewwithshared)
    * [`fn withZoom(value=1)`](#fn-optionsviewwithzoom)
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

Creates a new geomap panel with a title.
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


#### fn options.withBasemap

```jsonnet
options.withBasemap(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withBasemapMixin

```jsonnet
options.withBasemapMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withControls

```jsonnet
options.withControls(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withControlsMixin

```jsonnet
options.withControlsMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withLayers

```jsonnet
options.withLayers(value)
```

PARAMETERS:

* **value** (`array`)


#### fn options.withLayersMixin

```jsonnet
options.withLayersMixin(value)
```

PARAMETERS:

* **value** (`array`)


#### fn options.withTooltip

```jsonnet
options.withTooltip(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withTooltipMixin

```jsonnet
options.withTooltipMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withView

```jsonnet
options.withView(value)
```

PARAMETERS:

* **value** (`object`)


#### fn options.withViewMixin

```jsonnet
options.withViewMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### obj options.basemap


##### fn options.basemap.withConfig

```jsonnet
options.basemap.withConfig(value)
```

PARAMETERS:

* **value** (`string`)

Custom options depending on the type
##### fn options.basemap.withFilterData

```jsonnet
options.basemap.withFilterData(value)
```

PARAMETERS:

* **value** (`string`)

Defines a frame MatcherConfig that may filter data for the given layer
##### fn options.basemap.withLocation

```jsonnet
options.basemap.withLocation(value)
```

PARAMETERS:

* **value** (`object`)


##### fn options.basemap.withLocationMixin

```jsonnet
options.basemap.withLocationMixin(value)
```

PARAMETERS:

* **value** (`object`)


##### fn options.basemap.withName

```jsonnet
options.basemap.withName(value)
```

PARAMETERS:

* **value** (`string`)

configured unique display name
##### fn options.basemap.withOpacity

```jsonnet
options.basemap.withOpacity(value)
```

PARAMETERS:

* **value** (`integer`)

Common properties:
https://openlayers.org/en/latest/apidoc/module-ol_layer_Base-BaseLayer.html
Layer opacity (0-1)
##### fn options.basemap.withTooltip

```jsonnet
options.basemap.withTooltip(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Check tooltip (defaults to true)
##### fn options.basemap.withType

```jsonnet
options.basemap.withType(value)
```

PARAMETERS:

* **value** (`string`)


##### obj options.basemap.location


###### fn options.basemap.location.withGazetteer

```jsonnet
options.basemap.location.withGazetteer(value)
```

PARAMETERS:

* **value** (`string`)

Path to Gazetteer
###### fn options.basemap.location.withGeohash

```jsonnet
options.basemap.location.withGeohash(value)
```

PARAMETERS:

* **value** (`string`)

Field mappings
###### fn options.basemap.location.withLatitude

```jsonnet
options.basemap.location.withLatitude(value)
```

PARAMETERS:

* **value** (`string`)


###### fn options.basemap.location.withLongitude

```jsonnet
options.basemap.location.withLongitude(value)
```

PARAMETERS:

* **value** (`string`)


###### fn options.basemap.location.withLookup

```jsonnet
options.basemap.location.withLookup(value)
```

PARAMETERS:

* **value** (`string`)


###### fn options.basemap.location.withMode

```jsonnet
options.basemap.location.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"auto"`, `"geohash"`, `"coords"`, `"lookup"`


###### fn options.basemap.location.withWkt

```jsonnet
options.basemap.location.withWkt(value)
```

PARAMETERS:

* **value** (`string`)


#### obj options.controls


##### fn options.controls.withMouseWheelZoom

```jsonnet
options.controls.withMouseWheelZoom(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

let the mouse wheel zoom
##### fn options.controls.withShowAttribution

```jsonnet
options.controls.withShowAttribution(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Lower right
##### fn options.controls.withShowDebug

```jsonnet
options.controls.withShowDebug(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Show debug
##### fn options.controls.withShowMeasure

```jsonnet
options.controls.withShowMeasure(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Show measure
##### fn options.controls.withShowScale

```jsonnet
options.controls.withShowScale(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Scale options
##### fn options.controls.withShowZoom

```jsonnet
options.controls.withShowZoom(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Zoom (upper left)
#### obj options.tooltip


##### fn options.tooltip.withMode

```jsonnet
options.tooltip.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"none"`, `"details"`


#### obj options.view


##### fn options.view.withAllLayers

```jsonnet
options.view.withAllLayers(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.view.withId

```jsonnet
options.view.withId(value="zero")
```

PARAMETERS:

* **value** (`string`)
   - default value: `"zero"`


##### fn options.view.withLastOnly

```jsonnet
options.view.withLastOnly(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.view.withLat

```jsonnet
options.view.withLat(value=0)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `0`


##### fn options.view.withLayer

```jsonnet
options.view.withLayer(value)
```

PARAMETERS:

* **value** (`string`)


##### fn options.view.withLon

```jsonnet
options.view.withLon(value=0)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `0`


##### fn options.view.withMaxZoom

```jsonnet
options.view.withMaxZoom(value)
```

PARAMETERS:

* **value** (`integer`)


##### fn options.view.withMinZoom

```jsonnet
options.view.withMinZoom(value)
```

PARAMETERS:

* **value** (`integer`)


##### fn options.view.withPadding

```jsonnet
options.view.withPadding(value)
```

PARAMETERS:

* **value** (`integer`)


##### fn options.view.withShared

```jsonnet
options.view.withShared(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


##### fn options.view.withZoom

```jsonnet
options.view.withZoom(value=1)
```

PARAMETERS:

* **value** (`integer`)
   - default value: `1`


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