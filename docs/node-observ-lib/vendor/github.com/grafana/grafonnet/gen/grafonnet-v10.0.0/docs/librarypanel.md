# librarypanel

grafonnet.librarypanel

## Index

* [`fn withDescription(value)`](#fn-withdescription)
* [`fn withFolderUid(value)`](#fn-withfolderuid)
* [`fn withMeta(value)`](#fn-withmeta)
* [`fn withMetaMixin(value)`](#fn-withmetamixin)
* [`fn withModel(value)`](#fn-withmodel)
* [`fn withModelMixin(value)`](#fn-withmodelmixin)
* [`fn withName(value)`](#fn-withname)
* [`fn withSchemaVersion(value)`](#fn-withschemaversion)
* [`fn withType(value)`](#fn-withtype)
* [`fn withUid(value)`](#fn-withuid)
* [`fn withVersion(value)`](#fn-withversion)
* [`obj meta`](#obj-meta)
  * [`fn withConnectedDashboards(value)`](#fn-metawithconnecteddashboards)
  * [`fn withCreated(value)`](#fn-metawithcreated)
  * [`fn withCreatedBy(value)`](#fn-metawithcreatedby)
  * [`fn withCreatedByMixin(value)`](#fn-metawithcreatedbymixin)
  * [`fn withFolderName(value)`](#fn-metawithfoldername)
  * [`fn withFolderUid(value)`](#fn-metawithfolderuid)
  * [`fn withUpdated(value)`](#fn-metawithupdated)
  * [`fn withUpdatedBy(value)`](#fn-metawithupdatedby)
  * [`fn withUpdatedByMixin(value)`](#fn-metawithupdatedbymixin)
  * [`obj createdBy`](#obj-metacreatedby)
    * [`fn withAvatarUrl(value)`](#fn-metacreatedbywithavatarurl)
    * [`fn withId(value)`](#fn-metacreatedbywithid)
    * [`fn withName(value)`](#fn-metacreatedbywithname)
  * [`obj updatedBy`](#obj-metaupdatedby)
    * [`fn withAvatarUrl(value)`](#fn-metaupdatedbywithavatarurl)
    * [`fn withId(value)`](#fn-metaupdatedbywithid)
    * [`fn withName(value)`](#fn-metaupdatedbywithname)

## Fields

### fn withDescription

```jsonnet
withDescription(value)
```

PARAMETERS:

* **value** (`string`)

Panel description
### fn withFolderUid

```jsonnet
withFolderUid(value)
```

PARAMETERS:

* **value** (`string`)

Folder UID
### fn withMeta

```jsonnet
withMeta(value)
```

PARAMETERS:

* **value** (`object`)


### fn withMetaMixin

```jsonnet
withMetaMixin(value)
```

PARAMETERS:

* **value** (`object`)


### fn withModel

```jsonnet
withModel(value)
```

PARAMETERS:

* **value** (`object`)

TODO: should be the same panel schema defined in dashboard
Typescript: Omit<Panel, 'gridPos' | 'id' | 'libraryPanel'>;
### fn withModelMixin

```jsonnet
withModelMixin(value)
```

PARAMETERS:

* **value** (`object`)

TODO: should be the same panel schema defined in dashboard
Typescript: Omit<Panel, 'gridPos' | 'id' | 'libraryPanel'>;
### fn withName

```jsonnet
withName(value)
```

PARAMETERS:

* **value** (`string`)

Panel name (also saved in the model)
### fn withSchemaVersion

```jsonnet
withSchemaVersion(value)
```

PARAMETERS:

* **value** (`integer`)

Dashboard version when this was saved (zero if unknown)
### fn withType

```jsonnet
withType(value)
```

PARAMETERS:

* **value** (`string`)

The panel type (from inside the model)
### fn withUid

```jsonnet
withUid(value)
```

PARAMETERS:

* **value** (`string`)

Library element UID
### fn withVersion

```jsonnet
withVersion(value)
```

PARAMETERS:

* **value** (`integer`)

panel version, incremented each time the dashboard is updated.
### obj meta


#### fn meta.withConnectedDashboards

```jsonnet
meta.withConnectedDashboards(value)
```

PARAMETERS:

* **value** (`integer`)


#### fn meta.withCreated

```jsonnet
meta.withCreated(value)
```

PARAMETERS:

* **value** (`string`)


#### fn meta.withCreatedBy

```jsonnet
meta.withCreatedBy(value)
```

PARAMETERS:

* **value** (`object`)


#### fn meta.withCreatedByMixin

```jsonnet
meta.withCreatedByMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### fn meta.withFolderName

```jsonnet
meta.withFolderName(value)
```

PARAMETERS:

* **value** (`string`)


#### fn meta.withFolderUid

```jsonnet
meta.withFolderUid(value)
```

PARAMETERS:

* **value** (`string`)


#### fn meta.withUpdated

```jsonnet
meta.withUpdated(value)
```

PARAMETERS:

* **value** (`string`)


#### fn meta.withUpdatedBy

```jsonnet
meta.withUpdatedBy(value)
```

PARAMETERS:

* **value** (`object`)


#### fn meta.withUpdatedByMixin

```jsonnet
meta.withUpdatedByMixin(value)
```

PARAMETERS:

* **value** (`object`)


#### obj meta.createdBy


##### fn meta.createdBy.withAvatarUrl

```jsonnet
meta.createdBy.withAvatarUrl(value)
```

PARAMETERS:

* **value** (`string`)


##### fn meta.createdBy.withId

```jsonnet
meta.createdBy.withId(value)
```

PARAMETERS:

* **value** (`integer`)


##### fn meta.createdBy.withName

```jsonnet
meta.createdBy.withName(value)
```

PARAMETERS:

* **value** (`string`)


#### obj meta.updatedBy


##### fn meta.updatedBy.withAvatarUrl

```jsonnet
meta.updatedBy.withAvatarUrl(value)
```

PARAMETERS:

* **value** (`string`)


##### fn meta.updatedBy.withId

```jsonnet
meta.updatedBy.withId(value)
```

PARAMETERS:

* **value** (`integer`)


##### fn meta.updatedBy.withName

```jsonnet
meta.updatedBy.withName(value)
```

PARAMETERS:

* **value** (`string`)

