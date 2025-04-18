package openai

import (
	"context"
	"cs-assistant/model"
	"cs-assistant/utils"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAIUsecase struct {
	openai_repo model.OpenAIRepository
}

func NewOpenAIUsecase(openAIRepository model.OpenAIRepository) model.OpenAIUsecase {
	return &OpenAIUsecase{
		openai_repo: openAIRepository,
	}
}

func (u *OpenAIUsecase) GetChatGPTResponse(ctx context.Context, phone_number string, user_message string) (string, error) {

	// Ambil sesi aktif berdasarkan whatsapp_id (phone_number)
	chatMessages, err := u.openai_repo.GetActiveSession(ctx, phone_number)
	if err != nil {
		// Jika sesi tidak ditemukan, buat sesi baru
		fmt.Println("Sesi tidak ditemukan, membuat sesi baru")
		chatMessages = []openai.ChatCompletionMessage{}
	}

	if len(chatMessages) == 0 {
		currentDate := time.Now().Format("02-01-2006")
		systemPrompt := strings.ReplaceAll(string(utils.SystemPromptDefault), "{{currentDate}}", currentDate)
		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// Menambahkan pesan pengguna ke riwayat
	chatMessages = append(chatMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: user_message,
	})

	// Memanggil API ChatGPT untuk mendapatkan respons
	response, err := u.openai_repo.GetOpenAiResponse(ctx, chatMessages)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		assistantResponse := response.Choices[0].Message.Content

		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: assistantResponse,
		})

		// menyimpan sesi aktif ke database
		err := u.openai_repo.SaveChatHistory(ctx, phone_number, chatMessages)
		if err != nil {
			fmt.Println("Error saving chat history:", err)
			return "", err
		}

		return assistantResponse, nil
	} else {
		return "", fmt.Errorf("no response from OpenAI")
	}
}
