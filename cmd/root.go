package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addRoot() {

}

func rootFlags() {
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
	viper.BindPFlags(rootCmd.Flags())
}

var rootCmd = &cobra.Command{
	Use:           "ragctl",
	Short:         "RAG commandline",
	// SilenceUsage:  true,
	// SilenceErrors: true,
	// CompletionOptions: cobra.CompletionOptions{
	// 	DisableDefaultCmd: true,
	// },
	Run: func(cmd *cobra.Command, args []string) {
		if version, _ := cmd.Flags().GetBool("version"); version {
			fmt.Printf("Version: %v\n", "development")
			return
		}
		cmd.Print(cmd.UsageString())
	},
}
