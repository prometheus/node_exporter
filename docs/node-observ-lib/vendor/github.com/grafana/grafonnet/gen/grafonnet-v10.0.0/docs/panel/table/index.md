# table

grafonnet.panel.table

## Subpackages

* [options.sortBy](options/sortBy.md)
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
      * [`fn withAlign(value)`](#fn-fieldconfigdefaultscustomwithalign)
      * [`fn withCellOptions(value)`](#fn-fieldconfigdefaultscustomwithcelloptions)
      * [`fn withCellOptionsMixin(value)`](#fn-fieldconfigdefaultscustomwithcelloptionsmixin)
      * [`fn withDisplayMode(value)`](#fn-fieldconfigdefaultscustomwithdisplaymode)
      * [`fn withFilterable(value=true)`](#fn-fieldconfigdefaultscustomwithfilterable)
      * [`fn withHidden(value=true)`](#fn-fieldconfigdefaultscustomwithhidden)
      * [`fn withHideHeader(value=true)`](#fn-fieldconfigdefaultscustomwithhideheader)
      * [`fn withInspect(value=true)`](#fn-fieldconfigdefaultscustomwithinspect)
      * [`fn withMinWidth(value)`](#fn-fieldconfigdefaultscustomwithminwidth)
      * [`fn withWidth(value)`](#fn-fieldconfigdefaultscustomwithwidth)
      * [`obj cellOptions`](#obj-fieldconfigdefaultscustomcelloptions)
        * [`fn withTableAutoCellOptions(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtableautocelloptions)
        * [`fn withTableAutoCellOptionsMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtableautocelloptionsmixin)
        * [`fn withTableBarGaugeCellOptions(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablebargaugecelloptions)
        * [`fn withTableBarGaugeCellOptionsMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablebargaugecelloptionsmixin)
        * [`fn withTableColorTextCellOptions(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablecolortextcelloptions)
        * [`fn withTableColorTextCellOptionsMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablecolortextcelloptionsmixin)
        * [`fn withTableColoredBackgroundCellOptions(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablecoloredbackgroundcelloptions)
        * [`fn withTableColoredBackgroundCellOptionsMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablecoloredbackgroundcelloptionsmixin)
        * [`fn withTableImageCellOptions(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtableimagecelloptions)
        * [`fn withTableImageCellOptionsMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtableimagecelloptionsmixin)
        * [`fn withTableJsonViewCellOptions(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablejsonviewcelloptions)
        * [`fn withTableJsonViewCellOptionsMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablejsonviewcelloptionsmixin)
        * [`fn withTableSparklineCellOptions(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablesparklinecelloptions)
        * [`fn withTableSparklineCellOptionsMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionswithtablesparklinecelloptionsmixin)
        * [`obj TableAutoCellOptions`](#obj-fieldconfigdefaultscustomcelloptionstableautocelloptions)
          * [`fn withType()`](#fn-fieldconfigdefaultscustomcelloptionstableautocelloptionswithtype)
        * [`obj TableBarGaugeCellOptions`](#obj-fieldconfigdefaultscustomcelloptionstablebargaugecelloptions)
          * [`fn withMode(value)`](#fn-fieldconfigdefaultscustomcelloptionstablebargaugecelloptionswithmode)
          * [`fn withType()`](#fn-fieldconfigdefaultscustomcelloptionstablebargaugecelloptionswithtype)
          * [`fn withValueDisplayMode(value)`](#fn-fieldconfigdefaultscustomcelloptionstablebargaugecelloptionswithvaluedisplaymode)
        * [`obj TableColorTextCellOptions`](#obj-fieldconfigdefaultscustomcelloptionstablecolortextcelloptions)
          * [`fn withType()`](#fn-fieldconfigdefaultscustomcelloptionstablecolortextcelloptionswithtype)
        * [`obj TableColoredBackgroundCellOptions`](#obj-fieldconfigdefaultscustomcelloptionstablecoloredbackgroundcelloptions)
          * [`fn withMode(value)`](#fn-fieldconfigdefaultscustomcelloptionstablecoloredbackgroundcelloptionswithmode)
          * [`fn withType()`](#fn-fieldconfigdefaultscustomcelloptionstablecoloredbackgroundcelloptionswithtype)
        * [`obj TableImageCellOptions`](#obj-fieldconfigdefaultscustomcelloptionstableimagecelloptions)
          * [`fn withType()`](#fn-fieldconfigdefaultscustomcelloptionstableimagecelloptionswithtype)
        * [`obj TableJsonViewCellOptions`](#obj-fieldconfigdefaultscustomcelloptionstablejsonviewcelloptions)
          * [`fn withType()`](#fn-fieldconfigdefaultscustomcelloptionstablejsonviewcelloptionswithtype)
        * [`obj TableSparklineCellOptions`](#obj-fieldconfigdefaultscustomcelloptionstablesparklinecelloptions)
          * [`fn withAxisCenteredZero(value=true)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithaxiscenteredzero)
          * [`fn withAxisColorMode(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithaxiscolormode)
          * [`fn withAxisGridShow(value=true)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithaxisgridshow)
          * [`fn withAxisLabel(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithaxislabel)
          * [`fn withAxisPlacement(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithaxisplacement)
          * [`fn withAxisSoftMax(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithaxissoftmax)
          * [`fn withAxisSoftMin(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithaxissoftmin)
          * [`fn withAxisWidth(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithaxiswidth)
          * [`fn withBarAlignment(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithbaralignment)
          * [`fn withBarMaxWidth(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithbarmaxwidth)
          * [`fn withBarWidthFactor(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithbarwidthfactor)
          * [`fn withDrawStyle(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithdrawstyle)
          * [`fn withFillBelowTo(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithfillbelowto)
          * [`fn withFillColor(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithfillcolor)
          * [`fn withFillOpacity(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithfillopacity)
          * [`fn withGradientMode(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithgradientmode)
          * [`fn withHideFrom(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithhidefrom)
          * [`fn withHideFromMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithhidefrommixin)
          * [`fn withLineColor(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithlinecolor)
          * [`fn withLineInterpolation(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithlineinterpolation)
          * [`fn withLineStyle(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithlinestyle)
          * [`fn withLineStyleMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithlinestylemixin)
          * [`fn withLineWidth(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithlinewidth)
          * [`fn withPointColor(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithpointcolor)
          * [`fn withPointSize(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithpointsize)
          * [`fn withPointSymbol(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithpointsymbol)
          * [`fn withScaleDistribution(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithscaledistribution)
          * [`fn withScaleDistributionMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithscaledistributionmixin)
          * [`fn withShowPoints(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithshowpoints)
          * [`fn withSpanNulls(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithspannulls)
          * [`fn withStacking(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithstacking)
          * [`fn withStackingMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithstackingmixin)
          * [`fn withThresholdsStyle(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswiththresholdsstyle)
          * [`fn withThresholdsStyleMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswiththresholdsstylemixin)
          * [`fn withTransform(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithtransform)
          * [`fn withType()`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionswithtype)
          * [`obj hideFrom`](#obj-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionshidefrom)
            * [`fn withLegend(value=true)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionshidefromwithlegend)
            * [`fn withTooltip(value=true)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionshidefromwithtooltip)
            * [`fn withViz(value=true)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionshidefromwithviz)
          * [`obj lineStyle`](#obj-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionslinestyle)
            * [`fn withDash(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionslinestylewithdash)
            * [`fn withDashMixin(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionslinestylewithdashmixin)
            * [`fn withFill(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionslinestylewithfill)
          * [`obj scaleDistribution`](#obj-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsscaledistribution)
            * [`fn withLinearThreshold(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsscaledistributionwithlinearthreshold)
            * [`fn withLog(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsscaledistributionwithlog)
            * [`fn withType(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsscaledistributionwithtype)
          * [`obj stacking`](#obj-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsstacking)
            * [`fn withGroup(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsstackingwithgroup)
            * [`fn withMode(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsstackingwithmode)
          * [`obj thresholdsStyle`](#obj-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsthresholdsstyle)
            * [`fn withMode(value)`](#fn-fieldconfigdefaultscustomcelloptionstablesparklinecelloptionsthresholdsstylewithmode)
* [`obj libraryPanel`](#obj-librarypanel)
  * [`fn withName(value)`](#fn-librarypanelwithname)
  * [`fn withUid(value)`](#fn-librarypanelwithuid)
* [`obj options`](#obj-options)
  * [`fn withCellHeight(value)`](#fn-optionswithcellheight)
  * [`fn withFooter(value={"countRows": false,"reducer": [],"show": false})`](#fn-optionswithfooter)
  * [`fn withFooterMixin(value={"countRows": false,"reducer": [],"show": false})`](#fn-optionswithfootermixin)
  * [`fn withFrameIndex(value=0)`](#fn-optionswithframeindex)
  * [`fn withShowHeader(value=true)`](#fn-optionswithshowheader)
  * [`fn withShowTypeIcons(value=true)`](#fn-optionswithshowtypeicons)
  * [`fn withSortBy(value)`](#fn-optionswithsortby)
  * [`fn withSortByMixin(value)`](#fn-optionswithsortbymixin)
  * [`obj footer`](#obj-optionsfooter)
    * [`fn withTableFooterOptions(value)`](#fn-optionsfooterwithtablefooteroptions)
    * [`fn withTableFooterOptionsMixin(value)`](#fn-optionsfooterwithtablefooteroptionsmixin)
    * [`obj TableFooterOptions`](#obj-optionsfootertablefooteroptions)
      * [`fn withCountRows(value=true)`](#fn-optionsfootertablefooteroptionswithcountrows)
      * [`fn withEnablePagination(value=true)`](#fn-optionsfootertablefooteroptionswithenablepagination)
      * [`fn withFields(value)`](#fn-optionsfootertablefooteroptionswithfields)
      * [`fn withFieldsMixin(value)`](#fn-optionsfootertablefooteroptionswithfieldsmixin)
      * [`fn withReducer(value)`](#fn-optionsfootertablefooteroptionswithreducer)
      * [`fn withReducerMixin(value)`](#fn-optionsfootertablefooteroptionswithreducermixin)
      * [`fn withShow(value=true)`](#fn-optionsfootertablefooteroptionswithshow)
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

Creates a new table panel with a title.
### obj fieldConfig


#### obj fieldConfig.defaults


##### obj fieldConfig.defaults.custom


###### fn fieldConfig.defaults.custom.withAlign

```jsonnet
fieldConfig.defaults.custom.withAlign(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"auto"`, `"left"`, `"right"`, `"center"`

TODO -- should not be table specific! TODO docs
###### fn fieldConfig.defaults.custom.withCellOptions

```jsonnet
fieldConfig.defaults.custom.withCellOptions(value)
```

PARAMETERS:

* **value** (`string`)

Table cell options. Each cell has a display mode and other potential options for that display.
###### fn fieldConfig.defaults.custom.withCellOptionsMixin

```jsonnet
fieldConfig.defaults.custom.withCellOptionsMixin(value)
```

PARAMETERS:

* **value** (`string`)

Table cell options. Each cell has a display mode and other potential options for that display.
###### fn fieldConfig.defaults.custom.withDisplayMode

```jsonnet
fieldConfig.defaults.custom.withDisplayMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"auto"`, `"basic"`, `"color-background"`, `"color-background-solid"`, `"color-text"`, `"custom"`, `"gauge"`, `"gradient-gauge"`, `"image"`, `"json-view"`, `"lcd-gauge"`, `"sparkline"`

Internally, this is the "type" of cell that's being displayed in the table such as colored text, JSON, gauge, etc. The color-background-solid, gradient-gauge, and lcd-gauge modes are deprecated in favor of new cell subOptions
###### fn fieldConfig.defaults.custom.withFilterable

```jsonnet
fieldConfig.defaults.custom.withFilterable(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn fieldConfig.defaults.custom.withHidden

```jsonnet
fieldConfig.defaults.custom.withHidden(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn fieldConfig.defaults.custom.withHideHeader

```jsonnet
fieldConfig.defaults.custom.withHideHeader(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Hides any header for a column, usefull for columns that show some static content or buttons.
###### fn fieldConfig.defaults.custom.withInspect

```jsonnet
fieldConfig.defaults.custom.withInspect(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn fieldConfig.defaults.custom.withMinWidth

```jsonnet
fieldConfig.defaults.custom.withMinWidth(value)
```

PARAMETERS:

* **value** (`number`)


###### fn fieldConfig.defaults.custom.withWidth

```jsonnet
fieldConfig.defaults.custom.withWidth(value)
```

PARAMETERS:

* **value** (`number`)


###### obj fieldConfig.defaults.custom.cellOptions


####### fn fieldConfig.defaults.custom.cellOptions.withTableAutoCellOptions

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableAutoCellOptions(value)
```

PARAMETERS:

* **value** (`object`)

Auto mode table cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableAutoCellOptionsMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableAutoCellOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Auto mode table cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableBarGaugeCellOptions

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableBarGaugeCellOptions(value)
```

PARAMETERS:

* **value** (`object`)

Gauge cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableBarGaugeCellOptionsMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableBarGaugeCellOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Gauge cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableColorTextCellOptions

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableColorTextCellOptions(value)
```

PARAMETERS:

* **value** (`object`)

Colored text cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableColorTextCellOptionsMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableColorTextCellOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Colored text cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableColoredBackgroundCellOptions

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableColoredBackgroundCellOptions(value)
```

PARAMETERS:

* **value** (`object`)

Colored background cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableColoredBackgroundCellOptionsMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableColoredBackgroundCellOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Colored background cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableImageCellOptions

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableImageCellOptions(value)
```

PARAMETERS:

* **value** (`object`)

Json view cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableImageCellOptionsMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableImageCellOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Json view cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableJsonViewCellOptions

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableJsonViewCellOptions(value)
```

PARAMETERS:

* **value** (`object`)

Json view cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableJsonViewCellOptionsMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableJsonViewCellOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Json view cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableSparklineCellOptions

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableSparklineCellOptions(value)
```

PARAMETERS:

* **value** (`object`)

Sparkline cell options
####### fn fieldConfig.defaults.custom.cellOptions.withTableSparklineCellOptionsMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.withTableSparklineCellOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Sparkline cell options
####### obj fieldConfig.defaults.custom.cellOptions.TableAutoCellOptions


######## fn fieldConfig.defaults.custom.cellOptions.TableAutoCellOptions.withType

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableAutoCellOptions.withType()
```



####### obj fieldConfig.defaults.custom.cellOptions.TableBarGaugeCellOptions


######## fn fieldConfig.defaults.custom.cellOptions.TableBarGaugeCellOptions.withMode

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableBarGaugeCellOptions.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"basic"`, `"gradient"`, `"lcd"`

Enum expressing the possible display modes for the bar gauge component of Grafana UI
######## fn fieldConfig.defaults.custom.cellOptions.TableBarGaugeCellOptions.withType

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableBarGaugeCellOptions.withType()
```



######## fn fieldConfig.defaults.custom.cellOptions.TableBarGaugeCellOptions.withValueDisplayMode

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableBarGaugeCellOptions.withValueDisplayMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"color"`, `"hidden"`, `"text"`

Allows for the table cell gauge display type to set the gauge mode.
####### obj fieldConfig.defaults.custom.cellOptions.TableColorTextCellOptions


######## fn fieldConfig.defaults.custom.cellOptions.TableColorTextCellOptions.withType

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableColorTextCellOptions.withType()
```



####### obj fieldConfig.defaults.custom.cellOptions.TableColoredBackgroundCellOptions


######## fn fieldConfig.defaults.custom.cellOptions.TableColoredBackgroundCellOptions.withMode

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableColoredBackgroundCellOptions.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"basic"`, `"gradient"`

Display mode to the "Colored Background" display mode for table cells. Either displays a solid color (basic mode) or a gradient.
######## fn fieldConfig.defaults.custom.cellOptions.TableColoredBackgroundCellOptions.withType

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableColoredBackgroundCellOptions.withType()
```



####### obj fieldConfig.defaults.custom.cellOptions.TableImageCellOptions


######## fn fieldConfig.defaults.custom.cellOptions.TableImageCellOptions.withType

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableImageCellOptions.withType()
```



####### obj fieldConfig.defaults.custom.cellOptions.TableJsonViewCellOptions


######## fn fieldConfig.defaults.custom.cellOptions.TableJsonViewCellOptions.withType

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableJsonViewCellOptions.withType()
```



####### obj fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisCenteredZero

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisCenteredZero(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisColorMode

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisColorMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"series"`, `"text"`

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisGridShow

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisGridShow(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisLabel

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisLabel(value)
```

PARAMETERS:

* **value** (`string`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisPlacement

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisPlacement(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"auto"`, `"bottom"`, `"hidden"`, `"left"`, `"right"`, `"top"`

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisSoftMax

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisSoftMax(value)
```

PARAMETERS:

* **value** (`number`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisSoftMin

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisSoftMin(value)
```

PARAMETERS:

* **value** (`number`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisWidth

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withAxisWidth(value)
```

PARAMETERS:

* **value** (`number`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withBarAlignment

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withBarAlignment(value)
```

PARAMETERS:

* **value** (`number`)
   - valid values: `1`, `-1`, `0`

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withBarMaxWidth

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withBarMaxWidth(value)
```

PARAMETERS:

* **value** (`number`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withBarWidthFactor

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withBarWidthFactor(value)
```

PARAMETERS:

* **value** (`number`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withDrawStyle

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withDrawStyle(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"bars"`, `"line"`, `"points"`

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withFillBelowTo

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withFillBelowTo(value)
```

PARAMETERS:

* **value** (`string`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withFillColor

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withFillColor(value)
```

PARAMETERS:

* **value** (`string`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withFillOpacity

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withFillOpacity(value)
```

PARAMETERS:

* **value** (`number`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withGradientMode

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withGradientMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"hue"`, `"none"`, `"opacity"`, `"scheme"`

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withHideFrom

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withHideFrom(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withHideFromMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withHideFromMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineColor

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineColor(value)
```

PARAMETERS:

* **value** (`string`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineInterpolation

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineInterpolation(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"linear"`, `"smooth"`, `"stepAfter"`, `"stepBefore"`

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineStyle

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineStyle(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineStyleMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineStyleMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineWidth

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withLineWidth(value)
```

PARAMETERS:

* **value** (`number`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withPointColor

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withPointColor(value)
```

PARAMETERS:

* **value** (`string`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withPointSize

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withPointSize(value)
```

PARAMETERS:

* **value** (`number`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withPointSymbol

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withPointSymbol(value)
```

PARAMETERS:

* **value** (`string`)


######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withScaleDistribution

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withScaleDistribution(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withScaleDistributionMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withScaleDistributionMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withShowPoints

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withShowPoints(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"always"`, `"auto"`, `"never"`

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withSpanNulls

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withSpanNulls(value)
```

PARAMETERS:

* **value** (`["boolean", "number"]`)

Indicate if null values should be treated as gaps or connected. When the value is a number, it represents the maximum delta in the X axis that should be considered connected.  For timeseries, this is milliseconds
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withStacking

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withStacking(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withStackingMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withStackingMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withThresholdsStyle

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withThresholdsStyle(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withThresholdsStyleMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withThresholdsStyleMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withTransform

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withTransform(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"constant"`, `"negative-Y"`

TODO docs
######## fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withType

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.withType()
```



######## obj fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.hideFrom


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.hideFrom.withLegend

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.hideFrom.withLegend(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.hideFrom.withTooltip

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.hideFrom.withTooltip(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.hideFrom.withViz

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.hideFrom.withViz(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


######## obj fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.lineStyle


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.lineStyle.withDash

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.lineStyle.withDash(value)
```

PARAMETERS:

* **value** (`array`)


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.lineStyle.withDashMixin

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.lineStyle.withDashMixin(value)
```

PARAMETERS:

* **value** (`array`)


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.lineStyle.withFill

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.lineStyle.withFill(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"solid"`, `"dash"`, `"dot"`, `"square"`


######## obj fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.scaleDistribution


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.scaleDistribution.withLinearThreshold

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.scaleDistribution.withLinearThreshold(value)
```

PARAMETERS:

* **value** (`number`)


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.scaleDistribution.withLog

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.scaleDistribution.withLog(value)
```

PARAMETERS:

* **value** (`number`)


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.scaleDistribution.withType

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.scaleDistribution.withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"linear"`, `"log"`, `"ordinal"`, `"symlog"`

TODO docs
######## obj fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.stacking


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.stacking.withGroup

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.stacking.withGroup(value)
```

PARAMETERS:

* **value** (`string`)


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.stacking.withMode

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.stacking.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"none"`, `"normal"`, `"percent"`

TODO docs
######## obj fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.thresholdsStyle


######### fn fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.thresholdsStyle.withMode

```jsonnet
fieldConfig.defaults.custom.cellOptions.TableSparklineCellOptions.thresholdsStyle.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"area"`, `"dashed"`, `"dashed+area"`, `"line"`, `"line+area"`, `"off"`, `"series"`

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


#### fn options.withCellHeight

```jsonnet
options.withCellHeight(value)
```

PARAMETERS:

* **value** (`string`)

Controls the height of the rows
#### fn options.withFooter

```jsonnet
options.withFooter(value={"countRows": false,"reducer": [],"show": false})
```

PARAMETERS:

* **value** (`object`)
   - default value: `{"countRows": false,"reducer": [],"show": false}`

Controls footer options
#### fn options.withFooterMixin

```jsonnet
options.withFooterMixin(value={"countRows": false,"reducer": [],"show": false})
```

PARAMETERS:

* **value** (`object`)
   - default value: `{"countRows": false,"reducer": [],"show": false}`

Controls footer options
#### fn options.withFrameIndex

```jsonnet
options.withFrameIndex(value=0)
```

PARAMETERS:

* **value** (`number`)
   - default value: `0`

Represents the index of the selected frame
#### fn options.withShowHeader

```jsonnet
options.withShowHeader(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Controls whether the panel should show the header
#### fn options.withShowTypeIcons

```jsonnet
options.withShowTypeIcons(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Controls whether the header should show icons for the column types
#### fn options.withSortBy

```jsonnet
options.withSortBy(value)
```

PARAMETERS:

* **value** (`array`)

Used to control row sorting
#### fn options.withSortByMixin

```jsonnet
options.withSortByMixin(value)
```

PARAMETERS:

* **value** (`array`)

Used to control row sorting
#### obj options.footer


##### fn options.footer.withTableFooterOptions

```jsonnet
options.footer.withTableFooterOptions(value)
```

PARAMETERS:

* **value** (`object`)

Footer options
##### fn options.footer.withTableFooterOptionsMixin

```jsonnet
options.footer.withTableFooterOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Footer options
##### obj options.footer.TableFooterOptions


###### fn options.footer.TableFooterOptions.withCountRows

```jsonnet
options.footer.TableFooterOptions.withCountRows(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.footer.TableFooterOptions.withEnablePagination

```jsonnet
options.footer.TableFooterOptions.withEnablePagination(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


###### fn options.footer.TableFooterOptions.withFields

```jsonnet
options.footer.TableFooterOptions.withFields(value)
```

PARAMETERS:

* **value** (`array`)


###### fn options.footer.TableFooterOptions.withFieldsMixin

```jsonnet
options.footer.TableFooterOptions.withFieldsMixin(value)
```

PARAMETERS:

* **value** (`array`)


###### fn options.footer.TableFooterOptions.withReducer

```jsonnet
options.footer.TableFooterOptions.withReducer(value)
```

PARAMETERS:

* **value** (`array`)


###### fn options.footer.TableFooterOptions.withReducerMixin

```jsonnet
options.footer.TableFooterOptions.withReducerMixin(value)
```

PARAMETERS:

* **value** (`array`)


###### fn options.footer.TableFooterOptions.withShow

```jsonnet
options.footer.TableFooterOptions.withShow(value=true)
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