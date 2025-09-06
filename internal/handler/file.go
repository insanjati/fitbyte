package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/insanjati/fitbyte/internal/storage"
)

type FileHandler struct {
	storage *storage.MinIOStorage
}

func NewFileHandler(minioStorage *storage.MinIOStorage) *FileHandler {
	return &FileHandler{
		storage: minioStorage,
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Validate file size (max 10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Validate file type (jpeg/jpg/png)
	contentType := header.Header.Get("Content-Type")
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png"}

	isValidType := false
	for _, allowedType := range allowedTypes {
		if strings.Contains(contentType, allowedType) {
			isValidType = true
			break
		}
	}

	if !isValidType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG and PNG files are allowed"})
		return
	}

	// Upload to MinIO
	uri, err := h.storage.UploadFile(c.Request.Context(), file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uri": uri})
}
