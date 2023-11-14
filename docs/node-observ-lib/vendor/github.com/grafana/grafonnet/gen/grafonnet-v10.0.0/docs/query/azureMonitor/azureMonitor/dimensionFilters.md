# dimensionFilters



## Index

* [`fn withDimension(value)`](#fn-withdimension)
* [`fn withFilter(value)`](#fn-withfilter)
* [`fn withFilters(value)`](#fn-withfilters)
* [`fn withFiltersMixin(value)`](#fn-withfiltersmixin)
* [`fn withOperator(value)`](#fn-withoperator)

## Fields

### fn withDimension

```jsonnet
withDimension(value)
```

PARAMETERS:

* **value** (`string`)

Name of Dimension to be filtered on.
### fn withFilter

```jsonnet
withFilter(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated filter is deprecated in favour of filters to support multiselect.
### fn withFilters

```jsonnet
withFilters(value)
```

PARAMETERS:

* **value** (`array`)

Values to match with the filter.
### fn withFiltersMixin

```jsonnet
withFiltersMixin(value)
```

PARAMETERS:

* **value** (`array`)

Values to match with the filter.
### fn withOperator

```jsonnet
withOperator(value)
```

PARAMETERS:

* **value** (`string`)

String denoting the filter operation. Supports 'eq' - equals,'ne' - not equals, 'sw' - starts with. Note that some dimensions may not support all operators.