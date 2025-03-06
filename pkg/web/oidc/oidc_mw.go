package oidc

import (
	"net/http"
)

// Handle makes shure it is authenticated
func (om *Mux) Handle(pattern string, handler http.Handler) {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		_, err := om.GetSession(w, r)
		if err != nil {
			om.slog.Warn("Not authorised", "err", err)
			http.Redirect(w, r,  om.loginPath, http.StatusTemporaryRedirect)
			return
		}
		handler.ServeHTTP(w, r)
	}
	om.mux.HandleFunc(pattern, handleFunc)
}

// HandleFunc makes shure it is authenticated
func (om *Mux) HandleFunc(pattern string, handler http.HandlerFunc) {
	om.Handle(pattern, handler)
}
