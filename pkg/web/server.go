package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/rag"
)

// Server is the struct holding the webserver
type Server struct {
	slog *slog.Logger

	baseURL string

	httpSrv *http.Server
	mux     *http.ServeMux

	rag        *rag.Manager
	lastEmbedd map[string]time.Time
	docCache   docChace
}

// New creates a new webserver
func New(slog *slog.Logger, rag *rag.Manager) *Server {
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

	return srv
}

// OIDC Client
// https://github.com/zitadel/oidc

// Run starts the webserver in foreground
func (srv *Server) Run(ctx context.Context) error {

	if err := srv.schedulePeriodicVecDBUpdates(ctx); err != nil {
		slog.Error("Cannot start periodic embedding", "err", err)
	}

	srv.routes()

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
