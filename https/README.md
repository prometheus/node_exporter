# HTTPS Package for Prometheus

The `https` directory contains a Go package and a sample configuration file for running `node_exporter` with HTTPS instead of HTTP.
When running a server with TLS use the flag `--web.config`

e.g. `./node_exporter --web.config="web-config.yml"`
If the config is kept within the https directory.

The config file should be written in YAML format, and is reloaded on each connection to check for new certificates and/or authentication policy.

##Sample Config:
```
tlsConfig :
  # Certificate and key files for server to use to authenticate to client
  tlsCertPath : <filename>
  tlsKeyPath : <filename>

  # Server policy for client authentication. Maps to ClientAuth Policies
  # For more detail on clientAuth options: [ClientAuthType](https://golang.org/pkg/crypto/tls/#ClientAuthType)
  [ clientAuth : <string> | default = "NoClientCert" ]

  # CA certificate for client certificate authentication to the server
  [ clientCAs : <filename> ]
```
