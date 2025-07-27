package services

import (
	"context"
	"errors"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/rizkirmdhnnn/sweetlife-backend-go/config"
	"github.com/rizkirmdhnnn/sweetlife-backend-go/repositories"
)

// StorageBucketService is a service that provides methods to interact with storage bucket.
type StorageBucketService interface {
	UploadFile(file *multipart.FileHeader, path, fileName string) (string, error)
	DeleteFile(path string) error
}

// storageBucketService is a service that implements StorageBucketService.
type storageBucketService struct {
	storageRepo repositories.StorageBucketRepository
}

// NewStorageBucketService creates a new instance of StorageBucketService.
func NewStorageBucketService(storageRepo repositories.StorageBucketRepository) StorageBucketService {
	return &storageBucketService{
		storageRepo: storageRepo,
	}
}

// UploadFile implements StorageBucketService.
func (s *storageBucketService) UploadFile(file *multipart.FileHeader, path, fileName string) (string, error) {
	// Open file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Upload file
	url, err := s.storageRepo.UploadFile(context.Background(), config.ENV.STORAGE_BUCKET, path+fileName, src)
	if err != nil {
		return "", err
	}
	return url, nil
}

// DeleteFile implements StorageBucketService.
func (s *storageBucketService) DeleteFile(fileURL string) error {
	// Parse the URL to get the file path
	u, err := url.Parse(fileURL)
	if err != nil {
		return errors.New("invalid file URL")
	}

	// Ensure the URL belongs to the bucket
	if !strings.Contains(u.Host, "storage.googleapis.com") {
		return errors.New("file URL is not from Google Cloud Storage")
	}

	// Extract the file path from the URL
	filePath := strings.TrimPrefix(u.Path, "/sweetlife-go-new/") // Replace with your bucket root
	if filePath == "" {
		return errors.New("invalid file path")
	}

	// Delete the file
	if err := s.storageRepo.DeleteFile(context.Background(), config.ENV.STORAGE_BUCKET, filePath); err != nil {
		return err
	}

	return nil
}
