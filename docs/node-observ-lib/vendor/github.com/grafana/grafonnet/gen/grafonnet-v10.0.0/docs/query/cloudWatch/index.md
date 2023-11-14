# cloudWatch

grafonnet.query.cloudWatch

## Subpackages

* [CloudWatchLogsQuery.logGroups](CloudWatchLogsQuery/logGroups.md)
* [CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.parameters](CloudWatchMetricsQuery/sql/from/QueryEditorFunctionExpression/parameters.md)
* [CloudWatchMetricsQuery.sql.orderBy.parameters](CloudWatchMetricsQuery/sql/orderBy/parameters.md)
* [CloudWatchMetricsQuery.sql.select.parameters](CloudWatchMetricsQuery/sql/select/parameters.md)

## Index

* [`obj CloudWatchAnnotationQuery`](#obj-cloudwatchannotationquery)
  * [`fn withAccountId(value)`](#fn-cloudwatchannotationquerywithaccountid)
  * [`fn withActionPrefix(value)`](#fn-cloudwatchannotationquerywithactionprefix)
  * [`fn withAlarmNamePrefix(value)`](#fn-cloudwatchannotationquerywithalarmnameprefix)
  * [`fn withDatasource(value)`](#fn-cloudwatchannotationquerywithdatasource)
  * [`fn withDimensions(value)`](#fn-cloudwatchannotationquerywithdimensions)
  * [`fn withDimensionsMixin(value)`](#fn-cloudwatchannotationquerywithdimensionsmixin)
  * [`fn withHide(value=true)`](#fn-cloudwatchannotationquerywithhide)
  * [`fn withMatchExact(value=true)`](#fn-cloudwatchannotationquerywithmatchexact)
  * [`fn withMetricName(value)`](#fn-cloudwatchannotationquerywithmetricname)
  * [`fn withNamespace(value)`](#fn-cloudwatchannotationquerywithnamespace)
  * [`fn withPeriod(value)`](#fn-cloudwatchannotationquerywithperiod)
  * [`fn withPrefixMatching(value=true)`](#fn-cloudwatchannotationquerywithprefixmatching)
  * [`fn withQueryMode(value)`](#fn-cloudwatchannotationquerywithquerymode)
  * [`fn withQueryType(value)`](#fn-cloudwatchannotationquerywithquerytype)
  * [`fn withRefId(value)`](#fn-cloudwatchannotationquerywithrefid)
  * [`fn withRegion(value)`](#fn-cloudwatchannotationquerywithregion)
  * [`fn withStatistic(value)`](#fn-cloudwatchannotationquerywithstatistic)
  * [`fn withStatistics(value)`](#fn-cloudwatchannotationquerywithstatistics)
  * [`fn withStatisticsMixin(value)`](#fn-cloudwatchannotationquerywithstatisticsmixin)
* [`obj CloudWatchLogsQuery`](#obj-cloudwatchlogsquery)
  * [`fn withDatasource(value)`](#fn-cloudwatchlogsquerywithdatasource)
  * [`fn withExpression(value)`](#fn-cloudwatchlogsquerywithexpression)
  * [`fn withHide(value=true)`](#fn-cloudwatchlogsquerywithhide)
  * [`fn withId(value)`](#fn-cloudwatchlogsquerywithid)
  * [`fn withLogGroupNames(value)`](#fn-cloudwatchlogsquerywithloggroupnames)
  * [`fn withLogGroupNamesMixin(value)`](#fn-cloudwatchlogsquerywithloggroupnamesmixin)
  * [`fn withLogGroups(value)`](#fn-cloudwatchlogsquerywithloggroups)
  * [`fn withLogGroupsMixin(value)`](#fn-cloudwatchlogsquerywithloggroupsmixin)
  * [`fn withQueryMode(value)`](#fn-cloudwatchlogsquerywithquerymode)
  * [`fn withQueryType(value)`](#fn-cloudwatchlogsquerywithquerytype)
  * [`fn withRefId(value)`](#fn-cloudwatchlogsquerywithrefid)
  * [`fn withRegion(value)`](#fn-cloudwatchlogsquerywithregion)
  * [`fn withStatsGroups(value)`](#fn-cloudwatchlogsquerywithstatsgroups)
  * [`fn withStatsGroupsMixin(value)`](#fn-cloudwatchlogsquerywithstatsgroupsmixin)
* [`obj CloudWatchMetricsQuery`](#obj-cloudwatchmetricsquery)
  * [`fn withAccountId(value)`](#fn-cloudwatchmetricsquerywithaccountid)
  * [`fn withAlias(value)`](#fn-cloudwatchmetricsquerywithalias)
  * [`fn withDatasource(value)`](#fn-cloudwatchmetricsquerywithdatasource)
  * [`fn withDimensions(value)`](#fn-cloudwatchmetricsquerywithdimensions)
  * [`fn withDimensionsMixin(value)`](#fn-cloudwatchmetricsquerywithdimensionsmixin)
  * [`fn withExpression(value)`](#fn-cloudwatchmetricsquerywithexpression)
  * [`fn withHide(value=true)`](#fn-cloudwatchmetricsquerywithhide)
  * [`fn withId(value)`](#fn-cloudwatchmetricsquerywithid)
  * [`fn withLabel(value)`](#fn-cloudwatchmetricsquerywithlabel)
  * [`fn withMatchExact(value=true)`](#fn-cloudwatchmetricsquerywithmatchexact)
  * [`fn withMetricEditorMode(value)`](#fn-cloudwatchmetricsquerywithmetriceditormode)
  * [`fn withMetricName(value)`](#fn-cloudwatchmetricsquerywithmetricname)
  * [`fn withMetricQueryType(value)`](#fn-cloudwatchmetricsquerywithmetricquerytype)
  * [`fn withNamespace(value)`](#fn-cloudwatchmetricsquerywithnamespace)
  * [`fn withPeriod(value)`](#fn-cloudwatchmetricsquerywithperiod)
  * [`fn withQueryMode(value)`](#fn-cloudwatchmetricsquerywithquerymode)
  * [`fn withQueryType(value)`](#fn-cloudwatchmetricsquerywithquerytype)
  * [`fn withRefId(value)`](#fn-cloudwatchmetricsquerywithrefid)
  * [`fn withRegion(value)`](#fn-cloudwatchmetricsquerywithregion)
  * [`fn withSql(value)`](#fn-cloudwatchmetricsquerywithsql)
  * [`fn withSqlExpression(value)`](#fn-cloudwatchmetricsquerywithsqlexpression)
  * [`fn withSqlMixin(value)`](#fn-cloudwatchmetricsquerywithsqlmixin)
  * [`fn withStatistic(value)`](#fn-cloudwatchmetricsquerywithstatistic)
  * [`fn withStatistics(value)`](#fn-cloudwatchmetricsquerywithstatistics)
  * [`fn withStatisticsMixin(value)`](#fn-cloudwatchmetricsquerywithstatisticsmixin)
  * [`obj sql`](#obj-cloudwatchmetricsquerysql)
    * [`fn withFrom(value)`](#fn-cloudwatchmetricsquerysqlwithfrom)
    * [`fn withFromMixin(value)`](#fn-cloudwatchmetricsquerysqlwithfrommixin)
    * [`fn withGroupBy(value)`](#fn-cloudwatchmetricsquerysqlwithgroupby)
    * [`fn withGroupByMixin(value)`](#fn-cloudwatchmetricsquerysqlwithgroupbymixin)
    * [`fn withLimit(value)`](#fn-cloudwatchmetricsquerysqlwithlimit)
    * [`fn withOrderBy(value)`](#fn-cloudwatchmetricsquerysqlwithorderby)
    * [`fn withOrderByDirection(value)`](#fn-cloudwatchmetricsquerysqlwithorderbydirection)
    * [`fn withOrderByMixin(value)`](#fn-cloudwatchmetricsquerysqlwithorderbymixin)
    * [`fn withSelect(value)`](#fn-cloudwatchmetricsquerysqlwithselect)
    * [`fn withSelectMixin(value)`](#fn-cloudwatchmetricsquerysqlwithselectmixin)
    * [`fn withWhere(value)`](#fn-cloudwatchmetricsquerysqlwithwhere)
    * [`fn withWhereMixin(value)`](#fn-cloudwatchmetricsquerysqlwithwheremixin)
    * [`obj from`](#obj-cloudwatchmetricsquerysqlfrom)
      * [`fn withQueryEditorFunctionExpression(value)`](#fn-cloudwatchmetricsquerysqlfromwithqueryeditorfunctionexpression)
      * [`fn withQueryEditorFunctionExpressionMixin(value)`](#fn-cloudwatchmetricsquerysqlfromwithqueryeditorfunctionexpressionmixin)
      * [`fn withQueryEditorPropertyExpression(value)`](#fn-cloudwatchmetricsquerysqlfromwithqueryeditorpropertyexpression)
      * [`fn withQueryEditorPropertyExpressionMixin(value)`](#fn-cloudwatchmetricsquerysqlfromwithqueryeditorpropertyexpressionmixin)
      * [`obj QueryEditorFunctionExpression`](#obj-cloudwatchmetricsquerysqlfromqueryeditorfunctionexpression)
        * [`fn withName(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorfunctionexpressionwithname)
        * [`fn withParameters(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorfunctionexpressionwithparameters)
        * [`fn withParametersMixin(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorfunctionexpressionwithparametersmixin)
        * [`fn withType(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorfunctionexpressionwithtype)
      * [`obj QueryEditorPropertyExpression`](#obj-cloudwatchmetricsquerysqlfromqueryeditorpropertyexpression)
        * [`fn withProperty(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorpropertyexpressionwithproperty)
        * [`fn withPropertyMixin(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorpropertyexpressionwithpropertymixin)
        * [`fn withType(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorpropertyexpressionwithtype)
        * [`obj property`](#obj-cloudwatchmetricsquerysqlfromqueryeditorpropertyexpressionproperty)
          * [`fn withName(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorpropertyexpressionpropertywithname)
          * [`fn withType(value)`](#fn-cloudwatchmetricsquerysqlfromqueryeditorpropertyexpressionpropertywithtype)
    * [`obj groupBy`](#obj-cloudwatchmetricsquerysqlgroupby)
      * [`fn withExpressions(value)`](#fn-cloudwatchmetricsquerysqlgroupbywithexpressions)
      * [`fn withExpressionsMixin(value)`](#fn-cloudwatchmetricsquerysqlgroupbywithexpressionsmixin)
      * [`fn withType(value)`](#fn-cloudwatchmetricsquerysqlgroupbywithtype)
    * [`obj orderBy`](#obj-cloudwatchmetricsquerysqlorderby)
      * [`fn withName(value)`](#fn-cloudwatchmetricsquerysqlorderbywithname)
      * [`fn withParameters(value)`](#fn-cloudwatchmetricsquerysqlorderbywithparameters)
      * [`fn withParametersMixin(value)`](#fn-cloudwatchmetricsquerysqlorderbywithparametersmixin)
      * [`fn withType(value)`](#fn-cloudwatchmetricsquerysqlorderbywithtype)
    * [`obj select`](#obj-cloudwatchmetricsquerysqlselect)
      * [`fn withName(value)`](#fn-cloudwatchmetricsquerysqlselectwithname)
      * [`fn withParameters(value)`](#fn-cloudwatchmetricsquerysqlselectwithparameters)
      * [`fn withParametersMixin(value)`](#fn-cloudwatchmetricsquerysqlselectwithparametersmixin)
      * [`fn withType(value)`](#fn-cloudwatchmetricsquerysqlselectwithtype)
    * [`obj where`](#obj-cloudwatchmetricsquerysqlwhere)
      * [`fn withExpressions(value)`](#fn-cloudwatchmetricsquerysqlwherewithexpressions)
      * [`fn withExpressionsMixin(value)`](#fn-cloudwatchmetricsquerysqlwherewithexpressionsmixin)
      * [`fn withType(value)`](#fn-cloudwatchmetricsquerysqlwherewithtype)

## Fields

### obj CloudWatchAnnotationQuery


#### fn CloudWatchAnnotationQuery.withAccountId

```jsonnet
CloudWatchAnnotationQuery.withAccountId(value)
```

PARAMETERS:

* **value** (`string`)

The ID of the AWS account to query for the metric, specifying `all` will query all accounts that the monitoring account is permitted to query.
#### fn CloudWatchAnnotationQuery.withActionPrefix

```jsonnet
CloudWatchAnnotationQuery.withActionPrefix(value)
```

PARAMETERS:

* **value** (`string`)

Use this parameter to filter the results of the operation to only those alarms
that use a certain alarm action. For example, you could specify the ARN of
an SNS topic to find all alarms that send notifications to that topic.
e.g. `arn:aws:sns:us-east-1:123456789012:my-app-` would match `arn:aws:sns:us-east-1:123456789012:my-app-action`
but not match `arn:aws:sns:us-east-1:123456789012:your-app-action`
#### fn CloudWatchAnnotationQuery.withAlarmNamePrefix

```jsonnet
CloudWatchAnnotationQuery.withAlarmNamePrefix(value)
```

PARAMETERS:

* **value** (`string`)

An alarm name prefix. If you specify this parameter, you receive information
about all alarms that have names that start with this prefix.
e.g. `my-team-service-` would match `my-team-service-high-cpu` but not match `your-team-service-high-cpu`
#### fn CloudWatchAnnotationQuery.withDatasource

```jsonnet
CloudWatchAnnotationQuery.withDatasource(value)
```

PARAMETERS:

* **value** (`string`)

For mixed data sources the selected datasource is on the query level.
For non mixed scenarios this is undefined.
TODO find a better way to do this ^ that's friendly to schema
TODO this shouldn't be unknown but DataSourceRef | null
#### fn CloudWatchAnnotationQuery.withDimensions

```jsonnet
CloudWatchAnnotationQuery.withDimensions(value)
```

PARAMETERS:

* **value** (`object`)

A name/value pair that is part of the identity of a metric. For example, you can get statistics for a specific EC2 instance by specifying the InstanceId dimension when you search for metrics.
#### fn CloudWatchAnnotationQuery.withDimensionsMixin

```jsonnet
CloudWatchAnnotationQuery.withDimensionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

A name/value pair that is part of the identity of a metric. For example, you can get statistics for a specific EC2 instance by specifying the InstanceId dimension when you search for metrics.
#### fn CloudWatchAnnotationQuery.withHide

```jsonnet
CloudWatchAnnotationQuery.withHide(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

true if query is disabled (ie should not be returned to the dashboard)
Note this does not always imply that the query should not be executed since
the results from a hidden query may be used as the input to other queries (SSE etc)
#### fn CloudWatchAnnotationQuery.withMatchExact

```jsonnet
CloudWatchAnnotationQuery.withMatchExact(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Only show metrics that exactly match all defined dimension names.
#### fn CloudWatchAnnotationQuery.withMetricName

```jsonnet
CloudWatchAnnotationQuery.withMetricName(value)
```

PARAMETERS:

* **value** (`string`)

Name of the metric
#### fn CloudWatchAnnotationQuery.withNamespace

```jsonnet
CloudWatchAnnotationQuery.withNamespace(value)
```

PARAMETERS:

* **value** (`string`)

A namespace is a container for CloudWatch metrics. Metrics in different namespaces are isolated from each other, so that metrics from different applications are not mistakenly aggregated into the same statistics. For example, Amazon EC2 uses the AWS/EC2 namespace.
#### fn CloudWatchAnnotationQuery.withPeriod

```jsonnet
CloudWatchAnnotationQuery.withPeriod(value)
```

PARAMETERS:

* **value** (`string`)

The length of time associated with a specific Amazon CloudWatch statistic. Can be specified by a number of seconds, 'auto', or as a duration string e.g. '15m' being 15 minutes
#### fn CloudWatchAnnotationQuery.withPrefixMatching

```jsonnet
CloudWatchAnnotationQuery.withPrefixMatching(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Enable matching on the prefix of the action name or alarm name, specify the prefixes with actionPrefix and/or alarmNamePrefix
#### fn CloudWatchAnnotationQuery.withQueryMode

```jsonnet
CloudWatchAnnotationQuery.withQueryMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"Metrics"`, `"Logs"`, `"Annotations"`


#### fn CloudWatchAnnotationQuery.withQueryType

```jsonnet
CloudWatchAnnotationQuery.withQueryType(value)
```

PARAMETERS:

* **value** (`string`)

Specify the query flavor
TODO make this required and give it a default
#### fn CloudWatchAnnotationQuery.withRefId

```jsonnet
CloudWatchAnnotationQuery.withRefId(value)
```

PARAMETERS:

* **value** (`string`)

A unique identifier for the query within the list of targets.
In server side expressions, the refId is used as a variable name to identify results.
By default, the UI will assign A->Z; however setting meaningful names may be useful.
#### fn CloudWatchAnnotationQuery.withRegion

```jsonnet
CloudWatchAnnotationQuery.withRegion(value)
```

PARAMETERS:

* **value** (`string`)

AWS region to query for the metric
#### fn CloudWatchAnnotationQuery.withStatistic

```jsonnet
CloudWatchAnnotationQuery.withStatistic(value)
```

PARAMETERS:

* **value** (`string`)

Metric data aggregations over specified periods of time. For detailed definitions of the statistics supported by CloudWatch, see https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Statistics-definitions.html.
#### fn CloudWatchAnnotationQuery.withStatistics

```jsonnet
CloudWatchAnnotationQuery.withStatistics(value)
```

PARAMETERS:

* **value** (`array`)

@deprecated use statistic
#### fn CloudWatchAnnotationQuery.withStatisticsMixin

```jsonnet
CloudWatchAnnotationQuery.withStatisticsMixin(value)
```

PARAMETERS:

* **value** (`array`)

@deprecated use statistic
### obj CloudWatchLogsQuery


#### fn CloudWatchLogsQuery.withDatasource

```jsonnet
CloudWatchLogsQuery.withDatasource(value)
```

PARAMETERS:

* **value** (`string`)

For mixed data sources the selected datasource is on the query level.
For non mixed scenarios this is undefined.
TODO find a better way to do this ^ that's friendly to schema
TODO this shouldn't be unknown but DataSourceRef | null
#### fn CloudWatchLogsQuery.withExpression

```jsonnet
CloudWatchLogsQuery.withExpression(value)
```

PARAMETERS:

* **value** (`string`)

The CloudWatch Logs Insights query to execute
#### fn CloudWatchLogsQuery.withHide

```jsonnet
CloudWatchLogsQuery.withHide(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

true if query is disabled (ie should not be returned to the dashboard)
Note this does not always imply that the query should not be executed since
the results from a hidden query may be used as the input to other queries (SSE etc)
#### fn CloudWatchLogsQuery.withId

```jsonnet
CloudWatchLogsQuery.withId(value)
```

PARAMETERS:

* **value** (`string`)


#### fn CloudWatchLogsQuery.withLogGroupNames

```jsonnet
CloudWatchLogsQuery.withLogGroupNames(value)
```

PARAMETERS:

* **value** (`array`)

@deprecated use logGroups
#### fn CloudWatchLogsQuery.withLogGroupNamesMixin

```jsonnet
CloudWatchLogsQuery.withLogGroupNamesMixin(value)
```

PARAMETERS:

* **value** (`array`)

@deprecated use logGroups
#### fn CloudWatchLogsQuery.withLogGroups

```jsonnet
CloudWatchLogsQuery.withLogGroups(value)
```

PARAMETERS:

* **value** (`array`)

Log groups to query
#### fn CloudWatchLogsQuery.withLogGroupsMixin

```jsonnet
CloudWatchLogsQuery.withLogGroupsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Log groups to query
#### fn CloudWatchLogsQuery.withQueryMode

```jsonnet
CloudWatchLogsQuery.withQueryMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"Metrics"`, `"Logs"`, `"Annotations"`


#### fn CloudWatchLogsQuery.withQueryType

```jsonnet
CloudWatchLogsQuery.withQueryType(value)
```

PARAMETERS:

* **value** (`string`)

Specify the query flavor
TODO make this required and give it a default
#### fn CloudWatchLogsQuery.withRefId

```jsonnet
CloudWatchLogsQuery.withRefId(value)
```

PARAMETERS:

* **value** (`string`)

A unique identifier for the query within the list of targets.
In server side expressions, the refId is used as a variable name to identify results.
By default, the UI will assign A->Z; however setting meaningful names may be useful.
#### fn CloudWatchLogsQuery.withRegion

```jsonnet
CloudWatchLogsQuery.withRegion(value)
```

PARAMETERS:

* **value** (`string`)

AWS region to query for the logs
#### fn CloudWatchLogsQuery.withStatsGroups

```jsonnet
CloudWatchLogsQuery.withStatsGroups(value)
```

PARAMETERS:

* **value** (`array`)

Fields to group the results by, this field is automatically populated whenever the query is updated
#### fn CloudWatchLogsQuery.withStatsGroupsMixin

```jsonnet
CloudWatchLogsQuery.withStatsGroupsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Fields to group the results by, this field is automatically populated whenever the query is updated
### obj CloudWatchMetricsQuery


#### fn CloudWatchMetricsQuery.withAccountId

```jsonnet
CloudWatchMetricsQuery.withAccountId(value)
```

PARAMETERS:

* **value** (`string`)

The ID of the AWS account to query for the metric, specifying `all` will query all accounts that the monitoring account is permitted to query.
#### fn CloudWatchMetricsQuery.withAlias

```jsonnet
CloudWatchMetricsQuery.withAlias(value)
```

PARAMETERS:

* **value** (`string`)

Deprecated: use label
@deprecated use label
#### fn CloudWatchMetricsQuery.withDatasource

```jsonnet
CloudWatchMetricsQuery.withDatasource(value)
```

PARAMETERS:

* **value** (`string`)

For mixed data sources the selected datasource is on the query level.
For non mixed scenarios this is undefined.
TODO find a better way to do this ^ that's friendly to schema
TODO this shouldn't be unknown but DataSourceRef | null
#### fn CloudWatchMetricsQuery.withDimensions

```jsonnet
CloudWatchMetricsQuery.withDimensions(value)
```

PARAMETERS:

* **value** (`object`)

A name/value pair that is part of the identity of a metric. For example, you can get statistics for a specific EC2 instance by specifying the InstanceId dimension when you search for metrics.
#### fn CloudWatchMetricsQuery.withDimensionsMixin

```jsonnet
CloudWatchMetricsQuery.withDimensionsMixin(value)
```

PARAMETERS:

* **value** (`object`)

A name/value pair that is part of the identity of a metric. For example, you can get statistics for a specific EC2 instance by specifying the InstanceId dimension when you search for metrics.
#### fn CloudWatchMetricsQuery.withExpression

```jsonnet
CloudWatchMetricsQuery.withExpression(value)
```

PARAMETERS:

* **value** (`string`)

Math expression query
#### fn CloudWatchMetricsQuery.withHide

```jsonnet
CloudWatchMetricsQuery.withHide(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

true if query is disabled (ie should not be returned to the dashboard)
Note this does not always imply that the query should not be executed since
the results from a hidden query may be used as the input to other queries (SSE etc)
#### fn CloudWatchMetricsQuery.withId

```jsonnet
CloudWatchMetricsQuery.withId(value)
```

PARAMETERS:

* **value** (`string`)

ID can be used to reference other queries in math expressions. The ID can include numbers, letters, and underscore, and must start with a lowercase letter.
#### fn CloudWatchMetricsQuery.withLabel

```jsonnet
CloudWatchMetricsQuery.withLabel(value)
```

PARAMETERS:

* **value** (`string`)

Change the time series legend names using dynamic labels. See https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/graph-dynamic-labels.html for more details.
#### fn CloudWatchMetricsQuery.withMatchExact

```jsonnet
CloudWatchMetricsQuery.withMatchExact(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

Only show metrics that exactly match all defined dimension names.
#### fn CloudWatchMetricsQuery.withMetricEditorMode

```jsonnet
CloudWatchMetricsQuery.withMetricEditorMode(value)
```

PARAMETERS:

* **value** (`integer`)
   - valid values: `0`, `1`


#### fn CloudWatchMetricsQuery.withMetricName

```jsonnet
CloudWatchMetricsQuery.withMetricName(value)
```

PARAMETERS:

* **value** (`string`)

Name of the metric
#### fn CloudWatchMetricsQuery.withMetricQueryType

```jsonnet
CloudWatchMetricsQuery.withMetricQueryType(value)
```

PARAMETERS:

* **value** (`integer`)
   - valid values: `0`, `1`


#### fn CloudWatchMetricsQuery.withNamespace

```jsonnet
CloudWatchMetricsQuery.withNamespace(value)
```

PARAMETERS:

* **value** (`string`)

A namespace is a container for CloudWatch metrics. Metrics in different namespaces are isolated from each other, so that metrics from different applications are not mistakenly aggregated into the same statistics. For example, Amazon EC2 uses the AWS/EC2 namespace.
#### fn CloudWatchMetricsQuery.withPeriod

```jsonnet
CloudWatchMetricsQuery.withPeriod(value)
```

PARAMETERS:

* **value** (`string`)

The length of time associated with a specific Amazon CloudWatch statistic. Can be specified by a number of seconds, 'auto', or as a duration string e.g. '15m' being 15 minutes
#### fn CloudWatchMetricsQuery.withQueryMode

```jsonnet
CloudWatchMetricsQuery.withQueryMode(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"Metrics"`, `"Logs"`, `"Annotations"`


#### fn CloudWatchMetricsQuery.withQueryType

```jsonnet
CloudWatchMetricsQuery.withQueryType(value)
```

PARAMETERS:

* **value** (`string`)

Specify the query flavor
TODO make this required and give it a default
#### fn CloudWatchMetricsQuery.withRefId

```jsonnet
CloudWatchMetricsQuery.withRefId(value)
```

PARAMETERS:

* **value** (`string`)

A unique identifier for the query within the list of targets.
In server side expressions, the refId is used as a variable name to identify results.
By default, the UI will assign A->Z; however setting meaningful names may be useful.
#### fn CloudWatchMetricsQuery.withRegion

```jsonnet
CloudWatchMetricsQuery.withRegion(value)
```

PARAMETERS:

* **value** (`string`)

AWS region to query for the metric
#### fn CloudWatchMetricsQuery.withSql

```jsonnet
CloudWatchMetricsQuery.withSql(value)
```

PARAMETERS:

* **value** (`object`)


#### fn CloudWatchMetricsQuery.withSqlExpression

```jsonnet
CloudWatchMetricsQuery.withSqlExpression(value)
```

PARAMETERS:

* **value** (`string`)

When the metric query type is `metricQueryType` is set to `Query`, this field is used to specify the query string.
#### fn CloudWatchMetricsQuery.withSqlMixin

```jsonnet
CloudWatchMetricsQuery.withSqlMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn CloudWatchMetricsQuery.withStatistic

```jsonnet
CloudWatchMetricsQuery.withStatistic(value)
```

PARAMETERS:

* **value** (`string`)

Metric data aggregations over specified periods of time. For detailed definitions of the statistics supported by CloudWatch, see https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Statistics-definitions.html.
#### fn CloudWatchMetricsQuery.withStatistics

```jsonnet
CloudWatchMetricsQuery.withStatistics(value)
```

PARAMETERS:

* **value** (`array`)

@deprecated use statistic
#### fn CloudWatchMetricsQuery.withStatisticsMixin

```jsonnet
CloudWatchMetricsQuery.withStatisticsMixin(value)
```

PARAMETERS:

* **value** (`array`)

@deprecated use statistic
#### obj CloudWatchMetricsQuery.sql


##### fn CloudWatchMetricsQuery.sql.withFrom

```jsonnet
CloudWatchMetricsQuery.sql.withFrom(value)
```

PARAMETERS:

* **value** (`object`)

FROM part of the SQL expression
##### fn CloudWatchMetricsQuery.sql.withFromMixin

```jsonnet
CloudWatchMetricsQuery.sql.withFromMixin(value)
```

PARAMETERS:

* **value** (`object`)

FROM part of the SQL expression
##### fn CloudWatchMetricsQuery.sql.withGroupBy

```jsonnet
CloudWatchMetricsQuery.sql.withGroupBy(value)
```

PARAMETERS:

* **value** (`object`)


##### fn CloudWatchMetricsQuery.sql.withGroupByMixin

```jsonnet
CloudWatchMetricsQuery.sql.withGroupByMixin(value)
```

PARAMETERS:

* **value** (`object`)


##### fn CloudWatchMetricsQuery.sql.withLimit

```jsonnet
CloudWatchMetricsQuery.sql.withLimit(value)
```

PARAMETERS:

* **value** (`integer`)

LIMIT part of the SQL expression
##### fn CloudWatchMetricsQuery.sql.withOrderBy

```jsonnet
CloudWatchMetricsQuery.sql.withOrderBy(value)
```

PARAMETERS:

* **value** (`object`)


##### fn CloudWatchMetricsQuery.sql.withOrderByDirection

```jsonnet
CloudWatchMetricsQuery.sql.withOrderByDirection(value)
```

PARAMETERS:

* **value** (`string`)

The sort order of the SQL expression, `ASC` or `DESC`
##### fn CloudWatchMetricsQuery.sql.withOrderByMixin

```jsonnet
CloudWatchMetricsQuery.sql.withOrderByMixin(value)
```

PARAMETERS:

* **value** (`object`)


##### fn CloudWatchMetricsQuery.sql.withSelect

```jsonnet
CloudWatchMetricsQuery.sql.withSelect(value)
```

PARAMETERS:

* **value** (`object`)


##### fn CloudWatchMetricsQuery.sql.withSelectMixin

```jsonnet
CloudWatchMetricsQuery.sql.withSelectMixin(value)
```

PARAMETERS:

* **value** (`object`)


##### fn CloudWatchMetricsQuery.sql.withWhere

```jsonnet
CloudWatchMetricsQuery.sql.withWhere(value)
```

PARAMETERS:

* **value** (`object`)


##### fn CloudWatchMetricsQuery.sql.withWhereMixin

```jsonnet
CloudWatchMetricsQuery.sql.withWhereMixin(value)
```

PARAMETERS:

* **value** (`object`)


##### obj CloudWatchMetricsQuery.sql.from


###### fn CloudWatchMetricsQuery.sql.from.withQueryEditorFunctionExpression

```jsonnet
CloudWatchMetricsQuery.sql.from.withQueryEditorFunctionExpression(value)
```

PARAMETERS:

* **value** (`object`)


###### fn CloudWatchMetricsQuery.sql.from.withQueryEditorFunctionExpressionMixin

```jsonnet
CloudWatchMetricsQuery.sql.from.withQueryEditorFunctionExpressionMixin(value)
```

PARAMETERS:

* **value** (`object`)


###### fn CloudWatchMetricsQuery.sql.from.withQueryEditorPropertyExpression

```jsonnet
CloudWatchMetricsQuery.sql.from.withQueryEditorPropertyExpression(value)
```

PARAMETERS:

* **value** (`object`)


###### fn CloudWatchMetricsQuery.sql.from.withQueryEditorPropertyExpressionMixin

```jsonnet
CloudWatchMetricsQuery.sql.from.withQueryEditorPropertyExpressionMixin(value)
```

PARAMETERS:

* **value** (`object`)


###### obj CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression


####### fn CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.withName

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.withName(value)
```

PARAMETERS:

* **value** (`string`)


####### fn CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.withParameters

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.withParameters(value)
```

PARAMETERS:

* **value** (`array`)


####### fn CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.withParametersMixin

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.withParametersMixin(value)
```

PARAMETERS:

* **value** (`array`)


####### fn CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.withType

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorFunctionExpression.withType(value)
```

PARAMETERS:

* **value** (`string`)


###### obj CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression


####### fn CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.withProperty

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.withProperty(value)
```

PARAMETERS:

* **value** (`object`)


####### fn CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.withPropertyMixin

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.withPropertyMixin(value)
```

PARAMETERS:

* **value** (`object`)


####### fn CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.withType

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.withType(value)
```

PARAMETERS:

* **value** (`string`)


####### obj CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.property


######## fn CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.property.withName

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.property.withName(value)
```

PARAMETERS:

* **value** (`string`)


######## fn CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.property.withType

```jsonnet
CloudWatchMetricsQuery.sql.from.QueryEditorPropertyExpression.property.withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"string"`


##### obj CloudWatchMetricsQuery.sql.groupBy


###### fn CloudWatchMetricsQuery.sql.groupBy.withExpressions

```jsonnet
CloudWatchMetricsQuery.sql.groupBy.withExpressions(value)
```

PARAMETERS:

* **value** (`array`)


###### fn CloudWatchMetricsQuery.sql.groupBy.withExpressionsMixin

```jsonnet
CloudWatchMetricsQuery.sql.groupBy.withExpressionsMixin(value)
```

PARAMETERS:

* **value** (`array`)


###### fn CloudWatchMetricsQuery.sql.groupBy.withType

```jsonnet
CloudWatchMetricsQuery.sql.groupBy.withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"and"`, `"or"`


##### obj CloudWatchMetricsQuery.sql.orderBy


###### fn CloudWatchMetricsQuery.sql.orderBy.withName

```jsonnet
CloudWatchMetricsQuery.sql.orderBy.withName(value)
```

PARAMETERS:

* **value** (`string`)


###### fn CloudWatchMetricsQuery.sql.orderBy.withParameters

```jsonnet
CloudWatchMetricsQuery.sql.orderBy.withParameters(value)
```

PARAMETERS:

* **value** (`array`)


###### fn CloudWatchMetricsQuery.sql.orderBy.withParametersMixin

```jsonnet
CloudWatchMetricsQuery.sql.orderBy.withParametersMixin(value)
```

PARAMETERS:

* **value** (`array`)


###### fn CloudWatchMetricsQuery.sql.orderBy.withType

```jsonnet
CloudWatchMetricsQuery.sql.orderBy.withType(value)
```

PARAMETERS:

* **value** (`string`)


##### obj CloudWatchMetricsQuery.sql.select


###### fn CloudWatchMetricsQuery.sql.select.withName

```jsonnet
CloudWatchMetricsQuery.sql.select.withName(value)
```

PARAMETERS:

* **value** (`string`)


###### fn CloudWatchMetricsQuery.sql.select.withParameters

```jsonnet
CloudWatchMetricsQuery.sql.select.withParameters(value)
```

PARAMETERS:

* **value** (`array`)


###### fn CloudWatchMetricsQuery.sql.select.withParametersMixin

```jsonnet
CloudWatchMetricsQuery.sql.select.withParametersMixin(value)
```

PARAMETERS:

* **value** (`array`)


###### fn CloudWatchMetricsQuery.sql.select.withType

```jsonnet
CloudWatchMetricsQuery.sql.select.withType(value)
```

PARAMETERS:

* **value** (`string`)


##### obj CloudWatchMetricsQuery.sql.where


###### fn CloudWatchMetricsQuery.sql.where.withExpressions

```jsonnet
CloudWatchMetricsQuery.sql.where.withExpressions(value)
```

PARAMETERS:

* **value** (`array`)


###### fn CloudWatchMetricsQuery.sql.where.withExpressionsMixin

```jsonnet
CloudWatchMetricsQuery.sql.where.withExpressionsMixin(value)
```

PARAMETERS:

* **value** (`array`)


###### fn CloudWatchMetricsQuery.sql.where.withType

```jsonnet
CloudWatchMetricsQuery.sql.where.withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"and"`, `"or"`

