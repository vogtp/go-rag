package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vogtp/go-angular"
)

func (srv *Server) routes() {

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	srv.chi.Use(middleware.Timeout(60 * time.Second))

	srv.chi.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	fsys, err := fs.Sub(assetData, "ng/intrasearch/dist/intrasearch/browser")
	if err != nil {
		panic(err)
	}
	ngFS := angular.FileSystem(fsys)
	srv.chi.Handle("/", http.FileServer(ngFS))
	srv.chi.Handle("/static/", http.StripPrefix(srv.baseURL, http.FileServer(http.FS(assetData))))

	srv.openAiAPI("/api")
	srv.chi.HandleFunc("/vecdb/", srv.vecDBlist)
	srv.chi.HandleFunc("/vecdb/{collection}", srv.vecDBsearch)
	srv.chi.HandleFunc("/summary/{uuid}", srv.handleSummary)
	// srv.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/search/", http.StatusTemporaryRedirect)
	// })
}

func (srv *Server) openAiAPI(apiBasePath string) {
	srv.slog.Info("Registering openAI API", "basePath", apiBasePath)

	srv.chi.Route(apiBasePath, func(r chi.Router) {
		r.HandleFunc(fmt.Sprintf("POST %scompletions", apiBasePath), srv.completionHandler)
		r.HandleFunc(fmt.Sprintf("POST %schat/completions", apiBasePath), srv.chatCompletionHandler)
		r.HandleFunc(fmt.Sprintf("GET %smodels", apiBasePath), srv.modelsHandler)
		r.HandleFunc(fmt.Sprintf("GET %smodels/{model}", apiBasePath), srv.modelsHandler)
	})
}
