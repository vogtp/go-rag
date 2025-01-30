package web

import (
	"fmt"
	"log/slog"
	"net/http"

	chroma "github.com/amikos-tech/chroma-go"

	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func (srv Server) vecDBlist(w http.ResponseWriter, r *http.Request) {
	slog := slog.With("url", r.URL)
	slog.Info("Collection list requested")
	var data = struct {
		*commonData
		Collections []*chroma.Collection
		Path        string
	}{
		commonData: srv.common("List collections", r),
		Path:       r.URL.Path,
	}
	ctx := r.Context()
	client, err := vecdb.New(ctx, slog)
	if err != nil {
		slog.Error("Failed to create vectorDB client", "err", err)
		http.Error(w, fmt.Sprintf("Failed to create vectorDB client: %v", err), http.StatusInternalServerError)
		return
	}
	cols, err := client.ListCollections(ctx)
	if err != nil {
		slog.Error("Error listing vectorDB collections", "err", err)
		http.Error(w, fmt.Sprintf("Error listing vectorDB collections: %v", err), http.StatusInternalServerError)
		return
	}
	// data.Collections = make([]string, len(cols))
	// for i, c := range cols {
	// 	data.Collections[i] = c.Name
	// }
	data.Collections = cols
	srv.render(w, r, "vecdb_list.gohtml", data)
}
