package services

import (
	"context"
	"mime/multipart"

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
func (s *storageBucketService) DeleteFile(path string) error {
	// TODO: delete file from storage
	panic("unimplemented")
}
