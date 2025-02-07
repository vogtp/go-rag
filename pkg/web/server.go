package web

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/rag"
)

// Server is the struct holding the webserver
type Server struct {
	slog *slog.Logger

	baseURL string

	httpSrv *http.Server
	chi     *chi.Mux

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
	srv.chi = chi.NewRouter()

	return srv
}

// OIDC Client
// https://github.com/zitadel/oidc

// Run starts the webserver in foreground
func (srv *Server) Run(ctx context.Context) error {

	if err := srv.schedulePeriodicVecDBUpdates(ctx); err != nil {
		slog.Error("Cannot start periodic embedding", "err", err)
	}

	logger := httplog.NewLogger("httplog-example", httplog.Options{
		LogLevel: slog.LevelDebug,
		// JSON:             true,
		Concise: true,
		// RequestHeaders:   true,
		// ResponseHeaders:  true,
		MessageFieldName: "message",
		LevelFieldName:   "severity",
		TimeFieldFormat:  time.RFC3339,
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})
	logger.Logger = srv.slog
	srv.chi.Use(httplog.RequestLogger(logger))
	srv.chi.Use(middleware.CleanPath)
	srv.chi.Use(middleware.RequestID)
	srv.chi.Use(middleware.RealIP)
	//srv.chi.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{}))
	srv.chi.Use(middleware.Recoverer)

	srv.routes()

	srv.slog.Warn("Listen for incoming requests")
	srv.httpSrv.Handler = srv.chi
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
