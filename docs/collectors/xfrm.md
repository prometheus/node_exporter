# xfrm collector

The xfrm collector exposes metrics about xfrm.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_xfrm_acquire_error_packets_total | State hasn’t been fully acquired before use | n/a |
| node_xfrm_fwd_hdr_error_packets_total | Forward routing of a packet is not allowed | n/a |
| node_xfrm_in_buffer_error_packets_total | No buffer is left | n/a |
| node_xfrm_in_error_packets_total | All errors not matched by other | n/a |
| node_xfrm_in_hdr_error_packets_total | Header error | n/a |
| node_xfrm_in_no_pols_packets_total | No policy is found for states e.g. Inbound SAs are correct but no SP is found | n/a |
| node_xfrm_in_no_states_packets_total | No state is found i.e. Either inbound SPI, address, or IPsec protocol at SA is wrong | n/a |
| node_xfrm_in_pol_block_packets_total | Policy discards | n/a |
| node_xfrm_in_pol_error_packets_total | Policy error | n/a |
| node_xfrm_in_state_expired_packets_total | State is expired | n/a |
| node_xfrm_in_state_invalid_packets_total | State is invalid | n/a |
| node_xfrm_in_state_mismatch_packets_total | State has mismatch option e.g. UDP encapsulation type is mismatch | n/a |
| node_xfrm_in_state_mode_error_packets_total | Transformation mode specific error | n/a |
| node_xfrm_in_state_proto_error_packets_total | Transformation protocol specific error e.g. SA key is wrong | n/a |
| node_xfrm_in_state_seq_error_packets_total | Sequence error i.e. Sequence number is out of window | n/a |
| node_xfrm_in_tmpl_mismatch_packets_total | No matching template for states e.g. Inbound SAs are correct but SP rule is wrong | n/a |
| node_xfrm_out_bundle_check_error_packets_total | Bundle check error | n/a |
| node_xfrm_out_bundle_gen_error_packets_total | Bundle generation error | n/a |
| node_xfrm_out_error_packets_total | All errors which is not matched others | n/a |
| node_xfrm_out_no_states_packets_total | No state is found | n/a |
| node_xfrm_out_pol_block_packets_total | Policy discards | n/a |
| node_xfrm_out_pol_dead_packets_total | Policy is dead | n/a |
| node_xfrm_out_pol_error_packets_total | Policy error | n/a |
| node_xfrm_out_state_expired_packets_total | State is expired | n/a |
| node_xfrm_out_state_invalid_packets_total | State is invalid, perhaps expired | n/a |
| node_xfrm_out_state_mode_error_packets_total | Transformation mode specific error | n/a |
| node_xfrm_out_state_proto_error_packets_total | Transformation protocol specific error | n/a |
| node_xfrm_out_state_seq_error_packets_total | Sequence error i.e. Sequence number overflow | n/a |
