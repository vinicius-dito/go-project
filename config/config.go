package config

import (
	"context"
	"os"
)

type Config struct {
	GCPProjectId             string
	GCPProjectLocation       string
	FirestoreUsersCollection string
	BigQueryDataset          string
	BigQueryUsersTable       string
}

func LoadServerConfig(ctx context.Context) Config {
	return Config{
		GCPProjectId:             os.Getenv("GCP_PROJECT_ID"),
		GCPProjectLocation:       os.Getenv("GCP_PROJECT_LOCATION"),
		FirestoreUsersCollection: os.Getenv("FIRESTORE_USERS_COLLECTION"),
		BigQueryDataset:          os.Getenv("BIG_QUERY_DATASET"),
		BigQueryUsersTable:       os.Getenv("BIG_QUERY_USERS_TABLE"),
	}
}
