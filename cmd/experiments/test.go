package experiments

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Add s
func Add(rootCmd *cobra.Command) {
	testCmd.AddCommand(chromaCmd)
	chromaCmd.AddCommand(chromaColCmd)
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
