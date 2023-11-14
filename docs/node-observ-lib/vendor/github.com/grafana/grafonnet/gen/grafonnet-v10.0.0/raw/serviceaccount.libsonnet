// This file is generated, do not manually edit.
{
  '#': { help: 'grafonnet.serviceaccount', name: 'serviceaccount' },
  '#withAccessControl': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'AccessControl metadata associated with a given resource.' } },
  withAccessControl(value): { accessControl: value },
  '#withAccessControlMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'object' }], help: 'AccessControl metadata associated with a given resource.' } },
  withAccessControlMixin(value): { accessControl+: value },
  '#withAvatarUrl': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: "AvatarUrl is the service account's avatar URL. It allows the frontend to display a picture in front\nof the service account." } },
  withAvatarUrl(value): { avatarUrl: value },
  '#withId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'ID is the unique identifier of the service account in the database.' } },
  withId(value): { id: value },
  '#withIsDisabled': { 'function': { args: [{ default: true, enums: null, name: 'value', type: 'boolean' }], help: 'IsDisabled indicates if the service account is disabled.' } },
  withIsDisabled(value=true): { isDisabled: value },
  '#withLogin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Login of the service account.' } },
  withLogin(value): { login: value },
  '#withName': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'string' }], help: 'Name of the service account.' } },
  withName(value): { name: value },
  '#withOrgId': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'OrgId is the ID of an organisation the service account belongs to.' } },
  withOrgId(value): { orgId: value },
  '#withRole': { 'function': { args: [{ default: null, enums: ['Admin', 'Editor', 'Viewer'], name: 'value', type: 'string' }], help: "OrgRole is a Grafana Organization Role which can be 'Viewer', 'Editor', 'Admin'." } },
  withRole(value): { role: value },
  '#withTeams': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Teams is a list of teams the service account belongs to.' } },
  withTeams(value): { teams: (if std.isArray(value)
                              then value
                              else [value]) },
  '#withTeamsMixin': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'array' }], help: 'Teams is a list of teams the service account belongs to.' } },
  withTeamsMixin(value): { teams+: (if std.isArray(value)
                                    then value
                                    else [value]) },
  '#withTokens': { 'function': { args: [{ default: null, enums: null, name: 'value', type: 'integer' }], help: 'Tokens is the number of active tokens for the service account.\nTokens are used to authenticate the service account against Grafana.' } },
  withTokens(value): { tokens: value },
}
