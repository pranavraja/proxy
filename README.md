# Prerequisites

- [Go 1.7](https://golang.org/doc/install)

# Setup

    go get

Generate a certificate and private key, in PEM-encoded format, for this proxy
to use. You will end up with a cert.crt and cert.key file. Place these two
files in `$GOPATH/src/github.com/pranavraja/proxy`.

# Running

    go run cmd/proxy/main.go

This will start a proxy server on port 8080 which logs response status codes and times with pretty colors.

By default the proxy is set up to intercept all HTTPS traffic (see below).

# Reflections on trusting trust

This proxy is set up to always mitm HTTPS traffic. You will need to get any
devices you use to trust the `cert.crt` file you generated. The easiest way to
do this is to serve the cert.crt file from a static file server on the same
network and to visit the URL in your device's web browser. This should prompt
the OS to download and install the certificate as a root CA.

Alternatively, you can pass the `-nomitm` flag to disable HTTPS MITM functionality.

# Useful flags

    go run cmd/proxy/main.go -addr :8888

Listen on an alternative port (in this case 8888)

    go run cmd/proxy/main.go -dump out.txt -only hostname.org

Dump the full HTTP requests and responses, in wire format, to and from hostname.org into out.txt

    go run cmd/proxy/main.go -mitmonly hostname1.org,hostname2.org

Restrict HTTPS interception to the provided hosts.

    go run cmd/proxy/main.go -nomitm

Don't intercept HTTPS traffic at all.

# Advanced usage

You can provide a preprocessor function for the request and response. This
allows you to do any crazy mapping logic you can think of.

See `cmd/proxy/map.go` for an example of preprocessing requests.


