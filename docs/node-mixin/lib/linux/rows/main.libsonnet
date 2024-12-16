local g = import '../../g.libsonnet';
local commonlib = import 'common-lib/common/main.libsonnet';
{
  new(this):: {
    linux: (import './linux.libsonnet').new(this),
    use: (import './use.libsonnet').new(this),
  },
}
