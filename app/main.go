package main

import (
	"bufio"
	"context"
	"cs-assistant/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var initialized bool

func LoadEnv(path ...string) {
	if initialized {
		return
	}
	if len(path) > 0 {
		godotenv.Load(path[0])
		return
	}
	godotenv.Load()
	initialized = true
}

func main() {
	LoadEnv()
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	// Prompt sistem awal
	rawSystemPrompt := utils.SystemPromptDefault

	// Dapatkan waktu saat ini
	now := time.Now()
	currentDate := now.Format("02-01-2006")

	systemPrompt := strings.ReplaceAll(string(rawSystemPrompt), "{{currentDate}}", currentDate)

	// Simpan riwayat percakapan
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("You: ")
		userMessage, _ := reader.ReadString('\n')
		userMessage = strings.TrimSpace(userMessage)

		if strings.ToLower(userMessage) == "exit" || strings.ToLower(userMessage) == "quit" {
			break
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userMessage,
		})

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       openai.GPT4o,
				MaxTokens:   150,
				Messages:    messages,
				Temperature: 1.0,
				TopP:        1.0,
			},
		)

		// Tangani error jika ada
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Tampilkan hasil dari API
		if len(resp.Choices) > 0 {
			assistantResponse := resp.Choices[0].Message.Content
			fmt.Println("\nAi Response :")
			fmt.Println(assistantResponse)

			// Tambahkan respons asisten ke riwayat percakapan
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: assistantResponse,
			})
		} else {
			fmt.Println("Tidak ada respons yang diterima dari API.")
		}
	}
}
