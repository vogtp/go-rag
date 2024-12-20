package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
)

func addchroma() {
	testCmd.AddCommand(chromaCmd)
	chromaCmd.AddCommand(chromaColCmd)
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
	Aliases: []string{"c"},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

var chromaColCmd = &cobra.Command{
	Use:     "own <collection_name>",
	Short:   "Run chroma own text example",
	Aliases: []string{"c", "collection", "col", "o"},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		return chromaVecDBOwn(cmd.Context(), args[0])
	},
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
