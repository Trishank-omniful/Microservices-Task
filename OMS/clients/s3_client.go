package clients

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3ClientInterface interface {
	UploadCSV(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
}

type s3Client struct {
	client     *s3.Client
	bucketName string
}

func NewS3Client(client *s3.Client, bucketName string) *s3Client {
	return &s3Client{
		client: client, bucketName: bucketName,
	}
}

func (c *s3Client) UploadCSV(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	fileName := fmt.Sprintf("%s-%s%s", uuid.New().String(), time.Now().Format("20060102150405"), filepath.Ext(fileHeader.Filename))
	key := fmt.Sprintf("bulk-orders/%s", fileName)

	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String("text/csv"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return key, nil
}
