#!/usr/bin/env bash

set -euf -o pipefail

collectors=$(cat << COLLECTORS
  conntrack
  diskstats
  entropy
  filefd
  hwmon
  ksmd
  loadavg
  mdadm
  meminfo
  meminfo_numa
  netdev
  netstat
  nfs
  sockstat
  stat
  textfile
  bonding
  megacli
COLLECTORS
)

cd "$(dirname $0)"

port="$((10000 + (RANDOM % 10000)))"
tmpdir=$(mktemp -d /tmp/node_exporter_e2e_test.XXXXXX)

skip_re="^(go_|node_exporter_|process_|node_textfile_mtime)"

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
  -collector.procfs="collector/fixtures/proc" \
  -collector.sysfs="collector/fixtures/sys" \
  -collectors.enabled="$(echo ${collectors} | tr ' ' ',')" \
  -collector.textfile.directory="collector/fixtures/textfile/two_metric_files/" \
  -collector.megacli.command="collector/fixtures/megacli" \
  -web.listen-address "127.0.0.1:${port}" \
  -log.level="debug" > "${tmpdir}/node_exporter.log" 2>&1 &

echo $! > "${tmpdir}/node_exporter.pid"

finish() {
  if [ ${verbose} -ne 0 ]
  then
    echo "LOG ====================="
    cat "${tmpdir}/node_exporter.log"
    echo "========================="
  fi

  if [ ${update} -ne 0 ]
  then
    cp "${tmpdir}/e2e-output.txt" "collector/fixtures/e2e-output.txt"
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
  if which curl > /dev/null 2>&1
  then
    curl -s -f "$@"
  elif which wget > /dev/null 2>&1
  then
    wget -O - "$@"
  else
    echo "Neither curl nor wget found"
    exit 1
  fi
}

sleep 1

get "127.0.0.1:${port}/metrics" > "${tmpdir}/e2e-output.txt"

diff -u \
  <(grep -E -v "${skip_re}" "collector/fixtures/e2e-output.txt") \
  <(grep -E -v "${skip_re}" "${tmpdir}/e2e-output.txt")
