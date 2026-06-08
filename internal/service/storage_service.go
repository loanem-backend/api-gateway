package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/loanem-backend/api-gateway/pkg/storage"
)

type StorageService interface {
	StoreInstrumentPicture(ctx context.Context, instrumentID int32, file FileInfo) (string, error)
}

type storageService struct {
	storage *storage.S3Client
}

func NewStorageService(sc *storage.S3Client) StorageService {
	return &storageService{
		storage: sc,
	}
}

func (s *storageService) StoreInstrumentPicture(ctx context.Context, instrumentID int32, file FileInfo) (string, error) {
	key, err := constructKeyFromFileName(FileKindInstrument, instrumentID, file)
	if err != nil {
		return "", err
	}

	objectExists, err := s.checkObjectExists(ctx, key)
	if err != nil {
		return "", err
	}

	if objectExists {
		if _, err := s.storage.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.storage.Bucket),
			Key:    aws.String(key),
		}); err != nil {
			return "", fmt.Errorf("delete object: %w", err)
		}
	}

	if _, err := s.storage.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.storage.Bucket),
		Key:         aws.String(key),
		Body:        *file.File,
		ContentType: aws.String(file.Header.Header.Get("Content-Type")),
	}); err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}

	return key, nil
}

func (s *storageService) checkObjectExists(ctx context.Context, key string) (bool, error) {
	if _, err := s.storage.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.storage.Bucket),
		Key:    aws.String(key),
	}); err != nil {
		var notFoundErr *types.NotFound
		var noSuchKeyErr *types.NoSuchKey
		if errors.As(err, &notFoundErr) || errors.As(err, &noSuchKeyErr) {
			return false, nil
		}

		return false, fmt.Errorf("failed heading object: %w", err)
	}

	return true, nil
}
