# rule



## Subpackages

* [data](data.md)

## Index

* [`fn withAnnotations(value)`](#fn-withannotations)
* [`fn withAnnotationsMixin(value)`](#fn-withannotationsmixin)
* [`fn withCondition(value)`](#fn-withcondition)
* [`fn withData(value)`](#fn-withdata)
* [`fn withDataMixin(value)`](#fn-withdatamixin)
* [`fn withExecErrState(value)`](#fn-withexecerrstate)
* [`fn withFor(value)`](#fn-withfor)
* [`fn withIsPaused(value=true)`](#fn-withispaused)
* [`fn withLabels(value)`](#fn-withlabels)
* [`fn withLabelsMixin(value)`](#fn-withlabelsmixin)
* [`fn withName(value)`](#fn-withname)
* [`fn withNoDataState(value)`](#fn-withnodatastate)

## Fields

### fn withAnnotations

```jsonnet
withAnnotations(value)
```

PARAMETERS:

* **value** (`object`)


### fn withAnnotationsMixin

```jsonnet
withAnnotationsMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withCondition

```jsonnet
withCondition(value)
```

PARAMETERS:

* **value** (`string`)


### fn withData

```jsonnet
withData(value)
```

PARAMETERS:

* **value** (`array`)


### fn withDataMixin

```jsonnet
withDataMixin(value)
```

PARAMETERS:

* **value** (`array`)


### fn withExecErrState

```jsonnet
withExecErrState(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"OK"`, `"Alerting"`, `"Error"`


### fn withFor

```jsonnet
withFor(value)
```

PARAMETERS:

* **value** (`integer`)

A Duration represents the elapsed time between two instants
as an int64 nanosecond count. The representation limits the
largest representable duration to approximately 290 years.
### fn withIsPaused

```jsonnet
withIsPaused(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


### fn withLabels

```jsonnet
withLabels(value)
```

PARAMETERS:

* **value** (`object`)


### fn withLabelsMixin

```jsonnet
withLabelsMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withName

```jsonnet
withName(value)
```

PARAMETERS:

* **value** (`string`)


### fn withNoDataState

```jsonnet
withNoDataState(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"Alerting"`, `"NoData"`, `"OK"`

