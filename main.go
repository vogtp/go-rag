package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/vogtp/rag/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		select {
		case s := <- signals:
			slog.Warn("Got signal", "sig", s)
			cancel()
		case <-ctx.Done():
		}
	}()

	cobra.CheckErr(cmd.New().ExecuteContext(ctx))
}
