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

v="${TMPDIR:-/tmp}/virt.working"
rm -f "${v}" || true
touch "${v}"
if [ -x /usr/sbin/virt-what ]
then
  platforms=$(echo $( virt-what ) | tr '\n' ' ')
  echo '# HELP virt_platform reports one series per detected virtualization type. If no type is detected, the type is "none".' >>"${v}"
  echo '# TYPE virt_platform gauge' >>"${v}"
  count=0
  for platform in ${platforms}; do
    if [[ -z "${platform}" ]]; then
      continue
    fi
    count=$(( count + 1 ))
    echo "virt_platform{type=\"${platform}\"} 1" >>"${v}"
  done

  # Attempt AWS detection on older (m4) and newer (m5) instance types, taken from https://serverfault.com/a/903599
  if ([ -f /sys/hypervisor/uuid ] && [ `head -c 3 /sys/hypervisor/uuid` == "ec2" ]) ||
      ([ -r /sys/devices/virtual/dmi/id/product_uuid ] && [ `head -c 3 /sys/devices/virtual/dmi/id/product_uuid` == "EC2" ]); then
    count=$(( count + 1 ))
    echo "virt_platform{type=\"aws\"} 1" >>"${v}"
  fi
  # Attempt GCP detection
  if dmidecode -s bios-vendor | grep Google; then
    count=$(( count + 1 ))
    echo "virt_platform{type=\"gcp\"} 1" >>"${v}"
  fi

  if [[ "${count}" -eq 0 ]]; then
    echo "virt_platform{type=\"none\"} 1" >>"${v}"
  fi
fi
mv "${v}" virt.prom
