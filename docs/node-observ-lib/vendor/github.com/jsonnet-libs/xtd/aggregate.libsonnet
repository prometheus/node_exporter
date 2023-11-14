local d = import 'doc-util/main.libsonnet';

{
  local this = self,

  '#': d.pkg(
    name='aggregate',
    url='github.com/jsonnet-libs/xtd/aggregate.libsonnet',
    help=|||
      `aggregate` implements helper functions to aggregate arrays of objects into objects with arrays.

      Example:

      ```jsonnet
      local apps = [
        {
          appid: 'id1',
          name: 'yo',
          id: i,
        }
        for i in std.range(0, 10)
      ];

      aggregate.byKeys(apps, ['appid', 'name']);
      ```

      Output:

      ```json
      {
         "id1": {
            "yo": [
               {
                  "appid": "id1",
                  "id": 0,
                  "name": "yo"
               },
               {
                  "appid": "id1",
                  "id": 1,
                  "name": "yo"
               },
               ...
            ]
         }
      }
      ```
    |||,
  ),

  '#byKey':: d.fn(
    |||
      `byKey` aggregates an array by the value of `key`
    |||,
    [
      d.arg('arr', d.T.array),
      d.arg('key', d.T.string),
    ]
  ),
  byKey(arr, key):
    // find all values of key
    local values = std.set([
      item[key]
      for item in arr
    ]);

    // create the aggregate for the value of each key
    {
      [value]: [
        item
        for item in std.filter(
          function(x)
            x[key] == value,
          arr
        )
      ]
      for value in values
    },

  '#byKeys':: d.fn(
    |||
      `byKey` aggregates an array by iterating over `keys`, each item in `keys` nests the
      aggregate one layer deeper.
    |||,
    [
      d.arg('arr', d.T.array),
      d.arg('keys', d.T.array),
    ]
  ),
  byKeys(arr, keys):
    local aggregate = self.byKey(arr, keys[0]);
    // if last key in keys
    if std.length(keys) == 1

    // then return aggregate
    then aggregate

    // else aggregate with remaining keys
    else {
      [k]: this.byKeys(aggregate[k], keys[1:])
      for k in std.objectFields(aggregate)
    },

}
