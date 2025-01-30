package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/scraper"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

var testScaperCmd = &cobra.Command{
	Use:     "scraper",
	Short:   "Test scraper",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return scapper2vecDB(cmd.Context(), args)
	},
}

func scapper2vecDB(ctx context.Context, args []string) error {
	scrap, err := scraper.New(
		scraper.WithBlacklist([]string{
			"ueber-uns", "about-us",
			"aktuelles",
			"servicekatalog",
			"news",
			"shared-elements",
			"event",
		}),
	)
	if err != nil {
		return fmt.Errorf("cannot create scrapper: %w", err)
	}

	docsChannel := make(chan vecdb.EmbeddDocument, 10)

	go func() {
		if err := scrap.Call(ctx, "https://its.unibas.ch/de/anleitungen/", docsChannel); err != nil {
			slog.Error("Cannot scrap", "err", err)
		}
	}()

	client, err := vecdb.New(ctx, slog.Default(), vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
	if err != nil {
		return fmt.Errorf("Failed to create vector DB: %w", err)
	}

	return client.Embedd(ctx, "its-anleitung-scrap", docsChannel)
}
