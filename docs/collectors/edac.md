# edac collector

The edac collector exposes metrics about edac.

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_edac_correctable_errors_total | Total correctable memory errors. | controller |
| node_edac_csrow_correctable_errors_total | Total correctable memory errors for this csrow. | controller, csrow |
| node_edac_csrow_uncorrectable_errors_total | Total uncorrectable memory errors for this csrow. | controller, csrow |
| node_edac_uncorrectable_errors_total | Total uncorrectable memory errors. | controller |
