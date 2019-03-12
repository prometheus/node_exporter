# Https Package for prometheus

The `https` directory contains files and a template config for the implementation of tls.
When running a server with tls use the flag `--web.tls-config`
Where the path is from where the exporter was run.

e.g. `./node_exporter --web.tls-config="https/tls-config.yml"`
If the config is kept within the https directory 

### TLS Config Layout

```
#TLS CONFIG YAML
  # Main config options for tls

tlsConfig :

  # Paths to Cert File & Key file from base directory
  # Both required for valid tls
  # Paths set as string values
  tlsCertPath : ""
  tlsKeyPath : ""

  # RootCA's should be a string path to the set of root certificate authorities
  # if nil it will use the host's root CA set
  rootCAs : ~

  # ServerName used to verify hostname on returned certs
  # unless Insecure Skip Verify is true
  serverName : ~

  # ClientAuth declares the policy the server will follow for client auth
  # Accepts the following string values and maps to ClientAuth Policies
  # NoClientCert                -
  # RequestClientCert           -
  # RequireAnyClientCert        -
  # VerifyClientCertIfGiven     -
  # RequireAndVerifyClientCert  -
  clientAuth : ~

  # ClientCa's accepts a string path to the set of CA's
  clientCAs : ~

  # InsecureSkipVerify controls whether a client verifies the servers cert chain and hostname
  # Boolean value - TLS insecure if true so should only be set as true for testing
  insecureSkipVerify : ~
```
