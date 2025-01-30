package web

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/confluence"
)

func (srv Server) schedulePeriodicVecDBUpdates(ctx context.Context) error {
	updateIntervall := viper.GetDuration(cfg.VecDBUpdateIntervall)
	if updateIntervall < time.Hour {
		slog.Warn("Not starting periodic vector DB updates since update intervall is too short", "updateIntervall", updateIntervall)
		return nil
	}
	ticker := time.NewTicker(updateIntervall)
	go func() {
		if err := srv.embeddConfluence(ctx); err != nil {
			srv.slog.Error("Cannot embedd confluence", "err", err)
		}
		slog.Info("Finished vector DB update")
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}()
	return nil
}

func (srv *Server) embeddConfluence(ctx context.Context) error {
	collectionName := "confluence"
	c, err := confluence.GetDocuments(ctx, srv.slog)
	if err != nil {
		return err
	}
	client, err := vecdb.New(ctx, srv.slog, vecdb.WithOllamaAddress(cfg.GetOllamaHost(ctx)))
	if err != nil {
		return fmt.Errorf("Failed to create vector DB: %w", err)
	}

	if err := client.Embedd(ctx, collectionName, c); err != nil {
		return fmt.Errorf("confluence embedding failed: %w", err)
	}
	srv.lastEmbedd[collectionName] = time.Now()
	return nil
}
