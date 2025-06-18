package main

import (
	"log"
	"time"

	"github.com/Trishank-Omniful/Onboarding-Task/controllers"
	"github.com/Trishank-Omniful/Onboarding-Task/db"
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

	// Uncomment this line to run migrations
	// db.Migrate()

	gormDB := db.GetDB()

	config := &redis.Config{
		Hosts:       []string{"localhost:6379"},
		PoolSize:    50,
		MinIdleConn: 10,
	}

	client := redis.NewClient(config)
	log.Print("Redis Connected at PORT:6379")

	log.Print("Starting IMS at PORT:8000")

	server := http.InitializeServer(":8000", 10*time.Second, 10*time.Second, 70*time.Second, false)

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	hubRepo := repository.NewHubRepository(gormDB, client)
	hubController := controllers.NewHubController(hubRepo)
	IMS := server.Engine.Group("/api/v1")
	routes.RegisterHubRoutes(IMS, hubController)

	if err := server.StartServer("IMS"); err != nil {
		log.Fatal("Could Not start Server: ", err)
	}

}
