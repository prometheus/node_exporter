# azureMonitor

grafonnet.query.azureMonitor

## Subpackages

* [azureMonitor.dimensionFilters](azureMonitor/dimensionFilters.md)
* [azureMonitor.resources](azureMonitor/resources.md)
* [azureTraces.filters](azureTraces/filters.md)

## Index

* [`fn withAzureLogAnalytics(value)`](#fn-withazureloganalytics)
* [`fn withAzureLogAnalyticsMixin(value)`](#fn-withazureloganalyticsmixin)
* [`fn withAzureMonitor(value)`](#fn-withazuremonitor)
* [`fn withAzureMonitorMixin(value)`](#fn-withazuremonitormixin)
* [`fn withAzureResourceGraph(value)`](#fn-withazureresourcegraph)
* [`fn withAzureResourceGraphMixin(value)`](#fn-withazureresourcegraphmixin)
* [`fn withAzureTraces(value)`](#fn-withazuretraces)
* [`fn withAzureTracesMixin(value)`](#fn-withazuretracesmixin)
* [`fn withDatasource(value)`](#fn-withdatasource)
* [`fn withGrafanaTemplateVariableFn(value)`](#fn-withgrafanatemplatevariablefn)
* [`fn withGrafanaTemplateVariableFnMixin(value)`](#fn-withgrafanatemplatevariablefnmixin)
* [`fn withHide(value=true)`](#fn-withhide)
* [`fn withNamespace(value)`](#fn-withnamespace)
* [`fn withQueryType(value)`](#fn-withquerytype)
* [`fn withRefId(value)`](#fn-withrefid)
* [`fn withRegion(value)`](#fn-withregion)
* [`fn withResource(value)`](#fn-withresource)
* [`fn withResourceGroup(value)`](#fn-withresourcegroup)
* [`fn withSubscription(value)`](#fn-withsubscription)
* [`fn withSubscriptions(value)`](#fn-withsubscriptions)
* [`fn withSubscriptionsMixin(value)`](#fn-withsubscriptionsmixin)
* [`obj azureLogAnalytics`](#obj-azureloganalytics)
  * [`fn withQuery(value)`](#fn-azureloganalyticswithquery)
  * [`fn withResource(value)`](#fn-azureloganalyticswithresource)
  * [`fn withResources(value)`](#fn-azureloganalyticswithresources)
  * [`fn withResourcesMixin(value)`](#fn-azureloganalyticswithresourcesmixin)
  * [`fn withResultFormat(value)`](#fn-azureloganalyticswithresultformat)
  * [`fn withWorkspace(value)`](#fn-azureloganalyticswithworkspace)
* [`obj azureMonitor`](#obj-azuremonitor)
  * [`fn withAggregation(value)`](#fn-azuremonitorwithaggregation)
  * [`fn withAlias(value)`](#fn-azuremonitorwithalias)
  * [`fn withAllowedTimeGrainsMs(value)`](#fn-azuremonitorwithallowedtimegrainsms)
  * [`fn withAllowedTimeGrainsMsMixin(value)`](#fn-azuremonitorwithallowedtimegrainsmsmixin)
  * [`fn withCustomNamespace(value)`](#fn-azuremonitorwithcustomnamespace)
  * [`fn withDimension(value)`](#fn-azuremonitorwithdimension)
  * [`fn withDimensionFilter(value)`](#fn-azuremonitorwithdimensionfilter)
  * [`fn withDimensionFilters(value)`](#fn-azuremonitorwithdimensionfilters)
  * [`fn withDimensionFiltersMixin(value)`](#fn-azuremonitorwithdimensionfiltersmixin)
  * [`fn withMetricDefinition(value)`](#fn-azuremonitorwithmetricdefinition)
  * [`fn withMetricName(value)`](#fn-azuremonitorwithmetricname)
  * [`fn withMetricNamespace(value)`](#fn-azuremonitorwithmetricnamespace)
  * [`fn withRegion(value)`](#fn-azuremonitorwithregion)
  * [`fn withResourceGroup(value)`](#fn-azuremonitorwithresourcegroup)
  * [`fn withResourceName(value)`](#fn-azuremonitorwithresourcename)
  * [`fn withResourceUri(value)`](#fn-azuremonitorwithresourceuri)
  * [`fn withResources(value)`](#fn-azuremonitorwithresources)
  * [`fn withResourcesMixin(value)`](#fn-azuremonitorwithresourcesmixin)
  * [`fn withTimeGrain(value)`](#fn-azuremonitorwithtimegrain)
  * [`fn withTimeGrainUnit(value)`](#fn-azuremonitorwithtimegrainunit)
  * [`fn withTop(value)`](#fn-azuremonitorwithtop)
* [`obj azureResourceGraph`](#obj-azureresourcegraph)
  * [`fn withQuery(value)`](#fn-azureresourcegraphwithquery)
  * [`fn withResultFormat(value)`](#fn-azureresourcegraphwithresultformat)
* [`obj azureTraces`](#obj-azuretraces)
  * [`fn withFilters(value)`](#fn-azuretraceswithfilters)
  * [`fn withFiltersMixin(value)`](#fn-azuretraceswithfiltersmixin)
  * [`fn withOperationId(value)`](#fn-azuretraceswithoperationid)
  * [`fn withQuery(value)`](#fn-azuretraceswithquery)
  * [`fn withResources(value)`](#fn-azuretraceswithresources)
  * [`fn withResourcesMixin(value)`](#fn-azuretraceswithresourcesmixin)
  * [`fn withResultFormat(value)`](#fn-azuretraceswithresultformat)
  * [`fn withTraceTypes(value)`](#fn-azuretraceswithtracetypes)
  * [`fn withTraceTypesMixin(value)`](#fn-azuretraceswithtracetypesmixin)
* [`obj grafanaTemplateVariableFn`](#obj-grafanatemplatevariablefn)
  * [`fn withAppInsightsGroupByQuery(value)`](#fn-grafanatemplatevariablefnwithappinsightsgroupbyquery)
  * [`fn withAppInsightsGroupByQueryMixin(value)`](#fn-grafanatemplatevariablefnwithappinsightsgroupbyquerymixin)
  * [`fn withAppInsightsMetricNameQuery(value)`](#fn-grafanatemplatevariablefnwithappinsightsmetricnamequery)
  * [`fn withAppInsightsMetricNameQueryMixin(value)`](#fn-grafanatemplatevariablefnwithappinsightsmetricnamequerymixin)
  * [`fn withMetricDefinitionsQuery(value)`](#fn-grafanatemplatevariablefnwithmetricdefinitionsquery)
  * [`fn withMetricDefinitionsQueryMixin(value)`](#fn-grafanatemplatevariablefnwithmetricdefinitionsquerymixin)
  * [`fn withMetricNamesQuery(value)`](#fn-grafanatemplatevariablefnwithmetricnamesquery)
  * [`fn withMetricNamesQueryMixin(value)`](#fn-grafanatemplatevariablefnwithmetricnamesquerymixin)
  * [`fn withMetricNamespaceQuery(value)`](#fn-grafanatemplatevariablefnwithmetricnamespacequery)
  * [`fn withMetricNamespaceQueryMixin(value)`](#fn-grafanatemplatevariablefnwithmetricnamespacequerymixin)
  * [`fn withResourceGroupsQuery(value)`](#fn-grafanatemplatevariablefnwithresourcegroupsquery)
  * [`fn withResourceGroupsQueryMixin(value)`](#fn-grafanatemplatevariablefnwithresourcegroupsquerymixin)
  * [`fn withResourceNamesQuery(value)`](#fn-grafanatemplatevariablefnwithresourcenamesquery)
  * [`fn withResourceNamesQueryMixin(value)`](#fn-grafanatemplatevariablefnwithresourcenamesquerymixin)
  * [`fn withSubscriptionsQuery(value)`](#fn-grafanatemplatevariablefnwithsubscriptionsquery)
  * [`fn withSubscriptionsQueryMixin(value)`](#fn-grafanatemplatevariablefnwithsubscriptionsquerymixin)
  * [`fn withUnknownQuery(value)`](#fn-grafanatemplatevariablefnwithunknownquery)
  * [`fn withUnknownQueryMixin(value)`](#fn-grafanatemplatevariablefnwithunknownquerymixin)
  * [`fn withWorkspacesQuery(value)`](#fn-grafanatemplatevariablefnwithworkspacesquery)
  * [`fn withWorkspacesQueryMixin(value)`](#fn-grafanatemplatevariablefnwithworkspacesquerymixin)
  * [`obj AppInsightsGroupByQuery`](#obj-grafanatemplatevariablefnappinsightsgroupbyquery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnappinsightsgroupbyquerywithkind)
    * [`fn withMetricName(value)`](#fn-grafanatemplatevariablefnappinsightsgroupbyquerywithmetricname)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnappinsightsgroupbyquerywithrawquery)
  * [`obj AppInsightsMetricNameQuery`](#obj-grafanatemplatevariablefnappinsightsmetricnamequery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnappinsightsmetricnamequerywithkind)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnappinsightsmetricnamequerywithrawquery)
  * [`obj MetricDefinitionsQuery`](#obj-grafanatemplatevariablefnmetricdefinitionsquery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnmetricdefinitionsquerywithkind)
    * [`fn withMetricNamespace(value)`](#fn-grafanatemplatevariablefnmetricdefinitionsquerywithmetricnamespace)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnmetricdefinitionsquerywithrawquery)
    * [`fn withResourceGroup(value)`](#fn-grafanatemplatevariablefnmetricdefinitionsquerywithresourcegroup)
    * [`fn withResourceName(value)`](#fn-grafanatemplatevariablefnmetricdefinitionsquerywithresourcename)
    * [`fn withSubscription(value)`](#fn-grafanatemplatevariablefnmetricdefinitionsquerywithsubscription)
  * [`obj MetricNamesQuery`](#obj-grafanatemplatevariablefnmetricnamesquery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnmetricnamesquerywithkind)
    * [`fn withMetricNamespace(value)`](#fn-grafanatemplatevariablefnmetricnamesquerywithmetricnamespace)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnmetricnamesquerywithrawquery)
    * [`fn withResourceGroup(value)`](#fn-grafanatemplatevariablefnmetricnamesquerywithresourcegroup)
    * [`fn withResourceName(value)`](#fn-grafanatemplatevariablefnmetricnamesquerywithresourcename)
    * [`fn withSubscription(value)`](#fn-grafanatemplatevariablefnmetricnamesquerywithsubscription)
  * [`obj MetricNamespaceQuery`](#obj-grafanatemplatevariablefnmetricnamespacequery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnmetricnamespacequerywithkind)
    * [`fn withMetricNamespace(value)`](#fn-grafanatemplatevariablefnmetricnamespacequerywithmetricnamespace)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnmetricnamespacequerywithrawquery)
    * [`fn withResourceGroup(value)`](#fn-grafanatemplatevariablefnmetricnamespacequerywithresourcegroup)
    * [`fn withResourceName(value)`](#fn-grafanatemplatevariablefnmetricnamespacequerywithresourcename)
    * [`fn withSubscription(value)`](#fn-grafanatemplatevariablefnmetricnamespacequerywithsubscription)
  * [`obj ResourceGroupsQuery`](#obj-grafanatemplatevariablefnresourcegroupsquery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnresourcegroupsquerywithkind)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnresourcegroupsquerywithrawquery)
    * [`fn withSubscription(value)`](#fn-grafanatemplatevariablefnresourcegroupsquerywithsubscription)
  * [`obj ResourceNamesQuery`](#obj-grafanatemplatevariablefnresourcenamesquery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnresourcenamesquerywithkind)
    * [`fn withMetricNamespace(value)`](#fn-grafanatemplatevariablefnresourcenamesquerywithmetricnamespace)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnresourcenamesquerywithrawquery)
    * [`fn withResourceGroup(value)`](#fn-grafanatemplatevariablefnresourcenamesquerywithresourcegroup)
    * [`fn withSubscription(value)`](#fn-grafanatemplatevariablefnresourcenamesquerywithsubscription)
  * [`obj SubscriptionsQuery`](#obj-grafanatemplatevariablefnsubscriptionsquery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnsubscriptionsquerywithkind)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnsubscriptionsquerywithrawquery)
  * [`obj UnknownQuery`](#obj-grafanatemplatevariablefnunknownquery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnunknownquerywithkind)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnunknownquerywithrawquery)
  * [`obj WorkspacesQuery`](#obj-grafanatemplatevariablefnworkspacesquery)
    * [`fn withKind(value)`](#fn-grafanatemplatevariablefnworkspacesquerywithkind)
    * [`fn withRawQuery(value)`](#fn-grafanatemplatevariablefnworkspacesquerywithrawquery)
    * [`fn withSubscription(value)`](#fn-grafanatemplatevariablefnworkspacesquerywithsubscription)

## Fields

### fn withAzureLogAnalytics

```jsonnet
withAzureLogAnalytics(value)
```

PARAMETERS:

* **value** (`object`)

Azure Monitor Logs sub-query properties
### fn withAzureLogAnalyticsMixin

```jsonnet
withAzureLogAnalyticsMixin(value)
```

PARAMETERS:

* **value** (`object`)

Azure Monitor Logs sub-query properties
### fn withAzureMonitor

```jsonnet
withAzureMonitor(value)
```

PARAMETERS:

* **value** (`object`)


### fn withAzureMonitorMixin

```jsonnet
withAzureMonitorMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withAzureResourceGraph

```jsonnet
withAzureResourceGraph(value)
```

PARAMETERS:

* **value** (`object`)


### fn withAzureResourceGraphMixin

```jsonnet
withAzureResourceGraphMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withAzureTraces

```jsonnet
withAzureTraces(value)
```

PARAMETERS:

* **value** (`object`)

Application Insights Traces sub-query properties
### fn withAzureTracesMixin

```jsonnet
withAzureTracesMixin(value)
```

PARAMETERS:

* **value** (`object`)

Application Insights Traces sub-query properties
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
### fn withGrafanaTemplateVariableFn

```jsonnet
withGrafanaTemplateVariableFn(value)
```

PARAMETERS:

* **value** (`object`)


### fn withGrafanaTemplateVariableFnMixin

```jsonnet
withGrafanaTemplateVariableFnMixin(value)
```

PARAMETERS:

* **value** (`object`)


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
### fn withNamespace

```jsonnet
withNamespace(value)
```

PARAMETERS:

* **value** (`string`)


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
### fn withRegion

```jsonnet
withRegion(value)
```

PARAMETERS:

* **value** (`string`)

Azure Monitor query type.
queryType: #AzureQueryType
### fn withResource

```jsonnet
withResource(value)
```

PARAMETERS:

* **value** (`string`)


### fn withResourceGroup

```jsonnet
withResourceGroup(value)
```

PARAMETERS:

* **value** (`string`)

Template variables params. These exist for backwards compatiblity with legacy template variables.
### fn withSubscription

```jsonnet
withSubscription(value)
```

PARAMETERS:

* **value** (`string`)

Azure subscription containing the resource(s) to be queried.
### fn withSubscriptions

```jsonnet
withSubscriptions(value)
```

PARAMETERS:

* **value** (`array`)

Subscriptions to be queried via Azure Resource Graph.
### fn withSubscriptionsMixin

```jsonnet
withSubscriptionsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Subscriptions to be queried via Azure Resource Graph.
### obj azureLogAnalytics


#### fn azureLogAnalytics.withQuery

```jsonnet
azureLogAnalytics.withQuery(value)
```

PARAMETERS:

* **value** (`string`)

KQL query to be executed.
#### fn azureLogAnalytics.withResource

```jsonnet
azureLogAnalytics.withResource(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated Use resources instead
#### fn azureLogAnalytics.withResources

```jsonnet
azureLogAnalytics.withResources(value)
```

PARAMETERS:

* **value** (`array`)

Array of resource URIs to be queried.
#### fn azureLogAnalytics.withResourcesMixin

```jsonnet
azureLogAnalytics.withResourcesMixin(value)
```

PARAMETERS:

* **value** (`array`)

Array of resource URIs to be queried.
#### fn azureLogAnalytics.withResultFormat

```jsonnet
azureLogAnalytics.withResultFormat(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"table"`, `"time_series"`, `"trace"`


#### fn azureLogAnalytics.withWorkspace

```jsonnet
azureLogAnalytics.withWorkspace(value)
```

PARAMETERS:

* **value** (`string`)

Workspace ID. This was removed in Grafana 8, but remains for backwards compat
### obj azureMonitor


#### fn azureMonitor.withAggregation

```jsonnet
azureMonitor.withAggregation(value)
```

PARAMETERS:

* **value** (`string`)

The aggregation to be used within the query. Defaults to the primaryAggregationType defined by the metric.
#### fn azureMonitor.withAlias

```jsonnet
azureMonitor.withAlias(value)
```

PARAMETERS:

* **value** (`string`)

Aliases can be set to modify the legend labels. e.g. {{ resourceGroup }}. See docs for more detail.
#### fn azureMonitor.withAllowedTimeGrainsMs

```jsonnet
azureMonitor.withAllowedTimeGrainsMs(value)
```

PARAMETERS:

* **value** (`array`)

Time grains that are supported by the metric.
#### fn azureMonitor.withAllowedTimeGrainsMsMixin

```jsonnet
azureMonitor.withAllowedTimeGrainsMsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Time grains that are supported by the metric.
#### fn azureMonitor.withCustomNamespace

```jsonnet
azureMonitor.withCustomNamespace(value)
```

PARAMETERS:

* **value** (`string`)

Used as the value for the metricNamespace property when it's different from the resource namespace.
#### fn azureMonitor.withDimension

```jsonnet
azureMonitor.withDimension(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated This property was migrated to dimensionFilters and should only be accessed in the migration
#### fn azureMonitor.withDimensionFilter

```jsonnet
azureMonitor.withDimensionFilter(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated This property was migrated to dimensionFilters and should only be accessed in the migration
#### fn azureMonitor.withDimensionFilters

```jsonnet
azureMonitor.withDimensionFilters(value)
```

PARAMETERS:

* **value** (`array`)

Filters to reduce the set of data returned. Dimensions that can be filtered on are defined by the metric.
#### fn azureMonitor.withDimensionFiltersMixin

```jsonnet
azureMonitor.withDimensionFiltersMixin(value)
```

PARAMETERS:

* **value** (`array`)

Filters to reduce the set of data returned. Dimensions that can be filtered on are defined by the metric.
#### fn azureMonitor.withMetricDefinition

```jsonnet
azureMonitor.withMetricDefinition(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated Use metricNamespace instead
#### fn azureMonitor.withMetricName

```jsonnet
azureMonitor.withMetricName(value)
```

PARAMETERS:

* **value** (`string`)

The metric to query data for within the specified metricNamespace. e.g. UsedCapacity
#### fn azureMonitor.withMetricNamespace

```jsonnet
azureMonitor.withMetricNamespace(value)
```

PARAMETERS:

* **value** (`string`)

metricNamespace is used as the resource type (or resource namespace).
It's usually equal to the target metric namespace. e.g. microsoft.storage/storageaccounts
Kept the name of the variable as metricNamespace to avoid backward incompatibility issues.
#### fn azureMonitor.withRegion

```jsonnet
azureMonitor.withRegion(value)
```

PARAMETERS:

* **value** (`string`)

The Azure region containing the resource(s).
#### fn azureMonitor.withResourceGroup

```jsonnet
azureMonitor.withResourceGroup(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated Use resources instead
#### fn azureMonitor.withResourceName

```jsonnet
azureMonitor.withResourceName(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated Use resources instead
#### fn azureMonitor.withResourceUri

```jsonnet
azureMonitor.withResourceUri(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated Use resourceGroup, resourceName and metricNamespace instead
#### fn azureMonitor.withResources

```jsonnet
azureMonitor.withResources(value)
```

PARAMETERS:

* **value** (`array`)

Array of resource URIs to be queried.
#### fn azureMonitor.withResourcesMixin

```jsonnet
azureMonitor.withResourcesMixin(value)
```

PARAMETERS:

* **value** (`array`)

Array of resource URIs to be queried.
#### fn azureMonitor.withTimeGrain

```jsonnet
azureMonitor.withTimeGrain(value)
```

PARAMETERS:

* **value** (`string`)

The granularity of data points to be queried. Defaults to auto.
#### fn azureMonitor.withTimeGrainUnit

```jsonnet
azureMonitor.withTimeGrainUnit(value)
```

PARAMETERS:

* **value** (`string`)

@deprecated
#### fn azureMonitor.withTop

```jsonnet
azureMonitor.withTop(value)
```

PARAMETERS:

* **value** (`string`)

Maximum number of records to return. Defaults to 10.
### obj azureResourceGraph


#### fn azureResourceGraph.withQuery

```jsonnet
azureResourceGraph.withQuery(value)
```

PARAMETERS:

* **value** (`string`)

Azure Resource Graph KQL query to be executed.
#### fn azureResourceGraph.withResultFormat

```jsonnet
azureResourceGraph.withResultFormat(value)
```

PARAMETERS:

* **value** (`string`)

Specifies the format results should be returned as. Defaults to table.
### obj azureTraces


#### fn azureTraces.withFilters

```jsonnet
azureTraces.withFilters(value)
```

PARAMETERS:

* **value** (`array`)

Filters for property values.
#### fn azureTraces.withFiltersMixin

```jsonnet
azureTraces.withFiltersMixin(value)
```

PARAMETERS:

* **value** (`array`)

Filters for property values.
#### fn azureTraces.withOperationId

```jsonnet
azureTraces.withOperationId(value)
```

PARAMETERS:

* **value** (`string`)

Operation ID. Used only for Traces queries.
#### fn azureTraces.withQuery

```jsonnet
azureTraces.withQuery(value)
```

PARAMETERS:

* **value** (`string`)

KQL query to be executed.
#### fn azureTraces.withResources

```jsonnet
azureTraces.withResources(value)
```

PARAMETERS:

* **value** (`array`)

Array of resource URIs to be queried.
#### fn azureTraces.withResourcesMixin

```jsonnet
azureTraces.withResourcesMixin(value)
```

PARAMETERS:

* **value** (`array`)

Array of resource URIs to be queried.
#### fn azureTraces.withResultFormat

```jsonnet
azureTraces.withResultFormat(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"table"`, `"time_series"`, `"trace"`


#### fn azureTraces.withTraceTypes

```jsonnet
azureTraces.withTraceTypes(value)
```

PARAMETERS:

* **value** (`array`)

Types of events to filter by.
#### fn azureTraces.withTraceTypesMixin

```jsonnet
azureTraces.withTraceTypesMixin(value)
```

PARAMETERS:

* **value** (`array`)

Types of events to filter by.
### obj grafanaTemplateVariableFn


#### fn grafanaTemplateVariableFn.withAppInsightsGroupByQuery

```jsonnet
grafanaTemplateVariableFn.withAppInsightsGroupByQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withAppInsightsGroupByQueryMixin

```jsonnet
grafanaTemplateVariableFn.withAppInsightsGroupByQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withAppInsightsMetricNameQuery

```jsonnet
grafanaTemplateVariableFn.withAppInsightsMetricNameQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withAppInsightsMetricNameQueryMixin

```jsonnet
grafanaTemplateVariableFn.withAppInsightsMetricNameQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withMetricDefinitionsQuery

```jsonnet
grafanaTemplateVariableFn.withMetricDefinitionsQuery(value)
```

PARAMETERS:

* **value** (`object`)

@deprecated Use MetricNamespaceQuery instead
#### fn grafanaTemplateVariableFn.withMetricDefinitionsQueryMixin

```jsonnet
grafanaTemplateVariableFn.withMetricDefinitionsQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)

@deprecated Use MetricNamespaceQuery instead
#### fn grafanaTemplateVariableFn.withMetricNamesQuery

```jsonnet
grafanaTemplateVariableFn.withMetricNamesQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withMetricNamesQueryMixin

```jsonnet
grafanaTemplateVariableFn.withMetricNamesQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withMetricNamespaceQuery

```jsonnet
grafanaTemplateVariableFn.withMetricNamespaceQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withMetricNamespaceQueryMixin

```jsonnet
grafanaTemplateVariableFn.withMetricNamespaceQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withResourceGroupsQuery

```jsonnet
grafanaTemplateVariableFn.withResourceGroupsQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withResourceGroupsQueryMixin

```jsonnet
grafanaTemplateVariableFn.withResourceGroupsQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withResourceNamesQuery

```jsonnet
grafanaTemplateVariableFn.withResourceNamesQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withResourceNamesQueryMixin

```jsonnet
grafanaTemplateVariableFn.withResourceNamesQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withSubscriptionsQuery

```jsonnet
grafanaTemplateVariableFn.withSubscriptionsQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withSubscriptionsQueryMixin

```jsonnet
grafanaTemplateVariableFn.withSubscriptionsQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withUnknownQuery

```jsonnet
grafanaTemplateVariableFn.withUnknownQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withUnknownQueryMixin

```jsonnet
grafanaTemplateVariableFn.withUnknownQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withWorkspacesQuery

```jsonnet
grafanaTemplateVariableFn.withWorkspacesQuery(value)
```

PARAMETERS:

* **value** (`object`)


#### fn grafanaTemplateVariableFn.withWorkspacesQueryMixin

```jsonnet
grafanaTemplateVariableFn.withWorkspacesQueryMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### obj grafanaTemplateVariableFn.AppInsightsGroupByQuery


##### fn grafanaTemplateVariableFn.AppInsightsGroupByQuery.withKind

```jsonnet
grafanaTemplateVariableFn.AppInsightsGroupByQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"AppInsightsGroupByQuery"`


##### fn grafanaTemplateVariableFn.AppInsightsGroupByQuery.withMetricName

```jsonnet
grafanaTemplateVariableFn.AppInsightsGroupByQuery.withMetricName(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.AppInsightsGroupByQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.AppInsightsGroupByQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.AppInsightsMetricNameQuery


##### fn grafanaTemplateVariableFn.AppInsightsMetricNameQuery.withKind

```jsonnet
grafanaTemplateVariableFn.AppInsightsMetricNameQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"AppInsightsMetricNameQuery"`


##### fn grafanaTemplateVariableFn.AppInsightsMetricNameQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.AppInsightsMetricNameQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.MetricDefinitionsQuery


##### fn grafanaTemplateVariableFn.MetricDefinitionsQuery.withKind

```jsonnet
grafanaTemplateVariableFn.MetricDefinitionsQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"MetricDefinitionsQuery"`


##### fn grafanaTemplateVariableFn.MetricDefinitionsQuery.withMetricNamespace

```jsonnet
grafanaTemplateVariableFn.MetricDefinitionsQuery.withMetricNamespace(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricDefinitionsQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.MetricDefinitionsQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricDefinitionsQuery.withResourceGroup

```jsonnet
grafanaTemplateVariableFn.MetricDefinitionsQuery.withResourceGroup(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricDefinitionsQuery.withResourceName

```jsonnet
grafanaTemplateVariableFn.MetricDefinitionsQuery.withResourceName(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricDefinitionsQuery.withSubscription

```jsonnet
grafanaTemplateVariableFn.MetricDefinitionsQuery.withSubscription(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.MetricNamesQuery


##### fn grafanaTemplateVariableFn.MetricNamesQuery.withKind

```jsonnet
grafanaTemplateVariableFn.MetricNamesQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"MetricNamesQuery"`


##### fn grafanaTemplateVariableFn.MetricNamesQuery.withMetricNamespace

```jsonnet
grafanaTemplateVariableFn.MetricNamesQuery.withMetricNamespace(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricNamesQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.MetricNamesQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricNamesQuery.withResourceGroup

```jsonnet
grafanaTemplateVariableFn.MetricNamesQuery.withResourceGroup(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricNamesQuery.withResourceName

```jsonnet
grafanaTemplateVariableFn.MetricNamesQuery.withResourceName(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricNamesQuery.withSubscription

```jsonnet
grafanaTemplateVariableFn.MetricNamesQuery.withSubscription(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.MetricNamespaceQuery


##### fn grafanaTemplateVariableFn.MetricNamespaceQuery.withKind

```jsonnet
grafanaTemplateVariableFn.MetricNamespaceQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"MetricNamespaceQuery"`


##### fn grafanaTemplateVariableFn.MetricNamespaceQuery.withMetricNamespace

```jsonnet
grafanaTemplateVariableFn.MetricNamespaceQuery.withMetricNamespace(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricNamespaceQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.MetricNamespaceQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricNamespaceQuery.withResourceGroup

```jsonnet
grafanaTemplateVariableFn.MetricNamespaceQuery.withResourceGroup(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricNamespaceQuery.withResourceName

```jsonnet
grafanaTemplateVariableFn.MetricNamespaceQuery.withResourceName(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.MetricNamespaceQuery.withSubscription

```jsonnet
grafanaTemplateVariableFn.MetricNamespaceQuery.withSubscription(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.ResourceGroupsQuery


##### fn grafanaTemplateVariableFn.ResourceGroupsQuery.withKind

```jsonnet
grafanaTemplateVariableFn.ResourceGroupsQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"ResourceGroupsQuery"`


##### fn grafanaTemplateVariableFn.ResourceGroupsQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.ResourceGroupsQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.ResourceGroupsQuery.withSubscription

```jsonnet
grafanaTemplateVariableFn.ResourceGroupsQuery.withSubscription(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.ResourceNamesQuery


##### fn grafanaTemplateVariableFn.ResourceNamesQuery.withKind

```jsonnet
grafanaTemplateVariableFn.ResourceNamesQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"ResourceNamesQuery"`


##### fn grafanaTemplateVariableFn.ResourceNamesQuery.withMetricNamespace

```jsonnet
grafanaTemplateVariableFn.ResourceNamesQuery.withMetricNamespace(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.ResourceNamesQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.ResourceNamesQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.ResourceNamesQuery.withResourceGroup

```jsonnet
grafanaTemplateVariableFn.ResourceNamesQuery.withResourceGroup(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.ResourceNamesQuery.withSubscription

```jsonnet
grafanaTemplateVariableFn.ResourceNamesQuery.withSubscription(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.SubscriptionsQuery


##### fn grafanaTemplateVariableFn.SubscriptionsQuery.withKind

```jsonnet
grafanaTemplateVariableFn.SubscriptionsQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"SubscriptionsQuery"`


##### fn grafanaTemplateVariableFn.SubscriptionsQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.SubscriptionsQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.UnknownQuery


##### fn grafanaTemplateVariableFn.UnknownQuery.withKind

```jsonnet
grafanaTemplateVariableFn.UnknownQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"UnknownQuery"`


##### fn grafanaTemplateVariableFn.UnknownQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.UnknownQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


#### obj grafanaTemplateVariableFn.WorkspacesQuery


##### fn grafanaTemplateVariableFn.WorkspacesQuery.withKind

```jsonnet
grafanaTemplateVariableFn.WorkspacesQuery.withKind(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"WorkspacesQuery"`


##### fn grafanaTemplateVariableFn.WorkspacesQuery.withRawQuery

```jsonnet
grafanaTemplateVariableFn.WorkspacesQuery.withRawQuery(value)
```

PARAMETERS:

* **value** (`string`)


##### fn grafanaTemplateVariableFn.WorkspacesQuery.withSubscription

```jsonnet
grafanaTemplateVariableFn.WorkspacesQuery.withSubscription(value)
```

PARAMETERS:

* **value** (`string`)

