#!/bin/bash
#
#
# Description: Expose metrics for detected platform virtualization
# Dependencies: virt-what (packages)
#
# The script creates one metric series for each detected virtualization platform
# reported by virt-what.
#
# Author: Clayton Coleman <smarterclayton@gmail.com>

set -o errexit
set -o nounset
set -o pipefail

v="${TMPDIR}/virt.working"
rm -f "${v}" || true
touch "${v}"
if [ -x /usr/sbin/virt-what ]
then
  for platform in $( virt-what ); do
    if [[ -z "${platform}" ]]; then
      continue
    fi
    echo "# HELP virt_platform reports one series per detected virtualization type" >"${v}"
    echo "# TYPE virt_platform gauge" >"${v}"
    echo "virt_platform{type=\"${platform}\"} 1" >"${v}"
  done
fi
mv "${v}" virt.prom