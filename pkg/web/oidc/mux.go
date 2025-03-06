package oidc

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
)

// NewMux creates a new OIDC authenticated mux
func NewMux(ctx context.Context, slog *slog.Logger, mux *http.ServeMux, addr string, cfg Config) (*Mux, error) {
	om := &Mux{
		slog: slog.With("oidc", "oidc"),
		mux:  mux,
		addr: addr,
		cfg:  cfg,
	}
	redURI, err := url.Parse(cfg.RedirectURI)
	if err != nil {
		return nil, err
	}
	om.callbackPath = redURI.Path
	om.loginPath = cfg.LoginPath
	if len(cfg.LoginPath) < 1 {
		om.loginPath = "/login"
	}
	om.scopes = cfg.Scopes
	if len(om.scopes) < 1 {
		om.scopes = []string{"openid", "profile", "email"}
	}
	om.responseMode = cfg.ResponseMode
	if len(om.responseMode) < 1 {
		om.responseMode = "code token"
	}
	if err := om.init(ctx, slog); err != nil {
		return nil, err
	}
	return om, nil
}

// Mux OIDC authenticated mux
type Mux struct {
	slog *slog.Logger

	mux  *http.ServeMux
	addr string
	cfg  Config

	loginPath    string
	callbackPath string
	scopes []string
	responseMode string

	cookieHandler  *httphelper.CookieHandler
	providerOIDC   rp.RelyingParty
	stateOIDC      func() string
	urlOptionsOIDC []rp.URLParamOpt
}
