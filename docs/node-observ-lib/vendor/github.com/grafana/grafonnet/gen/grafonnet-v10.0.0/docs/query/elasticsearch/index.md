# elasticsearch

grafonnet.query.elasticsearch

## Subpackages

* [bucketAggs](bucketAggs/index.md)
* [metrics](metrics/index.md)

## Index

* [`fn withAlias(value)`](#fn-withalias)
* [`fn withBucketAggs(value)`](#fn-withbucketaggs)
* [`fn withBucketAggsMixin(value)`](#fn-withbucketaggsmixin)
* [`fn withDatasource(value)`](#fn-withdatasource)
* [`fn withHide(value=true)`](#fn-withhide)
* [`fn withMetrics(value)`](#fn-withmetrics)
* [`fn withMetricsMixin(value)`](#fn-withmetricsmixin)
* [`fn withQuery(value)`](#fn-withquery)
* [`fn withQueryType(value)`](#fn-withquerytype)
* [`fn withRefId(value)`](#fn-withrefid)
* [`fn withTimeField(value)`](#fn-withtimefield)

## Fields

### fn withAlias

```jsonnet
withAlias(value)
```

PARAMETERS:

* **value** (`string`)

Alias pattern
### fn withBucketAggs

```jsonnet
withBucketAggs(value)
```

PARAMETERS:

* **value** (`array`)

List of bucket aggregations
### fn withBucketAggsMixin

```jsonnet
withBucketAggsMixin(value)
```

PARAMETERS:

* **value** (`array`)

List of bucket aggregations
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
### fn withMetrics

```jsonnet
withMetrics(value)
```

PARAMETERS:

* **value** (`array`)

List of metric aggregations
### fn withMetricsMixin

```jsonnet
withMetricsMixin(value)
```

PARAMETERS:

* **value** (`array`)

List of metric aggregations
### fn withQuery

```jsonnet
withQuery(value)
```

PARAMETERS:

* **value** (`string`)

Lucene query
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
### fn withTimeField

```jsonnet
withTimeField(value)
```

PARAMETERS:

* **value** (`string`)

Name of time field