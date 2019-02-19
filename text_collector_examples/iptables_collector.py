#!/usr/bin/env python3

import subprocess
import re
import os
import argparse

# Change to your collector.textfile.directory
prom_outfile = '/var/lib/prometheus/node-exporter/iptables_collector.prom'

tables = ['filter', 'nat', 'mangle', 'raw']
re_chain = re.compile('^Chain')
re_header = re.compile('^pkts')
re_blankline = re.compile('^\s+$')

iptables_packet_lines = []
iptables_byte_lines = []

for table in tables:
  cmd = ['/sbin/iptables', '-L', '-n', '-v', '-x', '-t', table]
  proc = subprocess.Popen(cmd, stdout=subprocess.PIPE)
  for line in proc.stdout.readlines():
    line = line.decode('utf8')

    if re_blankline.match(str(line)):
      continue

    line_pieces = line.split()

    if re_chain.match(str(line_pieces[0])):
      l_chain_name = line_pieces[1]
      continue

    if re_header.match(str(line_pieces[0])):
      continue

    l_packets = line_pieces[0]
    l_bytes = line_pieces[1]
    l_target = line_pieces[2]
    l_prot = line_pieces[3]
    l_in = line_pieces[5]
    l_out = line_pieces[6]
    l_src = line_pieces[7]
    l_dest = line_pieces[8]
    l_options = ' '.join(line_pieces[9:]).replace('"','\\"')

    iptables_packet_lines.append('iptables_packets_total{table="%s",chain="%s",target="%s",prot="%s",in="%s",out="%s",src="%s",dest="%s",opt="%s"} %s' % (table,l_chain_name,l_target,l_prot,l_in,l_out,l_src,l_dest,l_options,l_packets))
    iptables_byte_lines.append('iptables_bytes_total{table="%s",chain="%s",target="%s",prot="%s",in="%s",out="%s",src="%s",dest="%s",opt="%s"} %s' % (table,l_chain_name,l_target,l_prot,l_in,l_out,l_src,l_dest,l_options,l_bytes))


parser = argparse.ArgumentParser()
parser.add_argument('--debug', action='store_true', help='debug iptables_collector output')
args = parser.parse_args()

if args.debug:
    print('# HELP iptables_packets_total packet counters for iptable rules.')
    print('# TYPE iptables_packets_total counter')
    for line in iptables_packet_lines:
        print(line)

    print('# HELP iptables_bytes_total byte counters for iptable rules.')
    print('# TYPE iptables_bytes_total counter')
    for line in iptables_byte_lines:
        print(line)
else:
    with open(prom_outfile, 'w') as prom_out:
        print('# HELP iptables_packets_total packet counters for iptable rules.', file=prom_out)
        print('# TYPE iptables_packets_total counter', file=prom_out)
        for line in iptables_packet_lines:
            print(line, file=prom_out)

        print('# HELP iptables_bytes_total byte counters for iptable rules.', file=prom_out)
        print('# TYPE iptables_bytes_total counter', file=prom_out)
        for line in iptables_byte_lines:
            print(line, file=prom_out)
