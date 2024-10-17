package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/chroma"
)

func addchroma() {
	testCmd.AddCommand(chromaCmd)
}

func chromaFlags() {
	addFlagOllamaUrl(chromaCmd)
	viper.BindPFlags(chromaCmd.Flags())
}

var chromaCmd = &cobra.Command{
	Use:     "chroma",
	Short:   "Run chroma stuff",
	Aliases: []string{"r"},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return chromaVecDB(cmd.Context())
	},
}

func chromaVecDB(ctx context.Context) error {
	chromeURL := "http://localhost:8000"
	index := "vogtp_test_rag"

	llm, e, err := getEmbedding("gemma")
	if err != nil {
		return err
	}

	store, err := chroma.New(
		chroma.WithChromaURL(chromeURL),
		chroma.WithNameSpace(index),
		chroma.WithEmbedder(e),
	)
	if err != nil {
		return fmt.Errorf("cannot create chroma client: %w", err)
	}

	data := []schema.Document{
		{PageContent: "Tokyo", Metadata: map[string]any{"population": 9.7, "area": 622}},
		{PageContent: "Kyoto", Metadata: map[string]any{"population": 1.46, "area": 828}},
		{PageContent: "Hiroshima", Metadata: map[string]any{"population": 1.2, "area": 905}},
		{PageContent: "Kazuno", Metadata: map[string]any{"population": 0.04, "area": 707}},
		{PageContent: "Nagoya", Metadata: map[string]any{"population": 2.3, "area": 326}},
		{PageContent: "Toyota", Metadata: map[string]any{"population": 0.42, "area": 918}},
		{PageContent: "Fukuoka", Metadata: map[string]any{"population": 1.59, "area": 341}},
		{PageContent: "Paris", Metadata: map[string]any{"population": 11, "area": 105}},
		{PageContent: "London", Metadata: map[string]any{"population": 9.5, "area": 1572}},
		{PageContent: "Santiago", Metadata: map[string]any{"population": 6.9, "area": 641}},
		{PageContent: "Buenos Aires", Metadata: map[string]any{"population": 15.5, "area": 203}},
		{PageContent: "Rio de Janeiro", Metadata: map[string]any{"population": 13.7, "area": 1200}},
		{PageContent: "Sao Paulo", Metadata: map[string]any{"population": 22.6, "area": 1523}},
	}

	ids, err := store.AddDocuments(ctx, data)
	if err != nil {
		return fmt.Errorf("cannot add docs: %w", err)
	}
	slog.Info("Added docs to vecDB", "cnt", len(ids))
	docs, err := store.SimilaritySearch(ctx, "yo", 2,
		vectorstores.WithScoreThreshold(0.5),
	)
	if err != nil {
		return fmt.Errorf("cannot search the docs: %w", err)
	}
	fmt.Printf("found docs: %v\n", docs)

	result, err := chains.Run(
		ctx,
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(store, 5, vectorstores.WithScoreThreshold(0.8)),
		),
		"City with a population of more than 5",
	)
	fmt.Println(result)
	return nil
}

func getEmbedding(model string) (llms.Model, *embeddings.EmbedderImpl, error) {
	llm, err := getOllamaClient(model)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create ollama client: %w", err)
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create ollama embedder: %w", err)
	}
	return llms.Model(llm), e, nil
}
