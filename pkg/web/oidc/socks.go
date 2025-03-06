package oidc

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"golang.org/x/net/proxy"
)

func getSocksProxy(slog *slog.Logger) http.RoundTripper {
	proxyAddr := viper.GetString(cfg.HTTPProxy)
	if len(proxyAddr) < 1 {
		return http.DefaultTransport
	}
	slog = slog.With("socks.proxy", proxyAddr)
	tr, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		slog.Warn("Cannot connect to socks", "proxyAddr", proxyAddr, "err", err)
		return http.DefaultTransport
	}

	return &http.Transport{
		DialContext: func(_ context.Context, network, address string) (net.Conn, error) {
			slog.Info("SOCKS request", "network", network, "address", address)
			return tr.Dial(network, address)
		},
	}
}
