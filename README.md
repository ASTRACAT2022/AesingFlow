# AesingFlow

AesingFlow is a Go library for an authenticated, multiplexed application protocol over QUIC v1 and TLS 1.3. This repository currently contains the protocol core and deliberately does not implement a proxy, TUN device, VPN client, or traffic impersonation.

See [architecture](docs/architecture.md), [wire protocol](docs/protocol.md), [security](docs/security.md), and [testing](docs/testing.md).

## SOCKS5 proxy example

The optional TCP proxy commands are documented in [proxy/README.md](proxy/README.md).
They run an authenticated AesingFlow server exit and a local SOCKS5 endpoint on
`127.0.0.1:8010`; they are not a TUN/VPN implementation.

## Quick check

```sh
go test ./...
go test -race ./...
go vet ./...
```
