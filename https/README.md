The `https` directory contains files and a template config for the implementation of tls.
When running a server with tls use the flag --web.tls-config=" /path to config/.yml "
Where the path is from where the exporter was run.

i.e. ./node_exporter --web.tls-config="https/tls-config"
If the config is kept within the https directory 

The layout of the config file should be as below:

#TLS CONFIG YAML
  # Paths to Cert File & Key file from base directory
  # Both required for valid tls
  # Paths set as string values
tlsCertPath : ""
tlsKeyPath : ""

  # Main config options for tls
  # Defaults for all options are nil values
tlsConfig :

  # Root CA's should be a string path to the set of root certificate authorities
  # if nil it will use the host's root CA set
  rootCAs : ~

  # Server Name used to verify hostname on returned certs
  # unless Insecure Skip Verify is true
  serverName : ~

  # Client auth declares the policy the server will follow for client auth
  # Accepts the following string values and maps to ClientAuth Policies
  # NoClientCert                -
  # RequestClientCert           -
  # RequireAnyClientCert        -
  # VerifyClientCertIfGiven     -
  # RequireAndVerifyClientCert  -
  clientAuth : ~

  # Client Ca's accepts a string path to the set of CA's
  clientCAs : ~

  # Controls whether a client verifies the servers cert chain and hostname
  # Boolean value - TLS insecure if true so should only be set as true for testing
  insecureSkipVerify : ~
