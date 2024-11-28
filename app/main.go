package main

import (
	"context"
	"cs-assistant/chatgpt"
	"cs-assistant/whatsapp"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	godotenv.Load(".env")
	ctx := context.Background()

	waClient, err := whatsapp.NewWhatsappmeowClient()
	if err != nil {
		return
	}

	// Membuat client ChatGPT
	chatgptClient := chatgpt.NewChatGPTClient(ctx)

	if err := waClient.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to WhatsApp: %v", err)
	}

	// Set event handler untuk WhatsApp
	waClient.SetEventsHandler(ctx, chatgptClient)

	defer func() {
		// menutup koneksi WhatsApp saat program selesai
		if err := waClient.Disconnect(); err != nil {
			log.Printf("Error while disconnecting WhatsApp: %v", err)
		}
	}()

	setupSignalHandler()
}

func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
