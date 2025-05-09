package whatsapp

import (
	"context"
	"cs-assistant/model"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/u2takey/go-utils/rand"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func checkMessage(messageTimestamp time.Time) error {
	currentTime := time.Now()

	// Batas Waktu 10 menit
	allowedDuration := 10 * time.Minute
	// allowedDuration := 5 * time.Second

	// Hitung selisih waktu antara waktu sekarang dan waktu pesan
	timeDifference := currentTime.Sub(messageTimestamp)

	// Jika pesan lebih dari 10 menit, jangan diproses
	if timeDifference > allowedDuration {
		return errors.New("pesan kadaluarsa")
	}

	return nil
}

func (w *WhatsappmeowClient) SetEventsHandler(ctx context.Context, openai_uc model.OpenAIUsecase) {
	w.client.AddEventHandler(func(evt interface{}) {
		switch e := evt.(type) {
		case *events.Message:
			// cek waktu pesan
			if err := checkMessage(e.Info.Timestamp); err != nil {
				return
			}

			var message string
			if e.Message.ExtendedTextMessage.GetText() != "" {
				fmt.Println("Mesasage from Whatsapp Web")
				message = e.Message.ExtendedTextMessage.GetText()
			} else if e.Message.GetConversation() != "" {
				fmt.Println("Mesasage from Whatsapp Mobile")
				message = e.Message.GetConversation()
			}

			fmt.Println("New message :", message)

			// cek apakah pesan kosong, jika iya, abaikan
			switch {
			case message == "":
				fmt.Println("Message is empty")
				return
			case e.Info.IsFromMe:
				fmt.Println("Message is from me")
				return
			case e.Info.IsGroup:
				fmt.Println("Message is group")
				return
			case e.Info.Sender.IsEmpty():
				fmt.Println("Sender is empty")
				return
			}

			// tandai pesan sebagai dibaca
			w.MarkMessageAsReadAndTypingStatus(e.Info.ID, e.Info.Chat, e.Info.Sender)
			// Kirim pesan ke ChatGPT
			phone_number := e.Info.Chat.User

			response, err := openai_uc.GetChatGPTResponse(ctx, phone_number, message)
			if err != nil {
				log.Println(err)
				return
			}

			if response == "" {
				fmt.Println("Response is empty")
				return
			}

			w.SendMessage(ctx, e.Info.Chat, e.Info.Sender, response)
		}
	})
}

func (w *WhatsappmeowClient) MarkMessageAsReadAndTypingStatus(msgID types.MessageID, chatJID, senderJID types.JID) {
	// Tentukan timestamp (waktu saat ini)
	timestamp := time.Now()

	// Menandai pesan sebagai dibaca (Centang 2 biru)
	err := w.client.MarkRead([]types.MessageID{msgID}, timestamp, chatJID, senderJID)
	if err != nil {
		// fmt.Println("Error marking message as read:", err)
		return
	}

	// sedang mengetik pesan
	err = w.client.SendChatPresence(senderJID, "composing", "")
	if err != nil {
		log.Println(err)
	}

	// Delay acak antara 1 hingga 3 detik
	randomSleepDuration := time.Duration(rand.Intn(3)+1) * time.Second
	time.Sleep(randomSleepDuration) // Tidur selama waktu acak yang dihasilkan
}
