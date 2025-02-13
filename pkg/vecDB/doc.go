package vecdb

import (
	"log/slog"
	"time"

	"github.com/tmc/langchaingo/textsplitter"
)

type EmbeddDocument struct {
	// IDMetaKey is the key for the identifing unique mata data
	IDMetaKey string
	// IDMetaValue is the value for the identifing unique mata data
	IDMetaValue string

	URL string

	Modified time.Time
	Title    string
	Document string

	MetaData map[string]any
}

func (e EmbeddDocument) Split(slog *slog.Logger) []string {
	splitter := textsplitter.NewMarkdownTextSplitter(textsplitter.WithChunkSize(1*1024))
	s, err := splitter.SplitText(e.Document)
	if err != nil {
		// Markdown splitter never throws error
		slog.Warn("cannot split text", "err", err)
	}
	s = append(s, e.Document)
	return s
}
