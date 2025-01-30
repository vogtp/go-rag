package rag

import (
	"context"

	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
)

type Model interface {
	GetName() string
	GetLLMName() string
	String() string

	GenerateContent(ctx context.Context, messages []llms.MessageContent, temperature float64, streamingFunc StreamingFunc) (string, error)
	ToOpenAI() openai.Model
}

type StreamingFunc func(ctx context.Context, chunk []byte) error
