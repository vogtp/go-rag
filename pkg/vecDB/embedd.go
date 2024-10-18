package vecdb

import (
	"context"
	"fmt"
	"time"

	"github.com/amikos-tech/chroma-go/types"
)

const (
	// MetaIDKey is the name of the key which identifies the unique value
	MetaIDKey = "IDkey"

	MetaPath    = "path"
	MetaIsRag   = "RAG"
	MetaCreated = "created"
	MetaUpdated = "updated"
	MetaURL     = "URL"
	MetaTitle   = "title"
)

func parseTime(t string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", t)
}

func (v *VecDB) Embedd(ctx context.Context, collectionName string, in <-chan EmbeddDocument) error {
	slog := v.slog.With("collection", collectionName)
	slogBase := slog
	slog.Info("Starting embedding")
	embedFunc, err := v.GetEmbeddingFunc()
	if err != nil {
		return fmt.Errorf("Error creating ollama embedding function: %s \n", err)
	}

	coll, err := v.CreateCollection(ctx, collectionName, map[string]interface{}{MetaIsRag: true, MetaCreated: time.Now().Unix})
	if err != nil {
		return fmt.Errorf("Failed to create collection: %v", err)
	}
	docUpdated := 0

	for d := range in {
		slog = slogBase.With(d.IDMetaKey, d.IDMetaValue)
		res, err := coll.Get(ctx, map[string]interface{}{d.IDMetaKey: d.IDMetaValue}, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("cannot check for existing docs: %w", err)
		}
		skipFile := len(res.Documents) > 0
		for _, m := range res.Metadatas {
			if u, ok := m[MetaUpdated].(string); ok {
				t, err := parseTime(u)
				if err != nil {
					slog.Info("Cannot parse update time", "time", u)
					skipFile = false
				}
				if d.Modified.After(t) {
					skipFile = false
					slog.Info("File was modified")
					break
				}
			} else {
				skipFile = false
				slog.Warn("cannot read meta data as string", "meta", MetaUpdated, "value", m[MetaUpdated])
			}
		}
		if skipFile {
			slog.Debug("document allready exists and not updated")
			return nil
		}

		rs, err := types.NewRecordSet(
			types.WithEmbeddingFunction(embedFunc),
			types.WithIDGenerator(types.NewULIDGenerator()),
		)

		if err != nil {
			slog.Warn("cannot create record set", "err", err)
			return fmt.Errorf("error creating record set: %s \n", err)
		}

		metadata := []types.Option{types.WithDocument(d.Document), types.WithMetadata(d.IDMetaKey, d.IDMetaValue), types.WithMetadata(MetaIDKey, d.IDMetaKey), types.WithMetadata(MetaUpdated, d.Modified.String())}
		if len(d.URL) > 0 {
			metadata = append(metadata, types.WithMetadata(MetaURL, d.URL))
		}
		if len(d.Title) > 0 {
			metadata = append(metadata, types.WithMetadata(MetaTitle, d.Title))
		}
		for k, v := range d.MetaData {
			metadata = append(metadata, types.WithMetadata(k, v))
		}
		rs.WithRecord(metadata...)

		_, err = rs.BuildAndValidate(ctx)
		if err != nil {
			slog.Debug("cannot validate document", "err", err, "rs", rs)
			slog.Warn("document not validated", "err", err)
			continue
			//return fmt.Errorf("error validating record set: %s \n", err)
		}
		// Add the records to the collection
		ids := rs.GetIDs()
		if len(ids) == len(res.Ids) {
			ids = res.Ids
			slog.Debug("Using IDs from existing document")
		}
		_, err = coll.Upsert(ctx, rs.GetEmbeddings(), rs.GetMetadatas(), rs.GetDocuments(), ids)
		if err != nil {
			slog.Warn("cannot add document", "err", err)
			return fmt.Errorf("Error adding documents: %s \n", err)
		}
		docUpdated++
	}

	// Count the number of documents in the collection
	countDocs, qrerr := coll.Count(ctx)
	if qrerr != nil {
		return fmt.Errorf("Error counting documents: %s \n", qrerr)
	}

	slog.Info("Finished embedding", "docsCount", countDocs, "docsUpdates", docUpdated)

	return nil
}
