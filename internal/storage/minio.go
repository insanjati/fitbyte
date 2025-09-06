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

func NewMinIOStorage(config *MinIOConfig) *MinIOStorage {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize MinIO client: %v", err))
	}

	storage := &MinIOStorage{
		client: client,
		config: config,
	}

	// Create bucket if it doesn't exist
	storage.ensureBucket()

	return storage
}

func (s *MinIOStorage) ensureBucket() {
	ctx := context.Background()
	exists, err := s.client.BucketExists(ctx, s.config.BucketName)
	if err != nil {
		panic(fmt.Sprintf("Error checking bucket existence: %v", err))
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			panic(fmt.Sprintf("Error creating bucket: %v", err))
		}
	}
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
