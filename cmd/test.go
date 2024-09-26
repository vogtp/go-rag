package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		return fmt.Errorf("nothing to test")
	},
}
