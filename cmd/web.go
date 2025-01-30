package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/pkg/rag"
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
		rag, err := rag.New(cmd.Context(), slog)
		if err != nil {
			return fmt.Errorf("cannot start rag backend: %w", err)
		}
		api := web.New(slog, rag)
		return api.Run(cmd.Context(), ":4444")
	},
}
