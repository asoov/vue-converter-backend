package interfaces

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient defines the interface for an OpenAI client.
type OpenAIClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}
