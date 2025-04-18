package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToMariaDB() (*gorm.DB, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Mendapatkan nilai dari variabel lingkungan
	dbUser := os.Getenv("MARIADB_USER")
	dbPassword := os.Getenv("MARIADB_PASSWORD")
	dbName := os.Getenv("MARIADB_DATABASE")
	dbHost := "mariadb"
	dbPort := "3306"

	// Format DSN untuk GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Membuka koneksi dengan GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Dapatkan koneksi database SQL underlying
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm.DB: %v", err)
	}

	// Set pengaturan connection pool
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test koneksi
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Connected to MariaDB database successfully with GORM!")
	return db, nil
}

// func ConnectToMySQL() (*sql.DB, error) {
// 	godotenv.Load()

// 	// Mendapatkan nilai dari variabel lingkungan
// 	dbUser := os.Getenv("MARIADB_USER")
// 	dbPassword := os.Getenv("MARIADB_PASSWORD")
// 	dbName := os.Getenv("MARIADB_DATABASE")
// 	dbHost := "mariadb"
// 	dbPort := "3306"

// 	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s", dbUser, dbPassword, dbHost, dbPort, dbName)
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check if the connection is successful
// 	err = db.Ping()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Optional: Add connection pool settings
// 	db.SetMaxOpenConns(10)
// 	db.SetMaxIdleConns(5)
// 	db.SetConnMaxLifetime(time.Hour)

// 	log.Println("Connected to mySQL database successfully!")
// 	return db, nil
// }
