# REST api server

## Deployment

#### SERVER_ADDR
SERVER_ADDR defines the address the Alexandria API service will listen to.
This is required.

```
SERVER_ADDR=:443
```

#### TLS_CERT / TLS_KEY / TLS_CACERT / MTLS_ENABLED
TLS_CERT and TLS_KEY define the file path where the certificates used to
serve https traffic are stored. These are optional.

When MTLS_ENABLED is set to true, TLS_CACERT defines the path to the CA certifcate
of the client used to enable mutual TLS.

For the local dev environment, these are set to
```
TLS_CERT="/etc/certs/arcade.pem"
TLS_KEY="/etc/certs/arcade_key.pem"
TLS_CACERT="/etc/certs/rootCA.pem"
MTS_ENABLED="true"
```


#### ALLOWED_ORIGINS
ALLOWED_ORIGINS sets the allowed origins for CORS requests, as used in the
'Allow-Access-Control-Origin' HTTP header. Note: Passing in a string "*"
will allow any domain.

#### ALLOWED_METHODS
ALLOWED_METHODS can be used to explicitly allow methods in the
Access-Control-Allow-Methods header. This is a replacement operation so you
must also pass GET, HEAD, and POST if you wish to support those methods.

#### ALLOWED_HEADERS
ALLOWED_HEADERS adds the provided headers to the list of allowed headers in a
CORS request. This is an append operation so the headers Accept,
Accept-Language, and Content-Language are always allowed. Content-Type must be
explicitly declared if accepting Content-Types other than
application/x-www-form-urlencoded, multipart/form-data, or
text/plain.

#### PPROF_ENABLED
PPROF_ENABLED adds the pprof endpoints to the server.
See: https://pkg.go.dev/net/http/pprof and https://go.dev/blog/pprof

```
PPROF_ENABLED="true"
```

#### READ_TIMEOUT
READ_TIMEOUT sets the server's read timeout. The default value is 5 seconds.
```
READ_TIMEOUT=5s
```

#### WRITE_TIMEOUT
WRITE_TIMEOUT sets the server's write timeout. The default value is 10 seconds.
```
WRITE_TIMEOUT=15s
```

#### SHUTDOWN_TIMEOUT
SHUTDOWN_TIMEOUT sets the server's shutdown timeout. The defailt value is 10 seconds.
```
SHUTDOWN_TIMEOUT=10s
```
