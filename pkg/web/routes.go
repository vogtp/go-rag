package web

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/vogtp/go-angular"
)

func (srv *Server) routes() error {

	fsys, err := fs.Sub(assetData, "ng/intrasearch/dist/intrasearch/browser")
	if err != nil {
		panic(err)
	}
	ngFS := angular.FileSystem(fsys)
	srv.oidcMux.Handle("/", http.FileServer(ngFS))
	srv.mux.Handle("/static/", http.StripPrefix(srv.baseURL, http.FileServer(http.FS(assetData))))

	srv.openAiAPI("/api/")
	srv.mux.HandleFunc("/vecdb/", srv.vecDBlist)
	srv.mux.HandleFunc("/vecdb/{collection}", srv.vecDBsearch)
	srv.mux.HandleFunc("/summary/{uuid}", srv.handleSummary)
	// srv.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/search/", http.StatusTemporaryRedirect)
	// })
	return nil
}

func (srv *Server) openAiAPI(apiBasePath string) {
	srv.slog.Info("Registering openAI API", "basePath", apiBasePath)

	srv.mux.HandleFunc(fmt.Sprintf("POST %scompletions", apiBasePath), srv.completionHandler)
	srv.mux.HandleFunc(fmt.Sprintf("POST %schat/completions", apiBasePath), srv.chatCompletionHandler)
	srv.mux.HandleFunc(fmt.Sprintf("GET %smodels", apiBasePath), srv.modelsHandler)
	srv.mux.HandleFunc(fmt.Sprintf("GET %smodels/{model}", apiBasePath), srv.modelsHandler)

}
