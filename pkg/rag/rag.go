package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
)

type Manager struct {
	slog *slog.Logger

	models []Model
}

func New(ctx context.Context, slog *slog.Logger) *Manager {
	m := Manager{
		slog: slog,

		models: []Model{{Name: "llama3.1", LLMName: viper.GetString(cfg.ModelDefault)}},
	}

	if err := m.updateModelsFromChroma(ctx, slog); err != nil {
		slog.Warn("Cannot get collections from chroma", "err", err)
	}

	return &m
}

func (m *Manager) updateModelsFromChroma(ctx context.Context, slog *slog.Logger) error {
	v, err := vecdb.New(ctx, slog)
	if err != nil {
		return fmt.Errorf("cannot connect to chroma: %w", err)
	}
	collections, err := v.ListCollections(ctx)
	if err != nil {
		return fmt.Errorf("cannot list chroma collections: %w", err)
	}
	model := viper.GetString(cfg.ModelDefault)
	for _, c := range collections {
		m.models = append(m.models, Model{Name: c.Name, Collection: c.Name, LLMName: model})
	}
	return nil
}

func (m Manager) Models() []Model {
	return m.models
}

func (m Manager) Model(name string) (*Model, error) {
	m.slog.Info("Query model", "model", name)
	decoded, err := url.QueryUnescape(name)
	if err != nil {
		decoded = name
	}
	for _, model := range m.models {
		m.slog.Debug("looking for model", "model", decoded, "cur", model.Name)
		if strings.EqualFold(model.Name, decoded) {
			m.slog.Debug("found model", "model", decoded)
			return &model, nil
		}
	}
	return nil, fmt.Errorf("model %s not found", name)
}
