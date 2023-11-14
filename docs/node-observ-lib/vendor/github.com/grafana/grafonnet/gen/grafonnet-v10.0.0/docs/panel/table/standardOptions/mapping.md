# mapping



## Index

* [`obj RangeMap`](#obj-rangemap)
  * [`fn withOptions(value)`](#fn-rangemapwithoptions)
  * [`fn withOptionsMixin(value)`](#fn-rangemapwithoptionsmixin)
  * [`fn withType(value)`](#fn-rangemapwithtype)
  * [`obj options`](#obj-rangemapoptions)
    * [`fn withFrom(value)`](#fn-rangemapoptionswithfrom)
    * [`fn withResult(value)`](#fn-rangemapoptionswithresult)
    * [`fn withResultMixin(value)`](#fn-rangemapoptionswithresultmixin)
    * [`fn withTo(value)`](#fn-rangemapoptionswithto)
    * [`obj result`](#obj-rangemapoptionsresult)
      * [`fn withColor(value)`](#fn-rangemapoptionsresultwithcolor)
      * [`fn withIcon(value)`](#fn-rangemapoptionsresultwithicon)
      * [`fn withIndex(value)`](#fn-rangemapoptionsresultwithindex)
      * [`fn withText(value)`](#fn-rangemapoptionsresultwithtext)
* [`obj RegexMap`](#obj-regexmap)
  * [`fn withOptions(value)`](#fn-regexmapwithoptions)
  * [`fn withOptionsMixin(value)`](#fn-regexmapwithoptionsmixin)
  * [`fn withType(value)`](#fn-regexmapwithtype)
  * [`obj options`](#obj-regexmapoptions)
    * [`fn withPattern(value)`](#fn-regexmapoptionswithpattern)
    * [`fn withResult(value)`](#fn-regexmapoptionswithresult)
    * [`fn withResultMixin(value)`](#fn-regexmapoptionswithresultmixin)
    * [`obj result`](#obj-regexmapoptionsresult)
      * [`fn withColor(value)`](#fn-regexmapoptionsresultwithcolor)
      * [`fn withIcon(value)`](#fn-regexmapoptionsresultwithicon)
      * [`fn withIndex(value)`](#fn-regexmapoptionsresultwithindex)
      * [`fn withText(value)`](#fn-regexmapoptionsresultwithtext)
* [`obj SpecialValueMap`](#obj-specialvaluemap)
  * [`fn withOptions(value)`](#fn-specialvaluemapwithoptions)
  * [`fn withOptionsMixin(value)`](#fn-specialvaluemapwithoptionsmixin)
  * [`fn withType(value)`](#fn-specialvaluemapwithtype)
  * [`obj options`](#obj-specialvaluemapoptions)
    * [`fn withMatch(value)`](#fn-specialvaluemapoptionswithmatch)
    * [`fn withPattern(value)`](#fn-specialvaluemapoptionswithpattern)
    * [`fn withResult(value)`](#fn-specialvaluemapoptionswithresult)
    * [`fn withResultMixin(value)`](#fn-specialvaluemapoptionswithresultmixin)
    * [`obj result`](#obj-specialvaluemapoptionsresult)
      * [`fn withColor(value)`](#fn-specialvaluemapoptionsresultwithcolor)
      * [`fn withIcon(value)`](#fn-specialvaluemapoptionsresultwithicon)
      * [`fn withIndex(value)`](#fn-specialvaluemapoptionsresultwithindex)
      * [`fn withText(value)`](#fn-specialvaluemapoptionsresultwithtext)
* [`obj ValueMap`](#obj-valuemap)
  * [`fn withOptions(value)`](#fn-valuemapwithoptions)
  * [`fn withOptionsMixin(value)`](#fn-valuemapwithoptionsmixin)
  * [`fn withType(value)`](#fn-valuemapwithtype)

## Fields

### obj RangeMap


#### fn RangeMap.withOptions

```jsonnet
RangeMap.withOptions(value)
```

PARAMETERS:

* **value** (`object`)


#### fn RangeMap.withOptionsMixin

```jsonnet
RangeMap.withOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn RangeMap.withType

```jsonnet
RangeMap.withType(value)
```

PARAMETERS:

* **value** (`string`)


#### obj RangeMap.options


##### fn RangeMap.options.withFrom

```jsonnet
RangeMap.options.withFrom(value)
```

PARAMETERS:

* **value** (`number`)

to and from are `number | null` in current ts, really not sure what to do
##### fn RangeMap.options.withResult

```jsonnet
RangeMap.options.withResult(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
##### fn RangeMap.options.withResultMixin

```jsonnet
RangeMap.options.withResultMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
##### fn RangeMap.options.withTo

```jsonnet
RangeMap.options.withTo(value)
```

PARAMETERS:

* **value** (`number`)


##### obj RangeMap.options.result


###### fn RangeMap.options.result.withColor

```jsonnet
RangeMap.options.result.withColor(value)
```

PARAMETERS:

* **value** (`string`)


###### fn RangeMap.options.result.withIcon

```jsonnet
RangeMap.options.result.withIcon(value)
```

PARAMETERS:

* **value** (`string`)


###### fn RangeMap.options.result.withIndex

```jsonnet
RangeMap.options.result.withIndex(value)
```

PARAMETERS:

* **value** (`integer`)


###### fn RangeMap.options.result.withText

```jsonnet
RangeMap.options.result.withText(value)
```

PARAMETERS:

* **value** (`string`)


### obj RegexMap


#### fn RegexMap.withOptions

```jsonnet
RegexMap.withOptions(value)
```

PARAMETERS:

* **value** (`object`)


#### fn RegexMap.withOptionsMixin

```jsonnet
RegexMap.withOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn RegexMap.withType

```jsonnet
RegexMap.withType(value)
```

PARAMETERS:

* **value** (`string`)


#### obj RegexMap.options


##### fn RegexMap.options.withPattern

```jsonnet
RegexMap.options.withPattern(value)
```

PARAMETERS:

* **value** (`string`)


##### fn RegexMap.options.withResult

```jsonnet
RegexMap.options.withResult(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
##### fn RegexMap.options.withResultMixin

```jsonnet
RegexMap.options.withResultMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
##### obj RegexMap.options.result


###### fn RegexMap.options.result.withColor

```jsonnet
RegexMap.options.result.withColor(value)
```

PARAMETERS:

* **value** (`string`)


###### fn RegexMap.options.result.withIcon

```jsonnet
RegexMap.options.result.withIcon(value)
```

PARAMETERS:

* **value** (`string`)


###### fn RegexMap.options.result.withIndex

```jsonnet
RegexMap.options.result.withIndex(value)
```

PARAMETERS:

* **value** (`integer`)


###### fn RegexMap.options.result.withText

```jsonnet
RegexMap.options.result.withText(value)
```

PARAMETERS:

* **value** (`string`)


### obj SpecialValueMap


#### fn SpecialValueMap.withOptions

```jsonnet
SpecialValueMap.withOptions(value)
```

PARAMETERS:

* **value** (`object`)


#### fn SpecialValueMap.withOptionsMixin

```jsonnet
SpecialValueMap.withOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn SpecialValueMap.withType

```jsonnet
SpecialValueMap.withType(value)
```

PARAMETERS:

* **value** (`string`)


#### obj SpecialValueMap.options


##### fn SpecialValueMap.options.withMatch

```jsonnet
SpecialValueMap.options.withMatch(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"true"`, `"false"`


##### fn SpecialValueMap.options.withPattern

```jsonnet
SpecialValueMap.options.withPattern(value)
```

PARAMETERS:

* **value** (`string`)


##### fn SpecialValueMap.options.withResult

```jsonnet
SpecialValueMap.options.withResult(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
##### fn SpecialValueMap.options.withResultMixin

```jsonnet
SpecialValueMap.options.withResultMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO docs
##### obj SpecialValueMap.options.result


###### fn SpecialValueMap.options.result.withColor

```jsonnet
SpecialValueMap.options.result.withColor(value)
```

PARAMETERS:

* **value** (`string`)


###### fn SpecialValueMap.options.result.withIcon

```jsonnet
SpecialValueMap.options.result.withIcon(value)
```

PARAMETERS:

* **value** (`string`)


###### fn SpecialValueMap.options.result.withIndex

```jsonnet
SpecialValueMap.options.result.withIndex(value)
```

PARAMETERS:

* **value** (`integer`)


###### fn SpecialValueMap.options.result.withText

```jsonnet
SpecialValueMap.options.result.withText(value)
```

PARAMETERS:

* **value** (`string`)


### obj ValueMap


#### fn ValueMap.withOptions

```jsonnet
ValueMap.withOptions(value)
```

PARAMETERS:

* **value** (`object`)


#### fn ValueMap.withOptionsMixin

```jsonnet
ValueMap.withOptionsMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn ValueMap.withType

```jsonnet
ValueMap.withType(value)
```

PARAMETERS:

* **value** (`string`)

