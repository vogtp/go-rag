package filesystem

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func Generate(ctx context.Context, path string) chan vecdb.EmbeddDocument {
	out := make(chan vecdb.EmbeddDocument, 3)
	go walkPath(ctx, out, path)
	return out
}

func walkPath(ctx context.Context, out chan vecdb.EmbeddDocument, path string) {
	defer close(out)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, _ error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if d.IsDir() {
			return nil
		}

		slog := slog.With("path", path)
		i, err := d.Info()
		if err != nil {
			slog.Warn("cannot get info", "err", err)
			return err
		}
		doc := vecdb.EmbeddDocument{
			IDMetaKey:   vecdb.MetaPath,
			IDMetaValue: path,
			Modified:    i.ModTime(),
		}

		doc.MetaData = make(map[string]any)
		doc.MetaData[vecdb.MetaPath] = path
		doc.MetaData[vecdb.MetaUpdated] = i.ModTime().String()

		slog.Debug("adding document to chroma")
		txt, err := os.ReadFile(path)
		if err != nil {
			slog.Warn("cannot read document", "err", err)
			return err
		}
		doc.Document = string(txt)
		out <- doc
		return nil
	})
	if err != nil {
		slog.Error("Error walking path", "err", err, "path", path)
	}
}
