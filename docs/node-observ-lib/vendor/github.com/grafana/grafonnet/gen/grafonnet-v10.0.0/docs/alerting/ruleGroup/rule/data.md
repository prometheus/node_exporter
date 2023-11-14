# data



## Index

* [`fn withDatasourceUid(value)`](#fn-withdatasourceuid)
* [`fn withModel(value)`](#fn-withmodel)
* [`fn withModelMixin(value)`](#fn-withmodelmixin)
* [`fn withQueryType(value)`](#fn-withquerytype)
* [`fn withRefId(value)`](#fn-withrefid)
* [`fn withRelativeTimeRange(value)`](#fn-withrelativetimerange)
* [`fn withRelativeTimeRangeMixin(value)`](#fn-withrelativetimerangemixin)
* [`obj relativeTimeRange`](#obj-relativetimerange)
  * [`fn withFrom(value)`](#fn-relativetimerangewithfrom)
  * [`fn withTo(value)`](#fn-relativetimerangewithto)

## Fields

### fn withDatasourceUid

```jsonnet
withDatasourceUid(value)
```

PARAMETERS:

* **value** (`string`)

Grafana data source unique identifier; it should be '__expr__' for a Server Side Expression operation.
### fn withModel

```jsonnet
withModel(value)
```

PARAMETERS:

* **value** (`object`)

JSON is the raw JSON query and includes the above properties as well as custom properties.
### fn withModelMixin

```jsonnet
withModelMixin(value)
```

PARAMETERS:

* **value** (`object`)

JSON is the raw JSON query and includes the above properties as well as custom properties.
### fn withQueryType

```jsonnet
withQueryType(value)
```

PARAMETERS:

* **value** (`string`)

QueryType is an optional identifier for the type of query.
It can be used to distinguish different types of queries.
### fn withRefId

```jsonnet
withRefId(value)
```

PARAMETERS:

* **value** (`string`)

RefID is the unique identifier of the query, set by the frontend call.
### fn withRelativeTimeRange

```jsonnet
withRelativeTimeRange(value)
```

PARAMETERS:

* **value** (`object`)

RelativeTimeRange is the per query start and end time
for requests.
### fn withRelativeTimeRangeMixin

```jsonnet
withRelativeTimeRangeMixin(value)
```

PARAMETERS:

* **value** (`object`)

RelativeTimeRange is the per query start and end time
for requests.
### obj relativeTimeRange


#### fn relativeTimeRange.withFrom

```jsonnet
relativeTimeRange.withFrom(value)
```

PARAMETERS:

* **value** (`integer`)

A Duration represents the elapsed time between two instants
as an int64 nanosecond count. The representation limits the
largest representable duration to approximately 290 years.
#### fn relativeTimeRange.withTo

```jsonnet
relativeTimeRange.withTo(value)
```

PARAMETERS:

* **value** (`integer`)

A Duration represents the elapsed time between two instants
as an int64 nanosecond count. The representation limits the
largest representable duration to approximately 290 years.