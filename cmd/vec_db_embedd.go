package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/confluence"
	"github.com/vogtp/rag/pkg/vecDB/filesystem"
)

var vecDbEmbbedCmd = &cobra.Command{
	Use:   "embedd",
	Short: "Embbed to content to a collection",
}

var vecDbEmbbedPathCmd = &cobra.Command{
	Use:   "embedd path <collection> <path>",
	Short: "Embbed to content of a path to a collection",

	Aliases: []string{"path", "p", "dir"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}
		collectionName := args[0]
		path := args[1]
		start := time.Now()
		defer func(t time.Time) {
			fmt.Printf("Updating collection %s took %s\n", collectionName, time.Since(t))
		}(start)
		ctx := cmd.Context()
		client, err := vecdb.New(ctx, slog.Default(), vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
		if err != nil {
			return fmt.Errorf("Failed to create vector DB: %w", err)
		}

		return client.Embedd(ctx, collectionName, filesystem.Generate(ctx, path))
	},
}

var vecDbEmbbedConfluenceCmd = &cobra.Command{
	Use:     "confluence",
	Short:   "Embbed confluence spaces into a collection",
	Aliases: []string{"conf"},
	RunE: func(cmd *cobra.Command, args []string) error {
		slog := slog.Default()
		ctx := cmd.Context()
		collectionName := "confluence"
		c, err := confluence.GetDocuments(ctx, slog)
		if err != nil {
			return err
		}
		client, err := vecdb.New(ctx, slog, vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
		if err != nil {
			return fmt.Errorf("Failed to create vector DB: %w", err)
		}

		return client.Embedd(ctx, collectionName, c)
	},
}
