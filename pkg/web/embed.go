package web

import (
	"embed"
	"html/template"
)


var (
	//go:embed templates static ng/intrasearch/dist/intrasearch/browser
	assetData embed.FS
	templates = template.Must(template.ParseFS(assetData, "templates/*.gohtml", "templates/common/*.gohtml"))
)