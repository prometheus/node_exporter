{
  local d = self,

  '#':
    d.pkg(
      name='doc-util',
      url='github.com/jsonnet-libs/docsonnet/doc-util',
      help=|||
        `doc-util` provides a Jsonnet interface for `docsonnet`,
         a Jsonnet API doc generator that uses structured data instead of comments.
      |||,
      filename=std.thisFile,
    )
    + d.package.withUsageTemplate(
      'local d = import "%(import)s"'
    ),

  package:: {
    '#new':: d.fn(|||
      `new` creates a new package

      Arguments:

      * given `name`
      * source `url` for jsonnet-bundler and the import
      * `help` text
      * `filename` for the import, defaults to blank for backward compatibility
      * `version` for jsonnet-bundler install, defaults to `master` just like jsonnet-bundler
    |||, [
      d.arg('name', d.T.string),
      d.arg('url', d.T.string),
      d.arg('help', d.T.string),
      d.arg('filename', d.T.string, ''),
      d.arg('version', d.T.string, 'master'),
    ]),
    new(name, url, help, filename='', version='master')::
      {
        name: name,
        help:
          help
          + std.get(self, 'installTemplate', '') % self
          + std.get(self, 'usageTemplate', '') % self,
        'import':
          if filename != ''
          then url + '/' + filename
          else url,
        url: url,
        filename: filename,
        version: version,

      }
      + self.withInstallTemplate(
        'jb install %(url)s@%(version)s'
      )
      + self.withUsageTemplate(
        'local %(name)s = import "%(import)s"'
      ),

    '#newSub':: d.fn(|||
      `newSub` creates a package without the preconfigured install/usage templates.

      Arguments:

      * given `name`
      * `help` text
    |||, [
      d.arg('name', d.T.string),
      d.arg('help', d.T.string),
    ]),
    newSub(name, help)::
      {
        name: name,
        help: help,
      },

    withInstallTemplate(template):: {
      installTemplate:
        if template != null
        then
          |||

            ## Install

            ```
            %s
            ```
          ||| % template
        else '',
    },

    withUsageTemplate(template):: {
      usageTemplate:
        if template != null
        then
          |||

            ## Usage

            ```jsonnet
            %s
            ```
          ||| % template
        else '',
    },
  },

  '#pkg':: self.package['#new'] + d.func.withHelp('`new` is a shorthand for `package.new`'),
  pkg:: self.package.new,

  '#object': d.obj('Utilities for documenting Jsonnet objects (`{ }`).'),
  object:: {
    '#new': d.fn('new creates a new object, optionally with description and fields', [d.arg('help', d.T.string), d.arg('fields', d.T.object)]),
    new(help='', fields={}):: { object: {
      help: help,
      fields: fields,
    } },

    '#withFields': d.fn('The `withFields` modifier overrides the fields property of an already created object', [d.arg('fields', d.T.object)]),
    withFields(fields):: { object+: {
      fields: fields,
    } },
  },

  '#obj': self.object['#new'] + d.func.withHelp('`obj` is a shorthand for `object.new`'),
  obj:: self.object.new,

  '#func': d.obj('Utilities for documenting Jsonnet methods (functions of objects)'),
  func:: {
    '#new': d.fn('new creates a new function, optionally with description and arguments', [d.arg('help', d.T.string), d.arg('args', d.T.array)]),
    new(help='', args=[]):: { 'function': {
      help: help,
      args: args,
    } },

    '#withHelp': d.fn('The `withHelp` modifier overrides the help text of that function', [d.arg('help', d.T.string)]),
    withHelp(help):: { 'function'+: {
      help: help,
    } },

    '#withArgs': d.fn('The `withArgs` modifier overrides the arguments of that function', [d.arg('args', d.T.array)]),
    withArgs(args):: { 'function'+: {
      args: args,
    } },
  },

  '#fn': self.func['#new'] + d.func.withHelp('`fn` is a shorthand for `func.new`'),
  fn:: self.func.new,

  '#argument': d.obj('Utilities for creating function arguments'),
  argument:: {
    '#new': d.fn(|||
      `new` creates a new function argument, taking the `name`, the `type`. Optionally it
      can take a `default` value and `enum`-erate potential values.

      Examples:

      ```jsonnet
      [
        d.argument.new('foo', d.T.string),
        d.argument.new('bar', d.T.string, default='loo'),
        d.argument.new('baz', d.T.number, enums=[1,2,3]),
      ]
      ```
    |||, [
      d.arg('name', d.T.string),
      d.arg('type', d.T.string),
      d.arg('default', d.T.any),
      d.arg('enums', d.T.array),
    ]),
    new(name, type, default=null, enums=null): {
      name: name,
      type: type,
      default: default,
      enums: enums,
    },
    '#fromSchema': d.fn(|||
      `fromSchema` creates a new function argument, taking a JSON `schema` to describe the type information for this argument.

      Examples:

      ```jsonnet
      [
        d.argument.fromSchema('foo', { type: 'string' }),
        d.argument.fromSchema('bar', { type: 'string', default='loo' }),
        d.argument.fromSchema('baz', { type: 'number', enum=[1,2,3] }),
      ]
      ```
    |||, [
      d.arg('name', d.T.string),
      d.arg('schema', d.T.object),
    ]),
    fromSchema(name, schema): {
      name: name,
      schema: schema,
    },
  },
  '#arg': self.argument['#new'] + self.func.withHelp('`arg` is a shorthand for `argument.new`'),
  arg:: self.argument.new,

  '#value': d.obj('Utilities for documenting plain Jsonnet values (primitives)'),
  value:: {
    '#new': d.fn('new creates a new object of given type, optionally with description and default value', [d.arg('type', d.T.string), d.arg('help', d.T.string), d.arg('default', d.T.any)]),
    new(type, help='', default=null): { value: {
      help: help,
      type: type,
      default: default,
    } },
  },
  '#val': self.value['#new'] + self.func.withHelp('`val` is a shorthand for `value.new`'),
  val:: self.value.new,

  // T contains constants for the Jsonnet types
  T:: {
    '#string': d.val(d.T.string, 'argument of type "string"'),
    string: 'string',

    '#number': d.val(d.T.string, 'argument of type "number"'),
    number: 'number',
    int: self.number,
    integer: self.number,

    '#boolean': d.val(d.T.string, 'argument of type "boolean"'),
    boolean: 'bool',
    bool: self.boolean,

    '#object': d.val(d.T.string, 'argument of type "object"'),
    object: 'object',

    '#array': d.val(d.T.string, 'argument of type "array"'),
    array: 'array',

    '#any': d.val(d.T.string, 'argument of type "any"'),
    any: 'any',

    '#null': d.val(d.T.string, 'argument of type "null"'),
    'null': 'null',
    nil: self['null'],

    '#func': d.val(d.T.string, 'argument of type "func"'),
    func: 'function',
    'function': self.func,
  },

  '#render': d.fn(
    |||
      `render` converts the docstrings to human readable Markdown files.

      Usage:

      ```jsonnet
      // docs.jsonnet
      d.render(import 'main.libsonnet')
      ```

      Call with: `jsonnet -S -c -m docs/ docs.jsonnet`
    |||,
    args=[
      d.arg('obj', d.T.object),
    ]
  ),
  render:: (import './render.libsonnet').render,

}
