#!/usr/local/bin/bash

set -euf -o pipefail

enabled_collectors=$(cat << COLLECTORS
  cpu
  meminfo
  netdev
  textfile
COLLECTORS
)
disabled_collectors=$(cat << COLLECTORS
  filesystem
  uname
  zfs
COLLECTORS
)
cd "$(dirname $0)"

port="$((10000 + (RANDOM % 10000)))"
tmpdir=$(mktemp -d /tmp/node_exporter_e2e_test.XXXXXX)

skip_re="^(go_|node_exporter_build_info|node_boot_time_seconds|node_cpu_seconds_total|node_scrape_collector_duration_seconds|node_network_|node_memory_|node_exec_|node_time_|node_load)"

arch="$(uname -m)"

fixture='collector/fixtures/e2e-output-freebsd.txt';

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
  --path.rootfs="collector/fixtures" \
  --path.procfs="collector/fixtures/proc" \
  --path.sysfs="collector/fixtures/sys" \
  --path.udev.data="collector/fixtures/udev/data" \
  $(for c in ${enabled_collectors}; do echo --collector.${c}  ; done) \
  $(for c in ${disabled_collectors}; do echo --no-collector.${c}  ; done) \
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