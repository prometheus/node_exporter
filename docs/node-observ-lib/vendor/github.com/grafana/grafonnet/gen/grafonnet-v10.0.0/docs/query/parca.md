# parca

grafonnet.query.parca

## Index

* [`fn withDatasource(value)`](#fn-withdatasource)
* [`fn withHide(value=true)`](#fn-withhide)
* [`fn withLabelSelector(value="{}")`](#fn-withlabelselector)
* [`fn withProfileTypeId(value)`](#fn-withprofiletypeid)
* [`fn withQueryType(value)`](#fn-withquerytype)
* [`fn withRefId(value)`](#fn-withrefid)

## Fields

### fn withDatasource

```jsonnet
withDatasource(value)
```

PARAMETERS:

* **value** (`string`)

For mixed data sources the selected datasource is on the query level.
For non mixed scenarios this is undefined.
TODO find a better way to do this ^ that's friendly to schema
TODO this shouldn't be unknown but DataSourceRef | null
### fn withHide

```jsonnet
withHide(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

true if query is disabled (ie should not be returned to the dashboard)
Note this does not always imply that the query should not be executed since
the results from a hidden query may be used as the input to other queries (SSE etc)
### fn withLabelSelector

```jsonnet
withLabelSelector(value="{}")
```

PARAMETERS:

* **value** (`string`)
   - default value: `"{}"`

Specifies the query label selectors.
### fn withProfileTypeId

```jsonnet
withProfileTypeId(value)
```

PARAMETERS:

* **value** (`string`)

Specifies the type of profile to query.
### fn withQueryType

```jsonnet
withQueryType(value)
```

PARAMETERS:

* **value** (`string`)

Specify the query flavor
TODO make this required and give it a default
### fn withRefId

```jsonnet
withRefId(value)
```

PARAMETERS:

* **value** (`string`)

A unique identifier for the query within the list of targets.
In server side expressions, the refId is used as a variable name to identify results.
By default, the UI will assign A->Z; however setting meaningful names may be useful.