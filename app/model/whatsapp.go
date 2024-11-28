package model

import (
	"context"
	"cs-assistant/chatgpt"

	"go.mau.fi/whatsmeow/types"
)

type WhatsAppClient interface {
	Connect(ctx context.Context) error
	Disconnect() error
	SetEventsHandler(ctx context.Context, chatgptClient *chatgpt.ChatGPTClient)
	SendMessage(ctx context.Context, infoChat, infoSender types.JID, message string) error
}
