package web

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-angular"
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
	docChace   docChace
}

// New creates a new webserver
func New(slog *slog.Logger, rag *rag.Manager) *Server {
	a := &Server{
		slog:       slog,
		rag:        rag,
		lastEmbedd: make(map[string]time.Time),
		docChace:   newDocCache(),
	}
	a.httpSrv = &http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	return a
}

// Run starts the webserver in foreground
func (srv *Server) Run(ctx context.Context) error {
	addr := viper.GetString(cfg.WebListen)
	srv.slog = srv.slog.With("listem_addr", addr)
	srv.httpSrv.Addr = addr
	srv.mux = http.NewServeMux()
	if err := srv.schedulePeriodicVecDBUpdates(ctx); err != nil {
		slog.Error("Cannot start periodic embedding", "err", err)
	}
	fsys, err := fs.Sub(assetData, "ng/intrasearch/dist/intrasearch/browser")
	if err != nil {
		panic(err)
	}
	ngFS := angular.FileSystem(fsys)
	srv.mux.Handle("/", http.FileServer(ngFS))
	srv.mux.Handle("/static/", http.StripPrefix(srv.baseURL, http.FileServer(http.FS(assetData))))

	srv.openAiAPI("/api")
	srv.mux.HandleFunc("/search/", srv.vecDBlist)
	srv.mux.HandleFunc("/search/{collection}", srv.vecDBsearch)
	srv.mux.HandleFunc("/summary/{uuid}", srv.handleSummary)
	// srv.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/search/", http.StatusTemporaryRedirect)
	// })

	srv.slog.Warn("Listen for incoming requests")
	srv.httpSrv.Handler = srv.mux
	go srv.closeOnCtxDone(ctx)
	return srv.httpSrv.ListenAndServe()
}

func (srv *Server) openAiAPI(apiBasePath string) {
	if !strings.HasSuffix(apiBasePath, "/") {
		apiBasePath = fmt.Sprintf("%s/", apiBasePath)
	}
	if !strings.HasPrefix(apiBasePath, "/") {
		apiBasePath = fmt.Sprintf("/%s", apiBasePath)
	}
	srv.slog.Info("Registering openAI API", "basePath", apiBasePath)
	srv.mux.HandleFunc(fmt.Sprintf("POST %scompletions", apiBasePath), srv.completionHandler)
	srv.mux.HandleFunc(fmt.Sprintf("POST %schat/completions", apiBasePath), srv.chatCompletionHandler)
	srv.mux.HandleFunc(fmt.Sprintf("GET %smodels", apiBasePath), srv.modelsHandler)
	srv.mux.HandleFunc(fmt.Sprintf("GET %smodels/{model}", apiBasePath), srv.modelsHandler)
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
