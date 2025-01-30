package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/vogtp/rag/pkg/rag"
)

// Server is the struct holding the webserver
type Server struct {
	slog *slog.Logger

	httpSrv *http.Server
	mux     *http.ServeMux

	rag *rag.Manager
}

// New creates a new webserver
func New(slog *slog.Logger, rag *rag.Manager) *Server {
	a := &Server{
		slog: slog,
		rag:  rag,
	}
	a.httpSrv = &http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	return a
}

// Run starts the webserver in foreground
func (srv *Server) Run(ctx context.Context, addr string) error {
	srv.slog = srv.slog.With("listem_addr", addr)
	srv.httpSrv.Addr = addr

	srv.mux = http.NewServeMux()

	srv.openAiAPI("/api")

	srv.slog.Warn("Listen for incoming requests")
	srv.httpSrv.Handler = srv.mux
	go srv.closeOnCtxDone(ctx)
	return srv.httpSrv.ListenAndServe()
}

func (srv *Server) openAiAPI(basePath string) {
	if !strings.HasSuffix(basePath, "/") {
		basePath = fmt.Sprintf("%s/", basePath)
	}
	if !strings.HasPrefix(basePath, "/") {
		basePath = fmt.Sprintf("/%s", basePath)
	}
	srv.slog.Info("Registering openAI API","basePath",basePath)
	srv.mux.HandleFunc(fmt.Sprintf("POST %scompletions", basePath), srv.completionHandler)
	srv.mux.HandleFunc(fmt.Sprintf("POST %schat/completions", basePath), srv.chatCompletionHandler)
	srv.mux.HandleFunc(fmt.Sprintf("GET %smodels", basePath), srv.modelsHandler)
	srv.mux.HandleFunc(fmt.Sprintf("GET %smodels/{model}", basePath), srv.modelsHandler)
}

func (srv *Server) closeOnCtxDone(ctx context.Context) {
	<-ctx.Done()
	sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.httpSrv.Shutdown(sdCtx); err != nil {
		slog.Error("cannot shutdown the webserver", "err", err)
	}
}

func (Server) setStreamHeaders(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")
}
