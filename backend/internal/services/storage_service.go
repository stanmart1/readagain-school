package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"readagain/internal/utils"
)

type StorageService struct {
	uploadDir      string
	imageOptimizer *utils.ImageOptimizer
}

func NewStorageService(uploadDir string) *StorageService {
	return &StorageService{
		uploadDir:      uploadDir,
		imageOptimizer: utils.NewImageOptimizer(800, 1200, 80),
	}
}

func (s *StorageService) UploadFile(file *multipart.FileHeader, folder string, allowedTypes []string, maxSize int64) (string, error) {
	if file.Size > maxSize {
		return "", utils.NewBadRequestError(fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", maxSize))
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.isAllowedType(ext, allowedTypes) {
		return "", utils.NewBadRequestError(fmt.Sprintf("File type %s is not allowed", ext))
	}

	src, err := file.Open()
	if err != nil {
		return "", utils.NewInternalServerError("Failed to open uploaded file", err)
	}
	defer src.Close()

	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	folderPath := filepath.Join(s.uploadDir, folder)

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", utils.NewInternalServerError("Failed to create upload directory", err)
	}

	filePath := filepath.Join(folderPath, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", utils.NewInternalServerError("Failed to create file", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", utils.NewInternalServerError("Failed to save file", err)
	}

	isImage := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
	if isImage {
		if err := s.imageOptimizer.OptimizeImage(filePath); err != nil {
			utils.ErrorLogger.Printf("Failed to optimize image: %v", err)
		}
	}

	relativePath := filepath.Join(folder, filename)
	return relativePath, nil
}

func (s *StorageService) DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	fullPath := filepath.Join(s.uploadDir, filePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(fullPath); err != nil {
		return utils.NewInternalServerError("Failed to delete file", err)
	}

	return nil
}

func (s *StorageService) GetFileURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	return "/uploads/" + filePath
}

func (s *StorageService) isAllowedType(ext string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (s *StorageService) UploadBookCover(file *multipart.FileHeader) (string, error) {
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".webp"}
	maxSize := int64(5 * 1024 * 1024)
	return s.UploadFile(file, "covers", allowedTypes, maxSize)
}

func (s *StorageService) UploadBookFile(file *multipart.FileHeader) (string, error) {
	allowedTypes := []string{".epub", ".pdf"}
	maxSize := int64(50 * 1024 * 1024)
	return s.UploadFile(file, "books", allowedTypes, maxSize)
}
