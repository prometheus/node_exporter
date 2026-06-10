# conntrack collector

The conntrack collector exposes metrics about conntrack.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_nf_conntrack_entries | Number of currently allocated flow entries for connection tracking. | n/a |
| node_nf_conntrack_entries_limit | Maximum size of connection tracking table. | n/a |
| node_nf_conntrack_stat_drop | Number of packets dropped due to conntrack failure. | n/a |
| node_nf_conntrack_stat_early_drop | Number of dropped conntrack entries to make room for new ones, if maximum table size was reached. | n/a |
| node_nf_conntrack_stat_found | Number of searched entries which were successful. | n/a |
| node_nf_conntrack_stat_ignore | Number of packets seen which are already connected to a conntrack entry. | n/a |
| node_nf_conntrack_stat_insert | Number of entries inserted into the list. | n/a |
| node_nf_conntrack_stat_insert_failed | Number of entries for which list insertion was attempted but failed. | n/a |
| node_nf_conntrack_stat_invalid | Number of packets seen which can not be tracked. | n/a |
| node_nf_conntrack_stat_search_restart | Number of conntrack table lookups which had to be restarted due to hashtable resizes. | n/a |
