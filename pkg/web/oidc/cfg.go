package oidc

// Config configures the oidc middleware
type Config struct {
	ClientID     string
	ClientSecret string
	Issuer       string
	RedirectURI  string
	LoginPath    string
	Scopes       []string
	ResponseMode string
	KeyPath      string
}
