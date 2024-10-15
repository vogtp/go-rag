package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/vogtp/rag/pkg/rag"
)

type API struct {
	slog *slog.Logger

	httpSrv *http.Server

	rag *rag.Manager
}

func New(slog *slog.Logger, rag *rag.Manager) *API {
	a := &API{
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

func (a *API) Run(ctx context.Context, addr string) error {
	a.slog = a.slog.With("address", addr)
	a.httpSrv.Addr = addr

	mux := http.NewServeMux()
	a.httpSrv.Handler = mux
	mux.HandleFunc("POST /completions", a.completionHandler)
	mux.HandleFunc("POST /chat/completions", a.chatCompletionHandler)
	mux.HandleFunc("GET /models", a.modelsHandler)
	mux.HandleFunc("GET /models/{model}", a.modelsHandler)
	go a.closeOnCtxDone(ctx)

	a.slog.Warn("Listen for incoming requests")
	return a.httpSrv.ListenAndServe()
}

func (a *API) closeOnCtxDone(ctx context.Context) {
	<-ctx.Done()
	sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.httpSrv.Shutdown(sdCtx); err != nil {
		slog.Error("cannot shutdown the webserver", "err", err)
	}
}

func (API) setStreamHeaders(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")
}
