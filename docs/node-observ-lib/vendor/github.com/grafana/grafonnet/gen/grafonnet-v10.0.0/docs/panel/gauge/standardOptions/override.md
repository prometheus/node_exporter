# override

Overrides allow you to customize visualization settings for specific fields or
series. This is accomplished by adding an override rule that targets
a particular set of fields and that can each define multiple options.

```jsonnet
override.byType.new('number')
+ override.byType.withPropertiesFromOptions(
  panel.standardOptions.withDecimals(2)
  + panel.standardOptions.withUnit('s')
)
```


## Index

* [`obj byName`](#obj-byname)
  * [`fn new(value)`](#fn-bynamenew)
  * [`fn withPropertiesFromOptions(options)`](#fn-bynamewithpropertiesfromoptions)
  * [`fn withProperty(id, value)`](#fn-bynamewithproperty)
* [`obj byQuery`](#obj-byquery)
  * [`fn new(value)`](#fn-byquerynew)
  * [`fn withPropertiesFromOptions(options)`](#fn-byquerywithpropertiesfromoptions)
  * [`fn withProperty(id, value)`](#fn-byquerywithproperty)
* [`obj byRegexp`](#obj-byregexp)
  * [`fn new(value)`](#fn-byregexpnew)
  * [`fn withPropertiesFromOptions(options)`](#fn-byregexpwithpropertiesfromoptions)
  * [`fn withProperty(id, value)`](#fn-byregexpwithproperty)
* [`obj byType`](#obj-bytype)
  * [`fn new(value)`](#fn-bytypenew)
  * [`fn withPropertiesFromOptions(options)`](#fn-bytypewithpropertiesfromoptions)
  * [`fn withProperty(id, value)`](#fn-bytypewithproperty)
* [`obj byValue`](#obj-byvalue)
  * [`fn new(value)`](#fn-byvaluenew)
  * [`fn withPropertiesFromOptions(options)`](#fn-byvaluewithpropertiesfromoptions)
  * [`fn withProperty(id, value)`](#fn-byvaluewithproperty)

## Fields

### obj byName


#### fn byName.new

```jsonnet
byName.new(value)
```

PARAMETERS:

* **value** (`string`)

`new` creates a new override of type `byName`.
#### fn byName.withPropertiesFromOptions

```jsonnet
byName.withPropertiesFromOptions(options)
```

PARAMETERS:

* **options** (`object`)

`withPropertiesFromOptions` takes an object with properties that need to be
overridden. See example code above.

#### fn byName.withProperty

```jsonnet
byName.withProperty(id, value)
```

PARAMETERS:

* **id** (`string`)
* **value** (`any`)

`withProperty` adds a property that needs to be overridden. This function can
be called multiple time, adding more properties.

### obj byQuery


#### fn byQuery.new

```jsonnet
byQuery.new(value)
```

PARAMETERS:

* **value** (`string`)

`new` creates a new override of type `byQuery`.
#### fn byQuery.withPropertiesFromOptions

```jsonnet
byQuery.withPropertiesFromOptions(options)
```

PARAMETERS:

* **options** (`object`)

`withPropertiesFromOptions` takes an object with properties that need to be
overridden. See example code above.

#### fn byQuery.withProperty

```jsonnet
byQuery.withProperty(id, value)
```

PARAMETERS:

* **id** (`string`)
* **value** (`any`)

`withProperty` adds a property that needs to be overridden. This function can
be called multiple time, adding more properties.

### obj byRegexp


#### fn byRegexp.new

```jsonnet
byRegexp.new(value)
```

PARAMETERS:

* **value** (`string`)

`new` creates a new override of type `byRegexp`.
#### fn byRegexp.withPropertiesFromOptions

```jsonnet
byRegexp.withPropertiesFromOptions(options)
```

PARAMETERS:

* **options** (`object`)

`withPropertiesFromOptions` takes an object with properties that need to be
overridden. See example code above.

#### fn byRegexp.withProperty

```jsonnet
byRegexp.withProperty(id, value)
```

PARAMETERS:

* **id** (`string`)
* **value** (`any`)

`withProperty` adds a property that needs to be overridden. This function can
be called multiple time, adding more properties.

### obj byType


#### fn byType.new

```jsonnet
byType.new(value)
```

PARAMETERS:

* **value** (`string`)

`new` creates a new override of type `byType`.
#### fn byType.withPropertiesFromOptions

```jsonnet
byType.withPropertiesFromOptions(options)
```

PARAMETERS:

* **options** (`object`)

`withPropertiesFromOptions` takes an object with properties that need to be
overridden. See example code above.

#### fn byType.withProperty

```jsonnet
byType.withProperty(id, value)
```

PARAMETERS:

* **id** (`string`)
* **value** (`any`)

`withProperty` adds a property that needs to be overridden. This function can
be called multiple time, adding more properties.

### obj byValue


#### fn byValue.new

```jsonnet
byValue.new(value)
```

PARAMETERS:

* **value** (`string`)

`new` creates a new override of type `byValue`.
#### fn byValue.withPropertiesFromOptions

```jsonnet
byValue.withPropertiesFromOptions(options)
```

PARAMETERS:

* **options** (`object`)

`withPropertiesFromOptions` takes an object with properties that need to be
overridden. See example code above.

#### fn byValue.withProperty

```jsonnet
byValue.withProperty(id, value)
```

PARAMETERS:

* **id** (`string`)
* **value** (`any`)

`withProperty` adds a property that needs to be overridden. This function can
be called multiple time, adding more properties.
