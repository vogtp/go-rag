package rag

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var _ Model = (*OllamaModel)(nil)

type OllamaModel struct {
	Name    string
	LLMName string

	// Collection      string
	OwnedBy string
	// Path            string
	// Template        string
	// Documents_Count int
}

func (m OllamaModel) GetName() string {
	return m.Name
}

func (m OllamaModel) GetLLMName() string {
	return m.LLMName
}

func (m OllamaModel) ToOpenAI() openai.Model {
	return openai.Model{
		// CreatedAt:  0,
		ID:      m.Name,
		Object:  "OllamaModel",
		OwnedBy: m.OwnedBy,
		Parent:  m.LLMName,
	}
}

// Depricated -> implement the calling functions (or even better generate models)
func (m OllamaModel) getOllamaClient() (llms.Model, error) {
	url := viper.GetString("url")
	slog.Info("connecting to ollama", "OllamaModel", m.LLMName, "url", url)
	return ollama.New(
		ollama.WithModel(m.LLMName),
		ollama.WithServerURL(url),
	)
}

func (m OllamaModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, temperature float64, streamingFunc StreamingFunc) (string, error) {
	llm, err := m.getOllamaClient()
	if err != nil {
		return "", fmt.Errorf("cannot get ollama client: %w", err)
	}

	chains.LoadCondenseQuestionGenerator(llm)

	resp, err := llm.GenerateContent(ctx, messages, llms.WithTemperature(temperature), llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		return streamingFunc(ctx, chunk)
	}))

	respString := ""
	if len(resp.Choices) > 0 {
		respString = resp.Choices[0].Content
	}

	return respString, err
}
