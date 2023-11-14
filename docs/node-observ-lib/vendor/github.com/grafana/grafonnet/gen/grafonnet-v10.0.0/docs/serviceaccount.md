# serviceaccount

grafonnet.serviceaccount

## Index

* [`fn withAccessControl(value)`](#fn-withaccesscontrol)
* [`fn withAccessControlMixin(value)`](#fn-withaccesscontrolmixin)
* [`fn withAvatarUrl(value)`](#fn-withavatarurl)
* [`fn withId(value)`](#fn-withid)
* [`fn withIsDisabled(value=true)`](#fn-withisdisabled)
* [`fn withLogin(value)`](#fn-withlogin)
* [`fn withName(value)`](#fn-withname)
* [`fn withOrgId(value)`](#fn-withorgid)
* [`fn withRole(value)`](#fn-withrole)
* [`fn withTeams(value)`](#fn-withteams)
* [`fn withTeamsMixin(value)`](#fn-withteamsmixin)
* [`fn withTokens(value)`](#fn-withtokens)

## Fields

### fn withAccessControl

```jsonnet
withAccessControl(value)
```

PARAMETERS:

* **value** (`object`)

AccessControl metadata associated with a given resource.
### fn withAccessControlMixin

```jsonnet
withAccessControlMixin(value)
```

PARAMETERS:

* **value** (`object`)

AccessControl metadata associated with a given resource.
### fn withAvatarUrl

```jsonnet
withAvatarUrl(value)
```

PARAMETERS:

* **value** (`string`)

AvatarUrl is the service account's avatar URL. It allows the frontend to display a picture in front
of the service account.
### fn withId

```jsonnet
withId(value)
```

PARAMETERS:

* **value** (`integer`)

ID is the unique identifier of the service account in the database.
### fn withIsDisabled

```jsonnet
withIsDisabled(value=true)
```

PARAMETERS:

* **value** (`boolean`)
   - default value: `true`

IsDisabled indicates if the service account is disabled.
### fn withLogin

```jsonnet
withLogin(value)
```

PARAMETERS:

* **value** (`string`)

Login of the service account.
### fn withName

```jsonnet
withName(value)
```

PARAMETERS:

* **value** (`string`)

Name of the service account.
### fn withOrgId

```jsonnet
withOrgId(value)
```

PARAMETERS:

* **value** (`integer`)

OrgId is the ID of an organisation the service account belongs to.
### fn withRole

```jsonnet
withRole(value)
```

PARAMETERS:

* **value** (`string`)
   - valid values: `"Admin"`, `"Editor"`, `"Viewer"`

OrgRole is a Grafana Organization Role which can be 'Viewer', 'Editor', 'Admin'.
### fn withTeams

```jsonnet
withTeams(value)
```

PARAMETERS:

* **value** (`array`)

Teams is a list of teams the service account belongs to.
### fn withTeamsMixin

```jsonnet
withTeamsMixin(value)
```

PARAMETERS:

* **value** (`array`)

Teams is a list of teams the service account belongs to.
### fn withTokens

```jsonnet
withTokens(value)
```

PARAMETERS:

* **value** (`integer`)

Tokens is the number of active tokens for the service account.
Tokens are used to authenticate the service account against Grafana.