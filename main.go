package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/cmd"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	cobra.CheckErr(cmd.New().ExecuteContext(ctx))
}
