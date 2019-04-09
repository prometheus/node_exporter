#!/usr/bin/env python
#
# Description: Extract chronyd metrics from chronyc -c.
# Author: Aanchal Malhotra <aanchal4@bu.edu>
#
# Works with chrony version 2.4 and higher

import subprocess
import sys

chrony_sourcestats_cmd = ['chronyc', '-c', 'sourcestats']
chrony_source_cmd = ['chronyc', '-c', 'sources']
metrics_fields = [
    "Name/IP Address",
    "NP",
    "NR",
    "Span",
    "Frequency",
    "Freq Skew",
    "Offset",
    "Std Dev"]

status_types = {'x': 0, '?': 1, '-': 2, '+': 3, '*': 4}

metrics_source = {
    "*": "synchronized (system peer)",
    "+": "synchronized",
    "?": "unreachable",
    "x": "Falseticker",
    "-": "reference clock"}

metrics_mode = {
    '^': "server",
    '=': "peer",
    "#": "reference clock"}


def get_cmdoutput(command):
    proc = subprocess.Popen(command, stdout=subprocess.PIPE)
    out, err = proc.communicate()
    return_code = proc.poll()
    if return_code:
        raise RuntimeError('Call to "{}" returned error: \
    {}'.format(command, return_code))
    return out


def printPrometheusformat(metric, values):
    print("# HELP chronyd_%s chronyd metric for %s" % (metric, metric))
    print("# TYPE chronyd_%s gauge" % (metric))
    for labels in values:
        if labels is None:
            print("chronyd_%s %f" % (metric, values[labels]))
        else:
            print("chronyd_%s{%s} %f" % (metric, labels, values[labels]))


def main(argv):
    peer_status_metrics = {}
    offset_metrics = {}
    freq_skew_metrics = {}
    freq_metrics = {}
    std_dev_metrics = {}
    chrony_sourcestats = get_cmdoutput(chrony_sourcestats_cmd)
    for line in chrony_sourcestats.split('\n'):
        if (len(line)) > 0:
            x = line.split(',')
            common_labels = "remote=\"%s\"" % (x[0])
            freq_metrics[common_labels] = float(x[4])
            freq_skew_metrics[common_labels] = float(x[5])
            std_dev_metrics[common_labels] = float(x[7])

    printPrometheusformat('freq_skew_ppm', freq_skew_metrics)
    printPrometheusformat('freq_ppm', freq_metrics)
    printPrometheusformat('std_dev_seconds', std_dev_metrics)

    chrony_source = get_cmdoutput(chrony_source_cmd)
    for line in chrony_source.split('\n'):
        if (len(line)) > 0:
            x = line.split(',')
            stratum = x[3]
            mode = metrics_mode[x[0]]
            common_labels = "remote=\"%s\"" % (x[2])
            peer_labels = "%s,stratum=\"%s\",mode=\"%s\"" % (
                common_labels,
                stratum,
                mode,
            )
            peer_status_metrics[peer_labels] = float(status_types[x[1]])
            offset_metrics[common_labels] = float(x[8])

    printPrometheusformat('peer_status', peer_status_metrics)
    printPrometheusformat('offset_seconds', offset_metrics)


if __name__ == "__main__":
    main(sys.argv[1:])
