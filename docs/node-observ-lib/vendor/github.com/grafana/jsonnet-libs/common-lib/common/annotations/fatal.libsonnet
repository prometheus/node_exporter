local g = import '../g.libsonnet';
local annotation = g.dashboard.annotation;
local base = import './base.libsonnet';

// Show fatal or critical events as annotations
base {
  new(
    title,
    target,
  ):
    super.new(title, target)
    + annotation.withIconColor('light-purple')
    + annotation.withHide(true),
}
