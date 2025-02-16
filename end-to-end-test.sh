#!/usr/bin/env bash

set -euf -o pipefail

# Allow setting GOHOSTOS for debugging purposes.
GOHOSTOS=${GOHOSTOS:-$(go env GOHOSTOS)}

# Allow setting arch for debugging purposes.
arch=${arch:-$(uname -m)}

maybe_flag_search_scope() {
  local collector=$1
  os_aux_os=""
  if [[ $GOHOSTOS =~ ^(freebsd|openbsd|netbsd|solaris|dragonfly)$ ]]; then
    os_aux_os=" ${collector}_bsd.go"
  fi
  echo "${collector}_common.go ${collector}.go ${collector}_${GOHOSTOS}.go ${collector}_${GOHOSTOS}_${arch}.go${os_aux_os}"
}

supported_collectors() {
  local collectors=$1
  local supported=""
  for collector in ${collectors}; do
    for filename in $(maybe_flag_search_scope "${collector}"); do
      file="collector/${filename}"
      if ./tools/tools match ${file} > /dev/null 2>&1; then
        if grep -h -E -o -- "registerCollector\(" ${file} > /dev/null 2>&1; then
          supported="${supported} ${collector}"
        fi
        break
      fi
    done
  done
  echo "${supported}" | tr ' ' '\n' | sort | uniq
}

enabled_collectors=$(cat << COLLECTORS
  arp
  bcache
  bonding
  btrfs
  buddyinfo
  cgroups
  conntrack
  cpu
  cpufreq
  cpu_vulnerabilities
  diskstats
  dmi
  drbd
  edac
  entropy
  fibrechannel
  filefd
  hwmon
  infiniband
  interrupts
  ipvs
  ksmd
  lnstat
  loadavg
  mdadm
  meminfo
  meminfo_numa
  mountstats
  netdev
  netstat
  nfs
  nfsd
  pressure
  processes
  qdisc
  rapl
  schedstat
  slabinfo
  sockstat
  softirqs
  stat
  sysctl
  textfile
  thermal_zone
  udp_queues
  vmstat
  watchdog
  wifi
  xfrm
  xfs
  zfs
  zoneinfo
COLLECTORS
)
supported_enabled_collectors=$(supported_collectors "${enabled_collectors}")

disabled_collectors=$(cat << COLLECTORS
  selinux
  filesystem
  timex
  uname
COLLECTORS
)
supported_disabled_collectors=$(supported_collectors "${disabled_collectors}")

cd "$(dirname $0)"

port="$((10000 + (RANDOM % 10000)))"
tmpdir=$(mktemp -d /tmp/node_exporter_e2e_test.XXXXXX)

skip_re="^(go_|node_exporter_build_info|node_scrape_collector_duration_seconds|process_|node_textfile_mtime_seconds|node_time_(zone|seconds)|node_network_(receive|transmit)_(bytes|packets)_total)"

case "${arch}" in
  aarch64|ppc64le) fixture_metrics='collector/fixtures/e2e-64k-page-output.txt' ;;
  *) fixture_metrics='collector/fixtures/e2e-output.txt' ;;
esac

# Only test CPU info collection on x86_64.
case "${arch}" in
  x86_64)
    cpu_info_collector='--collector.cpu.info'
    cpu_info_bugs='^(cpu_meltdown|spectre_.*|mds)$'
    cpu_info_flags='^(aes|avx.?|constant_tsc)$'
    ;;
  *)
    cpu_info_collector='--no-collector.cpu.info'
    cpu_info_bugs=''
    cpu_info_flags=''
    ;;
esac

keep=0; update=0; verbose=0
while getopts 'hkuv' opt
do
  case "$opt" in
    k)
      keep=1
      ;;
    u)
      update=1
      ;;
    v)
      verbose=1
      set -x
      ;;
    *)
      echo "Usage: $0 [-k] [-u] [-v]"
      echo "  -k: keep temporary files and leave node_exporter running"
      echo "  -u: update fixture_metrics"
      echo "  -v: verbose output"
      exit 1
      ;;
  esac
done

if [ ! -x ./node_exporter ]
then
    echo './node_exporter not found. Consider running `go build` first.' >&2
    exit 1
fi

collector_flags=$(cat << FLAGS
  ${cpu_info_collector}
  --collector.arp.device-exclude=nope
  --collector.bcache.priorityStats
  --collector.cpu.info.bugs-include=${cpu_info_bugs}
  --collector.cpu.info.flags-include=${cpu_info_flags}
  --collector.hwmon.chip-include=(applesmc|coretemp|hwmon4|nct6779)
  --collector.netclass.ignore-invalid-speed
  --collector.netclass.ignored-devices=(dmz|int)
  --collector.netdev.device-include=lo
  --collector.qdisc.device-include=(wlan0|eth0)
  --collector.qdisc.fixtures=collector/fixtures/qdisc/
  --collector.stat.softirq
  --collector.sysctl.include-info=kernel.seccomp.actions_avail
  --collector.sysctl.include=fs.file-nr
  --collector.sysctl.include=fs.file-nr:total,current,max
  --collector.sysctl.include=kernel.threads-max
  --collector.textfile.directory=collector/fixtures/textfile/two_metric_files/
  --collector.wifi.fixtures=collector/fixtures/wifi
  --no-collector.arp.netlink
FLAGS
)

# Handle supported --[no-]collector.<name> flags. These are not hardcoded.
_filtered_collector_flags=""
for flag in ${collector_flags}; do
  collector=$(echo "${flag}" | cut -d"." -f2)
  # If the flag is associated with an enabled-by-default collector, include it.
  enabled_by_default=0
  for filename in $(maybe_flag_search_scope "${collector}") ; do
      file="collector/${filename}"
      if grep -h -E -o -- "registerCollector\(.*, defaultEnabled" ${file} > /dev/null 2>&1; then
        _filtered_collector_flags="${_filtered_collector_flags} ${flag}"
        enabled_by_default=1
        break
      fi
  done
  if [ ${enabled_by_default} -eq 1 ]; then
    continue
  fi
  # If the flag is associated with an enabled-list collector, include it.
  if echo "${supported_enabled_collectors} ${supported_disabled_collectors}" | grep -q -w "${collector}"; then
    _filtered_collector_flags="${_filtered_collector_flags} ${flag}"
  fi
done

# Handle supported --[no-]collector.<name>.<collector> flags. These are hardcoded and matched by the expression below.
filtered_collector_flags=""
# Check flags of all supported collectors further down their sub-collectors (beyond the 2nd ".").
for flag in ${_filtered_collector_flags}; do
  # Iterate through all possible files where the flag may be defined.
  flag_collector="$(echo "${flag}" | cut -d"." -f2)"
  for filename in $(maybe_flag_search_scope "${flag_collector}") ; do
    file="collector/${filename}"
    # Move to next iteration if the current file is not included under the build context.
    if ! ./tools/tools match "$file" > /dev/null 2>&1; then
     continue
    fi
    # Flag has the format: --[no-]collector.<name>.<collector>.
    if [ -n "$(echo ${flag} | cut -d"." -f3)" ]; then
      # Check if the flag is used in the file.
      trimmed_flag=$(echo "${flag}" | tr -d "\"' " | cut -d"=" -f1 | cut -c 3-)
      if [[ $trimmed_flag =~ ^no- ]]; then
        trimmed_flag=$(echo $trimmed_flag | cut -c 4-)
      fi
      if grep -h -E -o -- "kingpin.Flag\(\"${trimmed_flag}" ${file} > /dev/null 2>&1; then
        filtered_collector_flags="${filtered_collector_flags} ${flag}"
      else
       continue
      fi
    # Flag has the format: --[no-]collector.<name>.
    else
      # Flag is supported by the host.
      filtered_collector_flags="${filtered_collector_flags} ${flag}"
    fi
  done
done

# Check for ignored flags.
ignored_flags=""
for flag in ${collector_flags}; do
  flag=$(echo "${flag}" | tr -d " ")
  if ! echo "${filtered_collector_flags}" | grep -q -F -- "${flag}" > /dev/null 2>&1; then
    ignored_flags="${ignored_flags} ${flag}"
  fi
done

echo "ENABLED COLLECTORS======="
echo "${supported_enabled_collectors:1}" | tr ' ' '\n' | sort
echo "========================="

echo "DISABLED COLLECTORS======"
echo "${supported_disabled_collectors:1}" | tr ' ' '\n' | sort
echo "========================="

echo "IGNORED FLAGS============"
echo "${ignored_flags:1}"| tr ' ' '\n' | sort | uniq
echo "========================="

./node_exporter \
  --path.rootfs="collector/fixtures" \
  --path.procfs="collector/fixtures/proc" \
  --path.sysfs="collector/fixtures/sys" \
  --path.udev.data="collector/fixtures/udev/data" \
  $(for c in ${supported_enabled_collectors}; do echo --collector.${c}  ; done) \
  $(for c in ${supported_disabled_collectors}; do echo --no-collector.${c}  ; done) \
  ${filtered_collector_flags} \
  --web.listen-address "127.0.0.1:${port}" \
  --log.level="debug" > "${tmpdir}/node_exporter.log" 2>&1 &

echo $! > "${tmpdir}/node_exporter.pid"

generated_metrics="${tmpdir}/e2e-output.txt"
for os in freebsd openbsd netbsd solaris dragonfly darwin; do
  if [ "${GOHOSTOS}" = "${os}" ]; then
    generated_metrics="${tmpdir}/e2e-output-${GOHOSTOS}.txt"
    fixture_metrics="${fixture_metrics::-4}-${GOHOSTOS}.txt"
  fi
done

finish() {
  if [ $? -ne 0 -o ${verbose} -ne 0 ]
  then
    cat << EOF >&2
LOG =====================
$(cat "${tmpdir}/node_exporter.log")
=========================
EOF
  fi

  if [ ${update} -ne 0 ]
  then
    cp "${generated_metrics}" "${fixture_metrics}"
  fi

  if [ ${keep} -eq 0 ]
  then
    kill -9 "$(cat ${tmpdir}/node_exporter.pid)"
    # This silences the "Killed" message
    set +e
    wait "$(cat ${tmpdir}/node_exporter.pid)" > /dev/null 2>&1
    rm -rf "${tmpdir}"
  fi
}

trap finish EXIT

get() {
  if command -v curl > /dev/null 2>&1
  then
    curl -s -f "$@"
  elif command -v wget > /dev/null 2>&1
  then
    wget -O - "$@"
  else
    echo "Neither curl nor wget found"
    exit 1
  fi
}

sleep 1

get "127.0.0.1:${port}/metrics" | grep --text -E -v "${skip_re}" > "${generated_metrics}"

# The following ignore-list is only applicable to the VMs used to run E2E tests on platforms for which containerized environments are not available.
# However, owing to this, there are some non-deterministic metrics that end up generating samples, unlike their containerized counterparts, for e.g., node_network_receive_bytes_total. 
non_deterministic_metrics=$(cat << METRICS
  node_boot_time_seconds
  node_cpu_frequency_hertz
  node_cpu_seconds_total
  node_disk_io_time_seconds_total
  node_disk_read_bytes_total
  node_disk_read_sectors_total
  node_disk_read_time_seconds_total
  node_disk_reads_completed_total
  node_disk_write_time_seconds_total
  node_disk_writes_completed_total
  node_disk_written_bytes_total
  node_disk_written_sectors_total
  node_exec_context_switches_total
  node_exec_device_interrupts_total
  node_exec_forks_total
  node_exec_software_interrupts_total
  node_exec_system_calls_total
  node_exec_traps_total
  node_interrupts_total
  node_load1
  node_load15
  node_load5
  node_memory_active_bytes
  node_memory_buffer_bytes
  node_memory_cache_bytes
  node_memory_compressed_bytes
  node_memory_free_bytes
  node_memory_inactive_bytes
  node_memory_internal_bytes
  node_memory_laundry_bytes
  node_memory_purgeable_bytes
  node_memory_size_bytes
  node_memory_swapped_in_bytes_total
  node_memory_swapped_out_bytes_total
  node_memory_wired_bytes
  node_netstat_tcp_receive_packets_total
  node_netstat_tcp_transmit_packets_total
  node_network_receive_bytes_total
  node_network_receive_multicast_total
  node_network_transmit_multicast_total
METRICS
)

# Remove non-deterministic metrics from the generated metrics file (as we run their workflows in VMs).
for os in freebsd openbsd netbsd solaris dragonfly darwin; do
  if [ "${GOHOSTOS}" = "${os}" ]; then
    for metric in ${non_deterministic_metrics}; do
      sed -i "/${metric}/d" "${generated_metrics}"
    done
  fi
done

diff -u \
  "${fixture_metrics}" \
  "${generated_metrics}"
