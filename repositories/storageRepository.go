package repositories

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

// thanks copilot for generating documentation

// StorageBucketRepository defines the interface for interacting with storage buckets.
// It provides methods to create buckets, upload and delete files, and manage folders within the buckets.
type StorageBucketRepository interface {
	// CreateBucket creates a new bucket in the repository.
	// Parameters:
	// - ctx: The context for the request.
	// - bucketName: The name of the bucket to create.
	// - projectID: The ID of the project where the bucket will be created.
	// - attr: The attributes of the bucket to be created.
	// Returns:
	// - error: An error if the bucket creation fails, otherwise nil.
	CreateBucket(ctx context.Context, bucketName, projectID string, attr *storage.BucketAttrs) error

	// UploadFile uploads a file to the specified bucket.
	// Parameters:
	// - ctx: The context for the request.
	// - bucketName: The name of the bucket where the file will be uploaded.
	// - objectName: The name of the object to be created in the bucket.
	// - file: The file to be uploaded.
	// Returns:
	// - string: The URL of the uploaded file.
	// - error: An error if the file upload fails, otherwise nil.
	UploadFile(ctx context.Context, bucketName, objectName string, file io.Reader) (string, error)

	// DeleteFile deletes a file from the specified bucket.
	// Parameters:
	// - ctx: The context for the request.
	// - bucketName: The name of the bucket where the file is located.
	// - objectName: The name of the object to be deleted.
	// Returns:
	// - error: An error if the file deletion fails, otherwise nil.
	DeleteFile(ctx context.Context, bucketName, objectName string) error

	// NewFolder creates a new folder in the specified bucket.
	// Parameters:
	// - ctx: The context for the request.
	// - bucketName: The name of the bucket where the folder will be created.
	// - folderName: The name of the folder to be created.
	// Returns:
	// - error: An error if the folder creation fails, otherwise nil.
	NewFolder(ctx context.Context, bucketName, folderName string) error

	// DeleteFolder deletes a folder from the specified bucket.
	// Parameters:
	// - ctx: The context for the request.
	// - bucketName: The name of the bucket where the folder is located.
	// - folderName: The name of the folder to be deleted.
	// Returns:
	// - error: An error if the folder deletion fails, otherwise nil.
	DeleteFolder(ctx context.Context, bucketName, folderName string) error
}

// storageBucketRepository is a repository for interacting with storage buckets.
type storageBucketRepository struct {
	client *storage.Client
}

// NewStorageBucketService creates a new instance of StorageBucketRepository.
func NewStorageBucketService(client *storage.Client) StorageBucketRepository {
	return &storageBucketRepository{
		client: client,
	}
}

// CreateBucket implements StorageBucketRepository.
func (r *storageBucketRepository) CreateBucket(ctx context.Context, bucketName string, projectID string, attr *storage.BucketAttrs) error {
	if err := r.client.Bucket(bucketName).Create(ctx, projectID, attr); err != nil {
		return err
	}
	return nil
}

// DeleteFile implements StorageBucketRepository.
func (r *storageBucketRepository) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	if err := r.client.Bucket(bucketName).Object(objectName).Delete(ctx); err != nil {
		return err
	}
	return nil
}

// DeleteFolder implements StorageBucketRepository.
func (r *storageBucketRepository) DeleteFolder(ctx context.Context, bucketName string, folderName string) error {
	panic("unimplemented")
}

// NewFolder implements StorageBucketRepository.
func (r *storageBucketRepository) NewFolder(ctx context.Context, bucketName string, folderName string) error {
	panic("unimplemented")
}

// UploadFile implements StorageBucketRepository.
func (r *storageBucketRepository) UploadFile(ctx context.Context, bucketName, objectName string, file io.Reader) (string, error) {
	bucket := r.client.Bucket(bucketName)
	writer := bucket.Object(objectName).NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	// URL generation (adjust as needed)
	fileURL := "https://storage.googleapis.com/" + bucketName + "/" + objectName
	return fileURL, nil
}
