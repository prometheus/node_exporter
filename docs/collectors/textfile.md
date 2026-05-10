# textfile collector

The textfile collector exposes metrics about textfile.

## Configuration Flags

| Flag | Description | Default |
| --- | --- | --- |
| collector.textfile.directory | Directory to read text files with metrics from, supports glob matching. (repeatable) |  |

## Metrics

| Metric | Description | Labels |
| --- | --- | --- |
| node_textfile_mtime_seconds | Unixtime mtime of textfiles successfully read. | file |
| node_textfile_scrape_error | 1 if there was an error opening or reading a file, 0 otherwise | n/a |
