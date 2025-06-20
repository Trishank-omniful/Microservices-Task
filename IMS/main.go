package main

import (
	"log"
	"time"

	"github.com/Trishank-Omniful/Onboarding-Task/constants"
	"github.com/Trishank-Omniful/Onboarding-Task/controllers"
	"github.com/Trishank-Omniful/Onboarding-Task/db"
	"github.com/Trishank-Omniful/Onboarding-Task/middleware"
	"github.com/Trishank-Omniful/Onboarding-Task/repository"
	"github.com/Trishank-Omniful/Onboarding-Task/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/redis"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Cannot load env", err)
	}

	db.Connect()

	// Uncomment this line to drop all tables
	db.DropAll()

	// Uncomment this line to run migrations
	db.Migrate()

	// Uncomment this line to clean Dirty migrations
	// db.CleanDirtyMigration()

	// Uncomment this line to run native Down migrations
	// db.NativeMigrationDown()

	// Uncomment this line to run native Up migrations
	// db.NativeMigrationUp()

	// Uncomment this line to Seed the Database with Dummy Data
	db.Seed()

	gormDB := db.GetDB()

	config := &redis.Config{
		Hosts:       []string{constants.RedisHost},
		PoolSize:    constants.RedisPoolSize,
		MinIdleConn: constants.RedisMinIdleConn,
	}

	client := redis.NewClient(config)
	log.Print("Redis Connected at PORT:6379")

	log.Print("Starting IMS at PORT:8000")

	server := http.InitializeServer(
		constants.ServerPort,
		time.Duration(constants.ServerReadTimeout)*time.Second,
		time.Duration(constants.ServerWriteTimeout)*time.Second,
		time.Duration(constants.ServerIdleTimeout)*time.Second,
		false,
	)

	server.Engine.Use(middleware.CORSMiddleware())
	server.Engine.Use(middleware.LoggingMiddleware())
	server.Engine.Use(middleware.ValidationMiddleware())

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok", "service": "IMS"})
	})

	hubRepo := repository.NewHubRepository(gormDB, client)
	hubController := controllers.NewHubController(hubRepo)
	IMS := server.Engine.Group("/api/v1")
	routes.RegisterHubRoutes(IMS, hubController)

	skuRepo := repository.NewSkuRepository(gormDB, client)
	skuController := controllers.NewSkuController(skuRepo)
	routes.RegisterSkuRoutes(IMS, skuController)

	inventoryRepo := repository.NewInventoryRepository(gormDB, client, skuRepo, hubRepo)
	inventoryController := controllers.NewInventoryController(inventoryRepo)
	routes.RegisterInventoryRoutes(IMS, inventoryController)

	if err := server.StartServer("IMS"); err != nil {
		log.Fatal("Could Not start Server: ", err)
	}

}
