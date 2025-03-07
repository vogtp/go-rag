package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/rag"
	"github.com/vogtp/rag/pkg/web/oidc"
)

// Server is the struct holding the webserver
type Server struct {
	slog *slog.Logger

	baseURL string

	httpSrv *http.Server
	mux     *http.ServeMux
	oidcMux *oidc.Mux

	rag        *rag.Manager
	lastEmbedd map[string]time.Time
	docCache   docChace
}

// New creates a new webserver
func New(ctx context.Context, slog *slog.Logger, rag *rag.Manager) (*Server, error) {
	srv := &Server{
		slog:       slog,
		rag:        rag,
		lastEmbedd: make(map[string]time.Time),
		docCache:   newDocCache(),
	}
	srv.httpSrv = &http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	addr := viper.GetString(cfg.WebListen)
	srv.slog = srv.slog.With("listem_addr", addr)
	srv.httpSrv.Addr = addr
	srv.mux = http.NewServeMux()
	oidcCfg := oidc.Config{
		ClientID:     viper.GetString(cfg.OIDCClientID),
		ClientSecret: viper.GetString(cfg.OIDCClientSecret),
		Issuer:       viper.GetString(cfg.OIDCIssuer),
		RedirectURI:  viper.GetString(cfg.OIDCRedirectURI),
	}
	om, err := oidc.NewMux(ctx, srv.slog, srv.mux, addr, oidcCfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create oidc mux: %w", err)
	}
	srv.oidcMux = om
	return srv, nil
}

// OIDC Client
// https://github.com/zitadel/oidc

// Run starts the webserver in foreground
func (srv *Server) Run(ctx context.Context) error {

	if err := srv.schedulePeriodicVecDBUpdates(ctx); err != nil {
		slog.Error("Cannot start periodic embedding", "err", err)
	}

	if err := srv.routes(); err != nil {
		return err
	}

	srv.slog.Warn("Listen for incoming requests")
	srv.httpSrv.Handler = srv.mux
	go srv.closeOnCtxDone(ctx)
	return srv.httpSrv.ListenAndServe()
}

func (srv *Server) closeOnCtxDone(ctx context.Context) {
	<-ctx.Done()
	sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.httpSrv.Shutdown(sdCtx); err != nil {
		slog.Error("cannot shutdown the webserver", "err", err)
	}
}

func (*Server) setStreamHeaders(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")
}
