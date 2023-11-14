# layers



## Index

* [`fn withConfig(value)`](#fn-withconfig)
* [`fn withFilterData(value)`](#fn-withfilterdata)
* [`fn withLocation(value)`](#fn-withlocation)
* [`fn withLocationMixin(value)`](#fn-withlocationmixin)
* [`fn withName(value)`](#fn-withname)
* [`fn withOpacity(value)`](#fn-withopacity)
* [`fn withTooltip(value=true)`](#fn-withtooltip)
* [`fn withType(value)`](#fn-withtype)
* [`obj location`](#obj-location)
  * [`fn withGazetteer(value)`](#fn-locationwithgazetteer)
  * [`fn withGeohash(value)`](#fn-locationwithgeohash)
  * [`fn withLatitude(value)`](#fn-locationwithlatitude)
  * [`fn withLongitude(value)`](#fn-locationwithlongitude)
  * [`fn withLookup(value)`](#fn-locationwithlookup)
  * [`fn withMode(value)`](#fn-locationwithmode)
  * [`fn withWkt(value)`](#fn-locationwithwkt)

## Fields

### fn withConfig

```jsonnet
withConfig(value)
```

PARAMETERS:

* **value** (`string`)

Custom options depending on the type
### fn withFilterData

```jsonnet
withFilterData(value)
```

PARAMETERS:

* **value** (`string`)

Defines a frame MatcherConfig that may filter data for the given layer
### fn withLocation

```jsonnet
withLocation(value)
```

PARAMETERS:

* **value** (`object`)


### fn withLocationMixin

```jsonnet
withLocationMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withName

```jsonnet
withName(value)
```

PARAMETERS:

* **value** (`string`)

configured unique display name
### fn withOpacity

```jsonnet
withOpacity(value)
```

PARAMETERS:

* **value** (`integer`)

Common properties:
https://openlayers.org/en/latest/apidoc/module-ol_layer_Base-BaseLayer.html
Layer opacity (0-1)
### fn withTooltip

```jsonnet
withTooltip(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Check tooltip (defaults to true)
### fn withType

```jsonnet
withType(value)
```

PARAMETERS:

* **value** (`string`)


### obj location


#### fn location.withGazetteer

```jsonnet
location.withGazetteer(value)
```

PARAMETERS:

* **value** (`string`)

Path to Gazetteer
#### fn location.withGeohash

```jsonnet
location.withGeohash(value)
```

PARAMETERS:

* **value** (`string`)

Field mappings
#### fn location.withLatitude

```jsonnet
location.withLatitude(value)
```

PARAMETERS:

* **value** (`string`)


#### fn location.withLongitude

```jsonnet
location.withLongitude(value)
```

PARAMETERS:

* **value** (`string`)


#### fn location.withLookup

```jsonnet
location.withLookup(value)
```

PARAMETERS:

* **value** (`string`)


#### fn location.withMode

```jsonnet
location.withMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"auto"`, `"geohash"`, `"coords"`, `"lookup"`


#### fn location.withWkt

```jsonnet
location.withWkt(value)
```

PARAMETERS:

* **value** (`string`)

