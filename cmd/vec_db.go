package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/filesystem"
)

func addVecDB() {
	rootCmd.AddCommand(vecDbCmd)
	vecDbCmd.AddCommand(vecDbEmbbedCmd)
	vecDbCmd.AddCommand(vecDbRmCmd)
	vecDbCmd.AddCommand(vecDbLsCmd)
	vecDbCmd.AddCommand(vecDbSearchCmd)
}

var vecDbCmd = &cobra.Command{
	Use:   "wecdb",
	Short: "manage the vector DB",

	Aliases: []string{"v", "vec"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

var vecDbEmbbedCmd = &cobra.Command{
	Use: "embedd <collection> <path> <ollama_url>",
	//Short: "create a modelfile  and build it automatically",

	Aliases: []string{"e", "create"},
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

		client, err := vecdb.New(slog.Default(), vecdb.WithOllamaAddress(getOllamaHost(args)))
		if err != nil {
			return fmt.Errorf("Failed to create vector DB: %w", err)
		}

		return client.Embedd(cmd.Context(), collectionName, filesystem.Generate(cmd.Context(), path))
	},
}

func getOllamaHost(args []string) string {
	if len(args) < 3 {
		return "http://llama-1.its.unibas.ch:11434"
	}
	return args[2]
}

var vecDbSearchCmd = &cobra.Command{
	Use: "search <collection> <path> <ollama_url>",
	//Short: "create a modelfile  and build it automatically",

	Aliases: []string{"s"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}
		collectionName := args[0]
		search := args[1]
		start := time.Now()
		defer func(t time.Time) {
			fmt.Printf("Updating collection %s took %s\n", collectionName, time.Since(t))
		}(start)

		client, err := vecdb.New(slog.Default(), vecdb.WithOllamaAddress(getOllamaHost(args)))
		if err != nil {
			return fmt.Errorf("Failed to create vector DB: %w", err)
		}

		res, err := client.Query(cmd.Context(), collectionName, []string{search}, 5)
		if err != nil {
			return fmt.Errorf("Failed to connect to vector DB: %w", err)
		}
		for _, r := range res[0].Documents {
			fmt.Printf("Docu: %+v\n", r)
		}
		return nil
	},
}

var vecDbRmCmd = &cobra.Command{
	Use: "rm  <collection>",
	//Short: "create a modelfile  and build it automatically",

	Aliases: []string{"del", "remove"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		client, err := vecdb.New(slog.Default())
		if err != nil {
			return fmt.Errorf("Failed to create client: %w", err)
		}
		for _, a := range args {
			if err := client.DeleteCollection(cmd.Context(), a); err != nil {
				return err
			}
		}
		return nil
	},
}

var vecDbLsCmd = &cobra.Command{
	Use: "ls",
	//Short: "create a modelfile  and build it automatically",

	Aliases: []string{"list", "show"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := vecdb.New(slog.Default())
		if err != nil {
			return fmt.Errorf("Failed to create client: %w", err)
		}
		cols, err := client.ListCollections(cmd.Context())
		if err != nil {
			return err
		}
		for _, c := range cols {
			fmt.Printf("%s\n", c.Name)
		}
		return nil
	},
}
