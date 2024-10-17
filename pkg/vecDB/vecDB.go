package vecdb

import (
	"context"
	"fmt"
	"log/slog"

	chroma "github.com/amikos-tech/chroma-go"
	ollamaEmbedd "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
)

type VecDB struct {
	slog            *slog.Logger
	chromaAddr      string
	chroma          *chroma.Client
	embedFunc       *ollamaEmbedd.OllamaEmbeddingFunction
	ollamaAddr      string
	embeddingsModel string
}

func New(slog *slog.Logger, opts ...Option) (*VecDB, error) {

	v := &VecDB{
		slog:       slog,
		ollamaAddr: "http://localhost:11434/",
		chromaAddr: "http://localhost:8000",
		//embeddingsModel: "nomic-embed-text",
		embeddingsModel: viper.GetString(cfg.ModelEmbedding),
	}
	for _, o := range opts {
		o(v)
	}
	v.slog = slog.With("chroma_addr", v.chromaAddr)

	client, err := chroma.NewClient(v.chromaAddr)
	if err != nil {
		return nil, fmt.Errorf("Failed to create chroma client: %w", err)
	}
	v.slog.Debug("Connected to chroma")
	v.chroma = client
	return v, nil
}

func (v *VecDB) CreateCollection(ctx context.Context, name string, metadata map[string]interface{}) (*chroma.Collection, error) {
	embedFunc, err := v.GetEmbeddingFunc()
	if err != nil {
		return nil, err
	}
	return v.chroma.CreateCollection(ctx, name, nil, true, embedFunc, types.L2)
}
func (v *VecDB) GetCollection(ctx context.Context, name string) (*chroma.Collection, error) {
	embedFunc, err := v.GetEmbeddingFunc()
	if err != nil {
		return nil, err
	}
	return v.chroma.GetCollection(ctx, name, embedFunc)
}

func (v *VecDB) GetEmbeddingFunc() (*ollamaEmbedd.OllamaEmbeddingFunction, error) {
	if v.embedFunc != nil {
		return v.embedFunc, nil
	}
	embedFunc, err := ollamaEmbedd.NewOllamaEmbeddingFunction(
		ollamaEmbedd.WithBaseURL(v.ollamaAddr),
		ollamaEmbedd.WithModel(v.embeddingsModel),
	)
	if err != nil {
		return nil, fmt.Errorf("Error creating ollama embedding function: %s \n", err)
	}
	v.embedFunc = embedFunc
	return embedFunc, nil
}

func (v *VecDB) DeleteCollection(ctx context.Context, collectionName string) error {
	_, err := v.chroma.DeleteCollection(ctx, collectionName)
	if err != nil {
		err = fmt.Errorf("cannot delete collection %s: %w", collectionName, err)
	}
	return err
}

func (v *VecDB) ListCollections(ctx context.Context) ([]*chroma.Collection, error) {
	return v.chroma.ListCollections(ctx)
}
