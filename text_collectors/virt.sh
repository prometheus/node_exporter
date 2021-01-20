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

v=$(mktemp --suffix=.virt.working)
chmod uga+r "${v}"

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
  if ! grep -q -s -F '{type="aws"}' "${v}"; then
    if ([ -f /sys/hypervisor/uuid ] && [ `head -c 3 /sys/hypervisor/uuid` == "ec2" ]) ||
        ([ -r /sys/devices/virtual/dmi/id/product_uuid ] && [ `head -c 3 /sys/devices/virtual/dmi/id/product_uuid` == "EC2" ]); then
      count=$(( count + 1 ))
      echo "virt_platform{type=\"aws\"} 1" >>"${v}"
    fi
  fi
  # Attempt GCP detection
  if dmidecode -s bios-vendor | grep Google; then
    count=$(( count + 1 ))
    echo "virt_platform{type=\"gcp\"} 1" >>"${v}"
  fi
  # Attempt OpenStack detection
  if dmidecode | grep "Product Name: OpenStack Compute"; then
    count=$(( count + 1 ))
    echo "virt_platform{type=\"openstack\"} 1" >>"${v}"
  fi

  if [[ "${count}" -eq 0 ]]; then
    line="$( printf 'virt_platform{type="none",bios_vendor="%q",bios_version="%q",system_manufacturer="%q",system_product_name="%q",system_version="%q",baseboard_manufacturer="%q",baseboard_product_name="%q"} 1' \
      "$( dmidecode -s bios-vendor )" "$( dmidecode -s bios-version )" \
      "$( dmidecode -s system-manufacturer )" "$( dmidecode -s system-product-name )" "$( dmidecode -s system-version )" \
      "$( dmidecode -s baseboard-manufacturer )" "$( dmidecode -s baseboard-product-name )" )"
    # Remove all escape characters as they are incorrect (as per https://github.com/prometheus/common/blob/master/expfmt/text_parse.go#L566-L571)
    echo "${line//\\/ }" >>"${v}"
  fi
fi
mv "${v}" virt.prom
