package chatgpt

import (
	"context"
	"cs-assistant/utils"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type TicketRequest struct {
	Tujuan  string `json:"tujuan"`
	Tanggal string `json:"tanggal"`
	Waktu   string `json:"waktu"`
}

type ChatGPTClient struct {
	client   *openai.Client
	messages []openai.ChatCompletionMessage // Menyimpan riwayat pesan
}

func NewChatGPTClient(ctx context.Context) *ChatGPTClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	return &ChatGPTClient{
		client:   client,
		messages: []openai.ChatCompletionMessage{}, // Inisialisasi riwayat pesan
	}
}

func (c *ChatGPTClient) GetChatGPTResponse(ctx context.Context, user_message string) (string, error) {
	now := time.Now()
	currentDate := now.Format("02-01-2006")
	systemPrompt := strings.ReplaceAll(string(utils.SystemPromptDefault), "{{currentDate}}", currentDate)

	// Menambahkan pesan sistem pertama kali
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	})

	// Menambahkan pesan pengguna
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: user_message,
	})
	fmt.Println("Panjang pesan :", len(c.messages))

	// sudah di coba dan sangat tidak efektif
	// Hanya mengirimkan 3 pesan terakhir, jika lebih dari itu, kita bisa membatasi
	// if len(c.messages) > 7 {
	// 	c.messages = c.messages[len(c.messages)-7:] // Ambil 3 pesan terakhir
	// }

	// Memanggil API ChatGPT untuk mendapatkan respons
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		MaxTokens:   30,
		Messages:    c.messages, // Riwayat pesan untuk sesi percakapan
		Temperature: 1.0,
		TopP:        1.0,
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return "", err
	}

	if len(resp.Choices) > 0 {
		assistantResponse := resp.Choices[0].Message.Content

		// Menambahkan respons asisten ke riwayat percakapan
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: assistantResponse,
		})

		var requestData TicketRequest
		err := json.Unmarshal([]byte(assistantResponse), &requestData)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return assistantResponse, nil
		}

		fmt.Println("Data tujuan  :", requestData.Tujuan)
		fmt.Println("Data tanggal :", requestData.Tanggal)
		fmt.Println("Data waktu   :", requestData.Waktu)

		// proses data dari Database
		return "", nil
	} else {
		fmt.Println("Tidak ada respons yang diterima dari API.")
		return "", nil
	}
}
