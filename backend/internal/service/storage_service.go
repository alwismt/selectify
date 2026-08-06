package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"alwis.dev/selectify/internal/logger"
)

type StorageService interface {
	UploadFile(ctx context.Context, file io.Reader, objectKey string, contentType string) error
	DeleteFile(ctx context.Context, objectKey string) error
}

type storageService struct {
	client        *s3.Client
	bucketName    string
	publicBaseURL string
}

func NewStorageService() StorageService {
	endpoint := os.Getenv("API_S3_ENDPOINT")
	bucketName := os.Getenv("API_S3_BUCKET_NAME")
	accessKeyID := os.Getenv("API_S3_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("API_S3_SECRET_ACCESS_KEY")
	publicBaseURL := os.Getenv("API_S3_PUBLIC_BASE_URL")

	if endpoint == "" {
		panic("API_S3_ENDPOINT is required")
	}
	if bucketName == "" {
		panic("API_S3_BUCKET_NAME is required")
	}
	if accessKeyID == "" {
		panic("API_S3_ACCESS_KEY_ID is required")
	}
	if secretAccessKey == "" {
		panic("API_S3_SECRET_ACCESS_KEY is required")
	}
	if publicBaseURL == "" {
		panic("API_S3_PUBLIC_BASE_URL is required")
	}
	awsConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load S3 configuration: %v", err))
	}
	client := s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(strings.TrimRight(endpoint, "/"))
		options.UsePathStyle = true
	})
	return &storageService{
		client:        client,
		bucketName:    bucketName,
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
	}

}

func (s *storageService) UploadFile(ctx context.Context, file io.Reader, objectKey string, contentType string) error {
	if file == nil {
		return fmt.Errorf("file is required")
	}
	if objectKey == "" {
		return fmt.Errorf("object key is required")
	}
	if contentType == "" {
		return fmt.Errorf("content type is required")
	}

	objectKey = strings.TrimLeft(objectKey, "/")
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return logger.Errorf(ctx, err, "failed to upload %q to bucket %q", objectKey, s.bucketName)
	}
	return nil

}

func (s *storageService) DeleteFile(ctx context.Context, objectKey string) error {
	if objectKey == "" {
		return fmt.Errorf("object key is required")
	}

	objectKey = strings.TrimLeft(objectKey, "/")
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return logger.Errorf(ctx, err, "failed to delete %q from bucket %q", objectKey, s.bucketName)
	}

	return nil
}
