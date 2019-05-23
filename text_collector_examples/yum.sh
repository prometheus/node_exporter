#!/bin/bash
#
# Description: Expose metrics from yum updates.
#
# Author: Slawomir Gonet <slawek@otwiera.cz>
# 
# Based on apt.sh by Ben Kochie <superq@gmail.com>

upgrades=$(/usr/bin/yum -q check-updates | awk 'BEGIN { mute=1 } /Obsoleting Packages/ { mute=0 } mute { print }' | egrep '^\w+\.\w+' | awk '{print $3}' | sort | uniq -c | awk '{print "yum_upgrades_pending{origin=\""$2"\"} "$1}')

echo '# HELP yum_upgrades_pending Yum package pending updates by origin.'
echo '# TYPE yum_upgrades_pending gauge'
if [[ -n "${upgrades}" ]] ; then
  echo "${upgrades}"
else
  echo 'yum_upgrades_pending{origin=""} 0'
fi

