package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/rag"
	"github.com/vogtp/rag/pkg/server"
)

func addTest() {
	rootCmd.AddCommand(testCmd)
}

func testFlags() {
	testCmd.Flags().BoolP("test", "t", false, "Show version information")
	viper.BindPFlags(testCmd.Flags())
}

var testCmd = &cobra.Command{
	Use:     "test",
	Short:   "Test stuff",
	Aliases: []string{"t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		slog := slog.Default()
		rag := rag.New(slog)
		api := server.New(slog, rag)
		return api.Run(cmd.Context(), ":4444")
	},
}
