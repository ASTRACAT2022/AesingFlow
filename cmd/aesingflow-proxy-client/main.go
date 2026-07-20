// aesingflow-proxy-client runs a local SOCKS5 proxy on macOS or another OS.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ASTRACAT2022/aesingflow/pkg/aesingflow"
	"github.com/ASTRACAT2022/aesingflow/proxy"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:8010", "local SOCKS5 listen address")
	server := flag.String("server", "", "AesingFlow server host:port")
	serverName := flag.String("server-name", "", "TLS certificate name (defaults to server host)")
	caFile := flag.String("ca", "", "server CA certificate in PEM format")
	token := flag.String("token", "", "AesingFlow access token")
	flag.Parse()
	if *server == "" || *caFile == "" || *token == "" {
		fmt.Fprintln(os.Stderr, "-server, -ca, and -token are required")
		os.Exit(2)
	}
	pem, err := os.ReadFile(*caFile)
	if err != nil {
		slog.Error("read CA", "error", err)
		os.Exit(1)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(pem) {
		slog.Error("invalid CA PEM")
		os.Exit(1)
	}
	name := *serverName
	if name == "" {
		name, _, err = net.SplitHostPort(*server)
		if err != nil {
			slog.Error("invalid server address", "error", err)
			os.Exit(2)
		}
	}
	client, err := aesingflow.NewClient(aesingflow.ClientConfig{Address: *server, TLSConfig: &tls.Config{RootCAs: roots, ServerName: name}, Token: *token, ConnectTimeout: 15 * time.Second})
	if err != nil {
		slog.Error("create AesingFlow client", "error", err)
		os.Exit(1)
	}
	service, err := proxy.NewClient(proxy.ClientConfig{ListenAddress: *listen, Client: client})
	if err != nil {
		slog.Error("create SOCKS5 proxy", "error", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err = service.ListenAndServe(ctx); err != nil {
		slog.Error("SOCKS5 proxy stopped", "error", err)
		os.Exit(1)
	}
}
