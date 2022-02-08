hsts-proxy
========
Tiny HTTP reverse proxy which adds HSTS header. This proxy adds the
`Strict-Transport-Security` headers which is required for PCI-DSS compliance to
the origin server's response and returns it to the user.

```bash
go build

# Run the proxy
./hsts-proxy

# Test
curl --resolve 'ifconfig.co:80:127.0.0.1' http://ifconfig.co -v
# Strict-Transport-Security header presents in the response
```
