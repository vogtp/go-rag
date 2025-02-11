package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/rag"
	"github.com/vogtp/rag/pkg/vecDB/chroma"
	"github.com/vogtp/rag/pkg/web"
)

func addWeb() {
	rootCmd.AddCommand(webCmd)
	webCmd.AddCommand(webStartCmd)
}

var webCmd = &cobra.Command{
	Use:     "web",
	Short:   "Manage RAG web server",
	Aliases: []string{"w", "rag", "r"},
}

var webStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start RAG web server",
	//Aliases: []string{"w", "rag", "r"},
	RunE: func(cmd *cobra.Command, args []string) error {
		slog := slog.Default()
		ctx := cmd.Context()
		_, err := startChroma(ctx, slog)
		if err != nil {
			return fmt.Errorf("chroma would not start: %w", err)
		}
		
		rag, err := rag.New(ctx, slog)
		if err != nil {
			return fmt.Errorf("cannot start rag backend: %w", err)
		}
		api := web.New(slog, rag)
		return api.Run(cmd.Context())
	},
}

func startChroma(ctx context.Context, slog *slog.Logger) (func(ctx context.Context) error, error) {
	c, err := chroma.NewContainer(slog)
	if err != nil {
		return nil, fmt.Errorf("cannot create chroma container: %w", err)
	}
	return c.EnsureStarted(ctx, viper.GetInt(cfg.ChromaPort))
}
