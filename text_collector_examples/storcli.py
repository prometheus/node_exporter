#!/usr/bin/env python

# Script to parse StorCLI's JSON output and expose
# MegaRAID health as Prometheus metrics.
#
# Tested against StorCLI 'Ver 1.14.12 Nov 25, 2014'.
#
# StorCLI reference manual:
# http://docs.avagotech.com/docs/12352476
#
# Advanced Software Options (ASO) not exposed as metrics currently.
#
# JSON key abbreviations used by StorCLI are documented in the standard command
# output, i.e.  when you omit the trailing 'J' from the command.

import argparse
import json
import subprocess

DESCRIPTION = """Parses StorCLI's JSON output and exposes MegaRAID health as
    Prometheus metrics."""
VERSION = '0.0.1'

METRIC_PREFIX = 'megaraid_'
METRIC_CONTROLLER_LABELS = '{{controller="{}", model="{}"}}'


def main(args):
    data = json.loads(get_storcli_json(args.storcli_path))

    # It appears that the data we need will always be present in the first
    # item in the Controllers array
    status = data['Controllers'][0]

    metrics = {
            'status_code': status['Command Status']['Status Code'],
            'controllers': status['Response Data']['Number of Controllers'],
            }

    for name, value in metrics.iteritems():
        print("{}{} {}".format(METRIC_PREFIX, name, value))

    controller_info = []
    controller_metrics = {}
    overview = []

    try:
        overview = status['Response Data']['System Overview']
    except KeyError:
        pass

    for controller in overview:
        controller_index = controller['Ctl']
        model = controller['Model']
        controller_info.append(METRIC_CONTROLLER_LABELS.format(controller_index, model))

        controller_metrics = {
                # FIXME: Parse dimmer switch options
                # 'dimmer_switch':          controller['DS'],

                'battery_backup_healthy':   int(controller['BBU'] == 'Opt'),
                'degraded':                 int(controller['Hlth'] == 'Dgd'),
                'drive_groups':             controller['DGs'],
                'emergency_hot_spare':      int(controller['EHS'] == 'Y'),
                'failed':                   int(controller['Hlth'] == 'Fld'),
                'healthy':                  int(controller['Hlth'] == 'Opt'),
                'physical_drives':          controller['PDs'],
                'ports':                    controller['Ports'],
                'scheduled_patrol_read':    int(controller['sPR'] == 'On'),
                'virtual_drives':           controller['VDs'],

                # Reverse StorCLI's logic to make metrics consistent
                'drive_groups_optimal':     int(controller['DNOpt'] == 0),
                'virtual_drives_optimal':   int(controller['VNOpt'] == 0),
                }

    for name, value in controller_metrics.iteritems():
        print('{}{}{{controller="{}"}} {}'.format(METRIC_PREFIX, name, controller_index, value))

    for labels in controller_info:
        print('{}{}{} {}'.format(METRIC_PREFIX, 'controller_info', labels, 1))


def get_storcli_json(storcli_path):
    storcli_cmd = [storcli_path, 'show', 'all', 'J']
    proc = subprocess.Popen(storcli_cmd, stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE)
    return proc.communicate()[0]

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=DESCRIPTION,
                                        formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument('--storcli_path',
                        default='/opt/MegaRAID/storcli/storcli64',
                        help='path to StorCLi binary')
    parser.add_argument('--version',
                        action='version',
                        version='%(prog)s {}'.format(VERSION))
    args = parser.parse_args()

    main(args)
