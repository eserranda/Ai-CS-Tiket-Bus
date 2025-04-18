package main

import (
	"context"
	"cs-assistant/openai"
	"cs-assistant/service"
	"cs-assistant/utils"
	"cs-assistant/whatsapp"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func bootstrap() *gorm.DB {
	mysql_chan := make(chan *gorm.DB, 1)

	errgrp := errgroup.Group{}

	errgrp.Go(func() error {
		db, err := utils.ConnectToMariaDB()
		if err != nil {
			return err
		}
		mysql_chan <- db

		// // Migrate database
		// if err := db.AutoMigrate(&chatgpt.ChatSession{}); err != nil {
		// 	return err
		// }

		return nil
	})

	if err := errgrp.Wait(); err != nil {
		panic(err)
	}

	return <-mysql_chan
}

func main() {
	godotenv.Load(".env")
	ctx := context.Background()
	mysql_db := bootstrap()

	wa_client, err := whatsapp.NewWhatsappmeowClient()
	if err != nil {
		log.Fatalf("Failed to create WhatsApp client: %v", err)
		return
	}

	if err := wa_client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to WhatsApp: %v", err)
		return
	}
	defer wa_client.Disconnect()

	// Menginisialisasi service OpenAI dan database connection
	openAIService := service.NewOpenAIService()

	openai_repository := openai.NewOpenAIRepository(mysql_db, openAIService)
	openai_usecase := openai.NewOpenAIUsecase(openai_repository)

	wa_client.SetEventsHandler(ctx, openai_usecase)

	setupSignalHandler()
}

func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
