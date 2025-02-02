package web

import (
	"net/http"

	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

// Error is a wrapper for srv.Error
func (srv *Server) Error(w http.ResponseWriter, r *http.Request, errStr string, code int) {
	data := struct {
		*commonData
		Error     string
		Code      int
		Documents []vecdb.QueryDocument
	}{
		commonData: srv.common(errStr, r),
		Error:      errStr,
		Code:       code,
	}
	// http.Error(w, errStr, code)
	w.WriteHeader(code)
	srv.render(w, r, "error.gohtml", data)
}
