package service

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
)

type FileStorage interface {
	Upload(ctx context.Context, objectKey, contentType string, file multipart.File, size int64) (string, error)
}

type FileService struct {
	db      *sql.DB
	storage FileStorage
}

func NewFileService(db *sql.DB, storage FileStorage) *FileService {
	return &FileService{db: db, storage: storage}
}

func (s *FileService) UploadFile(ctx context.Context, userID uuid.UUID, fh *multipart.FileHeader) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	objectKey := fmt.Sprintf("uploads/%s_%s", userID.String(), fh.Filename)

	key, err := s.storage.Upload(ctx, objectKey, fh.Header.Get("Content-Type"), file, fh.Size)
	if err != nil {
		return "", err
	}

	// Save metadata in DB
	_, err = s.db.ExecContext(ctx, `
        INSERT INTO files (user_id, object_key, mime_type, size)
        VALUES ($1, $2, $3, $4)`,
		userID, key, fh.Header.Get("Content-Type"), fh.Size)
	if err != nil {
		return "", err
	}

	return key, nil
}
