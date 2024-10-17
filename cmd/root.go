package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addRoot() {

}

func rootFlags() {
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
	if err:=viper.BindPFlags(rootCmd.Flags()); err != nil{
		slog.Warn("Cannot bind root flags","err",err)
	}
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
