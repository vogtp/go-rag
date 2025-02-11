package confluence

import (
	"log/slog"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

func parsePage(slog *slog.Logger, p string) string {
	// retrun html2text.HTML2Text(p)
	markdown, err := htmltomarkdown.ConvertString(p)
	if err != nil {
		slog.Error("cannot encode html to markdown", "err", err)
		return p
	}
	return markdown
}
