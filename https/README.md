# HTTPS Package for prometheus

The `https` directory contains files and a template config for the implementation of tls.
When running a server with tls use the flag `--web.tls-config`
Where the path is from where the exporter was run.

e.g. `./node_exporter --web.tls-config="https/tls-config.yml"`
If the config is kept within the https directory 

The config file should is written in YAML format.
The layout is outlined below, with optional parameters in brackets.

### TLS Config Layout

```
#TLS CONFIG YAML
  # Main config options for tls

tlsConfig :

  # Paths to Cert File & Key file from base directory
  # Both required for valid tls
  # Paths set as string values
  # These are reloaded on initial connection and
  tlsCertPath : <filename>
  tlsKeyPath : <filename>

  # ClientAuth declares the policy the server will follow for client auth
  # Accepts the following string values and maps to ClientAuth Policies
  # NoClientCert                
  # RequestClientCert           
  # RequireAnyClientCert        
  # VerifyClientCertIfGiven     
  # RequireAndVerifyClientCert  
  [ clientAuth : <string> | default = "NoClientCert" ]

  # ClientCa's accepts a string path to the set of CA's
  [ clientCAs : <filename> ]
  
```
