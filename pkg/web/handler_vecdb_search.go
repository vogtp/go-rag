package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

type queryDoc struct {
	*vecdb.QueryDocument
	UUID uuid.UUID
}

func (srv Server) vecDBsearch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	collection := r.PathValue("collection")
	if err := r.ParseForm(); err != nil {
		srv.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	query := r.FormValue("query")
	maxResStr := r.FormValue("maxResults")
	slog := srv.slog.With("collection", collection, "query", query, "maxResults", maxResStr)
	maxResults, err := strconv.Atoi(maxResStr)
	if err != nil {
		slog.Warn("Cannot convert max Results to int", "err", err)
		maxResults = 10
	}

	slog = srv.slog.With("collection", collection, "query", query, "maxResults", maxResults)
	slog.Info("Collection search requested")
	var data = struct {
		*commonData
		Collection string
		Query      string
		Documents  []queryDoc
	}{
		commonData: srv.common(fmt.Sprintf("Search: %s", collection), r),
		Collection: collection,
		Query:      query,
	}
	if !srv.lastEmbedd[collection].IsZero() {
		data.StatusMessage = fmt.Sprintf("Last %s update: %v", collection, srv.lastEmbedd[collection])
	}
	if len(query) > 0 {
		ctx := r.Context()
		docs, err := searchVecDB(ctx, slog, collection, query, maxResults)
		if err != nil {
			slog.Error("Cannot query vecDB", "err", err)
			srv.Error(w, r, err.Error(), http.StatusInternalServerError)
			return
		}
		cmpFunc := func(a, b vecdb.QueryDocument) bool {
			return a.URL == b.URL
		}
		docs = slices.CompactFunc(docs, cmpFunc)

		data.Documents = make([]queryDoc, len(docs))
		for i, d := range docs {
			qd := queryDoc{
				QueryDocument: &d,
				UUID:          uuid.New(),
			}
			data.Documents[i] = qd
			srv.docChace.add(&qd)
		}
	}
	data.StatusMessage = fmt.Sprintf("Duration %v - %v", time.Since(start).Truncate(time.Second), data.StatusMessage)
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
