package vecdb

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	chroma "github.com/amikos-tech/chroma-go"
	ollamaEmbedd "github.com/amikos-tech/chroma-go/pkg/embeddings/ollama"
	"github.com/amikos-tech/chroma-go/types"
	ollamaAPI "github.com/ollama/ollama/api"
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

func New(ctx context.Context, slog *slog.Logger, opts ...Option) (*VecDB, error) {
	chromaPort := "8000"
	v := &VecDB{
		slog:       slog,
		chromaAddr: fmt.Sprintf("http://localhost:%s", chromaPort),
		//embeddingsModel: "nomic-embed-text",
		embeddingsModel: viper.GetString(cfg.ModelEmbedding),
	}
	for _, o := range opts {
		o(v)
	}
	if len(v.ollamaAddr) < 1 {
		for _, o := range viper.GetStringSlice(cfg.OllamaHosts) {
			u, err := url.Parse(o)
			if err != nil {
				slog.Warn("Cannot parse ollama url", "url", o, "err", err)
				continue
			}
			c := ollamaAPI.NewClient(u, http.DefaultClient)
			if err := c.Heartbeat(ctx); err != nil {
				slog.Warn("Cannot connect to ollama", "url", o, "err", err)
			}
			v.ollamaAddr = o
			break
		}
	}
	if len(v.ollamaAddr) < 1 {
		return nil, fmt.Errorf("no ollama address given")
	}
	if len(v.chromaAddr) < 1 {
		return nil, fmt.Errorf("no chroma address given")
	}
	v.slog = slog.With("chroma_addr", v.chromaAddr, "ollama_addr", v.ollamaAddr)

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
	v.slog.Info("Loading embedding function", "embeddingsModel", v.embeddingsModel)
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
