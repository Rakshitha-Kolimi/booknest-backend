package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global DB variable
var DB *gorm.DB

func Connect() {
	// You can load these from environment variables or .env file
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Example: "postgresql://user:password@localhost:5432/booknest?sslmode=disable"
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		host, user, password, dbName, port,
	)

	// GORM config with pgx driver
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx", // Use pgx instead of default
		DSN:        dsn,
		PreferSimpleProtocol: true, // disables prepared statements
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ failed to get sqlDB: %v", err)
	}

	// Set connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Connected to PostgreSQL successfully!")
	DB = db
}
