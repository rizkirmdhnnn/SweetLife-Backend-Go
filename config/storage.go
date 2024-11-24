package config

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var Client *storage.Client

func LoadStorageBucket(pathServiceAccount string) {
	ctx := context.Background()
	// Create a new storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(pathServiceAccount))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	Client = client
}
