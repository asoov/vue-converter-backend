package services

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type ChatClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}
