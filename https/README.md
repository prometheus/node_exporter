# HTTPS Package for Prometheus

The `https` directory contains a Go package and a sample configuration file for
running `node_exporter` with HTTPS instead of HTTP. We currently support TLS 1.3
and TLS 1.2.

To run a server with TLS, use the flag `--web.config`.

e.g. `./node_exporter --web.config="web-config.yml"`
If the config is kept within the https directory.

The config file should be written in YAML format, and is reloaded on each connection to check for new certificates and/or authentication policy.

## Sample Config

```
tls_config:
  # Certificate and key files for server to use to authenticate to client
  cert_file: <filename>
  key_file: <filename>

  # Server policy for client authentication. Maps to ClientAuth Policies
  # For more detail on clientAuth options: [ClientAuthType](https://golang.org/pkg/crypto/tls/#ClientAuthType)
  [ client_auth_type: <string> | default = "NoClientCert" ]

  # CA certificate for client certificate authentication to the server
  [ client_ca_file: <filename> ]

# List of usernames and hashed passwords that have full access to the web
# server via basic authentication. If empty, no basic authentication is
# required. Passwords are hashed with bcrypt.
basic_auth_users:
  [ <username>: <password> ... ]
```

## About bcrypt

There are several tools out there to generate bcrypt passwords, e.g.
[htpasswd](https://httpd.apache.org/docs/2.4/programs/htpasswd.html):

`htpasswd -nBC 10 "" | tr -d ':\n`

That command will prompt you for a password and output the hashed password,
which will look something like:
`$2y$10$X0h1gDsPszWURQaxFh.zoubFi6DXncSjhoQNJgRrnGs7EsimhC7zG`

The cost (10 in the example) influences the time it takes for computing the
hash. A higher cost will en up slowing down the authentication process.
Depending on the machine, a cost of 10 will take about ~70ms where a cost of
18 can take up to a few seconds. That hash will be computed on every
password-protected request.
