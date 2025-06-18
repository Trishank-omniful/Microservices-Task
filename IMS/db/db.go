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
	log.Println("Migrating...")
	err := db.AutoMigrate(&models.Hub{}, &models.SKU{}, &models.Inventory{})

	if err != nil {
		log.Print("Failed to Auto Migrate: ", err)
	}
	log.Print("Migration Success")
}

func GetDB() *gorm.DB {
	return db
}

func DropAll() {
	log.Println("Dropping all tables...")
	err := db.Migrator().DropTable(&models.Inventory{}, &models.SKU{}, &models.Hub{})
	if err != nil {
		log.Println("Failed to drop tables: ", err)
		return
	}
	log.Println("All tables dropped successfully")
}

func Seed() {

	log.Println("Seeding The Database...")

	for i := 1; i <= 100; i++ {

		hub := models.Hub{
			Name:         fmt.Sprintf("test_name_%d", i),
			Address:      fmt.Sprintf("test_address_%d", i),
			City:         fmt.Sprintf("city_%d", i),
			State:        fmt.Sprintf("state_%d", i),
			Country:      "India",
			PostalCode:   fmt.Sprintf("400%03d", i),
			ContactName:  fmt.Sprintf("contact_%d", i),
			ContactEmail: fmt.Sprintf("contact_%d@example.com", i),
		}
		db.Create(&hub)

		sku := models.SKU{
			Code:        fmt.Sprintf("sku_code_%d", i),
			Name:        fmt.Sprintf("sku_name_%d", i),
			Description: fmt.Sprintf("Description for SKU %d", i),
			TenantId:    fmt.Sprintf("tenant_%d", i),
			SellerId:    fmt.Sprintf("seller_%d", i),
			Category:    "General",
			Price:       models.ToNullFloat64(float64(100 + i)),
		}
		db.Create(&sku)

		inventory := models.Inventory{
			HubID:    hub.ID,
			SKUID:    sku.ID,
			Quantity: 10 * i,
		}
		db.Create(&inventory)
	}
	log.Println("Seeded Hubs, SKUs, and Inventory records.")
}
