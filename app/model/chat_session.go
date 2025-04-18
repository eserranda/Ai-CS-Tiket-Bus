package model

import "time"

type ChatSession struct {
	ID         uint   `gorm:"primaryKey"`
	WhatsappID string `gorm:"not null"`
	Status     string `gorm:"default:'active'"`
	Messages   string `gorm:"type:jsonb"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Migrate the schema (GORM akan membuat tabel jika belum ada)
// db.AutoMigrate(&ChatSession{})
