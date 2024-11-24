package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var Client *storage.Client

func LoadStorageBucket() {
	ctx := context.Background()

	// Decode base64 encoded credentials
	credentials, err := base64.StdEncoding.DecodeString(ENV.GOOGLE_CREDENTIALS_BASE64)
	if err != nil {
		fmt.Printf("Failed to decode credentials: %v\n", err)
		return
	}

	// Create a new storage client
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	Client = client
}
