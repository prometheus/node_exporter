#!/usr/bin/env python

import subprocess
import re


tables = ['filter', 'nat', 'mangle', 'raw']
re_chain = re.compile('^Chain')
re_header = re.compile('^pkts')
re_blankline = re.compile('^\s+$')

iptables_packet_lines = []
iptables_byte_lines = []

for table in tables:
  cmd = ['iptables', '-L', '-v', '-x', '-t', table]
  proc = subprocess.Popen(cmd, stdout=subprocess.PIPE)
  for line in proc.stdout.readlines():
    if re_blankline.match(line):
      continue
  
    line_pieces = line.split()
  
    if re_chain.match(line_pieces[0]):
      l_chain_name = line_pieces[1]
      continue 
  
    if re_header.match(line_pieces[0]):
      continue
  
    l_packets = line_pieces[0]
    l_bytes = line_pieces[1]
    l_target = line_pieces[2]
    l_prot = line_pieces[3]
    l_in = line_pieces[5]
    l_out = line_pieces[6]
    l_src = line_pieces[7]
    l_dest = line_pieces[8]
    l_options = ' '.join(line_pieces[9:])
  
    iptables_packet_lines.append('iptables_packets{table="%s",chain="%s",target="%s",prot="%s",in="%s",out="%s",src="%s",dest="%s",opt="%s"} %s' % (table,l_chain_name,l_target,l_prot,l_in,l_out,l_src,l_dest,l_options,l_packets))
    iptables_byte_lines.append('iptables_bytes{table="%s",chain="%s",target="%s",prot="%s",in="%s",out="%s",src="%s",dest="%s",opt="%s"} %s' % (table,l_chain_name,l_target,l_prot,l_in,l_out,l_src,l_dest,l_options,l_bytes))

print '# HELP iptables_packets packet counters for iptable rules.'
print '# TYPE iptables_packets counter'
for line in iptables_packet_lines:
  print line

print '# HELP iptables_bytes byte counters for iptable rules.'
print '# TYPE iptables_bytes counter'
for line in iptables_byte_lines:
  print line
