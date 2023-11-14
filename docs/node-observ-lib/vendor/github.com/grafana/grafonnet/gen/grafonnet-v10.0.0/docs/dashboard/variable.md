# variable

Example usage:

```jsonnet
local g = import 'g.libsonnet';
local var = g.dashboard.variable;

local customVar =
  var.custom.new(
    'myOptions',
    values=['a', 'b', 'c', 'd'],
  )
  + var.custom.generalOptions.withDescription(
    'This is a variable for my custom options.'
  )
  + var.custom.selectionOptions.withMulti();

local queryVar =
  var.query.new('queryOptions')
  + var.query.queryTypes.withLabelValues(
    'up',
    'instance',
  )
  + var.query.withDatasource(
    type='prometheus',
    uid='mimir-prod',
  )
  + var.query.selectionOptions.withIncludeAll();


g.dashboard.new('my dashboard')
+ g.dashboard.withVariables([
  customVar,
  queryVar,
])
```


## Index

* [`obj adhoc`](#obj-adhoc)
  * [`fn new(name, type, uid)`](#fn-adhocnew)
  * [`fn newFromDatasourceVariable(name, variable)`](#fn-adhocnewfromdatasourcevariable)
  * [`obj generalOptions`](#obj-adhocgeneraloptions)
    * [`fn withCurrent(key, value="<same-as-key>")`](#fn-adhocgeneraloptionswithcurrent)
    * [`fn withDescription(value)`](#fn-adhocgeneraloptionswithdescription)
    * [`fn withLabel(value)`](#fn-adhocgeneraloptionswithlabel)
    * [`fn withName(value)`](#fn-adhocgeneraloptionswithname)
    * [`obj showOnDashboard`](#obj-adhocgeneraloptionsshowondashboard)
      * [`fn withLabelAndValue()`](#fn-adhocgeneraloptionsshowondashboardwithlabelandvalue)
      * [`fn withNothing()`](#fn-adhocgeneraloptionsshowondashboardwithnothing)
      * [`fn withValueOnly()`](#fn-adhocgeneraloptionsshowondashboardwithvalueonly)
* [`obj constant`](#obj-constant)
  * [`fn new(name, value)`](#fn-constantnew)
  * [`obj generalOptions`](#obj-constantgeneraloptions)
    * [`fn withCurrent(key, value="<same-as-key>")`](#fn-constantgeneraloptionswithcurrent)
    * [`fn withDescription(value)`](#fn-constantgeneraloptionswithdescription)
    * [`fn withLabel(value)`](#fn-constantgeneraloptionswithlabel)
    * [`fn withName(value)`](#fn-constantgeneraloptionswithname)
    * [`obj showOnDashboard`](#obj-constantgeneraloptionsshowondashboard)
      * [`fn withLabelAndValue()`](#fn-constantgeneraloptionsshowondashboardwithlabelandvalue)
      * [`fn withNothing()`](#fn-constantgeneraloptionsshowondashboardwithnothing)
      * [`fn withValueOnly()`](#fn-constantgeneraloptionsshowondashboardwithvalueonly)
* [`obj custom`](#obj-custom)
  * [`fn new(name, values)`](#fn-customnew)
  * [`obj generalOptions`](#obj-customgeneraloptions)
    * [`fn withCurrent(key, value="<same-as-key>")`](#fn-customgeneraloptionswithcurrent)
    * [`fn withDescription(value)`](#fn-customgeneraloptionswithdescription)
    * [`fn withLabel(value)`](#fn-customgeneraloptionswithlabel)
    * [`fn withName(value)`](#fn-customgeneraloptionswithname)
    * [`obj showOnDashboard`](#obj-customgeneraloptionsshowondashboard)
      * [`fn withLabelAndValue()`](#fn-customgeneraloptionsshowondashboardwithlabelandvalue)
      * [`fn withNothing()`](#fn-customgeneraloptionsshowondashboardwithnothing)
      * [`fn withValueOnly()`](#fn-customgeneraloptionsshowondashboardwithvalueonly)
  * [`obj selectionOptions`](#obj-customselectionoptions)
    * [`fn withIncludeAll(value=true, customAllValue)`](#fn-customselectionoptionswithincludeall)
    * [`fn withMulti(value=true)`](#fn-customselectionoptionswithmulti)
* [`obj datasource`](#obj-datasource)
  * [`fn new(name, type)`](#fn-datasourcenew)
  * [`fn withRegex(value)`](#fn-datasourcewithregex)
  * [`obj generalOptions`](#obj-datasourcegeneraloptions)
    * [`fn withCurrent(key, value="<same-as-key>")`](#fn-datasourcegeneraloptionswithcurrent)
    * [`fn withDescription(value)`](#fn-datasourcegeneraloptionswithdescription)
    * [`fn withLabel(value)`](#fn-datasourcegeneraloptionswithlabel)
    * [`fn withName(value)`](#fn-datasourcegeneraloptionswithname)
    * [`obj showOnDashboard`](#obj-datasourcegeneraloptionsshowondashboard)
      * [`fn withLabelAndValue()`](#fn-datasourcegeneraloptionsshowondashboardwithlabelandvalue)
      * [`fn withNothing()`](#fn-datasourcegeneraloptionsshowondashboardwithnothing)
      * [`fn withValueOnly()`](#fn-datasourcegeneraloptionsshowondashboardwithvalueonly)
  * [`obj selectionOptions`](#obj-datasourceselectionoptions)
    * [`fn withIncludeAll(value=true, customAllValue)`](#fn-datasourceselectionoptionswithincludeall)
    * [`fn withMulti(value=true)`](#fn-datasourceselectionoptionswithmulti)
* [`obj interval`](#obj-interval)
  * [`fn new(name, values)`](#fn-intervalnew)
  * [`fn withAutoOption(count, minInterval)`](#fn-intervalwithautooption)
  * [`obj generalOptions`](#obj-intervalgeneraloptions)
    * [`fn withCurrent(key, value="<same-as-key>")`](#fn-intervalgeneraloptionswithcurrent)
    * [`fn withDescription(value)`](#fn-intervalgeneraloptionswithdescription)
    * [`fn withLabel(value)`](#fn-intervalgeneraloptionswithlabel)
    * [`fn withName(value)`](#fn-intervalgeneraloptionswithname)
    * [`obj showOnDashboard`](#obj-intervalgeneraloptionsshowondashboard)
      * [`fn withLabelAndValue()`](#fn-intervalgeneraloptionsshowondashboardwithlabelandvalue)
      * [`fn withNothing()`](#fn-intervalgeneraloptionsshowondashboardwithnothing)
      * [`fn withValueOnly()`](#fn-intervalgeneraloptionsshowondashboardwithvalueonly)
* [`obj query`](#obj-query)
  * [`fn new(name, query="")`](#fn-querynew)
  * [`fn withDatasource(type, uid)`](#fn-querywithdatasource)
  * [`fn withDatasourceFromVariable(variable)`](#fn-querywithdatasourcefromvariable)
  * [`fn withRegex(value)`](#fn-querywithregex)
  * [`fn withSort(i=0, type="alphabetical", asc=true, caseInsensitive=false)`](#fn-querywithsort)
  * [`obj generalOptions`](#obj-querygeneraloptions)
    * [`fn withCurrent(key, value="<same-as-key>")`](#fn-querygeneraloptionswithcurrent)
    * [`fn withDescription(value)`](#fn-querygeneraloptionswithdescription)
    * [`fn withLabel(value)`](#fn-querygeneraloptionswithlabel)
    * [`fn withName(value)`](#fn-querygeneraloptionswithname)
    * [`obj showOnDashboard`](#obj-querygeneraloptionsshowondashboard)
      * [`fn withLabelAndValue()`](#fn-querygeneraloptionsshowondashboardwithlabelandvalue)
      * [`fn withNothing()`](#fn-querygeneraloptionsshowondashboardwithnothing)
      * [`fn withValueOnly()`](#fn-querygeneraloptionsshowondashboardwithvalueonly)
  * [`obj queryTypes`](#obj-queryquerytypes)
    * [`fn withLabelValues(label, metric="")`](#fn-queryquerytypeswithlabelvalues)
  * [`obj refresh`](#obj-queryrefresh)
    * [`fn onLoad()`](#fn-queryrefreshonload)
    * [`fn onTime()`](#fn-queryrefreshontime)
  * [`obj selectionOptions`](#obj-queryselectionoptions)
    * [`fn withIncludeAll(value=true, customAllValue)`](#fn-queryselectionoptionswithincludeall)
    * [`fn withMulti(value=true)`](#fn-queryselectionoptionswithmulti)
* [`obj textbox`](#obj-textbox)
  * [`fn new(name, default="")`](#fn-textboxnew)
  * [`obj generalOptions`](#obj-textboxgeneraloptions)
    * [`fn withCurrent(key, value="<same-as-key>")`](#fn-textboxgeneraloptionswithcurrent)
    * [`fn withDescription(value)`](#fn-textboxgeneraloptionswithdescription)
    * [`fn withLabel(value)`](#fn-textboxgeneraloptionswithlabel)
    * [`fn withName(value)`](#fn-textboxgeneraloptionswithname)
    * [`obj showOnDashboard`](#obj-textboxgeneraloptionsshowondashboard)
      * [`fn withLabelAndValue()`](#fn-textboxgeneraloptionsshowondashboardwithlabelandvalue)
      * [`fn withNothing()`](#fn-textboxgeneraloptionsshowondashboardwithnothing)
      * [`fn withValueOnly()`](#fn-textboxgeneraloptionsshowondashboardwithvalueonly)

## Fields

### obj adhoc


#### fn adhoc.new

```jsonnet
adhoc.new(name, type, uid)
```

PARAMETERS:

* **name** (`string`)
* **type** (`string`)
* **uid** (`string`)

`new` creates an adhoc template variable for datasource with `type` and `uid`.
#### fn adhoc.newFromDatasourceVariable

```jsonnet
adhoc.newFromDatasourceVariable(name, variable)
```

PARAMETERS:

* **name** (`string`)
* **variable** (`object`)

Same as `new` but selecting the datasource from another template variable.
#### obj adhoc.generalOptions


##### fn adhoc.generalOptions.withCurrent

```jsonnet
adhoc.generalOptions.withCurrent(key, value="<same-as-key>")
```

PARAMETERS:

* **key** (`any`)
* **value** (`any`)
   - default value: `"<same-as-key>"`

`withCurrent` sets the currently selected value of a variable. If key and value are different, both need to be given.

##### fn adhoc.generalOptions.withDescription

```jsonnet
adhoc.generalOptions.withDescription(value)
```

PARAMETERS:

* **value** (`string`)


##### fn adhoc.generalOptions.withLabel

```jsonnet
adhoc.generalOptions.withLabel(value)
```

PARAMETERS:

* **value** (`string`)


##### fn adhoc.generalOptions.withName

```jsonnet
adhoc.generalOptions.withName(value)
```

PARAMETERS:

* **value** (`string`)


##### obj adhoc.generalOptions.showOnDashboard


###### fn adhoc.generalOptions.showOnDashboard.withLabelAndValue

```jsonnet
adhoc.generalOptions.showOnDashboard.withLabelAndValue()
```



###### fn adhoc.generalOptions.showOnDashboard.withNothing

```jsonnet
adhoc.generalOptions.showOnDashboard.withNothing()
```



###### fn adhoc.generalOptions.showOnDashboard.withValueOnly

```jsonnet
adhoc.generalOptions.showOnDashboard.withValueOnly()
```



### obj constant


#### fn constant.new

```jsonnet
constant.new(name, value)
```

PARAMETERS:

* **name** (`string`)
* **value** (`string`)

`new` creates a hidden constant template variable.
#### obj constant.generalOptions


##### fn constant.generalOptions.withCurrent

```jsonnet
constant.generalOptions.withCurrent(key, value="<same-as-key>")
```

PARAMETERS:

* **key** (`any`)
* **value** (`any`)
   - default value: `"<same-as-key>"`

`withCurrent` sets the currently selected value of a variable. If key and value are different, both need to be given.

##### fn constant.generalOptions.withDescription

```jsonnet
constant.generalOptions.withDescription(value)
```

PARAMETERS:

* **value** (`string`)


##### fn constant.generalOptions.withLabel

```jsonnet
constant.generalOptions.withLabel(value)
```

PARAMETERS:

* **value** (`string`)


##### fn constant.generalOptions.withName

```jsonnet
constant.generalOptions.withName(value)
```

PARAMETERS:

* **value** (`string`)


##### obj constant.generalOptions.showOnDashboard


###### fn constant.generalOptions.showOnDashboard.withLabelAndValue

```jsonnet
constant.generalOptions.showOnDashboard.withLabelAndValue()
```



###### fn constant.generalOptions.showOnDashboard.withNothing

```jsonnet
constant.generalOptions.showOnDashboard.withNothing()
```



###### fn constant.generalOptions.showOnDashboard.withValueOnly

```jsonnet
constant.generalOptions.showOnDashboard.withValueOnly()
```



### obj custom


#### fn custom.new

```jsonnet
custom.new(name, values)
```

PARAMETERS:

* **name** (`string`)
* **values** (`array`)

`new` creates a custom template variable.

The `values` array accepts an object with key/value keys, if it's not an object
then it will be added as a string.

Example:
```
[
  { key: 'mykey', value: 'myvalue' },
  'myvalue',
  12,
]

#### obj custom.generalOptions


##### fn custom.generalOptions.withCurrent

```jsonnet
custom.generalOptions.withCurrent(key, value="<same-as-key>")
```

PARAMETERS:

* **key** (`any`)
* **value** (`any`)
   - default value: `"<same-as-key>"`

`withCurrent` sets the currently selected value of a variable. If key and value are different, both need to be given.

##### fn custom.generalOptions.withDescription

```jsonnet
custom.generalOptions.withDescription(value)
```

PARAMETERS:

* **value** (`string`)


##### fn custom.generalOptions.withLabel

```jsonnet
custom.generalOptions.withLabel(value)
```

PARAMETERS:

* **value** (`string`)


##### fn custom.generalOptions.withName

```jsonnet
custom.generalOptions.withName(value)
```

PARAMETERS:

* **value** (`string`)


##### obj custom.generalOptions.showOnDashboard


###### fn custom.generalOptions.showOnDashboard.withLabelAndValue

```jsonnet
custom.generalOptions.showOnDashboard.withLabelAndValue()
```



###### fn custom.generalOptions.showOnDashboard.withNothing

```jsonnet
custom.generalOptions.showOnDashboard.withNothing()
```



###### fn custom.generalOptions.showOnDashboard.withValueOnly

```jsonnet
custom.generalOptions.showOnDashboard.withValueOnly()
```



#### obj custom.selectionOptions


##### fn custom.selectionOptions.withIncludeAll

```jsonnet
custom.selectionOptions.withIncludeAll(value=true, customAllValue)
```

PARAMETERS:

* **value** (`bool`)
   - default value: `true`
* **customAllValue** (`bool`)

`withIncludeAll` enables an option to include all variables.

Optionally you can set a `customAllValue`.

##### fn custom.selectionOptions.withMulti

```jsonnet
custom.selectionOptions.withMulti(value=true)
```

PARAMETERS:

* **value** (`bool`)
   - default value: `true`

Enable selecting multiple values.
### obj datasource


#### fn datasource.new

```jsonnet
datasource.new(name, type)
```

PARAMETERS:

* **name** (`string`)
* **type** (`string`)

`new` creates a datasource template variable.
#### fn datasource.withRegex

```jsonnet
datasource.withRegex(value)
```

PARAMETERS:

* **value** (`string`)

`withRegex` filter for which data source instances to choose from in the
variable value list. Example: `/^prod/`

#### obj datasource.generalOptions


##### fn datasource.generalOptions.withCurrent

```jsonnet
datasource.generalOptions.withCurrent(key, value="<same-as-key>")
```

PARAMETERS:

* **key** (`any`)
* **value** (`any`)
   - default value: `"<same-as-key>"`

`withCurrent` sets the currently selected value of a variable. If key and value are different, both need to be given.

##### fn datasource.generalOptions.withDescription

```jsonnet
datasource.generalOptions.withDescription(value)
```

PARAMETERS:

* **value** (`string`)


##### fn datasource.generalOptions.withLabel

```jsonnet
datasource.generalOptions.withLabel(value)
```

PARAMETERS:

* **value** (`string`)


##### fn datasource.generalOptions.withName

```jsonnet
datasource.generalOptions.withName(value)
```

PARAMETERS:

* **value** (`string`)


##### obj datasource.generalOptions.showOnDashboard


###### fn datasource.generalOptions.showOnDashboard.withLabelAndValue

```jsonnet
datasource.generalOptions.showOnDashboard.withLabelAndValue()
```



###### fn datasource.generalOptions.showOnDashboard.withNothing

```jsonnet
datasource.generalOptions.showOnDashboard.withNothing()
```



###### fn datasource.generalOptions.showOnDashboard.withValueOnly

```jsonnet
datasource.generalOptions.showOnDashboard.withValueOnly()
```



#### obj datasource.selectionOptions


##### fn datasource.selectionOptions.withIncludeAll

```jsonnet
datasource.selectionOptions.withIncludeAll(value=true, customAllValue)
```

PARAMETERS:

* **value** (`bool`)
   - default value: `true`
* **customAllValue** (`bool`)

`withIncludeAll` enables an option to include all variables.

Optionally you can set a `customAllValue`.

##### fn datasource.selectionOptions.withMulti

```jsonnet
datasource.selectionOptions.withMulti(value=true)
```

PARAMETERS:

* **value** (`bool`)
   - default value: `true`

Enable selecting multiple values.
### obj interval


#### fn interval.new

```jsonnet
interval.new(name, values)
```

PARAMETERS:

* **name** (`string`)
* **values** (`array`)

`new` creates an interval template variable.
#### fn interval.withAutoOption

```jsonnet
interval.withAutoOption(count, minInterval)
```

PARAMETERS:

* **count** (`number`)
* **minInterval** (`string`)

`withAutoOption` adds an options to dynamically calculate interval by dividing
time range by the count specified.

`minInterval' has to be either unit-less or end with one of the following units:
"y, M, w, d, h, m, s, ms".

#### obj interval.generalOptions


##### fn interval.generalOptions.withCurrent

```jsonnet
interval.generalOptions.withCurrent(key, value="<same-as-key>")
```

PARAMETERS:

* **key** (`any`)
* **value** (`any`)
   - default value: `"<same-as-key>"`

`withCurrent` sets the currently selected value of a variable. If key and value are different, both need to be given.

##### fn interval.generalOptions.withDescription

```jsonnet
interval.generalOptions.withDescription(value)
```

PARAMETERS:

* **value** (`string`)


##### fn interval.generalOptions.withLabel

```jsonnet
interval.generalOptions.withLabel(value)
```

PARAMETERS:

* **value** (`string`)


##### fn interval.generalOptions.withName

```jsonnet
interval.generalOptions.withName(value)
```

PARAMETERS:

* **value** (`string`)


##### obj interval.generalOptions.showOnDashboard


###### fn interval.generalOptions.showOnDashboard.withLabelAndValue

```jsonnet
interval.generalOptions.showOnDashboard.withLabelAndValue()
```



###### fn interval.generalOptions.showOnDashboard.withNothing

```jsonnet
interval.generalOptions.showOnDashboard.withNothing()
```



###### fn interval.generalOptions.showOnDashboard.withValueOnly

```jsonnet
interval.generalOptions.showOnDashboard.withValueOnly()
```



### obj query


#### fn query.new

```jsonnet
query.new(name, query="")
```

PARAMETERS:

* **name** (`string`)
* **query** (`string`)
   - default value: `""`

Create a query template variable.

`query` argument is optional, this can also be set with `query.queryTypes`.

#### fn query.withDatasource

```jsonnet
query.withDatasource(type, uid)
```

PARAMETERS:

* **type** (`string`)
* **uid** (`string`)

Select a datasource for the variable template query.
#### fn query.withDatasourceFromVariable

```jsonnet
query.withDatasourceFromVariable(variable)
```

PARAMETERS:

* **variable** (`object`)

Select the datasource from another template variable.
#### fn query.withRegex

```jsonnet
query.withRegex(value)
```

PARAMETERS:

* **value** (`string`)

`withRegex` can extract part of a series name or metric node segment. Named
capture groups can be used to separate the display text and value
([see examples](https://grafana.com/docs/grafana/latest/variables/filter-variables-with-regex#filter-and-modify-using-named-text-and-value-capture-groups)).

#### fn query.withSort

```jsonnet
query.withSort(i=0, type="alphabetical", asc=true, caseInsensitive=false)
```

PARAMETERS:

* **i** (`number`)
   - default value: `0`
* **type** (`string`)
   - default value: `"alphabetical"`
* **asc** (`bool`)
   - default value: `true`
* **caseInsensitive** (`bool`)
   - default value: `false`

Choose how to sort the values in the dropdown.

This can be called as `withSort(<number>) to use the integer values for each
option. If `i==0` then it will be ignored and the other arguments will take
precedence.

The numerical values are:

- 1 - Alphabetical (asc)
- 2 - Alphabetical (desc)
- 3 - Numerical (asc)
- 4 - Numerical (desc)
- 5 - Alphabetical (case-insensitive, asc)
- 6 - Alphabetical (case-insensitive, desc)

#### obj query.generalOptions


##### fn query.generalOptions.withCurrent

```jsonnet
query.generalOptions.withCurrent(key, value="<same-as-key>")
```

PARAMETERS:

* **key** (`any`)
* **value** (`any`)
   - default value: `"<same-as-key>"`

`withCurrent` sets the currently selected value of a variable. If key and value are different, both need to be given.

##### fn query.generalOptions.withDescription

```jsonnet
query.generalOptions.withDescription(value)
```

PARAMETERS:

* **value** (`string`)


##### fn query.generalOptions.withLabel

```jsonnet
query.generalOptions.withLabel(value)
```

PARAMETERS:

* **value** (`string`)


##### fn query.generalOptions.withName

```jsonnet
query.generalOptions.withName(value)
```

PARAMETERS:

* **value** (`string`)


##### obj query.generalOptions.showOnDashboard


###### fn query.generalOptions.showOnDashboard.withLabelAndValue

```jsonnet
query.generalOptions.showOnDashboard.withLabelAndValue()
```



###### fn query.generalOptions.showOnDashboard.withNothing

```jsonnet
query.generalOptions.showOnDashboard.withNothing()
```



###### fn query.generalOptions.showOnDashboard.withValueOnly

```jsonnet
query.generalOptions.showOnDashboard.withValueOnly()
```



#### obj query.queryTypes


##### fn query.queryTypes.withLabelValues

```jsonnet
query.queryTypes.withLabelValues(label, metric="")
```

PARAMETERS:

* **label** (`string`)
* **metric** (`string`)
   - default value: `""`

Construct a Prometheus template variable using `label_values()`.
#### obj query.refresh


##### fn query.refresh.onLoad

```jsonnet
query.refresh.onLoad()
```


Refresh label values on dashboard load.
##### fn query.refresh.onTime

```jsonnet
query.refresh.onTime()
```


Refresh label values on time range change.
#### obj query.selectionOptions


##### fn query.selectionOptions.withIncludeAll

```jsonnet
query.selectionOptions.withIncludeAll(value=true, customAllValue)
```

PARAMETERS:

* **value** (`bool`)
   - default value: `true`
* **customAllValue** (`bool`)

`withIncludeAll` enables an option to include all variables.

Optionally you can set a `customAllValue`.

##### fn query.selectionOptions.withMulti

```jsonnet
query.selectionOptions.withMulti(value=true)
```

PARAMETERS:

* **value** (`bool`)
   - default value: `true`

Enable selecting multiple values.
### obj textbox


#### fn textbox.new

```jsonnet
textbox.new(name, default="")
```

PARAMETERS:

* **name** (`string`)
* **default** (`string`)
   - default value: `""`

`new` creates a textbox template variable.
#### obj textbox.generalOptions


##### fn textbox.generalOptions.withCurrent

```jsonnet
textbox.generalOptions.withCurrent(key, value="<same-as-key>")
```

PARAMETERS:

* **key** (`any`)
* **value** (`any`)
   - default value: `"<same-as-key>"`

`withCurrent` sets the currently selected value of a variable. If key and value are different, both need to be given.

##### fn textbox.generalOptions.withDescription

```jsonnet
textbox.generalOptions.withDescription(value)
```

PARAMETERS:

* **value** (`string`)


##### fn textbox.generalOptions.withLabel

```jsonnet
textbox.generalOptions.withLabel(value)
```

PARAMETERS:

* **value** (`string`)


##### fn textbox.generalOptions.withName

```jsonnet
textbox.generalOptions.withName(value)
```

PARAMETERS:

* **value** (`string`)


##### obj textbox.generalOptions.showOnDashboard


###### fn textbox.generalOptions.showOnDashboard.withLabelAndValue

```jsonnet
textbox.generalOptions.showOnDashboard.withLabelAndValue()
```



###### fn textbox.generalOptions.showOnDashboard.withNothing

```jsonnet
textbox.generalOptions.showOnDashboard.withNothing()
```



###### fn textbox.generalOptions.showOnDashboard.withValueOnly

```jsonnet
textbox.generalOptions.showOnDashboard.withValueOnly()
```


