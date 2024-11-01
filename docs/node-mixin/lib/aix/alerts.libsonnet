{
  new(this, parentPrometheus):
    {
      groups:
        //keep only alerts listed in alertsKeep
        std.filter(
          function(group) std.length(group.rules) > 0,
          [
            {
              name: group.name,
              rules: [
                rule
                for rule in group.rules
                if std.length(std.find(rule.alert, this.config.alertsKeep)) > 0
              ],
            }
            for group in parentPrometheus.alerts.groups
          ],

        ),

    },
}
