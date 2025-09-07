package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOConfig struct {
	Endpoint       string
	AccessKey      string
	SecretKey      string
	BucketName     string
	PublicEndpoint string
	UseSSL         bool
}

type MinIOStorage struct {
	client *minio.Client
	config *MinIOConfig
}

func NewMinIOStorage(config *MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	storage := &MinIOStorage{
		client: client,
		config: config,
	}

	// Create bucket if it doesn't exist
	if err := storage.ensureBucket(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return storage, nil
}

func (s *MinIOStorage) ensureBucket() error {
	ctx := context.Background()
	exists, err := s.client.BucketExists(ctx, s.config.BucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
	}

	return nil
}

func (s *MinIOStorage) UploadFile(ctx context.Context, file io.Reader, header *multipart.FileHeader) (string, error) {
	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	objectName := filepath.Join("uploads", filename)

	// Upload file
	_, err := s.client.PutObject(ctx, s.config.BucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	// Return public URL
	return fmt.Sprintf("%s/%s/%s", s.config.PublicEndpoint, s.config.BucketName, objectName), nil
}
