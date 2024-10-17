package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
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
	chromaCmd.AddCommand(chromaDelCmd)
}

const (
	deleteCollectionFlag = "delete"
	chromeURL            = "http://localhost:8000"
	index                = "vogtp_test_rag"
)

func chromaFlags() {
	addFlagOllamaUrl(chromaCmd)
	chromaCmd.PersistentFlags().Bool(deleteCollectionFlag, false, "delete collection")
	if err := viper.BindPFlags(chromaCmd.PersistentFlags()); err != nil {
		slog.Warn("cannot bind chroma flags", "err", err)
	}
}

var chromaCmd = &cobra.Command{
	Use:     "chroma",
	Short:   "Run chroma stuff",
	Aliases: []string{"r"},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return chromaVecDB(cmd)
	},
}

var chromaDelCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete collection",
	Aliases: []string{"del", "d"},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		_, e, err := getEmbedding("mxbai-embed-large")
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
		if err := store.RemoveCollection(); err != nil {
			return fmt.Errorf("cannot remove collection: %w", err)
		}
		return nil
	},
}

func chromaVecDB(cmd *cobra.Command) error {
	ctx := cmd.Context()

	llm, err := getOllamaClient("llama3.1")
	if err != nil {
		return fmt.Errorf("cannot load llama3.1: %w", err)
	}

	_, e, err := getEmbedding("mxbai-embed-large")
	if err != nil {
		return err
	}
	index := uuid.New().String()
	store, err := chroma.New(
		chroma.WithChromaURL(chromeURL),
		chroma.WithNameSpace(index),
		chroma.WithEmbedder(e),
		chroma.WithDistanceFunction(types.COSINE),
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
		{PageContent: "London", Metadata: map[string]any{"population": 19.5, "area": 1572}},
		{PageContent: "Santiago", Metadata: map[string]any{"population": 6.9, "area": 641}},
		{PageContent: "Buenos Aires", Metadata: map[string]any{"population": 15.5, "area": 203}},
		{PageContent: "Rio de Janeiro", Metadata: map[string]any{"population": 13.7, "area": 1200}},
		{PageContent: "Sao Paulo", Metadata: map[string]any{"population": 22.6, "area": 1523}},
	}
	for i, d := range data {
		meta := make([]string, 0)
		for k, v := range d.Metadata {
			meta = append(meta, fmt.Sprintf("%s=%v", k, v))
		}
		d.PageContent = fmt.Sprintf("%s \nMetadata: %s", d.PageContent, strings.Join(meta, ", "))
		data[i] = d
	}
	for _, d := range data {
		fmt.Printf("%+v\n", d)
	}

	ids, err := store.AddDocuments(ctx, data)
	if err != nil {
		slog.Warn("Cannot add docs", "cnt", len(data), "err", err)
		//return fmt.Errorf("cannot add docs: %w", err)
	}
	slog.Info("Added docs to vecDB", "cnt", len(ids), "ids", strings.Join(ids, ","))

	for _, question := range []string{"london population", "london", "City with a population of more than 15"} {
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
