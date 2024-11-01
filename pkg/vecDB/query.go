package vecdb

import (
	"context"
	"fmt"
)

type QueryResult struct {
	Question  string
	Documents []QueryDocument
}
type QueryDocument struct {
	Content  string
	Modified string
	URL      string
	Title    string
}

func (v *VecDB) Query(ctx context.Context, collection string, queryTexts []string, nResults int32) ([]QueryResult, error) {
	v.slog.Info("Query vecDB", "collection", collection, "queryTexts", queryTexts, "embeddingsModel", v.embeddingsModel)
	col, err := v.GetCollection(ctx, collection)
	if err != nil {
		return nil, fmt.Errorf("cannot get collection %s: %w", collection, err)
	}
	qr, err := col.Query(ctx, queryTexts, 2, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	results := make([]QueryResult, len(queryTexts))
	for idx, question := range queryTexts {
		res := QueryResult{
			Question:  question,
			Documents: make([]QueryDocument, len(qr.Documents[idx])),
		}
		for i := range qr.Documents[idx] {
			doc := QueryDocument{
				Content: qr.Documents[idx][i],
			}
			metaData := qr.Metadatas[idx][i]
			if m, ok := metaData[MetaUpdated].(string); ok {
				doc.Modified = m
			}
			if u, ok := metaData[MetaURL].(string); ok {
				doc.URL = u
			}
			if t, ok := metaData[MetaTitle].(string); ok {
				doc.Title = t
			}

			res.Documents[i] = doc
		}
		results[idx] = res
	}

	return results, nil
}
