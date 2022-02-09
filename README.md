hsts-proxy
========
Tiny HTTP reverse proxy which adds HSTS header. This proxy adds the
`Strict-Transport-Security` headers which is required for PCI-DSS compliance to
the origin server's response and returns it to the user.

```bash
go build
# It uses 80 port directly
sudo setcap 'cap_net_bind_service=+ep' hsts-proxy

# Run the proxy
./hsts-proxy

# Test
curl --resolve 'ifconfig.co:80:127.0.0.1' http://ifconfig.co -v
# Strict-Transport-Security header presents in the response
```

Use [`build-all`] to create binaries of all supported targets.

### How to install
```bash
# Install the binary
sudo curl -Lo /usr/local/bin/hsts-proxy \
  https://github.com/iamport/hsts-proxy/releases/download/v1.0.0/hsts-proxy-linux-amd64

# Register hsts-proxy service
sudo tee /etc/systemd/system/hsts-proxy.service <<'EOF'
[Unit]
Description=Tiny HTTP reverse proxy which adds HSTS header
Wants=network-online.target

[Service]
AmbientCapabilities=CAP_NET_BIND_SERVICE
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
DynamicUser=yes
ExecStart=/usr/local/bin/hsts-proxy

[Install]
WantedBy=multi-user.target
EOF

# Enable and start hsts-proxy service
sudo systemctl enable --now hsts-proxy
```

[`build-all`]: build-all
