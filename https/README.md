# HTTPS Package for Prometheus

The `https` directory contains files and a template config for the implementation of TLS.
When running a server with tls use the flag `--web.tls-config`
Where the path is from where the exporter was run.

e.g. `./node_exporter --web.tls-config="https/tls-config.yml"`
If the config is kept within the https directory.

The config file should is written in YAML format, and is reloaded on each connection to check for new certificates and/or authentication policy.