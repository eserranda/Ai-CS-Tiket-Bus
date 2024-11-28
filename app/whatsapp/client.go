package whatsapp

import (
	"context"
	"cs-assistant/model"
	"fmt"
	"os"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsappmeowClient struct {
	client *whatsmeow.Client
}

func NewWhatsappmeowClient() (model.WhatsAppClient, error) {
	logger := waLog.Stdout("Database", "DEBUG", true)

	// Menghubungkan ke database menggunakan sqlstore
	sql, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s.db?_foreign_keys=on", os.Getenv("WHATSAPP_DB_NAME")), logger)
	if err != nil {
		return nil, err
	}

	// Mendapatkan perangkat pertama dari store
	deviceStore, err := sql.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	// Membuat Whatsmeow client
	client := whatsmeow.NewClient(deviceStore, logger)

	return &WhatsappmeowClient{
		client: client,
	}, nil
}

func (w *WhatsappmeowClient) Connect(ctx context.Context) error {
	if w.client.Store.ID == nil {
		qrChan, _ := w.client.GetQRChannel(ctx)
		err := w.client.Connect()
		if err != nil {
			return err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("Scan the QR Code with your WhatsApp app:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := w.client.Connect()
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *WhatsappmeowClient) Disconnect() error {
	w.client.Disconnect()

	return nil
}
