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
  platforms="$( virt-what )"
  echo '# HELP virt_platform reports one series per detected virtualization type. If no type is detected, the type is "none".' >>"${v}"
  echo '# TYPE virt_platform gauge' >>"${v}"
  for platform in "${platforms}"; do
    if [[ -z "${platform}" ]]; then
      continue
    fi
    echo "virt_platform{type=\"${platform}\"} 1" >>"${v}"
  done
  if [[ -z "${platforms}" ]]; then
    echo "virt_platform{type=\"none\"} 1" >>"${v}"
  fi
fi
mv "${v}" virt.prom