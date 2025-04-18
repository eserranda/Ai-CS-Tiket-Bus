package model

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIRepository interface {
	GetOpenAiResponse(ctx context.Context, message []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error)
	GetActiveSession(ctx context.Context, phone_number string) ([]openai.ChatCompletionMessage, error)
	SaveChatHistory(ctx context.Context, phone_number string, chatMessages []openai.ChatCompletionMessage) error
}

type OpenAIUsecase interface {
	GetChatGPTResponse(ctx context.Context, phone_number string, message string) (string, error)
}
