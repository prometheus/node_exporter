#!/usr/bin/env bash

set -euf -o pipefail

enabled_collectors=$(cat << COLLECTORS
  arp
  bcache
  buddyinfo
  conntrack
  cpu
  cpufreq
  diskstats
  drbd
  edac
  entropy
  filefd
  hwmon
  infiniband
  interrupts
  ipvs
  ksmd
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
  qdisc
  schedstat
  sockstat
  stat
  thermal_zone
  textfile
  bonding
  vmstat
  wifi
  xfs
  zfs
  processes
COLLECTORS
)
disabled_collectors=$(cat << COLLECTORS
  filesystem
  time
  timex
  uname
COLLECTORS
)
cd "$(dirname $0)"

port="$((10000 + (RANDOM % 10000)))"
tmpdir=$(mktemp -d /tmp/node_exporter_e2e_test.XXXXXX)

skip_re="^(go_|node_exporter_build_info|node_scrape_collector_duration_seconds|process_|node_textfile_mtime_seconds)"

arch="$(uname -m)"

case "${arch}" in
  aarch64|ppc64le) fixture='collector/fixtures/e2e-64k-page-output.txt' ;;
  *) fixture='collector/fixtures/e2e-output.txt' ;;
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
      echo "  -u: update fixture"
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

./node_exporter \
  --path.procfs="collector/fixtures/proc" \
  --path.sysfs="collector/fixtures/sys" \
  $(for c in ${enabled_collectors}; do echo --collector.${c}  ; done) \
  $(for c in ${disabled_collectors}; do echo --no-collector.${c}  ; done) \
  --collector.textfile.directory="collector/fixtures/textfile/two_metric_files/" \
  --collector.wifi.fixtures="collector/fixtures/wifi" \
  --collector.qdisc.fixtures="collector/fixtures/qdisc/" \
  --collector.netclass.ignored-devices="(bond0|dmz|int)" \
  --web.listen-address "127.0.0.1:${port}" \
  --log.level="debug" > "${tmpdir}/node_exporter.log" 2>&1 &

echo $! > "${tmpdir}/node_exporter.pid"

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
    cp "${tmpdir}/e2e-output.txt" "${fixture}"
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

get "127.0.0.1:${port}/metrics" | grep -E -v "${skip_re}" > "${tmpdir}/e2e-output.txt"

diff -u \
  "${fixture}" \
  "${tmpdir}/e2e-output.txt"
