#!/usr/bin/env python3

# Collect per-device btrfs filesystem errors.
# Designed to work on Debian and Centos 6 (with python2.6).

import collections
import glob
import os
import re
import subprocess

def get_btrfs_mount_points():
    """List all btrfs mount points.

    Yields:
        (string) filesystem mount points.
    """
    with open("/proc/mounts") as f:
        for line in f:
            parts = line.split()
            if parts[2] == "btrfs":
                yield parts[1]

def get_btrfs_errors(mountpoint):
    """Get per-device errors for a btrfs mount point.

    Args:
        mountpoint: (string) path to a mount point.

    Yields:
        (device, error_type, error_count) tuples, where:
            device: (string) path to block device.
            error_type: (string) type of btrfs error.
            error_count: (int) number of btrfs errors of a given type.
    """
    p = subprocess.Popen(["btrfs", "device", "stats", mountpoint],
                         stdout=subprocess.PIPE)
    (stdout, stderr) = p.communicate()
    if p.returncode != 0:
        raise RuntimeError("btrfs returned exit code %d" % p.returncode)
    for line in stdout.splitlines():
        if line == '':
            continue
        # Sample line:
        # [/dev/vdb1].flush_io_errs   0
        m = re.search(r"^\[([^\]]+)\]\.(\S+)\s+(\d+)$", line.decode("utf-8"))
        if not m:
            raise RuntimeError("unexpected output from btrfs: '%s'" % line)
        yield m.group(1), m.group(2), int(m.group(3))

def btrfs_error_metrics():
    """Collect btrfs error metrics.

    Returns:
        a list of strings to be exposed as Prometheus metrics.
    """
    metric = "node_btrfs_errors_total"
    contents = [
        "# TYPE %s counter" % metric,
        "# HELP %s number of btrfs errors" % metric,
    ]
    errors_by_device = collections.defaultdict(dict)
    for mountpoint in get_btrfs_mount_points():
        for device, error_type, error_count in get_btrfs_errors(mountpoint):
            contents.append(
                '%s{mountpoint="%s",device="%s",type="%s"} %d' %
                (metric, mountpoint, device, error_type, error_count))

    if len(contents) > 2:
        # return metrics if there are actual btrfs filesystems found
        # (i.e. `contents` contains more than just TYPE and HELP).
        return contents

def btrfs_allocation_metrics():
    """Collect btrfs allocation metrics.

    Returns:
        a list of strings to be exposed as Prometheus metrics.
    """
    prefix = 'node_btrfs_allocation'
    metric_to_filename = {
        'size_bytes': 'total_bytes',
        'used_bytes': 'bytes_used',
        'reserved_bytes': 'bytes_reserved',
        'pinned_bytes': 'bytes_pinned',
        'disk_size_bytes': 'disk_total',
        'disk_used_bytes': 'disk_used',
    }
    contents = []
    for m, f in metric_to_filename.items():
        contents += [
            "# TYPE %s_%s gauge" % (prefix, m),
            "# HELP %s_%s btrfs allocation data (%s)" % (prefix, m, f),
        ]

    for alloc in glob.glob("/sys/fs/btrfs/*/allocation"):
        fs = alloc.split('/')[4]
        for type_ in ('data', 'metadata', 'system'):
            for m, f in metric_to_filename.items():
                filename = os.path.join(alloc, type_, f)
                with open(filename) as f:
                    value = int(f.read().strip())
                    contents.append('%s_%s{fs="%s",type="%s"} %d' % (
                        prefix, m, fs, type_, value))
    if len(contents) > 2*len(metric_to_filename):
        return contents

if __name__ == "__main__":
    contents = ((btrfs_error_metrics() or []) +
                (btrfs_allocation_metrics() or []))

    print("\n".join(contents))
