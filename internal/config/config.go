package config

import (
	"context"
	"fmt"
	"os"

	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/loanem-backend/api-gateway/pkg/storage"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func InitStorageClient() *storage.S3Client {
	cfg, err := s3config.LoadDefaultConfig(context.Background(),
		s3config.WithRegion(os.Getenv("STORAGE_REGION")),
	)
	if err != nil {
		panic(fmt.Errorf("failed setting up storage client: %w", err))
	}

	client := s3.NewFromConfig(cfg)

	return &storage.S3Client{
		Client:        client,
		PresignClient: s3.NewPresignClient(client),
		Bucket:        os.Getenv("STORAGE_BUCKET"),
	}
}
