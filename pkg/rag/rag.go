package rag

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

type Manager struct {
	slog *slog.Logger

	models []Model
}

func New(slog *slog.Logger) *Manager {
	m := Manager{
		slog: slog,

		models: []Model{{Name: "Dummy 1", LLMName: "llama3.1"}, {Name: "Dummy 2", LLMName: "llama3.1"}},
	}
	return &m
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
