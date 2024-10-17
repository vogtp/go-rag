package cmd

import (
	"fmt"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"github.com/vogtp/rag/pkg/cfg"
)

func chromaVecDBCol(cmd *cobra.Command, index string) error {
	ctx := cmd.Context()

	model := viper.GetString(cfg.ModelEmbedding)
	llm, err := getOllamaClient(model)
	if err != nil {
		return fmt.Errorf("cannot load embedding model %s: %w", model, err)
	}

	_, e, err := getEmbedding("mxbai-embed-large")
	if err != nil {
		return err
	}
	store, err := chroma.New(
		chroma.WithChromaURL(chromeURL),
		chroma.WithNameSpace(index),
		chroma.WithEmbedder(e),
		chroma.WithDistanceFunction(types.COSINE),
	)
	if err != nil {
		return fmt.Errorf("cannot create chroma client: %w", err)
	}
	questions := []string{
		"get an account",
		"lost my password",
		"cannot connect",
	}
	for _, question := range questions {
		docs, err := store.SimilaritySearch(ctx, question, 3, vectorstores.WithScoreThreshold(0.3))

		if err != nil {
			return fmt.Errorf("cannot search the docs: %w", err)
		}
		fmt.Printf("**************\nQuestion: %s\nDocs: %v\n", question, len(docs))
		for i, d := range docs {
			fmt.Printf("Doc %v score: %v -> %v %v\n", i, d.Score, d.PageContent, d.Metadata)
		}

		result, err := chains.Run(
			ctx,
			chains.NewRetrievalQAFromLLM(
				llm,
				vectorstores.ToRetriever(
					store,
					5,
					//vectorstores.WithNameSpace(index),
					//vectorstores.WithScoreThreshold(0.8),
				),
			),
			question,
			//"City with a population of more than 5",
		)
		if err != nil {
			return fmt.Errorf("cannot run chain: %w", err)
		}
		fmt.Printf("Res: %v\n", result)
	}
	return nil
}
