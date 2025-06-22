package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func LoadAWSConfig(ctx context.Context) (*s3.Client, error) {
	localStackEndpointStr := os.Getenv("LOCALSTACK_ENDPOINT")
	if localStackEndpointStr == "" {
		localStackEndpointStr = "http://localhost:4566"
		log.Printf("LOCALSTACK_ENDPOINT not set, defaulting to %s", localStackEndpointStr)
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
		log.Printf("AWS_REGION not set, defaulting to %s", awsRegion)
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config for LocalStack: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = &localStackEndpointStr
		o.UsePathStyle = true
	})

	return s3Client, nil
}

func LoadAWSSQSClient(ctx context.Context) (*sqs.Client, error) {
	localStackEndpointStr := os.Getenv("LOCALSTACK_ENDPOINT")
	if localStackEndpointStr == "" {
		localStackEndpointStr = "http://localhost:4566"
		log.Printf("LOCALSTACK_ENDPOINT not set, defaulting to %s", localStackEndpointStr)
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1"
		log.Printf("AWS_REGION not set, defaulting to %s", awsRegion)
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config for LocalStack SQS: %w", err)
	}

	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = &localStackEndpointStr
	})

	return sqsClient, nil
}
