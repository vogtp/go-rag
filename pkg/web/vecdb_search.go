package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func (srv Server) vecDBsearch(w http.ResponseWriter, r *http.Request) {
	collection := r.PathValue("collection")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query := r.FormValue("query")
	slog := slog.With("collection", collection, "query", query)
	slog.Info("Collection search requested")
	var data = struct {
		*commonData
		Collection string
		Query      string
		Documents  []vecdb.QueryDocument
	}{
		commonData: srv.common(fmt.Sprintf("Search: %s", collection), r),
		Collection: collection,
		Query:      query,
	}
	ctx := r.Context()
	maxResults := 15
	docs, err := searchVecDB(ctx, slog, collection, query, maxResults)
	if err != nil {
		slog.Error("Cannot query vecDB", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Documents = docs
	srv.render(w, r, "vecdb_search.gohtml", data)
}

func searchVecDB(ctx context.Context, slog *slog.Logger, collection string, query string, maxResults int) ([]vecdb.QueryDocument, error) {
	client, err := vecdb.New(ctx, slog, vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
	if err != nil {
		return nil, fmt.Errorf("Failed to create vector DB: %w", err)
	}

	res, err := client.Query(ctx, collection, []string{query}, int32(maxResults))
	if err != nil {
		return nil, fmt.Errorf("Failed to query vector DB: %w", err)
	}
	return res[0].Documents, nil
}
