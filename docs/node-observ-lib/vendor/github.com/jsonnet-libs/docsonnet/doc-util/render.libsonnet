{
  local root = self,

  render(obj):
    assert std.isObject(obj) && '#' in obj : 'error: object is not a docsonnet package';
    local package = self.package(obj);
    package.toFiles(),

  findPackages(obj, path=[]): {
    local find(obj, path, parentWasPackage=true) =
      std.foldl(
        function(acc, k)
          acc
          + (
            // If matches a package but warn if also has an object docstring
            if '#' in obj[k] && '#' + k in obj
               && !std.objectHasAll(obj[k]['#'], 'ignore')
            then std.trace(
              'warning: %s both defined as object and package' % k,
              [root.package(obj[k], path + [k], parentWasPackage)]
            )
            // If matches a package, return it
            else if '#' in obj[k]
                    && !std.objectHasAll(obj[k]['#'], 'ignore')
            then [root.package(obj[k], path + [k], parentWasPackage)]
            // If not, keep looking
            else find(obj[k], path + [k], parentWasPackage=false)
          ),
        std.filter(
          function(k)
            !std.startsWith(k, '#')
            && std.isObject(obj[k]),
          std.objectFieldsAll(obj)
        ),
        []
      ),

    packages: find(obj, path),

    hasPackages(): std.length(self.packages) > 0,

    toIndex(relativeTo=[]):
      if self.hasPackages()
      then
        std.join('\n', [
          '* ' + p.link(relativeTo)
          for p in self.packages
        ])
        + '\n'
      else '',

    toFiles():
      std.foldl(
        function(acc, p)
          acc
          + { [p.path]: p.toString() }
          + p.packages.toFiles(),
        self.packages,
        {}
      ),
  },

  package(obj, path=[], parentWasPackage=true): {
    local this = self,
    local doc = obj['#'],

    packages: root.findPackages(obj, path),
    fields: root.fields(obj),

    local pathsuffix =
      (if self.packages.hasPackages()
       then '/index.md'
       else '.md'),

    // filepath on disk
    path:
      std.join('/', path)
      + pathsuffix,

    link(relativeTo):
      local relativepath = root.util.getRelativePath(path, relativeTo);
      '[%s](%s)' % [
        std.join('.', relativepath),
        std.join('/', relativepath)
        + pathsuffix,
      ],

    toFiles():
      { 'README.md': this.toString() }
      + self.packages.toFiles(),

    toString():
      std.join(
        '\n',
        [
          '# ' + doc.name + '\n',
          std.get(doc, 'help', ''),
          '',
        ]
        + (if self.packages.hasPackages()
           then [
             '## Subpackages\n\n'
             + self.packages.toIndex(path),
           ]
           else [])
        + (if self.fields.hasFields()
           then [
             '## Index\n\n'
             + self.fields.toIndex()
             + '\n## Fields\n'
             + self.fields.toString(),
           ]
           else [])
      ),
  },

  fields(obj, path=[]): {
    values: root.findValues(obj, path),
    functions: root.findFunctions(obj, path),
    objects: root.findObjects(obj, path),

    hasFields():
      std.any([
        self.values.hasFields(),
        self.functions.hasFields(),
        self.objects.hasFields(),
      ]),

    toIndex():
      std.join('', [
        self.functions.toIndex(),
        self.objects.toIndex(),
      ]),

    toString():
      std.join('', [
        self.values.toString(),
        self.functions.toString(),
        self.objects.toString(),
      ]),
  },

  findObjects(obj, path=[]): {
    local keys =
      std.filter(
        root.util.filter('object', obj),
        std.objectFieldsAll(obj)
      ),

    local undocumentedKeys =
      std.filter(
        function(k)
          std.all([
            !std.startsWith(k, '#'),
            std.isObject(obj[k]),
            !std.objectHasAll(obj[k], 'ignore'),
            !('#' + k in obj),  // not documented in parent
            !('#' in obj[k]),  // not a sub package
          ]),
        std.objectFieldsAll(obj)
      ),

    objects:
      std.foldl(
        function(acc, k)
          acc + [
            root.obj(
              root.util.realkey(k),
              obj[k],
              obj[root.util.realkey(k)],
              path,
            ),
          ],
        keys,
        []
      )
      + std.foldl(
        function(acc, k)
          local o = root.obj(
            k,
            { object: { help: '' } },
            obj[k],
            path,
          );
          acc
          + (if o.fields.hasFields()
             then [o]
             else []),
        undocumentedKeys,
        []
      ),

    hasFields(): std.length(self.objects) > 0,

    toIndex():
      if self.hasFields()
      then
        std.join('', [
          std.join(
            '',
            [' ' for d in std.range(0, (std.length(path) * 2) - 1)]
            + ['* ', f.link]
            + ['\n']
            + (if f.fields.hasFields()
               then [f.fields.toIndex()]
               else [])
          )
          for f in self.objects
        ])
      else '',

    toString():
      if self.hasFields()
      then
        std.join('', [
          o.toString()
          for o in self.objects
        ])
      else '',
  },

  obj(name, doc, obj, path): {
    fields: root.fields(obj, path + [name]),

    path: std.join('.', path + [name]),
    fragment: root.util.fragment(std.join('', path + [name])),
    link: '[`obj %s`](#obj-%s)' % [name, self.fragment],

    toString():
      std.join(
        '\n',
        [root.util.title('obj ' + self.path, std.length(path) + 2)]
        + (if std.get(doc.object, 'help', '') != ''
           then [doc.object.help]
           else [])
        + [self.fields.toString()]
      ),
  },

  findFunctions(obj, path=[]): {
    local keys =
      std.filter(
        root.util.filter('function', obj),
        std.objectFieldsAll(obj)
      ),

    functions:
      std.foldl(
        function(acc, k)
          acc + [
            root.func(
              root.util.realkey(k),
              obj[k],
              path,
            ),
          ],
        keys,
        []
      ),

    hasFields(): std.length(self.functions) > 0,

    toIndex():
      if self.hasFields()
      then
        std.join('', [
          std.join(
            '',
            [' ' for d in std.range(0, (std.length(path) * 2) - 1)]
            + ['* ', f.link]
            + ['\n']
          )
          for f in self.functions
        ])
      else '',

    toString():
      if self.hasFields()
      then
        std.join('', [
          f.toString()
          for f in self.functions
        ])
      else '',
  },

  func(name, doc, path): {
    path: std.join('.', path + [name]),
    fragment: root.util.fragment(std.join('', path + [name])),
    link: '[`fn %s(%s)`](#fn-%s)' % [name, self.args, self.fragment],

    local getType(arg) =
      local type =
        if 'schema' in arg
        then std.get(arg.schema, 'type', '')
        else std.get(arg, 'type', '');
      if std.isArray(type)
      then std.join(',', [t for t in type])
      else type,

    // Use BelRune as default can be 'null' as a value. Only supported for arg.schema, arg.default didn't support this, not sure how to support without breaking asssumptions downstream.
    local BelRune = std.char(7),
    local getDefault(arg) =
      if 'schema' in arg
      then std.get(arg.schema, 'default', BelRune)
      else
        local d = std.get(arg, 'default', BelRune);
        if d == null
        then BelRune
        else d,

    local getEnum(arg) =
      if 'schema' in arg
      then std.get(arg.schema, 'enum', [])
      else
        local d = std.get(arg, 'enums', []);
        if d == null
        then []
        else d,

    args:
      std.join(', ', [
        local default = getDefault(arg);
        if default != BelRune
        then std.join('=', [
          arg.name,
          std.manifestJson(default),
        ])
        else arg.name
        for arg in doc['function'].args
      ]),


    args_list:
      if std.length(doc['function'].args) > 0
      then
        '\nPARAMETERS:\n\n'
        + std.join('\n', [
          '* **%s** (`%s`)' % [arg.name, getType(arg)]
          + (
            local default = getDefault(arg);
            if default != BelRune
            then '\n   - default value: `%s`' % std.manifestJson(default)
            else ''
          )
          + (
            local enum = getEnum(arg);
            if enum != []
            then
              '\n   - valid values: %s' %
              std.join(', ', [
                '`%s`' % std.manifestJson(item)
                for item in enum
              ])
            else ''
          )
          for arg in doc['function'].args
        ])
      else '',

    toString():
      std.join('\n', [
        root.util.title('fn ' + self.path, std.length(path) + 2),
        |||
          ```jsonnet
          %s(%s)
          ```
          %s
        ||| % [self.path, self.args, self.args_list],
        std.get(doc['function'], 'help', ''),
      ]),
  },

  findValues(obj, path=[]): {
    local keys =
      std.filter(
        root.util.filter('value', obj),
        std.objectFieldsAll(obj)
      ),

    values:
      std.foldl(
        function(acc, k)
          acc + [
            root.val(
              root.util.realkey(k),
              obj[k],
              obj[root.util.realkey(k)],
              path,
            ),
          ],
        keys,
        []
      ),

    hasFields(): std.length(self.values) > 0,

    toString():
      if self.hasFields()
      then
        std.join('\n', [
          '* ' + f.toString()
          for f in self.values
        ]) + '\n'
      else '',
  },

  val(name, doc, obj, path): {
    toString():
      std.join(' ', [
        '`%s`' % std.join('.', path + [name]),
        '(`%s`):' % doc.value.type,
        '`"%s"`' % obj,
        '-',
        std.get(doc.value, 'help', ''),
      ]),
  },

  util: {
    realkey(key):
      assert std.startsWith(key, '#') : 'Key %s not a docstring key' % key;
      key[1:],
    title(title, depth=0):
      std.join(
        '',
        ['\n']
        + ['#' for i in std.range(0, depth)]
        + [' ', title, '\n']
      ),
    fragment(title):
      std.asciiLower(
        std.strReplace(
          std.strReplace(title, '.', '')
          , ' ', '-'
        )
      ),
    filter(type, obj):
      function(k)
        std.all([
          std.startsWith(k, '#'),
          std.isObject(obj[k]),
          !std.objectHasAll(obj[k], 'ignore'),
          type in obj[k],
          root.util.realkey(k) in obj,
        ]),

    getRelativePath(path, relativeTo):
      local shortest = std.min(std.length(relativeTo), std.length(path));

      local commonIndex =
        std.foldl(
          function(acc, i) (
            if acc.stop
            then acc
            else
              acc + {
                // stop count if path diverges
                local stop = relativeTo[i] != path[i],
                stop: stop,
                count+: if stop then 0 else 1,
              }
          ),
          std.range(0, shortest - 1),
          { stop: false, count: 0 }
        ).count;

      local _relativeTo = relativeTo[commonIndex:];
      local _path = path[commonIndex:];

      // prefix for relative difference
      local prefix = ['..' for i in std.range(0, std.length(_relativeTo) - 1)];

      // return path with prefix
      prefix + _path,
  },
}
