package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/rag"
	"github.com/vogtp/rag/pkg/server"
	"github.com/vogtp/rag/pkg/vecDB/confluence"
)

func addTest() {
	testCmd.AddCommand(testRagCmd)
	testCmd.AddCommand(testScaperCmd)
	testCmd.AddCommand(testConfluenceCmd)
	rootCmd.AddCommand(testCmd)
}

func testFlags() {
	testCmd.Flags().BoolP("test", "t", false, "Show version information")
	if err := viper.BindPFlags(testCmd.Flags()); err != nil {
		slog.Warn("Cannot bind test cmd flags", "err", err)
	}
}

var testCmd = &cobra.Command{
	Use:     "test",
	Short:   "Test stuff",
	Aliases: []string{"t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}
var testScaperCmd = &cobra.Command{
	Use:     "scraper",
	Short:   "Test scraper",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		//_, err := rag.WebScrapRag(cmd.Context(), "https://its.unibas.ch
		_, err := rag.WebScrapRag(cmd.Context(), "https://its.unibas.ch/de/anleitungen/e-mail-kalender/mail-zugang/")

		return err
	},
}

var testRagCmd = &cobra.Command{
	Use:     "rag",
	Short:   "Start RAG server",
	Aliases: []string{"r"},
	RunE: func(cmd *cobra.Command, args []string) error {
		slog := slog.Default()
		rag := rag.New(slog)
		api := server.New(slog, rag)
		return api.Run(cmd.Context(), ":4444")
	},
}

var testConfluenceCmd = &cobra.Command{
	Use: "confluence",
	// Short:   "Start RAG server",
	Aliases: []string{"conf"},
	RunE: func(cmd *cobra.Command, args []string) error {
		slog := slog.Default()

		c, err := confluence.Generate(cmd.Context(), slog)
		if err != nil {
			return err
		}
		for doc := range c {
			fmt.Printf("Doc %v Size: %v\n", doc.Title, len(doc.Document))
		}
		return nil
	},
}
