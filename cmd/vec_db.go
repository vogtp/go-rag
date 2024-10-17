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
	vecDbCmd.AddCommand(vecDbColLsCmd)
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
	Use:   "embedd <collection> <path> <ollama_url>",
	Short: "Embbed to content of a path to a collection",

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
	Use:   "search <collection> <path> <ollama_url>",
	Short: "Search in a collection",

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
	Use:   "rm  <collection>",
	Short: "Delete collections.  Sepatate by space or use all to delete all",

	Aliases: []string{"del", "remove", "delete"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		client, err := vecdb.New(slog.Default())
		if err != nil {
			return fmt.Errorf("Failed to create client: %w", err)
		}
		if args[0] == "all" {
			cols, err := client.ListCollections(cmd.Context())
			if err != nil {
				return err
			}
			for _, c := range cols {
				slog.Info("Deleting collection", "name", c.Name)
				if err := client.DeleteCollection(cmd.Context(), c.Name); err != nil {
					return err
				}
			}
			return nil
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
	Use:   "ls",
	Short: "List all collection",

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

var vecDbColLsCmd = &cobra.Command{
	Use:   "col <collection_name>",
	Short: "List collection documents",

	Aliases: []string{"c", "collection"},
	Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		colName := args[0]
		ctx := cmd.Context()
		client, err := vecdb.New(slog.Default())
		if err != nil {
			return fmt.Errorf("Failed to create client: %w", err)
		}
		col, err := client.GetCollection(ctx, colName)
		if err != nil {
			return err
		}

		res, err := col.GetWithOptions(ctx)
		if err != nil {
			return fmt.Errorf("cannot get collection documents: %w", err)
		}
		for i, d := range res.Metadatas {
			// fmt.Printf("%s %s\n", d[vecdb.MetaPath], d[vecdb.MetaUpdated])
			fmt.Printf("  ID: %v Len: %v Meta: %v\n", res.Ids[i],len(res.Documents[i]), d)
		}
		fmt.Printf("Found %v docs in collection %s\n", len(res.Documents), colName)
		return nil
	},
}
