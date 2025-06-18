package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Trishank-Omniful/Onboarding-Task/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     string
	port     string
	name     string
	username string
	password string
	db       *gorm.DB
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Cannot load env", err)
	}
	host = os.Getenv("POSTGRES_HOST")
	port = os.Getenv("POSTGRES_PORT")
	name = os.Getenv("POSTGRES_DB")
	username = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")

	if host == "" || port == "" || name == "" || username == "" || password == "" {
		fmt.Print("One or more required environment variables are not set", host, port, name, username, password)
	}
}

func Connect() {

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable`, host, username, password, name, port)

	connect, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Print("Cannot Connect to DB: ", err)
		return
	}

	db = connect
	log.Print("Database Connected Successfully")
}

func Migrate() {
	err := db.AutoMigrate(&models.Hub{}, &models.SKU{}, &models.Inventory{})

	if err != nil {
		log.Print("Failed to Auto Migrate: ", err)
	}
	log.Print("Migration Success")
}

func GetDB() *gorm.DB {
	return db
}
