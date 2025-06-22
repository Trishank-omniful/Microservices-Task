package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Trishank-omniful/Onboarding-Task/config"
	"github.com/Trishank-omniful/Onboarding-Task/constants"
	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/http"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to Load env file")
	}
}

func main() {
	ctx := context.Background()

	s3, err := config.LoadAWSConfig(ctx)

	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	s3Bucket := os.Getenv("S3_BUCKET_NAME")
	if s3Bucket == "" {
		log.Printf("S3 Bucket not set. Switching to default: TEST")
		s3Bucket = "TEST"
	}

	server := http.InitializeServer(
		constants.ServerPort,
		time.Duration(constants.ServerReadTimeout)*time.Second,
		time.Duration(constants.ServerWriteTimeout)*time.Second,
		time.Duration(constants.ServerIdleTimeout)*time.Second,
		false,
	)

	oms := server.Engine.Group("/api/v1/oms")

}
