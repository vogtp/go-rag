package rag

import (
	"fmt"
	"log/slog"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Model struct {
	Name    string
	LLMName string

	Collection      string
	OwnedBy         string
	Path            string
	Template        string
	Documents_Count int
}

func (m Model) ToOpenAI() openai.Model {
	return openai.Model{
		// CreatedAt:  0,
		ID:      m.Name,
		Object:  "model",
		OwnedBy: m.OwnedBy,
		Parent:  m.LLMName,
	}
}

func (m Model) GetLLM() (llms.Model, error) {
	llm, err := getOllamaClient(m.LLMName)
	if err != nil {
		return nil, fmt.Errorf("cannot create ollama connection: %w", err)
	}
	return llm, nil
}

func getOllamaClient(model string) (*ollama.LLM, error) {
	url := viper.GetString("url")
	slog.Info("connecting to ollama", "model", model, "url", url)
	return ollama.New(
		ollama.WithModel(model),
		ollama.WithServerURL(url),
	)
}
