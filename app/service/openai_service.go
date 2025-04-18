package service

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService() *OpenAIService {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	return &OpenAIService{
		client: client,
	}
}

func (s *OpenAIService) GetClient() *openai.Client {
	return s.client
}
