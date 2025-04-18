package openai

import (
	"context"
	"cs-assistant/model"
	"cs-assistant/service"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type OpenAIRepository struct {
	client *openai.Client
	db     *gorm.DB
}

func NewOpenAIRepository(db *gorm.DB, openAIService *service.OpenAIService) model.OpenAIRepository {
	return &OpenAIRepository{
		client: openAIService.GetClient(), // Mengambil client dari service
		db:     db,
	}
}

func (r *OpenAIRepository) GetActiveSession(ctx context.Context, whatsappID string) ([]openai.ChatCompletionMessage, error) {
	var session model.ChatSession

	result := r.db.WithContext(ctx).
		Where("whatsapp_id = ? AND status = ?", whatsappID, "active").
		Order("created_at DESC").
		First(&session)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("Sesi tidak ditemukan")
			return nil, nil // Jika sesi tidak ditemukan
		}
		return nil, result.Error
	}

	fmt.Println("Sesi ditemukan")

	var chatMessages []openai.ChatCompletionMessage
	if err := json.Unmarshal([]byte(session.Messages), &chatMessages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %v", err)
	}

	// fmt.Println("Riwayat pesan dari database:", chatMessages)

	return chatMessages, nil
}

func (r *OpenAIRepository) SaveChatHistory(ctx context.Context, whatsappID string, messages []openai.ChatCompletionMessage) error {
	// Konversi messages menjadi JSON
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("error marshalling messages: %w", err)
	}

	// Coba update session aktif yang sudah ada
	result := r.db.WithContext(ctx).
		Model(&model.ChatSession{}).
		Where("whatsapp_id = ? AND status = ?", whatsappID, "active").
		Updates(map[string]interface{}{
			"messages":   string(messagesJSON),
			"updated_at": time.Now(),
		})

	// Jika tidak ada record yang diupdate, buat baru
	if result.RowsAffected == 0 {
		fmt.Println("Session tidak ditemukan, membuat sesi baru")
		newSession := model.ChatSession{
			WhatsappID: whatsappID,
			Messages:   string(messagesJSON),
			Status:     "active",
		}

		if err := r.db.WithContext(ctx).Create(&newSession).Error; err != nil {
			return fmt.Errorf("error creating new session: %w", err)
		}
	} else if result.Error != nil {
		return fmt.Errorf("error updating session: %w", result.Error)
	}

	fmt.Println("Riwayat pesan berhasil disimpan ke database")

	return nil
}

func (r *OpenAIRepository) GetOpenAiResponse(ctx context.Context, chatMessages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
	// Memanggil API ChatGPT untuk mendapatkan respons
	resp, err := r.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		MaxTokens:   30,
		Messages:    chatMessages,
		Temperature: 1.0,
		TopP:        1.0,
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return openai.ChatCompletionResponse{}, err
	}

	return resp, nil
}
