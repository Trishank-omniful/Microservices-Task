package main

import (
	"log"
	"time"

	"github.com/Trishank-Omniful/Onboarding-Task/db"
	"github.com/Trishank-Omniful/Onboarding-Task/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Cannot load env", err)
	}

	db.Connect()

	// Uncomment this line to run migrations
	db.Migrate()

	log.Print("Starting IMS at PORT:8000")

	server := http.InitializeServer(":8000", 10*time.Second, 10*time.Second, 70*time.Second, false)

	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	IMS := server.Engine.Group("/api/v1")
	routes.RegisterHubRoutes(IMS)

	if err := server.StartServer("IMS"); err != nil {
		log.Fatal("Could Not start Server: ", err)
	}

}
