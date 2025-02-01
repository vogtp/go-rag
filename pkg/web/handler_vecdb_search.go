package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

func (srv Server) vecDBsearch(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	collection := r.PathValue("collection")
	if err := r.ParseForm(); err != nil {
		srv.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	query := r.FormValue("query")
	maxResStr := r.FormValue("maxResults")
	summariseStr := r.FormValue("summary")
	slog := srv.slog.With("collection", collection, "query", query, "maxResults", maxResStr, "summary", summariseStr)
	maxResults, err := strconv.Atoi(maxResStr)
	if err != nil {
		slog.Warn("Cannot convert max Results to int", "err", err)
		maxResults = 10
	}

	summarise, err := strconv.ParseBool(summariseStr)
	if err != nil {
		slog.Debug("Cannot parse summery as bool", "err", err)
		summarise = false
	}
	// if summarise && maxResults > 5 {
	// 	slog.Info("summary does not support too many results")
	// 	maxResults = 5
	// }
	slog = srv.slog.With("collection", collection, "query", query, "maxResults", maxResults, "summary", summarise)
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
		if summarise {
			llm, err := getOllamaClient(ctx, viper.GetString(cfg.ModelDefault))
			if err == nil {
				for i, d := range docs {
					content := []llms.MessageContent{
						llms.TextParts(llms.ChatMessageTypeSystem, "You are a technical analyst who provides short high level summary of the input text.\nDo not use more than 50 words.\nJust provide the short summary with no addional comment."),
						llms.TextParts(llms.ChatMessageTypeHuman, d.Content),
					}
					completion, err := llm.GenerateContent(ctx, content, llms.WithTemperature(0.001))
					if err != nil {
						slog.Warn("Cannot gernerate ollama content", "err", err)
					}
					docs[i].Content = completion.Choices[0].Content
				}
			} else {
				slog.Warn("Cannot connect to ollama", "err", err)

			}
		} 
		data.Documents = docs
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
