# testData

grafonnet.query.testData

## Subpackages

* [csvWave](csvWave.md)

## Index

* [`fn withAlias(value)`](#fn-withalias)
* [`fn withChannel(value)`](#fn-withchannel)
* [`fn withCsvContent(value)`](#fn-withcsvcontent)
* [`fn withCsvFileName(value)`](#fn-withcsvfilename)
* [`fn withCsvWave(value)`](#fn-withcsvwave)
* [`fn withCsvWaveMixin(value)`](#fn-withcsvwavemixin)
* [`fn withDatasource(value)`](#fn-withdatasource)
* [`fn withErrorType(value)`](#fn-witherrortype)
* [`fn withHide(value=true)`](#fn-withhide)
* [`fn withLabels(value)`](#fn-withlabels)
* [`fn withLevelColumn(value=true)`](#fn-withlevelcolumn)
* [`fn withLines(value)`](#fn-withlines)
* [`fn withNodes(value)`](#fn-withnodes)
* [`fn withNodesMixin(value)`](#fn-withnodesmixin)
* [`fn withPoints(value)`](#fn-withpoints)
* [`fn withPointsMixin(value)`](#fn-withpointsmixin)
* [`fn withPulseWave(value)`](#fn-withpulsewave)
* [`fn withPulseWaveMixin(value)`](#fn-withpulsewavemixin)
* [`fn withQueryType(value)`](#fn-withquerytype)
* [`fn withRawFrameContent(value)`](#fn-withrawframecontent)
* [`fn withRefId(value)`](#fn-withrefid)
* [`fn withScenarioId(value)`](#fn-withscenarioid)
* [`fn withSeriesCount(value)`](#fn-withseriescount)
* [`fn withSim(value)`](#fn-withsim)
* [`fn withSimMixin(value)`](#fn-withsimmixin)
* [`fn withSpanCount(value)`](#fn-withspancount)
* [`fn withStream(value)`](#fn-withstream)
* [`fn withStreamMixin(value)`](#fn-withstreammixin)
* [`fn withStringInput(value)`](#fn-withstringinput)
* [`fn withUsa(value)`](#fn-withusa)
* [`fn withUsaMixin(value)`](#fn-withusamixin)
* [`obj nodes`](#obj-nodes)
  * [`fn withCount(value)`](#fn-nodeswithcount)
  * [`fn withType(value)`](#fn-nodeswithtype)
* [`obj pulseWave`](#obj-pulsewave)
  * [`fn withOffCount(value)`](#fn-pulsewavewithoffcount)
  * [`fn withOffValue(value)`](#fn-pulsewavewithoffvalue)
  * [`fn withOnCount(value)`](#fn-pulsewavewithoncount)
  * [`fn withOnValue(value)`](#fn-pulsewavewithonvalue)
  * [`fn withTimeStep(value)`](#fn-pulsewavewithtimestep)
* [`obj sim`](#obj-sim)
  * [`fn withConfig(value)`](#fn-simwithconfig)
  * [`fn withConfigMixin(value)`](#fn-simwithconfigmixin)
  * [`fn withKey(value)`](#fn-simwithkey)
  * [`fn withKeyMixin(value)`](#fn-simwithkeymixin)
  * [`fn withLast(value=true)`](#fn-simwithlast)
  * [`fn withStream(value=true)`](#fn-simwithstream)
  * [`obj key`](#obj-simkey)
    * [`fn withTick(value)`](#fn-simkeywithtick)
    * [`fn withType(value)`](#fn-simkeywithtype)
    * [`fn withUid(value)`](#fn-simkeywithuid)
* [`obj stream`](#obj-stream)
  * [`fn withBands(value)`](#fn-streamwithbands)
  * [`fn withNoise(value)`](#fn-streamwithnoise)
  * [`fn withSpeed(value)`](#fn-streamwithspeed)
  * [`fn withSpread(value)`](#fn-streamwithspread)
  * [`fn withType(value)`](#fn-streamwithtype)
  * [`fn withUrl(value)`](#fn-streamwithurl)
* [`obj usa`](#obj-usa)
  * [`fn withFields(value)`](#fn-usawithfields)
  * [`fn withFieldsMixin(value)`](#fn-usawithfieldsmixin)
  * [`fn withMode(value)`](#fn-usawithmode)
  * [`fn withPeriod(value)`](#fn-usawithperiod)
  * [`fn withStates(value)`](#fn-usawithstates)
  * [`fn withStatesMixin(value)`](#fn-usawithstatesmixin)

## Fields

### fn withAlias

```jsonnet
withAlias(value)
```

PARAMETERS:

* **value** (`string`)


### fn withChannel

```jsonnet
withChannel(value)
```

PARAMETERS:

* **value** (`string`)


### fn withCsvContent

```jsonnet
withCsvContent(value)
```

PARAMETERS:

* **value** (`string`)


### fn withCsvFileName

```jsonnet
withCsvFileName(value)
```

PARAMETERS:

* **value** (`string`)


### fn withCsvWave

```jsonnet
withCsvWave(value)
```

PARAMETERS:

* **value** (`array`)


### fn withCsvWaveMixin

```jsonnet
withCsvWaveMixin(value)
```

PARAMETERS:

* **value** (`array`)


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
### fn withErrorType

```jsonnet
withErrorType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"server_panic"`, `"frontend_exception"`, `"frontend_observable"`


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
### fn withLabels

```jsonnet
withLabels(value)
```

PARAMETERS:

* **value** (`string`)


### fn withLevelColumn

```jsonnet
withLevelColumn(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


### fn withLines

```jsonnet
withLines(value)
```

PARAMETERS:

* **value** (`integer`)


### fn withNodes

```jsonnet
withNodes(value)
```

PARAMETERS:

* **value** (`object`)


### fn withNodesMixin

```jsonnet
withNodesMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withPoints

```jsonnet
withPoints(value)
```

PARAMETERS:

* **value** (`array`)


### fn withPointsMixin

```jsonnet
withPointsMixin(value)
```

PARAMETERS:

* **value** (`array`)


### fn withPulseWave

```jsonnet
withPulseWave(value)
```

PARAMETERS:

* **value** (`object`)


### fn withPulseWaveMixin

```jsonnet
withPulseWaveMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withQueryType

```jsonnet
withQueryType(value)
```

PARAMETERS:

* **value** (`string`)

Specify the query flavor
TODO make this required and give it a default
### fn withRawFrameContent

```jsonnet
withRawFrameContent(value)
```

PARAMETERS:

* **value** (`string`)


### fn withRefId

```jsonnet
withRefId(value)
```

PARAMETERS:

* **value** (`string`)

A unique identifier for the query within the list of targets.
In server side expressions, the refId is used as a variable name to identify results.
By default, the UI will assign A->Z; however setting meaningful names may be useful.
### fn withScenarioId

```jsonnet
withScenarioId(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"random_walk"`, `"slow_query"`, `"random_walk_with_error"`, `"random_walk_table"`, `"exponential_heatmap_bucket_data"`, `"linear_heatmap_bucket_data"`, `"no_data_points"`, `"datapoints_outside_range"`, `"csv_metric_values"`, `"predictable_pulse"`, `"predictable_csv_wave"`, `"streaming_client"`, `"simulation"`, `"usa"`, `"live"`, `"grafana_api"`, `"arrow"`, `"annotations"`, `"table_static"`, `"server_error_500"`, `"logs"`, `"node_graph"`, `"flame_graph"`, `"raw_frame"`, `"csv_file"`, `"csv_content"`, `"trace"`, `"manual_entry"`, `"variables-query"`


### fn withSeriesCount

```jsonnet
withSeriesCount(value)
```

PARAMETERS:

* **value** (`integer`)


### fn withSim

```jsonnet
withSim(value)
```

PARAMETERS:

* **value** (`object`)


### fn withSimMixin

```jsonnet
withSimMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withSpanCount

```jsonnet
withSpanCount(value)
```

PARAMETERS:

* **value** (`integer`)


### fn withStream

```jsonnet
withStream(value)
```

PARAMETERS:

* **value** (`object`)


### fn withStreamMixin

```jsonnet
withStreamMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withStringInput

```jsonnet
withStringInput(value)
```

PARAMETERS:

* **value** (`string`)


### fn withUsa

```jsonnet
withUsa(value)
```

PARAMETERS:

* **value** (`object`)


### fn withUsaMixin

```jsonnet
withUsaMixin(value)
```

PARAMETERS:

* **value** (`object`)


### obj nodes


#### fn nodes.withCount

```jsonnet
nodes.withCount(value)
```

PARAMETERS:

* **value** (`integer`)


#### fn nodes.withType

```jsonnet
nodes.withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"random"`, `"response"`, `"random edges"`


### obj pulseWave


#### fn pulseWave.withOffCount

```jsonnet
pulseWave.withOffCount(value)
```

PARAMETERS:

* **value** (`integer`)


#### fn pulseWave.withOffValue

```jsonnet
pulseWave.withOffValue(value)
```

PARAMETERS:

* **value** (`number`)


#### fn pulseWave.withOnCount

```jsonnet
pulseWave.withOnCount(value)
```

PARAMETERS:

* **value** (`integer`)


#### fn pulseWave.withOnValue

```jsonnet
pulseWave.withOnValue(value)
```

PARAMETERS:

* **value** (`number`)


#### fn pulseWave.withTimeStep

```jsonnet
pulseWave.withTimeStep(value)
```

PARAMETERS:

* **value** (`integer`)


### obj sim


#### fn sim.withConfig

```jsonnet
sim.withConfig(value)
```

PARAMETERS:

* **value** (`object`)


#### fn sim.withConfigMixin

```jsonnet
sim.withConfigMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn sim.withKey

```jsonnet
sim.withKey(value)
```

PARAMETERS:

* **value** (`object`)


#### fn sim.withKeyMixin

```jsonnet
sim.withKeyMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn sim.withLast

```jsonnet
sim.withLast(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


#### fn sim.withStream

```jsonnet
sim.withStream(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`


#### obj sim.key


##### fn sim.key.withTick

```jsonnet
sim.key.withTick(value)
```

PARAMETERS:

* **value** (`number`)


##### fn sim.key.withType

```jsonnet
sim.key.withType(value)
```

PARAMETERS:

* **value** (`string`)


##### fn sim.key.withUid

```jsonnet
sim.key.withUid(value)
```

PARAMETERS:

* **value** (`string`)


### obj stream


#### fn stream.withBands

```jsonnet
stream.withBands(value)
```

PARAMETERS:

* **value** (`integer`)


#### fn stream.withNoise

```jsonnet
stream.withNoise(value)
```

PARAMETERS:

* **value** (`integer`)


#### fn stream.withSpeed

```jsonnet
stream.withSpeed(value)
```

PARAMETERS:

* **value** (`integer`)


#### fn stream.withSpread

```jsonnet
stream.withSpread(value)
```

PARAMETERS:

* **value** (`integer`)


#### fn stream.withType

```jsonnet
stream.withType(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"signal"`, `"logs"`, `"fetch"`


#### fn stream.withUrl

```jsonnet
stream.withUrl(value)
```

PARAMETERS:

* **value** (`string`)


### obj usa


#### fn usa.withFields

```jsonnet
usa.withFields(value)
```

PARAMETERS:

* **value** (`array`)


#### fn usa.withFieldsMixin

```jsonnet
usa.withFieldsMixin(value)
```

PARAMETERS:

* **value** (`array`)


#### fn usa.withMode

```jsonnet
usa.withMode(value)
```

PARAMETERS:

* **value** (`string`)


#### fn usa.withPeriod

```jsonnet
usa.withPeriod(value)
```

PARAMETERS:

* **value** (`string`)


#### fn usa.withStates

```jsonnet
usa.withStates(value)
```

PARAMETERS:

* **value** (`array`)


#### fn usa.withStatesMixin

```jsonnet
usa.withStatesMixin(value)
```

PARAMETERS:

* **value** (`array`)

